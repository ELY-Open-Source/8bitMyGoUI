package datetime

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type StopwatchOptions struct {
	Disabled bool
}

type StopwatchResult struct {
	Element                         ui.Element
	Started, Stopped, Lapped, Reset bool
}

type watchKey struct{}

// watch is time banked, since when it runs, and laps.
type watch struct {
	banked time.Duration
	since  time.Time
	laps   []time.Duration
}

// tenths prints d as MM:SS.T.
func tenths(d time.Duration) string {
	return fmt.Sprintf("%s.%d", clock(d), int(d/(100*time.Millisecond))%10)
}

// Stopwatch times with START, STOP, LAP, RESET; laps below.
func Stopwatch(c *ui.Context, text string, o StopwatchOptions) StopwatchResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(4)).Role(ui.RoleGroup).Label(text)
	w := ui.Local(col, watchKey{}, func() watch { return watch{} })
	now := k.Now()
	running := !w.since.IsZero()
	elapsed := w.banked
	if running {
		elapsed += k.Elapsed(&w.since)
		k.Tick(c, 10)
	}
	res := StopwatchResult{Element: col}
	col.Children(func() {
		s := tenths(elapsed)
		ui.Box(c).Size(k.D(pixel.Measure(s)*2), k.D(pixel.LineHeight*2)).Role(ui.RoleText).Label(s).Draw(func(p *ui.Painter, r ui.Rect) {
			k.Canvas(p, r).Text(0, 0, s, pixel.Ink, 2)
		})
		ui.Row(c).Gap(k.D(4)).Children(func() {
			toggle := k.Msg("start")
			if running {
				toggle = k.Msg("stop")
			}
			if input.Button(c, toggle, input.ButtonOptions{Size: input.Small, Variant: input.Primary, Disabled: o.Disabled}).Clicked {
				if running {
					w.banked, w.since, res.Stopped = elapsed, time.Time{}, true
				} else {
					w.since, res.Started = now, true
				}
			}
			if input.Button(c, k.Msg("lap"), input.ButtonOptions{Size: input.Small, Disabled: o.Disabled || !running}).Clicked {
				w.laps, res.Lapped = append(w.laps, elapsed), true
			}
			if input.Button(c, k.Msg("reset"), input.ButtonOptions{Size: input.Small, Disabled: o.Disabled || running || elapsed == 0}).Clicked {
				*w, res.Reset = watch{}, true
			}
		})
		for i := len(w.laps) - 1; i >= max(0, len(w.laps)-5); i-- {
			pixel.Text(c, fmt.Sprintf("%s %d  %s", k.Msg("lap"), i+1, tenths(w.laps[i])), pixel.TextOptions{Slot: pixel.Faint})
		}
	})
	if res.Started || res.Stopped || res.Lapped || res.Reset {
		slog.Debug("8bit: stopwatch", "label", text, "elapsed", elapsed, "laps", len(w.laps))
	}
	return res
}
