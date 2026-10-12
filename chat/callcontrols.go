package chat

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CallControlsOptions struct {
	Disabled bool
}

type CallControlsResult struct {
	Element         *ui.Element
	Changed, HungUp bool
}

// CallControls toggles mute and camera, and hangs up.
func CallControls(c *ui.Context, text string, muted, camera *bool, o CallControlsOptions) CallControlsResult {
	k := pixel.K(c)
	var res CallControlsResult
	// A slash crosses a sprite that is off.
	button := func(base func() *ui.Element, sp *pixel.Sprite, fill, ink pixel.Slot, off bool) bool {
		shell, b, st := pixel.Shell(c, func() *ui.Element { return base().Disabled(o.Disabled).FocusRing(false) })
		shell.Size(k.D(22), k.D(22)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, 22, 22, pixel.Raised, fill)
			g.Sprite(5, 5, sp, 0, ink, 1)
			if off {
				g.Line(4, 17, 17, 4, pixel.Bad)
			}
			g.Finish(st)
		})
		return b.Changed() || b.Clicked()
	}
	res.Element = ui.Row(c).Gap(k.D(4)).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).Label(text).Children(func() {
		if button(func() *ui.Element { return ui.ToggleBase(c, muted).Label(k.Msg("mute")) }, icon.Mic, pixel.Face, pixel.Ink, *muted) {
			res.Changed = true
			slog.Debug("8bit: call mute", "label", text, "muted", *muted)
		}
		if button(func() *ui.Element { return ui.ToggleBase(c, camera).Label(k.Msg("camera")) }, icon.Camera, pixel.Face, pixel.Ink, !*camera) {
			res.Changed = true
			slog.Debug("8bit: call camera", "label", text, "on", *camera)
		}
		if button(func() *ui.Element { return ui.ButtonBase(c).Label(k.Msg("hang_up")) }, icon.Phone, pixel.Bad, pixel.AccentInk, false) {
			res.HungUp = true
			slog.Debug("8bit: call hung up", "label", text)
		}
	})
	return res
}
