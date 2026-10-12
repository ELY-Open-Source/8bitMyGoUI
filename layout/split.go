package layout

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

const gripW = 6

type SplitOptions struct {
	// Vertical stacks the panes top and bottom.
	Vertical bool
	// Min is the smallest pane in dots; zero means 24.
	Min      int
	Label    string
	Disabled bool
}

type SplitResult struct {
	Element *ui.Element
	Changed bool
}

// Split shows two panes and a grip moving *size dots.
func Split(c *ui.Context, size *int, o SplitOptions, first, second func()) SplitResult {
	k := pixel.K(c)
	least := pixel.Or("Split Min", o.Min, 24)
	var e *ui.Element
	var total float32
	back, fwd := ui.KeyLeft, ui.KeyRight
	if o.Vertical {
		e = ui.Column(c)
		total, back, fwd = e.Bounds().H, ui.KeyUp, ui.KeyDown
	} else {
		e = ui.Row(c)
		total = e.Bounds().W
	}
	limit := func(n int) int {
		if total > 0 {
			n = min(n, int(total/k.Dot())-gripW-least)
		}
		return max(n, least)
	}
	start := *size
	e.AlignItems(ui.Stretch).Children(func() {
		pane := ui.Box(c).Clip().Children(first)
		shell, base, st := pixel.Shell(c, func() *ui.Element {
			cur := ui.CursorResizeEW
			if o.Vertical {
				cur = ui.CursorResizeNS
			}
			return ui.Box(c).Focusable().Role(ui.RoleSplitter).Label(label(k, o.Label, "resize")).Disabled(o.Disabled).Cursor(cur)
		})
		if dx, dy, ok := drag(k, base); ok {
			if o.Vertical {
				dx = dy
			}
			*size += dx
		}
		if !o.Disabled && pixel.Shortcut(base, 0, back) {
			*size -= 4
		}
		if !o.Disabled && pixel.Shortcut(base, 0, fwd) {
			*size += 4
		}
		*size = limit(*size)
		if o.Vertical {
			pane.FillWidth().Height(k.D(*size))
			shell.FillWidth().Height(k.D(gripW))
		} else {
			pane.FillHeight().Width(k.D(*size))
			shell.FillHeight().Width(k.D(gripW))
		}
		shell.Draw(func(p *ui.Painter, r ui.Rect) { grip(k.Canvas(p, r), st, o.Vertical) })
		ui.Box(c).Grow(1).Clip().Children(second)
	})
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		now := r.W
		if o.Vertical {
			now = r.H
		}
		if now != total {
			c.Invalidate()
		}
	})
	res := SplitResult{Element: e, Changed: *size != start}
	if res.Changed {
		slog.Debug("8bit: split moved", "size", *size)
	}
	return res
}

// grip is a bar with three bumps across its length.
func grip(g pixel.Canvas, st pixel.State, vertical bool) {
	fill := pixel.Face
	if st.Hover || st.Press {
		fill = pixel.Lit
	}
	g.Rect(0, 0, g.W, g.H, fill)
	for i := -1; i <= 1; i++ {
		x, y := g.W/2-1, g.H/2-1+i*4
		if vertical {
			x, y = g.W/2-1+i*4, g.H/2-1
		}
		g.Dot(x, y, pixel.Lit)
		g.Dot(x+1, y+1, pixel.Shade)
	}
	if st.Disabled {
		g.Checker(0, 0, g.W, g.H, pixel.Face)
	}
	if st.Focus {
		g.FocusMark(0, 0, g.W, g.H)
	}
}
