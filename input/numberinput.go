package input

import (
	"log/slog"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type NumberInputOptions struct {
	Min, Max float64
	// Step is zero for 1.
	Step     float64
	Width    int
	Disabled bool
}

type numberKey struct{}

type numberState struct {
	text    string
	focused bool
	// seen is the value last shown, to notice the caller's.
	seen float64
}

type NumberInputResult = TextInputResult

// NumberInput edits *value as text; bad text reddens the edge.
func NumberInput(c *ui.Context, label string, value *float64, o NumberInputOptions) NumberInputResult {
	lo, hi := bounds(o.Min, o.Max)
	step := o.Step
	if step <= 0 {
		step = 1
	}
	wrap := ui.Box(c)
	st := ui.Local(wrap, numberKey{}, func() numberState { return numberState{text: fmtValue(*value), seen: *value} })
	var res TextInputResult
	wrap.Children(func() {
		k := pixel.K(c)
		set := func(v float64) {
			if v = max(lo, min(hi, v)); v != *value {
				*value, res.Changed = v, true
			}
			st.text, st.seen = fmtValue(*value), *value
		}
		bad := func() bool {
			v, err := strconv.ParseFloat(st.text, 64)
			return err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v < lo || v > hi
		}
		if *value != st.seen || !st.focused && !bad() && st.text != fmtValue(*value) {
			st.text, st.seen = fmtValue(*value), *value
		}
		var ed pixel.Edit
		ed, res.Element = field(c, &st.text, pixel.EditorOptions{Label: label, Disabled: o.Disabled}, o.Width, bad, nil, func() {
			ui.Column(c).Children(func() {
				if arrowBox(c, k, tinyUp, o.Disabled) {
					set(*value + step)
				}
				if arrowBox(c, k, tinyDown, o.Disabled) {
					set(*value - step)
				}
			})
		})
		if v, err := strconv.ParseFloat(st.text, 64); ed.Changed && !bad() && err == nil && v != *value {
			*value, st.seen, res.Changed = v, v, true
		}
		st.focused = ed.Focused
		if ed.Input.Shortcut(0, ui.KeyUp) {
			set(*value + step)
		}
		if ed.Input.Shortcut(0, ui.KeyDown) {
			set(*value - step)
		}
		res.Input, res.Submitted, res.Focused = ed.Input, ed.Submitted, ed.Focused
	})
	if res.Changed {
		slog.Debug("8bit: number input", "label", label, "value", *value)
	}
	res.Element = wrap
	return res
}

// arrowBox is a half-height step button the pointer presses.
func arrowBox(c *ui.Context, k pixel.Kit, sp *pixel.Sprite, disabled bool) bool {
	b := ui.Box(c).Size(k.D(10), k.D(pixel.LineHeight/2)).Role(ui.RoleNone)
	down := b.Pressed()
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		dy := 0
		if down {
			dy = 1
		}
		g.Sprite((g.W-sp.W)/2, (g.H-sp.H)/2+dy, sp, 0, pixel.Faint, 1)
	})
	return b.Clicked() && !disabled
}
