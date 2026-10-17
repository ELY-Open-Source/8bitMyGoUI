package layout

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type StatusBarOptions struct{}

// StatusBar shows segments split by sunken dividers; the first grows.
func StatusBar(c *ui.Context, segments []string, o StatusBarOptions) ui.Element {
	k := pixel.K(c)
	return ui.Row(c).FillWidth().Height(k.D(pixel.LineHeight + 3)).AlignItems(ui.Stretch).Role(ui.RoleStatus).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Face)
		g.HLine(0, 0, g.W, pixel.Edge)
		g.HLine(0, 1, g.W, pixel.Lit)
	}).Children(func() {
		for i, s := range segments {
			if i > 0 {
				ui.Box(c).Width(k.D(2)).Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					g.VLine(0, 3, g.H-5, pixel.Shade)
					g.VLine(1, 3, g.H-5, pixel.Lit)
				})
			}
			seg := ui.Box(c).Padding(k.D(2), k.D(6), 0, k.D(6)).Children(func() {
				pixel.Text(c, s, pixel.TextOptions{})
			})
			if i == 0 {
				seg.Grow(1).Clip()
			}
		}
	})
}
