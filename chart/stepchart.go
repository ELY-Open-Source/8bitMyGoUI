package chart

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type StepChartOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
	Slot          pixel.Slot
}

// StepChart holds each value flat, then steps to the next.
func StepChart(c *ui.Context, text string, labels []string, values []float64, o StepChartOptions) ui.Element {
	k := pixel.K(c)
	w, h := size("StepChart", o.Width, o.Height, 160, 80)
	p := fitSeries(labels, []Series{{Name: text, Values: values}}, w, h)
	ink := pick(o.Slot, 0)
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(read(labels, values)).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		p.grid(g, labels)
		n := len(values)
		for i, v := range values {
			x0, x1, y := p.x+i*p.w/n, p.x+(i+1)*p.w/n, p.at(v)
			g.HLine(x0, y, x1-x0, ink)
			if i+1 < n {
				y1 := p.at(values[i+1])
				g.VLine(x1, min(y, y1), abs(y-y1)+1, ink)
			}
		}
	})
}
