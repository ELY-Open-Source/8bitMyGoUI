package game_test

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/game"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestHealthBarTrailHoldsThenDrains(t *testing.T) {
	wd := &world{clk: newClock()}
	hp := 80
	tt := view(300, 60, wd, func(c *ui.Context) { game.HealthBar(c, "HP", hp, 100, game.HealthBarOptions{}) })
	// The fill's middle row is 10; full is 92 dots.
	row := func(s pixel.Slot) int { return count(t, tt.Image(), 0, 10, 96, 1, s) }
	if row(pixel.Good) != 73 {
		t.Fatalf("80 HP fills %d dots", row(pixel.Good))
	}
	hp = 40
	tt.Frame()
	if row(pixel.Warn) != 36 || row(pixel.AccentLit) != 37 {
		t.Fatalf("at the hit: fill %d, trail %d", row(pixel.Warn), row(pixel.AccentLit))
	}
	for _, step := range []struct {
		after time.Duration
		trail int
	}{{300 * time.Millisecond, 37}, {200 * time.Millisecond, 19}, {500 * time.Millisecond, 0}} {
		wd.clk.t = wd.clk.t.Add(step.after)
		tt.Frame()
		if row(pixel.AccentLit) != step.trail {
			t.Fatalf("%v later the trail is %d dots", step.after, row(pixel.AccentLit))
		}
	}
	hp = 90
	tt.Frame()
	if row(pixel.Good) != 82 || row(pixel.AccentLit) != 0 {
		t.Fatal("a heal left a trail")
	}
	// A heal during a live trail clears it too.
	hp = 30
	tt.Frame()
	hp = 60
	tt.Frame()
	if row(pixel.AccentLit) != 0 {
		t.Fatal("a heal kept a live trail")
	}
	keeps(t, tt.Image())
	// A new full drops the old trail.
	full := 100
	tt = view(300, 60, wd, func(c *ui.Context) { game.HealthBar(c, "HP", hp, full, game.HealthBarOptions{}) })
	hp, full = 40, 50
	tt.Frame()
	if row(pixel.AccentLit) != 0 || count(t, tt.Image(), 96, 0, 54, 30, pixel.AccentLit) != 0 {
		t.Fatal("the trail outlived a new full")
	}
	panics(t, "a bar narrower than its labels", func() {
		view(200, 60, wd, func(c *ui.Context) { game.HealthBar(c, "HP", 5, 100, game.HealthBarOptions{Width: 20}) })
	})
	panics(t, "hp past full", func() {
		view(200, 60, wd, func(c *ui.Context) { game.HealthBar(c, "HP", 5, 4, game.HealthBarOptions{}) })
	})
}

func TestHealthBarReducedHasNoTrailAndNoBlink(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		hp := 90
		tt := view(300, 60, wd, func(c *ui.Context) { game.HealthBar(c, "HP", hp, 100, game.HealthBarOptions{}) })
		hp = 20
		tt.Frame()
		if trail := count(t, tt.Image(), 0, 10, 96, 1, pixel.AccentLit); (trail == 0) != reduce {
			t.Fatalf("reduced %v: trail of %d dots", reduce, trail)
		}
		// Low HP blinks at 2 Hz: Bad, then Lit.
		seen := map[pixel.Slot]bool{}
		for range 2 {
			seen[at(t, tt.Image(), 4, 10)] = true
			wd.clk.t = wd.clk.t.Add(250 * time.Millisecond)
			tt.Frame()
		}
		if blinked := seen[pixel.Lit] && seen[pixel.Bad]; blinked == reduce {
			t.Fatalf("reduced %v: low HP showed %v", reduce, seen)
		}
	}
}

