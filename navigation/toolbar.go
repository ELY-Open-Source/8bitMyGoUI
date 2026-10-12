package navigation

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ToolbarOptions struct{}

// Toolbar holds tool buttons as one Tab stop; arrows move.
func Toolbar(c *ui.Context, text string, o ToolbarOptions, content func()) *ui.Element {
	k := pixel.K(c)
	return ui.Row(c).AlignSelf(ui.Start).Gap(k.D(2)).Padding(k.D(3)).AlignItems(ui.Center).Role(ui.RoleToolbar).Label(text).FocusGroup(ui.Horizontal).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Raised, pixel.Face)
	}).Children(content)
}
