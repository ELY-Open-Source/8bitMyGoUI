package feedback

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type BusyOverlayOptions struct{}

type busyKey struct{}

// BusyOverlay veils busy content, whose controls then take no input.
func BusyOverlay(c *ui.Context, busy bool, o BusyOverlayOptions, content func()) ui.Element {
	k := pixel.K(c)
	outer := ui.Box(c)
	// held is whether the veil had the focus last frame.
	held := ui.Local(outer, busyKey{}, func() bool { return false })
	if !busy && *held {
		// From the root, parked controls inside take it back.
		c.Root().Focus()
		*held = false
	}
	if busy {
		// Busy content is one Tab stop: the veil.
		outer.FocusGroup(ui.Vertical)
	}
	return outer.Children(func() {
		s := pixel.Current(c)
		s.Inert = s.Inert || busy
		before := c.Root().Focused()
		var body ui.Element
		pixel.Scope(c, s, func() { body = ui.Box(c).Children(content) })
		if !busy {
			return
		}
		veil := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Focusable().FocusRing(false).Role(ui.RoleProgress).Label(k.Msg("wait"))
		if body.FocusWithin() || !before && c.Root().Focused() {
			veil.Focus()
		}
		*held = veil.Focused()
		focus := veil.FocusVisible()
		veil.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Checker(0, 0, g.W, g.H, pixel.Back)
			// The spinner doubles in size where it fits.
			sc := max(1, min(2, min(g.W, g.H)/14))
			beads(g, (g.W-12*sc)/2, (g.H-12*sc)/2, sc)
			if focus {
				g.FocusMark(0, 0, g.W, g.H)
			}
		})
	})
}
