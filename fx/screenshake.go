package fx

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type ScreenShakeOptions struct {
	// Cue shakes the content each time it changes.
	Cue int
	// Amount is the widest jolt in dots; zero means 3.
	Amount int
}

// quakes are the jolts, a step each, in Amount thirds.
var quakes = [...][2]int{{3, -2}, {-3, 2}, {2, 3}, {-3, -2}, {2, 1}, {-1, -2}, {1, 1}, {-1, 0}, {0, 0}}

// ScreenShake jolts its content both ways, settling as it goes.
func ScreenShake(c *ui.Context, o ScreenShakeOptions, content func()) ui.Element {
	k := pixel.K(c)
	a := pixel.Or("ScreenShake Amount", o.Amount, 3)
	e := ui.Box(c)
	q := quakes[k.Shot(c, e, o.Cue, 15, len(quakes)-1)]
	dx, dy := q[0]*a/3, q[1]*a/3
	return e.Padding(k.D(a+dy), k.D(a-dx), k.D(a-dy), k.D(a+dx)).Children(content)
}
