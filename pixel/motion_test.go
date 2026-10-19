package pixel_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestCycledRemapsTheBand(t *testing.T) {
	sp := pixel.MustSprite("agcr")
	band := []pixel.Slot{pixel.Accent, pixel.Good, pixel.Info}
	for _, tc := range []struct {
		draw func(g pixel.Canvas) pixel.Canvas
		want []pixel.Slot
	}{
		{func(g pixel.Canvas) pixel.Canvas { return g.Cycled(band, 1) }, []pixel.Slot{pixel.Good, pixel.Info, pixel.Accent, pixel.Bad}},
		{func(g pixel.Canvas) pixel.Canvas { return g.Cycled(band, -1) }, []pixel.Slot{pixel.Info, pixel.Accent, pixel.Good, pixel.Bad}},
		// The later band maps first, the earlier one after it.
		{func(g pixel.Canvas) pixel.Canvas {
			return g.Cycled(band[:2], 1).Cycled(band[1:], 1)
		}, []pixel.Slot{pixel.Good, pixel.Info, pixel.Accent, pixel.Bad}},
	} {
		d := canvas(t, 4, 1, func(g pixel.Canvas) { tc.draw(g).Sprite(0, 0, sp, 0, pixel.Ink, 1) })
		for i, want := range tc.want {
			if got := d.at(i, 0); got != want {
				t.Fatalf("dot %d drew %v, want %v", i, got, want)
			}
		}
	}
	// The original canvas keeps its slots.
	d := canvas(t, 2, 1, func(g pixel.Canvas) {
		g.Cycled(band, 1)
		g.Sprite(0, 0, sp, 0, pixel.Ink, 1)
	})
	if d.at(0, 0) != pixel.Accent {
		t.Fatal("Cycled changed the canvas it came from")
	}
	for name, band := range map[string][]pixel.Slot{"repeat": {pixel.Good, pixel.Good}, "overflow": {pixel.Good, pixel.Slots}, "empty": nil} {
		mustPanic(t, name, func() { canvas(t, 1, 1, func(g pixel.Canvas) { g.Cycled(band, 1) }) })
	}
}

func TestSpriteAtAndSpriteOf(t *testing.T) {
	sp := pixel.MustSprite("a#.\ng.y\n\nyyy\nyyy")
	for _, tc := range []struct {
		frame, x, y int
		want        pixel.Slot
		ok          bool
	}{{0, 0, 0, pixel.Accent, true}, {0, 1, 0, pixel.Warn, true}, {0, 2, 0, 0, false}, {0, 2, 1, pixel.Warn, true}, {1, 1, 1, pixel.Warn, true}, {2, 0, 1, pixel.Good, true}} {
		if got, ok := sp.At(tc.frame, tc.x, tc.y, pixel.Warn); got != tc.want || ok != tc.ok {
			t.Fatalf("At(%d, %d, %d) = %v %v, want %v %v", tc.frame, tc.x, tc.y, got, ok, tc.want, tc.ok)
		}
	}
	mustPanic(t, "At outside", func() { sp.At(0, 3, 0, pixel.Ink) })
	of := pixel.SpriteOf(3, 1, func(x, y int) (pixel.Slot, bool) { return pixel.Slot(x + 11), x != 1 })
	d := canvas(t, 3, 1, func(g pixel.Canvas) { g.Sprite(0, 0, of, 0, pixel.Ink, 1) })
	if d.at(0, 0) != pixel.Good || d.at(1, 0) != pixel.Back || d.at(2, 0) != pixel.Bad {
		t.Fatal("SpriteOf drew the wrong dots")
	}
	mustPanic(t, "SpriteOf empty", func() { pixel.SpriteOf(0, 1, func(x, y int) (pixel.Slot, bool) { return 0, true }) })
	mustPanic(t, "SpriteOf slot", func() { pixel.SpriteOf(1, 1, func(x, y int) (pixel.Slot, bool) { return pixel.Slots, true }) })
}

func TestShotSpendsCuesUnderReduceMotion(t *testing.T) {
	at := time.Unix(1_700_000_000, 0)
	cue, step := 0, -1
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }})
		step = pixel.K(c).Shot(c, ui.Box(c), cue, 10, 5)
	}, 10, 10)
	frame := func(want int) {
		t.Helper()
		tt.Frame()
		if step != want {
			t.Fatalf("cue %d: step %d, want %d", cue, step, want)
		}
	}
	// No cue yet: the shot is over.
	frame(5)
	cue = 1
	frame(0)
	at = at.Add(200 * time.Millisecond)
	frame(2)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	cue = 2
	frame(5)
	tt.SetPreferences(ui.Preferences{ReduceMotion: false, TextScale: 1})
	frame(5)
	cue = 3
	frame(0)
}

func TestTowardSettlesUnderReduceMotion(t *testing.T) {
	at := time.Unix(1_700_000_000, 0)
	on, step := false, -1
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }})
		step = pixel.K(c).Toward(c, ui.Box(c), on, 8, 4)
	}, 10, 10)
	on = true
	tt.Frame()
	at = at.Add(125 * time.Millisecond)
	tt.Frame()
	// Reduce Motion jumps to the end, and the end stays.
	for _, reduce := range []bool{true, false} {
		tt.SetPreferences(ui.Preferences{ReduceMotion: reduce, TextScale: 1})
		tt.Frame()
		if step != 4 {
			t.Fatalf("reduce %v: step %d, want 4", reduce, step)
		}
	}
}

