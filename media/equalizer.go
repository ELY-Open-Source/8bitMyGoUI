package media

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type EqualizerOptions struct {
	// Height is each fader's travel; zero means 64.
	Height   int
	Disabled bool
}

type EqualizerResult struct {
	Element *ui.Element
	// Band is the band changed this frame, else -1.
	Band int
}

// Equalizer sets each band's gain, ±12 dB, on faders.
func Equalizer(c *ui.Context, text string, bands []string, gains []float64, o EqualizerOptions) EqualizerResult {
	k := pixel.K(c)
	if len(bands) == 0 || len(bands) != len(gains) {
		panic(fmt.Sprintf("media: Equalizer has %d bands and %d gains", len(bands), len(gains)))
	}
	for i, g := range gains {
		if !(g >= -12 && g <= 12) {
			panic(fmt.Sprintf("media: Equalizer band %q gain %v outside ±12", bands[i], g))
		}
	}
	res := EqualizerResult{Band: -1}
	res.Element = ui.Row(c).Gap(k.D(6)).Role(ui.RoleGroup).Label(text).Children(func() {
		for i, b := range bands {
			ui.Column(c).Gap(k.D(2)).AlignItems(ui.Center).Children(func() {
				if input.Fader(c, b, &gains[i], input.FaderOptions{Min: -12, Max: 12, Height: o.Height, Disabled: o.Disabled}).Changed {
					res.Band = i
					slog.Debug("8bit: equalizer", "label", text, "band", b, "gain", gains[i])
				}
				ui.Box(c).Size(k.D(pixel.MicroWidth(b)), k.D(5)).Draw(func(p *ui.Painter, r ui.Rect) {
					k.Canvas(p, r).Micro(0, 0, b, pixel.Faint, 1)
				})
			})
		}
	})
	return res
}
