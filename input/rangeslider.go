package input

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type RangeSliderOptions struct {
	Min, Max float64
	Width    int
	Disabled bool
}

type RangeSliderResult = SliderResult

// RangeSlider sets *low and *high with two uncrossing knobs.
func RangeSlider(c *ui.Context, text string, low, high *float64, o RangeSliderOptions) RangeSliderResult {
	k := pixel.K(c)
	lo, hi := bounds(o.Min, o.Max)
	w := pixel.Or("RangeSlider Width", o.Width, 96)
	was := [2]float64{*low, *high}
	e := ui.Box(c).Size(k.D(w), k.D(12)).Role(ui.RoleGroup).Label(text)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		a, b := knobAt(*low, lo, hi, g.W-6), knobAt(*high, lo, hi, g.W-6)
		g.Frame(0, 4, g.W, 4, pixel.Sunken, pixel.Well)
		g.Rect(a+3, 5, b-a, 2, pixel.Accent)
	}).Children(func() {
		knob(c, k, e, low, rail{lo: lo, hi: *high, label: text + " low", disabled: o.Disabled}, lo, hi, w)
		knob(c, k, e, high, rail{lo: *low, hi: hi, label: text + " high", disabled: o.Disabled}, lo, hi, w)
	})
	if *low > *high {
		panic(fmt.Sprintf("input: RangeSlider low %v above high %v", *low, *high))
	}
	changed := was != [2]float64{*low, *high}
	if changed {
		slog.Debug("8bit: range slider", "label", text, "low", *low, "high", *high)
	}
	return SliderResult{Element: e, Changed: changed}
}

// knob is one draggable end, fenced by r.lo and r.hi.
func knob(c *ui.Context, k pixel.Kit, track *ui.Element, v *float64, r rail, lo, hi float64, w int) {
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleSlider).Label(r.label).Description(fmtValue(*v)).Disabled(r.disabled)
	})
	set := func(x float64) { *v = max(r.lo, min(r.hi, x)) }
	if !r.disabled {
		if _, _, held := base.Dragged(); held {
			px, _, _ := base.PointerPosition()
			t := track.Bounds()
			set(lo + float64((base.Bounds().X+px-t.X)/max(1, t.W))*(hi-lo))
		}
		step := (hi - lo) / 20
		if base.Shortcut(0, ui.KeyLeft) || base.Shortcut(0, ui.KeyDown) {
			set(*v - step)
		}
		if base.Shortcut(0, ui.KeyRight) || base.Shortcut(0, ui.KeyUp) {
			set(*v + step)
		}
	}
	shell.Absolute().Left(k.D(knobAt(*v, lo, hi, w-6))).Top(k.D(1)).Size(k.D(6), k.D(10)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, 6, 10, pixel.Raised, pixel.Face)
		finish(g, st)
	})
}
