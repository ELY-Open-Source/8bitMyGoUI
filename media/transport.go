package media

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TransportOptions struct {
	Disabled bool
}

type TransportResult struct {
	Element *ui.Element
	// Changed is play or pause; the rest are presses.
	Changed, Stopped, Previous, Next bool
}

// Transport is ⏮ ▶/❚❚ ■ ⏭; play toggles *playing.
func Transport(c *ui.Context, text string, playing *bool, o TransportOptions) TransportResult {
	k := pixel.K(c)
	var res TransportResult
	res.Element = ui.Row(c).Gap(k.D(2)).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).Label(text).Children(func() {
		press := func(sp *pixel.Sprite, key string) bool {
			return input.IconButton(c, sp, k.Msg(key), input.IconButtonOptions{Disabled: o.Disabled}).Clicked
		}
		res.Previous = press(prev, "prev_track")
		sp := icon.Play
		if *playing {
			sp = icon.Pause
		}
		res.Changed = input.Toggle(c, "", playing, input.ToggleOptions{Icon: sp, Label: k.Msg("play"), Disabled: o.Disabled}).Changed
		// Stop ends play too.
		if res.Stopped = press(icon.Stop, "stop"); res.Stopped && *playing {
			*playing, res.Changed = false, true
		}
		res.Next = press(icon.Next, "next_track")
	})
	if res.Changed || res.Stopped || res.Previous || res.Next {
		slog.Debug("8bit: transport", "label", text, "playing", *playing, "stopped", res.Stopped, "previous", res.Previous, "next", res.Next)
	}
	return res
}
