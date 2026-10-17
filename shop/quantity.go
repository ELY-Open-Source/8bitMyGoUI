package shop

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type QuantityOptions struct {
	// Min is the least; Max the most, zero meaning 99.
	Min, Max int
	Disabled bool
}

type QuantityResult struct {
	Element ui.Element
	Changed bool
}

// Quantity steps − n + within its limits.
func Quantity(c *ui.Context, text string, value *int, o QuantityOptions) QuantityResult {
	k := pixel.K(c)
	lo, hi := o.Min, pixel.Or("Quantity Max", o.Max, 99)
	if lo < 0 || hi < lo || *value < lo || *value > hi {
		panic(fmt.Sprintf("shop: Quantity %d outside %d to %d", *value, lo, hi))
	}
	shell, base, st := pixel.Shell(c, func() ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleStepper).Label(text).Disabled(o.Disabled)
	})
	old, nw := *value, pixel.Measure(strconv.Itoa(hi))+6
	if !base.IsDisabled() {
		step := 0
		// Enter and Space are taken, so clicks come from pointers.
		pixel.Shortcut(base, 0, ui.KeyEnter)
		pixel.Shortcut(base, 0, ui.KeySpace)
		if pixel.Shortcut(base, 0, ui.KeyLeft) || pixel.Shortcut(base, 0, ui.KeyDown) {
			step--
		}
		if pixel.Shortcut(base, 0, ui.KeyRight) || pixel.Shortcut(base, 0, ui.KeyUp) {
			step++
		}
		if base.Clicked() {
			if x, _, _ := base.PointerPosition(); int(x/k.Dot()) < 14 {
				step--
			} else if int(x/k.Dot()) >= 14+nw {
				step++
			}
		}
		// One at a time toward the limit, so nothing overflows.
		for ; step > 0 && *value < hi; step-- {
			*value++
		}
		for ; step < 0 && *value > lo; step++ {
			*value--
		}
		if pixel.Shortcut(base, 0, ui.KeyHome) {
			*value = lo
		}
		if pixel.Shortcut(base, 0, ui.KeyEnd) {
			*value = hi
		}
	}
	n := *value
	base.Description(strconv.Itoa(n))
	shell.Size(k.D(28+nw), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		// A side at its limit fades; it goes no further.
		for _, b := range []struct {
			x  int
			sp *pixel.Sprite
			on bool
		}{{0, icon.Minus, n > lo}, {14 + nw, icon.Plus, n < hi}} {
			ink := pixel.Faint
			if b.on {
				ink = pixel.Ink
			}
			g.Frame(b.x, 1, 14, 14, pixel.Raised, pixel.Face)
			g.Sprite(b.x+1, 2, b.sp, 0, ink, 1)
		}
		s := strconv.Itoa(n)
		g.Text(14+(nw-pixel.Measure(s))/2, 0, s, pixel.Ink, 1)
		g.Finish(st)
	})
	res := QuantityResult{Element: shell, Changed: *value != old}
	if res.Changed {
		slog.Debug("8bit: quantity", "label", text, "value", *value)
	}
	return res
}
