package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/shop"
)

var shopEntries = []Entry{
	{436, "PriceTag", "价签", "shop", "A coin and a number. On sale, the old price gets a line through it.", demoPriceTag},
	{437, "Wallet", "钱包", "shop", "Gold and gems. Earn some and watch the digits roll.", demoWallet},
	{438, "ShopItem", "商品卡", "shop", "A ware on the shelf: sprite, name, price, BUY. Short of gold, BUY sits out.", demoShopItem},
	{439, "Cart", "购物车", "shop", "What you are about to buy, how many, and how much it all costs.", demoCart},
	{440, "Receipt", "小票", "shop", "Printed dot by dot, holes down the sides, torn at the bottom.", demoReceipt},
	{441, "Checkout", "结账", "shop", "CART, PAY, CONFIRM, and the shopkeeper says THANK YOU!", demoCheckout},
	{442, "CurrencyConvert", "兑换", "shop", "A hundred gold buys a gem. Pick how many and CONVERT.", demoCurrencyConvert},
	{443, "SaleBadge", "促销角标", "shop", "SALE on a starburst, rim flashing, as shops have always shouted.", demoSaleBadge},
}

type shopDemo struct {
	gold, gems, step int
	cart             []shop.CartLine
}

func shopState(c *ui.Context) *shopDemo {
	return ui.Local(c.Root(), demoKey("shop"), func() shopDemo {
		return shopDemo{gold: 950, gems: 3, cart: []shop.CartLine{
			{Name: "POTION", Price: 30, Qty: 3}, {Name: "ETHER", Price: 45, Qty: 1}, {Name: "PHOENIX FEATHER", Price: 200, Qty: 2},
		}}
	})
}

func demoPriceTag(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		shop.PriceTag(c, 120, shop.PriceTagOptions{})
		shop.PriceTag(c, 75, shop.PriceTagOptions{Was: 150})
		shop.PriceTag(c, 3, shop.PriceTagOptions{Currency: shop.Gems})
	})
}

func demoWallet(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := shopState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		shop.Wallet(c, s.gold, s.gems, shop.WalletOptions{Digits: 5})
		if input.Button(c, "FIND GOLD", input.ButtonOptions{Size: input.Small}).Clicked {
			s.gold += 137
			log(fmt.Sprintf("gold %d", s.gold))
		}
	})
}

func demoShopItem(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := shopState(c)
	ui.Row(c).Gap(k.D(6)).Children(func() {
		for _, w := range []shop.Ware{
			{Name: "POTION", Sprite: icon.Potion, Price: 30},
			{Name: "SWORD", Sprite: icon.Sword, Price: 3000},
			{Name: "SHIELD", Sprite: icon.Shield, Price: 80, Was: 120},
		} {
			if shop.ShopItem(c, w, shop.ShopItemOptions{Short: s.gold < w.Price}).Bought {
				s.gold -= w.Price
				log(fmt.Sprintf("bought %s, %d gold left", w.Name, s.gold))
			}
		}
	})
}

func demoCart(c *ui.Context, log func(string)) {
	s := shopState(c)
	r := shop.Cart(c, "Cart", s.cart, shop.CartOptions{})
	if r.Changed {
		log("quantity changed")
	}
	if r.Removed >= 0 {
		log("removed " + s.cart[r.Removed].Name)
		s.cart = append(s.cart[:r.Removed:r.Removed], s.cart[r.Removed+1:]...)
	}
}

func demoReceipt(c *ui.Context, log func(string)) {
	shop.Receipt(c, "ITEM SHOP", []shop.ReceiptLine{
		{Name: "POTION", Qty: 3, Price: 30}, {Name: "ETHER", Qty: 1, Price: 45}, {Name: "TENT", Qty: 1, Price: 120},
	}, shop.ReceiptOptions{Time: time.Date(1997, 1, 31, 18, 42, 0, 0, time.UTC)})
}

func demoCheckout(c *ui.Context, log func(string)) {
	s := shopState(c)
	r := shop.Checkout(c, "Checkout", &s.step, shop.CheckoutOptions{}, func() {
		pixel.Text(c, []string{"3 ITEMS, 375 GOLD", "PAY WITH: GOLD", "DELIVER TO: THE INN"}[s.step], pixel.TextOptions{})
	})
	if r.Changed {
		log(fmt.Sprintf("step %d", s.step+1))
	}
	if r.Done {
		log("order placed")
	}
}

func demoCurrencyConvert(c *ui.Context, log func(string)) {
	s := shopState(c)
	if n := shop.CurrencyConvert(c, "Exchange", &s.gold, &s.gems, shop.CurrencyConvertOptions{}).Converted; n > 0 {
		log(fmt.Sprintf("bought %d gems, %d gold left", n, s.gold))
	}
}

func demoSaleBadge(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(8)).Children(func() {
		shop.SaleBadge(c, 50, shop.SaleBadgeOptions{})
		shop.SaleBadge(c, 10, shop.SaleBadgeOptions{})
	})
}
