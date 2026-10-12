package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type FunnelOptions struct {
	// Width is in dots; zero means 160.
	Width int
}

// Funnel narrows stage by stage, each bar with its count.
func Funnel(c *ui.Context, text string, stages []Slice, o FunnelOptions) *ui.Element {
	k := pixel.K(c)
	if len(stages) == 0 {
		panic("chart: Funnel has no stages")
	}
	w := pixel.Or("Funnel Width", o.Width, 160)
	top := stages[0].Value
	for i, s := range stages {
		span([]float64{s.Value})
		if s.Value < 0 || i > 0 && s.Value > stages[i-1].Value {
			panic(fmt.Sprintf("chart: Funnel stage %q has %v; stages must narrow from %v", s.Name, s.Value, top))
		}
	}
	gutter := 0
	for _, s := range stages {
		gutter = max(gutter, pixel.Measure(s.Name)+6)
	}
	return ui.Box(c).Size(k.D(w), k.D(len(stages)*pixel.LineHeight)).Role(ui.RoleImage).Label(text).Description(readSlices(stages)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		room := w - gutter
		for i, s := range stages {
			bw := max(2, int(part(s.Value, top)*float64(room)))
			x, y := gutter+(room-bw)/2, i*pixel.LineHeight
			g.Text(0, y, s.Name, pixel.Faint, 1)
			g.Rect(x, y+2, bw, pixel.LineHeight-4, pick(s.Slot, i))
			// The count sits on its bar, or beside thin ones.
			n := num(s.Value)
			if tinyW(n)+4 <= bw {
				tiny(g, x+(bw-tinyW(n))/2, y+6, n, pixel.AccentInk)
			} else {
				tiny(g, x+bw+3, y+6, n, pixel.Ink)
			}
		}
	})
}
