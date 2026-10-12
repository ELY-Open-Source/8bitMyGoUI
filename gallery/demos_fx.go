package gallery

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/fx"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var fxEntries = []Entry{
	{291, "PaletteCycle", "调色板循环", "fx", "Four slots take turns, so the rings seem to flow.", demoPaletteCycle},
	{292, "Sky", "天空", "fx", "Blue by day; by night, stars that blink.", demoSky},
	{293, "SunMoon", "日月", "fx", "Up at six, down at six. The moon takes the night shift.", demoSunMoon},
	{294, "Clouds", "云", "fx", "Far clouds dawdle, near clouds hurry.", demoClouds},
	{295, "Starfield", "星空", "fx", "Three layers of stars rush past. Hyperspace at 12 fps.", demoStarfield},
	{296, "Rain", "雨", "fx", "Slanted streaks, each landing with a tiny splash.", demoRain},
	{297, "Snow", "雪", "fx", "Flakes sway down; near ones are big and quick.", demoSnow},
	{298, "Petals", "樱花瓣", "fx", "Petals ride the breeze, turning as they fall.", demoPetals},
	{299, "Leaves", "落叶", "fx", "Leaves swing side to side on the way down.", demoLeaves},
	{300, "Fireflies", "萤火虫", "fx", "Little lights wander, glow, and rest.", demoFireflies},
	{301, "Waterfall", "瀑布", "fx", "Four slots cycling down make the water fall.", demoWaterfallFx},
	{302, "Ocean", "海浪", "fx", "Wave bands roll in by cycling, with glints on top.", demoOcean},
	{303, "Fire", "火焰", "fx", "Three flicker frames, heat rising through them.", demoFire},
	{304, "Lava", "熔岩", "fx", "Hot bands flow; bubbles swell and pop. Do not touch.", demoLava},
	{305, "Aurora", "极光", "fx", "Curtains of light, their colors sliding sideways.", demoAurora},
	{306, "Lightning", "闪电", "fx", "A flash, then a forked bolt. No two strikes alike.", demoLightning},
	{307, "Mountains", "山脉", "fx", "Three ridges scroll at three speeds. Parallax!", demoMountains},
	{308, "SeasonScene", "四季场景", "fx", "Pick a season and an hour; the world repaints itself.", demoSeasonScene},
}

type fxDemo struct {
	season, hour                               int
	strike, shake, flash, mosaic, glitch, boom int
	fade, iris, dissolve                       bool
}

func fxState(c *ui.Context) *fxDemo {
	p := pixel.K(c).Palette()
	return ui.Local(c.Root(), demoKey("fx"), func() fxDemo { return fxDemo{season: int(p.Season), hour: hourOf(p)} })
}

// hourOf suits p: morning by day, late evening by night.
func hourOf(p *pixel.Palette) int {
	if p.Night {
		return 21
	}
	return 10
}

// paletteOf is the settings in force with season's palette.
func paletteOf(c *ui.Context, s pixel.Season, night bool) pixel.Settings {
	set := pixel.Current(c)
	set.Palette = pixel.Find(s, night)
	return set
}

// over stacks an effect on a same-sized sky.
func over(c *ui.Context, w, h int, effect func()) {
	layout.Layers(c, layout.LayersOptions{}, func() { fx.Sky(c, fx.SkyOptions{Width: w, Height: h}) }, effect)
}

// portalBand rings the portal, outside in.
var portalBand = []pixel.Slot{pixel.Accent, pixel.AccentLit, pixel.Magic, pixel.Info}

var portal = pixel.SpriteOf(15, 15, func(x, y int) (pixel.Slot, bool) {
	d := math.Hypot(float64(x-7), float64(y-7))
	return portalBand[int(d)%len(portalBand)], d < 7.5
})

func demoPaletteCycle(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		fx.PaletteCycle(c, portal, portalBand, fx.PaletteCycleOptions{Scale: 3, Label: "portal"})
		fx.PaletteCycle(c, portal, portalBand, fx.PaletteCycleOptions{Scale: 2, FPS: 4, Label: "slow portal"})
	})
}

func demoSky(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(4)).Children(func() {
		for _, night := range []bool{false, true} {
			pixel.Scope(c, paletteOf(c, k.Palette().Season, night), func() {
				fx.Sky(c, fx.SkyOptions{Width: 68, Height: 64, Label: map[bool]string{false: "day sky", true: "night sky"}[night]})
			})
		}
	})
}

