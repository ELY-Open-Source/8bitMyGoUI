package text

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type ShakeTextOptions struct {
	Slot  pixel.Slot
	Scale int
	Seed  int
}

// ShakeText jitters each character a dot, seeded, at 12 fps.
func ShakeText(c *ui.Context, s string, o ShakeTextOptions) ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	return block(c, k, s, pixel.Measure(s)*sc+2, pixel.LineHeight*sc+2, func(g pixel.Canvas) {
		step := g.Step(12)
		each(s, sc, func(i int, r rune, x int) {
			dx, dy := shake(k.Reduce(), o.Seed, i, step)
			g.Text(1+x+dx, 1+dy, string(r), o.Slot, sc)
		})
	})
}

// shake is a −1…1 offset, or none under Reduce Motion.
func shake(reduce bool, seed, i, step int) (dx, dy int) {
	if reduce {
		return 0, 0
	}
	n := noise(seed, i, step)
	return int(n%3) - 1, int(n/3%3) - 1
}
