package chat

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DateSeparatorOptions struct {
	// Width is in dots; zero means 200.
	Width int
}

// DateSeparator rules a line with TODAY, YESTERDAY or the date.
func DateSeparator(c *ui.Context, day time.Time, o DateSeparatorOptions) ui.Element {
	k := pixel.K(c)
	now := k.Now()
	// At midnight today becomes yesterday.
	k.At(c, time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()))
	return rule(c, k, dayName(k, day.In(now.Location()), now), pixel.Or("DateSeparator Width", o.Width, 200), pixel.Faint, pixel.Edge)
}

// dayName is TODAY, YESTERDAY, or the locale's short date.
func dayName(k pixel.Kit, day, now time.Time) string {
	same := func(a, b time.Time) bool { return a.Year() == b.Year() && a.YearDay() == b.YearDay() }
	switch {
	case same(day, now):
		return k.Msg("today")
	case same(day, now.AddDate(0, 0, -1)):
		return k.Msg("yesterday")
	}
	wd, months := strings.Fields(k.Msg("days")), strings.Fields(k.Msg("months"))
	return fmt.Sprintf(k.Msg("day_date"), wd[day.Weekday()], day.Day(), months[day.Month()-1], int(day.Month()))
}
