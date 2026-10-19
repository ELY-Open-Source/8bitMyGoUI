package fx_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/fx"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestPaletteCycleRotatesTheBand(t *testing.T) {
	band := []pixel.Slot{pixel.Accent, pixel.Good, pixel.Info}
	sp := pixel.MustSprite("agcr")
	clk := newClock()
	tt := view(8, 2, pixel.SpringDay, clk, false, func(c *ui.Context) {
		fx.PaletteCycle(c, sp, band, fx.PaletteCycleOptions{})
	})
	for range 4 {
		n, d := clk.step(8), dots(tt.Image(), pixel.SpringDay, 4, 1)
		for i := range band {
			if want := band[(i+n)%3]; d[0][i] != int(want) {
				t.Fatalf("step %d: dot %d is %d, want %v", n, i, d[0][i], want)
			}
		}
		if d[0][3] != int(pixel.Bad) {
			t.Fatal("a slot outside the band changed")
		}
		clk.t = clk.t.Add(125 * time.Millisecond)
		tt.Frame()
	}
	mustPanic(t, "one slot", func() {
		view(8, 2, pixel.SpringDay, clk, false, func(c *ui.Context) { fx.PaletteCycle(c, sp, band[:1], fx.PaletteCycleOptions{}) })
	})
}

func TestSkyBandsByDayStarsByNight(t *testing.T) {
	d := dots(view(80, 48, pixel.SpringDay, newClock(), false, func(c *ui.Context) {
		fx.Sky(c, fx.SkyOptions{Width: 80, Height: 48})
	}).Image(), pixel.SpringDay, 80, 48)
	// By day the top band is plain Info: no stars.
	if count(d, 0, 0, 80, 24, pixel.Info) != 80*24 || count(d, 0, 33, 80, 1, pixel.AccentLit) != 80 || count(d, 0, 47, 80, 1, pixel.Lit) != 80 {
		t.Fatal("the day sky must run Info, AccentLit, Lit")
	}
	n := dots(view(80, 48, pixel.WinterNight, newClock(), false, func(c *ui.Context) {
		fx.Sky(c, fx.SkyOptions{Width: 80, Height: 48})
	}).Image(), pixel.WinterNight, 80, 48)
	if count(n, 0, 0, 80, 32, pixel.Ink) == 0 || count(n, 0, 32, 80, 16, pixel.Ink) != 0 {
		t.Fatal("stars belong to the top two thirds of the night")
	}
}

func TestSunMoonRidesItsArc(t *testing.T) {
	at := func(hour, x, y int) int {
		d := dots(view(160, 96, pixel.SpringDay, newClock(), false, func(c *ui.Context) {
			fx.SunMoon(c, hour, fx.SunMoonOptions{})
		}).Image(), pixel.SpringDay, 160, 96)
		if hour >= 6 && hour < 18 == (count(d, 0, 0, 160, 96, pixel.Warn) == 0) {
			t.Fatalf("hour %d: the sun is out by day only", hour)
		}
		return d[y][x]
	}
	// Noon and midnight top the arc; six sits low left.
	for _, tc := range []struct {
		hour, x, y int
		want       pixel.Slot
	}{{12, 85, 10, pixel.Warn}, {6, 15, 63, pixel.Warn}, {0, 85, 10, pixel.Lit}, {18, 15, 63, pixel.Lit}} {
		if got := at(tc.hour, tc.x, tc.y); got != int(tc.want) {
			t.Errorf("hour %d: dot %d,%d is %d, want %v", tc.hour, tc.x, tc.y, got, tc.want)
		}
	}
	for _, hour := range []int{-1, 24} {
		mustPanic(t, "hour", func() {
			view(16, 16, pixel.SpringDay, newClock(), false, func(c *ui.Context) { fx.SunMoon(c, hour, fx.SunMoonOptions{}) })
		})
	}
}

// ridge is each column's top s dot, or the height.
func ridge(d [][]int, s pixel.Slot) []int {
	out := make([]int, len(d[0]))
	for x := range out {
		out[x] = len(d)
		for y := range d {
			if d[y][x] == int(s) {
				out[x] = y
				break
			}
		}
	}
	return out
}

