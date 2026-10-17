package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LavaOptions struct {
	// Width and Height are dots; zeros mean 160 by 48.
	Width, Height int
	Label         string
}

// Lava flows in cycling veins while bubbles swell and pop.
func Lava(c *ui.Context, o LavaOptions) ui.Element {
	k := pixel.K(c)
	w, h := size("Lava", o.Width, o.Height, 160, 48)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { lava(g, k, g.Step(6)) })
	})
}

// lava's veins flow and cycle row by row; crust drifts.
func lava(g pixel.Canvas, k pixel.Kit, step int) {
	g.Rect(0, 0, g.W, g.H, pixel.Bad)
	band := pixel.Cycle([]pixel.Slot{pixel.Warn, pixel.Warn, light(k), pixel.Warn, pixel.Bad, pixel.Bad}, -step)
	for j := 0; 4*j < g.H; j++ {
		s, dir := band[j%len(band)], 1-2*(j%2)
		for i := range g.W/14 + 1 {
			x := wrap(rnd(i, 130+j)+dir*step, g.W+16) - 8
			g.Rect(x, 4*j+1+rnd(i, 140+j)%2, 4+rnd(i, 150+j)%8, 1, s)
		}
	}
	crust := pick(k, pixel.Edge, pixel.Shade)
	for i := range g.W * g.H / 300 {
		g.Rect(wrap(rnd(i, 160)+step/2, g.W), rnd(i, 161)%g.H, 2+rnd(i, 162)%3, 1, crust)
	}
	for i := range max(2, g.W/24) {
		t := step + rnd(i, 130)
		// Each life of a bubble rises somewhere new.
		x, y := rnd(t/12*64+i, 131)%g.W, 4+rnd(t/12*64+i, 132)%max(1, g.H-8)
		switch wrap(t, 12) {
		case 4:
			g.Dot(x, y, pixel.Warn)
		case 5:
			g.Disc(x, y, 1, pixel.Warn)
		case 6:
			g.Disc(x, y, 2, pixel.Warn)
			g.Dot(x-1, y-1, light(k))
		case 7:
			g.Circle(x, y, 3, pixel.Warn)
			g.Dot(x-1, y-1, light(k))
		case 8:
			for _, d := range [][2]int{{-4, 0}, {4, 0}, {0, -4}, {0, 4}} {
				g.Dot(x+d[0], y+d[1], pixel.Warn)
			}
		}
	}
}
