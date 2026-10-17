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
	Element ui.Element
	// Canceled is NO, Escape or a click outside.
	Confirmed, Canceled bool
}

type popKey struct{}

// popping is open, and held the focus last frame.
type popping struct{ open, held bool }

// Popconfirm asks question in a bubble before its button acts.
func Popconfirm(c *ui.Context, text, question string, o PopconfirmOptions) PopconfirmResult {
	k := pixel.K(c)
	wrap := ui.Box(c)
	pop := ui.Local(wrap, popKey{}, func() popping { return popping{} })
	open := &pop.open
	var b input.ButtonResult
	wrap.Children(func() { b = input.Button(c, text, input.ButtonOptions{Disabled: o.Disabled}) })
	if b.Clicked {
		*open = true
	}
	if o.Disabled || pixel.Current(c).Inert {
		// A focus in the closing popup parks on the button.
		if *open && pop.held {
			b.Base.Focus()
			pixel.Park(c, b.Base, true)
		}
		*open = false
	}
	was := *open
	res := PopconfirmResult{Element: wrap}
	pop.held = false
	Popover(c, question, b.Element, open, PopoverOptions{}, func() {
		ui.Column(c).Gap(k.D(4)).Children(func() {
			pixel.Text(c, question, pixel.TextOptions{})
			ui.Row(c).Gap(k.D(4)).Justify(ui.End).Children(func() {
				no := input.Button(c, k.Msg("no"), input.ButtonOptions{Size: input.Small})
				yes := input.Button(c, k.Msg("yes"), input.ButtonOptions{Size: input.Small, Variant: input.Primary})
				if no.Clicked {
					*open = false
				}
				if yes.Clicked {
					res.Confirmed, *open = true, false
				}
				// Focused, not FocusWithin: MyGo links parents after frames.
				pop.held = no.Base.Focused() || yes.Base.Focused()
			})
		})
	})
	res.Canceled = was && !*open && !res.Confirmed
	if res.Confirmed || res.Canceled {
		slog.Debug("8bit: popconfirm", "label", text, "yes", res.Confirmed)
	}
	return res
}
