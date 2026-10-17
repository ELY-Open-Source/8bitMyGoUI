package code

import (
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CommandHistoryOptions struct {
	// Width and Height are dots; zeros mean 200 by 120.
	Width, Height int
	Disabled      bool
}

type CommandHistoryResult struct {
	Element ui.Element
	// Rerun is Enter or a double click on a command.
	Changed, Rerun bool
}

// CommandHistory numbers past commands; Enter reruns one.
func CommandHistory(c *ui.Context, text string, commands []string, selected *int, o CommandHistoryOptions) CommandHistoryResult {
	k := pixel.K(c)
	var r data.ListResult
	e := ui.Box(c).Width(k.D(pixel.Or("CommandHistory Width", o.Width, 200))).Children(func() {
		r = data.List(c, text, len(commands), selected, data.ListOptions{Height: o.Height, Disabled: o.Disabled}, func(i int) string { return commands[i] }, func(i int) {
			ui.Row(c).Gap(k.D(4)).Children(func() {
				pixel.Text(c, strconv.Itoa(i+1), pixel.TextOptions{Slot: pixel.Faint})
				pixel.Text(c, commands[i], pixel.TextOptions{MaxLines: 1})
			})
		})
	})
	if r.Submitted {
		slog.Debug("8bit: command history rerun", "label", text, "index", *selected)
	}
	return CommandHistoryResult{Element: e, Changed: r.Changed, Rerun: r.Submitted}
}
