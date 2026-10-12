package shop_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/shop"
)

func TestCheckoutEndsInThanks(t *testing.T) {
	step, done := 0, 0
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		if shop.Checkout(c, "Checkout", &step, shop.CheckoutOptions{}, func() {
			pixel.Text(c, []string{"cart", "pay", "confirm"}[step], pixel.TextOptions{})
		}).Done {
			done++
		}
	})
	click(t, tt, "BACK")
	click(t, tt, "NEXT ▶")
	click(t, tt, "NEXT ▶")
	if step != 2 || !tt.HasText("confirm") {
		t.Fatalf("BACK at the start, then NEXT twice: step %d", step)
	}
	click(t, tt, "PLACE ORDER")
	if step != 3 || done != 1 || !tt.HasText("THANK YOU!") || tt.HasText("confirm") {
		t.Fatalf("PLACE ORDER: step %d, done %d, texts %q", step, done, tt.Texts())
	}
	mustPanic(t, "a fifth step", func() {
		s := 4
		view(600, 300, newWorld(), func(c *ui.Context) { shop.Checkout(c, "C", &s, shop.CheckoutOptions{}, func() {}) })
	})
}

func TestCurrencyConvertTradesAtTheRate(t *testing.T) {
	gold, gems, got := 350, 1, 0
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		got += shop.CurrencyConvert(c, "Swap", &gold, &gems, shop.CurrencyConvertOptions{}).Converted
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	click(t, tt, "CONVERT")
	if gold != 50 || gems != 4 || got != 3 {
		t.Fatalf("End, CONVERT left %d gold, %d gems, converted %d", gold, gems, got)
	}
	if !tt.HasText("NOT ENOUGH") {
		t.Fatalf("50 gold buys no gem: %q", tt.Texts())
	}
	click(t, tt, "CONVERT")
	if gold != 50 {
		t.Fatal("CONVERT spent gold it did not have")
	}
}

func TestPricingTiersChooseAndIgnoreTurnOffKeys(t *testing.T) {
	tiers := [3]shop.Tier{{Price: 0, Perks: []string{"ADS"}}, {Price: 5}, {Price: 9, Perks: []string{"A HAT"}}}
	for _, inert := range []bool{false, true} {
		chosen, off := -1, false
		wd := newWorld()
		tt := view(600, 300, wd, func(c *ui.Context) {
			shop.PricingTiers(c, "Plans", tiers, &chosen, shop.PricingTiersOptions{Disabled: off && !inert})
		})
		click(t, tt, "SILVER, 5 GOLD")
		if chosen != 1 {
			t.Fatalf("clicking SILVER chose %d", chosen)
		}
		// Space, sent while on, lands as the plans turn off.
		tt.Key(0, ui.KeyRight)
		off, wd.inert = true, inert
		tt.Key(0, ui.KeySpace)
		if chosen != 1 {
			t.Errorf("inert %v: Space as the plans turned off chose %d", inert, chosen)
		}
	}
	mustPanic(t, "a fourth plan", func() {
		ch := 3
		view(600, 300, newWorld(), func(c *ui.Context) { shop.PricingTiers(c, "P", tiers, &ch, shop.PricingTiersOptions{}) })
	})
}

func TestDailyRewardClaimsToday(t *testing.T) {
	marks, claims := []bool{true, false, true, false, false, false, false}, 0
	tt := view(700, 300, newWorld(), func(c *ui.Context) {
		if shop.DailyReward(c, "Daily", 4, marks, shop.DailyRewardOptions{}).Claimed {
			claims++
		}
	})
	click(t, tt, "CLAIM")
	click(t, tt, "CLAIM")
	if !marks[3] || claims != 1 || !tt.HasText("SEE YOU TOMORROW!") {
		t.Fatalf("CLAIM: %v, texts %q", marks, tt.Texts())
	}
	mustPanic(t, "a claim from the future", func() {
		future := []bool{false, false, false, false, true, false, false}
		view(700, 300, newWorld(), func(c *ui.Context) { shop.DailyReward(c, "D", 4, future, shop.DailyRewardOptions{}) })
	})
}

func TestBidTickerCrownsTheFirstHighest(t *testing.T) {
	wd := newWorld()
	bids := []shop.Bid{{Bidder: "Aria", Amount: 150, At: wd.clk.t}, {Bidder: "Brom", Amount: 150, At: wd.clk.t}, {Bidder: "Cass", Amount: 90, At: wd.clk.t}}
	tt := view(600, 300, wd, func(c *ui.Context) { shop.BidTicker(c, "Bids", bids, shop.BidTickerOptions{}) })
	img := tt.Image()
	// The crown sits left of the first to bid 150.
	for name, want := range map[string]bool{"Aria": true, "Brom": false, "Cass": false} {
		x, y, _, _ := dots(find(t, tt, name))
		if got := count(img, x-16, y, 12, 12, pixel.Warn) > 0; got != want {
			t.Errorf("%s crowned %v, want %v", name, got, want)
		}
	}
}

