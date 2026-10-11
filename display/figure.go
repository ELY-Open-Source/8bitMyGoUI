package display

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type FigureOptions struct {
	// Slot inks the caption; zero means Faint.
	Slot pixel.Slot
}

// Figure shows content over a centered caption.
func Figure(c *ui.Context, caption string, o FigureOptions, content func()) *ui.Element {
	k := pixel.K(c)
	return ui.Column(c).Gap(k.D(3)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(caption).Children(func() {
		content()
		pixel.Text(c, caption, pixel.TextOptions{Slot: slotOr(o.Slot, pixel.Faint)})
	})
}
