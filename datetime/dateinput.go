package datetime

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DateInputOptions struct {
	Disabled bool
}

type DateInputResult struct {
	Element *ui.Element
	Changed bool
}

type dateKey struct{}

// popping is whether the calendar is open and focused.
type popping struct{ open, held bool }

// DateInput shows YYYY-MM-DD; a press opens a calendar under it.
func DateInput(c *ui.Context, text string, date *time.Time, o DateInputOptions) DateInputResult {
	k := pixel.K(c)
	old := *date
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.ButtonBase(c).Role(ui.RolePopUpButton).Label(text).Disabled(o.Disabled)
	})
	pop := ui.Local(base, dateKey{}, func() popping { return popping{} })
	open, was := &pop.open, pop.open
	if base.Clicked() {
		*open = !*open
	}
	// Closing with the field, its focus parks as the field's.
	if o.Disabled || pixel.Current(c).Inert {
		if *open && pop.held {
			base.Focus()
			pixel.Park(c, base, true)
		}
		*open = false
	}
	shown := date.Format("2006-01-02")
	base.Description(shown)
	shell.Size(k.D(pixel.Measure(shown)+24), k.D(pixel.LineHeight+6)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
		g.Text(4, 3, shown, pixel.Ink, 1)
		g.Sprite(g.W-icon.Calendar.W-4, (g.H-icon.Calendar.H)/2, icon.Calendar, 0, pixel.Faint, 1)
		g.Finish(st)
	})
	ui.PopoverBase(c, base, open, func(panel *ui.Element) {
		panel.Padding(k.D(5)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
		}).Children(func() {
			e, grid, picked := month(c, text, date, false, true, func(d time.Time) look { return look{chosen: sameDay(d, *date)} })
			// The grid holds focus before its chain links.
			pop.held = grid.Focused() || e.FocusWithin()
			if picked {
				*open = false
			}
		})
	})
	if *open != was {
		slog.Debug("8bit: date input popup", "label", text, "open", *open)
	}
	res := DateInputResult{Element: shell, Changed: !date.Equal(old)}
	if res.Changed {
		slog.Debug("8bit: date input", "label", text, "date", date.Format("2006-01-02"))
	}
	return res
}
