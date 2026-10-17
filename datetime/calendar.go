package datetime

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"
)

type CalendarOptions struct {
	Disabled bool
}

type CalendarResult struct {
	Element ui.Element
	// Picked is a click, Enter or Space on a day.
	Changed, Picked bool
}

// Calendar moves *date through a month grid; today is marked.
func Calendar(c *ui.Context, text string, date *time.Time, o CalendarOptions) CalendarResult {
	old := *date
	e, _, picked := month(c, text, date, o.Disabled, false, func(d time.Time) look { return look{chosen: sameDay(d, *date)} })
	res := CalendarResult{Element: e, Changed: !date.Equal(old), Picked: picked}
	if res.Changed || res.Picked {
		slog.Debug("8bit: calendar", "label", text, "date", date.Format("2006-01-02"), "picked", picked)
	}
	return res
}
