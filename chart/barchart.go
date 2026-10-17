package chart

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type BarChartOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
	// Slot fills the bars; zero means Accent.
	Slot pixel.Slot
}

// BarChart stands bars on micro-font axes; hover reads a value.
func BarChart(c *ui.Context, text string, labels []string, values []float64, o BarChartOptions) ui.Element {
	k := pixel.K(c)
	if len(labels) != len(values) || len(values) == 0 {
		panic(fmt.Sprintf("chart: BarChart has %d labels for %d values", len(labels), len(values)))
	}
	w, h := size("BarChart", o.Width, o.Height, 160, 80)
	lo, hi := span(values)
	p := fit(min(0, lo), max(0, hi), w, h, labels)
	e := ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(read(labels, values))
	hover := -1
	if x, y, over := e.PointerPosition(); over {
		if dx := int(x/k.Dot()) - p.x; dx >= 0 && dx < p.w && int(y/k.Dot()) < p.y+p.h {
			hover = dx * len(values) / p.w
		}
	}
	ink := pick(o.Slot, 0)
	return e.Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		p.grid(g, labels)
		share := p.w / len(values)
		bw, zero := max(1, share*2/3), p.at(0)
		for i, v := range values {
			x, top := p.slot(i, len(values))-bw/2, p.at(v)
			g.Rect(x, min(top, zero), bw, max(1, abs(zero-top)), ink)
		}
		if hover >= 0 {
			readout(g, p.slot(hover, len(values)), p.at(values[hover]), num(values[hover]), w)
		}
	})
}

// readout boxes s above x, y, kept inside w dots.
func readout(g pixel.Canvas, x, y int, s string, w int) {
	bw := tinyW(s) + 4
	bx, by := max(0, min(w-bw, x-bw/2)), max(0, y-11)
	g.Frame(bx, by, bw, 9, pixel.Flat, pixel.Lit)
	tiny(g, bx+2, by+2, s, pixel.Ink)
}

// read lists values for assistive technology: "A 3, B 5".
func read(labels []string, values []float64) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = labels[i] + " " + num(v)
	}
	return strings.Join(parts, ", ")
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
