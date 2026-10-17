package media

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TunerOptions struct {
	Disabled bool
}

type TunerResult struct {
	Element ui.Element
	Changed bool
}

// Tuner sets *mhz on the FM dial by tenths.
func Tuner(c *ui.Context, text string, mhz *float64, o TunerOptions) TunerResult {
	k := pixel.K(c)
	const lo, hi = 87.5, 108.0
	if !(*mhz >= lo && *mhz <= hi) {
		panic(fmt.Sprintf("media: Tuner %v MHz outside FM", *mhz))
	}
	var res TunerResult
	// The glass holds the scale; the scale is the slider.
	res.Element = ui.Box(c).Size(k.D(176), k.D(32)).PaddingX(k.D(8)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, 176, 32, pixel.Sunken, pixel.Well)
		g.Dither(2, 2, 172, 28, pixel.Well, pixel.Warn, 3)
		g.HLine(8, 24, 162, pixel.Ink)
		g.Micro(150, 26, "MHZ", pixel.Faint, 1)
	}).Children(func() {
		shell, st, changed := input.SlideBase(c, text, mhz, input.RailOptions{Min: lo, Max: hi, Step: 0.1, Disabled: o.Disabled})
		res.Changed = changed
		at := *mhz
		shell.Size(k.D(160), k.D(32)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			// The dot under the pointer is the frequency it sets.
			x := func(f float64) int { return int((f - lo) / (hi - lo) * 160) }
			for f := 88; f <= 108; f++ {
				h := 3
				if f%4 == 0 {
					h = 6
					s := fmt.Sprint(f)
					g.Micro(x(float64(f))-pixel.MicroWidth(s)/2, 4, s, pixel.Ink, 1)
				}
				g.VLine(x(float64(f)), 24-h, h, pixel.Ink)
			}
			g.Rect(x(at), 3, 2, 26, pixel.Bad)
			g.Finish(st)
		})
	})
	if res.Changed {
		slog.Debug("8bit: tuner", "label", text, "mhz", *mhz)
	}
	return res
}
