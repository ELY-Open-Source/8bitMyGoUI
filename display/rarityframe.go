package display

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type RarityFrameOptions struct {
	// Padding is inside the border; zero means 4.
	Padding int
}

// RarityFrame borders children by rarity; legendary carries a glint.
func RarityFrame(c *ui.Context, r Rarity, o RarityFrameOptions, children func()) *ui.Element {
	k := pixel.K(c)
	ink := r.Slot()
	pad := 3 + pixel.Or("RarityFrame Padding", o.Padding, 4)
	return ui.Column(c).Padding(k.D(pad)).Role(ui.RoleGroup).Label(r.name(k)).Draw(func(p *ui.Painter, r2 ui.Rect) {
		g := k.Canvas(p, r2)
		if g.W < 6 || g.H < 6 {
			return
		}
		g.Frame(0, 0, g.W, g.H, pixel.Flat, ink)
		g.Rect(3, 3, g.W-6, g.H-6, pixel.Face)
		g.Rect(2, 2, g.W-4, 1, pixel.Edge)
		g.Rect(2, g.H-3, g.W-4, 1, pixel.Edge)
		g.Rect(2, 2, 1, g.H-4, pixel.Edge)
		g.Rect(g.W-3, 2, 1, g.H-4, pixel.Edge)
		if r >= Epic {
			for _, at := range [][2]int{{0, 0}, {g.W - 2, 0}, {0, g.H - 2}, {g.W - 2, g.H - 2}} {
				g.Rect(at[0], at[1], 2, 2, pixel.Lit)
			}
		}
	}).DrawOver(func(p *ui.Painter, r2 ui.Rect) {
		if r != Legendary {
			return
		}
		g := k.Canvas(p, r2)
		if g.W < 6 || g.H < 6 {
			return
		}
		per := 2*(g.W+g.H) - 4
		at := g.Step(20) * 2 % per
		for i := range 6 {
			x, y := around(g.W, g.H, (at+i)%per)
			g.Dot(x, y, pixel.Lit)
		}
	}).Children(children)
}

// around is the edge dot d steps clockwise from 0,0.
func around(w, h, d int) (int, int) {
	switch {
	case d < w:
		return d, 0
	case d < w+h-1:
		return w - 1, d - w + 1
	case d < 2*w+h-2:
		return w - 1 - (d - w - h + 2), h - 1
	default:
		return 0, h - 1 - (d - 2*w - h + 3)
	}
}
