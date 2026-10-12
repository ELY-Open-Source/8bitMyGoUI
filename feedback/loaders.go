package feedback

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DotsLoaderOptions struct{}

// DotsLoader bounces three dots in turn.
func DotsLoader(c *ui.Context, text string, o DotsLoaderOptions) *ui.Element {
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

type BarsLoaderOptions struct{}

// heights rise and fall; each bar runs a step behind.
var heights = [8]int{3, 6, 9, 12, 9, 6, 3, 3}

// BarsLoader raises and drops five bars.
func BarsLoader(c *ui.Context, text string, o BarsLoaderOptions) *ui.Element {
	k := pixel.K(c)
	return ui.Box(c).Role(ui.RoleProgress).Label(text).Size(k.D(23), k.D(12)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		f := g.Busy(10)
		for i := range 5 {
			h := heights[(f+8-i)%8]
			g.Rect(i*5, 12-h, 3, h, pixel.Accent)
		}
	})
}
