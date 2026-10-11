package layout

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type AspectBoxOptions struct {
	// W:H is the ratio; zero sides mean 1.
	W, H int
}

// AspectBox fills the width and keeps W:H in whole dots.
func AspectBox(c *ui.Context, o AspectBoxOptions, children func()) *ui.Element {
	k := pixel.K(c)
	rw, rh := orDef("AspectBox W", o.W, 1), orDef("AspectBox H", o.H, 1)
	e := ui.Box(c).FillWidth().Clip()
	w, check := pixel.WrapWidth(c, e)
	return e.Height(k.D(w * rh / rw)).Draw(func(p *ui.Painter, r ui.Rect) { check(r) }).Children(children)
}