func TestTradeWindowNeedsBothAndResetsOnChange(t *testing.T) {
	mine := []data.Item{{Sprite: icon.Sword, Name: "Sword"}}
	var trade shop.Trade
	traded := 0
	tt := view(700, 400, newWorld(), func(c *ui.Context) {
		if shop.TradeWindow(c, "Trade", "Brom", mine, []data.Item{{Sprite: icon.Gem, Name: "Gem", Count: 3}}, &trade, shop.TradeWindowOptions{}).Traded {
			traded++
		}
	})
	click(t, tt, "CONFIRM")
	trade.Theirs = true
	tt.Frame()
	tt.Frame()
	if traded != 1 || !tt.HasText("TRADE DONE!") {
		t.Fatalf("both confirmed: traded %d, texts %q", traded, tt.Texts())
	}
	// A new item on my side takes both confirms back.
	mine = append(mine, data.Item{Sprite: icon.Potion, Name: "Potion"})
	tt.Frame()
	if trade.Mine || trade.Theirs {
		t.Fatalf("a changed offer kept %+v", trade)
	}
	mustPanic(t, "nine items", func() {
		var tr shop.Trade
		view(700, 400, newWorld(), func(c *ui.Context) {
			shop.TradeWindow(c, "T", "B", make([]data.Item, 9), nil, &tr, shop.TradeWindowOptions{})
		})
	})
}

func TestVendingMachineVendsByCode(t *testing.T) {
	var slots [16]shop.Snack
	slots[6] = shop.Snack{Name: "LEAF", Sprite: icon.Leaf, Price: 30, Stock: 1}
	slots[7] = shop.Snack{Name: "DROP", Sprite: icon.Drop, Price: 90, Stock: 0}
	credit, vended := 100, -1
	tt := view(700, 400, newWorld(), func(c *ui.Context) {
		if v := shop.VendingMachine(c, "Machine", slots, &credit, shop.VendingMachineOptions{}).Vended; v >= 0 {
			vended = v
		}
	})
	// Keys sit four a row: A-D, 1-4, CLR, OK.
	r := find(t, tt, "KEYPAD")
	key := func(i int) { tt.ClickAt(r.X+2*float32(i%4*22+11), r.Y+2*float32(i/4*18+9)) }
	key(1)
	key(6)
	key(9)
	if vended != 6 || credit != 70 || !tt.HasText("ENJOY!") {
		t.Fatalf("B3 OK: vended %d, credit %d, texts %q", vended, credit, tt.Texts())
	}
	key(1)
	key(7)
	key(9)
	if !tt.HasText("SOLD OUT") {
		t.Fatalf("B4 has none: %q", tt.Texts())
	}
	key(9)
	if !tt.HasText("CODE?") {
		t.Fatalf("OK with no code: %q", tt.Texts())
	}
}

// offPalette counts dots no palette color matches.
func offPalette(img *image.RGBA) int {
	b := img.Bounds()
	return count(img, 0, 0, b.Dx()/2, b.Dy()/2, none)
}

func TestShopStaysOnThePaletteWhenOff(t *testing.T) {
	builds := map[string]func(c *ui.Context, off bool){
		"item": func(c *ui.Context, off bool) {
			shop.ShopItem(c, shop.Ware{Name: "P", Sprite: icon.Potion, Price: 1}, shop.ShopItemOptions{Disabled: off})
		},
		"quantity": func(c *ui.Context, off bool) {
			n := 1
			shop.Quantity(c, "Q", &n, shop.QuantityOptions{Disabled: off})
		},
		"cart": func(c *ui.Context, off bool) {
			shop.Cart(c, "C", []shop.CartLine{{Name: "P", Price: 1, Qty: 1}}, shop.CartOptions{Disabled: off})
		},
		"checkout": func(c *ui.Context, off bool) {
			s := 1
			shop.Checkout(c, "C", &s, shop.CheckoutOptions{Disabled: off}, func() {})
		},
		"convert": func(c *ui.Context, off bool) {
			g, m := 500, 0
			shop.CurrencyConvert(c, "X", &g, &m, shop.CurrencyConvertOptions{Disabled: off})
		},
		"tiers": func(c *ui.Context, off bool) {
			ch := 0
			shop.PricingTiers(c, "T", [3]shop.Tier{}, &ch, shop.PricingTiersOptions{Disabled: off})
		},
		"coupon": func(c *ui.Context, off bool) {
			shop.Coupon(c, "OFF", "X", shop.CouponOptions{Disabled: off})
		},
		"daily": func(c *ui.Context, off bool) {
			shop.DailyReward(c, "D", 2, make([]bool, 7), shop.DailyRewardOptions{Disabled: off})
		},
		"trade": func(c *ui.Context, off bool) {
			var tr shop.Trade
			shop.TradeWindow(c, "T", "B", nil, nil, &tr, shop.TradeWindowOptions{Disabled: off})
		},
		"vending": func(c *ui.Context, off bool) {
			cr := 0
			shop.VendingMachine(c, "V", [16]shop.Snack{}, &cr, shop.VendingMachineOptions{Disabled: off})
		},
	}
	for name, build := range builds {
		for _, inert := range []bool{false, true} {
			wd := newWorld()
			wd.inert = inert
			wd.clk.t = wd.clk.t.Add(time.Second)
			tt := view(800, 500, wd, func(c *ui.Context) { build(c, !inert) })
			if n := offPalette(tt.Image()); n != 0 {
				t.Errorf("%s, inert %v: %d dots off the palette", name, inert, n)
			}
		}
	}
}
