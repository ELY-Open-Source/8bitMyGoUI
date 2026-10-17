package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PauseMenuOptions struct {
	Disabled bool
}

type PauseMenuResult struct {
	Element ui.Element
	// Chosen is RESUME 0, OPTIONS 1, QUIT 2, or -1.
	Chosen int
}

type pauseKey struct{}

// PauseMenu veils the game under PAUSED; Escape resumes.
func PauseMenu(c *ui.Context, paused *bool, o PauseMenuOptions) PauseMenuResult {
	k := pixel.K(c)
	home := ui.Box(c)
	cur := ui.Local(home, pauseKey{}, func() int { return 0 })
	was := *paused
	res := PauseMenuResult{Chosen: -1}
	items := []string{k.Msg("resume"), k.Msg("options"), k.Msg("quit")}
	res.Element = ui.DialogBase(c, paused, func(back, panel ui.Element) {
		deep, light := dark(k)
		dim(k, back, deep, 10)
		title := k.Msg("paused")
		panel.Label(title).Focusable().FocusRing(false).AlignItems(ui.Center).Gap(k.D(10)).Children(func() {
			pixel.Harbor(c, panel, func() {
				ui.Box(c).Size(k.D(pixel.Measure(title)*3), k.D(pixel.LineHeight*3)).Role(ui.RoleText).Label(title).Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					if g.Step(4)%4 < 2 || k.Reduce() {
						g.Text(0, 0, title, light, 3)
					}
				})
				ui.Box(c).Padding(k.D(6), k.D(10)).Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
				}).Children(func() {
					if i := choose(c, title, items, cur, false, o.Disabled); i >= 0 {
						res.Chosen = i
						*paused = i != 0
					}
				})
			})
		})
	})
	if was && !*paused && res.Chosen < 0 {
		res.Chosen = 0
	}
	if !*paused {
		*cur = 0
	}
	if res.Chosen >= 0 {
		slog.Debug("8bit: pause menu", "chosen", items[res.Chosen])
	}
	return res
}
