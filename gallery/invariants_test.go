package gallery

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// render draws view in a w×h dot window.
func render(view func(c *ui.Context), w, h int) *image.RGBA {
	return ui.Render(view, 2*w, 2*h, 1)
}

func colors(ps ...*pixel.Palette) map[color.RGBA]bool {
	out := map[color.RGBA]bool{}
	for _, p := range ps {
		for _, c := range p.Colors {
			out[color.RGBA{c.R, c.G, c.B, 255}] = true
		}
	}
	return out
}

// keepsStyle fails on off-palette colors or off-grid dots.
func keepsStyle(t *testing.T, name string, img *image.RGBA, allowed map[color.RGBA]bool) {
	t.Helper()
	b := img.Bounds()
	for y := 0; y < b.Dy(); y += 2 {
		for x := 0; x < b.Dx(); x += 2 {
			c := img.RGBAAt(x, y)
			if !allowed[c] {
				t.Fatalf("%s: pixel %d,%d is %v, off the palette", name, x, y, c)
			}
			if img.RGBAAt(x+1, y) != c || img.RGBAAt(x, y+1) != c || img.RGBAAt(x+1, y+1) != c {
				t.Fatalf("%s: dot at %d,%d is split across the grid", name, x/2, y/2)
			}
		}
	}
}

func shot(t *testing.T, name string, img *image.RGBA) {
	dir := os.Getenv("PIXEL_SHOTS")
	if dir == "" {
		return
	}
	f, err := os.Create(filepath.Join(dir, name+".png"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}

func TestInvariantsEntries(t *testing.T) {
	for _, e := range Entries() {
		for i, pal := range pixel.Palettes {
			allowed := colors(pal)
			if strings.Contains(e.Code(), "pixel.Scope(") {
				allowed = colors(pixel.Palettes...)
			}
			img := render(func(c *ui.Context) {
				pixel.Use(c, pixel.Settings{Palette: pal})
				ui.Column(c).Fill().Padding(pixel.K(c).D(8)).Children(func() { e.Demo(c, func(string) {}) })
			}, 320, 240)
			keepsStyle(t, e.Slug()+" in "+pal.Name, img, allowed)
			if i == 0 {
				shot(t, e.Slug(), img)
			}
		}
	}
}

func TestInvariantsPages(t *testing.T) {
	pages := map[string]string{"boot": "", "home": "/", "category": "/cat/foundations", "entry": "/f/palettes", "lost": "/nowhere", "search": "?pal"}
	for name, path := range pages {
		for i, pal := range pixel.Palettes {
			a := New("/")
			a.palette = i
			switch {
			case strings.HasPrefix(path, "?"):
				a.started, a.query = true, path[1:]
			case path != "":
				a.Go(path)
			}
			allowed := colors(pal)
			if name == "entry" {
				allowed = colors(pixel.Palettes...)
			}
			img := render(a.View, 640, 400)
			keepsStyle(t, name+" in "+pal.Name, img, allowed)
			if i == 0 {
				shot(t, "page-"+name, img)
			}
		}
	}
}
