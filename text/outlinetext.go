package text

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type OutlineTextOptions struct {
	// Slot is the fill; zero means Lit.
	Slot  pixel.Slot
	Scale int
}

// OutlineText rings every glyph dot of s with Edge.
func OutlineText(c *ui.Context, s string, o OutlineTextOptions) *ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	fill := slotOr(o.Slot, pixel.Lit)
	return block(c, k, s, pixel.Measure(s)*sc+2, pixel.LineHeight*sc+2, func(g pixel.Canvas) {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				g.Text(1+dx, 1+dy, s, pixel.Edge, sc)
			}
		}
		g.Text(1, 1, s, fill, sc)
	})
}
