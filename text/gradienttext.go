package text

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type GradientTextOptions struct {
	// Top is zero for AccentLit; Bottom is zero for Accent.
	Top, Bottom pixel.Slot
	Scale       int
}

// GradientText dithers glyph dots from Top down to Bottom.
func GradientText(c *ui.Context, s string, o GradientTextOptions) *ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	top, bottom := slotOr(o.Top, pixel.AccentLit), slotOr(o.Bottom, pixel.Accent)
	return block(c, k, s, pixel.Measure(s)*sc, pixel.LineHeight*sc, func(g pixel.Canvas) {
		g.TextBlend(0, 0, s, top, bottom, sc)
	})
}
