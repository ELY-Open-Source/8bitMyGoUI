package layout

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type ScreenOptions struct {
	// Vignette dithers Shade into the edges.
	Vignette bool
	Padding  int
}

// Screen fills the window with Back; children form a column.
func Screen(c *ui.Context, o ScreenOptions, children func()) ui.Element {
	k := pixel.K(c)
	pad := pixel.Or("Screen Padding", o.Padding, 0)
	return ui.Column(c).Fill().Padding(k.D(pad)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Back)
		if o.Vignette {
			vignette(g)
		}
	}).Children(children)
}

// vignette darkens four bands toward the edges.
func vignette(g pixel.Canvas) {
	band := max(2, min(g.W, g.H)/32)
	for i := range 4 {
		in, lv := i*band, 8-2*i
		g.Veil(in, in, g.W-2*in, band, pixel.Shade, lv)
		g.Veil(in, g.H-in-band, g.W-2*in, band, pixel.Shade, lv)
		g.Veil(in, in+band, band, g.H-2*in-2*band, pixel.Shade, lv)
		g.Veil(g.W-in-band, in+band, band, g.H-2*in-2*band, pixel.Shade, lv)
	}
}
