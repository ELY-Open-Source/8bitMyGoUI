package display

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type SignpostOptions struct {
	Disabled bool
}

type SignpostResult struct {
	Element ui.Element
	// Clicked is the board clicked this frame, or -1.
	Clicked int
}

// Signpost is a post with one to three arrow boards.
func Signpost(c *ui.Context, boards []string, o SignpostOptions) SignpostResult {
	if len(boards) < 1 || len(boards) > 3 {
		panic(fmt.Sprintf("display: Signpost needs 1–3 boards, got %d", len(boards)))
	}
	k := pixel.K(c)
	w := 0
	for _, b := range boards {
		w = max(w, pixel.Measure(b))
	}
	w += 18
	res := SignpostResult{Clicked: -1}
	res.Element = ui.Column(c).Width(k.D(w+8)).Gap(k.D(3)).AlignItems(ui.Center).Padding(0, 0, k.D(20), 0).Role(ui.RoleGroup).FocusGroup(ui.Vertical).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(g.W/2-2, 2, 4, g.H-2, pixel.Shade)
		g.VLine(g.W/2-2, 2, g.H-2, pixel.Edge)
		g.VLine(g.W/2+1, 2, g.H-2, pixel.Edge)
	}).Children(func() {
		for i, b := range boards {
			if board(c, k, b, w, i%2 == 1, o.Disabled) {
				res.Clicked = i
			}
		}
	})
	if res.Clicked >= 0 {
		slog.Debug("8bit: signpost", "board", boards[res.Clicked])
	}
	return res
}

// board is one wooden arrow, pointing left if asked.
func board(c *ui.Context, k pixel.Kit, label string, w int, left, disabled bool) bool {
	shell, base, st := pixel.Shell(c, func() ui.Element { return ui.ButtonBase(c).Label(label).Disabled(disabled) })
	h := pixel.LineHeight + 2
	shell.Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		dy := 0
		if st.Press {
			dy = 1
		}
		for j := range h - dy {
			cut := abs(j - (h-1)/2)
			x0, x1 := 0, w-cut
			if left {
				x0, x1 = cut, w
			}
			fill := pixel.Warn
			if j == 0 || j == h-1-dy {
				fill = pixel.Edge
			}
			g.HLine(x0, j+dy, x1-x0, fill)
			g.Dot(x0, j+dy, pixel.Edge)
			g.Dot(x1-1, j+dy, pixel.Edge)
		}
		g.Veil(1, 1+dy, w-2, h-3, pixel.Shade, 4)
		tx := 6
		if left {
			tx = w - 6 - pixel.Measure(label)
		}
		g.Text(tx, 1+dy, label, pixel.Ink, 1)
		if st.Hover {
			g.Veil(1, 1+dy, w-2, h-3, pixel.Lit, 4)
		}
		if st.Disabled {
			g.Checker(0, 0, w, h, pixel.Face)
		}
		if st.Focus {
			g.FocusMark(0, 0, w, h)
		}
	})
	return base.Clicked()
}
