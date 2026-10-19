package game

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Ware is an item for sale; Owned counts those held.
type Ware struct {
	Item         data.Item
	Price, Owned int
}

type ShopOptions struct {
	Disabled bool
}

type ShopResult struct {
	Element ui.Element
	// Bought and Sold are this frame's trades, or -1.
	Bought, Sold int
}

// Shop buys at price, sells at half; dear wares fade.
func Shop(c *ui.Context, text string, wares []Ware, gold *int, o ShopOptions) ShopResult {
	k := pixel.K(c)
	if *gold < 0 {
		panic(fmt.Sprintf("game: Shop gold %d", *gold))
	}
	name := 0
	for _, w := range wares {
		if w.Price < 0 || w.Owned < 0 || w.Item.Sprite == nil {
			panic(fmt.Sprintf("game: ware %q costs %d, owned %d", w.Item.Name, w.Price, w.Owned))
		}
		name = max(name, pixel.Measure(w.Item.Name))
	}
	col := ui.Column(c).Gap(k.D(3)).Role(ui.RoleGroup).Label(text)
	res := ShopResult{Element: col, Bought: -1, Sold: -1}
	col.Children(func() {
		CoinCounter(c, k.Msg("gold"), *gold, CoinCounterOptions{})
		for i := range wares {
			w := &wares[i]
			dear := w.Price > *gold
			ink := pixel.Ink
			if dear {
				ink = pixel.Faint
			}
			ui.Row(c).Gap(k.D(5)).AlignItems(ui.Center).Children(func() {
				data.ItemSlot(c, w.Item, data.ItemSlotOptions{})
				ui.Box(c).Width(k.D(name + 4)).Children(func() { pixel.Text(c, w.Item.Name, pixel.TextOptions{Slot: ink}) })
				ui.Row(c).Width(k.D(40)).Gap(k.D(2)).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(k.D(12), k.D(12)).Draw(func(p *ui.Painter, r ui.Rect) {
						k.Canvas(p, r).Sprite(0, 0, icon.Coin, 0, pixel.Ink, 1)
					})
					pixel.Text(c, strconv.Itoa(w.Price), pixel.TextOptions{Slot: ink})
				})
				ui.Box(c).Width(k.D(24)).Children(func() { pixel.Text(c, "×"+strconv.Itoa(w.Owned), pixel.TextOptions{Slot: pixel.Faint}) })
				if input.Button(c, k.Msg("buy"), input.ButtonOptions{Size: input.Small, Label: k.Msg("buy") + " " + w.Item.Name, Disabled: o.Disabled || dear}).Clicked {
					*gold, w.Owned, res.Bought = *gold-w.Price, w.Owned+1, i
				}
				if input.Button(c, k.Msg("sell"), input.ButtonOptions{Size: input.Small, Label: k.Msg("sell") + " " + w.Item.Name, Disabled: o.Disabled || w.Owned == 0}).Clicked {
					*gold, w.Owned, res.Sold = *gold+w.Price/2, w.Owned-1, i
				}
			})
		}
	})
	if res.Bought >= 0 || res.Sold >= 0 {
		slog.Debug("8bit: shop", "label", text, "bought", res.Bought, "sold", res.Sold, "gold", *gold)
	}
	return res
}
