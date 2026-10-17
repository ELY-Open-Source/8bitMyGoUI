package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ButtonGroupOptions struct {
	Variant  Variant
	Size     Size
	Disabled bool
	// Label names the group for assistive technology.
	Label string
}

type ButtonGroupResult struct {
	Element ui.Element
	// Clicked is the button clicked this frame, or -1.
	Clicked int
}

// ButtonGroup joins buttons as one Tab stop; arrows move.
func ButtonGroup(c *ui.Context, labels []string, o ButtonGroupOptions) ButtonGroupResult {
	k := pixel.K(c)
	h, pad, top, sc := o.Size.metrics()
	res := ButtonGroupResult{Clicked: -1}
	res.Element = ui.Row(c).Role(ui.RoleToolbar).Label(o.Label).FocusGroup(ui.Horizontal).Children(func() {
		for i, l := range labels {
			shell, base, st := pixel.Shell(c, func() ui.Element { return ui.ButtonBase(c).Label(l).Disabled(o.Disabled) })
			w := pixel.Measure(l)*sc + 2*pad
			if i > 0 {
				shell.Margin(0, 0, 0, k.D(-1))
			}
			shell.Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				ink, dy := face(g, o.Variant, st, false)
				g.Text(pad, top+dy, l, ink, sc)
				g.Finish(st)
			})
			if base.Clicked() {
				res.Clicked = i
			}
		}
	})
	if res.Clicked >= 0 {
		slog.Debug("8bit: button group", "label", labels[res.Clicked])
	}
	return res
}
