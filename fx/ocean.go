package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type OceanOptions struct {
	// Width and Height are dots; zeros mean 160 by 64.
	Width, Height int
	Label         string
}

// Ocean rolls wave crests toward you by cycling their slots.
func Ocean(c *ui.Context, o OceanOptions) ui.Element {
	k := pixel.K(c)
	w, h := size("Ocean", o.Width, o.Height, 160, 64)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { ocean(g, k, g.Step(6)) })
	})
}

// ocean's crest rows grow sparser and longer downward.
func ocean(g pixel.Canvas, k pixel.Kit, step int) {
	g.Rect(0, 0, g.W, g.H, pixel.Info)
	g.Dither(0, 0, g.W, 2, pixel.Info, light(k), 6)
	// Crest row j shows the band turned back each step.
	band := pixel.Cycle([]pixel.Slot{pixel.Info, pixel.Info, pixel.AccentLit, light(k)}, -step)
	for y, j := 3, 0; y < g.H; j++ {
		if s := band[j%len(band)]; s != pixel.Info {
			for i := range max(1, g.W/12) {
				g.HLine(rnd(i, 120+j)%g.W, y, 2+j/2, s)
			}
		}
		y += 2 + j/2
	}
}
