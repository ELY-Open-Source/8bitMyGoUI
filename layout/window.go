package layout

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

const tabH = pixel.LineHeight + 2

type WindowOptions struct {
	// Title shows in a tab above the body.
	Title string
	// Padding is inside the bevel; zero means 6.
	Padding int
}

// Window is a beveled window with an optional title tab.
func Window(c *ui.Context, o WindowOptions, children func()) *ui.Element {
	k := pixel.K(c)
	pad := 3 + pixel.Or("Window Padding", o.Padding, 6)
	top := 0
	e := ui.Column(c)
	if o.Title != "" {
		top = tabH - 1
		e.Role(ui.RoleGroup).Label(o.Title)
	}
	return e.Padding(k.D(top+pad), k.D(pad), k.D(pad), k.D(pad)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		if g.W < 12 || g.H < top+6 {
			return
		}
		g.Frame(0, top, g.W, g.H-top, pixel.Window, pixel.Face)
		if o.Title == "" {
			return
		}
		w := min(pixel.Measure(o.Title)+12, g.W-8)
		g.Frame(4, 0, w, tabH, pixel.Raised, pixel.Face)
		g.Rect(5, tabH-2, w-2, 4, pixel.Face)
		g.Clip(4, 0, w-2, tabH, func() { g.Text(10, 1, o.Title, pixel.Ink, 1) })
	}).Children(children)
}
