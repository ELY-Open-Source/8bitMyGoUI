package display

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var (
	rise = pixel.MustSprite(`
		..#..
		.###.
		#####
	`)
	fall = pixel.MustSprite(`
		#####
		.###.
		..#..
	`)
)

type StatOptions struct {
	// Delta shows under the value, ▲ in Good.
	Delta string
	// Down turns the delta to ▼ in Bad.
	Down bool
}

// Stat is a big value with a label and delta.
func Stat(c *ui.Context, label, value string, o StatOptions) *ui.Element {
	k := pixel.K(c)
	return ui.Column(c).Gap(k.D(1)).Role(ui.RoleGroup).Label(label).Children(func() {
		pixel.Text(c, label, pixel.TextOptions{Slot: pixel.Faint})
		pixel.Text(c, value, pixel.TextOptions{Scale: 2})
		if o.Delta != "" {
			sp, ink := rise, pixel.Good
			if o.Down {
				sp, ink = fall, pixel.Bad
			}
			ui.Row(c).Gap(k.D(3)).AlignItems(ui.Center).Children(func() {
				Icon(c, sp, IconOptions{Slot: ink})
				pixel.Text(c, o.Delta, pixel.TextOptions{Slot: ink})
			})
		}
	})
}
