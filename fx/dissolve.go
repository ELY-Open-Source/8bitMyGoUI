package fx

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type DissolveOptions struct {
	// Slot replaces the dots; zero means the darkest slot.
	Slot pixel.Slot
}

// Dissolve switches content's dots to a slot in shuffled order.
func Dissolve(c *ui.Context, covered bool, o DissolveOptions, content func()) ui.Element {
	k := pixel.K(c)
	s := slotOr(o.Slot, dark(k))
	e := ui.Column(c)
	step := k.Toward(c, e, covered, 12, 8)
	return cover(k, e, func(g pixel.Canvas) { dissolve(g, s, step) }, content)
}

// dissolve covers dots whose fixed rank, 0–7, is below step.
func dissolve(g pixel.Canvas, s pixel.Slot, step int) {
	on := func(x, y int) bool { return rnd(x+y<<16, 150)%8 < step }
	for y := range g.H {
		for x := 0; x < g.W; x++ {
			if !on(x, y) {
				continue
			}
			end := x + 1
			for end < g.W && on(end, y) {
				end++
			}
			g.HLine(x, y, end-x, s)
			x = end
		}
	}
}
