package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type WaterfallOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
}

// Waterfall walks a running total: Good up, Bad down.
func Waterfall(c *ui.Context, text string, labels []string, deltas []float64, o WaterfallOptions) ui.Element {
	k := pixel.K(c)
	if len(labels) != len(deltas) || len(deltas) == 0 {
		panic(fmt.Sprintf("chart: Waterfall has %d labels for %d steps", len(labels), len(deltas)))
	}
	w, h := size("Waterfall", o.Width, o.Height, 160, 80)
	totals := make([]float64, len(deltas)+1)
	for i, d := range deltas {
		totals[i+1] = totals[i] + d
	}
	lo, hi := span(totals)
	p := fit(lo, hi, w, h, labels)
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(read(labels, deltas)).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		p.grid(g, labels)
		bw := max(1, p.w/len(deltas)*2/3)
		for i, d := range deltas {
			a, b := p.at(totals[i]), p.at(totals[i+1])
			ink := pixel.Good
			if d < 0 {
				ink = pixel.Bad
			}
			x := p.slot(i, len(deltas)) - bw/2
			g.Rect(x, min(a, b), bw, max(1, abs(a-b)), ink)
			if i+1 < len(deltas) {
				g.Checker(x+bw, b, p.slot(i+1, len(deltas))-bw/2-x-bw, 1, pixel.Edge)
			}
		}
	})
}
