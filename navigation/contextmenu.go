package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type ContextMenuOptions struct {
	Disabled bool
}

type ContextMenuResult struct {
	Element ui.Element
	// Chosen is the item taken this frame, or -1.
	Chosen int
}

type contextKey struct{}

// spot is the menu and where it opened, in DIPs.
type spot struct {
	m    menu
	x, y float32
}

// ContextMenu opens on right click or the menu key.
func ContextMenu(c *ui.Context, text string, items []MenuItem, o ContextMenuOptions, content func()) ContextMenuResult {
	k := pixel.K(c)
	off := o.Disabled || pixel.Current(c).Inert
	area := ui.Box(c).FocusRing(false).Role(ui.RoleGroup).Label(text)
	// The area takes focus, so the menu key reaches it.
	if off {
		area.Description(k.Msg("disabled"))
	} else {
		area.Focusable()
	}
	s := ui.Local(area, contextKey{}, func() spot { return spot{m: menu{lit: -1}} })
	pixel.Park(c, area, off)
	if off {
		s.m.open = false
	}
	res := ContextMenuResult{Element: area, Chosen: -1}
	if area.RightClicked() && !off {
		s.x, s.y, _ = area.PointerPosition()
		s.m.show(items)
	}
	key, shiftF10 := pixel.Shortcut(area, 0, ui.KeyContextMenu), pixel.Shortcut(area, ui.Shift, ui.KeyF10)
	if (key || shiftF10) && !off {
		// At the pointer when over the area, else its corner.
		x, y, over := area.PointerPosition()
		if !over {
			x, y = 0, 0
		}
		s.x, s.y = x, y
		s.m.show(items)
	}
	// Right clicks keep the focus, so the window takes keys.
	res.Chosen = s.m.keys(c.Shortcut, items)
	focus := area.FocusVisible()
	area.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focus {
			g := k.Canvas(p, r)
			g.FocusMark(0, 0, g.W, g.H)
		}
	}).Children(content)
	panel, clicked := s.m.list(c, area, items)
	if clicked >= 0 {
		res.Chosen = clicked
	}
	if panel.Valid() {
		b := area.Bounds()
		panel.Left(b.X + s.x).Top(b.Y + s.y)
	}
	if res.Chosen >= 0 {
		slog.Debug("8bit: context menu", "item", items[res.Chosen].Label)
	}
	return res
}
