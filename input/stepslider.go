package input

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type StepSliderOptions struct {
	Min, Max float64
	// Step is zero for a tenth of the range.
	Step     float64
	Width    int
	Disabled bool
}

type StepSliderResult = SliderResult

// StepSlider snaps to steps, with ticks and its value.
func StepSlider(c *ui.Context, text string, value *float64, o StepSliderOptions) StepSliderResult {
	k := pixel.K(c)
	lo, hi := bounds(o.Min, o.Max)
	step := o.Step
	if step <= 0 {
		step = (hi - lo) / 10
	}
	w := pixel.Or("StepSlider Width", o.Width, 96)
	var res SliderResult
	row := ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		shell, st, changed := slide(c, value, rail{lo: lo, hi: hi, step: step, label: text, disabled: o.Disabled})
		shell.Size(k.D(w), k.D(16)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			x := knobAt(*value, lo, hi, w-6)
			groove(g, x+3, 4)
			for v := lo; v <= hi+step/2; v += step {
				g.Dot(knobAt(v, lo, hi, w-6)+3, 13, pixel.Faint)
			}
			g.Frame(x, 1, 6, 10, pixel.Raised, pixel.Face)
			g.Finish(st)
		})
		res = report(shell, changed, "step slider", text, *value)
		pixel.Text(c, fmtValue(*value), pixel.TextOptions{})
	})
	res.Element = row
	return res
}
