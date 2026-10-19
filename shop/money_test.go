package shop_test

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/shop"
)

func TestPriceTagStrikesASale(t *testing.T) {
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		ui.Column(c).Children(func() {
			shop.PriceTag(c, 120, shop.PriceTagOptions{})
			shop.PriceTag(c, 75, shop.PriceTagOptions{Was: 150})
		})
	})
	find(t, tt, "120 GOLD")
	x, y, w, _ := dots(find(t, tt, "75 GOLD, WAS 150"))
	// The strike runs Bad through the old price's middle.
	if count(tt.Image(), x+13, y+8, w-13, 1, pixel.Bad) < 10 {
		t.Fatal("the old price is not struck")
	}
	for name, o := range map[string]shop.PriceTagOptions{"a raise": {Was: 50}, "an unknown currency": {Currency: 2}} {
		mustPanic(t, name, func() {
			view(400, 200, newWorld(), func(c *ui.Context) { shop.PriceTag(c, 75, o) })
		})
	}
}

func TestWalletSaysBoth(t *testing.T) {
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		shop.Wallet(c, 950, 3, shop.WalletOptions{Digits: 5})
	})
	find(t, tt, "WALLET")
	mustPanic(t, "a negative price", func() {
		view(400, 200, newWorld(), func(c *ui.Context) { shop.PriceTag(c, -5, shop.PriceTagOptions{}) })
	})
}

func TestShopItemBuysUnlessShort(t *testing.T) {
	short, bought := false, 0
	tt := view(400, 300, newWorld(), func(c *ui.Context) {
		if shop.ShopItem(c, shop.Ware{Name: "POTION", Sprite: icon.Potion, Price: 30}, shop.ShopItemOptions{Short: short}).Bought {
			bought++
		}
	})
	click(t, tt, "BUY POTION")
	short = true
	tt.Frame()
	click(t, tt, "BUY POTION")
	if bought != 1 || !tt.HasText("NOT ENOUGH") {
		t.Fatalf("bought %d times; texts %q", bought, tt.Texts())
	}
}

func TestQuantityStepsWithinLimits(t *testing.T) {
	n := 2
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		shop.Quantity(c, "Potions", &n, shop.QuantityOptions{Min: 1, Max: 3})
	})
	r := find(t, tt, "Potions")
	tt.ClickAt(r.X+2*4, r.Y+r.H/2)
	tt.ClickAt(r.X+2*4, r.Y+r.H/2)
	if n != 1 {
		t.Fatalf("two clicks on - from 2 gave %d", n)
	}
	// Pointer on +, Enter and Space still do nothing.
	tt.Move(r.X+r.W-2*4, r.Y+r.H/2)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyRight)
	if n != 2 {
		t.Fatalf("Enter, Space, Right gave %d", n)
	}
	tt.Key(0, ui.KeyEnd)
	if n != 3 {
		t.Fatalf("End gave %d", n)
	}
	tt.Key(0, ui.KeyUp)
	tt.ClickAt(r.X+r.W-2*4, r.Y+r.H/2)
	if n != 3 {
		t.Fatalf("Up and + past the top gave %d", n)
	}
	mustPanic(t, "a value past Max", func() {
		m := 4
		view(400, 200, newWorld(), func(c *ui.Context) { shop.Quantity(c, "Q", &m, shop.QuantityOptions{Max: 3}) })
	})
}

