package chart

import (
	"fmt"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type BoxPlotOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
	Slot          pixel.Slot
}

// quartile is the q-th quartile of sorted values, interpolated.
func quartile(sorted []float64, q float64) float64 {
	at := q * float64(len(sorted)-1)
	i := int(at)
	if i+1 >= len(sorted) {
		return sorted[i]
	}
	return sorted[i] + (sorted[i+1]-sorted[i])*(at-float64(i))
}

// BoxPlot boxes each group's quartiles; whiskers reach its ends.
func BoxPlot(c *ui.Context, text string, groups []Series, o BoxPlotOptions) *ui.Element {
	k := pixel.K(c)
	if len(groups) == 0 {
		panic("chart: BoxPlot has no groups")
	}
	w, h := size("BoxPlot", o.Width, o.Height, 160, 80)
	labels, rows := make([]string, len(groups)), make([][]float64, len(groups))
	for i, s := range groups {
		if len(s.Values) == 0 {
			panic(fmt.Sprintf("chart: BoxPlot group %q is empty", s.Name))
		}
		labels[i], rows[i] = s.Name, slices.Sorted(slices.Values(s.Values))
	}
	lo, hi := span(rows...)
	p := fit(lo, hi, w, h, labels)
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(fives(labels, rows)).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		p.grid(g, labels)
		bw := max(3, p.w/len(groups)/2)
		for i, row := range rows {
			x := p.slot(i, len(groups))
			q1, mid, q3 := p.at(quartile(row, 0.25)), p.at(quartile(row, 0.5)), p.at(quartile(row, 0.75))
			top, bottom := p.at(row[len(row)-1]), p.at(row[0])
			g.VLine(x, top, q3-top, pixel.Edge)
			g.VLine(x, q1, bottom-q1+1, pixel.Edge)
			g.HLine(x-bw/4, top, bw/2+1, pixel.Edge)
			g.HLine(x-bw/4, bottom, bw/2+1, pixel.Edge)
			g.Frame(x-bw/2, q3, bw, max(3, q1-q3+1), pixel.Flat, pick(groups[i].Slot, i))
			g.HLine(x-bw/2, mid, bw, pixel.Ink)
		}
	})
}

// fives lists each sorted row's five-number summary.
func fives(labels []string, rows [][]float64) string {
	parts := make([]string, len(rows))
	for i, row := range rows {
		parts[i] = labels[i] + " " + nums([]float64{row[0], quartile(row, 0.25), quartile(row, 0.5), quartile(row, 0.75), row[len(row)-1]})
	}
	return strings.Join(parts, "; ")
}
