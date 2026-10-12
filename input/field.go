package input

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type FieldOptions struct {
	// Description shows below the control in Faint.
	Description string
	// Error shows below it in Bad, with a ! sign.
	Error string
}

// Field puts a label above a control, notes below.
func Field(c *ui.Context, label string, o FieldOptions, control func()) *ui.Element {
	k := pixel.K(c)
	return ui.Column(c).Gap(k.D(2)).Role(ui.RoleGroup).Label(label).Children(func() {
		pixel.Text(c, label, pixel.TextOptions{})
		control()
		if o.Description != "" {
			pixel.Text(c, o.Description, pixel.TextOptions{Slot: pixel.Faint, Wrap: true})
		}
		if o.Error != "" {
			ui.Row(c).Gap(k.D(3)).AlignItems(ui.Center).Children(func() {
				display.Icon(c, icon.Warning, display.IconOptions{Slot: pixel.Bad})
				pixel.Text(c, o.Error, pixel.TextOptions{Slot: pixel.Bad, Wrap: true}).Grow(1)
			})
		}
	})
}
