package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CommandMenuOptions struct {
	Disabled bool
}

type CommandMenuResult struct {
	Element *ui.Element
	// Chosen is the item taken this frame, or -1.
	Chosen   int
	Canceled bool
}

type cursorKey struct{}

// CommandMenu is a JRPG command window with a bobbing ▶.
func CommandMenu(c *ui.Context, text string, items []MenuItem, o CommandMenuOptions) CommandMenuResult {
	k := pixel.K(c)
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleList).Label(text).Disabled(o.Disabled)
	})
	cur := ui.Local(base, cursorKey{}, func() int { return 0 })
	res := CommandMenuResult{Element: shell, Chosen: -1}
	n := len(items)
	pick := func(i int) {
		if !items[i].Disabled {
			res.Chosen = i
		}
	}
	if base.Shortcut(0, ui.KeyDown) {
		*cur = (*cur + 1) % n
	}
	if base.Shortcut(0, ui.KeyUp) {
		*cur = (*cur + n - 1) % n
	}
	enter, space := base.Shortcut(0, ui.KeyEnter), base.Shortcut(0, ui.KeySpace)
	if enter || space {
		pick(*cur)
	}
	res.Canceled = base.Shortcut(0, ui.KeyEscape)
	// Enter and Space are taken, so clicks come from pointers.
	if base.Clicked() {
		_, y, _ := base.PointerPosition()
		if i := int(y/k.Dot()-4) / pixel.LineHeight; i >= 0 && i < n {
			*cur = i
			pick(i)
		}
	}
	base.Description(items[*cur].Label)
	w := 0
	for _, it := range items {
		w = max(w, pixel.Measure(it.Label))
	}
	shell.Size(k.D(w+28), k.D(n*pixel.LineHeight+8)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
		for i, it := range items {
			y, ink := 4+i*pixel.LineHeight, pixel.Ink
			if it.Disabled {
				ink = pixel.Faint
			}
			if i == *cur {
				g.Sprite(6+g.Step(4)%2, y+5, pointer, 0, pixel.Accent, 1)
			}
			g.Text(16, y, it.Label, ink, 1)
		}
		g.Finish(st)
	})
	switch {
	case res.Chosen >= 0:
		slog.Debug("8bit: command", "label", text, "item", items[res.Chosen].Label)
	case res.Canceled:
		slog.Debug("8bit: command canceled", "label", text)
	}
	return res
}
