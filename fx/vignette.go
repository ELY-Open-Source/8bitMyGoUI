package fx

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type VignetteOptions struct {
	// Slot darkens the edges; zero means the darkest slot.
	Slot pixel.Slot
}

// Vignette dithers a slot inward from edges; corners go solid.
func Vignette(c *ui.Context, o VignetteOptions, content func()) ui.Element {
	k := pixel.K(c)
	s := slotOr(o.Slot, dark(k))
	return cover(k, ui.Column(c), func(g pixel.Canvas) { vignette(g, s) }, content)
}

// vignette's level grows with the elliptic distance out.
func vignette(g pixel.Canvas, s pixel.Slot) {
	level := func(x, y int) int {
		ex := float64(2*x+1-g.W) / float64(g.W)
		ey := float64(2*y+1-g.H) / float64(g.H)
		return min(16, max(0, int((ex*ex+ey*ey-0.6)*16/1.2)))
	}
	for y := range g.H {
		for x := 0; x < g.W; {
			l, end := level(x, y), x+1
			for end < g.W && level(end, y) == l {
				end++
			}
			if l > 0 {
				g.Veil(x, y, end-x, 1, s, l)
			}
			x = end
		}
	}
}
