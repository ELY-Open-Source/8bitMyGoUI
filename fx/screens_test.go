package fx_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/fx"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestScanlinesShadeHalfOfEveryOtherRow(t *testing.T) {
	d := dots(view(40, 20, pixel.SpringDay, newClock(), false, func(c *ui.Context) {
		fx.Scanlines(c, fx.ScanlinesOptions{}, func() { solid(c, 40, 20, pixel.Info) })
	}).Image(), pixel.SpringDay, 40, 20)
	for y := range 20 {
		if want := 20 * (y % 2); count(d, 0, y, 40, 1, pixel.Shade) != want {
			t.Fatalf("row %d holds %d Shade dots, want %d", y, count(d, 0, y, 40, 1, pixel.Shade), want)
		}
	}
}

func TestVignetteDarkensTowardTheEdges(t *testing.T) {
	d := dots(view(80, 48, pixel.SpringDay, newClock(), false, func(c *ui.Context) {
		fx.Vignette(c, fx.VignetteOptions{}, func() { solid(c, 80, 48, pixel.Lit) })
	}).Image(), pixel.SpringDay, 80, 48)
	if d[0][0] != int(pixel.Ink) || d[47][79] != int(pixel.Ink) || count(d, 24, 14, 32, 20, pixel.Ink) != 0 {
		t.Fatal("corners must be solid and the middle clear")
	}
	if edge, in := count(d, 0, 0, 1, 48, pixel.Ink), count(d, 12, 0, 1, 48, pixel.Ink); edge <= in || in == 0 {
		t.Fatalf("the edge column holds %d dark dots, one further in %d", edge, in)
	}
}

// cued builds an effect whose cue the test changes.
func cued(clk *clock, reduce bool, w, h int, build func(c *ui.Context, cue int)) (*ui.Tester, *int) {
	cue := new(int)
	return view(w, h, pixel.SpringDay, clk, reduce, func(c *ui.Context) { build(c, *cue) }), cue
}

func TestScreenShakeJoltsThenSettles(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		clk := newClock()
		tt, cue := cued(clk, reduce, 40, 30, func(c *ui.Context, cue int) {
			fx.ScreenShake(c, fx.ScreenShakeOptions{Cue: cue}, func() { solid(c, 20, 10, pixel.Info) })
		})
		at := func() (float32, float32) {
			r, ok := tt.Find("content")
			if !ok {
				t.Fatal("no content")
			}
			return r.X, r.Y
		}
		x0, y0 := at()
		*cue++
		tt.Frame()
		// The first jolt is three dots right and two up.
		if x, y := at(); !reduce && (x != x0+6 || y != y0-4) || reduce && (x != x0 || y != y0) {
			t.Fatalf("reduce %v: content moved from %v,%v to %v,%v", reduce, x0, y0, x, y)
		}
		clk.t = clk.t.Add(time.Second)
		tt.Frame()
		if x, y := at(); x != x0 || y != y0 {
			t.Fatalf("reduce %v: the shake did not settle", reduce)
		}
	}
}

func TestFlashFillsAStepThenHalfAStep(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		clk := newClock()
		tt, cue := cued(clk, reduce, 40, 20, func(c *ui.Context, cue int) {
			fx.Flash(c, fx.FlashOptions{Cue: cue}, func() { solid(c, 40, 20, pixel.Info) })
		})
		lit := func() int { return count(dots(tt.Image(), pixel.SpringDay, 40, 20), 0, 0, 40, 20, pixel.Lit) }
		if lit() != 0 {
			t.Fatal("flashed before its cue")
		}
		*cue++
		want := []int{800, 400, 0}
		if reduce {
			want = []int{0, 0, 0}
		}
		for _, w := range want {
			tt.Frame()
			if lit() != w {
				t.Fatalf("reduce %v: %d Lit dots, want %d", reduce, lit(), w)
			}
			clk.t = clk.t.Add(100 * time.Millisecond)
		}
	}
}

