package text

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CodeOptions struct {
	Slot pixel.Slot
}

// Code is an inline chip on Well, bordered in Edge.
func Code(c *ui.Context, s string, o CodeOptions) ui.Element {
	k := pixel.K(c)
	return block(c, k, s, pixel.Measure(s)+6, pixel.LineHeight+2, func(g pixel.Canvas) {
		g.Frame(0, 0, g.W, g.H, pixel.Flat, pixel.Well)
		g.Text(3, 1, s, o.Slot, 1)
	})
}
