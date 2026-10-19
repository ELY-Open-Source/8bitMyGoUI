package overlay_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func newClock() *clock { return &clock{time.Unix(3_000_000, 0)} }

// view builds on a fake clock; inert veils every control.
func view(w, h int, clk *clock, inert bool, build func(c *ui.Context)) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: clk.now, Inert: inert})
		build(c)
	}, w, h)
}

func find(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("no %q in %q", label, tt.Texts())
	}
	return r
}

func click(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	if err := tt.Click(label); err != nil {
		t.Fatal(err)
	}
}

func shows(tt *ui.Tester, label string) bool {
	_, ok := tt.Find(label)
	return ok
}

// slot is the palette slot of pixel x, y.
func slot(t *testing.T, img *image.RGBA, x, y int) pixel.Slot {
	t.Helper()
	c := img.RGBAAt(x, y)
	for s, pc := range pixel.SpringDay.Colors {
		if pc.R == c.R && pc.G == c.G && pc.B == c.B {
			return pixel.Slot(s)
		}
	}
	t.Fatalf("pixel %d,%d is %v, off the palette", x, y, c)
	return 0
}

// keeps fails on off-palette colors or split dots.
func keeps(t *testing.T, img *image.RGBA) {
	t.Helper()
	b := img.Bounds()
	for y := 0; y+1 < b.Dy(); y += 2 {
		for x := 0; x+1 < b.Dx(); x += 2 {
			slot(t, img, x, y)
			c := img.RGBAAt(x, y)
			if img.RGBAAt(x+1, y) != c || img.RGBAAt(x, y+1) != c || img.RGBAAt(x+1, y+1) != c {
				t.Fatalf("dot at %d,%d is split", x/2, y/2)
			}
		}
	}
}

// count tallies dots of slot s in r.
func count(t *testing.T, img *image.RGBA, r ui.Rect, s pixel.Slot) int {
	t.Helper()
	n := 0
	for y := int(r.Y); y < int(r.Y+r.H); y += 2 {
		for x := int(r.X); x < int(r.X+r.W); x += 2 {
			if slot(t, img, x, y) == s {
				n++
			}
		}
	}
	return n
}
