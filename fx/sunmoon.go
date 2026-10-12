package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SunMoonOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// SunMoon arcs the sun by day, the moon by night.
func SunMoon(c *ui.Context, hour int, o SunMoonOptions) *ui.Element {
	checkHour("SunMoon", hour)
	k := pixel.K(c)
	w, h := size("SunMoon", o.Width, o.Height, 160, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { sunMoon(g, k, hour, g.Step(4)) })
	})
}

// arc centers the body: low at the ends, high midway.
func arc(w, h, hour int) (int, int) {
	t := wrap(hour-6, 12)
	horizon, top := h*3/4, 10
	x := 10 + (w-20)*(2*t+1)/24
	return x, horizon - (horizon-top)*(2*t+1)*(23-2*t)/143
}

// sunMoon draws a sun with turning rays, or cratered moon.
func sunMoon(g pixel.Canvas, k pixel.Kit, hour, step int) {
	x, y := arc(g.W, g.H, hour)
	if hour < 6 || hour >= 18 {
		g.Disc(x, y, 5, light(k))
		g.Dot(x-2, y-1, pixel.Faint)
		g.Rect(x+1, y+1, 2, 1, pixel.Faint)
		g.Dot(x+2, y-3, pixel.Faint)
		return
	}
	// Rays turn an eighth each step: straight, then slanted.
	for i := range 4 {
		dx, dy := [4]int{1, 0, -1, 0}[i], [4]int{0, 1, 0, -1}[i]
		if step%2 == 1 {
			dx, dy = [4]int{1, 1, -1, -1}[i], [4]int{1, -1, 1, -1}[i]
		}
		g.Line(x+7*dx, y+7*dy, x+8*dx, y+8*dy, pixel.Warn)
	}
	g.Disc(x, y, 5, pixel.Warn)
	g.Rect(x-2, y-2, 2, 2, light(k))
}
