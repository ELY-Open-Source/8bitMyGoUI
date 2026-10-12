package game

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LootPopupOptions struct {
	Disabled bool
}

type LootPopupResult struct {
	Element *ui.Element
	// Closed is OK, Escape or a click outside.
	Closed bool
}

// LootPopup announces YOU GOT: an item, rising, until OK.
func LootPopup(c *ui.Context, it data.Item, open *bool, o LootPopupOptions) LootPopupResult {
	k := pixel.K(c)
	if it.Sprite == nil {
		panic(fmt.Sprintf("game: LootPopup item %q has no sprite", it.Name))
	}
	title := fmt.Sprintf(k.Msg("you_got"), it.String())
	w := max(140, pixel.Measure(title)+24)
	r := overlay.Dialog(c, title, open, overlay.DialogOptions{Actions: []string{k.Msg("ok")}, Width: w, Disabled: o.Disabled}, func() {
		sp := it.Sprite
		e := ui.Box(c).Size(k.D(w-24), k.D(3*sp.H+12)).Role(ui.RoleImage).Label(it.String())
		// The item rises 8 dots; reduced, it starts up.
		rise := 8
		if !k.Reduce() {
			rise = k.Since(c, start(k, e, 0).at, 15, 8)
		}
		e.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			x := (g.W - 3*sp.W) / 2
			g.Sprite(x, g.H-3*sp.H-rise-2, sp, 0, pixel.Ink, 3)
			if rise == 8 {
				for _, d := range [][2]int{{-4, 2}, {3 * sp.W, 4}, {-2, 3 * sp.H / 2}, {3*sp.W + 2, 3 * sp.H / 2}} {
					g.Rect(x+d[0], d[1], 2, 2, pixel.Warn)
				}
			}
		})
	})
	return LootPopupResult{Element: r.Element, Closed: r.Action == 0 || r.Closed}
}
