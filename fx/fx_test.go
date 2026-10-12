package fx_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/fx"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func newClock() *clock { return &clock{time.Unix(1_800_000_000, 0)} }

// step is the clock's frame at fps, as pixel counts.
func (c *clock) step(fps int64) int { return int(c.t.UnixNano() / (int64(time.Second) / fps)) }

// tick is one step at fps, just past the edge.
func tick(fps int) time.Duration { return time.Second/time.Duration(fps) + time.Microsecond }

// view builds w×h dots on a fake clock, 2px each.
func view(w, h int, pal *pixel.Palette, clk *clock, reduce bool, build func(c *ui.Context)) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pal, Now: clk.now})
		build(c)
	}, 2*w, 2*h)
	if reduce {
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tt.Frame()
	}
	return tt
}

// dots reads w×h dots as slots; -1 is off-palette.
func dots(img *image.RGBA, pal *pixel.Palette, w, h int) [][]int {
	out := make([][]int, h)
	for y := range h {
		out[y] = make([]int, w)
		for x := range w {
			c, s := img.RGBAAt(2*x, 2*y), -1
			for i, pc := range pal.Colors {
				if pc.R == c.R && pc.G == c.G && pc.B == c.B {
					s = i
					break
				}
			}
			out[y][x] = s
		}
	}
	return out
}

// count is how many dots of a box hold s.
func count(d [][]int, x, y, w, h int, s pixel.Slot) int {
	n := 0
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			if d[j][i] == int(s) {
				n++
			}
		}
	}
	return n
}

func same(a, b [][]int) bool {
	for y := range a {
		for x := range a[y] {
			if a[y][x] != b[y][x] {
				return false
			}
		}
	}
	return true
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

// solid is a one-slot w×h box under effects.
func solid(c *ui.Context, w, h int, s pixel.Slot) {
	k := pixel.K(c)
	ui.Box(c).Size(k.D(w), k.D(h)).Label("content").Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, s)
	})
}

var portal = pixel.MustSprite("agcagc\ncagcag")

func TestAmbientEffectsMoveButHoldStillUnderReduceMotion(t *testing.T) {
	for name, build := range map[string]func(c *ui.Context){
		"PaletteCycle": func(c *ui.Context) {
			fx.PaletteCycle(c, portal, []pixel.Slot{pixel.Accent, pixel.Good, pixel.Info}, fx.PaletteCycleOptions{})
		},
		"SunMoon":    func(c *ui.Context) { fx.SunMoon(c, 12, fx.SunMoonOptions{Width: 80, Height: 48}) },
		"Clouds":     func(c *ui.Context) { fx.Clouds(c, fx.CloudsOptions{Width: 80, Height: 48}) },
		"Starfield":  func(c *ui.Context) { fx.Starfield(c, fx.StarfieldOptions{Width: 80, Height: 48}) },
		"Rain":       func(c *ui.Context) { fx.Rain(c, fx.RainOptions{Width: 80, Height: 48}) },
		"Snow":       func(c *ui.Context) { fx.Snow(c, fx.SnowOptions{Width: 80, Height: 48}) },
		"Petals":     func(c *ui.Context) { fx.Petals(c, fx.PetalsOptions{Width: 80, Height: 48}) },
		"Leaves":     func(c *ui.Context) { fx.Leaves(c, fx.LeavesOptions{Width: 80, Height: 48}) },
		"Fireflies":  func(c *ui.Context) { fx.Fireflies(c, fx.FirefliesOptions{Width: 80, Height: 48}) },
		"Waterfall":  func(c *ui.Context) { fx.Waterfall(c, fx.WaterfallOptions{Width: 40, Height: 48}) },
		"Ocean":      func(c *ui.Context) { fx.Ocean(c, fx.OceanOptions{Width: 80, Height: 48}) },
		"Fire":       func(c *ui.Context) { fx.Fire(c, fx.FireOptions{}) },
		"Lava":       func(c *ui.Context) { fx.Lava(c, fx.LavaOptions{Width: 80, Height: 48}) },
		"Aurora":     func(c *ui.Context) { fx.Aurora(c, fx.AuroraOptions{Width: 80, Height: 48}) },
		"Mountains":  func(c *ui.Context) { fx.Mountains(c, fx.MountainsOptions{Width: 80, Height: 48}) },
		"Season":     func(c *ui.Context) { fx.SeasonScene(c, pixel.Winter, 22, fx.SeasonSceneOptions{Width: 80, Height: 48}) },
		"Sparkle":    func(c *ui.Context) { fx.Sparkle(c, fx.SparkleOptions{}) },
		"Confetti":   func(c *ui.Context) { fx.Confetti(c, fx.ConfettiOptions{Width: 80, Height: 48}) },
		"NightSky":   func(c *ui.Context) { fx.Sky(c, fx.SkyOptions{Width: 80, Height: 48}) },
		"NightSnows": func(c *ui.Context) { fx.Snow(c, fx.SnowOptions{Width: 80, Height: 48}) },
	} {
		// A day sky holds still: only stars move.
		pal := pixel.SpringDay
		if name == "NightSky" || name == "NightSnows" || name == "Season" {
			pal = pixel.WinterNight
		}
		for _, reduce := range []bool{false, true} {
			clk := newClock()
			tt := view(80, 48, pal, clk, reduce, build)
			first := dots(tt.Image(), pal, 80, 48)
			moved := false
			for _, d := range []time.Duration{300, 700, 1300} {
				clk.t = clk.t.Add(d * time.Millisecond)
				tt.Frame()
				moved = moved || !same(first, dots(tt.Image(), pal, 80, 48))
			}
			if moved == reduce {
				t.Errorf("%s, reduce %v: moved %v", name, reduce, moved)
			}
		}
	}
}
