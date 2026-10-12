package feedback

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PulseOptions struct {
	// Slot colors the ring; zero means Accent.
	Slot pixel.Slot
}

// Pulse rings its child with a ring growing in steps.
func Pulse(c *ui.Context, o PulseOptions, content func()) *ui.Element {
	k := pixel.K(c)
	slot := o.Slot
	if slot == 0 {
		slot = pixel.Accent
	}
	return ui.Box(c).AlignSelf(ui.Start).Padding(k.D(6)).DrawOver(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		f := g.Step(6) % 6
		x, level := 5-f, 16-f*3
		w, h := g.W-2*x, g.H-2*x
		g.Veil(x, x, w, 1, slot, level)
		g.Veil(x, x+h-1, w, 1, slot, level)
		g.Veil(x, x, 1, h, slot, level)
		g.Veil(x+w-1, x, 1, h, slot, level)
	}).Children(content)
}
