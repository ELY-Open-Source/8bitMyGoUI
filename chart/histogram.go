package chart

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type HistogramOptions struct {
	// Bins is roughly how many, zero for 8; edges round.
	Bins int
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
	Slot          pixel.Slot
}

// Histogram counts values into bins, side by side.
func Histogram(c *ui.Context, text string, values []float64, o HistogramOptions) *ui.Element {
	k := pixel.K(c)
	n := pixel.Or("Histogram Bins", o.Bins, 8)
	w, h := size("Histogram", o.Width, o.Height, 160, 80)
	if len(values) == 0 {
		panic("chart: Histogram has no values")
	}
	lo, hi := span(values)
	// A lone value gets a bin sized to itself.
	run := hi - lo
	if run == 0 {
		run = max(math.Abs(lo), 1)
	}
	step := nice(run / float64(n))
	from := math.Floor(lo/step) * step
	bins := int((hi-from)/step) + 1
	counts, labels := make([]float64, bins), make([]string, bins)
	for _, v := range values {
		counts[int((v-from)/step)]++
	}
	for i := range bins {
		labels[i] = tick(from+float64(i)*step, step)
	}
	_, top := span(counts)
	p := fit(0, top, w, h, labels)
	ink := pick(o.Slot, 0)
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(read(labels, counts)).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		p.grid(g, labels)
		for i, v := range counts {
			x0, x1 := p.x+i*p.w/bins, p.x+(i+1)*p.w/bins
			g.Rect(x0, p.at(v), x1-x0-1, p.at(0)-p.at(v), ink)
		}
	})
}
