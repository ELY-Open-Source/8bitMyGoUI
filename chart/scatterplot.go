package chart

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ScatterPlotOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
	Slot          pixel.Slot
}

// ScatterPlot puts each point as a 2×2 square.
func ScatterPlot(c *ui.Context, text string, points []Point, o ScatterPlotOptions) *ui.Element {
	k := pixel.K(c)
	w, h := size("ScatterPlot", o.Width, o.Height, 160, 80)
	f := fitPoints(points, w, h, 2)
	ink := pick(o.Slot, 0)
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(where(points, false)).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		f.grid(g)
		for _, p := range points {
			g.Rect(f.ax(p.X), f.at(p.Y)-1, 2, 2, ink)
		}
	})
}
