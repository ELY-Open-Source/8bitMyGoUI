package feedback

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var (
	bang = pixel.MustSprite(`
		..###..
		.##.##.
		###.###
		###.###
		#######
		###.###
		.#####.
	`)
	tick = pixel.MustSprite(`
		......#
		.....##
		#...##.
		##.##..
		.###...
		..#....
		.......
	`)
)

type InlineMessageOptions struct {
	Kind Kind
}

// InlineMessage is a small ✓ or ! field note.
func InlineMessage(c *ui.Context, message string, o InlineMessageOptions) *ui.Element {
	k := pixel.K(c)
	slot, _ := o.Kind.look()
	sp := bang
	if o.Kind == Success {
		sp = tick
	}
	return ui.Row(c).Gap(k.D(3)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(message).Children(func() {
		ui.Box(c).Size(k.D(sp.W), k.D(sp.H)).Draw(func(p *ui.Painter, r ui.Rect) {
			k.Canvas(p, r).Sprite(0, 0, sp, 0, slot, 1)
		})
		pixel.Text(c, message, pixel.TextOptions{Slot: slot})
	})
}
