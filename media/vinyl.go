package media

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type VinylOptions struct {
	Playing bool
}

// Vinyl turns a record while playing; the arm swings on.
func Vinyl(c *ui.Context, text string, o VinylOptions) *ui.Element {
	k := pixel.K(c)
	e := art(c, k, text, 72, 64)
	swing := float64(k.Toward(c, e, o.Playing, 12, 6)) / 6
	a := float64(turned(c, k, e, o.Playing, 8)%8) * math.Pi / 4
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		body, line := dark(k)
		g.Disc(32, 32, 29, body)
		for _, rr := range []int{25, 21, 17, 13} {
			g.Circle(32, 32, rr, line)
		}
		// Light stays put while the record turns under it.
		g.Line(14, 17, 18, 13, pixel.Lit)
		g.Line(46, 51, 50, 47, pixel.Lit)
		g.Disc(32, 32, 9, pixel.Accent)
		g.Rect(31+int(math.Round(5*math.Cos(a))), 31+int(math.Round(5*math.Sin(a))), 2, 2, pixel.AccentLit)
		g.Disc(32, 32, 1, pixel.Lit)
		// The arm swings from rest, 90°, onto grooves, 118°.
		arm := (90 + 28*swing) * math.Pi / 180
		hx, hy := 66+int(math.Round(30*math.Cos(arm))), 6+int(math.Round(30*math.Sin(arm)))
		g.Line(66, 6, hx, hy, pixel.Edge)
		g.Rect(hx-1, hy-1, 3, 3, pixel.Edge)
		g.Disc(66, 6, 3, pixel.Edge)
		g.Disc(66, 6, 2, pixel.Face)
	})
}
