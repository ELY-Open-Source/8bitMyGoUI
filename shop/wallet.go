package shop

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

type WalletOptions struct {
	// Digits pads each balance; zero fits the number.
	Digits int
}

// Wallet shows gold and gem balances that roll on change.
func Wallet(c *ui.Context, gold, gems int, o WalletOptions) *ui.Element {
	k := pixel.K(c)
	g, m := amount("Wallet gold", gold), amount("Wallet gems", gems)
	_, gw := coin(k, Gold)
	_, mw := coin(k, Gems)
	return ui.Column(c).Gap(k.D(2)).Role(ui.RoleGroup).Label(k.Msg("wallet")).Description(g + " " + gw + ", " + m + " " + mw).Children(func() {
		for _, row := range []struct {
			cur Currency
			n   int
		}{{Gold, gold}, {Gems, gems}} {
			sp, _ := coin(k, row.cur)
			ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(k.D(12), k.D(12)).Draw(func(p *ui.Painter, r ui.Rect) {
					k.Canvas(p, r).Sprite(0, 0, sp, 0, pixel.Ink, 1)
				})
				text.Odometer(c, row.n, text.OdometerOptions{Digits: o.Digits})
			})
		}
	})
}
