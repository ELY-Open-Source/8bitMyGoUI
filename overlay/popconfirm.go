package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PopconfirmOptions struct {
	Disabled bool
}

type PopconfirmResult struct {
	Element *ui.Element
	// Canceled is NO, Escape or a click outside.
	Confirmed, Canceled bool
}

type popKey struct{}

// Popconfirm asks question in a bubble before its button acts.
func Popconfirm(c *ui.Context, text, question string, o PopconfirmOptions) PopconfirmResult {
	k := pixel.K(c)
	wrap := ui.Box(c)
	open := ui.Local(wrap, popKey{}, func() bool { return false })
	var b input.ButtonResult
	wrap.Children(func() { b = input.Button(c, text, input.ButtonOptions{Disabled: o.Disabled}) })
	if b.Clicked {
		*open = true
	}
	// A popup above a veil closes with its control.
	if o.Disabled || pixel.Current(c).Inert {
		*open = false
	}
	was := *open
	res := PopconfirmResult{Element: wrap}
	Popover(c, question, b.Element, open, PopoverOptions{}, func() {
		ui.Column(c).Gap(k.D(4)).Children(func() {
			pixel.Text(c, question, pixel.TextOptions{})
			ui.Row(c).Gap(k.D(4)).Justify(ui.End).Children(func() {
				if input.Button(c, k.Msg("no"), input.ButtonOptions{Size: input.Small}).Clicked {
					*open = false
				}
				if input.Button(c, k.Msg("yes"), input.ButtonOptions{Size: input.Small, Variant: input.Primary}).Clicked {
					res.Confirmed, *open = true, false
				}
			})
		})
	})
	res.Canceled = was && !*open && !res.Confirmed
	if res.Confirmed || res.Canceled {
		slog.Debug("8bit: popconfirm", "label", text, "yes", res.Confirmed)
	}
	return res
}
