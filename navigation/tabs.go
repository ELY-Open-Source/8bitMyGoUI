// Package navigation moves the player between screens and places.
package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TabsOptions struct {
	Disabled bool
}

type TabsResult struct {
	Element ui.Element
	Changed bool
}

// Tabs are folder tabs; the chosen one joins its panel.
func Tabs(c *ui.Context, text string, labels []string, selected *int, o TabsOptions, panel func()) TabsResult {
	return tabs(c, text, labels, selected, false, o, panel)
}

// tabs lays tabs beside a panel; side means left.
func tabs(c *ui.Context, text string, labels []string, selected *int, side bool, o TabsOptions, panel func()) TabsResult {
	k := pixel.K(c)
	th := pixel.LineHeight + 6
	tw := 0
	for _, l := range labels {
		tw = max(tw, pixel.Measure(l)+12)
	}
	// size is tab i's extent along the strip.
	size := func(i int) int {
		if side {
			return th
		}
		return pixel.Measure(labels[i]) + 12
	}
	outer := ui.Column(c)
	if side {
		outer = ui.Row(c).FillWidth().AlignItems(ui.Stretch)
	}
	res := TabsResult{Element: outer}
	outer.Children(func() {
		// The base turns before it learns it is off.
		v := *selected
		p := ui.TabsBase(c, &v, len(labels))
		if side {
			p.List.Column().Width(k.D(tw))
		}
		p.List.Label(text).Children(func() {
			for i, l := range labels {
				shell, _, st := pixel.Shell(c, func() ui.Element { return p.Tab(i).Label(l).Disabled(o.Disabled) })
				w, h := size(i), th
				if side {
					w = tw
				}
				shell.Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
					folder(k.Canvas(p, r), l, i == *selected, side, st)
				})
			}
		})
		res.Changed = keep(c, v, selected, o.Disabled)
		at := 0
		for i := range *selected {
			at += size(i)
		}
		body := ui.Column(c).Padding(k.D(8)).Gap(k.D(4))
		if side {
			body.Grow(1)
		}
		body.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			if g.W < 3 || g.H < 3 {
				return
			}
			g.Frame(0, 0, g.W, g.H, pixel.Raised, pixel.Face)
			// The chosen tab opens into the panel.
			if side {
				g.Rect(0, at+1, 2, th-2, pixel.Face)
			} else {
				g.Rect(at+1, 0, size(*selected)-2, 2, pixel.Face)
			}
		}).Children(panel)
	})
	if res.Changed {
		slog.Debug("8bit: tabs", "label", text, "selected", *selected)
	}
	return res
}

// folder draws a tab: the chosen one stands taller, open-ended.
func folder(g pixel.Canvas, l string, on, side bool, st pixel.State) {
	fill, ink, in := pixel.Shade, pixel.Faint, 2
	if on {
		fill, ink, in = pixel.Face, pixel.Ink, 0
	} else if st.Hover {
		fill = pixel.Lit
	}
	if side {
		g.VLine(in, 1, g.H-2, pixel.Edge)
		g.HLine(in+1, 0, g.W-in-1, pixel.Edge)
		g.HLine(in+1, g.H-1, g.W-in-1, pixel.Edge)
		g.Rect(in+1, 1, g.W-in-1, g.H-2, fill)
		g.Text(in+6, 3, l, ink, 1)
		if !on {
			g.VLine(g.W-1, 0, g.H, pixel.Edge)
		}
	} else {
		g.HLine(1, in, g.W-2, pixel.Edge)
		g.VLine(0, in+1, g.H-in-1, pixel.Edge)
		g.VLine(g.W-1, in+1, g.H-in-1, pixel.Edge)
		g.Rect(1, in+1, g.W-2, g.H-in-1, fill)
		g.Text(6, in+3, l, ink, 1)
		if !on {
			g.HLine(0, g.H-1, g.W, pixel.Edge)
		}
	}
	g.Finish(st)
}

// keep takes v as *selected while the control is on.
func keep(c *ui.Context, v int, selected *int, disabled bool) bool {
	if v == *selected || disabled || pixel.Current(c).Inert {
		return false
	}
	*selected = v
	return true
}