func TestFadeDitherStepsThroughBayerLevels(t *testing.T) {
	clk, covered := newClock(), true
	tt := view(32, 16, pixel.SpringDay, clk, false, func(c *ui.Context) {
		fx.FadeDither(c, covered, fx.FadeDitherOptions{}, func() { solid(c, 32, 16, pixel.Lit) })
	})
	dark := func() int { return count(dots(tt.Image(), pixel.SpringDay, 32, 16), 0, 0, 32, 16, pixel.Ink) }
	// It starts covered; uncovering drops two levels a step.
	covered = false
	for level := 16; level >= 0; level -= 2 {
		tt.Frame()
		if dark() != 32*level {
			t.Fatalf("level %d: %d dark dots", level, dark())
		}
		clk.t = clk.t.Add(tick(12))
	}
}

func TestIrisWipeClosesOnTheMiddle(t *testing.T) {
	clk, closed := newClock(), false
	tt := view(64, 32, pixel.SpringDay, clk, false, func(c *ui.Context) {
		fx.IrisWipe(c, closed, fx.IrisWipeOptions{}, func() { solid(c, 64, 32, pixel.Lit) })
	})
	closed = true
	last := -1
	for step := range 9 {
		tt.Frame()
		d := dots(tt.Image(), pixel.SpringDay, 64, 32)
		n := count(d, 0, 0, 64, 32, pixel.Ink)
		switch {
		case n <= last && step < 8:
			t.Fatalf("step %d: the iris did not close further", step)
		case count(d, 0, 0, 32, 32, pixel.Ink) != count(d, 32, 0, 32, 32, pixel.Ink):
			t.Fatalf("step %d: the iris is off the middle", step)
		case step < 8 && d[16][32] != int(pixel.Lit), step == 8 && n != 64*32:
			t.Fatalf("step %d: the middle shows %d", step, d[16][32])
		}
		last = n
		clk.t = clk.t.Add(tick(12))
	}
}

func TestDissolveCoversDotsInShuffledOrder(t *testing.T) {
	clk, covered := newClock(), false
	tt := view(64, 32, pixel.SpringDay, clk, false, func(c *ui.Context) {
		fx.Dissolve(c, covered, fx.DissolveOptions{}, func() { solid(c, 64, 32, pixel.Lit) })
	})
	covered = true
	var prev [][]int
	for step := range 9 {
		tt.Frame()
		d := dots(tt.Image(), pixel.SpringDay, 64, 32)
		if n := count(d, 0, 0, 64, 32, pixel.Ink); n < 2048*step/8-120 || n > 2048*step/8+120 {
			t.Fatalf("step %d covered %d of 2048", step, n)
		}
		for y := range d {
			for x := range d[y] {
				if prev != nil && prev[y][x] == int(pixel.Ink) && d[y][x] != int(pixel.Ink) {
					t.Fatalf("step %d uncovered %d,%d", step, x, y)
				}
			}
		}
		// Halfway, every row and column is part covered.
		if step == 4 {
			for y := range 32 {
				if n := count(d, 0, y, 64, 1, pixel.Ink); n == 0 || n == 64 {
					t.Fatalf("row %d has %d of 64 covered", y, n)
				}
			}
			for x := range 64 {
				if n := count(d, x, 0, 1, 32, pixel.Ink); n == 0 || n == 32 {
					t.Fatalf("column %d has %d of 32 covered", x, n)
				}
			}
		}
		prev = d
		clk.t = clk.t.Add(tick(12))
	}
}

// quilt is 16×16 dots of four slots; neighbors differ.
var quilt = pixel.SpriteOf(16, 16, func(x, y int) (pixel.Slot, bool) {
	return []pixel.Slot{pixel.Good, pixel.Info, pixel.Bad, pixel.Warn}[(x+y*3)%4], true
})

func TestMosaicGrowsBlocksThenSharpens(t *testing.T) {
	clk := newClock()
	tt, cue := cued(clk, false, 16, 16, func(c *ui.Context, cue int) { fx.Mosaic(c, quilt, fx.MosaicOptions{Cue: cue}) })
	plain := dots(tt.Image(), pixel.SpringDay, 16, 16)
	*cue++
	for step, b := range []int{1, 2, 3, 4, 6, 8, 6, 4, 3, 2, 1} {
		tt.Frame()
		d := dots(tt.Image(), pixel.SpringDay, 16, 16)
		for y := range 16 {
			for x := range 16 {
				if d[y][x] != d[y/b*b][x/b*b] {
					t.Fatalf("step %d: block %d is not one slot at %d,%d", step, b, x, y)
				}
			}
		}
		if b == 1 != same(d, plain) {
			t.Fatalf("step %d: blocks of %d, picture changed %v", step, b, !same(d, plain))
		}
		clk.t = clk.t.Add(100 * time.Millisecond)
	}
}

