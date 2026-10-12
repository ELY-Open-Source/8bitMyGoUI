package project

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// TimeEntry is time on a task; zero Stop is running.
type TimeEntry struct {
	Task        string
	Start, Stop time.Time
}

type TimeTrackerOptions struct {
	// Width is in dots; zero means 220.
	Width    int
	Disabled bool
}

type TimeTrackerResult struct {
	Element          *ui.Element
	Started, Stopped bool
}

// TimeTracker starts and stops entries and totals them.
func TimeTracker(c *ui.Context, text, task string, entries *[]TimeEntry, o TimeTrackerOptions) TimeTrackerResult {
	k := pixel.K(c)
	es, now := *entries, k.Now()
	if task == "" {
		panic("project: TimeTracker has no task")
	}
	for i, e := range es {
		// Only the last entry may still run.
		if e.Task == "" || e.Start.IsZero() || e.Stop.IsZero() && i < len(es)-1 || !e.Stop.IsZero() && e.Stop.Before(e.Start) {
			panic(fmt.Sprintf("project: TimeTracker entry %d of %q runs %v to %v", i, e.Task, e.Start, e.Stop))
		}
	}
	running := len(es) > 0 && es[len(es)-1].Stop.IsZero()
	var res TimeTrackerResult
	res.Element = ui.Column(c).Width(k.D(pixel.Or("TimeTracker Width", o.Width, 220))).Gap(k.D(4)).Role(ui.RoleGroup).Label(text).Children(func() {
		ui.Row(c).FillWidth().Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
			shown := time.Duration(0)
			if running {
				// A start ahead of now, as clocks drift, reads zero.
				shown = max(0, now.Sub(es[len(es)-1].Start))
				k.At(c, now.Add(time.Second-shown%time.Second))
			}
			pixel.Text(c, hms(shown), pixel.TextOptions{Scale: 2}).Grow(1)
			label, v := k.Msg("start"), input.Primary
			if running {
				label, v = k.Msg("stop"), input.Danger
			}
			if !input.Button(c, label, input.ButtonOptions{Variant: v, Disabled: o.Disabled}).Clicked {
				return
			}
			if running {
				// A stop before a drifted start ends it there.
				last := &(*entries)[len(es)-1]
				last.Stop, res.Stopped = now, true
				if now.Before(last.Start) {
					last.Stop = last.Start
				}
			} else {
				*entries, res.Started = append(es, TimeEntry{Task: task, Start: now}), true
			}
			slog.Debug("8bit: time tracker", "label", text, "running", res.Started, "entries", len(*entries))
		})
		// Newest first; the total counts the running one too.
		total := time.Duration(0)
		for i := len(es) - 1; i >= 0; i-- {
			e := es[i]
			d, ink := e.Stop.Sub(e.Start), pixel.Faint
			if e.Stop.IsZero() {
				d, ink = max(0, now.Sub(e.Start)), pixel.Accent
			}
			// Centuries of entries saturate instead of wrapping.
			if d > math.MaxInt64-total {
				d = math.MaxInt64 - total
			}
			total += d
			// New rows come in on top; each keeps a key.
			ui.Row(c).Key(i).FillWidth().Gap(k.D(4)).Children(func() {
				pixel.Text(c, e.Task, pixel.TextOptions{Wrap: true, MaxLines: 1}).Grow(1)
				pixel.Text(c, hms(d), pixel.TextOptions{Slot: ink})
			})
		}
		ui.Row(c).FillWidth().Gap(k.D(4)).Children(func() {
			pixel.Text(c, k.Msg("total"), pixel.TextOptions{}).Grow(1)
			pixel.Text(c, hms(total), pixel.TextOptions{})
		})
	})
	return res
}
