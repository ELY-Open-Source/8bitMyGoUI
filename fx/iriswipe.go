package fx

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type IrisWipeOptions struct {
	// Slot fills outside the iris; zero means the darkest slot.
	Slot pixel.Slot
}

// IrisWipe closes a circle on content's middle, or opens it.
func IrisWipe(c *ui.Context, closed bool, o IrisWipeOptions, content func()) *ui.Element {
	k := pixel.K(c)
	s := slotOr(o.Slot, dark(k))
	e := ui.Column(c)
	step := k.Toward(c, e, closed, 12, 8)
	return cover(k, e, func(g pixel.Canvas) { iris(g, s, step) }, content)
}

// iris fills outside a circle, gone by step 8.
func iris(g pixel.Canvas, s pixel.Slot, step int) {
	if step == 0 {
		return
	}
	cx, cy := g.W/2, g.H/2
	r := int(math.Ceil(math.Hypot(float64(g.W)/2, float64(g.H)/2))) * (8 - step) / 8
	for y := range g.H {
		dy := y - cy
		if dy*dy >= r*r {
			g.HLine(0, y, g.W, s)
			continue
		}
		dx := int(math.Sqrt(float64(r*r - dy*dy)))
		g.HLine(0, y, cx-dx, s)
		g.HLine(cx+dx, y, g.W-cx-dx, s)
	}
}
