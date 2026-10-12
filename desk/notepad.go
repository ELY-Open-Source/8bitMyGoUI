package desk

import (
	"log/slog"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/navigation"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type NotepadOptions struct {
	// Width is in dots, zero 240; Lines zero means 6.
	Width, Lines int
	Disabled     bool
}

type NotepadResult struct {
	Element *ui.Element
	// Saved is FILE, SAVE chosen this frame.
	Changed, Saved bool
}

// Notepad is a text area under FILE and EDIT menus.
func Notepad(c *ui.Context, text string, value *string, o NotepadOptions) NotepadResult {
	k := pixel.K(c)
	w, lines := pixel.Or("Notepad Width", o.Width, 240), pixel.Or("Notepad Lines", o.Lines, 6)
	var res NotepadResult
	res.Element = ui.Column(c).Width(k.D(w)).Gap(k.D(2)).Role(ui.RoleGroup).Label(text).Children(func() {
		m := navigation.MenuBar(c, k.Msg("menu_bar"), []navigation.Menu{
			{Title: k.Msg("file_menu"), Items: []navigation.MenuItem{{Label: k.Msg("new_file")}, {Label: k.Msg("save")}}},
			{Title: k.Msg("edit_menu"), Items: []navigation.MenuItem{{Label: k.Msg("time_date")}}},
		}, navigation.MenuBarOptions{Disabled: o.Disabled})
		switch [2]int{m.Menu, m.Item} {
		case [2]int{0, 0}:
			*value, res.Changed = "", true
		case [2]int{0, 1}:
			res.Saved = true
		case [2]int{1, 0}:
			// As old Notepad's F5: the time and date, typed in.
			*value, res.Changed = *value+k.Now().Format("15:04 2006-01-02"), true
		}
		f := input.TextArea(c, text, value, input.TextAreaOptions{Width: w, Lines: lines, Disabled: o.Disabled})
		res.Changed = res.Changed || f.Changed
	})
	if res.Changed || res.Saved {
		slog.Debug("8bit: notepad", "label", text, "runes", utf8.RuneCountInString(*value), "saved", res.Saved)
	}
	return res
}
