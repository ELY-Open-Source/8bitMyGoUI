package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type BulletChartOptions struct {
	// Width is in dots; zero means 160.
	Width int
}

// BulletChart lays a bar on rising ranges and a target.
func BulletChart(c *ui.Context, text string, value, target float64, ranges []float64, o BulletChartOptions) *ui.Element {
	k := pixel.K(c)
	if len(ranges) == 0 {
		panic("chart: BulletChart needs ranges")
	}
	span(ranges, []float64{value, target})
	for i, r := range ranges {
		if r <= 0 || i > 0 && r <= ranges[i-1] {
			panic(fmt.Sprintf("chart: BulletChart ranges %v must rise above zero", ranges))
		}
	}
	end := ranges[len(ranges)-1]
	if value < 0 || value > end || target < 0 || target > end {
		panic(fmt.Sprintf("chart: BulletChart value %v or target %v outside 0–%v", value, target, end))
	}
	w := pixel.Or("BulletChart Width", o.Width, 160)
	x := func(v float64) int { return int(v / end * float64(w-1)) }
	return ui.Box(c).Size(k.D(w), k.D(14)).Role(ui.RoleMeter).Label(text).Description(num(value) + " / " + num(target)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		from := 0
		for i, band := range ranges {
			to := x(band) + 1
			g.Dither(from, 0, to-from, 14, pixel.Well, pixel.Shade, 12-12*i/len(ranges))
			from = to
		}
		g.Rect(0, 5, max(1, x(value)), 4, pixel.Ink)
		g.Rect(min(x(target), w-2), 2, 2, 10, pixel.Accent)
	})
}
