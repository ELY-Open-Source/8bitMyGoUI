package datetime

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type YearPickerOptions struct {
	Disabled bool
}

type YearPickerResult struct {
	Element ui.Element
	Changed bool
}

type yearKey struct{}

// YearPicker lays out a decade's years; ◀ ▶ turn decades.
func YearPicker(c *ui.Context, text string, year *int, o YearPickerOptions) YearPickerResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(2))
	st := ui.Local(col, yearKey{}, func() paging { return paging{at: *year - *year%10, cur: *year % 10} })
	res := YearPickerResult{Element: col}
	col.Children(func() {
		pager(c, fmt.Sprintf("%d–%d", st.at, st.at+9), 5*28, o.Disabled, func(d int) {
			st.at += 10 * d
			slog.Debug("8bit: year picker turned", "label", text, "decade", st.at)
		})
		shell, _, ps, picked := pixel.Pad(c, text, 5, 10, 28, cell, &st.cur, o.Disabled)
		if picked >= 0 {
			*year, res.Changed = st.at+picked, true
		}
		at, cur, chosen := st.at, st.cur, *year
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range 10 {
				x, y := i%5*28, i/5*cell
				if ps.Focus && i == cur {
					g.Cursor(x, y, 28, cell)
				}
				if at+i == chosen {
					g.Frame(x+1, y, 26, cell, pixel.Sunken, pixel.Well)
				}
				n := strconv.Itoa(at + i)
				g.Text(x+(28-pixel.Measure(n))/2, y, n, pixel.Ink, 1)
			}
			g.Finish(ps)
		})
	})
	if res.Changed {
		slog.Debug("8bit: year picker", "label", text, "year", *year)
	}
	return res
}
