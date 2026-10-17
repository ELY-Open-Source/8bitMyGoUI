package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

type ConfirmDialogOptions struct {
	Disabled bool
}

type ConfirmDialogResult struct {
	Element ui.Element
	// Canceled is NO, Escape or a click outside.
	Confirmed, Canceled bool
}

type confirmKey struct{}

// ConfirmDialog asks ARE YOU SURE? with YES / NO.
func ConfirmDialog(c *ui.Context, message string, open *bool, o ConfirmDialogOptions) ConfirmDialogResult {
	k := pixel.K(c)
	home := ui.Box(c)
	cur := ui.Local(home, confirmKey{}, func() int { return 1 })
	was := *open
	var res ConfirmDialogResult
	res.Element = ui.DialogBase(c, open, func(back, panel ui.Element) {
		dim(k, back, pixel.Back, 8)
		panel.Width(k.D(180)).Label(k.Msg("sure")).Description(message).Focusable().FocusRing(false).Children(func() {
			pixel.Harbor(c, panel, func() {
				layout.Window(c, layout.WindowOptions{Title: k.Msg("sure")}, func() {
					text.Paragraph(c, message, text.ParagraphOptions{})
					ui.Box(c).Height(k.D(6))
					switch choose(c, k.Msg("sure"), []string{k.Msg("yes"), k.Msg("no")}, cur, false, o.Disabled) {
					case 0:
						res.Confirmed, *open = true, false
					case 1:
						*open = false
					}
				})
			})
		})
	})
	res.Canceled = was && !*open && !res.Confirmed
	if !*open {
		*cur = 1
	}
	if res.Confirmed || res.Canceled {
		slog.Debug("8bit: confirm", "message", message, "yes", res.Confirmed)
	}
	return res
}
