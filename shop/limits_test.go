package shop_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/shop"
)

func TestTradeWindowSeesEveryChangeAndLogsTheTrade(t *testing.T) {
	mine := []data.Item{{Sprite: icon.Sword, Name: "A", Count: 1}}
	trade := shop.Trade{Mine: true, Theirs: true}
	tt := view(700, 400, newWorld(), func(c *ui.Context) {
		shop.TradeWindow(c, "Trade", "Brom", mine, nil, &trade, shop.TradeWindowOptions{})
	})
	// Same name and count, new sprite: a new offer.
	trade = shop.Trade{Mine: true, Theirs: true}
	mine[0].Sprite = icon.Shield
	tt.Frame()
	if trade.Mine || trade.Theirs {
		t.Fatalf("a swapped sprite kept %+v", trade)
	}
	trade.Mine = true
	tt.Frame()
	log := logged(func() {
		trade.Theirs = true
		tt.Frame()
	})
	if !strings.Contains(log, "trade") {
		t.Fatalf("the partner's confirm finished the trade unlogged:\n%s", log)
	}
}

func TestLimitsHoldAtTheEdgeOfInt(t *testing.T) {
	lines := []shop.CartLine{{Name: "BIG", Price: math.MaxInt/2 + 1, Qty: 1}}
	tt := view(600, 300, newWorld(), func(c *ui.Context) { shop.Cart(c, "Cart", lines, shop.CartOptions{}) })
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if lines[0].Qty != 1 {
		t.Fatalf("a raise past int took the cart to %d", lines[0].Qty)
	}
	n := math.MaxInt
	tt = view(600, 200, newWorld(), func(c *ui.Context) {
		shop.Quantity(c, "Q", &n, shop.QuantityOptions{Max: math.MaxInt})
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if n != math.MaxInt {
		t.Fatalf("+ at the top wrapped to %d", n)
	}
}

func TestCheckoutTurnsItsContentOff(t *testing.T) {
	step, pressed := 0, 0
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		shop.Checkout(c, "Checkout", &step, shop.CheckoutOptions{Disabled: true}, func() {
			if input.Button(c, "EDIT CART", input.ButtonOptions{}).Clicked {
				pressed++
			}
		})
	})
	click(t, tt, "EDIT CART")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if pressed != 0 {
		t.Fatalf("a disabled checkout's content took %d presses", pressed)
	}
}

func TestVendingKeysSpeakTheLocale(t *testing.T) {
	credit := 0
	tt := view(700, 400, newWorld(), func(c *ui.Context) {
		shop.VendingMachine(c, "Machine", [16]shop.Snack{}, &credit, shop.VendingMachineOptions{})
	})
	x, y, _, _ := dots(find(t, tt, "KEYPAD"))
	// OK is the last key: a Good check.
	if count(tt.Image(), x+22, y+36, 22, 18, pixel.Good) == 0 {
		t.Fatal("OK is not drawn as a check")
	}
}

func TestCheckoutKeepsFocusAndVeilsWhenOff(t *testing.T) {
	step := 0
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		shop.Checkout(c, "Checkout", &step, shop.CheckoutOptions{}, func() {
			// Step two says more, which would move the buttons.
			for range step + 1 {
				pixel.Text(c, "line", pixel.TextOptions{})
			}
		})
	})
	click(t, tt, "NEXT ▶")
	tt.Key(0, ui.KeyEnter)
	if step != 2 {
		t.Fatalf("NEXT, then Enter on it, reached step %d", step)
	}
	shot := func(off bool) int {
		tt := view(600, 300, newWorld(), func(c *ui.Context) {
			s := 0
			shop.Checkout(c, "Checkout", &s, shop.CheckoutOptions{Disabled: off}, func() {
				input.Button(c, "EDIT CART", input.ButtonOptions{})
			})
		})
		x, y, w, h := dots(find(t, tt, "EDIT CART"))
		return count(tt.Image(), x, y, w, h, pixel.Face)
	}
	// Off, the content wears the Face checker veil.
	if on, off := shot(false), shot(true); on == off {
		t.Fatalf("EDIT CART looks the same off: %d Face dots", on)
	}
}

func TestCheckoutStepsStartFresh(t *testing.T) {
	step, sel := 0, 0
	rows := make([]string, 60)
	for i := range rows {
		rows[i] = fmt.Sprintf("row %d", i)
	}
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		shop.Checkout(c, "Checkout", &step, shop.CheckoutOptions{}, func() {
			// Each step shows a long list in the same place.
			names := func(i int) string { return fmt.Sprint(step, rows[i]) }
			data.List(c, "List", len(rows), &sel, data.ListOptions{Height: 60}, names, func(i int) {
				pixel.Text(c, names(i), pixel.TextOptions{})
			})
		})
	})
	click(t, tt, "List")
	tt.Key(0, ui.KeyEnd)
	click(t, tt, "NEXT ▶")
	// The new step's list starts at its top.
	if _, ok := tt.Find("1row 0"); !ok {
		t.Fatalf("step two's list kept the old scroll: %q", tt.Texts())
	}
}
