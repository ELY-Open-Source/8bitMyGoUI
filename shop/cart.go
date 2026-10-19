package shop

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// CartLine is a ware in the cart and how many.
type CartLine struct {
	// Name is unique in the cart; it keys the row.
	Name  string
	Price int
	Qty   int
}

type CartOptions struct {
	Currency Currency
	// Width is in dots; zero means 260.
	Width    int
	Disabled bool
}

type CartResult struct {
	Element ui.Element
	Changed bool
	// Removed is the line whose ✕ was pressed, else -1.
	Removed int
}

// Cart lists lines with quantities and their total.
func Cart(c *ui.Context, text string, lines []CartLine, o CartOptions) CartResult {
	k := pixel.K(c)
	names, total := map[string]bool{}, 0
	for _, l := range lines {
		amount("Cart price", l.Price)
		// Checked against what is left, so totals cannot overflow.
		if l.Name == "" || names[l.Name] || l.Qty < 1 || l.Qty > 99 || l.Qty > 0 && l.Price > (math.MaxInt-total)/l.Qty {
			panic(fmt.Sprintf("shop: Cart line %q, twice %v, %d at %d", l.Name, names[l.Name], l.Qty, l.Price))
		}
		names[l.Name], total = true, total+l.Price*l.Qty
	}
	res := CartResult{Removed: -1}
	res.Element = ui.Column(c).Width(k.D(pixel.Or("Cart Width", o.Width, 260))).Gap(k.D(3)).Role(ui.RoleList).Label(text).Children(func() {
		for i := range lines {
			l := &lines[i]
			ui.Row(c).Key(l.Name).FillWidth().Gap(k.D(4)).AlignItems(ui.Center).Role(ui.RoleListItem).Children(func() {
				pixel.Text(c, l.Name, pixel.TextOptions{Wrap: true, MaxLines: 1}).Grow(1)
				// A raise may not carry the total past an int.
				most := 99
				if l.Price > 0 {
					most = min(99, (math.MaxInt-(total-l.Price*l.Qty))/l.Price)
				}
				res.Changed = Quantity(c, l.Name, &l.Qty, QuantityOptions{Min: 1, Max: most, Disabled: o.Disabled}).Changed || res.Changed
				money(c, k, o.Currency, l.Price*l.Qty, pixel.Ink)
				if input.IconButton(c, icon.Close, k.Msg("remove")+" "+l.Name, input.IconButtonOptions{Size: input.Small, Disabled: o.Disabled}).Clicked {
					res.Removed = i
				}
			})
		}
		if len(lines) == 0 {
			pixel.Text(c, k.Msg("empty"), pixel.TextOptions{Slot: pixel.Faint})
		}
		ui.Row(c).FillWidth().Gap(k.D(4)).Children(func() {
			pixel.Text(c, k.Msg("total"), pixel.TextOptions{}).Grow(1)
			money(c, k, o.Currency, total, pixel.Accent)
		})
	})
	if res.Removed >= 0 {
		slog.Debug("8bit: cart remove", "label", text, "line", res.Removed)
	}
	return res
}
