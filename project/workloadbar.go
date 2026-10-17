package project

import (
	"fmt"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/feedback"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Load is a member's hours of work.
type Load struct {
	Name  string
	Hours float64
}

type WorkloadBarOptions struct {
	// Capacity is each member's hours; zero means 40.
	Capacity float64
	// Width is in dots; zero means 240.
	Width int
}

// WorkloadBar shows each member's load against capacity.
func WorkloadBar(c *ui.Context, text string, loads []Load, o WorkloadBarOptions) ui.Element {
	k := pixel.K(c)
	cp := o.Capacity
	if cp == 0 {
		cp = 40
	}
	if !(cp > 0) || math.IsInf(cp, 0) || len(loads) == 0 {
		panic(fmt.Sprintf("project: WorkloadBar of %d members, capacity %v", len(loads), cp))
	}
	num := func(v float64) string { return strconv.FormatFloat(v, 'g', 4, 64) }
	nameW, valW, vals := 0, 0, make([]string, len(loads))
	for i, l := range loads {
		if l.Name == "" || !(l.Hours >= 0) || math.IsInf(l.Hours, 0) {
			panic(fmt.Sprintf("project: WorkloadBar member %q at %v hours", l.Name, l.Hours))
		}
		vals[i] = fmt.Sprintf(k.Msg("hours_of"), num(l.Hours), num(cp))
		nameW, valW = max(nameW, pixel.Measure(l.Name)), max(valW, pixel.Measure(vals[i]))
	}
	w := pixel.Or("WorkloadBar Width", o.Width, 240)
	// Face, name, meter and hours share a row, three gaps.
	meterW := w - 16 - nameW - valW - 12
	if meterW < 12 {
		panic(fmt.Sprintf("project: WorkloadBar Width %d leaves a %d-dot meter", w, meterW))
	}
	return ui.Column(c).Width(k.D(w)).Gap(k.D(3)).Role(ui.RoleGroup).Label(text).Children(func() {
		for i, l := range loads {
			ui.Row(c).FillWidth().Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
				display.Avatar(c, l.Name, display.AvatarOptions{})
				show(c, k, l.Name, pixel.Ink, nameW)
				// Shares of capacity: past 0.8 warns, past 1 is Bad.
				feedback.Meter(c, l.Name, l.Hours/cp, feedback.MeterOptions{Max: 1.25, Warning: 0.8, Critical: 1, Width: meterW})
				ui.Box(c).Width(k.D(valW)).Children(func() { pixel.Text(c, vals[i], pixel.TextOptions{Slot: pixel.Faint}) })
			})
		}
	})
}
