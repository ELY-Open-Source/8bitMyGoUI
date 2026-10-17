package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DropdownMenuOptions struct {
	Disabled bool
}

type DropdownMenuResult struct {
	Element ui.Element
	// Chosen is the item taken this frame, or -1.
	Chosen int
}

// DropdownMenu is a button opening a menu of items.
func DropdownMenu(c *ui.Context, text string, items []MenuItem, o DropdownMenuOptions) DropdownMenuResult {
	off := o.Disabled || pixel.Current(c).Inert
	b := input.Button(c, text+" ▼", input.ButtonOptions{Label: text, Disabled: o.Disabled})
	m := ui.Local(b.Element, menuKey{}, func() menu { return menu{lit: -1} })
	// The button keeps the focus and parks it while off.
	if off {
		m.open = false
	}
	res := DropdownMenuResult{Element: b.Element, Chosen: -1}
	switch {
	case b.Clicked && m.open:
		m.open = false
	case b.Clicked, !m.open && !off && pixel.Shortcut(b.Element, 0, ui.KeyDown):
		m.show(items)
	}
	res.Chosen = m.keys(b.Element.Shortcut, items)
	if _, clicked := m.list(c, b.Element, items); clicked >= 0 {
		res.Chosen = clicked
	}
	if res.Chosen >= 0 {
		slog.Debug("8bit: dropdown", "label", text, "item", items[res.Chosen].Label)
	}
	return res
}
