package tabletop

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type PlayingCardOptions struct {
	// Back shows the card face down.
	Back bool
}

// PlayingCard is a card: rank and suit, or its back.
func PlayingCard(c *ui.Context, cd Card, o PlayingCardOptions) ui.Element {
	k := pixel.K(c)
	cd.check("PlayingCard")
	name := cd.name(k)
	if o.Back {
		name = k.Msg("card_back")
	}
	return ui.Box(c).Size(k.D(cardW), k.D(cardH)).Role(ui.RoleImage).Label(name).Draw(func(p *ui.Painter, r ui.Rect) {
		drawCard(k.Canvas(p, r), 0, 0, cd, o.Back)
	})
}
