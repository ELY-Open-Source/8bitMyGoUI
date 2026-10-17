package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type StarfieldOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// Starfield scrolls three layers of stars, near ones fastest.
func Starfield(c *ui.Context, o StarfieldOptions) ui.Element {
	k := pixel.K(c)
	w, h := size("Starfield", o.Width, o.Height, 160, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { starfield(g, k, g.Step(12)) })
	})
}

// starfield moves layer l one dot every 4>>l steps.
func starfield(g pixel.Canvas, k pixel.Kit, step int) {
	g.Rect(0, 0, g.W, g.H, dark(k))
	for l, s := range []pixel.Slot{pixel.Faint, pixel.AccentLit, light(k)} {
		shift := step / (4 >> l)
		for i := range g.W * g.H / (120 + 80*l) {
			x, y := wrap(rnd(i, 20+l)-shift, g.W), rnd(i, 23+l)%g.H
			switch t := wrap(step+rnd(i, 26), 12); {
			case l == 2 && t < 2:
				g.Rect(x-1, y, 3, 1, s)
				g.Rect(x, y-1, 1, 3, s)
			case t != 6:
				g.Dot(x, y, s)
			}
		}
	}
}
