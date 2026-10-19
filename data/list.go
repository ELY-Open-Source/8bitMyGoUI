package data

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type ListOptions struct {
	// Height is in dots; zero means 120.
	Height   int
	Disabled bool
}

type ListResult struct {
	Element ui.Element
	// Submitted is Enter or a double click on a row.
	Changed, Submitted bool
}

type listKey struct{}

// List is a virtual list; ▶ marks the chosen row.
func List(c *ui.Context, text string, n int, selected *int, o ListOptions, name func(i int) string, row func(i int)) ListResult {
	k := pixel.K(c)
	off := o.Disabled || pixel.Current(c).Inert
	box := ui.Box(c)
	st := ui.Local(box, listKey{}, func() ui.ListState { return ui.ListState{} })
	// Names read aloud and let typing choose a row.
	st.Selected, st.Label = selected, name
	if off {
		st.Selected = nil
	}
	var e ui.Element
	box.Children(func() {
		e = ui.List(c, st, n, pixel.Lazy(c, func(i int) {
			on := i == *selected
			// The chosen row is Face: Ink reads on it always.
			ui.Row(c).FillWidth().Padding(k.D(1), k.D(2), k.D(1), k.D(12)).Draw(func(p *ui.Painter, r ui.Rect) {
				if on {
					g := k.Canvas(p, r)
					g.Rect(0, 0, g.W, g.H, pixel.Face)
					g.Sprite(3, (g.H-pointer.H)/2, pointer, 0, pixel.Accent, 1)
				}
			}).Children(func() { row(i) })
		})).Height(k.D(pixel.Or("List Height", o.Height, 120))).Label(text)
		veil(c, off)
	})
	pixel.Park(c, e, off)
	if off {
		e.Description(k.Msg("disabled"))
	}
	well(c, box, pixel.State{Focus: e.FocusRing(false).FocusVisible(), Disabled: o.Disabled})
	res := ListResult{Element: box, Changed: e.Changed(), Submitted: e.Submitted()}
	if res.Changed {
		slog.Debug("8bit: list", "label", text, "selected", *selected)
	}
	return res
}
