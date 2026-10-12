package tabletop

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DiceTrayOptions struct {
	// Cue rolls them all, as Dice's Cue does.
	Cue int
}

// DiceTray holds several dice on felt, and their total.
func DiceTray(c *ui.Context, text string, values []int, o DiceTrayOptions) *ui.Element {
	k := pixel.K(c)
	if len(values) == 0 {
		panic("tabletop: DiceTray has no dice")
	}
	total := 0
	for _, v := range values {
		total += v
	}
	sum := fmt.Sprintf("%s %d", k.Msg("total"), total)
	return ui.Column(c).Gap(k.D(4)).AlignItems(ui.Center).AlignSelf(ui.Start).Role(ui.RoleGroup).Label(text).Description(sum).Children(func() {
		ui.Row(c).Gap(k.D(2)).Padding(k.D(6)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Good)
			g.Veil(2, 2, g.W-4, g.H-4, pixel.Shade, 4)
		}).Children(func() {
			for i, v := range values {
				Dice(c, fmt.Sprintf("%s %d", text, i+1), v, DiceOptions{Cue: o.Cue})
			}
		})
		pixel.Text(c, sum, pixel.TextOptions{})
	})
}
