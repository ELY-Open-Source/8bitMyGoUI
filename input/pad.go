package input

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// pad is a focusable grid returning the picked cell.
func pad(c *ui.Context, text string, cols, n, cw, ch int, cur *int, disabled bool) (shell, base *ui.Element, st pixel.State, picked int) {
	k := pixel.K(c)
	shell, base, st = pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Disabled(disabled)
	})
	shell.Size(k.D(cols*cw), k.D((n+cols-1)/cols*ch))
	for _, m := range []struct {
		key ui.Key
		to  func() int
	}{
		{ui.KeyLeft, func() int { return *cur - 1 }}, {ui.KeyRight, func() int { return *cur + 1 }},
		{ui.KeyUp, func() int { return *cur - cols }}, {ui.KeyDown, func() int { return *cur + cols }},
		{ui.KeyHome, func() int { return 0 }}, {ui.KeyEnd, func() int { return n - 1 }},
	} {
		if to := m.to(); base.Shortcut(0, m.key) && to >= 0 && to < n {
			*cur = to
		}
	}
	picked = -1
	enter, space := base.Shortcut(0, ui.KeyEnter), base.Shortcut(0, ui.KeySpace)
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

// cursor lights a cell with an Accent ring.
func cursor(g pixel.Canvas, x, y, w, h int) {
	g.Rect(x, y, w, h, pixel.Accent)
	g.Rect(x+1, y+1, w-2, h-2, pixel.Lit)
}
