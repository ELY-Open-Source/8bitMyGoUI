package data

import (
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SortableOptions struct {
	// Width is in dots; zero means 160.
	Width    int
	Disabled bool
}

type SortableResult struct {
	Element ui.Element
	Changed bool
}

type sortKey struct{}

// sorting is the cursor row; held while pressed.
type sorting struct {
	cur  int
	held bool
}

// Sortable reorders rows by dragging, or Alt with the arrows.
func Sortable(c *ui.Context, text string, items *[]string, o SortableOptions) SortableResult {
	k := pixel.K(c)
	n := len(*items)
	off := o.Disabled || pixel.Current(c).Inert
	shell, base, st := pixel.Shell(c, func() ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleList).Label(text).Disabled(o.Disabled)
	})
	s := ui.Local(base, sortKey{}, func() sorting { return sorting{} })
	res := SortableResult{Element: shell}
	if n == 0 {
		return res
	}
	s.cur = max(0, min(s.cur, n-1))
	move := func(from, to int) {
		if to < 0 || to >= n || to == from {
			return
		}
		it := (*items)[from]
		*items = slices.Insert(slices.Delete(*items, from, from+1), to, it)
		s.cur, res.Changed = to, true
	}
	if pixel.Shortcut(base, 0, ui.KeyDown) {
		s.cur = min(s.cur+1, n-1)
	}
	if pixel.Shortcut(base, 0, ui.KeyUp) {
		s.cur = max(s.cur-1, 0)
	}
	if pixel.Shortcut(base, ui.Alt, ui.KeyDown) {
		move(s.cur, s.cur+1)
	}
	if pixel.Shortcut(base, ui.Alt, ui.KeyUp) {
		move(s.cur, s.cur-1)
	}
	// A press takes a row; dragging carries it.
	_, _, held := base.Dragged()
	if held = held && !off; held {
		_, y, _ := base.PointerPosition()
		at := max(0, min(n-1, int(y/k.Dot())/pixel.LineHeight))
		if !s.held {
			s.cur = at
		}
		move(s.cur, at)
	}
	s.held = held
	base.Description((*items)[s.cur])
	shown, cur, held := append([]string(nil), *items...), s.cur, s.held
	shell.Size(k.D(pixel.Or("Sortable Width", o.Width, 160)), k.D(n*pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, it := range shown {
			y, ink := i*pixel.LineHeight, pixel.Ink
			switch {
			case i == cur && held:
				g.Rect(0, y, g.W, pixel.LineHeight, pixel.Lit)
				g.HLine(0, y+pixel.LineHeight-1, g.W, pixel.Shade)
			case i == cur:
				g.Rect(0, y, g.W, pixel.LineHeight, pixel.Accent)
				ink = pixel.AccentInk
			}
			for d := range 3 {
				g.Rect(4, y+4+d*3, 2, 2, ink)
				g.Rect(7, y+4+d*3, 2, 2, ink)
			}
			g.Text(14, y, it, ink, 1)
		}
		g.Finish(st)
	})
	if res.Changed {
		slog.Debug("8bit: sorted", "label", text, "items", *items)
	}
	return res
}
