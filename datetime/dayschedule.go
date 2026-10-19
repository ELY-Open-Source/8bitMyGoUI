package datetime

import (
	"fmt"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type DayScheduleOptions struct {
	// From and To bound the hours; zeros mean 8–18.
	From, To int
}

// hourH is an hour row's height; gutter holds labels.
const hourH, gutter = pixel.LineHeight + 4, 24

// DaySchedule places events on hour rows by start and length.
func DaySchedule(c *ui.Context, events []Event, o DayScheduleOptions) ui.Element {
	k := pixel.K(c)
	from, to := o.From, o.To
	if from == 0 && to == 0 {
		from, to = 8, 18
	}
	if from < 0 || to > 24 || from >= to {
		panic(fmt.Sprintf("datetime: DaySchedule hours %d to %d", from, to))
	}
	sorted := slices.Clone(events)
	slices.SortStableFunc(sorted, func(a, b Event) int { return a.Start.Compare(b.Start) })
	// Chains of overlapping events share the width in lanes.
	lanes, shared := make([]int, len(sorted)), make([]int, len(sorted))
	var ends []time.Time
	var last time.Time
	head := 0
	for i, e := range sorted {
		if i > 0 && e.Start.Compare(last) >= 0 {
			for j := head; j < i; j++ {
				shared[j] = len(ends)
			}
			head, ends = i, nil
		}
		lane := slices.IndexFunc(ends, func(end time.Time) bool { return end.Compare(e.Start) <= 0 })
		if lane < 0 {
			lane, ends = len(ends), append(ends, time.Time{})
		}
		lanes[i], ends[lane] = lane, e.Start.Add(e.Length)
		if end := ends[lane]; end.Compare(last) > 0 || i == head {
			last = end
		}
	}
	for j := head; j < len(sorted); j++ {
		shared[j] = len(ends)
	}
	box := ui.Box(c).FillWidth().Height(k.D((to - from) * hourH)).Role(ui.RoleList)
	w, check := pixel.WrapWidth(c, box)
	return box.Draw(func(p *ui.Painter, r ui.Rect) {
		check(r)
		g := k.Canvas(p, r)
		for h := from; h < to; h++ {
			y := (h - from) * hourH
			g.Micro(0, y+1, fmt.Sprintf("%02d:00", h), pixel.Faint, 1)
			g.Checker(gutter, y, g.W-gutter, 1, pixel.Edge)
		}
	}).Children(func() {
		for i, e := range sorted {
			// Minutes from the top row, cut to the hours shown.
			top := e.Start.Hour()*60 + e.Start.Minute() - from*60
			end := min(top+int(e.Length/time.Minute), (to-from)*60)
			if top = max(0, top); end <= top {
				continue
			}
			h := max(3, (end-top)*hourH/60)
			y := min(top*hourH/60, (to-from)*hourH-h)
			laneW := max(4, (w-gutter)/shared[i])
			ui.Box(c).Absolute().Left(k.D(gutter+lanes[i]*laneW)).Top(k.D(y)).Size(k.D(laneW-1), k.D(h)).Role(ui.RoleListItem).
				Label(e.Start.Format("15:04") + " " + e.Title).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Frame(0, 0, g.W, g.H, pixel.Raised, pixel.AccentLit)
				g.Clip(1, 1, g.W-2, g.H-2, func() { g.Text(3, 0, e.Title, pixel.Ink, 1) })
			})
		}
	})
}
