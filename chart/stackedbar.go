package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type StackedBarOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
}

// StackedBar stacks each category's series, one slot per series.
func StackedBar(c *ui.Context, text string, labels []string, series []Series, o StackedBarOptions) *ui.Element {
	k := pixel.K(c)
	if len(labels) == 0 {
		panic("chart: StackedBar has no labels")
	}
	w, h := size("StackedBar", o.Width, o.Height, 160, 80)
	totals := make([]float64, len(labels))
	for _, s := range series {
		if len(s.Values) != len(labels) {
			panic(fmt.Sprintf("chart: series %q has %d values for %d labels", s.Name, len(s.Values), len(labels)))
		}
		for i, v := range s.Values {
			if v < 0 {
				panic(fmt.Sprintf("chart: series %q stacks a negative %v", s.Name, v))
			}
			totals[i] += v
		}
	}
	_, hi := span(totals)
	p := fit(0, hi, w, h, labels)
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(readSeries(labels, series)).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		p.grid(g, labels)
		bw := max(1, p.w/len(labels)*2/3)
		for i := range labels {
			x, base := p.slot(i, len(labels))-bw/2, 0.0
			for j, s := range series {
				lo, top := p.at(base), p.at(base+s.Values[i])
				g.Rect(x, top, bw, lo-top, pick(s.Slot, j))
				base += s.Values[i]
			}
		}
	})
}
