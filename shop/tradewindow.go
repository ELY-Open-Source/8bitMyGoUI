package shop

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Trade is who has confirmed it.
type Trade struct {
	Mine, Theirs bool
}

type TradeWindowOptions struct {
	Disabled bool
}

type TradeWindowResult struct {
	Element *ui.Element
	Changed bool
	// Traded reports both confirms landing this frame.
	Traded bool
}

type tradeKey struct{}

// trade is the offers last seen and if both confirmed.
type trade struct {
	mine, theirs []data.Item
	both         bool
}

const tradeSlots = 8

// TradeWindow shows two offers; both sides must confirm.
func TradeWindow(c *ui.Context, text, partner string, mine, theirs []data.Item, t *Trade, o TradeWindowOptions) TradeWindowResult {
	k := pixel.K(c)
	if partner == "" || len(mine) > tradeSlots || len(theirs) > tradeSlots {
		panic(fmt.Sprintf("shop: TradeWindow with %q, %d and %d items", partner, len(mine), len(theirs)))
	}
	col := ui.Column(c).Gap(k.D(6)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(text)
	st := ui.Local(col, tradeKey{}, func() trade { return trade{mine: slices.Clone(mine), theirs: slices.Clone(theirs)} })
	var res TradeWindowResult
	// Any change to an offer withdraws both confirms.
	if !slices.Equal(mine, st.mine) || !slices.Equal(theirs, st.theirs) {
		st.mine, st.theirs = slices.Clone(mine), slices.Clone(theirs)
		if t.Mine || t.Theirs {
			t.Mine, t.Theirs, res.Changed = false, false, true
			slog.Debug("8bit: trade offer changed", "label", text)
		}
	}
	res.Element = col.Children(func() {
		ui.Row(c).Gap(k.D(8)).Children(func() {
			side(c, k, k.Msg("you"), mine, t.Mine)
			side(c, k, partner, theirs, t.Theirs)
		})
		label, v := k.Msg("confirm"), input.Primary
		if t.Mine {
			label, v = k.Msg("cancel"), input.Normal
		}
		if input.Button(c, label, input.ButtonOptions{Variant: v, Disabled: o.Disabled}).Clicked {
			t.Mine, res.Changed = !t.Mine, true
			slog.Debug("8bit: trade confirm", "label", text, "mine", t.Mine)
		}
		if t.Mine && t.Theirs {
			pixel.Text(c, k.Msg("traded"), pixel.TextOptions{Slot: pixel.Good, Scale: 2})
		}
	})
	// The partner may confirm between frames; compare with the last.
	res.Traded = t.Mine && t.Theirs && !st.both
	st.both = t.Mine && t.Theirs
	if res.Traded {
		slog.Debug("8bit: trade done", "label", text)
	}
	return res
}

// side is one offer: who, eight slots and a stamp.
func side(c *ui.Context, k pixel.Kit, who string, items []data.Item, ok bool) {
	said := make([]string, len(items))
	for i, it := range items {
		said[i] = it.String()
	}
	layout.Panel(c, layout.PanelOptions{Style: pixel.Sunken, Fill: pixel.Well}, func() {
		ui.Column(c).Gap(k.D(3)).Children(func() {
			pixel.Text(c, who, pixel.TextOptions{})
			w := 4*data.SlotSize + 3*2
			ui.Box(c).Size(k.D(w), k.D(2*data.SlotSize+2)).Role(ui.RoleList).Label(who).Description(strings.Join(said, ", ")).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				for i := range tradeSlots {
					var it data.Item
					if i < len(items) {
						it = items[i]
					}
					data.DrawSlot(g, i%4*(data.SlotSize+2), i/4*(data.SlotSize+2), it, data.ItemSlotOptions{})
				}
			})
			stamp, ink := k.Msg("waiting"), pixel.Faint
			if ok {
				stamp, ink = k.Msg("confirmed"), pixel.Good
			}
			pixel.Text(c, stamp, pixel.TextOptions{Slot: ink})
		})
	})
}
