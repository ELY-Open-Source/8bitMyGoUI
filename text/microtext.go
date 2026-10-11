package text

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MicroTextOptions struct {
	Slot  pixel.Slot
	Scale int
}

// MicroText draws s upper-cased in micro font; unknown runes panic.
func MicroText(c *ui.Context, s string, o MicroTextOptions) *ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	up := strings.ToUpper(s)
	return block(c, k, s, pixel.MicroWidth(up)*sc, pixel.MicroHeight*sc, func(g pixel.Canvas) {
		g.Micro(0, 0, up, o.Slot, sc)
	})
}
