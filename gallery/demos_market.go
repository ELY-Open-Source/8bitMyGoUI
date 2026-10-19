package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/shop"
)

var marketEntries = []Entry{
	{444, "Quantity", "数量", "shop", "- and + within limits. The side at its limit fades.", demoQuantity},
	{445, "PricingTiers", "价格方案", "shop", "Bronze, silver and gold. Gold comes with a hat.", demoPricingTiers},
	{446, "Coupon", "优惠码", "shop", "A ticket with notched edges. COPY takes the code.", demoCoupon},
	{447, "DailyReward", "每日福利", "shop", "Seven days, seven chests. Come back daily; today's is ready.", demoDailyReward},
	{448, "BidTicker", "竞价", "shop", "The latest bids on top; the highest wears the crown.", demoBidTicker},
	{449, "TradeWindow", "交易框", "shop", "Two offers. Both must confirm, and any change takes the confirms back.", demoTradeWindow},
	{450, "VendingMachine", "自动售货机", "shop", "Type a code like B3 and press OK. Credit first, snacks after.", demoVendingMachine},
}

type marketDemo struct {
	qty, tier, credit int
	days              []bool
	mine              []data.Item
	trade             shop.Trade
}

func marketState(c *ui.Context) *marketDemo {
	return ui.Local(c.Root(), demoKey("market"), func() marketDemo {
		return marketDemo{qty: 1, tier: -1, credit: 100, days: []bool{true, false, true, false, false, false, false}, mine: []data.Item{{Sprite: icon.Sword, Name: "Sword"}}}
	})
}

func demoQuantity(c *ui.Context, log func(string)) {
	s := marketState(c)
	if shop.Quantity(c, "Potions", &s.qty, shop.QuantityOptions{Min: 1, Max: 9}).Changed {
		log(fmt.Sprintf("%d potions", s.qty))
	}
}

func demoPricingTiers(c *ui.Context, log func(string)) {
	s := marketState(c)
	if shop.PricingTiers(c, "Plans", [3]shop.Tier{
		{Price: 0, Perks: []string{"1 SAVE SLOT", "ADS"}},
		{Price: 5, Perks: []string{"3 SAVE SLOTS", "NO ADS"}},
		{Price: 9, Perks: []string{"99 SLOTS", "NO ADS", "A HAT"}},
	}, &s.tier, shop.PricingTiersOptions{Currency: shop.Gems}).Changed {
		log([]string{"bronze", "silver", "gold"}[s.tier])
	}
}

func demoCoupon(c *ui.Context, log func(string)) {
	if shop.Coupon(c, "10% OFF", "UPUPDOWNDOWN", shop.CouponOptions{}).Copied {
		log("code copied")
	}
}

func demoDailyReward(c *ui.Context, log func(string)) {
	s := marketState(c)
	if shop.DailyReward(c, "Rewards", 4, s.days, shop.DailyRewardOptions{}).Claimed {
		log("claimed day 4")
	}
}

func demoBidTicker(c *ui.Context, log func(string)) {
	now := pixel.K(c).Now()
	shop.BidTicker(c, "Bids", []shop.Bid{
		{Bidder: "Aria", Amount: 120, At: now.Add(-9 * time.Minute)},
		{Bidder: "Brom", Amount: 150, At: now.Add(-6 * time.Minute)},
		{Bidder: "Cass", Amount: 140, At: now.Add(-40 * time.Second)},
	}, shop.BidTickerOptions{})
}

func demoTradeWindow(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := marketState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		r := shop.TradeWindow(c, "Trade", "Brom", s.mine, []data.Item{{Sprite: icon.Gem, Name: "Gem", Count: 3}}, &s.trade, shop.TradeWindowOptions{})
		if r.Changed {
			log(fmt.Sprintf("you %v, Brom %v", s.trade.Mine, s.trade.Theirs))
		}
		if r.Traded {
			log("traded")
		}
		ui.Row(c).Gap(k.D(6)).Children(func() {
			if input.Button(c, "BROM CONFIRMS", input.ButtonOptions{Size: input.Small}).Clicked {
				s.trade.Theirs = true
				log("Brom confirms")
			}
			if input.Button(c, "ADD POTION", input.ButtonOptions{Size: input.Small}).Clicked {
				s.mine = append(s.mine, data.Item{Sprite: icon.Potion, Name: "Potion"})
				log("you add a potion")
			}
		})
	})
}

// snacks fill the machine: A1 to D4, three of each.
var snacks = func() (out [16]shop.Snack) {
	sprites := []*pixel.Sprite{icon.Potion, icon.Gem, icon.Heart, icon.Star, icon.Coin, icon.Flower, icon.Leaf, icon.Drop, icon.Sun, icon.Moon, icon.Cloud, icon.Snowflake, icon.Bolt, icon.Gift, icon.Book, icon.Key}
	names := []string{"POTION", "GEM", "HEART", "STAR", "COIN", "FLOWER", "LEAF", "DROP", "SUN", "MOON", "CLOUD", "SNOW", "BOLT", "GIFT", "BOOK", "KEY"}
	for i := range out {
		out[i] = shop.Snack{Name: names[i], Sprite: sprites[i], Price: 10 + i%4*10, Stock: 3 - i%5/4*3}
	}
	return out
}()

func demoVendingMachine(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := marketState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		if i := shop.VendingMachine(c, "Machine", snacks, &s.credit, shop.VendingMachineOptions{}).Vended; i >= 0 {
			log(fmt.Sprintf("vended %s, credit %d", snacks[i].Name, s.credit))
		}
		if input.Button(c, "INSERT COIN", input.ButtonOptions{Size: input.Small}).Clicked {
			s.credit += 10
			log(fmt.Sprintf("credit %d", s.credit))
		}
	})
}
