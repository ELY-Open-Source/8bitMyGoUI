package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ConfettiOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// confettiSlots color the bits in turn.
var confettiSlots = []pixel.Slot{pixel.Accent, pixel.Good, pixel.Warn, pixel.Bad, pixel.Info, pixel.Magic}

// Confetti tumbles bits of every color down, turning.
func Confetti(c *ui.Context, o ConfettiOptions) *ui.Element {
	k := pixel.K(c)
	w, h := size("Confetti", o.Width, o.Height, 160, 96)
	return box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { confetti(g, g.Step(10)) })
	})
}

// confetti turns each bit flat, edge-on, upright, edge-on.
func confetti(g pixel.Canvas, step int) {
	for i := range g.W * g.H / 90 {
		speed := 1 + rnd(i, 90)%2
		y := wrap(rnd(i, 91)+step*speed, g.H+4) - 2
		x := wrap(rnd(i, 92)+sway(step+rnd(i, 93), 12, 1), g.W)
		s := confettiSlots[i%len(confettiSlots)]
		switch wrap(step+i, 4) {
		case 0:
			g.Rect(x, y, 2, 1, s)
		case 2:
			g.Rect(x, y, 1, 2, s)
		default:
			g.Dot(x, y, s)
		}
	}
}
