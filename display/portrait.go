package display

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PortraitOptions struct {
	// Size is the picture's side in dots; zero means 48.
	Size int
}

// Portrait frames a character's picture over a nameplate.
func Portrait(c *ui.Context, name string, o PortraitOptions, picture func()) ui.Element {
	k := pixel.K(c)
	size := pixel.Or("Portrait Size", o.Size, 48)
	return ui.Column(c).Width(k.D(size + 6)).Role(ui.RoleGroup).Label(name).Children(func() {
		ui.Box(c).Size(k.D(size+6), k.D(size+6)).Padding(k.D(3)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Double, pixel.Well)
		}).Children(func() {
			ui.Box(c).Size(k.D(size), k.D(size)).Clip().Children(picture)
		})
		ui.Box(c).Size(k.D(size+6), k.D(pixel.LineHeight+2)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Notched, pixel.Face)
			g.Clip(1, 0, g.W-2, g.H, func() { g.Text((g.W-pixel.Measure(name))/2, 1, name, pixel.Ink, 1) })
		})
	})
}
