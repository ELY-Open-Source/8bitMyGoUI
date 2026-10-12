package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type HeatmapOptions struct {
	// Cell is a cell's least side, zero for 10.
	Cell int
	// Low and High end the dither; zeros mean Well, Bad.
	Low, High pixel.Slot
}

// Heatmap shades cells by dither level between two slots.
func Heatmap(c *ui.Context, text string, rows, cols []string, values [][]float64, o HeatmapOptions) *ui.Element {
	k := pixel.K(c)
	if len(values) != len(rows) {
		panic(fmt.Sprintf("chart: Heatmap has %d rows of values for %d rows", len(values), len(rows)))
	}
	for i, r := range values {
		if len(r) != len(cols) {
			panic(fmt.Sprintf("chart: Heatmap row %q has %d values for %d columns", rows[i], len(r), len(cols)))
		}
	}
	cell := pixel.Or("Heatmap Cell", o.Cell, 10)
	low, high := o.Low, o.High
	if low == 0 && high == 0 {
		low, high = pixel.Well, pixel.Bad
	}
	lo, hi := span(values...)
	if hi == lo {
		hi = lo + 1
	}
	rs, cs := small(rows...), small(cols...)
	gutter := 0
	for _, r := range rows {
		gutter = max(gutter, markW(r, rs)+3)
	}
	// Cells grow to fit their labels.
	cw, ch := cell, max(cell, inkH(rs)+2)
	for _, col := range cols {
		cw = max(cw, markW(col, cs)+3)
	}
	head := labelH(cols)
	w, h := gutter+len(cols)*cw, head+len(rows)*ch
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(readSeries(cols, grid(rows, values))).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for j, col := range cols {
			mark(g, gutter+j*cw+(cw-markW(col, cs))/2, 1, col, cs, pixel.Faint)
		}
		for i, row := range values {
			mark(g, 0, head+i*ch+(ch-1-inkH(rs))/2, rows[i], rs, pixel.Faint)
			for j, v := range row {
				level := int((v - lo) / (hi - lo) * 16)
				g.Dither(gutter+j*cw, head+i*ch, cw-1, ch-1, low, high, level)
			}
		}
	})
}

// grid names each row of values as a series.
func grid(rows []string, values [][]float64) []Series {
	out := make([]Series, len(rows))
	for i, r := range rows {
		out[i] = Series{Name: r, Values: values[i]}
	}
	return out
}
