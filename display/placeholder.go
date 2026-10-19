package display

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type PlaceholderOptions struct {
	// W×H is zero for 64×48 dots.
	W, H int
}

// Placeholder is a dithered box with a cross and size.
func Placeholder(c *ui.Context, o PlaceholderOptions) ui.Element {
	k := pixel.K(c)
	w, h := pixel.Or("Placeholder W", o.W, 64), pixel.Or("Placeholder H", o.H, 48)
	size := fmt.Sprintf("%d×%d", w, h)
	return art(c, k, size, w, h, func(g pixel.Canvas) {
		if w < 3 || h < 3 {
			g.Dither(0, 0, w, h, pixel.Face, pixel.Shade, 4)
			return
		}
		g.Frame(0, 0, w, h, pixel.Flat, pixel.Face)
		g.Dither(1, 1, w-2, h-2, pixel.Face, pixel.Shade, 4)
		g.Line(1, 1, w-2, h-2, pixel.Shade)
		g.Line(w-2, 1, 1, h-2, pixel.Shade)
		tw := pixel.Measure(size)
		if tw+4 <= w && pixel.LineHeight <= h {
			g.Rect((w-tw)/2-2, (h-pixel.LineHeight)/2, tw+4, pixel.LineHeight, pixel.Face)
			g.Text((w-tw)/2, (h-pixel.LineHeight)/2, size, pixel.Faint, 1)
		}
	})
}
