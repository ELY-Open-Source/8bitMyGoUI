package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ScanlinesOptions struct {
	// Slot shades the rows; zero means Shade.
	Slot pixel.Slot
}

// Scanlines shades every other row of content at half strength.
func Scanlines(c *ui.Context, o ScanlinesOptions, content func()) ui.Element {
	k := pixel.K(c)
	s := slotOr(o.Slot, pixel.Shade)
	return cover(k, ui.Column(c), func(g pixel.Canvas) {
		for y := 1; y < g.H; y += 2 {
			g.Checker(0, y, g.W, 1, s)
		}
	}, content)
}
