package display

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type Rank uint8

const (
	Gold Rank = iota
	Silver
	Bronze
)

var ribbon = pixel.MustSprite(`
	aaa.....ccc
	aaaa...cccc
	.aaaa.cccc.
	..aaaccccc.
	...aaccc...
	....acc....
`)

type MedalOptions struct {
	Scale int
}

// Medal is a ribboned medal showing its rank.
func Medal(c *ui.Context, r Rank, o MedalOptions) ui.Element {
	if r > Bronze {
		panic(fmt.Sprintf("display: unknown Rank %d", r))
	}
	k, sc := pixel.K(c), scale(o.Scale)
	name := k.Msg([...]string{"gold", "silver", "bronze"}[r])
	return art(c, k, name, 11*sc, 18*sc, func(g pixel.Canvas) {
		g.Sprite(0, 0, ribbon, 0, pixel.Ink, sc)
		g.Disc(5*sc, 12*sc, 5*sc, pixel.Edge)
		g.Disc(5*sc, 12*sc, 5*sc-1, pixel.Shade)
		metal(g, r, 5*sc, 12*sc, 5*sc-2)
		g.Micro(5*sc-1, 12*sc-2, itoa(int(r)+1), pixel.Edge, 1)
	})
}

// metal fills a disc in the rank's colors.
func metal(g pixel.Canvas, r Rank, cx, cy, rad int) {
	switch r {
	case Gold:
		g.Disc(cx, cy, rad, pixel.Warn)
	case Silver:
		g.Disc(cx, cy, rad, pixel.Face)
	case Bronze:
		g.Disc(cx, cy, rad, pixel.Warn)
		g.Veil(cx-rad, cy-rad, 2*rad+1, 2*rad+1, pixel.Bad, 8)
		g.Circle(cx, cy, rad, pixel.Shade)
	}
	g.Dot(cx-rad/2, cy-rad/2, pixel.Lit)
}
