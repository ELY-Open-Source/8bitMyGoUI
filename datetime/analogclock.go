package datetime

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type AnalogClockOptions struct {
	// Size is the face's width in dots; zero means 48.
	Size int
}

// AnalogClock ticks each second; its hands are Bresenham lines.
func AnalogClock(c *ui.Context, o AnalogClockOptions) ui.Element {
	k := pixel.K(c)
	size := pixel.Or("AnalogClock Size", o.Size, 48)
	// Wake on the next whole second.
	now := k.Now()
	k.At(c, now.Truncate(time.Second).Add(time.Second))
	h, m, s := now.Clock()
	// hand ends at turn (0–1 from twelve) and length r.
	hand := func(cx, cy int, turn float64, r int) (int, int) {
		a := turn * 2 * math.Pi
		return cx + int(math.Round(float64(r)*math.Sin(a))), cy - int(math.Round(float64(r)*math.Cos(a)))
	}
	return ui.Box(c).Size(k.D(size), k.D(size)).Role(ui.RoleImage).Label(now.Format("15:04:05")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		cx, cy, rad := size/2, size/2, size/2-1
		g.Disc(cx, cy, rad, pixel.Face)
		g.Circle(cx, cy, rad, pixel.Edge)
		for i := range 12 {
			x, y := hand(cx, cy, float64(i)/12, rad-3)
			g.Dot(x, y, pixel.Edge)
		}
		x, y := hand(cx, cy, (float64(h%12)+float64(m)/60)/12, rad/2)
		g.Line(cx, cy, x, y, pixel.Ink)
		x, y = hand(cx, cy, (float64(m)+float64(s)/60)/60, rad-5)
		g.Line(cx, cy, x, y, pixel.Ink)
		x, y = hand(cx, cy, float64(s)/60, rad-4)
		g.Line(cx, cy, x, y, pixel.Accent)
		g.Dot(cx, cy, pixel.Accent)
	})
}
