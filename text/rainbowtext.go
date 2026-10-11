package text

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var rainbow = []pixel.Slot{pixel.Bad, pixel.Warn, pixel.Good, pixel.Info, pixel.Magic}

type RainbowTextOptions struct {
	// Band is zero for Bad, Warn, Good, Info, Magic.
	Band  []pixel.Slot
	Scale int
}

// RainbowText gives each character the next slot of a band.
func RainbowText(c *ui.Context, s string, o RainbowTextOptions) *ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	band := o.Band
	if len(band) == 0 {
		band = rainbow
	}
	return block(c, k, s, pixel.Measure(s)*sc, pixel.LineHeight*sc, func(g pixel.Canvas) {
		cur := pixel.Cycle(band, g.Step(8))
		each(s, sc, func(i int, r rune, x int) {
			g.Text(x, 0, string(r), cur[i%len(cur)], sc)
		})
	})
}
