package chat

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TypingIndicatorOptions struct {
	// Slot inks the dots; zero means Faint.
	Slot pixel.Slot
}

// TypingIndicator bounces three dots in a bubble while name types.
func TypingIndicator(c *ui.Context, name string, o TypingIndicatorOptions) *ui.Element {
	k := pixel.K(c)
	ink := o.Slot
	if ink == pixel.Ink {
		ink = pixel.Faint
	}
	return display.Bubble(c, display.BubbleOptions{Tail: display.Left, At: 4}, func() {
		ui.Box(c).Role(ui.RoleStatus).Label(fmt.Sprintf(k.Msg("typing"), name)).Size(k.D(17), k.D(7)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			// Each dot hops in turn, then a rest step.
			hop := g.Step(6) % 4
			for i := range 3 {
				y := 3
				if hop == i+1 {
					y = 0
				}
				g.Rect(i*7, y, 3, 3, ink)
			}
		})
	}).AlignSelf(ui.Start)
}
