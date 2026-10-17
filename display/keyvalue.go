package display

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type KeyValueOptions struct {
	// Slot inks the value; zero means Ink.
	Slot pixel.Slot
}

// KeyValue joins key and value with a dotted leader.
func KeyValue(c *ui.Context, key, value string, o KeyValueOptions) ui.Element {
	k := pixel.K(c)
	kw, vw := pixel.Measure(key), pixel.Measure(value)
	return ui.Box(c).FillWidth().MinWidth(k.D(kw + vw + 8)).Height(k.D(pixel.LineHeight)).Role(ui.RoleText).Label(key + " " + value).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Text(0, 0, key, pixel.Faint, 1)
		g.Checker(kw+3, pixel.LineHeight-5, g.W-kw-vw-6, 1, pixel.Faint)
		g.Text(g.W-vw, 0, value, o.Slot, 1)
	})
}
