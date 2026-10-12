package overlay

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

type AlertDialogOptions struct {
	// Icon heads the message; nil means a warning sign.
	Icon *pixel.Sprite
	// Buttons label the answers; nil means OK.
	Buttons []string
	// Default is the button Enter presses.
	Default  int
	Disabled bool
}

type AlertDialogResult struct {
	Element *ui.Element
	// Pressed is the button pressed this frame, or -1.
	Pressed int
	Closed  bool
}

// AlertDialog asks for an answer; Enter presses the default.
func AlertDialog(c *ui.Context, title, message string, open *bool, o AlertDialogOptions) AlertDialogResult {
	k := pixel.K(c)
	buttons, art := o.Buttons, o.Icon
	if buttons == nil {
		buttons = []string{k.Msg("ok")}
	}
	if art == nil {
		art = icon.Warning
	}
	if o.Default < 0 || o.Default >= len(buttons) {
		panic(fmt.Sprintf("overlay: AlertDialog Default %d of %d buttons", o.Default, len(buttons)))
	}
	was := *open
	res := AlertDialogResult{Pressed: -1}
	press := func(i int) { res.Pressed, *open = i, false }
	res.Element = ui.DialogBase(c, open, func(back, panel *ui.Element) {
		dim(k, back, pixel.Back, 8)
		panel.Width(k.D(220)).Role(ui.RoleAlertDialog).Label(title).Description(message).Focusable().FocusRing(false)
		if panel.Shortcut(0, ui.KeyEnter) && !o.Disabled && !pixel.Current(c).Inert {
			press(o.Default)
		}
		panel.Children(func() {
			pixel.Harbor(c, panel, func() {
				layout.Window(c, layout.WindowOptions{Title: title}, func() {
					ui.Row(c).Gap(k.D(8)).Children(func() {
						ui.Box(c).Size(k.D(art.W*2), k.D(art.H*2)).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
							k.Canvas(p, r).Sprite(0, 0, art, 0, pixel.Warn, 2)
						})
						ui.Column(c).Grow(1).Children(func() { text.Paragraph(c, message, text.ParagraphOptions{}) })
					})
					ui.Row(c).Gap(k.D(4)).Justify(ui.End).Padding(k.D(6), 0, 0, 0).Children(func() {
						for i, b := range buttons {
							v := input.Normal
							if i == o.Default {
								v = input.Primary
							}
							if input.Button(c, b, input.ButtonOptions{Size: input.Small, Variant: v, Disabled: o.Disabled}).Clicked {
								press(i)
							}
						}
					})
				})
			})
		})
	})
	res.Closed = was && !*open && res.Pressed < 0
	if res.Pressed >= 0 || res.Closed {
		slog.Debug("8bit: alert", "title", title, "pressed", res.Pressed, "closed", res.Closed)
	}
	return res
}
