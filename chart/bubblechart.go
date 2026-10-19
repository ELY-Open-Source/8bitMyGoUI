package chart

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type BubbleChartOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
	Slot          pixel.Slot
}

// BubbleChart draws points as discs, area by Size.
func BubbleChart(c *ui.Context, text string, points []Point, o BubbleChartOptions) ui.Element {
	k := pixel.K(c)
	w, h := size("BubbleChart", o.Width, o.Height, 160, 80)
	f := fitPoints(points, w, h, 8)
	most := 0.0
	for _, p := range points {
		if !(p.Size >= 0 && p.Size <= math.MaxFloat64) {
			panic(fmt.Sprintf("chart: bubble size %v is not a size", p.Size))
		}
		most = max(most, p.Size)
	}
	ink := pick(o.Slot, 0)
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(where(points, true)).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		f.grid(g)
		for _, p := range points {
			rad := 1 + int(math.Sqrt(part(p.Size, most))*6)
			g.Disc(f.ax(p.X), f.at(p.Y), rad, ink)
			g.Circle(f.ax(p.X), f.at(p.Y), rad, pixel.Edge)
		}
	})
}
