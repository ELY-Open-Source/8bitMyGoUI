package text

import (
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type WaveTextOptions struct {
	Slot  pixel.Slot
	Scale int
	// Amp is the bob in dots; zero means 2.
	Amp int
}

// WaveText bobs characters on a sine wave at 10 fps.
func WaveText(c *ui.Context, s string, o WaveTextOptions) ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	amp := pixel.Or("WaveText Amp", o.Amp, 2)
	return block(c, k, s, pixel.Measure(s)*sc, pixel.LineHeight*sc+2*amp, func(g pixel.Canvas) {
		step := g.Step(10)
		each(s, sc, func(i int, r rune, x int) {
			g.Text(x, amp+wave(amp, step+i), string(r), o.Slot, sc)
		})
	})
}

// wave is the bob at step n of an 8-cycle.
func wave(amp, n int) int {
	return int(math.Round(float64(amp) * math.Sin(float64(n%8)*math.Pi/4)))
}
