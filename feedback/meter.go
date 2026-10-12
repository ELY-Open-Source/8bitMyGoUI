package feedback

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MeterOptions struct {
	// Min and Max are zero for 0 to 1.
	Min, Max float64
	// Warning and Critical are levels; zero means none.
	Warning, Critical float64
	// Width is in dots; zero means 120.
	Width int
}

// Meter shows a level in Good, Warn or Bad.
func Meter(c *ui.Context, text string, value float64, o MeterOptions) *ui.Element {
	k := pixel.K(c)
	lo, hi := o.Min, o.Max
	if lo == 0 && hi == 0 {
		hi = 1
	}
	if hi <= lo {
		panic("feedback: Meter range is empty")
	}
	w := pixel.Or("Meter Width", o.Width, 120)
	at := func(v float64) int { return 2 + int(max(0, min(1, (v-lo)/(hi-lo)))*float64(w-4)) }
	ink := level(value, o.Warning, o.Critical)
	return ui.Box(c).Role(ui.RoleMeter).Label(text).Description(fmtNum(value)).Size(k.D(w), k.D(10)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, w, 10, pixel.Sunken, pixel.Well)
		g.Rect(2, 2, at(value)-2, 6, ink)
		for _, mark := range []float64{o.Warning, o.Critical} {
			if mark != 0 {
				g.VLine(at(mark), 1, 8, pixel.Edge)
			}
		}
	})
}
