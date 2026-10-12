package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type RadioGroupOptions struct {
	// Row lays the radios out across.
	Row      bool
	Disabled bool
	Label    string
}

type radioKey struct{}

type RadioGroupResult struct {
	Element *ui.Element
	Changed bool
}

// RadioGroup chooses *selected; one Tab stop, arrows choose.
func RadioGroup(c *ui.Context, labels []string, selected *int, o RadioGroupOptions) RadioGroupResult {
	k := pixel.K(c)
	var e *ui.Element
	if o.Row {
		e = ui.Row(c).Gap(k.D(10)).FocusGroup(ui.Horizontal)
	} else {
		e = ui.Column(c).Gap(k.D(2)).FocusGroup(ui.Vertical)
	}
	res := RadioGroupResult{Element: e}
	last := ui.Local(e, radioKey{}, func() int { return -1 })
	focus := -1
	e.Role(ui.RoleRadioGroup).Label(o.Label).Children(func() {
		for i, l := range labels {
			was := *selected == i
			on := was
			shell, base, st := pixel.Shell(c, func() *ui.Element {
				return ui.CheckboxBase(c, &on).Role(ui.RoleRadio).Label(l).Disabled(o.Disabled)
			})
			if base.Focused() {
				focus = i
			}
			// Arrow moves show focus and choose; Tab entry does not.
			moved := base.FocusVisible() && *last >= 0 && *last != i
			if !was && (base.Changed() || moved) {
				*selected, res.Changed = i, true
			}
			shell.Size(k.D(12+pixel.Measure(l)), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
				dot(k.Canvas(p, r), l, *selected == i, st)
			})
		}
	})
	*last = focus
	if res.Changed {
		slog.Debug("8bit: radio group", "selected", *selected)
	}
	return res
}
