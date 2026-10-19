package files_test

import (
	"bytes"
	"image"
	"log/slog"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

// world is the settings a test view runs under.
type world struct {
	clk           *clock
	reduce, inert bool
}

func newWorld() *world {
	return &world{clk: &clock{time.Date(2026, 9, 2, 10, 30, 0, 0, time.UTC)}}
}

// view builds on a fake clock, two pixels a dot.
func view(w, h int, wd *world, build func(c *ui.Context)) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: wd.clk.now, Inert: wd.inert, Reduce: wd.reduce})
		build(c)
	}, w, h)
	return tt
}

// later moves the clock on and builds a frame.
func later(tt *ui.Tester, wd *world, d time.Duration) {
	wd.clk.t = wd.clk.t.Add(d)
	tt.Frame()
}

func click(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	if err := tt.Click(label); err != nil {
		t.Fatalf("%v in %q", err, tt.Texts())
	}
}

func find(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("no %q in %q", label, tt.Texts())
	}
	return r
}

// none is a color outside the palette.
const none pixel.Slot = 255

// at is the slot at dot x, y, or none.
func at(img *image.RGBA, x, y int) pixel.Slot {
	c := img.RGBAAt(2*x, 2*y)
	for s, want := range pixel.SpringDay.Colors {
		if c.R == want.R && c.G == want.G && c.B == want.B {
			return pixel.Slot(s)
		}
	}
	return none
}

// count is how many dots in a box hold s.
func count(img *image.RGBA, x, y, w, h int, s pixel.Slot) int {
	n := 0
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			if at(img, i, j) == s {
				n++
			}
		}
	}
	return n
}

// dots is a found box in dots.
func dots(r ui.Rect) (x, y, w, h int) {
	return int(r.X / 2), int(r.Y / 2), int(r.W / 2), int(r.H / 2)
}

func mustPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s did not panic", name)
		}
	}()
	fn()
}

// logged runs fn with the debug log captured.
func logged(fn func()) string {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(old)
	fn()
	return buf.String()
}
