package chart

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SparklineOptions struct {
	// Width is in dots; zero means 48.
	Width int
	Slot  pixel.Slot
}

// Sparkline is a line one text line high, end marked.
func Sparkline(c *ui.Context, text string, values []float64, o SparklineOptions) *ui.Element {
	k := pixel.K(c)
	w := pixel.Or("Sparkline Width", o.Width, 48)
	if len(values) < 2 {
		panic("chart: Sparkline needs two values or more")
	}
	lo, hi := span(values)
	if hi == lo {
		hi = lo + 1
	}
	ink := pick(o.Slot, 0)
	at := func(i int) (int, int) {
		return i * (w - 2) / (len(values) - 1), pixel.LineHeight - 3 - int((values[i]-lo)/(hi-lo)*float64(pixel.LineHeight-6))
	}
	return ui.Box(c).Size(k.D(w), k.D(pixel.LineHeight)).Role(ui.RoleImage).Label(text).Description(nums(values)).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		for i := 1; i < len(values); i++ {
			x0, y0 := at(i - 1)
			x1, y1 := at(i)
			g.Line(x0, y0, x1, y1, ink)
		}
		x, y := at(len(values) - 1)
		g.Rect(x, y-1, 2, 2, pixel.Accent)
	})
}
