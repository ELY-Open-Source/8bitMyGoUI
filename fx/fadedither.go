package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type FadeDitherOptions struct {
	// Slot covers the content; zero means the darkest slot.
	Slot pixel.Slot
}

// FadeDither covers content through Bayer levels, or uncovers it.
func FadeDither(c *ui.Context, covered bool, o FadeDitherOptions, content func()) *ui.Element {
	k := pixel.K(c)
	s := slotOr(o.Slot, dark(k))
	e := ui.Column(c)
	level := 2 * k.Toward(c, e, covered, 12, 8)
	return cover(k, e, func(g pixel.Canvas) { g.Veil(0, 0, g.W, g.H, s, level) }, content)
}