func TestHugeHPKeepsItsColor(t *testing.T) {
	tt := view(500, 60, &world{clk: newClock()}, func(c *ui.Context) {
		game.HealthBar(c, "HP", math.MaxInt64/2, math.MaxInt64, game.HealthBarOptions{Width: 200})
	})
	// Half the largest HP is Warn, not wrapped to Bad.
	if n := count(t, tt.Image(), 0, 10, 200, 1, pixel.Warn); n != 98 {
		t.Fatalf("half of the largest HP fills %d Warn dots", n)
	}
	tt = view(500, 60, &world{clk: newClock()}, func(c *ui.Context) {
		game.ManaBar(c, "MP", 1, math.MaxInt64/2, game.ManaBarOptions{Width: 200, Tick: math.MaxInt64/4 + 1})
	})
	if n := count(t, tt.Image(), 2, 10, 196, 1, pixel.Edge); n != 1 {
		t.Fatalf("one tick at the middle came out as %d", n)
	}
	// Ticks stop at the top; dense ticks drop out.
	for _, tc := range []struct{ full, tick, want int }{{math.MaxInt64, math.MaxInt64 - 1, 1}, {1000, 10, 0}} {
		tt = view(500, 60, &world{clk: newClock()}, func(c *ui.Context) {
			game.ManaBar(c, "MP", 1, tc.full, game.ManaBarOptions{Width: 200, Tick: tc.tick})
		})
		if n := count(t, tt.Image(), 2, 10, 196, 1, pixel.Edge); n != tc.want {
			t.Fatalf("full %d, tick %d drew %d ticks", tc.full, tc.tick, n)
		}
	}
}

func TestManaBarTicks(t *testing.T) {
	tt := view(300, 60, &world{clk: newClock()}, func(c *ui.Context) { game.ManaBar(c, "MP", 35, 60, game.ManaBarOptions{}) })
	img := tt.Image()
	// 35 of 60 fills 53 dots; three ticks cut in.
	if info, ticks := count(t, img, 0, 10, 96, 1, pixel.Info), count(t, img, 2, 10, 92, 1, pixel.Edge); info != 50 || ticks != 5 {
		t.Fatalf("fill %d dots, %d ticks", info, ticks)
	}
}

func TestExpBarFlashesWhenItFills(t *testing.T) {
	for _, tc := range []struct {
		name          string
		level, xp     int
		reduce, flash bool
	}{
		{"a level up", 8, 10, false, true}, {"a fill", 7, 300, false, true},
		{"a reduced level up", 8, 10, true, false}, {"a lost level", 6, 10, false, false},
	} {
		wd := &world{clk: newClock(), reduce: tc.reduce}
		level, xp := 7, 290
		tt := view(300, 60, wd, func(c *ui.Context) { game.ExpBar(c, "EXP", level, xp, 300, game.ExpBarOptions{}) })
		level, xp = tc.level, tc.xp
		tt.Frame()
		// With no flash, the first frame already shows xp.
		if n, want := count(t, tt.Image(), 2, 10, 92, 1, pixel.Warn), tc.xp*92/300; !tc.flash && n != want {
			t.Fatalf("%s: the first frame fills %d dots, not %d", tc.name, n, want)
		}
		seen := map[pixel.Slot]bool{}
		for range 2 {
			img := tt.Image()
			seen[at(t, img, 4, 10)] = true
			if full := count(t, img, 2, 10, 92, 1, pixel.Lit) + count(t, img, 2, 10, 92, 1, pixel.Warn); tc.flash && full != 92 {
				t.Fatalf("%s: the flash fills %d dots", tc.name, full)
			}
			wd.clk.t = wd.clk.t.Add(100 * time.Millisecond)
			tt.Frame()
		}
		if flashed := seen[pixel.Lit] && seen[pixel.Warn]; flashed != tc.flash {
			t.Fatalf("%s: the bar showed %v", tc.name, seen)
		}
		// After any flash the bar shows xp.
		wd.clk.t = wd.clk.t.Add(time.Second)
		tt.Frame()
		if n, want := count(t, tt.Image(), 2, 10, 92, 1, pixel.Warn), tc.xp*92/300; n != want {
			t.Fatalf("%s: %d of 300 fills %d dots, not %d", tc.name, tc.xp, n, want)
		}
	}
}

