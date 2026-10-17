package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LineChartOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
}

// LineChart draws series as Bresenham lines with square markers.
func LineChart(c *ui.Context, text string, labels []string, series []Series, o LineChartOptions) ui.Element {
	k := pixel.K(c)
	w, h := size("LineChart", o.Width, o.Height, 160, 80)
	p := fitSeries(labels, series, w, h)
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(readSeries(labels, series)).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		p.grid(g, labels)
		lines(g, p, len(labels), series)
	})
}

// fitSeries checks series against labels and fits their range.
func fitSeries(labels []string, series []Series, w, h int) plot {
	rows := make([][]float64, len(series))
	for i, s := range series {
		if len(s.Values) != len(labels) {
			panic(fmt.Sprintf("chart: series %q has %d values for %d labels", s.Name, len(s.Values), len(labels)))
		}
		rows[i] = s.Values
	}
	lo, hi := span(rows...)
	return fit(lo, hi, w, h, labels)
}

// lines joins each series' points and marks them with squares.
func lines(g pixel.Canvas, p plot, n int, series []Series) {
	for j, s := range series {
		ink := pick(s.Slot, j)
		for i := 1; i < n; i++ {
			g.Line(p.slot(i-1, n), p.at(s.Values[i-1]), p.slot(i, n), p.at(s.Values[i]), ink)
		}
		for i, v := range s.Values {
			g.Rect(p.slot(i, n)-1, p.at(v)-1, 3, 3, ink)
		}
	}
}

// readSeries lists each series for assistive technology.
func readSeries(labels []string, series []Series) string {
	out := ""
	for i, s := range series {
		if i > 0 {
			out += "; "
		}
		out += s.Name + ": " + read(labels, s.Values)
	}
	return out
}
