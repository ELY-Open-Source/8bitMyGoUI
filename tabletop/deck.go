package tabletop

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DeckOptions struct {
	Disabled bool
}

type DeckResult struct {
	Element *ui.Element
	// Drawn is a card taken this frame.
	Drawn bool
}

// Deck is a face-down pile, counted; a click draws.
func Deck(c *ui.Context, text string, count int, o DeckOptions) DeckResult {
	k := pixel.K(c)
	// The badge holds three digits.
	if count < 0 || count > 999 {
		panic(fmt.Sprintf("tabletop: Deck of %d", count))
	}
	left := fmt.Sprintf(k.Msg("cards_left"), count)
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleButton).Label(text).Description(left).Disabled(o.Disabled || count == 0)
	})
	res := DeckResult{Element: shell, Drawn: !base.IsDisabled() && base.Clicked()}
	if res.Drawn {
		slog.Debug("8bit: deck drew", "label", text, "left", count-1)
	}
	// A deeper pile shows more edges, six at most.
	layers := min(count, 6)
	shell.Size(k.D(cardW+7), k.D(cardH+7)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		if count == 0 {
			g.Checker(0, 6, cardW, 1, pixel.Edge)
			g.Checker(0, cardH+5, cardW, 1, pixel.Edge)
			g.Checker(0, 6, 1, cardH, pixel.Edge)
			g.Checker(cardW-1, 6, 1, cardH, pixel.Edge)
		}
		for i := range layers {
			off := layers - 1 - i
			drawCard(g, off, 6-off, Card{Rank: 1}, true)
		}
		if count > 0 {
			n := strconv.Itoa(count)
			w := pixel.MicroWidth(n) + 4
			g.Rect(cardW+6-w, cardH, w, 7, pixel.Accent)
			g.Micro(cardW+8-w, cardH+1, n, pixel.AccentInk, 1)
		}
		g.Finish(st)
	})
	return res
}
