package chart

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PieChartOptions struct {
	// Size is the pie's width in dots; zero means 64.
	Size int
}

// PieChart rasterizes one sector per slice, dot by dot.
func PieChart(c *ui.Context, text string, slices []Slice, o PieChartOptions) ui.Element {
	k := pixel.K(c)
	n := pixel.Or("PieChart Size", o.Size, 64)
	_, cuts := shares(slices)
	return ui.Box(c).Size(k.D(n), k.D(n)).Role(ui.RoleImage).Label(text).Description(readSlices(slices)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		c, rad := n/2, n/2-1
		for i, s := range slices {
			sector(g, c, c, 0, rad, cuts[i], cuts[i+1], pick(s.Slot, i))
		}
		g.Circle(c, c, rad, pixel.Edge)
	})
}
