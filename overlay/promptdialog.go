package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/layout"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type PromptDialogOptions struct {
	Placeholder string
	MaxLength   int
	Disabled    bool
}

type PromptDialogResult struct {
	Element ui.Element
	// Canceled is CANCEL, Escape or a click outside.
	Submitted, Canceled bool
}

// PromptDialog asks for a line, as ENTER YOUR NAME.
func PromptDialog(c *ui.Context, title string, value *string, open *bool, o PromptDialogOptions) PromptDialogResult {
	k := pixel.K(c)
	was := *open
	var res PromptDialogResult
	res.Element = ui.DialogBase(c, open, func(back, panel ui.Element) {
		dim(k, back, pixel.Back, 8)
		panel.Width(k.D(220)).Label(title).Focusable().FocusRing(false).Children(func() {
			pixel.Harbor(c, panel, func() {
				layout.Window(c, layout.WindowOptions{Title: title}, func() {
					ui.Column(c).Gap(k.D(6)).AlignItems(ui.Stretch).Children(func() {
						f := input.TextInput(c, title, value, input.TextInputOptions{Placeholder: o.Placeholder, MaxLength: o.MaxLength, Width: 196, Disabled: o.Disabled})
						pixel.AutoFocus(c, f.Input)
						// A Shortcut, not Submitted: the frame builds again closed.
						if pixel.Shortcut(f.Element, 0, ui.KeyEnter) && !o.Disabled {
							res.Submitted, *open = true, false
						}
						actions(c, []string{k.Msg("cancel"), k.Msg("ok")}, o.Disabled, func(i int) {
							res.Submitted, *open = i == 1, false
						})
					})
				})
			})
		})
	})
	res.Canceled = was && !*open && !res.Submitted
	if res.Submitted || res.Canceled {
		slog.Debug("8bit: prompt", "title", title, "submitted", res.Submitted)
	}
	return res
}
