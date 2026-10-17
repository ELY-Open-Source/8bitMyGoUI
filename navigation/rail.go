package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type NavItem struct {
	Icon  *pixel.Sprite
	Label string
}

type NavRailOptions struct {
	Disabled bool
}

type NavRailResult struct {
	Element ui.Element
	Changed bool
}

// NavRail is an icon rail, named on hover or focus.
func NavRail(c *ui.Context, text string, items []NavItem, selected *int, o NavRailOptions) NavRailResult {
	k := pixel.K(c)
	res := NavRailResult{}
	// The bases turn before they learn they are off.
	v := *selected
	res.Element = ui.RadioGroup(c, func() {
		for i, it := range items {
			shell, _, st := pixel.Shell(c, func() ui.Element {
				return ui.RadioBase(c, &v, i).Label(it.Label).Disabled(o.Disabled)
			})
			shell.Size(k.D(22), k.D(22)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				ink := pixel.Ink
				switch {
				case i == *selected:
					g.Frame(0, 0, g.W, g.H, pixel.Notched, pixel.Accent)
					ink = pixel.AccentInk
				case st.Hover:
					g.Frame(0, 0, g.W, g.H, pixel.Notched, pixel.Lit)
				}
				g.Sprite(5, 5, it.Icon, 0, ink, 1)
				g.Finish(st)
			})
			if st.Hover || st.Focus {
				tip(c, k, shell.Bounds(), it.Label)
			}
		}
	}).AlignSelf(ui.Start).Gap(k.D(2)).Padding(k.D(3)).Label(text).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Face)
		g.VLine(g.W-1, 0, g.H, pixel.Edge)
	})
	if res.Changed = keep(c, v, selected, o.Disabled); res.Changed {
		slog.Debug("8bit: nav rail", "label", text, "selected", *selected)
	}
	return res
}

// tip names an item to its right, letting presses by.
func tip(c *ui.Context, k pixel.Kit, at ui.Rect, s string) {
	ui.Overlay(c, func() {
		w := pixel.Measure(s) + 8
		ui.Box(c).Absolute().Left(at.X+at.W+k.D(4)).Top(at.Y+(at.H-k.D(pixel.LineHeight+2))/2).
			Size(k.D(w), k.D(pixel.LineHeight+2)).PassThrough().Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Notched, pixel.Ink)
			g.Text(4, 1, s, pixel.Back, 1)
		})
	})
}
