package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/fx"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var screenFxEntries = []Entry{
	{309, "Scanlines", "扫描线", "fx", "Every other row, half shaded. Instant CRT.", demoScanlines},
	{310, "Vignette", "暗角", "fx", "The corners dither into the dark.", demoVignette},
	{311, "ScreenShake", "震屏", "fx", "Something big just landed.", demoScreenShake},
	{312, "Flash", "闪白", "fx", "A step of light, then half a step. Never under Reduce Motion.", demoFlash},
	{313, "FadeDither", "抖动淡出", "fx", "Fades through Bayer levels, never through opacity.", demoFadeDither},
	{314, "IrisWipe", "光圈转场", "fx", "The circle closes on the middle. That's all, folks.", demoIrisWipe},
	{315, "Dissolve", "溶解", "fx", "Dots give up one batch at a time, in shuffled order.", demoDissolve},
	{316, "Mosaic", "马赛克", "fx", "Blocks grow to eight dots, then the picture comes back.", demoMosaic},
	{317, "Glitch", "故障", "fx", "Rows slip and two colors trade places. Blow on the cartridge.", demoGlitch},
	{318, "Sparkle", "闪光", "fx", "Four stars take turns around a shiny thing.", demoSparkle},
	{319, "Explosion", "爆炸", "fx", "Fireball, smoke, debris. Stand back.", demoExplosion},
	{320, "Confetti", "彩纸", "fx", "Bits of every color tumble down. You won!", demoConfetti},
}

// stage is a little world in the palette's season.
func stage(c *ui.Context) {
	p := pixel.K(c).Palette()
	fx.SeasonScene(c, p.Season, hourOf(p), fx.SeasonSceneOptions{Width: 128, Height: 72, Label: "stage"})
}

func demoScanlines(c *ui.Context, log func(string)) {
	ui.Column(c).AlignItems(ui.Start).Children(func() {
		fx.Scanlines(c, fx.ScanlinesOptions{}, func() { stage(c) })
	})
}

func demoVignette(c *ui.Context, log func(string)) {
	ui.Column(c).AlignItems(ui.Start).Children(func() {
		fx.Vignette(c, fx.VignetteOptions{}, func() { stage(c) })
	})
}

func demoScreenShake(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := fxState(c)
	ui.Column(c).Gap(k.D(6)).AlignItems(ui.Start).Children(func() {
		fx.ScreenShake(c, fx.ScreenShakeOptions{Cue: s.shake}, func() { stage(c) })
		if small(c, "SHAKE", false) {
			s.shake++
			log("rumble!")
		}
	})
}

func demoFlash(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := fxState(c)
	ui.Column(c).Gap(k.D(6)).AlignItems(ui.Start).Children(func() {
		fx.Flash(c, fx.FlashOptions{Cue: s.flash}, func() { stage(c) })
		if small(c, "FLASH", false) {
			s.flash++
			log("flash")
		}
	})
}

// toggled shows content under a transition, with its toggle.
func toggled(c *ui.Context, text string, on *bool, log func(string), effect func(content func())) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(6)).AlignItems(ui.Start).Children(func() {
		effect(func() { stage(c) })
		if input.Toggle(c, text, on, input.ToggleOptions{Size: input.Small}).Changed {
			log(fmt.Sprintf("%s %v", text, *on))
		}
	})
}

func demoFadeDither(c *ui.Context, log func(string)) {
	s := fxState(c)
	toggled(c, "FADE OUT", &s.fade, log, func(content func()) { fx.FadeDither(c, s.fade, fx.FadeDitherOptions{}, content) })
}

func demoIrisWipe(c *ui.Context, log func(string)) {
	s := fxState(c)
	toggled(c, "CLOSE", &s.iris, log, func(content func()) { fx.IrisWipe(c, s.iris, fx.IrisWipeOptions{}, content) })
}

func demoDissolve(c *ui.Context, log func(string)) {
	s := fxState(c)
	toggled(c, "DISSOLVE", &s.dissolve, log, func(content func()) { fx.Dissolve(c, s.dissolve, fx.DissolveOptions{}, content) })
}

func demoMosaic(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := fxState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
			fx.Mosaic(c, display.Reduce(sunset, 48, 32, k.Palette()), fx.MosaicOptions{Cue: s.mosaic, Scale: 2, Label: "sunset"})
			fx.Mosaic(c, slime, fx.MosaicOptions{Cue: s.mosaic, Scale: 3, Label: "slime"})
		})
		if small(c, "PIXELATE", false) {
			s.mosaic++
			log("pixelate")
		}
	})
}

func demoGlitch(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := fxState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
			fx.Glitch(c, display.Reduce(sunset, 48, 32, k.Palette()), fx.GlitchOptions{Cue: s.glitch, Scale: 2, Label: "sunset"})
			fx.Glitch(c, slime, fx.GlitchOptions{Cue: s.glitch, Scale: 3, Label: "slime"})
		})
		if small(c, "GLITCH", false) {
			s.glitch++
			log("glitch")
		}
	})
}

func demoSparkle(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	layout.Layers(c, func() {
		ui.Box(c).Size(k.D(48), k.D(48)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
			display.Gem(c, display.Legendary, display.GemOptions{Scale: 2})
		})
	}, func() { fx.Sparkle(c, fx.SparkleOptions{Size: 48, Label: "sparkles"}) })
}

func demoExplosion(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := fxState(c)
	ui.Column(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		fx.Explosion(c, fx.ExplosionOptions{Cue: s.boom, Label: "blast"})
		if small(c, "BOOM", false) {
			s.boom++
			log("boom")
		}
	})
}

func demoConfetti(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	layout.Layers(c, func() {
		ui.Box(c).Size(k.D(140), k.D(64)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
			pixel.Text(c, "YOU WIN!", pixel.TextOptions{Scale: 2, Slot: pixel.Accent})
		})
	}, func() { fx.Confetti(c, fx.ConfettiOptions{Width: 140, Height: 64, Label: "confetti"}) })
}
