package input

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type StepperOptions struct {
	Min, Max int
	// Step is zero for 1.
	Step     int
	Disabled bool
}

type StepperResult struct {
	Element *ui.Element
	Changed bool
}

// Stepper is ◀ value ▶; held arrows repeat, keys step.
func Stepper(c *ui.Context, text string, value *int, o StepperOptions) StepperResult {
	if o.Max < o.Min {
		panic(fmt.Sprintf("input: Stepper range %d to %d is empty", o.Min, o.Max))
	}
	k := pixel.K(c)
	step := orDef("Stepper Step", o.Step, 1)
	old := *value
	set := func(v int) { *value = max(o.Min, min(o.Max, v)) }
	w := max(pixel.Measure(strconv.Itoa(o.Min)), pixel.Measure(strconv.Itoa(o.Max))) + 8
	e := ui.Row(c).Gap(k.D(2)).AlignItems(ui.Center).Children(func() {
		set(*value - step*held(c, k, tinyLeft, o.Disabled))
		shell, base, st := pixel.Shell(c, func() *ui.Element {
			return ui.Box(c).Focusable().Role(ui.RoleStepper).Label(text).Description(strconv.Itoa(*value)).Disabled(o.Disabled)
		})
		if !o.Disabled {
			if base.Shortcut(0, ui.KeyLeft) || base.Shortcut(0, ui.KeyDown) {
				set(*value - step)
			}
			if base.Shortcut(0, ui.KeyRight) || base.Shortcut(0, ui.KeyUp) {
				set(*value + step)
			}
		}
		shell.Size(k.D(w), k.D(pixel.LineHeight+4)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			shown := strconv.Itoa(*value)
			g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
			g.Text((g.W-pixel.Measure(shown))/2, 2, shown, pixel.Ink, 1)
			finish(g, st)
		})
		set(*value + step*held(c, k, tinyRight, o.Disabled))
		base.Description(strconv.Itoa(*value))
	})
	res := StepperResult{Element: e, Changed: *value != old}
	if res.Changed {
		slog.Debug("8bit: stepper", "label", text, "value", *value)
	}
	return res
}

type holdKey struct{}

type hold struct {
	since time.Time
	n     int
}

// held counts an arrow's steps: one, then ten a second.
func held(c *ui.Context, k pixel.Kit, sp *pixel.Sprite, disabled bool) int {
	b := ui.Box(c).Size(k.D(12), k.D(pixel.LineHeight+4)).Role(ui.RoleNone)
	down := b.Pressed() && !disabled
	h := ui.Local(b, holdKey{}, func() hold { return hold{} })
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		style, dy := pixel.Raised, 0
		if down {
			style, dy = pixel.Sunken, 1
		}
		g.Frame(0, 0, g.W, g.H, style, pixel.Face)
		g.Sprite((g.W-sp.W)/2, (g.H-sp.H)/2+dy, sp, 0, pixel.Ink, 1)
		if disabled {
			g.Checker(0, 0, g.W, g.H, pixel.Face)
		}
	})
	if !down {
		*h = hold{}
		return 0
	}
	k.Tick(c, 10)
	if h.since.IsZero() {
		*h = hold{since: k.Now(), n: 1}
		return 1
	}
	want := 1 + max(0, int((k.Now().Sub(h.since)-400*time.Millisecond)/(100*time.Millisecond)))
	n := want - h.n
	h.n = want
	return n
}
