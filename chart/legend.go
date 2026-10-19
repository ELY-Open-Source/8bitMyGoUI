package chart

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type LegendOptions struct {
	Disabled bool
}

type LegendResult struct {
	Element ui.Element
	Changed bool
}

// Legend names series by swatch; clicks hide or show one.
func Legend(c *ui.Context, text string, series []Series, hidden []bool, o LegendOptions) LegendResult {
	k := pixel.K(c)
	if len(hidden) != len(series) {
		panic(fmt.Sprintf("chart: Legend has %d series and %d hidden flags", len(series), len(hidden)))
	}
	row := ui.Row(c).Gap(k.D(8)).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(text)
	res := LegendResult{Element: row}
	row.Children(func() {
		for i, s := range series {
			shown := !hidden[i]
			shell, base, st := pixel.Shell(c, func() ui.Element { return ui.ToggleBase(c, &shown).Label(s.Name).Disabled(o.Disabled) })
			// The base turns before it learns it is off.
			if base.Changed() && !base.IsDisabled() {
				hidden[i], res.Changed = !shown, true
				slog.Debug("8bit: legend", "label", text, "series", s.Name, "shown", shown)
			}
			shown = !hidden[i]
			ink := pick(s.Slot, i)
			shell.Size(k.D(10+pixel.Measure(s.Name)), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				name := pixel.Ink
				if shown {
					g.Rect(0, 5, 6, 6, ink)
				} else {
					g.Rect(0, 5, 6, 6, pixel.Edge)
					g.Rect(1, 6, 4, 4, pixel.Well)
					name = pixel.Faint
				}
				g.Text(10, 0, s.Name, name, 1)
				g.Finish(st)
			})
		}
	})
	return res
}