func TestCartTotalsAndRemoves(t *testing.T) {
	lines := []shop.CartLine{{Name: "POTION", Price: 30, Qty: 3}, {Name: "ETHER", Price: 45, Qty: 1}}
	removed := -1
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		if r := shop.Cart(c, "Cart", lines, shop.CartOptions{}); r.Removed >= 0 {
			removed = r.Removed
		}
	})
	find(t, tt, "135 GOLD")
	// Tab reaches POTION's quantity first; Right adds one.
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if lines[0].Qty != 4 {
		t.Fatalf("Right on POTION's quantity left %+v", lines)
	}
	find(t, tt, "165 GOLD")
	click(t, tt, "REMOVE ETHER")
	if removed != 1 {
		t.Fatalf("REMOVE ETHER gave %d", removed)
	}
	for name, bad := range map[string][]shop.CartLine{
		"two POTION lines": {{Name: "P", Price: 1, Qty: 1}, {Name: "P", Price: 1, Qty: 1}},
		"a total past int": {{Name: "P", Price: 1 << 61, Qty: 8}},
	} {
		mustPanic(t, name, func() {
			view(600, 300, newWorld(), func(c *ui.Context) { shop.Cart(c, "C", bad, shop.CartOptions{}) })
		})
	}
}

func TestReceiptReadsLinesAndTotal(t *testing.T) {
	tt := view(500, 400, newWorld(), func(c *ui.Context) {
		shop.Receipt(c, "ITEM SHOP", []shop.ReceiptLine{{Name: "POTION", Qty: 3, Price: 30}, {Name: "TENT", Qty: 1, Price: 120}}, shop.ReceiptOptions{})
	})
	x, y, w, h := dots(find(t, tt, "ITEM SHOP"))
	// Feed holes run down both edges in Shade.
	if count(tt.Image(), x+3, y, 2, h, pixel.Shade) < 8 || count(tt.Image(), x+w-5, y, 2, h, pixel.Shade) < 8 {
		t.Fatal("no feed holes down the sides")
	}
	mustPanic(t, "a line of none", func() {
		view(500, 400, newWorld(), func(c *ui.Context) {
			shop.Receipt(c, "S", []shop.ReceiptLine{{Name: "A", Qty: 0, Price: 1}}, shop.ReceiptOptions{})
		})
	})
}

func TestSaleBadgeFlashesUnlessReduced(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = reduce
		tt := view(300, 200, wd, func(c *ui.Context) { shop.SaleBadge(c, 50, shop.SaleBadgeOptions{}) })
		x, y, w, h := dots(find(t, tt, "SALE -50%"))
		before := count(tt.Image(), x, y, w, h, pixel.Warn)
		later(tt, wd, 1000e6)
		after := count(tt.Image(), x, y, w, h, pixel.Warn)
		if flashed := before != after; flashed == reduce {
			t.Errorf("reduce %v: Warn dots went %d to %d", reduce, before, after)
		}
	}
	mustPanic(t, "100% off", func() {
		view(300, 200, newWorld(), func(c *ui.Context) { shop.SaleBadge(c, 100, shop.SaleBadgeOptions{}) })
	})
}

func TestCouponCopiesQuietly(t *testing.T) {
	var log string
	tt := view(600, 200, newWorld(), func(c *ui.Context) {
		shop.Coupon(c, "10% OFF", "SECRET42", shop.CouponOptions{})
	})
	log = logged(func() { click(t, tt, "COPY") })
	if tt.Clipboard() != "SECRET42" || !tt.HasText("COPIED!") || strings.Contains(log, "SECRET42") {
		t.Fatalf("clipboard %q; log:\n%s", tt.Clipboard(), log)
	}
}

func TestCouponSaysCopiedForTwoSeconds(t *testing.T) {
	wd := newWorld()
	tt := view(400, 200, wd, func(c *ui.Context) { shop.Coupon(c, "10% OFF", "UPUP", shop.CouponOptions{}) })
	click(t, tt, "COPY")
	for _, tc := range []struct {
		at   time.Duration
		says string
	}{{1500 * time.Millisecond, "COPIED!"}, {2100 * time.Millisecond, "COPY"}} {
		wd.clk.t = wd.clk.t.Add(tc.at - map[string]time.Duration{"COPIED!": 0, "COPY": 1500 * time.Millisecond}[tc.says])
		tt.Frame()
		if !tt.HasText(tc.says) {
			t.Errorf("after %v the button does not say %s", tc.at, tc.says)
		}
	}
}
