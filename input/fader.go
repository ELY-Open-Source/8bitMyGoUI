package input

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type FaderOptions struct {
	Min, Max float64
	// Height is the travel in dots; zero means 64.
	Height   int
	Disabled bool
}

type FaderResult = SliderResult

// Fader is a vertical slider with a cap and ticks.
func Fader(c *ui.Context, text string, value *float64, o FaderOptions) FaderResult {
	k := pixel.K(c)
	lo, hi := bounds(o.Min, o.Max)
	h := pixel.Or("Fader Height", o.Height, 64)
	shell, st, changed := slide(c, value, rail{lo: lo, hi: hi, vertical: true, label: text, disabled: o.Disabled})
	shell.Size(k.D(16), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for t := 0; t <= 4; t++ {
			y := 3 + t*(g.H-7)/4
			g.HLine(0, y, 3, pixel.Faint)
			g.HLine(13, y, 3, pixel.Faint)
		}
		g.Frame(6, 0, 4, g.H, pixel.Sunken, pixel.Well)
		y := g.H - 7 - knobAt(*value, lo, hi, g.H-7)
		g.Rect(7, y+3, 2, g.H-y-4, pixel.Accent)
		g.Frame(2, y, 12, 7, pixel.Raised, pixel.Face)
		g.HLine(4, y+3, 8, pixel.Edge)
		g.Finish(st)
	})
	return report(shell, changed, "fader", text, *value)
}
