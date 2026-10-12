package chart_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func newClock() *clock { return &clock{time.Date(2026, 10, 12, 9, 0, 0, 0, time.UTC)} }

// view builds charts on a fake clock.
func view(w, h int, clk *clock, build func(c *ui.Context)) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: clk.now})
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

// at is the slot of dot x, y.
func at(t *testing.T, img *image.RGBA, x, y int) pixel.Slot {
	t.Helper()
	c := img.RGBAAt(2*x, 2*y)
	for s, pc := range pixel.SpringDay.Colors {
		if pc.R == c.R && pc.G == c.G && pc.B == c.B {
			return pixel.Slot(s)
		}
	}
	t.Fatalf("dot %d,%d is %v, off the palette", x, y, c)
	return 0
}

// count tallies dots of s in a dot box.
func count(t *testing.T, img *image.RGBA, x, y, w, h int, s pixel.Slot) int {
	t.Helper()
	n := 0
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			if at(t, img, i, j) == s {
				n++
			}
		}
	}
	return n
}

// keeps fails on off-palette colors or split dots.
func keeps(t *testing.T, img *image.RGBA) {
	t.Helper()
	b := img.Bounds()
	for y := 0; y+1 < b.Dy(); y += 2 {
		for x := 0; x+1 < b.Dx(); x += 2 {
			at(t, img, x/2, y/2)
			c := img.RGBAAt(x, y)
			if img.RGBAAt(x+1, y) != c || img.RGBAAt(x, y+1) != c || img.RGBAAt(x+1, y+1) != c {
				t.Fatalf("dot at %d,%d is split", x/2, y/2)
			}
		}
	}
}

// panics fails unless fn panics.
func panics(t *testing.T, what string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatalf("%s did not panic", what)
		}
	}()
	fn()
}

// column is the top row and count of s.
func column(t *testing.T, img *image.RGBA, x, h int, s pixel.Slot) (top, n int) {
	t.Helper()
	top = -1
	for y := range h {
		if at(t, img, x, y) == s {
			if top < 0 {
				top = y
			}
			n++
		}
	}
	return top, n
}
