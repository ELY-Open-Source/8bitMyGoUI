package project

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/feedback"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Sprint is a stretch of work, counted in points.
type Sprint struct {
	Name        string
	Start, End  time.Time
	Done, Total int
}

type SprintPanelOptions struct {
	// Width is in dots; zero means 200.
	Width int
}

const day = 24 * time.Hour

// SprintPanel shows a sprint's points done and days left.
func SprintPanel(c *ui.Context, s Sprint, o SprintPanelOptions) ui.Element {
	k := pixel.K(c)
	if s.Name == "" || !s.Start.Before(s.End) || s.Total < 1 || s.Done < 0 || s.Done > s.Total {
		panic(fmt.Sprintf("project: SprintPanel %q from %v to %v, %d of %d", s.Name, s.Start, s.End, s.Done, s.Total))
	}
	w, now := pixel.Or("SprintPanel Width", o.Width, 200), k.Now()
	// Room for the percent; a zero bar would mean 120.
	bar := w - 4 - pixel.MicroWidth("100%")
	if bar < 8 {
		panic(fmt.Sprintf("project: SprintPanel Width %d leaves a %d-dot bar", w, bar))
	}
	rest := s.End.Sub(now)
	// Rounded up without adding, so far ends cannot overflow.
	left, ink := fmt.Sprintf(k.Msg("days_left"), 1+(rest-1)/day), pixel.Faint
	switch {
	case rest <= 0:
		left, ink = k.Msg("over"), pixel.Bad
	case rest <= day:
		left, ink = k.Msg("last_day"), pixel.Warn
	}
	if rest > 0 {
		// The count turns when a whole day runs out.
		k.At(c, now.Add(rest-(rest-1)/day*day))
	}
	// The panel hugs its column instead of the parent's width.
	return layout.Panel(c, layout.PanelOptions{Style: pixel.Raised}, func() {
		ui.Column(c).Width(k.D(w)).Gap(k.D(3)).Role(ui.RoleGroup).Label(s.Name).Children(func() {
			pixel.Text(c, s.Name, pixel.TextOptions{Wrap: true, MaxLines: 1})
			pixel.Text(c, s.Start.Format("01-02")+" → "+s.End.Format("01-02"), pixel.TextOptions{Slot: pixel.Faint})
			feedback.Progress(c, s.Name, float64(s.Done)/float64(s.Total), feedback.ProgressOptions{Width: bar, Percent: true})
			ui.Row(c).FillWidth().Gap(k.D(4)).Children(func() {
				pixel.Text(c, fmt.Sprintf(k.Msg("sprint_points"), s.Done, s.Total), pixel.TextOptions{}).Grow(1)
				pixel.Text(c, left, pixel.TextOptions{Slot: ink})
			})
		})
	}).AlignSelf(ui.Start)
}