func TestMountainRidgesScrollAtTheirSpeeds(t *testing.T) {
	clk := newClock()
	tt := view(80, 48, pixel.SpringDay, clk, false, func(c *ui.Context) {
		fx.Mountains(c, fx.MountainsOptions{Width: 80, Height: 48})
	})
	near, far := ridge(dots(tt.Image(), pixel.SpringDay, 80, 48), pixel.Edge), ridge(dots(tt.Image(), pixel.SpringDay, 80, 48), pixel.Shade)
	// The near ridge moves every two steps, the far eight.
	clk.t = clk.t.Add(250 * time.Millisecond)
	tt.Frame()
	if now := ridge(dots(tt.Image(), pixel.SpringDay, 80, 48), pixel.Edge); !equalShift(near, now, -1) {
		t.Fatalf("near ridge %v became %v", near, now)
	}
	clk.t = clk.t.Add(750 * time.Millisecond)
	tt.Frame()
	if now := ridge(dots(tt.Image(), pixel.SpringDay, 80, 48), pixel.Shade); !equalShift(far, now, 48) {
		t.Fatalf("far ridge %v became %v", far, now)
	}
}

// equalShift is whether b is a, one column left.
func equalShift(a, b []int, hidden int) bool {
	shown := 0
	for x := range len(a) - 1 {
		if b[x] == hidden || a[x+1] == hidden {
			continue
		}
		if b[x] != a[x+1] {
			return false
		}
		shown++
	}
	return shown > len(a)/2
}

func TestOceanCrestsLightUpInTurn(t *testing.T) {
	clk := newClock()
	tt := view(80, 48, pixel.SpringDay, clk, false, func(c *ui.Context) { fx.Ocean(c, fx.OceanOptions{Width: 80, Height: 48}) })
	// The first crest row runs Info, Lit, AccentLit, Info.
	for range 4 {
		d, n := dots(tt.Image(), pixel.SpringDay, 80, 48), clk.step(6)
		lit, pink := count(d, 0, 3, 80, 1, pixel.Lit), count(d, 0, 3, 80, 1, pixel.AccentLit)
		if want := [4][2]bool{{false, false}, {true, false}, {false, true}, {false, false}}[n%4]; (lit > 0) != want[0] || (pink > 0) != want[1] {
			t.Fatalf("step %d: %d Lit and %d AccentLit dots on the first crest", n, lit, pink)
		}
		clk.t = clk.t.Add(tick(6))
		tt.Frame()
	}
}

func TestWaterfallSegmentsFallThreeRowsAStep(t *testing.T) {
	clk := newClock()
	tt := view(40, 48, pixel.SpringDay, clk, false, func(c *ui.Context) {
		fx.Waterfall(c, fx.WaterfallOptions{Width: 40, Height: 48})
	})
	before := dots(tt.Image(), pixel.SpringDay, 40, 48)
	clk.t = clk.t.Add(100 * time.Millisecond)
	tt.Frame()
	after := dots(tt.Image(), pixel.SpringDay, 40, 48)
	lit := 0
	for x := 10; x < 30; x++ {
		for y := 3; y < 39; y++ {
			if before[y-3][x] == int(pixel.Lit) {
				lit++
			}
			if after[y][x] != before[y-3][x] {
				t.Fatalf("water at %d,%d did not fall three rows", x, y)
			}
		}
	}
	// One segment in six is lit.
	if lit < 20*36/6-40 || lit > 20*36/6+40 {
		t.Fatalf("%d lit water dots", lit)
	}
}