func TestMosaicBlocksTakeTheirCommonestSlot(t *testing.T) {
	// One Good dot in a field of Bad.
	sp := pixel.SpriteOf(8, 8, func(x, y int) (pixel.Slot, bool) {
		if x == 0 && y == 0 {
			return pixel.Good, true
		}
		return pixel.Bad, true
	})
	clk := newClock()
	tt, cue := cued(clk, false, 8, 8, func(c *ui.Context, cue int) { fx.Mosaic(c, sp, fx.MosaicOptions{Cue: cue}) })
	*cue++
	tt.Frame()
	clk.t = clk.t.Add(500*time.Millisecond + time.Microsecond)
	tt.Frame()
	if d := dots(tt.Image(), pixel.SpringDay, 8, 8); count(d, 0, 0, 8, 8, pixel.Bad) != 64 {
		t.Fatal("an 8-dot block must take its commonest slot")
	}
}

func TestGlitchScramblesAFewFramesThenSettles(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		clk := newClock()
		tt, cue := cued(clk, reduce, 22, 16, func(c *ui.Context, cue int) { fx.Glitch(c, quilt, fx.GlitchOptions{Cue: cue}) })
		plain := dots(tt.Image(), pixel.SpringDay, 22, 16)
		*cue++
		changed := false
		for range 6 {
			tt.Frame()
			changed = changed || !same(plain, dots(tt.Image(), pixel.SpringDay, 22, 16))
			clk.t = clk.t.Add(tick(15))
		}
		tt.Frame()
		if changed == reduce || !same(plain, dots(tt.Image(), pixel.SpringDay, 22, 16)) {
			t.Fatalf("reduce %v: glitched %v, then not back to the sprite", reduce, changed)
		}
	}
}

func TestCuesSpentUnderReduceMotionNeverReplay(t *testing.T) {
	for name, build := range map[string]func(c *ui.Context, cue int){
		"Lightning": func(c *ui.Context, cue int) { fx.Lightning(c, fx.LightningOptions{Cue: cue, Width: 32, Height: 32}) },
		"Explosion": func(c *ui.Context, cue int) { fx.Explosion(c, fx.ExplosionOptions{Cue: cue, Size: 32}) },
		"Mosaic":    func(c *ui.Context, cue int) { fx.Mosaic(c, quilt, fx.MosaicOptions{Cue: cue}) },
		"Glitch":    func(c *ui.Context, cue int) { fx.Glitch(c, quilt, fx.GlitchOptions{Cue: cue}) },
		"Flash": func(c *ui.Context, cue int) {
			fx.Flash(c, fx.FlashOptions{Cue: cue}, func() { solid(c, 32, 16, pixel.Info) })
		},
		"ScreenShake": func(c *ui.Context, cue int) {
			fx.ScreenShake(c, fx.ScreenShakeOptions{Cue: cue}, func() { solid(c, 20, 10, pixel.Info) })
		},
	} {
		clk := newClock()
		tt, cue := cued(clk, false, 40, 40, build)
		before := dots(tt.Image(), pixel.SpringDay, 40, 40)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		*cue++
		tt.Frame()
		tt.SetPreferences(ui.Preferences{ReduceMotion: false, TextScale: 1})
		for _, d := range []time.Duration{0, 300 * time.Millisecond} {
			clk.t = clk.t.Add(d)
			tt.Frame()
			if !same(before, dots(tt.Image(), pixel.SpringDay, 40, 40)) {
				t.Errorf("%s replayed a cue spent under Reduce Motion", name)
			}
		}
	}
}

