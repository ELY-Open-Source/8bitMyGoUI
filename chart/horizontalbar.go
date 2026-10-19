package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type HorizontalBarOptions struct {
	// Width is in dots, zero means 160; rows are lines.
	Width int
	// Slot fills the bars; zero means Accent.
	Slot pixel.Slot
}

// HorizontalBar lays labeled bars across, longest to the edge.
func HorizontalBar(c *ui.Context, text string, labels []string, values []float64, o HorizontalBarOptions) ui.Element {
	k := pixel.K(c)
	if len(labels) != len(values) || len(values) == 0 {
		panic(fmt.Sprintf("chart: HorizontalBar has %d labels for %d values", len(labels), len(values)))
	}
	w := pixel.Or("HorizontalBar Width", o.Width, 160)
	lo, hi := span(values)
	if lo < 0 {
		panic(fmt.Sprintf("chart: HorizontalBar value %v is negative", lo))
	}
	gutter, room := 0, 0
	for i, l := range labels {
		gutter, room = max(gutter, pixel.Measure(l)), max(room, tinyW(num(values[i])))
	}
	ink := pick(o.Slot, 0)
	return ui.Box(c).Size(k.D(w), k.D(len(values)*pixel.LineHeight)).Role(ui.RoleImage).Label(text).Description(read(labels, values)).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		span := w - gutter - room - 8
		for i, v := range values {
			y, bw := i*pixel.LineHeight, int(part(v, hi)*float64(span))
			g.Text(gutter-pixel.Measure(labels[i]), y, labels[i], pixel.Faint, 1)
			g.Rect(gutter+4, y+4, max(1, bw), pixel.LineHeight-8, ink)
			tiny(g, gutter+6+bw, y+6, num(v), pixel.Ink)
		}
	})
}
