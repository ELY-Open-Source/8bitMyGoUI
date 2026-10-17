package desk

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LockScreenOptions struct{}

type LockScreenResult struct {
	Element ui.Element
	// Unlocked is a key or click opening it this frame.
	Unlocked bool
}

// LockScreen covers the window: time, date, PRESS ANY KEY.
func LockScreen(c *ui.Context, text string, open *bool, o LockScreenOptions) LockScreenResult {
	k := pixel.K(c)
	now := k.Now()
	if *open {
		k.At(c, now.Truncate(time.Minute).Add(time.Minute))
	}
	clock := now.Format("15:04")
	days, months := strings.Fields(k.Msg("days")), strings.Fields(k.Msg("months"))
	date := fmt.Sprintf(k.Msg("day_date"), days[now.Weekday()], now.Day(), months[now.Month()-1], int(now.Month()))
	press := k.Msg("press_any_key")
	anchor, panel, woke := cover(c, k, text, open, false, func(g pixel.Canvas) {
		g.Rect(0, 0, g.W, g.H, pixel.Well)
		cy := g.H/2 - 40
		g.Text((g.W-4*pixel.Measure(clock))/2, cy, clock, pixel.Ink, 4)
		g.Text((g.W-pixel.Measure(date))/2, cy+4*pixel.LineHeight+4, date, pixel.Faint, 1)
		// The prompt blinks on 4 fps steps; reduced, it holds.
		if g.Step(4)%4 < 2 {
			g.Text((g.W-pixel.Measure(press))/2, g.H-3*pixel.LineHeight, press, pixel.Accent, 1)
		}
	})
	if panel.Valid() {
		panel.Description(clock + ", " + date + ", " + press)
	}
	if woke {
		slog.Debug("8bit: lock screen unlocked", "label", text)
	}
	return LockScreenResult{Element: anchor, Unlocked: woke}
}
