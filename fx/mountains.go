package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MountainsOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// Mountains scroll three ridges: far ones slow, near ones quick.
func Mountains(c *ui.Context, o MountainsOptions) *ui.Element {
	k := pixel.K(c)
	w, h := size("Mountains", o.Width, o.Height, 160, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { mountains(g, k, g.Step(8), 3) })
	})
}

// tri rises from 0 to a, then falls, over p.
func tri(t, p, a int) int {
	t, h := wrap(t, p), p/2
	if t < h {
		return a * t / h
	}
	return a * (p - t) / h
}

// mountains draws n ridges; ridge l moves every 8>>l steps.
func mountains(g pixel.Canvas, k pixel.Kit, step, n int) {
	slots := [3]pixel.Slot{pick(k, pixel.Shade, pixel.Face), pick(k, pixel.Faint, pixel.Back), pick(k, pixel.Edge, pixel.Well)}
	snowy := k.Palette().Season == pixel.Winter
	for l := range n {
		shift := step / (8 >> l)
		for x := range g.W {
			u := x + shift
			h := g.H*(7-2*l)/16 + tri(u+rnd(l, 100), 48+16*l, g.H/5) + tri(2*u+rnd(l, 101), 37, g.H/12) + rnd(u, 102+l)%2
			g.VLine(x, g.H-h, h, slots[l])
			// Winter caps the far peaks with snow.
			if snowy && l == 0 && h > g.H*9/16 {
				g.VLine(x, g.H-h, 2, light(k))
			}
		}
	}
}
