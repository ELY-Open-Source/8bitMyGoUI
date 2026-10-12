package media

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ScrubberOptions struct {
	// Width is the bar in dots; zero means 96.
	Width    int
	Disabled bool
}

type ScrubberResult struct {
	Element *ui.Element
	Changed bool
}

// Scrubber seeks *at within total, the time at either end.
func Scrubber(c *ui.Context, text string, at *time.Duration, total time.Duration, o ScrubberOptions) ScrubberResult {
	k := pixel.K(c)
	if total <= 0 || *at < 0 || *at > total {
		panic(fmt.Sprintf("media: Scrubber at %v of %v", *at, total))
	}
	// Seconds, so readers hear the place, not a fraction.
	secs := at.Seconds()
	var res ScrubberResult
	res.Element = ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
		times := func(d time.Duration) {
			s := mmss(d)
			ui.Box(c).Size(k.D(pixel.MicroWidth(s)), k.D(5)).Role(ui.RoleText).Label(s).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Micro(0, 0, s, pixel.Faint, 1)
			})
		}
		times(*at)
		if input.Slider(c, text, &secs, input.SliderOptions{Max: total.Seconds(), Width: o.Width, Disabled: o.Disabled}).Changed {
			// Never past the end, whatever the float says.
			*at, res.Changed = min(total, time.Duration(secs*float64(time.Second))), true
			slog.Debug("8bit: scrubber", "label", text, "at", *at)
		}
		times(total)
	})
	return res
}
