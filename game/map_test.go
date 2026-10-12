package game_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/game"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestMinimapTilesFogPlayerAndMarks(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		seen := [][]bool{{true, true, false}, {true, true, true}}
		tt := view(100, 60, wd, func(c *ui.Context) {
			game.Minimap(c, "Map", []string{`#.~`, `"^ `}, seen, game.Spot{X: 1, Y: 0}, []game.Marker{{Spot: game.Spot{X: 0, Y: 1}, Slot: pixel.Warn}}, game.MinimapOptions{})
		})
		img := tt.Image()
		// Tiles are 3 dots inside a 1-dot frame.
		for _, tc := range []struct {
			x, y int
			s    pixel.Slot
		}{{0, 0, pixel.Edge}, {0, 1, pixel.Warn}, {1, 1, pixel.Faint}, {2, 1, pixel.Back}} {
			if n := count(t, img, 1+tc.x*3, 1+tc.y*3, 3, 3, tc.s); n != 9 {
				t.Fatalf("tile %d,%d holds %d %v dots", tc.x, tc.y, n, tc.s)
			}
		}
		if fog := count(t, img, 7, 1, 3, 3, pixel.Edge); fog == 0 || fog == 9 {
			t.Fatalf("the fogged tile holds %d Edge dots", fog)
		}
		var lit []int
		for range 2 {
			lit = append(lit, count(t, tt.Image(), 4, 1, 3, 3, pixel.Accent))
			wd.clk.t = wd.clk.t.Add(250 * time.Millisecond)
			tt.Frame()
		}
		if blinked := lit[0] != lit[1]; blinked == reduce || lit[0]+lit[1] == 0 {
			t.Fatalf("reduced %v: the player lit %v", reduce, lit)
		}
	}
	panics(t, "a ragged map", func() {
		view(100, 60, &world{clk: newClock()}, func(c *ui.Context) {
			game.Minimap(c, "x", []string{"##", "#"}, nil, game.Spot{}, nil, game.MinimapOptions{})
		})
	})
	panics(t, "ragged fog", func() {
		view(100, 60, &world{clk: newClock()}, func(c *ui.Context) {
			game.Minimap(c, "x", []string{"##", "##"}, [][]bool{{true, true}, {true}}, game.Spot{}, nil, game.MinimapOptions{})
		})
	})
	panics(t, "a player off the map", func() {
		view(100, 60, &world{clk: newClock()}, func(c *ui.Context) {
			game.Minimap(c, "x", []string{"##"}, nil, game.Spot{X: 2}, nil, game.MinimapOptions{})
		})
	})
}

func TestCompassTurnsWithTheHeading(t *testing.T) {
	shot := func(heading float64) (*ui.Tester, int) {
		tt := view(300, 60, &world{clk: newClock()}, func(c *ui.Context) {
			game.Compass(c, "Heading", heading, []game.Bearing{{Degrees: 70, Slot: pixel.Warn}, {Degrees: 250, Slot: pixel.Bad}}, game.CompassOptions{})
		})
		return tt, count(t, tt.Image(), 0, 0, 120, 20, pixel.Ink)
	}
	tt, ink := shot(30)
	img := tt.Image()
	// N sits 20 dots left of middle, E 40 right.
	if count(t, img, 38, 0, 5, 5, pixel.Ink) == 0 || count(t, img, 98, 0, 5, 5, pixel.Ink) == 0 {
		t.Fatal("N or E is not where 30° puts it")
	}
	if count(t, img, 0, 0, 120, 20, pixel.Warn) != 9 || count(t, img, 0, 0, 120, 20, pixel.Bad) != 0 {
		t.Fatal("markers in and out of view are wrong")
	}
	if _, again := shot(390); again != ink {
		t.Fatal("390° differs from 30°")
	}
	// -360° is 0°, in view of a 350° heading.
	tt = view(300, 60, &world{clk: newClock()}, func(c *ui.Context) {
		game.Compass(c, "Heading", 350, []game.Bearing{{Degrees: -360, Slot: pixel.Magic}}, game.CompassOptions{})
	})
	if count(t, tt.Image(), 0, 0, 120, 20, pixel.Magic) != 9 {
		t.Fatal("a marker at -360° is missing")
	}
	keeps(t, img)
}

