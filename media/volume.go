package media

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type VolumeOptions struct {
	// Bars is how many steps; zero means 8.
	Bars     int
	Disabled bool
}

type VolumeResult struct {
	Element ui.Element
	Changed bool
}

// Volume sets *level in stepped bars; the speaker mutes.
func Volume(c *ui.Context, text string, level *int, muted *bool, o VolumeOptions) VolumeResult {
	k := pixel.K(c)
	n := pixel.Or("Volume Bars", o.Bars, 8)
	if n < 2 || *level < 0 || *level > n {
		panic(fmt.Sprintf("media: Volume level %d of %d bars", *level, n))
	}
	var res VolumeResult
	res.Element = ui.Row(c).Gap(k.D(3)).AlignItems(ui.Center).Children(func() {
		sp := icon.Volume
		if *muted {
			sp = icon.Mute
		}
		mute := input.Toggle(c, "", muted, input.ToggleOptions{Icon: sp, Label: k.Msg("mute"), Disabled: o.Disabled}).Changed
		v := float64(*level)
		shell, st, moved := input.SlideBase(c, text, &v, input.RailOptions{Max: float64(n), Step: 1, Disabled: o.Disabled})
		*level = int(v)
		lit, quiet := *level, *muted
		shell.Size(k.D(4*n), k.D(14)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range n {
				h := 4 + i*10/(n-1)
				s := pixel.Well
				switch {
				case i >= lit:
				case quiet:
					s = pixel.Faint
				case i >= n*3/4:
					s = pixel.Warn
				default:
					s = pixel.Good
				}
				g.Frame(4*i, 14-h, 3, h, pixel.Flat, s)
			}
			g.Finish(st)
		})
		res.Changed = mute || moved
	})
	if res.Changed {
		slog.Debug("8bit: volume", "label", text, "level", *level, "muted", *muted)
	}
	return res
}
