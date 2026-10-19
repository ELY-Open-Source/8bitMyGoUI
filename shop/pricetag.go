package shop

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type PriceTagOptions struct {
	// Was is the price before a sale; zero means none.
	Was      int
	Currency Currency
}

// PriceTag is a coin and amount; sales strike the old.
func PriceTag(c *ui.Context, price int, o PriceTagOptions) ui.Element {
	k := pixel.K(c)
	sp, word := coin(k, o.Currency)
	now, was := amount("PriceTag", price), amount("PriceTag Was", o.Was)
	if o.Was != 0 && o.Was <= price {
		panic(fmt.Sprintf("shop: PriceTag sale from %d to %d raises it", o.Was, price))
	}
	say, w, wasW := now+" "+word, 14+pixel.Measure(now), 0
	if o.Was > 0 {
		say += ", " + fmt.Sprintf(k.Msg("was"), was)
		wasW = pixel.Measure(was) + 6
	}
	return ui.Box(c).Size(k.D(w+wasW), k.D(pixel.LineHeight)).Role(ui.RoleText).Label(say).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Sprite(0, 2, sp, 0, pixel.Ink, 1)
		ink := pixel.Ink
		// On sale the old price is struck, the new red.
		if wasW > 0 {
			g.Text(14, 0, was, pixel.Faint, 1)
			g.HLine(13, pixel.LineHeight/2, wasW-4, pixel.Bad)
			ink = pixel.Bad
		}
		g.Text(14+wasW, 0, now, ink, 1)
	})
}
