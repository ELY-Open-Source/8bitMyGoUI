package datetime

import (
	"fmt"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// cell is a day's square in dots.
const cell = pixel.LineHeight

// look is how one day of a month grid draws.
type look struct{ chosen, between bool }

// addMonths moves t n months, keeping the day inside.
func addMonths(t time.Time, n int) time.Time {
	y, m, d := t.Date()
	first := time.Date(y, m+time.Month(n), 1, 0, 0, 0, 0, t.Location())
	last := first.AddDate(0, 1, -1).Day()
	return onDay(t, time.Date(first.Year(), first.Month(), min(d, last), 0, 0, 0, 0, t.Location()))
}

// month moves *cur by keys; clicks, Enter and Space pick.
func month(c *ui.Context, text string, cur *time.Time, disabled, focus bool, mark func(time.Time) look) (e, grid *ui.Element, picked bool) {
	k := pixel.K(c)
	e = ui.Column(c).Gap(k.D(2)).Role(ui.RoleGroup).Label(text)
	e.Children(func() {
		ui.Row(c).Width(k.D(7 * cell)).AlignItems(ui.Center).Children(func() {
			if input.IconButton(c, icon.ChevronLeft, k.Msg("prev"), input.IconButtonOptions{Size: input.Small, Disabled: disabled}).Clicked {
				*cur = addMonths(*cur, -1)
			}
			ui.Row(c).Grow(1).Justify(ui.Center).Children(func() {
				pixel.Text(c, fmt.Sprintf("%s %d", monthName(k, cur.Month()), cur.Year()), pixel.TextOptions{})
			})
			if input.IconButton(c, icon.ChevronRight, k.Msg("next"), input.IconButtonOptions{Size: input.Small, Disabled: disabled}).Clicked {
				*cur = addMonths(*cur, 1)
			}
		})
		names := weekdays(k)
		ui.Box(c).Size(k.D(7*cell), k.D(pixel.LineHeight)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, n := range names {
				// Micro letters leave room between two-letter names.
				if pixel.MicroOK(n) {
					g.Micro(i*cell+(cell-pixel.MicroWidth(n))/2, 6, n, pixel.Faint, 1)
				} else {
					g.Text(i*cell+(cell-pixel.Measure(n))/2, 0, n, pixel.Faint, 1)
				}
			}
		})
		shell, base, st := pixel.Shell(c, func() *ui.Element {
			return ui.Box(c).Focusable().Role(ui.RoleTable).Label(text).Disabled(disabled)
		})
		grid = base
		if focus {
			base.AutoFocus()
		}
		y, m, d := cur.Date()
		for _, mv := range []struct {
			key ui.Key
			to  time.Time
		}{
			{ui.KeyLeft, cur.AddDate(0, 0, -1)}, {ui.KeyRight, cur.AddDate(0, 0, 1)},
			{ui.KeyUp, cur.AddDate(0, 0, -7)}, {ui.KeyDown, cur.AddDate(0, 0, 7)},
			{ui.KeyPageUp, addMonths(*cur, -1)}, {ui.KeyPageDown, addMonths(*cur, 1)},
			{ui.KeyHome, cur.AddDate(0, 0, 1-d)}, {ui.KeyEnd, addMonths(cur.AddDate(0, 0, 1-d), 1).AddDate(0, 0, -1)},
		} {
			if base.Shortcut(0, mv.key) {
				*cur = mv.to
			}
		}
		enter, space := base.Shortcut(0, ui.KeyEnter), base.Shortcut(0, ui.KeySpace)
		picked = enter || space
		first := time.Date(y, m, 1, 0, 0, 0, 0, cur.Location())
		first = first.AddDate(0, 0, -int(first.Weekday()))
		// Enter and Space are taken, so clicks come from pointers.
		if base.Clicked() {
			px, py, _ := base.PointerPosition()
			if col, row := int(px/k.D(cell)), int(py/k.D(cell)); px >= 0 && py >= 0 && col < 7 && row < 6 {
				*cur, picked = onDay(*cur, first.AddDate(0, 0, row*7+col)), true
			}
		}
		base.Description(cur.Format("2006-01-02"))
		now, shown := today(c, k), cur.Month()
		shell.Size(k.D(7*cell), k.D(6*cell)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range 42 {
				day, x, y := first.AddDate(0, 0, i), i%7*cell, i/7*cell
				lk, ink := mark(day), pixel.Ink
				if day.Month() != shown {
					ink = pixel.Faint
				}
				if lk.between {
					g.Veil(x, y, cell, cell, pixel.Accent, 4)
				}
				if lk.chosen {
					g.Frame(x, y, cell, cell, pixel.Sunken, pixel.Well)
				}
				n := strconv.Itoa(day.Day())
				g.Text(x+(cell-pixel.Measure(n))/2, y, n, ink, 1)
				if sameDay(day, now) {
					g.HLine(x+3, y+cell-2, cell-6, pixel.Accent)
				}
				if st.Focus && sameDay(day, *cur) {
					g.Rect(x, y, cell, 1, pixel.Accent)
					g.Rect(x, y+cell-1, cell, 1, pixel.Accent)
					g.Rect(x, y, 1, cell, pixel.Accent)
					g.Rect(x+cell-1, y, 1, cell, pixel.Accent)
				}
			}
			g.Finish(st)
		})
	})
	return e, grid, picked
}
