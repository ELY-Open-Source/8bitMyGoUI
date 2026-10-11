package input

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ToggleGroupOptions struct {
	// Single keeps at most one toggle on.
	Single   bool
	Size     Size
	Disabled bool
	Label    string
}

type ToggleGroupResult struct {
	Element *ui.Element
	Changed bool
}

// ToggleGroup joins toggles; on[i] holds each one's state.
func ToggleGroup(c *ui.Context, labels []string, on []bool, o ToggleGroupOptions) ToggleGroupResult {
	if len(on) != len(labels) {
		panic(fmt.Sprintf("input: ToggleGroup has %d labels and %d states", len(labels), len(on)))
	}
	k := pixel.K(c)
	h, pad, top, sc := o.Size.metrics()
	res := ToggleGroupResult{}
	res.Element = ui.Row(c).Role(ui.RoleToolbar).Label(o.Label).FocusGroup(ui.Horizontal).Children(func() {
		for i, l := range labels {
			shell, base, st := pixel.Shell(c, func() *ui.Element { return ui.ToggleBase(c, &on[i]).Label(l).Disabled(o.Disabled) })
			if i > 0 {
				shell.Margin(0, 0, 0, k.D(-1))
			}
			shell.Size(k.D(pixel.Measure(l)*sc+2*pad), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				ink, dy := face(g, Normal, st, on[i])
				g.Text(pad, top+dy, l, ink, sc)
				finish(g, st)
			})
			if base.Changed() {
				res.Changed = true
				if o.Single && on[i] {
					for j := range on {
						on[j] = j == i
					}
				}
			}
		}
	})
	if res.Changed {
		slog.Debug("8bit: toggle group", "on", on)
	}
	return res
}
