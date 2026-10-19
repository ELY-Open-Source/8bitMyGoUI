package project

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type TaskListOptions struct {
	// Width is in dots; zero means 240.
	Width    int
	Disabled bool
}

type TaskListResult struct {
	Element ui.Element
	// Changed reports a task checked or a filter turned.
	Changed bool
}

type filterKey struct{}

// TaskList shows tasks under OPEN, DONE and URGENT filters.
func TaskList(c *ui.Context, text string, tasks []Task, o TaskListOptions) TaskListResult {
	k := pixel.K(c)
	col := ui.Column(c).Width(k.D(pixel.Or("TaskList Width", o.Width, 240))).Gap(k.D(4)).Role(ui.RoleGroup).Label(text)
	on := ui.Local(col, filterKey{}, func() []bool { return make([]bool, 3) })
	var res TaskListResult
	res.Element = col.Children(func() {
		if data.FilterBar(c, k.Msg("filters"), strings.Fields(k.Msg("task_filters")), *on, data.FilterBarOptions{Disabled: o.Disabled}).Changed {
			res.Changed = true
			slog.Debug("8bit: task list filters", "label", text, "on", *on)
		}
		shown := 0
		for i := range tasks {
			t := &tasks[i]
			// OPEN and DONE alone narrow; both or neither keep all.
			open, done, urgent := (*on)[0], (*on)[1], (*on)[2]
			if open != done && t.Done != done || urgent && t.Priority < High {
				continue
			}
			shown++
			// Rows come and go with filters; each keeps a key.
			ui.Box(c).Key(i).FillWidth().Children(func() {
				res.Changed = TaskItem(c, t, TaskItemOptions{Disabled: o.Disabled}).Changed || res.Changed
			})
		}
		if shown == 0 {
			pixel.Text(c, k.Msg("empty"), pixel.TextOptions{Slot: pixel.Faint})
		}
	})
	return res
}
