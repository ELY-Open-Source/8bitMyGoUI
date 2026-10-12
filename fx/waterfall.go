package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type WaterfallOptions struct {
	// Width and Height are dots; zeros mean 64 by 96.
	Width, Height int
	Label         string
}

// Waterfall pours between cliffs by cycling its water slots.
func Waterfall(c *ui.Context, o WaterfallOptions) *ui.Element {
	k := pixel.K(c)
	w, h := size("Waterfall", o.Width, o.Height, 64, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { waterfall(g, k, g.Step(10)) })
	})
}

// waterfall lights one three-row segment in six per column.
func waterfall(g pixel.Canvas, k pixel.Kit, step int) {
	l, r, pool := g.W/4, g.W-g.W/4, g.H-g.H/6
	rock, lit := pick(k, pixel.Faint, pixel.Face), pick(k, pixel.Shade, pixel.Lit)
	g.Rect(0, 0, l, g.H, rock)
	g.Rect(r, 0, g.W-r, g.H, rock)
	for i := range g.H / 3 {
		x := rnd(i, 110) % max(1, l-2)
		if i%2 == 1 {
			x += r + 1
		}
		g.HLine(x, rnd(i, 111)%pool, 2, lit)
	}
	g.Rect(l, 0, r-l, pool, pixel.Info)
	// Segments take the band turned back each step: falling.
	band := pixel.Cycle([]pixel.Slot{light(k), pixel.Info, pixel.Info, pixel.Info, pixel.Info, pixel.Info}, -step)
	for x := l; x < r; x++ {
		off := rnd(x, 112) % 6
		for seg := 0; 3*seg < pool; seg++ {
			if s := band[(seg+off)%6]; s != pixel.Info {
				g.VLine(x, 3*seg, min(3, pool-3*seg), s)
			}
		}
	}
	// Foam churns where it lands; ripples run outward.
	g.Rect(0, pool, g.W, g.H-pool, pixel.Info)
	g.Veil(l-3, pool, r-l+6, 3, light(k), 6+4*wrap(step, 2))
	for i := range 2 {
		d := wrap(step/2+4*i, 8)
		g.HLine(l-6-2*d, pool+4+2*i, 3, light(k))
		g.HLine(r+3+2*d, pool+4+2*i, 3, light(k))
	}
}
