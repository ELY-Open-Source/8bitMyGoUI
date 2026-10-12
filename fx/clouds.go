package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CloudsOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// Clouds drift by: far ones slowly, near ones faster.
func Clouds(c *ui.Context, o CloudsOptions) *ui.Element {
	k := pixel.K(c)
	w, h := size("Clouds", o.Width, o.Height, 160, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { clouds(g, g.Step(8)) })
	})
}

// clouds drifts two layers across the top third of g.
func clouds(g pixel.Canvas, step int) {
	for layer, sp := range []*pixel.Sprite{smallCloud, bigCloud} {
		// Far clouds shift a dot each four steps; near, two.
		shift, span := step/(4>>layer), g.W+sp.W
		for i := range max(1, g.W/(120-20*layer)) {
			x := wrap(rnd(i, 10+layer)-shift, span) - sp.W
			y := rnd(i, 12+layer)%max(1, g.H/3) + 4*layer
			g.Sprite(x, y, sp, 0, pixel.Ink, 1)
		}
	}
}
