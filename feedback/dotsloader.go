package feedback

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type DotsLoaderOptions struct{}

// DotsLoader bounces three dots in turn.
func DotsLoader(c *ui.Context, text string, o DotsLoaderOptions) ui.Element {
	k := pixel.K(c)
	return ui.Box(c).Role(ui.RoleProgress).Label(text).Size(k.D(18), k.D(8)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		up := g.Busy(6) % 4
		for i := range 3 {
			y, ink := 4, pixel.Faint
			if i == up {
				y, ink = 0, pixel.Accent
			}
			g.Rect(i*7, y, 4, 4, ink)
		}
	})
}
