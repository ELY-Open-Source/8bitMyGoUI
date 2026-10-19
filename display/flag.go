package display

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type FlagOptions struct {
	// Slot is the cloth; zero means Bad.
	Slot  pixel.Slot
	Scale int
	Label string
}

// Flag waves on its pole, stepped at 6 fps.
func Flag(c *ui.Context, o FlagOptions) ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	cloth := slotOr(o.Slot, pixel.Bad)
	return art(c, k, o.Label, 20*sc, 26*sc, func(g pixel.Canvas) {
		step := g.Step(6)
		g.Rect(sc, 2*sc, sc, 24*sc, pixel.Edge)
		g.Rect(0, 0, 3*sc, 2*sc, pixel.Warn)
		for col := range 16 {
			dy := [4]int{0, 1, 1, 0}[(step+col/4)%4]
			x, y := (2+col)*sc, (3+dy)*sc
			g.Rect(x, y, sc, 10*sc, cloth)
			if (step+col/4)%4 == 2 {
				g.Rect(x, y, sc, 10*sc, pixel.Shade)
				g.Veil(x, y, sc, 10*sc, cloth, 8)
			}
			g.Rect(x, y, sc, sc, pixel.Lit)
		}
	})
}
