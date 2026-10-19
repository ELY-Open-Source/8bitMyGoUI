package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type Menu struct {
	Title string
	Items []MenuItem
}

type MenuBarOptions struct {
	Disabled bool
}

type MenuBarResult struct {
	Element ui.Element
	// Menu and Item are this frame's choice, else -1.
	Menu, Item int
}

type barKey struct{}

// bar is the open menu, its state and the pointer.
type bar struct {
	open int
	m    menu
	ptr  [2]float32
}

// MenuBar is a bar of menus; arrows go across, down.
func MenuBar(c *ui.Context, text string, menus []Menu, o MenuBarOptions) MenuBarResult {
	k := pixel.K(c)
	off := o.Disabled || pixel.Current(c).Inert
	row := ui.Row(c).Gap(k.D(2)).Padding(k.D(2), k.D(4)).AlignItems(ui.Center).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).Label(text)
	b := ui.Local(row, barKey{}, func() bar { return bar{open: -1} })
	res := MenuBarResult{Element: row, Menu: -1, Item: -1}
	open := func(i int) {
		b.open = i
		b.m.show(menus[i].Items)
	}
	// The title keeps the focus and parks it while off.
	if off && b.open >= 0 {
		b.open = -1
	}
	if b.open >= 0 && !b.m.open {
		b.open = -1
	}
	n := len(menus)
	// Left and Right hand an open menu to its neighbor.
	if b.open >= 0 {
		left, right := pixel.Shortcut(row, 0, ui.KeyLeft), pixel.Shortcut(row, 0, ui.KeyRight)
		switch {
		case left:
			open((b.open + n - 1) % n)
		case right:
			open((b.open + 1) % n)
		}
	}
	// Only a moving pointer hands the bar to another title.
	px, py, over := row.PointerPosition()
	moved := over && [2]float32{px, py} != b.ptr
	b.ptr = [2]float32{px, py}
	var anchor ui.Element
	row.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Face)
		g.HLine(0, g.H-1, g.W, pixel.Edge)
	}).Children(func() {
		for i, mn := range menus {
			t := input.Button(c, mn.Title, input.ButtonOptions{Variant: input.Ghost, Size: input.Small, Disabled: o.Disabled})
			if _, _, on := t.Element.PointerPosition(); on && moved && b.open >= 0 && i != b.open {
				open(i)
			}
			if i == b.open {
				anchor = t.Element
			}
			switch {
			case t.Clicked && b.open == i:
				b.open = -1
			case t.Clicked:
				open(i)
				anchor = t.Element
			case b.open < 0 && !off && t.Element.FocusWithin() && pixel.Shortcut(row, 0, ui.KeyDown):
				open(i)
				anchor = t.Element
			}
		}
	})
	if b.open >= 0 && anchor.Valid() {
		items := menus[b.open].Items
		res.Item = b.m.keys(row.Shortcut, items)
		if _, clicked := b.m.list(c, anchor, items); clicked >= 0 {
			res.Item = clicked
		}
		if res.Item >= 0 {
			res.Menu = b.open
			b.open = -1
		}
	}
	if res.Item >= 0 {
		slog.Debug("8bit: menu bar", "menu", menus[res.Menu].Title, "item", menus[res.Menu].Items[res.Item].Label)
	}
	return res
}
