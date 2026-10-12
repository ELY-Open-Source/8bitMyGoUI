package tabletop

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CardHandOptions struct {
	Disabled bool
}

type CardHandResult struct {
	Element *ui.Element
	// Played is the card played this frame, else -1.
	Played int
}

type handKey struct{}

// hand is the cursor, and each card's last drawn rise.
type hand struct {
	cur  int
	rise []int
}

// fan is how far apart a hand's cards sit.
const fan = 18

// CardHand fans cards; the pointed or chosen lifts; clicks play.
func CardHand(c *ui.Context, text string, cards []Card, o CardHandOptions) CardHandResult {
	k := pixel.K(c)
	seen := map[Card]bool{}
	for _, cd := range cards {
		cd.check("CardHand")
		if seen[cd] {
			panic(fmt.Sprintf("tabletop: CardHand holds %+v twice", cd))
		}
		seen[cd] = true
	}
	col := ui.Column(c)
	h := ui.Local(col, handKey{}, func() hand { return hand{} })
	cur := &h.cur
	res := CardHandResult{Played: -1}
	n := len(cards)
	w := max(cardW, fan*(n-1)+cardW)
	res.Element = col.Children(func() {
		shell, base, st := pixel.Shell(c, func() *ui.Element {
			return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Disabled(o.Disabled)
		})
		*cur = max(0, min(*cur, n-1))
		// hit is the top card drawn at a point.
		hit := func(x, y float32) int {
			dx, dy := int(x/k.Dot()), int(y/k.Dot())
			for i := n - 1; i >= 0; i-- {
				top := 8
				if i < len(h.rise) {
					top -= h.rise[i]
				}
				if dx >= i*fan && dx < i*fan+cardW && dy >= top && dy < top+cardH {
					return i
				}
			}
			return -1
		}
		lit := -1
		if !base.IsDisabled() && n > 0 {
			for _, m := range []struct {
				key ui.Key
				to  int
			}{{ui.KeyLeft, *cur - 1}, {ui.KeyRight, *cur + 1}, {ui.KeyHome, 0}, {ui.KeyEnd, n - 1}} {
				if pixel.Shortcut(base, 0, m.key) {
					*cur = max(0, min(n-1, m.to))
				}
			}
			// Enter and Space are taken, so clicks come from pointers.
			if pixel.Shortcut(base, 0, ui.KeyEnter) || pixel.Shortcut(base, 0, ui.KeySpace) {
				res.Played = *cur
			}
			x, y, over := base.PointerPosition()
			if i := hit(x, y); base.Clicked() && i >= 0 {
				*cur, res.Played = i, i
			}
			switch {
			case over && st.Hover:
				lit = hit(x, y)
			case st.Focus:
				lit = *cur
			}
		}
		said := make([]string, n)
		for i, cd := range cards {
			said[i] = cd.name(k)
		}
		if n > 0 {
			base.Description(said[*cur] + "; " + strings.Join(said, ", "))
		}
		// Cards rise in three steps, keyed: cards come and go.
		rise := make([]int, n)
		for i, cd := range cards {
			rise[i] = 2 * k.Toward(c, ui.Box(c).Key(cd).Size(0, 0), i == lit, 15, 3)
		}
		h.rise = rise
		shell.Size(k.D(w), k.D(cardH+8)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, cd := range cards {
				drawCard(g, i*fan, 8-rise[i], cd, false)
			}
			g.Finish(st)
		})
	})
	if res.Played >= 0 {
		slog.Debug("8bit: card played", "label", text, "card", cards[res.Played].name(k))
	}
	return res
}
