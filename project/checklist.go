package project

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/feedback"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ChecklistOptions struct {
	// Width is in dots; zero means 200.
	Width    int
	Disabled bool
}

type ChecklistResult struct {
	Element *ui.Element
	Changed bool
}

// Checklist ticks off subtasks under a progress bar.
func Checklist(c *ui.Context, text string, items []string, done []bool, o ChecklistOptions) ChecklistResult {
	k := pixel.K(c)
	if len(items) == 0 || len(done) != len(items) {
		panic(fmt.Sprintf("project: Checklist has %d items and %d marks", len(items), len(done)))
	}
	w := pixel.Or("Checklist Width", o.Width, 200)
	// The bar leaves room for the widest count, "9/9".
	bar := w - 4 - pixel.Measure(fmt.Sprintf("%d/%d", len(items), len(items)))
	if bar < 8 {
		panic(fmt.Sprintf("project: Checklist Width %d leaves a %d-dot bar", w, bar))
	}
	n := 0
	for _, d := range done {
		if d {
			n++
		}
	}
	var res ChecklistResult
	res.Element = ui.Column(c).Width(k.D(w)).Gap(k.D(3)).Role(ui.RoleGroup).Label(text).Children(func() {
		ui.Row(c).FillWidth().Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
			feedback.Progress(c, text, float64(n)/float64(len(items)), feedback.ProgressOptions{Width: bar})
			pixel.Text(c, fmt.Sprintf("%d/%d", n, len(items)), pixel.TextOptions{Slot: pixel.Faint})
		})
		for i, it := range items {
			res.Changed = input.Checkbox(c, it, &done[i], input.CheckboxOptions{Disabled: o.Disabled}).Changed || res.Changed
		}
		// A full list says so, as a quest log would.
		if n == len(items) {
			pixel.Text(c, k.Msg("quest_complete"), pixel.TextOptions{Slot: pixel.Good})
		}
	})
	if res.Changed {
		slog.Debug("8bit: checklist", "label", text)
	}
	return res
}
