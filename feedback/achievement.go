package feedback

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type AchievementOptions struct {
	// Cue shows the card each time it changes.
	Cue int
	// Width is in dots; zero means 200.
	Width int
}

// Achievement slides a trophy card in, holds it, then out.
func Achievement(c *ui.Context, title string, o AchievementOptions) ui.Element {
	k := pixel.K(c)
	w := pixel.Or("Achievement Width", o.Width, 200)
	e := ui.Box(c)
	// Ten steps a second; in by 5, gone by 40.
	step := k.Since(c, k.Cue(e, o.Cue), 10, 40)
	if step >= 40 {
		return e
	}
	out := 0
	switch {
	case k.Reduce():
	case step < 5:
		out = 5 - step
	case step >= 35:
		out = step - 34
	}
	head := k.Msg("achievement")
	return e.Role(ui.RoleStatus).Label(head+" "+title).Size(k.D(w), k.D(40)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		x := out * w / 5
		g.Clip(0, 0, g.W, g.H, func() {
			g.Frame(x, 0, w, 40, pixel.Window, pixel.Face)
			g.Sprite(x+6, 8, icon.Trophy, 0, pixel.Warn, 2)
			g.Text(x+36, 4, head, pixel.Faint, 1)
			g.Text(x+36, 20, title, pixel.Ink, 1)
		})
	})
}
