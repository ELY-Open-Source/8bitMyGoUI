package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type AuroraOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// Aurora hangs curtains whose slots slide sideways.
func Aurora(c *ui.Context, o AuroraOptions) ui.Element {
	k := pixel.K(c)
	w, h := size("Aurora", o.Width, o.Height, 160, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { aurora(g, g.Step(8)) })
	})
}

// aurora's columns ripple on top and fade at the hem.
func aurora(g pixel.Canvas, step int) {
	band := pixel.Cycle([]pixel.Slot{pixel.Good, pixel.Info, pixel.Magic, pixel.Good, pixel.AccentLit}, step)
	for x := range g.W {
		top := g.H/8 + tri(x/2+step, 64, g.H/6) + tri(x+2*step, 36, 3)
		long := g.H/4 + tri(3*x/2-step, 50, g.H/6)
		s := band[x/4%len(band)]
		g.VLine(x, top, long*2/3, s)
		g.Veil(x, top+long*2/3, 1, long/3, s, 6)
	}
}
