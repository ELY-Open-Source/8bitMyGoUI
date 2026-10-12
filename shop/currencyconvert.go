package shop

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CurrencyConvertOptions struct {
	// Rate is gold per gem; zero means 100.
	Rate     int
	Disabled bool
}

type CurrencyConvertResult struct {
	Element *ui.Element
	// Converted is the gems bought this frame.
	Converted int
}

type convertKey struct{}

// CurrencyConvert trades gold for gems at a rate.
func CurrencyConvert(c *ui.Context, text string, gold, gems *int, o CurrencyConvertOptions) CurrencyConvertResult {
	k := pixel.K(c)
	rate := pixel.Or("CurrencyConvert Rate", o.Rate, 100)
	amount("CurrencyConvert gold", *gold)
	amount("CurrencyConvert gems", *gems)
	// No more gems than gold buys or an int holds.
	most := min(*gold/rate, math.MaxInt-*gems)
	col := ui.Column(c).Gap(k.D(4)).Role(ui.RoleGroup).Label(text)
	n := ui.Local(col, convertKey{}, func() int { return 1 })
	*n = max(1, min(*n, most))
	var res CurrencyConvertResult
	res.Element = col.Children(func() {
		pixel.Text(c, fmt.Sprintf(k.Msg("rate"), rate), pixel.TextOptions{Slot: pixel.Faint})
		if most < 1 {
			pixel.Text(c, k.Msg("not_enough"), pixel.TextOptions{Slot: pixel.Bad})
		} else {
			ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
				Quantity(c, k.Msg("gems"), n, QuantityOptions{Min: 1, Max: most, Disabled: o.Disabled})
				money(c, k, Gold, *n*rate, pixel.Ink)
			})
		}
		if input.Button(c, k.Msg("convert"), input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled || most < 1}).Clicked {
			*gold, *gems, res.Converted = *gold-*n*rate, *gems+*n, *n
			slog.Debug("8bit: convert", "label", text, "gems", *n)
		}
	})
	return res
}
