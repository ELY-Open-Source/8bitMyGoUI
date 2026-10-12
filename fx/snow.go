package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SnowOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// Snow drifts flakes down, swaying; near ones are big, quick.
func Snow(c *ui.Context, o SnowOptions) *ui.Element {
	k := pixel.K(c)
	w, h := size("Snow", o.Width, o.Height, 160, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { snow(g, k, g.Step(8)) })
	})
}

func snow(g pixel.Canvas, k pixel.Kit, step int) {
	s := light(k)
	for i := range g.W * g.H / 110 {
		near := rnd(i, 40)%3 == 0
		speed := 1
		if near {
			speed = 2
		}
		y := wrap(rnd(i, 41)+step*speed, g.H+4) - 2
		x := wrap(rnd(i, 42)+sway(step+rnd(i, 43), 16, 2), g.W)
		if near {
			g.Rect(x, y, 2, 2, s)
		} else {
			g.Dot(x, y, s)
		}
	}
}
