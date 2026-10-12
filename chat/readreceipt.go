package chat

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ReadReceiptOptions struct {
	// Scale is 1–4; zero means 1.
	Scale int
}

var check = pixel.MustSprite(`
....#
...##
#.##.
###..
.#...
`)

// ReadReceipt shows one check once sent, two once read.
func ReadReceipt(c *ui.Context, read bool, o ReadReceiptOptions) *ui.Element {
	k := pixel.K(c)
	sc := pixel.Or("ReadReceipt Scale", o.Scale, 1)
	if sc > 4 {
		panic(fmt.Sprintf("chat: ReadReceipt Scale %d outside 1–4", sc))
	}
	label, ink, w := k.Msg("sent"), pixel.Faint, check.W
	if read {
		label, ink, w = k.Msg("seen"), pixel.Info, check.W+3
	}
	return ui.Box(c).Role(ui.RoleImage).Label(label).Size(k.D(w*sc), k.D(check.H*sc)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Sprite(0, 0, check, 0, ink, sc)
		if read {
			g.Sprite(3*sc, 0, check, 0, ink, sc)
		}
	})
}
