package input

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type RatingOptions struct {
	// Max is the number of marks; zero means 5.
	Max int
	// Hearts draws hearts instead of stars.
	Hearts   bool
	Disabled bool
}

type RatingResult struct {
	Element *ui.Element
	Changed bool
}

type ratingKey struct{}

// Rating sets *value by click or arrows; reclicking clears.
func Rating(c *ui.Context, text string, value *int, o RatingOptions) RatingResult {
	k := pixel.K(c)
	n := orDef("Rating Max", o.Max, 5)
	full, empty := icon.Star, icon.StarEmpty
	if o.Hearts {
		full, empty = icon.Heart, icon.HeartEmpty
	}
	e := ui.Row(c).Gap(k.D(1)).FocusGroup(ui.Horizontal).Role(ui.RoleRadioGroup).Label(text)
	last := ui.Local(e, ratingKey{}, func() int { return -1 })
	focus := -1
	res := RatingResult{Element: e}
	e.Children(func() {
		for i := 1; i <= n; i++ {
			was := *value == i
			on := was
			shell, base, st := pixel.Shell(c, func() *ui.Element {
				return ui.CheckboxBase(c, &on).Role(ui.RoleRadio).Label(fmt.Sprintf("%d/%d", i, n)).Disabled(o.Disabled)
			})
			if base.Focused() {
				focus = i
			}
			// Arrow moves show focus and choose; clicks toggle.
			switch {
			case base.Changed() && was:
				*value, res.Changed = 0, true
			case base.Changed(), base.FocusVisible() && *last >= 0 && *last != i && !was:
				*value, res.Changed = i, true
			}
			shell.Size(k.D(full.W+2), k.D(full.H+2)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				sp, ink := empty, pixel.Faint
				if i <= *value {
					sp, ink = full, pixel.Accent
				}
				g.Sprite(1, 1, sp, 0, ink, 1)
				finish(g, st)
			})
		}
	})
	*last = focus
	e.Description(fmt.Sprintf("%d/%d", *value, n))
	if res.Changed {
		slog.Debug("8bit: rating", "label", text, "value", *value)
	}
	return res
}
