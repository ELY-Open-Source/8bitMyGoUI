package chart

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PolarAreaOptions struct {
	// Size is the chart's width in dots; zero means 64.
	Size int
}

// PolarArea gives slices equal turns, reaching out by value.
func PolarArea(c *ui.Context, text string, slices []Slice, o PolarAreaOptions) ui.Element {
	k := pixel.K(c)
	n := pixel.Or("PolarArea Size", o.Size, 64)
	shares(slices)
	most := 0.0
	for _, s := range slices {
		most = max(most, s.Value)
	}
	return ui.Box(c).Size(k.D(n), k.D(n)).Role(ui.RoleImage).Label(text).Description(readSlices(slices)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		c, rad := n/2, n/2-1
		g.Circle(c, c, rad, pixel.Shade)
		g.Circle(c, c, rad/2, pixel.Shade)
		step := 1 / float64(len(slices))
		for i, s := range slices {
			if s.Value > 0 {
				reach := max(1, int(s.Value/most*float64(rad)))
				sector(g, c, c, 0, reach, float64(i)*step, float64(i+1)*step, pick(s.Slot, i))
			}
		}
	})
}
