package pixel_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
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
