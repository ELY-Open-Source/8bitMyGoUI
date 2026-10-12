package pixel

import "github.com/egoist/mygo/ui"

// Pad is a focusable grid; it returns the picked cell.
func Pad(c *ui.Context, text string, cols, n, cw, ch int, cur *int, disabled bool) (shell, base *ui.Element, st State, picked int) {
	k := K(c)
	shell, base, st = Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Disabled(disabled)
	})
	shell.Size(k.D(cols*cw), k.D((n+cols-1)/cols*ch))
	picked = -1
	// Disabled or inert, even a key already sent waits.
	if base.IsDisabled() {
		return shell, base, st, picked
	}
	for _, m := range []struct {
		key ui.Key
		to  func() int
	}{
		{ui.KeyLeft, func() int { return *cur - 1 }}, {ui.KeyRight, func() int { return *cur + 1 }},
		{ui.KeyUp, func() int { return *cur - cols }}, {ui.KeyDown, func() int { return *cur + cols }},
		{ui.KeyHome, func() int { return 0 }}, {ui.KeyEnd, func() int { return n - 1 }},
	} {
		if to := m.to(); Shortcut(base, 0, m.key) && to >= 0 && to < n {
			*cur = to
		}
	}
	enter, space := Shortcut(base, 0, ui.KeyEnter), Shortcut(base, 0, ui.KeySpace)
	if enter || space {
		picked = *cur
	}
	// Enter and Space are taken, so clicks come from pointers.
	if base.Clicked() {
		x, y, _ := base.PointerPosition()
		if col, i := int(x/k.D(cw)), int(y/k.D(ch))*cols+int(x/k.D(cw)); x >= 0 && y >= 0 && col < cols && i < n {
			*cur, picked = i, i
		}
	}
	return shell, base, st, picked
}

// Cursor lights a cell with an Accent ring.
func (g Canvas) Cursor(x, y, w, h int) {
	g.Rect(x, y, w, h, Accent)
	g.Rect(x+1, y+1, w-2, h-2, Lit)
}
