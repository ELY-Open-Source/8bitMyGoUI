package feedback

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/text"
)

var openChest = pixel.MustSprite(`
	..eeeeeeeeeeee..
	.esssssssssssse.
	.elllllllllllle.
	..eeeeeeeeeeee..
	................
	eeeeeeeeeeeeeeee
	ewwwwwwwwwwwwwwe
	ewwwwwwwwwwwwwwe
	eeeeeeeyyeeeeeee
	essssssyysssssse
	esssssssssssssse
	eeeeeeeeeeeeeeee
`)

type EmptyStateOptions struct {
	// Action labels a button; Acted reports it.
	Action   string
	Disabled bool
}

type EmptyStateResult struct {
	Element ui.Element
	Acted   bool
}

// EmptyState shows NO ITEMS over an empty chest.
func EmptyState(c *ui.Context, message string, o EmptyStateOptions) EmptyStateResult {
	k := pixel.K(c)
	res := EmptyStateResult{}
	res.Element = ui.Column(c).Gap(k.D(6)).Padding(k.D(12)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(k.Msg("empty")).Children(func() {
		display.Sprite(c, openChest, display.SpriteOptions{Scale: 2})
		text.ShadowText(c, k.Msg("empty"), text.ShadowTextOptions{Slot: pixel.Faint, Scale: 2})
		pixel.Text(c, message, pixel.TextOptions{Slot: pixel.Faint})
		if o.Action != "" {
			res.Acted = input.Button(c, o.Action, input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled}).Clicked
		}
	})
	return res
}
