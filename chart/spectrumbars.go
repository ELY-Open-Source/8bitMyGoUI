package chart

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SpectrumBarsOptions struct {
	// Segments is each bar's height in blocks; zero means 10.
	Segments int
}

type peaksKey struct{}

// peaks are the falling dots, when built last, whether falling.
type peaks struct {
	at     []float64
	last   time.Time
	moving bool
}

// SpectrumBars lights segmented bars; peak dots fall behind them.
func SpectrumBars(c *ui.Context, text string, levels []float64, o SpectrumBarsOptions) *ui.Element {
	k := pixel.K(c)
	n := pixel.Or("SpectrumBars Segments", o.Segments, 10)
	for _, l := range levels {
		if !(l >= 0 && l <= 1) {
			panic(fmt.Sprintf("chart: SpectrumBars level %v outside 0–1", l))
		}
	}
	e := ui.Box(c).Size(k.D(len(levels)*6), k.D(n*3)).Role(ui.RoleImage).Label(text).Description(nums(levels))
	pk := ui.Local(e, peaksKey{}, func() peaks { return peaks{} })
	if len(pk.at) != len(levels) {
		pk.at = append([]float64(nil), levels...)
	}
	now := k.Now()
	// A peak falls a bar a second; reduced, it rests.
	fall := now.Sub(pk.last).Seconds()
	switch {
	case pk.last.IsZero() || k.Reduce():
		fall = 1
	case !pk.moving:
		// Resting peaks start falling now, not at the last build.
		fall = 0
	}
	pk.last, pk.moving = now, false
	for i, l := range levels {
		pk.at[i] = max(l, pk.at[i]-fall)
		pk.moving = pk.moving || pk.at[i] > l
	}
	if pk.moving {
		k.Tick(c, 20)
	}
	shown := append([]float64(nil), pk.at...)
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, l := range levels {
			lit := int(l*float64(n) + 0.5)
			for s := range n {
				y := (n - 1 - s) * 3
				ink := pixel.Good
				switch {
				case s >= n*8/10:
					ink = pixel.Bad
				case s >= n*6/10:
					ink = pixel.Warn
				}
				if s < lit {
					g.Rect(i*6, y, 5, 2, ink)
				} else {
					g.Checker(i*6, y, 5, 2, pixel.Shade)
				}
			}
			if top := int(shown[i]*float64(n) + 0.5); top > 0 {
				g.Rect(i*6, (n-top)*3, 5, 1, pixel.Ink)
			}
		}
	})
}
