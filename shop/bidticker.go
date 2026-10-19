package shop

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/datetime"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Bid is an offer by a bidder at a time.
type Bid struct {
	Bidder string
	Amount int
	At     time.Time
}

type BidTickerOptions struct {
	Currency Currency
	// Width is in dots; zero means 240.
	Width int
}

// BidTicker lists bids, newest first, crowning the highest.
func BidTicker(c *ui.Context, text string, bids []Bid, o BidTickerOptions) ui.Element {
	k := pixel.K(c)
	top := -1
	for i, b := range bids {
		amount("BidTicker", b.Amount)
		if b.Bidder == "" || b.At.IsZero() {
			panic(fmt.Sprintf("shop: BidTicker bid %d has no bidder or time", i))
		}
		// The first to reach the highest keeps the crown.
		if top < 0 || b.Amount > bids[top].Amount {
			top = i
		}
	}
	list := ui.Column(c).Width(k.D(pixel.Or("BidTicker Width", o.Width, 240))).Gap(k.D(2)).Role(ui.RoleList).Label(text)
	if top >= 0 {
		_, word := coin(k, o.Currency)
		list.Description(fmt.Sprintf(k.Msg("highest"), bids[top].Bidder, bids[top].Amount, word))
	}
	return list.Children(func() {
		for i := len(bids) - 1; i >= 0; i-- {
			b, crowned := bids[i], i == top
			// New bids come in on top; rows keep keys.
			ui.Row(c).Key(i).FillWidth().Gap(k.D(4)).AlignItems(ui.Center).Role(ui.RoleListItem).Children(func() {
				ui.Box(c).Size(k.D(12), k.D(12)).Draw(func(p *ui.Painter, r ui.Rect) {
					if crowned {
						k.Canvas(p, r).Sprite(0, 0, icon.Crown, 0, pixel.Ink, 1)
					}
				})
				ink := pixel.Ink
				if crowned {
					ink = pixel.Warn
				}
				pixel.Text(c, b.Bidder, pixel.TextOptions{Wrap: true, MaxLines: 1, Slot: ink}).Grow(1)
				money(c, k, o.Currency, b.Amount, ink)
				datetime.RelativeTime(c, b.At, datetime.RelativeTimeOptions{})
			})
		}
		if len(bids) == 0 {
			pixel.Text(c, k.Msg("empty"), pixel.TextOptions{Slot: pixel.Faint})
		}
	})
}
