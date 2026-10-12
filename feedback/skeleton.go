package feedback

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SkeletonOptions struct {
	// Lines is the text lines; zero means 3.
	Lines int
	// Avatar puts a square before the lines.
	Avatar bool
	// Width is in dots; zero means 160.
	Width int
}

// Skeleton stands in for content; a dithered shine sweeps it.
func Skeleton(c *ui.Context, o SkeletonOptions) *ui.Element {
	k := pixel.K(c)
	n := pixel.Or("Skeleton Lines", o.Lines, 3)
	w := pixel.Or("Skeleton Width", o.Width, 160)
	h := n*12 - 4
	return ui.Box(c).Role(ui.RoleStatus).Label(k.Msg("loading")).Size(k.D(w), k.D(max(h, 24))).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		x := 0
		var shapes [][4]int
		if o.Avatar {
			shapes, x = append(shapes, [4]int{0, 0, 24, 24}), 30
		}
		for i := range n {
			lw := w - x
			if i == n-1 && n > 1 {
				lw = lw * 3 / 5
			}
			shapes = append(shapes, [4]int{x, i * 12, lw, 8})
		}
		shine := g.Step(8)*6%(w+48) - 24
		for _, s := range shapes {
			g.Dither(s[0], s[1], s[2], s[3], pixel.Shade, pixel.Face, 4)
			g.Clip(s[0], s[1], s[2], s[3], func() { g.Veil(shine, s[1], 16, s[3], pixel.Lit, 6) })
		}
	})
}
