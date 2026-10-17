package chart

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DonutChartOptions struct {
	// Size is the ring's width in dots; zero means 64.
	Size int
}

// DonutChart is a pie with a hole holding the total.
func DonutChart(c *ui.Context, text string, slices []Slice, o DonutChartOptions) ui.Element {
	k := pixel.K(c)
	n := pixel.Or("DonutChart Size", o.Size, 64)
	total, cuts := shares(slices)
	sum := num(total)
	return ui.Box(c).Size(k.D(n), k.D(n)).Role(ui.RoleImage).Label(text).Description(readSlices(slices)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		c, rad := n/2, n/2-1
		for i, s := range slices {
			sector(g, c, c, rad/2, rad, cuts[i], cuts[i+1], pick(s.Slot, i))
		}
		g.Circle(c, c, rad, pixel.Edge)
		g.Circle(c, c, rad/2, pixel.Edge)
		tiny(g, c-tinyW(sum)/2+1, c-2, sum, pixel.Ink)
	})
}
