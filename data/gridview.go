package data

import (
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type GridViewOptions struct {
	// CellWidth and CellHeight are each item's dots; zeros mean 48.
	CellWidth, CellHeight int
	// Height is the view's in dots; zero means 140.
	Height   int
	Disabled bool
}

type GridViewResult struct {
	Element *ui.Element
	// Submitted is Enter or a double click on an item.
	Changed, Submitted bool
}

type gridKey struct{}

// grid holds the rows' place and clicks met during layout.
type grid struct {
	rows            ui.ListState
	clicked, double bool
}

// GridView is a virtual grid framing its chosen item.
func GridView(c *ui.Context, text string, n int, selected *int, o GridViewOptions, item func(i int)) GridViewResult {
	k := pixel.K(c)
	off := o.Disabled || pixel.Current(c).Inert
	cw, ch := pixel.Or("GridView CellWidth", o.CellWidth, 48), pixel.Or("GridView CellHeight", o.CellHeight, 48)
	box := ui.Box(c).FillWidth().Role(ui.RoleGroup).Label(text)
	if !off {
		box.Focusable()
	}
	pixel.Park(c, box, off)
	s := ui.Local(box, gridKey{}, func() grid { return grid{} })
	res := GridViewResult{Element: box, Changed: s.clicked, Submitted: s.double}
	s.clicked, s.double = false, false
	cols := 1
	box.Children(func() {
		holder := ui.Box(c).FillWidth()
		w, check := pixel.WrapWidth(c, holder)
		cols = max(1, w/cw)
		holder.Draw(func(p *ui.Painter, r ui.Rect) { check(r) }).Children(func() {
			ui.List(c, &s.rows, (n+cols-1)/cols, pixel.Lazy(c, func(r int) {
				ui.Row(c).Children(func() {
					for i := r * cols; i < min(n, (r+1)*cols); i++ {
						cell := ui.Box(c).Size(k.D(cw), k.D(ch))
						if !off && cell.Clicked() {
							*selected, s.clicked = i, true
							box.Focus()
							c.Invalidate()
						}
						if !off && cell.DoubleClicked() {
							s.double = true
						}
						on := i == *selected
						cell.Draw(func(p *ui.Painter, r ui.Rect) {
							if on {
								g := k.Canvas(p, r)
								g.Rect(0, 0, g.W, 2, pixel.Accent)
								g.Rect(0, g.H-2, g.W, 2, pixel.Accent)
								g.Rect(0, 0, 2, g.H, pixel.Accent)
								g.Rect(g.W-2, 0, 2, g.H, pixel.Accent)
							}
						}).Children(func() { item(i) })
					}
				})
			})).Width(k.D(cols * cw)).Height(k.D(pixel.Or("GridView Height", o.Height, 140)))
		})
		veil(c, off)
	})
	old := *selected
	for _, m := range []struct {
		key ui.Key
		to  int
	}{
		{ui.KeyLeft, old - 1}, {ui.KeyRight, old + 1}, {ui.KeyUp, old - cols}, {ui.KeyDown, old + cols},
		{ui.KeyHome, 0}, {ui.KeyEnd, n - 1},
	} {
		// The grid stays enabled, so off is checked here.
		if !off && pixel.Shortcut(box, 0, m.key) && m.to >= 0 && m.to < n {
			*selected = m.to
		}
	}
	if *selected != old {
		res.Changed = true
		s.rows.ScrollIntoView(*selected / cols)
	}
	if !off && pixel.Shortcut(box, 0, ui.KeyEnter) {
		res.Submitted = true
	}
	if off {
		box.Description(k.Msg("disabled"))
	} else if *selected >= 0 {
		box.Description(strconv.Itoa(*selected+1) + "/" + strconv.Itoa(n))
	}
	well(c, box, pixel.State{Focus: box.FocusRing(false).FocusVisible(), Disabled: o.Disabled})
	if res.Changed || res.Submitted {
		slog.Debug("8bit: grid view", "label", text, "selected", *selected, "submitted", res.Submitted)
	}
	return res
}
