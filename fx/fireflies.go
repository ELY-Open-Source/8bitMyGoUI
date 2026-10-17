package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type FirefliesOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// Fireflies wander and glow: bright, dim, then dark a while.
func Fireflies(c *ui.Context, o FirefliesOptions) ui.Element {
	k := pixel.K(c)
	w, h := size("Fireflies", o.Width, o.Height, 160, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { fireflies(g, g.Step(8)) })
	})
}

func fireflies(g pixel.Canvas, step int) {
	for i := range max(3, g.W*g.H/600) {
		x := rnd(i, 70)%g.W + sway(step+rnd(i, 71), 40+rnd(i, 72)%24, 8)
		y := rnd(i, 73)%g.H + sway(step+rnd(i, 74), 28+rnd(i, 75)%20, 5)
		switch t := wrap(step+rnd(i, 76), 16); {
		case t < 6:
			g.Rect(x-1, y, 3, 1, pixel.Good)
			g.Rect(x, y-1, 1, 3, pixel.Good)
			g.Dot(x, y, pixel.Warn)
		case t < 9:
			g.Dot(x, y, pixel.Good)
		}
	}
}
