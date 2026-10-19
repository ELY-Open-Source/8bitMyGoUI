package pixel_test

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// dots reads a Dot 2 frame, dot by dot.
type dots struct {
	t   *testing.T
	img *image.RGBA
	pal *pixel.Palette
}

// at is the slot of dot x, y.
func (d dots) at(x, y int) pixel.Slot {
	d.t.Helper()
	c := d.img.RGBAAt(2*x, 2*y)
	for _, o := range [][2]int{{1, 0}, {0, 1}, {1, 1}} {
		if d.img.RGBAAt(2*x+o[0], 2*y+o[1]) != c {
			d.t.Fatalf("dot %d,%d is not one color", x, y)
		}
	}
	for s, pc := range d.pal.Colors {
		if pc.R == c.R && pc.G == c.G && pc.B == c.B && c.A == 255 {
			return pixel.Slot(s)
		}
	}
	d.t.Fatalf("dot %d,%d has %v, outside the palette", x, y, c)
	return 0
}

// count is how many dots of an area hold s.
func (d dots) count(x, y, w, h int, s pixel.Slot) int {
	n := 0
	for j := range h {
		for i := range w {
			if d.at(x+i, y+j) == s {
				n++
			}
		}
	}
	return n
}

// canvas draws on a w×h dot box at the top-left.
func canvas(t *testing.T, w, h int, draw func(g pixel.Canvas)) dots {
	t.Helper()
	pal := pixel.SummerDay
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pal, Dot: 2})
		k := pixel.K(c)
		ui.Box(c).Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) { draw(k.Canvas(p, r)) })
	}, 2*w, 2*h)
	return dots{t: t, img: tt.Image(), pal: pal}
}
