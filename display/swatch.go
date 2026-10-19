package display

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type SwatchOptions struct {
	// Columns is zero for 4.
	Columns int
}

// Swatch shows the active palette's sixteen slots with names.
func Swatch(c *ui.Context, o SwatchOptions) ui.Element {
	k := pixel.K(c)
	cols := pixel.Or("Swatch Columns", o.Columns, 4)
	cell := 14 + pixel.Measure("AccentLit")
	return ui.Row(c).Wrap().Width(k.D(cols * (cell + 4))).Gap(k.D(4)).Role(ui.RoleList).Label(k.Palette().Name).Children(func() {
		for s := range pixel.Slot(pixel.Slots) {
			ui.Box(c).Size(k.D(cell), k.D(pixel.LineHeight)).Role(ui.RoleListItem).Label(s.String()).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Frame(0, 2, 12, 12, pixel.Flat, s)
				g.Text(14, 0, s.String(), pixel.Ink, 1)
			})
		}
	})
}