func TestFireHeatRisesARowAStep(t *testing.T) {
	clk := newClock()
	tt := view(32, 40, pixel.SpringDay, clk, false, func(c *ui.Context) { fx.Fire(c, fx.FireOptions{}) })
	for range 4 {
		d, n, glows := dots(tt.Image(), pixel.SpringDay, 32, 40), clk.step(10), 0
		for y := range d {
			for x := range d[y] {
				if d[y][x] == int(pixel.Lit) {
					glows++
					if row := y/2 - 2; (row+n)%4 != 0 {
						t.Fatalf("step %d: heat glows on flame row %d", n, row)
					}
				}
			}
		}
		if glows == 0 {
			t.Fatalf("step %d: no heat", n)
		}
		clk.t = clk.t.Add(100 * time.Millisecond)
		tt.Frame()
	}
}

func TestSeasonScenePaintsInItsSeasonsPalette(t *testing.T) {
	for _, tc := range []struct {
		season pixel.Season
		hour   int
		pal    *pixel.Palette
	}{{pixel.Spring, 10, pixel.SpringDay}, {pixel.Summer, 21, pixel.SummerNight}, {pixel.Autumn, 5, pixel.AutumnNight}, {pixel.Winter, 6, pixel.WinterDay}} {
		d := dots(view(80, 48, pixel.SpringDay, newClock(), false, func(c *ui.Context) {
			fx.SeasonScene(c, tc.season, tc.hour, fx.SeasonSceneOptions{Width: 80, Height: 48})
		}).Image(), tc.pal, 80, 48)
		for y := range d {
			for x := range d[y] {
				if d[y][x] < 0 {
					t.Fatalf("%s at %d: dot %d,%d is off its palette", tc.pal.Name, tc.hour, x, y)
				}
			}
		}
	}
	mustPanic(t, "season", func() {
		view(16, 16, pixel.SpringDay, newClock(), false, func(c *ui.Context) { fx.SeasonScene(c, 4, 10, fx.SeasonSceneOptions{}) })
	})
	mustPanic(t, "hour", func() {
		view(16, 16, pixel.SpringDay, newClock(), false, func(c *ui.Context) { fx.SeasonScene(c, pixel.Spring, 24, fx.SeasonSceneOptions{}) })
	})
}

func TestWeatherKeepsToItsSlots(t *testing.T) {
	for name, tc := range map[string]struct {
		build func(c *ui.Context)
		slots []pixel.Slot
	}{
		"Rain":      {func(c *ui.Context) { fx.Rain(c, fx.RainOptions{Width: 80, Height: 48}) }, []pixel.Slot{pixel.Faint}},
		"Snow":      {func(c *ui.Context) { fx.Snow(c, fx.SnowOptions{Width: 80, Height: 48}) }, []pixel.Slot{pixel.Lit}},
		"Petals":    {func(c *ui.Context) { fx.Petals(c, fx.PetalsOptions{Width: 80, Height: 48}) }, []pixel.Slot{pixel.Accent, pixel.AccentLit}},
		"Leaves":    {func(c *ui.Context) { fx.Leaves(c, fx.LeavesOptions{Width: 80, Height: 48}) }, []pixel.Slot{pixel.Warn, pixel.Bad, pixel.Accent}},
		"Fireflies": {func(c *ui.Context) { fx.Fireflies(c, fx.FirefliesOptions{Width: 80, Height: 48}) }, []pixel.Slot{pixel.Good, pixel.Warn}},
		"Confetti": {func(c *ui.Context) { fx.Confetti(c, fx.ConfettiOptions{Width: 80, Height: 48}) },
			[]pixel.Slot{pixel.Accent, pixel.Good, pixel.Warn, pixel.Bad, pixel.Info, pixel.Magic}},
	} {
		d := dots(view(80, 48, pixel.SpringDay, newClock(), false, tc.build).Image(), pixel.SpringDay, 80, 48)
		seen := 0
		for y := range d {
			for x, s := range d[y] {
				if s == int(pixel.Back) {
					continue
				}
				seen++
				ok := false
				for _, want := range tc.slots {
					ok = ok || s == int(want)
				}
				if !ok {
					t.Fatalf("%s drew slot %d at %d,%d", name, s, x, y)
				}
			}
		}
		if seen == 0 {
			t.Fatalf("%s drew nothing", name)
		}
	}
}
