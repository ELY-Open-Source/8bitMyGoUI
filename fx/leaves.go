package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LeavesOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// leafSlots color the leaves in turn.
var leafSlots = []pixel.Slot{pixel.Warn, pixel.Bad, pixel.Accent}

// Leaves swing down side to side, turning on their stems.
func Leaves(c *ui.Context, o LeavesOptions) ui.Element {
	k := pixel.K(c)
	w, h := size("Leaves", o.Width, o.Height, 160, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { leaves(g, g.Step(8)) })
	})
}

func leaves(g pixel.Canvas, step int) {
	for i := range max(2, g.W*g.H/500) {
		y := wrap(rnd(i, 60)+step, g.H+10) - 5
		x := wrap(rnd(i, 61)+step/4+sway(step+rnd(i, 62), 20, 5), g.W+10) - 5
		g.Sprite(x, y, leaf, (step+i)/3, leafSlots[i%len(leafSlots)], 1)
	}
}
