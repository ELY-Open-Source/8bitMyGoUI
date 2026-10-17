package data

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MatrixOptions struct {
	Disabled bool
}

type MatrixResult struct {
	Element ui.Element
	Changed bool
}

type matrixKey struct{}

// Matrix toggles cells under row and column headers.
func Matrix(c *ui.Context, text string, rows, cols []string, on [][]bool, o MatrixOptions) MatrixResult {
	k := pixel.K(c)
	if len(on) != len(rows) {
		panic(fmt.Sprintf("data: Matrix has %d rows, %d of values", len(rows), len(on)))
	}
	cw, rw, ch := 14, 0, pixel.LineHeight+2
	for _, l := range cols {
		cw = max(cw, pixel.Measure(l)+6)
	}
	for i, l := range rows {
		rw = max(rw, pixel.Measure(l)+6)
		if len(on[i]) != len(cols) {
			panic(fmt.Sprintf("data: Matrix row %q has %d values, not %d", l, len(on[i]), len(cols)))
		}
	}
	box := ui.Column(c)
	cur := ui.Local(box, matrixKey{}, func() int { return 0 })
	res := MatrixResult{Element: box}
	head := func(l string, w, h int) {
		ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleText).Label(l).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Text((g.W-pixel.Measure(l))/2, (g.H-pixel.LineHeight)/2, l, pixel.Faint, 1)
		})
	}
	box.Children(func() {
		ui.Row(c).Padding(0, 0, 0, k.D(rw)).Children(func() {
			for _, l := range cols {
				head(l, cw, pixel.LineHeight)
			}
		})
		ui.Row(c).Children(func() {
			ui.Column(c).Children(func() {
				for _, l := range rows {
					head(l, rw, ch)
				}
			})
			n := len(cols)
			*cur = max(0, min(*cur, len(rows)*n-1))
			shell, base, st, picked := pixel.Pad(c, text, n, len(rows)*n, cw, ch, cur, o.Disabled)
			if picked >= 0 {
				on[picked/n][picked%n] = !on[picked/n][picked%n]
				res.Changed = true
			}
			if *cur < len(rows)*n {
				r, cl, state := *cur/n, *cur%n, "off"
				if on[r][cl] {
					state = "on"
				}
				base.Description(rows[r] + " " + cols[cl] + " " + k.Msg(state))
			}
			shell.Draw(func(p *ui.Painter, rc ui.Rect) {
				g := k.Canvas(p, rc)
				for r := range rows {
					for cl := range cols {
						x, y := cl*cw, r*ch
						if st.Focus && r*n+cl == *cur {
							g.Cursor(x, y, cw, ch)
						}
						g.Frame(x+(cw-8)/2, y+(ch-8)/2, 8, 8, pixel.Sunken, pixel.Well)
						if on[r][cl] {
							g.Rect(x+(cw-8)/2+2, y+(ch-8)/2+2, 4, 4, pixel.Accent)
						}
					}
				}
				g.Finish(st)
			})
		})
	})
	if res.Changed {
		slog.Debug("8bit: matrix", "label", text, "cell", *cur)
	}
	return res
}
