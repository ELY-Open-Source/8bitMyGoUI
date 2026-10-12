package input

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type KnobOptions struct {
	// Min and Max are zero for 0 to 1.
	Min, Max float64
	// Step snaps the value; zero means none.
	Step     float64
	Disabled bool
}

type KnobResult = SliderResult

type knobKey struct{}

// Knob turns through 16 angles by vertical drag or arrows.
func Knob(c *ui.Context, text string, value *float64, o KnobOptions) KnobResult {
	k := pixel.K(c)
	lo, hi := bounds(o.Min, o.Max)
	r := rail{lo: lo, hi: hi, step: o.Step}
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleSlider).Label(text).Description(fmtValue(*value)).Disabled(o.Disabled)
	})
	old, set := *value, r.snap(value)
	// raw is the unrounded value a drag carries.
	raw := ui.Local(base, knobKey{}, func() float64 { return *value })
	if !o.Disabled {
		if _, dy, held := base.Dragged(); held {
			*raw = max(lo, min(hi, *raw-float64(dy/k.Dot())*(hi-lo)/64))
			set(*raw)
		} else {
			*raw = *value
		}
		r.keys(base, value, set)
	}
	shell.Size(k.D(24), k.D(24)).Draw(func(p *ui.Painter, rc ui.Rect) {
		g := k.Canvas(p, rc)
		g.Disc(12, 12, 11, pixel.Edge)
		g.Disc(12, 12, 10, pixel.Shade)
		g.Disc(11, 11, 9, pixel.Face)
		a := (-135 + 18*math.Round((*value-lo)/(hi-lo)*15)) * math.Pi / 180
		x, y := 12+int(math.Round(7*math.Sin(a))), 12-int(math.Round(7*math.Cos(a)))
		g.Line(12, 12, x, y, pixel.Accent)
		g.Dot(x, y, pixel.Ink)
		finish(g, st)
	})
	return report(shell, *value != old, "knob", text, *value)
}
