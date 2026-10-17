package media

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type FilmstripOptions struct {
	// Scale is dots per sprite dot, 1–4; zero means 2.
	Scale int
}

// Filmstrip lays a sprite's frames on film with sprocket holes.
func Filmstrip(c *ui.Context, text string, sp *pixel.Sprite, o FilmstripOptions) ui.Element {
	k := pixel.K(c)
	s := scale("Filmstrip", o.Scale, 2)
	n := sp.Frames()
	fw, fh := sp.W*s+4, sp.H*s+2
	w, h := n*fw+4, fh+14
	return art(c, k, text, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		body, line := dark(k)
		g.Rect(0, 0, w, h, body)
		for x := 2; x+3 <= w; x += 6 {
			g.Rect(x, 2, 3, 3, line)
			g.Rect(x, h-5, 3, 3, line)
		}
		for i := range n {
			x := 3 + i*fw
			g.Rect(x, 7, fw-2, fh, pixel.Lit)
			g.Sprite(x+1, 8, sp, i, pixel.Ink, s)
		}
	})
}
