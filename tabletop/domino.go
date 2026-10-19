package tabletop

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type DominoOptions struct {
	// Vertical stands the tile on end.
	Vertical bool
}

// Domino is a tile of two halves of 0–6 pips.
func Domino(c *ui.Context, a, b int, o DominoOptions) ui.Element {
	k := pixel.K(c)
	if a < 0 || a > 6 || b < 0 || b > 6 {
		panic(fmt.Sprintf("tabletop: Domino %d|%d", a, b))
	}
	w, h := 39, 20
	if o.Vertical {
		w, h = h, w
	}
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(fmt.Sprintf("%d | %d", a, b)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, w, h, pixel.Raised, pixel.Lit)
		if o.Vertical {
			g.HLine(3, 19, 14, pixel.Edge)
			pips(g, 1, 1, 18, a, pixel.Ink)
			pips(g, 1, 20, 18, b, pixel.Ink)
			return
		}
		g.VLine(19, 3, 14, pixel.Edge)
		pips(g, 1, 1, 18, a, pixel.Ink)
		pips(g, 20, 1, 18, b, pixel.Ink)
	})
}
