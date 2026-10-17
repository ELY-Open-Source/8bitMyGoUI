package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type FlashOptions struct {
	// Cue flashes each time it changes.
	Cue int
	// Slot fills the flash; zero means Lit.
	Slot pixel.Slot
}

// Flash fills content with a slot one step, then half.
func Flash(c *ui.Context, o FlashOptions, content func()) ui.Element {
	k := pixel.K(c)
	s := slotOr(o.Slot, pixel.Lit)
	e := ui.Column(c)
	// Reduce Motion never flashes.
	step := k.Shot(c, e, o.Cue, 10, 2)
	return cover(k, e, func(g pixel.Canvas) {
		switch step {
		case 0:
			g.Rect(0, 0, g.W, g.H, s)
		case 1:
			g.Veil(0, 0, g.W, g.H, s, 8)
		}
	}, content)
}
