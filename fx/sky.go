package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SkyOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// Sky dithers down the palette's sky slots; night shows stars.
func Sky(c *ui.Context, o SkyOptions) *ui.Element {
	k := pixel.K(c)
	w, h := size("Sky", o.Width, o.Height, 160, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { sky(g, k, g.Step(4)) })
	})
}

// sky fills g in bands cut at height percents.
func sky(g pixel.Canvas, k pixel.Kit, step int) {
	bands, cuts := []pixel.Slot{pixel.Info, pixel.AccentLit, pixel.Lit}, []int{0, 55, 85, 100}
	if k.Palette().Night {
		bands, cuts = []pixel.Slot{pixel.Well, pixel.Back, pixel.Face, pixel.Lit}, []int{0, 35, 60, 82, 100}
	}
	for i, s := range bands {
		y0, y1 := g.H*cuts[i]/100, g.H*cuts[i+1]/100
		g.Rect(0, y0, g.W, y1-y0, s)
	}
	// Seams dither five rows from one band into the next.
	for i := 1; i < len(bands); i++ {
		y := g.H * cuts[i] / 100
		for j := range 5 {
			g.Dither(0, y-2+j, g.W, 1, bands[i-1], bands[i], 2+3*j)
		}
	}
	if k.Palette().Night {
		stars(g, k, step)
	}
}

// stars twinkle over the top two thirds.
func stars(g pixel.Canvas, k pixel.Kit, step int) {
	tones := []pixel.Slot{light(k), light(k), pixel.AccentLit, light(k), pixel.Warn}
	for i := range g.W * g.H / 90 {
		x, y := rnd(i, 1)%g.W, rnd(i, 2)%max(1, g.H*2/3)
		s, t := tones[i%len(tones)], wrap(step+rnd(i, 3), 8)
		switch {
		case t == 0:
			g.Dot(x, y, pixel.Faint)
		case i%9 == 0 && t < 3:
			g.Rect(x-1, y, 3, 1, s)
			g.Rect(x, y-1, 1, 3, s)
		default:
			g.Dot(x, y, s)
		}
	}
}
