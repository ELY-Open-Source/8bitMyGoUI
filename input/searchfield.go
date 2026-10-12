package input

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SearchFieldOptions struct {
	Placeholder string
	Width       int
	Disabled    bool
}

type SearchFieldResult = TextInputResult

// SearchField has a magnifier; Escape clears, Enter submits.
func SearchField(c *ui.Context, label string, value *string, o SearchFieldOptions) SearchFieldResult {
	ed, e := field(c, value, pixel.EditorOptions{Label: label, Placeholder: o.Placeholder, Disabled: o.Disabled}, o.Width, nil, icon.Search, nil)
	if pixel.Shortcut(ed.Input, 0, ui.KeyEscape) && *value != "" {
		*value, ed.Changed = "", true
	}
	return result(e, ed, label)
}
