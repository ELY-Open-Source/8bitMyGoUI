package layout

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CRTFrameOptions struct {
	NoScanlines bool
	// Brand is micro text on the chin.
	Brand string
}

const crtSide, crtChin = 10, 16

// CRTFrame shows children on a monitor with scanlines.
func CRTFrame(c *ui.Context, o CRTFrameOptions, children func()) ui.Element {
	k := pixel.K(c)
	return ui.Column(c).Padding(k.D(crtSide), k.D(crtSide), k.D(crtChin), k.D(crtSide)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		box(g, pixel.Raised, pixel.Face)
		if g.W > 2*crtSide && g.H > crtSide+crtChin {
			g.Frame(crtSide-3, crtSide-3, g.W-2*crtSide+6, g.H-crtSide-crtChin+6, pixel.Sunken, pixel.Well)
		}
		g.Micro(crtSide, g.H-crtChin/2-2, o.Brand, pixel.Shade, 1)
		g.Rect(g.W-crtSide-6, g.H-crtChin/2-1, 3, 3, pixel.Good)
	}).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if o.NoScanlines {
			return
		}
		// Scanlines are Shade at half strength on every other row.
		g := k.Canvas(p, r)
		for y := crtSide + 1; y < g.H-crtChin; y += 2 {
			g.Checker(crtSide, y, g.W-2*crtSide, 1, pixel.Shade)
		}
	}).Children(children)
}
