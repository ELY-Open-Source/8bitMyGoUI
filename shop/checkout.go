package shop

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/navigation"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CheckoutOptions struct {
	Disabled bool
}

type CheckoutResult struct {
	Element ui.Element
	// Done reports the order placed this frame.
	Changed, Done bool
}

// Checkout walks CART, PAY and CONFIRM, then says THANK YOU!
func Checkout(c *ui.Context, text string, step *int, o CheckoutOptions, content func()) CheckoutResult {
	k := pixel.K(c)
	labels := strings.Fields(k.Msg("checkout_steps"))
	last := len(labels) - 1
	if *step < 0 || *step > last {
		panic(fmt.Sprintf("shop: Checkout step %d of %d", *step, len(labels)))
	}
	old := *step
	var res CheckoutResult
	res.Element = ui.Column(c).Gap(k.D(6)).Role(ui.RoleGroup).Label(text).Children(func() {
		navigation.Steps(c, labels, *step, navigation.StepsOptions{})
		if *step == last {
			pixel.Text(c, k.Msg("thank_you"), pixel.TextOptions{Slot: pixel.Good, Scale: 2})
			return
		}
		// A wrapper per step keeps NEXT in place; off, veiled.
		in := pixel.Current(c)
		in.Inert = in.Inert || o.Disabled
		ui.Column(c).Key(*step).Children(func() { pixel.Scope(c, in, content) }).DrawOver(func(p *ui.Painter, r ui.Rect) {
			if o.Disabled {
				g := k.Canvas(p, r)
				g.Checker(0, 0, g.W, g.H, pixel.Face)
			}
		})
		ui.Row(c).Gap(k.D(6)).Children(func() {
			if input.Button(c, k.Msg("back"), input.ButtonOptions{Disabled: o.Disabled || *step == 0}).Clicked {
				*step--
			}
			// The last step before THANK YOU! places the order.
			next := k.Msg("next_step")
			if *step == last-1 {
				next = k.Msg("place_order")
			}
			if input.Button(c, next, input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled}).Clicked {
				*step++
			}
		})
	})
	res.Changed, res.Done = *step != old, *step == last && old != last
	if res.Changed {
		slog.Debug("8bit: checkout", "label", text, "step", *step)
	}
	return res
}