func TestTowardStepsAndTurnsWhereItIs(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		at := time.Unix(1_700_000_000, 0)
		on, step := true, -1
		tt := ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }})
			e := ui.Box(c)
			step = pixel.K(c).Toward(c, e, on, 8, 4)
		}, 10, 10)
		if reduce {
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		}
		frame := func(d time.Duration, want int) {
			t.Helper()
			at = at.Add(d)
			tt.Frame()
			if !reduce && step != want || reduce && step != 4*boolInt(on) {
				t.Fatalf("reduce %v, on %v: step %d, want %d", reduce, on, step, want)
			}
		}
		// Shown first, it is already at its end.
		frame(0, 4)
		on = false
		frame(0, 4)
		frame(125*time.Millisecond, 3)
		frame(125*time.Millisecond, 2)
		// Turned at 2, it climbs back from 2.
		on = true
		frame(0, 2)
		frame(125*time.Millisecond, 3)
		frame(time.Second, 4)
	}
	mustPanic(t, "no steps", func() {
		ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay})
			pixel.K(c).Toward(c, ui.Box(c), true, 8, 0)
		}, 10, 10)
	})
}

func TestMotionStepsAtFourToFifteen(t *testing.T) {
	at := time.Unix(1_700_000_000, 0)
	run := func(fn func(c *ui.Context, k pixel.Kit)) {
		ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }})
			fn(c, pixel.K(c))
		}, 10, 10)
	}
	// S5: whole frames at 4–15 fps, for every clock.
	for _, fps := range []int{1, 2, 3, 16, 30} {
		mustPanic(t, fmt.Sprint("Tick at ", fps), func() { run(func(c *ui.Context, k pixel.Kit) { k.Tick(c, fps) }) })
		mustPanic(t, fmt.Sprint("Since at ", fps), func() { run(func(c *ui.Context, k pixel.Kit) { k.Since(c, at, fps, 9) }) })
		mustPanic(t, fmt.Sprint("Toward at ", fps), func() { run(func(c *ui.Context, k pixel.Kit) { k.Toward(c, ui.Box(c), true, fps, 2) }) })
		mustPanic(t, fmt.Sprint("Since from zero at ", fps), func() { run(func(c *ui.Context, k pixel.Kit) { k.Since(c, time.Time{}, fps, 9) }) })
		mustPanic(t, fmt.Sprint("Shot at ", fps), func() { run(func(c *ui.Context, k pixel.Kit) { k.Shot(c, ui.Box(c), 0, fps, 9) }) })
		// Reduce Motion spends a shot without counting.
		mustPanic(t, fmt.Sprint("a reduced Shot at ", fps), func() {
			ui.NewTester(func(c *ui.Context) {
				pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }, Reduce: true})
				pixel.K(c).Shot(c, ui.Box(c), 0, fps, 9)
			}, 10, 10)
		})
		mustPanic(t, fmt.Sprint("Step at ", fps), func() {
			run(func(c *ui.Context, k pixel.Kit) {
				ui.Box(c).Size(4, 4).Draw(func(p *ui.Painter, r ui.Rect) { k.Canvas(p, r).Step(fps) })
			})
		})
	}
	for _, fps := range []int{4, 15} {
		run(func(c *ui.Context, k pixel.Kit) {
			k.Tick(c, fps)
			k.Since(c, at, fps, 9)
			k.Toward(c, ui.Box(c), true, fps, 2)
		})
	}
}

func TestSinceHoldsAtZeroWhenTheClockStepsBack(t *testing.T) {
	at := time.Unix(1_700_000_000, 0)
	got := -1
	ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }})
		got = pixel.K(c).Since(c, at.Add(time.Hour), 8, 5)
	}, 10, 10)
	if got != 0 {
		t.Fatalf("a start an hour ahead counted %d steps", got)
	}
	// Thirty years ahead must not overflow into a full count.
	ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }})
		got = pixel.K(c).Since(c, at.AddDate(30, 0, 0), 15, 5)
	}, 10, 10)
	if got != 0 {
		t.Fatalf("a start thirty years ahead counted %d steps", got)
	}
	ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }})
		got = pixel.K(c).Since(c, at.AddDate(-30, 0, 0), 15, 5)
	}, 10, 10)
	if got != 5 {
		t.Fatalf("a start thirty years back counted %d of 5 steps", got)
	}
}

func TestElapsedRebasesWhenTheClockStepsBack(t *testing.T) {
	at := time.Unix(1_700_000_000, 0)
	start, got := at.Add(time.Hour), time.Duration(-1)
	ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }})
		got = pixel.K(c).Elapsed(&start)
	}, 10, 10)
	if got != 0 || !start.Equal(at) {
		t.Fatalf("elapsed %v, start moved to %v", got, start)
	}
}
