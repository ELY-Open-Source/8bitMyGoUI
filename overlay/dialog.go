package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DialogOptions struct {
	// Actions label buttons under the content; any one closes.
	Actions []string
	// Width is the window's in dots; zero means 220.
	Width    int
	Disabled bool
}

type DialogResult struct {
	// Element is the window, nil while closed.
	Element ui.Element
	// Action is the button clicked this frame, or -1.
	Action int
	// Closed is Escape or a click outside.
	Closed bool
}

// Dialog is a titled window over a checkered backdrop.
func Dialog(c *ui.Context, title string, open *bool, o DialogOptions, content func()) DialogResult {
	k := pixel.K(c)
	w := pixel.Or("Dialog Width", o.Width, 220)
	was := *open
	res := DialogResult{Action: -1}
	res.Element = ui.DialogBase(c, open, func(back, panel ui.Element) {
		dim(k, back, pixel.Back, 8)
		// The panel holds the focus when nothing inside can.
		panel.Width(k.D(w)).Label(title).Focusable().FocusRing(false).Children(func() {
			pixel.Harbor(c, panel, func() {
				layout.Window(c, layout.WindowOptions{Title: title}, func() {
					ui.Column(c).Gap(k.D(6)).AlignItems(ui.Stretch).Children(func() {
						body(c, o.Disabled, content)
						actions(c, o.Actions, o.Disabled, func(i int) { res.Action, *open = i, false })
					})
				})
			})
		})
	})
	res.Closed = was && !*open && res.Action < 0
	if res.Action >= 0 || res.Closed {
		slog.Debug("8bit: dialog", "title", title, "action", res.Action, "closed", res.Closed)
	}
	return res
}

// actions lines buttons up at the right; act takes clicks.
func actions(c *ui.Context, labels []string, disabled bool, act func(i int)) {
	if len(labels) == 0 {
		return
	}
	ui.Row(c).Gap(pixel.K(c).D(4)).Justify(ui.End).Children(func() {
		for i, l := range labels {
			if input.Button(c, l, input.ButtonOptions{Size: input.Small, Disabled: disabled}).Clicked {
				act(i)
			}
		}
	})
}
