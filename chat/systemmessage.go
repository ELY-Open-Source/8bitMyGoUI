package chat

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SystemMessageOptions struct {
	// Slot inks the line; zero means Faint.
	Slot pixel.Slot
}

// diamond flanks a system message.
var diamond = pixel.MustSprite(`
.#.
###
.#.
`)

// SystemMessage centers a notice, such as PLAYER 2 HAS JOINED.
func SystemMessage(c *ui.Context, text string, o SystemMessageOptions) *ui.Element {
	k := pixel.K(c)
	ink := o.Slot
	if ink == pixel.Ink {
		ink = pixel.Faint
	}
	w := tagW(text) + 16
	return ui.Box(c).AlignSelf(ui.Center).Role(ui.RoleText).Label(text).Size(k.D(w), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		y := (pixel.LineHeight - diamond.H) / 2
		g.Sprite(0, y, diamond, 0, pixel.Accent, 1)
		tag(g, 8, 0, text, ink)
		g.Sprite(w-diamond.W, y, diamond, 0, pixel.Accent, 1)
	})
}
