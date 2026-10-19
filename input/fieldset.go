package input

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type FieldsetOptions struct{}

// Fieldset frames a group of fields under a legend.
func Fieldset(c *ui.Context, legend string, o FieldsetOptions, content func()) ui.Element {
	k := pixel.K(c)
	lw := pixel.Measure(legend)
	return ui.Column(c).Padding(k.D(pixel.LineHeight), k.D(8), k.D(8), k.D(8)).Gap(k.D(6)).Role(ui.RoleGroup).Label(legend).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		top := pixel.LineHeight / 2
		g.Frame(0, top, g.W, g.H-top, pixel.Double, pixel.Face)
		g.Rect(6, top-1, lw+6, 4, pixel.Face)
		g.Text(9, 0, legend, pixel.Accent, 1)
	}).Children(content)
}
