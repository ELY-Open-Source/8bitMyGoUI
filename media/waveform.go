package media

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type WaveformOptions struct {
	// Height is in dots; zero means 32.
	Height int
}

// Waveform draws samples, ±1, as lines from the middle.
func Waveform(c *ui.Context, text string, samples []float64, o WaveformOptions) *ui.Element {
	k := pixel.K(c)
	if len(samples) == 0 {
		panic("media: Waveform has no samples")
	}
	for _, s := range samples {
		if !(s >= -1 && s <= 1) {
			panic(fmt.Sprintf("media: Waveform sample %v outside ±1", s))
		}
	}
	h := pixel.Or("Waveform Height", o.Height, 32)
	half := (h - 1) / 2
	return art(c, k, text, 2*len(samples)-1, h).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for x := 1; x < g.W; x += 2 {
			g.Dot(x, half, pixel.Faint)
		}
		for i, s := range samples {
			n := int(math.Round(s * float64(half)))
			g.VLine(2*i, min(half, half-n), max(n, -n)+1, pixel.Accent)
		}
	})
}
