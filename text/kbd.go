package text

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type KbdOptions struct {
	Pressed bool
}

// Kbd is a keycap with depth; Pressed sinks it.
func Kbd(c *ui.Context, label string, o KbdOptions) *ui.Element {
	k := pixel.K(c)
	face := pixel.LineHeight + 2
	return block(c, k, label, pixel.Measure(label)+10, face+3, func(g pixel.Canvas) {
		top := 0
		if o.Pressed {
			top = 2
		}
		g.Frame(0, 2, g.W, face+1, pixel.Flat, pixel.Shade)
		g.Frame(0, top, g.W, face, pixel.Raised, pixel.Face)
		g.Text(5, top+1, label, pixel.Ink, 1)
	})
}
