package layout

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ViewportOptions struct {
	// W×H is the virtual screen in dots; zero means 160×144.
	W, H int
}

// Viewport fills its parent and draws children at W×H dots.
func Viewport(c *ui.Context, o ViewportOptions, children func()) ui.Element {
	k := pixel.K(c)
	vw, vh := pixel.Or("Viewport W", o.W, 160), pixel.Or("Viewport H", o.H, 144)
	base := int(k.Dot())
	outer := ui.Box(c).Fill()
	was := outer.Bounds()
	d := fit(vw, vh, was.W, was.H, base)
	// Offsets stay whole outer dots, so both grids line up.
	x := (int(was.W)/base - vw*d/base) / 2 * base
	y := (int(was.H)/base - vh*d/base) / 2 * base
	return outer.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W != was.W || r.H != was.H {
			c.Invalidate()
		}
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Back)
	}).Children(func() {
		s := pixel.Current(c)
		s.Dot = d
		pixel.Scope(c, s, func() {
			ui.Box(c).Absolute().Left(float32(max(0, x))).Top(float32(max(0, y))).Size(float32(vw*d), float32(vh*d)).Clip().Children(children)
		})
	})
}

// fit picks the largest fitting multiple of base, ≤ 4.
func fit(w, h int, roomW, roomH float32, base int) int {
	d := base
	for d+base <= 4 && float32(w*(d+base)) <= roomW && float32(h*(d+base)) <= roomH {
		d += base
	}
	return d
}
