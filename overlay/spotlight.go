package overlay

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type SpotlightOptions struct{}

// Spotlight dithers all but target and captions it.
func Spotlight(c *ui.Context, target ui.Element, caption string, o SpotlightOptions) ui.Element {
	k := pixel.K(c)
	b := target.Bounds()
	deep, _ := dark(k)
	var veil ui.Element
	ui.Overlay(c, func() {
		veil = ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).PassThrough().Role(ui.RoleStatus).Label(caption)
		veil.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			// The hole is the target, two dots wider all round.
			x, y := int(b.X/k.Dot())-2, int(b.Y/k.Dot())-2
			w, h := int(b.W/k.Dot())+4, int(b.H/k.Dot())+4
			g.Veil(0, 0, g.W, y, deep, 12)
			g.Veil(0, y+h, g.W, g.H-y-h, deep, 12)
			g.Veil(0, y, x, h, deep, 12)
			g.Veil(x+w, y, g.W-x-w, h, deep, 12)
			g.FocusMark(x+2, y+2, w-4, h-4)
		})
	})
	w := pixel.Measure(caption) + 16
	bubble(c, ui.Rect{X: b.X, Y: b.Y - k.D(2), W: b.W, H: b.H + k.D(4)}, func() {
		pixel.Text(c, caption, pixel.TextOptions{})
	}, w)
	return veil
}
