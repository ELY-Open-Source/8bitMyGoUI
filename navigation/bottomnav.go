package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type BottomNavOptions struct {
	Disabled bool
}

type BottomNavResult struct {
	Element ui.Element
	Changed bool
}

// BottomNav is a bar of labeled icons for narrow screens.
func BottomNav(c *ui.Context, text string, items []NavItem, selected *int, o BottomNavOptions) BottomNavResult {
	k := pixel.K(c)
	res := BottomNavResult{}
	// The bases turn before they learn they are off.
	v := *selected
	res.Element = ui.RadioGroup(c, func() {
		for i, it := range items {
			shell, _, st := pixel.Shell(c, func() ui.Element {
				return ui.RadioBase(c, &v, i).Label(it.Label).Disabled(o.Disabled)
			})
			w := max(pixel.Measure(it.Label), 12) + 12
			shell.Size(k.D(w), k.D(16+pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				ink := pixel.Faint
				if i == *selected {
					ink = pixel.Accent
					g.HLine(2, 0, g.W-4, pixel.Accent)
				} else if st.Hover {
					ink = pixel.Ink
				}
				g.Sprite((g.W-it.Icon.W)/2, 2, it.Icon, 0, ink, 1)
				g.Text((g.W-pixel.Measure(it.Label))/2, 14, it.Label, ink, 1)
				g.Finish(st)
			})
		}
	}).Row().Gap(k.D(4)).Padding(k.D(2), k.D(6)).Label(text).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Face)
		g.HLine(0, 0, g.W, pixel.Edge)
	})
	if res.Changed = keep(c, v, selected, o.Disabled); res.Changed {
		slog.Debug("8bit: bottom nav", "label", text, "selected", *selected)
	}
	return res
}
