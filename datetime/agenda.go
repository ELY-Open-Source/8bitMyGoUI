package datetime

import (
	"fmt"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type AgendaOptions struct{}

// Agenda groups events by day, earliest first.
func Agenda(c *ui.Context, events []Event, o AgendaOptions) ui.Element {
	k := pixel.K(c)
	sorted := slices.Clone(events)
	slices.SortStableFunc(sorted, func(a, b Event) int { return a.Start.Compare(b.Start) })
	return ui.Column(c).Gap(k.D(2)).Role(ui.RoleList).Children(func() {
		var last time.Time
		for i, e := range sorted {
			if i == 0 || !sameDay(e.Start, last) {
				head := fmt.Sprintf("%s %d %s", weekdays(k)[e.Start.Weekday()], e.Start.Day(), monthName(k, e.Start.Month()))
				ui.Box(c).Padding(k.D(4), 0, 0, 0).Children(func() {
					pixel.Text(c, head, pixel.TextOptions{Slot: pixel.Faint})
				})
			}
			last = e.Start
			EventChip(c, e, EventChipOptions{})
		}
	})
}
