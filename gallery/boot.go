package gallery

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type bootKey struct{}

// boot assembles the logo dot by dot, then waits.
func (a *App) boot(c *ui.Context) {
	k := pixel.K(c)
	step := k.Tick(c, 12)
	first := ui.Local(c.Root(), bootKey{}, func() int { return step })
	shown := min(16, step-*first)
	if k.Reduce() {
		shown = 16
	}
	start := shown == 16 && (k.Reduce() || step%12 < 8)
	shell, base, _ := pixel.Shell(c, func() *ui.Element { return ui.ButtonBase(c).Label("Press start") })
	shell.Fill().Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		w := pixel.Measure(logo) * 4
		x, y := (g.W-w)/2, g.H/2-40
		shadow(g, x, y, logo, pixel.Accent, 4)
		g.Veil(x, y, w+4, 68, pixel.Back, 16-shown)
		if start {
			s := "PRESS START"
			g.Text((g.W-pixel.Measure(s)*2)/2, y+84, s, pixel.Ink, 2)
		}
		hint := "Enter · Space · Click     回车 · 空格 · 点击"
		g.Text((g.W-pixel.Measure(hint))/2, g.H-28, hint, pixel.Faint, 1)
	})
	base.AutoFocus()
	if base.Clicked() || c.Shortcut(0, ui.KeyEnter) || c.Shortcut(0, ui.KeySpace) {
		slog.Debug("8bit: gallery start")
		a.started = true
		c.Invalidate()
	}
}