func TestTransitionsKeepTheirEndAfterReduceMotion(t *testing.T) {
	clk, covered := newClock(), false
	tt := view(32, 16, pixel.SpringDay, clk, false, func(c *ui.Context) {
		fx.FadeDither(c, covered, fx.FadeDitherOptions{}, func() { solid(c, 32, 16, pixel.Lit) })
	})
	covered = true
	tt.Frame()
	clk.t = clk.t.Add(tick(12))
	tt.Frame()
	for _, reduce := range []bool{true, false} {
		tt.SetPreferences(ui.Preferences{ReduceMotion: reduce, TextScale: 1})
		tt.Frame()
		if n := count(dots(tt.Image(), pixel.SpringDay, 32, 16), 0, 0, 32, 16, pixel.Ink); n != 512 {
			t.Fatalf("reduce %v: %d of 512 covered", reduce, n)
		}
	}
}

// blank is whether a box shows only the page.
func blank(tt *ui.Tester, w, h int) bool {
	return count(dots(tt.Image(), pixel.SpringDay, w, h), 0, 0, w, h, pixel.Back) == w*h
}

func TestOneShotsRunOnlyOnTheirCue(t *testing.T) {
	for name, tc := range map[string]struct {
		fps   int
		steps int
		build func(c *ui.Context, cue int)
	}{
		"Explosion": {12, 8, func(c *ui.Context, cue int) { fx.Explosion(c, fx.ExplosionOptions{Cue: cue}) }},
		"Lightning": {15, 5, func(c *ui.Context, cue int) { fx.Lightning(c, fx.LightningOptions{Cue: cue, Width: 48, Height: 48}) }},
	} {
		for _, reduce := range []bool{false, true} {
			clk := newClock()
			tt, cue := cued(clk, reduce, 48, 48, tc.build)
			if !blank(tt, 48, 48) {
				t.Fatalf("%s showed before its cue", name)
			}
			*cue++
			for step := range tc.steps {
				tt.Frame()
				if blank(tt, 48, 48) != reduce {
					t.Fatalf("%s, reduce %v: step %d blank %v", name, reduce, step, !reduce)
				}
				clk.t = clk.t.Add(tick(tc.fps))
			}
			tt.Frame()
			if !blank(tt, 48, 48) {
				t.Fatalf("%s outlived its steps", name)
			}
		}
	}
}

func TestLightningFlashesThenForks(t *testing.T) {
	bolt := func(cue int) [][]int {
		clk := newClock()
		tt, c := cued(clk, false, 48, 48, func(c *ui.Context, cue int) {
			fx.Lightning(c, fx.LightningOptions{Cue: cue, Width: 48, Height: 48})
		})
		*c = cue
		tt.Frame()
		if d := dots(tt.Image(), pixel.SpringDay, 48, 48); count(d, 0, 0, 48, 48, pixel.Lit) != 48*48 {
			t.Fatal("a strike opens with a full flash")
		}
		clk.t = clk.t.Add(2 * tick(15))
		tt.Frame()
		return dots(tt.Image(), pixel.SpringDay, 48, 48)
	}
	a, b := bolt(1), bolt(2)
	if count(a, 0, 0, 48, 48, pixel.Warn) == 0 || count(a, 0, 47, 48, 1, pixel.Warn) == 0 {
		t.Fatal("the bolt must reach the ground")
	}
	if same(a, b) {
		t.Fatal("two strikes drew the same bolt")
	}
}

func TestSparkleCirclesItsMiddle(t *testing.T) {
	clk := newClock()
	tt := view(40, 40, pixel.SpringDay, clk, false, func(c *ui.Context) { fx.Sparkle(c, fx.SparkleOptions{}) })
	var quarters [4]int
	for range 8 {
		d := dots(tt.Image(), pixel.SpringDay, 40, 40)
		if count(d, 0, 0, 24, 24, pixel.Warn) != count(d, 0, 0, 40, 40, pixel.Warn) {
			t.Fatal("a star left the square")
		}
		for i, q := range [][2]int{{0, 0}, {12, 0}, {0, 12}, {12, 12}} {
			quarters[i] += count(d, q[0], q[1], 12, 12, pixel.Warn)
		}
		clk.t = clk.t.Add(tick(10))
		tt.Frame()
	}
	for i, n := range quarters {
		if n == 0 {
			t.Fatalf("no star twinkled in quarter %d", i)
		}
	}
}
