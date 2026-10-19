package fx

import (
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type SparkleOptions struct {
	// Size is the square's side in dots; zero means 24.
	Size int
	// Slot inks the stars; zero means Warn.
	Slot  pixel.Slot
	Label string
}

// glints are a star's frames by step; -1 is dark.
var glints = [...]int{0, 1, 2, 1, 0, -1, -1, -1}

// Sparkle twinkles four stars around its middle, in turn.
func Sparkle(c *ui.Context, o SparkleOptions) ui.Element {
	k := pixel.K(c)
	n := pixel.Or("Sparkle Size", o.Size, 24)
	s := slotOr(o.Slot, pixel.Warn)
	return box(c, k, o.Label, n, n).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { sparkle(g, s, g.Step(10)) })
	})
}

// sparkle moves each star to a new spot every cycle.
func sparkle(g pixel.Canvas, s pixel.Slot, step int) {
	for i := range 4 {
		t := step + 2*i
		f := glints[wrap(t, len(glints))]
		if f < 0 {
			continue
		}
		seed := t/len(glints)*4 + i
		a := float64(i)*math.Pi/2 + float64(rnd(seed, 170)%60-30)*math.Pi/180
		d := float64(g.W) * float64(20+rnd(seed, 171)%20) / 100
		x := g.W/2 + int(math.Round(d*math.Cos(a))) - 2
		y := g.H/2 + int(math.Round(d*math.Sin(a))) - 2
		g.Sprite(x, y, twinkle, f, s, 1)
	}
}
