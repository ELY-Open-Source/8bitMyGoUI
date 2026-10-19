// Package datetime picks, shows and counts dates and times.
package datetime

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Event is a timed entry of a day.
type Event struct {
	Start  time.Time
	Length time.Duration
	Title  string
}

// monthName is m's short name in the frame's locale.
func monthName(k pixel.Kit, m time.Month) string {
	return strings.Fields(k.Msg("months"))[m-1]
}

// weekdays are the day names, Sunday first.
func weekdays(k pixel.Kit) []string { return strings.Fields(k.Msg("weekdays")) }

// clock prints d as MM:SS, or H:MM:SS past an hour.
func clock(d time.Duration) string {
	s := int(d / time.Second)
	if s >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", s/3600, s/60%60, s%60)
	}
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}

// sameDay reports whether a and b share a date.
func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// onDay moves t to d's date, keeping clock and place.
func onDay(t, d time.Time) time.Time {
	y, m, day := d.Date()
	h, mi, s := t.Clock()
	return time.Date(y, m, day, h, mi, s, t.Nanosecond(), t.Location())
}

// today is now; the view rebuilds as the day turns.
func today(c *ui.Context, k pixel.Kit) time.Time {
	n := k.Now()
	k.At(c, time.Date(n.Year(), n.Month(), n.Day()+1, 0, 0, 0, 0, n.Location()))
	return n
}
