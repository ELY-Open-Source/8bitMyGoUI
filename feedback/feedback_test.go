package feedback_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/feedback"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

// view builds on a fake clock; reduce stops motion.
func view(w, h int, clk *clock, reduce bool, build func(c *ui.Context)) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: clk.now})
		build(c)
	}, w, h)
	if reduce {
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	}
	return tt
}

func newClock() *clock { return &clock{time.Unix(3_000_000, 0)} }

// wait moves the clock on and builds a frame.
func wait(tt *ui.Tester, clk *clock, d time.Duration) {
	clk.t = clk.t.Add(d)
	tt.Frame()
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

// slots are the palette slots of a box's dots.
func slots(t *testing.T, img *image.RGBA, r ui.Rect) []pixel.Slot {
	t.Helper()
	var out []pixel.Slot
	for y := int(r.Y); y < int(r.Y+r.H); y += 2 {
		for x := int(r.X); x < int(r.X+r.W); x += 2 {
			c, found := img.RGBAAt(x, y), false
			for s, pc := range pixel.SpringDay.Colors {
				if pc.R == c.R && pc.G == c.G && pc.B == c.B {
					out, found = append(out, pixel.Slot(s)), true
					break
				}
			}
			if !found {
				t.Fatalf("pixel %d,%d off the palette", x, y)
			}
		}
	}
	return out
}

func count(s []pixel.Slot, want pixel.Slot) int {
	n := 0
	for _, v := range s {
		if v == want {
			n++
		}
	}
	return n
}

func equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// moves reports whether the labeled box changes as time passes.
func moves(t *testing.T, tt *ui.Tester, clk *clock, label string, d time.Duration) bool {
	t.Helper()
	before := slots(t, tt.Image(), find(t, tt, label))
	wait(tt, clk, d)
	return !equal(before, slots(t, tt.Image(), find(t, tt, label)))
}

func TestProgressFillsChunksAndRuns(t *testing.T) {
	clk := newClock()
	v := 0.5
	tt := view(400, 200, clk, true, func(c *ui.Context) {
		feedback.Progress(c, "Load", v, feedback.ProgressOptions{Width: 103, Percent: true})
	})
	// 103 dots hold 20 chunks of 4×4 dots.
	if n := count(slots(t, tt.Image(), find(t, tt, "Load")), pixel.Accent); n != 10*4*4 {
		t.Fatalf("half a bar lit %d dots, want %d", n, 10*4*4)
	}
	v = -1
	tt.Frame()
	if !moves(t, tt, clk, "Load", 200*time.Millisecond) {
		t.Fatal("an unknown progress stood still under Reduce Motion")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("a value past 1 did not panic")
		}
	}()
	v = 1.5
	tt.Frame()
}

func TestLoadersKeepMovingUnderReduceMotion(t *testing.T) {
	clk := newClock()
	tt := view(400, 300, clk, true, func(c *ui.Context) {
		feedback.Spinner(c, "Spin", feedback.SpinnerOptions{Scale: 2})
		feedback.DotsLoader(c, "Dots", feedback.DotsLoaderOptions{})
		feedback.BarsLoader(c, "Bars", feedback.BarsLoaderOptions{})
		feedback.Hourglass(c, "Sand", feedback.HourglassOptions{})
	})
	for _, l := range []string{"Spin", "Dots", "Bars", "Sand"} {
		if !moves(t, tt, clk, l, 300*time.Millisecond) {
			t.Errorf("%s stood still", l)
		}
	}
}

func TestMeterLevels(t *testing.T) {
	for v, want := range map[float64]pixel.Slot{0.3: pixel.Good, 0.7: pixel.Warn, 0.95: pixel.Bad} {
		tt := view(300, 100, newClock(), false, func(c *ui.Context) {
			feedback.Meter(c, "HP", v, feedback.MeterOptions{Warning: 0.6, Critical: 0.9})
		})
		s := slots(t, tt.Image(), find(t, tt, "HP"))
		if count(s, want) == 0 {
			t.Errorf("%v shows no %v", v, want)
		}
	}
}

func TestSkeletonShineStopsUnderReduceMotion(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		clk := newClock()
		tt := view(400, 200, clk, reduce, func(c *ui.Context) {
			feedback.Skeleton(c, feedback.SkeletonOptions{Avatar: true})
		})
		if moved := moves(t, tt, clk, "NOW LOADING", 250*time.Millisecond); moved == reduce {
			t.Errorf("reduce %v: the shine moved %v", reduce, moved)
		}
	}
}

func TestLoadingScreenRotatesTips(t *testing.T) {
	clk := newClock()
	tt := view(600, 300, clk, false, func(c *ui.Context) {
		feedback.LoadingScreen(c, feedback.LoadingScreenOptions{Tips: []string{"Save often.", "Pet the slime."}})
	})
	r := find(t, tt, "NOW LOADING")
	tips := ui.Rect{X: r.X, Y: r.Y + 2*2*pixel.LineHeight + 16, W: r.W, H: 2 * 2 * pixel.LineHeight}
	before := slots(t, tt.Image(), tips)
	wait(tt, clk, 4*time.Second)
	if equal(before, slots(t, tt.Image(), tips)) {
		t.Fatal("the tip did not change after four seconds")
	}
}