func TestHeartsFillByQuarters(t *testing.T) {
	tt := view(300, 60, &world{clk: newClock()}, func(c *ui.Context) {
		game.Hearts(c, "Life", 1.75, 3, game.HeartsOptions{})
	})
	img := tt.Image()
	// Quadrant dots: bottom left, top left, top right, bottom right.
	quads := [][2]int{{4, 7}, {3, 4}, {8, 4}, {7, 7}}
	for i, want := range [][4]bool{{true, true, true, true}, {true, true, true, false}, {false, false, false, false}} {
		for q, d := range quads {
			if got := at(t, img, i*11+d[0], d[1]) == pixel.Bad; got != want[q] {
				t.Fatalf("heart %d quadrant %d filled %v", i, q, got)
			}
		}
	}
	keeps(t, img)
	// A lone quarter fills the bottom left only.
	tt = view(100, 40, &world{clk: newClock()}, func(c *ui.Context) { game.Hearts(c, "Life", 0.25, 1, game.HeartsOptions{}) })
	for q, d := range quads {
		if got := at(t, tt.Image(), d[0], d[1]) == pixel.Bad; got != (q == 0) {
			t.Fatalf("a quarter heart filled quadrant %d: %v", q, got)
		}
	}
	for _, bad := range [][2]float64{{0.3, 3}, {1, 2.5}, {4, 3}, {-1, 3}} {
		panics(t, "hearts out of quarters", func() {
			view(100, 40, &world{clk: newClock()}, func(c *ui.Context) { game.Hearts(c, "x", bad[0], bad[1], game.HeartsOptions{}) })
		})
	}
}

func TestLivesShowTheirIcon(t *testing.T) {
	tt := view(300, 100, &world{clk: newClock()}, func(c *ui.Context) {
		game.Lives(c, "Lives", 3, game.LivesOptions{})
		game.Lives(c, "Stars", 2, game.LivesOptions{Sprite: icon.Star})
	})
	lives, stars := find(t, tt, "Lives"), find(t, tt, "Stars")
	lx, ly := dots(lives)
	sx, sy := dots(stars)
	// The hero wears a Bad cap; the star is ink.
	if count(t, tt.Image(), lx, ly, 12, 16, pixel.Bad) == 0 || count(t, tt.Image(), sx, sy, 12, 16, pixel.Bad) != 0 {
		t.Fatal("life icons are wrong")
	}
	if int(lives.W)/2 != 12+2+pixel.Measure("×3") {
		t.Fatalf("lives are %v px wide", lives.W)
	}
}

func TestCoinSpinsUnlessReduced(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		tt := view(300, 60, wd, func(c *ui.Context) { game.CoinCounter(c, "Gold", 1250, game.CoinCounterOptions{}) })
		find(t, tt, "Gold")
		before := count(t, tt.Image(), 0, 0, 12, 16, pixel.Warn)
		wd.clk.t = wd.clk.t.Add(250 * time.Millisecond)
		tt.Frame()
		if spun := count(t, tt.Image(), 0, 0, 12, 16, pixel.Warn) != before; spun == reduce {
			t.Fatalf("reduced %v: the coin spun %v", reduce, spun)
		}
	}
}

func TestScoreBoardLightsANewBest(t *testing.T) {
	score := 500
	tt := view(400, 100, &world{clk: newClock()}, func(c *ui.Context) { game.ScoreBoard(c, score, 9900, game.ScoreBoardOptions{}) })
	if !tt.HasText("SCORE 500, HI-SCORE 9900") || count(t, tt.Image(), 0, 0, 200, 8, pixel.Warn) != 0 {
		t.Fatalf("texts %q", tt.Texts())
	}
	score = 12500
	tt.Frame()
	if !tt.HasText("SCORE 12500, HI-SCORE 12500") || count(t, tt.Image(), 0, 0, 200, 8, pixel.Warn) == 0 {
		t.Fatal("a new best is not lit")
	}
}

func TestTimerHUDHurriesUnderThirty(t *testing.T) {
	for _, tc := range []struct {
		left          time.Duration
		reduce, blink bool
		ink           pixel.Slot
	}{{87 * time.Second, false, false, pixel.Ink}, {12 * time.Second, false, true, pixel.Bad}, {12 * time.Second, true, false, pixel.Bad},
		{29500 * time.Millisecond, true, false, pixel.Bad}} {
		wd := &world{clk: newClock(), reduce: tc.reduce}
		tt := view(300, 100, wd, func(c *ui.Context) { game.TimerHUD(c, tc.left, game.TimerHUDOptions{}) })
		var shown []int
		for range 2 {
			shown = append(shown, count(t, tt.Image(), 0, 8, 60, 34, tc.ink))
			wd.clk.t = wd.clk.t.Add(250 * time.Millisecond)
			tt.Frame()
		}
		if blinked := (shown[0] == 0) != (shown[1] == 0); blinked != tc.blink || shown[0]+shown[1] == 0 {
			t.Fatalf("%v left, reduced %v: digit dots %v", tc.left, tc.reduce, shown)
		}
	}
}
