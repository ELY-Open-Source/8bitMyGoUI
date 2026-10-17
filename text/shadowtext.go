package text

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ShadowTextOptions struct {
	Slot  pixel.Slot
	Scale int
}

// ShadowText drops a 1-dot Shade shadow behind s.
func ShadowText(c *ui.Context, s string, o ShadowTextOptions) ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	return block(c, k, s, pixel.Measure(s)*sc+1, pixel.LineHeight*sc+1, func(g pixel.Canvas) {
		g.Text(1, 1, s, pixel.Shade, sc)
		g.Text(0, 0, s, o.Slot, sc)
	})
}
