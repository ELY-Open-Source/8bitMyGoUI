package agent

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TokenManaOptions struct {
	// Width is in dots; zero means 120.
	Width int
}

// TokenMana shows a token budget's remainder as MP.
func TokenMana(c *ui.Context, text string, used, budget int, o TokenManaOptions) ui.Element {
	k := pixel.K(c)
	if used < 0 || budget < 1 {
		panic(fmt.Sprintf("agent: TokenMana %d of %d", used, budget))
	}
	w := pixel.Or("TokenMana Width", o.Width, 120)
	left := max(0, budget-used)
	mp, value := k.Msg("mp"), tokens(left)+"/"+tokens(budget)
	if need := pixel.MicroWidth(mp) + 4 + pixel.MicroWidth(value); w < need {
		panic(fmt.Sprintf("agent: TokenMana Width %d is under %d", w, need))
	}
	ink, said := pixel.Ink, value
	if left == 0 {
		ink, said = pixel.Bad, value+", "+k.Msg("out_of_mana")
	}
	return ui.Box(c).Size(k.D(w), k.D(15)).Role(ui.RoleMeter).Label(text).Description(said).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Micro(0, 0, mp, pixel.Ink, 1)
		g.Micro(w-pixel.MicroWidth(value), 0, value, ink, 1)
		gauge(g, 0, 7, w, float64(left)/float64(budget), pixel.Info)
		// Ticks mark tenths; an empty well wears a Bad veil.
		for i := 1; i < 10; i++ {
			g.VLine(2+i*(w-4)/10, 9, 4, pixel.Edge)
		}
		if left == 0 {
			g.Checker(2, 9, w-4, 4, pixel.Bad)
		}
	})
}
