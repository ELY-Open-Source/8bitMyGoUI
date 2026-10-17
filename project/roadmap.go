package project

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Bar is work on a lane from Start to End.
type Bar struct {
	Title      string
	Start, End time.Time
	// Slot colors the bar's cap; zero picks one by order.
	Slot pixel.Slot
}

// Lane is a row of bars: a team or theme.
type Lane struct {
	Name string
	Bars []Bar
}

type RoadmapOptions struct {
	// Months shown; zero means 6.
	Months int
	// MonthWidth is in dots; zero means 44.
	MonthWidth int
}

var capSlots = []pixel.Slot{pixel.Info, pixel.Good, pixel.Warn, pixel.Magic, pixel.Bad}

const laneH = pixel.LineHeight + 8

// Roadmap lays bars on lanes under months, from from's.
func Roadmap(c *ui.Context, text string, from time.Time, lanes []Lane, o RoadmapOptions) ui.Element {
	k := pixel.K(c)
	months, mw := pixel.Or("Roadmap Months", o.Months, 6), pixel.Or("Roadmap MonthWidth", o.MonthWidth, 44)
	if len(lanes) == 0 || months > 120 {
		panic(fmt.Sprintf("project: Roadmap of %d lanes over %d months", len(lanes), months))
	}
	nameW, said := 0, make([]string, len(lanes))
	for i, l := range lanes {
		var bs []string
		for _, b := range l.Bars {
			if b.Title == "" || !b.Start.Before(b.End) {
				panic(fmt.Sprintf("project: Roadmap bar %q runs %v to %v", b.Title, b.Start, b.End))
			}
			bs = append(bs, b.Title+" "+b.Start.Format("01-02")+"–"+b.End.Format("01-02"))
		}
		if l.Name == "" {
			panic("project: Roadmap has a lane with no name")
		}
		nameW, said[i] = max(nameW, pixel.Measure(l.Name)), l.Name+": "+strings.Join(bs, ", ")
	}
	nameW += 8
	start := time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, from.Location())
	end := start.AddDate(0, months, 0)
	// x places a time in its month's column, within view.
	x := func(t time.Time) int {
		if !start.Before(t) {
			return nameW
		}
		if !t.Before(end) {
			return nameW + months*mw
		}
		m := (t.Year()-start.Year())*12 + int(t.Month()-start.Month())
		from := start.AddDate(0, m, 0)
		f := float64(t.Sub(from)) / float64(from.AddDate(0, 1, 0).Sub(from))
		return nameW + m*mw + int(f*float64(mw))
	}
	top := pixel.LineHeight + 4
	w, h, now := nameW+months*mw, top+len(lanes)*laneH, k.Now()
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(strings.Join(said, "; ")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for m := range months {
			mx := nameW + m*mw
			g.Text(mx+3, 0, month(k, start.AddDate(0, m, 0).Month()), pixel.Faint, 1)
			for y := top; y < h; y += 2 {
				g.Dot(mx, y, pixel.Shade)
			}
		}
		g.HLine(nameW, top-2, months*mw, pixel.Shade)
		for i, l := range lanes {
			y := top + i*laneH
			g.Text(0, y+4, l.Name, pixel.Ink, 1)
			for j, b := range l.Bars {
				// Bars outside the months shrink to nothing.
				x0, x1 := x(b.Start), x(b.End)
				if x1-x0 < 4 {
					continue
				}
				slot := b.Slot
				if slot == 0 {
					slot = capSlots[j%len(capSlots)]
				}
				g.Frame(x0, y+2, x1-x0, laneH-4, pixel.Raised, pixel.Face)
				g.Rect(x0+2, y+4, 3, laneH-8, slot)
				g.Clip(x0+7, y+2, x1-x0-9, laneH-4, func() { g.Text(x0+7, y+4, b.Title, pixel.Ink, 1) })
			}
		}
		// Today is a dashed Accent line, while in view.
		if !now.Before(start) && now.Before(end) {
			for tx, y := x(now), top-2; y < h; y += 3 {
				g.Rect(tx, y, 1, 2, pixel.Accent)
			}
		}
	})
}
