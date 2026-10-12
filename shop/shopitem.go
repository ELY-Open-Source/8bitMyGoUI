package shop

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Ware is something for sale.
type Ware struct {
	Name string
	// Sprite is at most 16 dots square.
	Sprite *pixel.Sprite
	Price  int
	// Was is the price before a sale; zero means none.
	Was      int
	Currency Currency
}

type ShopItemOptions struct {
	// Short says the wallet cannot pay; BUY turns off.
	Short    bool
	Disabled bool
}

type ShopItemResult struct {
	Element *ui.Element
	Bought  bool
}

// ShopItem is a ware's sprite, name, price and BUY.
func ShopItem(c *ui.Context, w Ware, o ShopItemOptions) ShopItemResult {
	k := pixel.K(c)
	if w.Name == "" || w.Sprite == nil || w.Sprite.W > 16 || w.Sprite.H > 16 {
		panic(fmt.Sprintf("shop: ShopItem %q needs a sprite of 16 dots at most", w.Name))
	}
	var res ShopItemResult
	res.Element = layout.Panel(c, layout.PanelOptions{Style: pixel.Raised}, func() {
		ui.Column(c).Width(k.D(96)).Gap(k.D(3)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(k.D(32), k.D(32)).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Sprite(16-w.Sprite.W, 16-w.Sprite.H, w.Sprite, 0, pixel.Ink, 2)
			})
			name := pixel.Wrap(w.Name, 92, 1)[0]
			ui.Box(c).FillWidth().Height(k.D(pixel.LineHeight)).Role(ui.RoleText).Label(w.Name).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Text((g.W-pixel.Measure(name))/2, 0, name, pixel.Ink, 1)
			})
			PriceTag(c, w.Price, PriceTagOptions{Was: w.Was, Currency: w.Currency})
			res.Bought = input.Button(c, k.Msg("buy"), input.ButtonOptions{Variant: input.Primary, Size: input.Small, Disabled: o.Disabled || o.Short, Label: k.Msg("buy") + " " + w.Name}).Clicked
			// The line is kept, so shelves of cards stay even.
			short := ""
			if o.Short {
				short = k.Msg("not_enough")
			}
			ui.Box(c).Size(k.D(pixel.Measure(short)), k.D(pixel.LineHeight)).Role(ui.RoleText).Label(short).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Text(0, 0, short, pixel.Bad, 1)
			})
		})
	})
	if res.Bought {
		slog.Debug("8bit: shop buy", "ware", w.Name)
	}
	return res
}
