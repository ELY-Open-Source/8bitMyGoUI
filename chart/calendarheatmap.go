package chart

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CalendarHeatmapOptions struct {
	// Slot is the busiest days' color; zero means Good.
	Slot pixel.Slot
}

// day is a 4×4 square plus its gap, in dots.
const day = 5

// CalendarHeatmap rates a year's UTC days 0–4, week by week.
func CalendarHeatmap(c *ui.Context, text string, year int, level func(d time.Time) int, o CalendarHeatmapOptions) ui.Element {
	k := pixel.K(c)
	ink := o.Slot
	if ink == 0 {
		ink = pixel.Good
	}
	first := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	skip := int(first.Weekday())
	levels := make([]int, first.AddDate(1, 0, 0).Sub(first)/(24*time.Hour))
	sums := make([]float64, 12)
	for i := range levels {
		d := first.AddDate(0, 0, i)
		l := level(d)
		if l < 0 || l > 4 {
			panic(fmt.Sprintf("chart: CalendarHeatmap level %d outside 0–4", l))
		}
		levels[i], sums[d.Month()-1] = l, sums[d.Month()-1]+float64(l)
	}
	months := strings.Fields(k.Msg("months"))
	sm := small(months...)
	head := labelH(months)
	weeks := (skip + len(levels) + 6) / 7
	// Assistive technology hears each month's summed levels.
	return ui.Box(c).Size(k.D(weeks*day), k.D(head+7*day)).Role(ui.RoleImage).Label(text + " " + strconv.Itoa(year)).Description(read(months, sums)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		// Month names top their first week, room permitting.
		next := 0
		for m, name := range months {
			x := (skip + time.Date(year, time.Month(m+1), 1, 0, 0, 0, 0, time.UTC).YearDay() - 1) / 7 * day
			if x >= next {
				mark(g, x, 1, name, sm, pixel.Faint)
				next = x + markW(name, sm) + 3
			}
		}
		for i, l := range levels {
			n := skip + i
			g.Dither(n/7*day, head+n%7*day, day-1, day-1, pixel.Shade, ink, l*4)
		}
	})
}
