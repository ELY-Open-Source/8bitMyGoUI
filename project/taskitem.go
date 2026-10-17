package project

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Task is a thing to do.
type Task struct {
	Title    string
	Done     bool
	Priority Level
	// Assignee is a party member's name; empty is nobody.
	Assignee string
}

type TaskItemOptions struct {
	Disabled bool
}

type TaskItemResult struct {
	Element ui.Element
	Changed bool
}

// TaskItem is a checkbox, title, priority and assignee.
func TaskItem(c *ui.Context, t *Task, o TaskItemOptions) TaskItemResult {
	k := pixel.K(c)
	if t.Title == "" {
		panic("project: TaskItem has no title")
	}
	var res TaskItemResult
	res.Element = ui.Row(c).FillWidth().Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		res.Changed = input.Checkbox(c, t.Title, &t.Done, input.CheckboxOptions{Disabled: o.Disabled}).Changed
		ui.Box(c).Grow(1)
		Priority(c, t.Priority, PriorityOptions{})
		// Nobody shows as an empty dotted seat.
		if t.Assignee == "" {
			ui.Box(c).Size(k.D(16), k.D(16)).Role(ui.RoleImage).Label(k.Msg("nobody")).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Veil(0, 0, 16, 16, pixel.Shade, 8)
			})
			return
		}
		display.Avatar(c, t.Assignee, display.AvatarOptions{})
	})
	if res.Changed {
		slog.Debug("8bit: task", "done", t.Done)
	}
	return res
}
