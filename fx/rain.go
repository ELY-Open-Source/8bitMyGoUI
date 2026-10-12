package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type RainOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// Rain slants streaks down; each splashes as it lands.
func Rain(c *ui.Context, o RainOptions) *ui.Element {
	k := pixel.K(c)
	w, h := size("Rain", o.Width, o.Height, 160, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { rain(g, k, g.Step(12)) })
	})
}

// rain drops fall three rows, a dot left, each step.
func rain(g pixel.Canvas, k pixel.Kit, step int) {
	s := pick(k, pixel.Faint, pixel.Info)
	fall := (g.H+4)/3 + 1
	for i := range g.W * g.H / 160 {
		p := wrap(step+rnd(i, 30), fall+3)
		x := wrap(rnd(i, 31)-p, g.W)
		switch {
		case p < fall:
			y := 3*p - 4
			g.Line(x+1, y, x, y+3, s)
		case p == fall:
			g.Dot(x-1, g.H-2, s)
			g.Dot(x+1, g.H-2, s)
		case p == fall+1:
			g.Dot(x-2, g.H-3, s)
			g.Dot(x+2, g.H-3, s)
		}
	}
}
