package datetime

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MonthPickerOptions struct {
	Disabled bool
}

type MonthPickerResult struct {
	Element ui.Element
	Changed bool
}

type monthKey struct{}

// paging is the year or decade shown, and a cell.
type paging struct {
	at, cur int
}

// MonthPicker lays out twelve months; ◀ ▶ turn years.
func MonthPicker(c *ui.Context, text string, month *time.Time, o MonthPickerOptions) MonthPickerResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(2))
	st := ui.Local(col, monthKey{}, func() paging { return paging{at: month.Year(), cur: int(month.Month()) - 1} })
	res := MonthPickerResult{Element: col}
	// Cells fit the widest name in the frame's locale.
	cw := 28
	for m := range 12 {
		cw = max(cw, pixel.Measure(monthName(k, time.Month(m+1)))+6)
	}
	col.Children(func() {
		pager(c, strconv.Itoa(st.at), 4*cw, o.Disabled, func(d int) {
			st.at += d
			slog.Debug("8bit: month picker turned", "label", text, "year", st.at)
		})
		shell, _, ps, picked := pixel.Pad(c, text, 4, 12, cw, cell, &st.cur, o.Disabled)
		if picked >= 0 {
			*month = time.Date(st.at, time.Month(picked+1), 1, 0, 0, 0, 0, month.Location())
			res.Changed = true
		}
		y, m := month.Year(), int(month.Month())-1
		at, cur := st.at, st.cur
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range 12 {
				x, yy := i%4*cw, i/4*cell
				if ps.Focus && i == cur {
					g.Cursor(x, yy, cw, cell)
				}
				if at == y && i == m {
					g.Frame(x+1, yy, cw-2, cell, pixel.Sunken, pixel.Well)
				}
				name := monthName(k, time.Month(i+1))
				g.Text(x+(cw-pixel.Measure(name))/2, yy, name, pixel.Ink, 1)
			}
			g.Finish(ps)
		})
	})
	if res.Changed {
		slog.Debug("8bit: month picker", "label", text, "month", month.Format("2006-01"))
	}
	return res
}

// pager heads a grid: ◀ title ▶ call turn.
func pager(c *ui.Context, title string, w int, disabled bool, turn func(d int)) {
	k := pixel.K(c)
	ui.Row(c).Width(k.D(w)).AlignItems(ui.Center).Children(func() {
		if input.IconButton(c, icon.ChevronLeft, k.Msg("prev"), input.IconButtonOptions{Size: input.Small, Disabled: disabled}).Clicked {
			turn(-1)
		}
		ui.Row(c).Grow(1).Justify(ui.Center).Children(func() { pixel.Text(c, title, pixel.TextOptions{}) })
		if input.IconButton(c, icon.ChevronRight, k.Msg("next"), input.IconButtonOptions{Size: input.Small, Disabled: disabled}).Clicked {
			turn(1)
		}
	})
}
