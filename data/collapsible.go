package data

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CollapsibleOptions struct {
	Disabled bool
}

type CollapsibleResult struct {
	Element *ui.Element
	Changed bool
}

// Collapsible shows content under its ▸ ▾ header while *open.
func Collapsible(c *ui.Context, title string, open *bool, o CollapsibleOptions, content func()) CollapsibleResult {
	k := pixel.K(c)
	col := ui.Column(c).AlignItems(ui.Stretch)
	res := CollapsibleResult{Element: col}
	col.Children(func() {
		// The base turns before it learns it is off.
		v := *open
		shell, base, st := pixel.Shell(c, func() *ui.Element {
			return ui.CollapsibleBase(c, &v).Trigger.Label(title).Disabled(o.Disabled)
		})
		if res.Changed = v != *open && !base.IsDisabled(); res.Changed {
			*open = v
		}
		shown := *open
		shell.Height(k.D(pixel.LineHeight + 6)).MinWidth(k.D(pixel.Measure(title) + 20)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			fill, dy := pixel.Face, 0
			switch {
			case st.Press:
				fill, dy = pixel.Shade, 1
			case st.Hover:
				fill = pixel.Lit
			}
			g.Frame(0, 0, g.W, g.H, pixel.Raised, fill)
			if shown {
				g.Sprite(5, 8+dy, down, 0, pixel.Ink, 1)
			} else {
				g.Sprite(6, 6+dy, pointer, 0, pixel.Ink, 1)
			}
			g.Text(14, 3+dy, title, pixel.Ink, 1)
			g.Finish(st)
		})
		if shown {
			ui.Column(c).Padding(k.D(4), k.D(6)).Children(content)
		}
	})
	if res.Changed {
		slog.Debug("8bit: collapsible", "title", title, "open", *open)
	}
	return res
}
