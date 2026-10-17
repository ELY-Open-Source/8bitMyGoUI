package display

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type IconOptions struct {
	Slot  pixel.Slot
	Scale int
	// Label names the icon; empty hides it.
	Label string
}

// Icon draws a sprite at ×1–×4 in an ink slot.
func Icon(c *ui.Context, sp *pixel.Sprite, o IconOptions) ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	return art(c, k, o.Label, sp.W*sc, sp.H*sc, func(g pixel.Canvas) { g.Sprite(0, 0, sp, 0, o.Slot, sc) })
}
