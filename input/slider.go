package input

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SliderOptions struct {
	// Min and Max are zero for 0 to 1.
	Min, Max float64
	// Width is the track in dots; zero means 96.
	Width    int
	Disabled bool
}

type SliderResult struct {
	Element *ui.Element
	Changed bool
}

// Slider sets *value by drag, arrows, Home and End.
func Slider(c *ui.Context, text string, value *float64, o SliderOptions) SliderResult {
	k := pixel.K(c)
	lo, hi := bounds(o.Min, o.Max)
	w := orDef("Slider Width", o.Width, 96)
	shell, st, changed := slide(c, value, rail{lo: lo, hi: hi, label: text, disabled: o.Disabled})
	shell.Size(k.D(w), k.D(12)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		x := knobAt(*value, lo, hi, g.W-6)
		groove(g, x+3, 4)
		g.Frame(x, 1, 6, 10, pixel.Raised, pixel.Face)
		finish(g, st)
	})
	return report(shell, changed, "slider", text, *value)
}

// rail maps a slide's pointer and keys to a value.
type rail struct {
	lo, hi, step float64
	vertical     bool
	label        string
	disabled     bool
}

// slide is a focusable track; pointer and keys set *v.
func slide(c *ui.Context, v *float64, r rail) (*ui.Element, pixel.State, bool) {
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleSlider).Label(r.label).Description(fmtValue(*v)).Disabled(r.disabled)
	})
	old := *v
	set := r.snap(v)
	if !r.disabled {
		if _, _, held := base.Dragged(); held {
			px, py, _ := base.PointerPosition()
			b := base.Bounds()
			frac := float64(px / max(1, b.W))
			if r.vertical {
				frac = 1 - float64(py/max(1, b.H))
			}
			set(r.lo + frac*(r.hi-r.lo))
		}
		r.keys(base, v, set)
	}
	return shell, st, *v != old
}

// snap sets *v rounded to steps and kept in range.
func (r rail) snap(v *float64) func(float64) {
	if r.step < 0 {
		panic(fmt.Sprintf("input: Step %v is negative", r.step))
	}
	top := r.hi
	if r.step > 0 {
		top = r.lo + math.Floor((r.hi-r.lo)/r.step+1e-9)*r.step
	}
	return func(x float64) {
		if r.step > 0 {
			x = r.lo + math.Round((x-r.lo)/r.step)*r.step
		}
		*v = max(r.lo, min(top, x))
	}
}

// keys step *v by arrows; Home and End jump.
func (r rail) keys(base *ui.Element, v *float64, set func(float64)) {
	step := r.step
	if step == 0 {
		step = (r.hi - r.lo) / 20
	}
	for _, kd := range []struct {
		key ui.Key
		to  func() float64
	}{
		{ui.KeyLeft, func() float64 { return *v - step }}, {ui.KeyDown, func() float64 { return *v - step }},
		{ui.KeyRight, func() float64 { return *v + step }}, {ui.KeyUp, func() float64 { return *v + step }},
		{ui.KeyHome, func() float64 { return r.lo }}, {ui.KeyEnd, func() float64 { return r.hi }},
	} {
		if base.Shortcut(0, kd.key) {
			set(kd.to())
		}
	}
}

// bounds are lo and hi; zeros mean 0 to 1.
func bounds(lo, hi float64) (float64, float64) {
	if lo == 0 && hi == 0 {
		return 0, 1
	}
	if hi <= lo {
		panic(fmt.Sprintf("input: range %v to %v is empty", lo, hi))
	}
	return lo, hi
}

// knobAt is the knob's offset for v along span dots.
func knobAt(v, lo, hi float64, span int) int {
	return int(math.Round((v - lo) / (hi - lo) * float64(span)))
}

// groove is the Well track, filled with Accent to x.
func groove(g pixel.Canvas, x, y int) {
	g.Frame(0, y, g.W, 4, pixel.Sunken, pixel.Well)
	g.Rect(1, y+1, x-1, 2, pixel.Accent)
}

// report logs a change and wraps the result.
func report(e *ui.Element, changed bool, kind, text string, v float64) SliderResult {
	if changed {
		slog.Debug("8bit: "+kind, "label", text, "value", v)
	}
	return SliderResult{Element: e, Changed: changed}
}
