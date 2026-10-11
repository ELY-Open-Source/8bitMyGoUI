package display

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var cup = pixel.MustSprite(`
	..eeeeeeeeeee..
	eeeyyyyyyyyyeee
	e.eyllyyyyyye.e
	e.eylyyyyyyye.e
	.eeyyyyyyyyyee.
	...eyyyyyyye...
	....eyyyyye....
	.....eyyye.....
	......eye......
	......eye......
	.....eyyye.....
	...eesssssee...
	...esssssssse..
`)

type TrophyOptions struct {
	Scale int
}

// Trophy is a gold cup over an engraved plaque.
func Trophy(c *ui.Context, plaque string, o TrophyOptions) *ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	up := strings.ToUpper(plaque)
	micro := pixel.MicroOK(up)
	tw, th := pixel.Measure(plaque), pixel.LineHeight
	if micro {
		tw, th = pixel.MicroWidth(up), pixel.MicroHeight
	}
	pw, ph := tw+6, th+4
	w := max(cup.W*sc, pw)
	return art(c, k, plaque, w, cup.H*sc+1+ph, func(g pixel.Canvas) {
		g.Sprite((w-cup.W*sc)/2, 0, cup, 0, pixel.Ink, sc)
		g.Frame((w-pw)/2, cup.H*sc+1, pw, ph, pixel.Raised, pixel.Shade)
		if micro {
			g.Micro((w-pw)/2+3, cup.H*sc+3, up, pixel.Ink, 1)
		} else {
			g.Text((w-pw)/2+3, cup.H*sc+3, plaque, pixel.Ink, 1)
		}
	})
}
