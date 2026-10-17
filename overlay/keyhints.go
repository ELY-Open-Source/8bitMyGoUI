package overlay

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

// KeyHint pairs a control with what it does.
type KeyHint struct {
	Key, Action string
}

type KeyHintsOptions struct{}

// KeyHints shows controls, as Ⓐ CONFIRM Ⓑ BACK.
func KeyHints(c *ui.Context, hints []KeyHint, o KeyHintsOptions) ui.Element {
	k := pixel.K(c)
	return ui.Row(c).Gap(k.D(10)).Padding(k.D(3), k.D(6)).AlignItems(ui.Center).Role(ui.RoleGroup).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Flat, pixel.Face)
	}).Children(func() {
		for _, h := range hints {
			ui.Row(c).Gap(k.D(3)).AlignItems(ui.Center).Role(ui.RoleText).Label(h.Key + " " + h.Action).Children(func() {
				if len([]rune(h.Key)) == 1 && pixel.MicroOK(h.Key) {
					ui.Box(c).Size(k.D(11), k.D(11)).Draw(func(p *ui.Painter, r ui.Rect) {
						g := k.Canvas(p, r)
						g.Disc(5, 5, 5, pixel.Ink)
						g.Micro((11-pixel.MicroWidth(h.Key))/2, (11-pixel.MicroHeight)/2, h.Key, pixel.Back, 1)
					})
				} else {
					text.Kbd(c, h.Key, text.KbdOptions{})
				}
				pixel.Text(c, h.Action, pixel.TextOptions{})
			})
		}
	})
}
