package overlay

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LightboxOptions struct {
	Disabled bool
}

type LightboxResult struct {
	Element *ui.Element
	// Changed turns the picture; Closed is ×, Escape or outside.
	Changed, Closed bool
}

// Lightbox shows pictures over a dark veil, with ◀ ▶.
func Lightbox(c *ui.Context, text string, n int, at *int, open *bool, o LightboxOptions, picture func(i int)) LightboxResult {
	k := pixel.K(c)
	if n < 1 || *at < 0 || *at >= n {
		panic(fmt.Sprintf("overlay: Lightbox picture %d of %d", *at, n))
	}
	was, old := *open, *at
	off := o.Disabled || pixel.Current(c).Inert
	turn := func(d int) { *at = (*at + d + n) % n }
	var res LightboxResult
	res.Element = ui.DialogBase(c, open, func(back, panel *ui.Element) {
		deep, light := dark(k)
		dim(k, back, deep, 13)
		left, right := pixel.Shortcut(panel, 0, ui.KeyLeft), pixel.Shortcut(panel, 0, ui.KeyRight)
		switch {
		case off:
		case left:
			turn(-1)
		case right:
			turn(1)
		}
		panel.Label(text).Focusable().FocusRing(false).AlignItems(ui.Center).Gap(k.D(6)).Children(func() {
			pixel.Harbor(c, panel, func() {
				ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
					if input.IconButton(c, icon.ChevronLeft, k.Msg("prev"), input.IconButtonOptions{Disabled: o.Disabled}).Clicked {
						turn(-1)
					}
					ui.Box(c).Key(*at).Children(func() { picture(*at) })
					if input.IconButton(c, icon.ChevronRight, k.Msg("next"), input.IconButtonOptions{Disabled: o.Disabled}).Clicked {
						turn(1)
					}
				})
				ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
					pixel.Text(c, fmt.Sprintf("%d / %d", *at+1, n), pixel.TextOptions{Slot: light})
					if input.IconButton(c, icon.Close, k.Msg("close"), input.IconButtonOptions{Size: input.Small, Disabled: o.Disabled}).Clicked {
						*open = false
					}
				})
			})
		})
	})
	res.Changed, res.Closed = *at != old, was && !*open
	if res.Changed || res.Closed {
		slog.Debug("8bit: lightbox", "label", text, "at", *at, "closed", res.Closed)
	}
	return res
}
