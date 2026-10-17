package datetime

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"
)

type DateRangeOptions struct {
	Disabled bool
}

type DateRangeResult struct {
	Element ui.Element
	Changed bool
}

type rangeKey struct{}

// ranging is the keys' day; tail picks the end next.
type ranging struct {
	cur  time.Time
	tail bool
}

// DateRange picks a start, then an end; days between dither.
func DateRange(c *ui.Context, text string, start, end *time.Time, o DateRangeOptions) DateRangeResult {
	home := ui.Box(c)
	st := ui.Local(home, rangeKey{}, func() ranging { return ranging{cur: *start} })
	e, _, picked := month(c, text, &st.cur, o.Disabled, false, func(d time.Time) look {
		day := time.Date(d.Year(), d.Month(), d.Day(), 12, 0, 0, 0, time.UTC)
		from := time.Date(start.Year(), start.Month(), start.Day(), 12, 0, 0, 0, time.UTC)
		to := time.Date(end.Year(), end.Month(), end.Day(), 12, 0, 0, 0, time.UTC)
		return look{chosen: sameDay(d, *start) || sameDay(d, *end), between: day.Compare(from) > 0 && day.Before(to)}
	})
	res := DateRangeResult{Element: e}
	if picked {
		switch d := st.cur; {
		case !st.tail:
			*start, *end = onDay(*start, d), onDay(*end, d)
		case d.Before(*start):
			*start, *end = onDay(*start, d), onDay(*end, *start)
		default:
			*end = onDay(*end, d)
		}
		st.tail, res.Changed = !st.tail, true
		slog.Debug("8bit: date range", "label", text, "start", start.Format("2006-01-02"), "end", end.Format("2006-01-02"))
	}
	return res
}
