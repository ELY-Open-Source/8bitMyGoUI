package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/feedback"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ProgressDialogOptions struct {
	Disabled bool
}

type ProgressDialogResult struct {
	Element *ui.Element
	// Canceled is CANCEL, Escape or a click outside.
	Canceled bool
}

// ProgressDialog holds the view on a progress bar with CANCEL.
func ProgressDialog(c *ui.Context, title string, value float64, open *bool, o ProgressDialogOptions) ProgressDialogResult {
	k := pixel.K(c)
	was := *open
	var res ProgressDialogResult
	res.Element = ui.DialogBase(c, open, func(back, panel *ui.Element) {
		dim(k, back, pixel.Back, 8)
		panel.Width(k.D(200)).Label(title).Focusable().FocusRing(false).Children(func() {
			pixel.Harbor(c, panel, func() {
				layout.Window(c, layout.WindowOptions{Title: title}, func() {
					ui.Column(c).Gap(k.D(6)).AlignItems(ui.Stretch).Children(func() {
						feedback.Progress(c, title, value, feedback.ProgressOptions{})
						actions(c, []string{k.Msg("cancel")}, o.Disabled, func(int) { *open = false })
					})
				})
			})
		})
	})
	res.Canceled = was && !*open
	if res.Canceled {
		slog.Debug("8bit: progress canceled", "title", title, "value", value)
	}
	return res
}
