package datetime

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type WeekStripOptions struct {
	Disabled bool
}

type WeekStripResult struct {
	Element *ui.Element
	Changed bool
}

// day is a week strip cell in dots.
const dayW, dayH = 22, 2 * pixel.LineHeight

// WeekStrip shows *date's week; Left and Right step days.
func WeekStrip(c *ui.Context, text string, date *time.Time, o WeekStripOptions) WeekStripResult {
	k := pixel.K(c)
	old := *date
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Disabled(o.Disabled)
	})
	if base.Shortcut(0, ui.KeyLeft) {
		*date = date.AddDate(0, 0, -1)
	}
	if base.Shortcut(0, ui.KeyRight) {
		*date = date.AddDate(0, 0, 1)
	}
	first := date.AddDate(0, 0, -int(date.Weekday()))
	// Enter and Space are taken, so clicks come from pointers.
	base.Shortcut(0, ui.KeyEnter)
	base.Shortcut(0, ui.KeySpace)
	if base.Clicked() {
		x, _, _ := base.PointerPosition()
		if i := int(x / k.D(dayW)); x >= 0 && i < 7 {
			*date = first.AddDate(0, 0, i)
		}
	}
	base.Description(date.Format("2006-01-02"))
	names, chosen, now := weekdays(k), *date, today(c, k)
	shell.Size(k.D(7*dayW), k.D(dayH)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i := range 7 {
			d, x := first.AddDate(0, 0, i), i*dayW
			if sameDay(d, chosen) {
				g.Frame(x, 0, dayW, dayH, pixel.Sunken, pixel.Well)
			}
			n := strconv.Itoa(d.Day())
			g.Text(x+(dayW-pixel.Measure(names[i]))/2, 0, names[i], pixel.Faint, 1)
			g.Text(x+(dayW-pixel.Measure(n))/2, pixel.LineHeight-2, n, pixel.Ink, 1)
			if sameDay(d, now) {
				g.HLine(x+4, dayH-2, dayW-8, pixel.Accent)
			}
		}
		g.Finish(st)
	})
	res := WeekStripResult{Element: shell, Changed: !date.Equal(old)}
	if res.Changed {
		slog.Debug("8bit: week strip", "label", text, "date", date.Format("2006-01-02"))
	}
	return res
}
