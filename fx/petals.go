package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PetalsOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// Petals fall on a breeze, turning as they go.
func Petals(c *ui.Context, o PetalsOptions) *ui.Element {
	k := pixel.K(c)
	w, h := size("Petals", o.Width, o.Height, 160, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { petals(g, g.Step(8)) })
	})
}

// petals fall a row a step, drifting left every three.
func petals(g pixel.Canvas, step int) {
	for i := range g.W * g.H / 260 {
		y := wrap(rnd(i, 50)+step, g.H+6) - 3
		x := wrap(rnd(i, 51)-step/3+sway(step+rnd(i, 52), 24, 3), g.W+4) - 2
		g.Sprite(x, y, petal, step/2+i, pixel.Ink, 1)
	}
}
