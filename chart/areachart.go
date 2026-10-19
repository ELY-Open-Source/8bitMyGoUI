package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type AreaChartOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
	Slot          pixel.Slot
}

// AreaChart draws a line with the area below it dithered.
func AreaChart(c *ui.Context, text string, labels []string, values []float64, o AreaChartOptions) ui.Element {
	k := pixel.K(c)
	w, h := size("AreaChart", o.Width, o.Height, 160, 80)
	if len(values) != len(labels) || len(values) < 2 {
		panic(fmt.Sprintf("chart: AreaChart has %d labels for %d values", len(labels), len(values)))
	}
	s := []Series{{Name: text, Values: values, Slot: pick(o.Slot, 0)}}
	lo, hi := span(values)
	p := fit(min(0, lo), max(0, hi), w, h, labels)
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(read(labels, values)).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		p.grid(g, labels)
		// The area spans from the line to zero, either way.
		n, base := len(values), p.at(0)
		for i := 1; i < n; i++ {
			x0, x1 := p.slot(i-1, n), p.slot(i, n)
			for x := x0; x < x1; x++ {
				y := p.at(values[i-1] + (values[i]-values[i-1])*float64(x-x0)/float64(x1-x0))
				g.Veil(x, min(y, base), 1, abs(y-base), s[0].Slot, 6)
			}
		}
		lines(g, p, n, s)
	})
}