func TestBossBarPhases(t *testing.T) {
	tt := view(500, 80, &world{clk: newClock()}, func(c *ui.Context) {
		game.BossBar(c, "SLIME KING", 500, 1000, game.BossBarOptions{Phases: []float64{0.25}})
	})
	img := tt.Image()
	// The quarter mark stands 49 dots into the 196-dot fill.
	if n := count(t, img, 0, 22, 200, 1, pixel.Bad); n != 97 || at(t, img, 51, 22) != pixel.Ink {
		t.Fatalf("half a 196-dot bar is %d Bad dots", n)
	}
	panics(t, "a bar narrower than the name", func() {
		view(500, 80, &world{clk: newClock()}, func(c *ui.Context) {
			game.BossBar(c, "SLIME KING", 1, 1, game.BossBarOptions{Width: 20})
		})
	})
	panics(t, "a phase at 1", func() {
		view(500, 80, &world{clk: newClock()}, func(c *ui.Context) {
			game.BossBar(c, "x", 1, 1, game.BossBarOptions{Phases: []float64{1}})
		})
	})
}

func TestStatusEffectsEdgeByKind(t *testing.T) {
	tt := view(200, 60, &world{clk: newClock()}, func(c *ui.Context) {
		game.StatusEffects(c, "Effects", []game.Effect{{Sprite: icon.Bolt, Name: "Haste", Turns: 3}, {Sprite: icon.Skull, Name: "Poison", Debuff: true}}, game.StatusEffectsOptions{})
	})
	img := tt.Image()
	if at(t, img, 5, 1) != pixel.Good || at(t, img, 17+5, 1) != pixel.Bad {
		t.Fatal("edges do not tell buffs from debuffs")
	}
	// Only a count above zero gets a badge.
	if count(t, img, 0, 9, 16, 7, pixel.Good) == 0 || count(t, img, 17, 9, 16, 7, pixel.Bad) != 0 {
		t.Fatal("turn badges are wrong")
	}
}

func TestDamageNumbersStayInTheirBox(t *testing.T) {
	hits := []game.Hit{{Value: 123456789, Crit: true, X: 155, Y: 2}, {Value: 7, Crit: true, X: 2, Y: 79}, {Value: 1234567890123456789, Crit: true, X: 80, Y: 40}}
	tt := view(400, 300, &world{clk: newClock()}, func(c *ui.Context) {
		ui.Column(c).Padding(pixel.K(c).D(20)).Children(func() {
			game.DamageNumbers(c, "Hits", &hits, game.DamageNumbersOptions{})
		})
	})
	img := tt.Image()
	// The box spans dots 20–179 across and 20–99 down.
	outside := count(t, img, 0, 0, 200, 20, pixel.Warn) + count(t, img, 180, 0, 20, 120, pixel.Warn) +
		count(t, img, 0, 100, 200, 20, pixel.Warn) + count(t, img, 0, 0, 20, 120, pixel.Warn)
	if n := outside; n != 0 || count(t, img, 20, 20, 160, 80, pixel.Warn) == 0 {
		t.Fatalf("%d crit dots outside the box", n)
	}
}

func TestDamageNumbersShiftInsideWhole(t *testing.T) {
	whole := func(x, y int) int {
		hits := []game.Hit{{Value: 12345, Crit: true, X: x, Y: y}}
		tt := view(400, 200, &world{clk: newClock()}, func(c *ui.Context) { game.DamageNumbers(c, "Hits", &hits, game.DamageNumbersOptions{}) })
		return count(t, tt.Image(), 0, 0, 160, 80, pixel.Warn)
	}
	// Near an edge a number shifts in, keeping its dots.
	if mid := whole(80, 40); mid == 0 || whole(158, 40) != mid || whole(80, 79) != mid || whole(1, 1) != mid {
		t.Fatal("a number near an edge lost dots")
	}
}

func TestDamageNumbersRiseAndExpire(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		hits := []game.Hit{{Value: 12, X: 40, Y: 40}, {Value: 99, Crit: true, X: 110, Y: 40}}
		tt := view(400, 200, wd, func(c *ui.Context) { game.DamageNumbers(c, "Hits", &hits, game.DamageNumbersOptions{}) })
		top := func(s pixel.Slot) int {
			for y := range 80 {
				if count(t, tt.Image(), 0, y, 160, 1, s) > 0 {
					return y
				}
			}
			return -1
		}
		t0, crit := top(pixel.Ink), top(pixel.Warn)
		wd.clk.t = wd.clk.t.Add(400 * time.Millisecond)
		tt.Frame()
		// Halfway through its life it has risen 6 dots.
		if rose := t0 - top(pixel.Ink); t0 < 0 || crit < 0 || (rose == 6) == reduce || rose != 6 && rose != 0 {
			t.Fatalf("reduced %v: the number rose %d from %d", reduce, rose, t0)
		}
		wd.clk.t = wd.clk.t.Add(500 * time.Millisecond)
		tt.Frame()
		if len(hits) != 0 || top(pixel.Ink) != -1 {
			t.Fatalf("reduced %v: %d hits outlived their time", reduce, len(hits))
		}
	}
}