func demoSunMoon(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := fxState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		pixel.Scope(c, paletteOf(c, k.Palette().Season, s.hour < 6 || s.hour >= 18), func() {
			over(c, 140, 72, func() {
				fx.SunMoon(c, s.hour, fx.SunMoonOptions{Width: 140, Height: 72, Label: fmt.Sprintf("%02d:00", s.hour)})
			})
		})
		if input.Stepper(c, "Hour", &s.hour, input.StepperOptions{Max: 23}).Changed {
			log(fmt.Sprintf("hour %02d", s.hour))
		}
	})
}

func demoClouds(c *ui.Context, log func(string)) {
	over(c, 140, 80, func() { fx.Clouds(c, fx.CloudsOptions{Width: 140, Height: 80, Label: "clouds"}) })
}

func demoStarfield(c *ui.Context, log func(string)) {
	fx.Starfield(c, fx.StarfieldOptions{Width: 140, Height: 80, Label: "stars"})
}

func demoRain(c *ui.Context, log func(string)) {
	over(c, 140, 80, func() { fx.Rain(c, fx.RainOptions{Width: 140, Height: 80, Label: "rain"}) })
}

func demoSnow(c *ui.Context, log func(string)) {
	over(c, 140, 80, func() { fx.Snow(c, fx.SnowOptions{Width: 140, Height: 80, Label: "snow"}) })
}

func demoPetals(c *ui.Context, log func(string)) {
	over(c, 140, 80, func() { fx.Petals(c, fx.PetalsOptions{Width: 140, Height: 80, Label: "petals"}) })
}

func demoLeaves(c *ui.Context, log func(string)) {
	over(c, 140, 80, func() { fx.Leaves(c, fx.LeavesOptions{Width: 140, Height: 80, Label: "leaves"}) })
}

func demoFireflies(c *ui.Context, log func(string)) {
	pixel.Scope(c, paletteOf(c, pixel.K(c).Palette().Season, true), func() {
		over(c, 140, 80, func() { fx.Fireflies(c, fx.FirefliesOptions{Width: 140, Height: 80, Label: "fireflies"}) })
	})
}

func demoWaterfallFx(c *ui.Context, log func(string)) {
	fx.Waterfall(c, fx.WaterfallOptions{Label: "waterfall"})
}

func demoOcean(c *ui.Context, log func(string)) {
	fx.Ocean(c, fx.OceanOptions{Width: 140, Label: "ocean"})
}

func demoFire(c *ui.Context, log func(string)) {
	fx.Fire(c, fx.FireOptions{Scale: 3, Label: "campfire"})
}

func demoLava(c *ui.Context, log func(string)) {
	fx.Lava(c, fx.LavaOptions{Width: 140, Label: "lava"})
}

func demoAurora(c *ui.Context, log func(string)) {
	pixel.Scope(c, paletteOf(c, pixel.K(c).Palette().Season, true), func() {
		layout.Layers(c, layout.LayersOptions{}, func() { fx.Sky(c, fx.SkyOptions{Width: 140, Height: 80}) },
			func() { fx.Aurora(c, fx.AuroraOptions{Width: 140, Height: 80, Label: "aurora"}) },
			func() { fx.Mountains(c, fx.MountainsOptions{Width: 140, Height: 80}) })
	})
}

func demoLightning(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := fxState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		pixel.Scope(c, paletteOf(c, k.Palette().Season, true), func() {
			over(c, 120, 80, func() { fx.Lightning(c, fx.LightningOptions{Cue: s.strike, Width: 120, Height: 80, Label: "storm"}) })
		})
		if small(c, "STRIKE", false) {
			s.strike++
			log(fmt.Sprintf("strike %d", s.strike))
		}
	})
}

func demoMountains(c *ui.Context, log func(string)) {
	over(c, 140, 80, func() { fx.Mountains(c, fx.MountainsOptions{Width: 140, Height: 80, Label: "mountains"}) })
}

func demoSeasonScene(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := fxState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		fx.SeasonScene(c, pixel.Season(s.season), s.hour, fx.SeasonSceneOptions{Width: 144, Height: 80, Label: "the world"})
		seasons := []string{"Spring", "Summer", "Autumn", "Winter"}
		if input.Segmented(c, "Season", seasons, &s.season, input.SegmentedOptions{}).Changed {
			log("season " + seasons[s.season])
		}
		if input.Stepper(c, "Hour", &s.hour, input.StepperOptions{Max: 23}).Changed {
			log(fmt.Sprintf("hour %02d", s.hour))
		}
	})
}
