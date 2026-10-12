package input

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type HoldButtonOptions struct {
	// Duration is the hold needed; zero means one second.
	Duration time.Duration
	// Danger fills in Bad instead of Accent.
	Danger   bool
	Disabled bool
}

type HoldButtonResult struct {
	Element   *ui.Element
	Confirmed bool
}

type holdButtonKey struct{}

// holding is a hold's start and its last key repeat.
type holding struct {
	since, key time.Time
	done       bool
}

// repeatGap outlasts a key's delay before it repeats.
const repeatGap = 600 * time.Millisecond

// HoldButton confirms only when held; its fill grows stepwise.
func HoldButton(c *ui.Context, text string, o HoldButtonOptions) HoldButtonResult {
	k := pixel.K(c)
	d := o.Duration
	if d < 0 {
		panic(fmt.Sprintf("input: HoldButton Duration %v", d))
	}
	if d == 0 {
		d = time.Second
	}
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.ButtonBase(c).Label(text).Description(k.Msg("hold")).Disabled(o.Disabled)
	})
	h := ui.Local(base, holdButtonKey{}, func() holding { return holding{} })
	if base.IsDisabled() {
		*h = holding{}
	}
	now := k.Now()
	// Held keys repeat; each repeat moves the hold on.
	enter, space := pixel.Shortcut(base, 0, ui.KeyEnter), pixel.Shortcut(base, 0, ui.KeySpace)
	if enter || space {
		if k.Elapsed(&h.key) > repeatGap {
			h.since = now
		}
		h.key = now
	}
	keyed := k.Elapsed(&h.key) <= repeatGap
	var held time.Duration
	switch {
	case st.Press:
		if h.since.IsZero() {
			h.since = now
		}
		held = k.Elapsed(&h.since)
	case keyed:
		held = max(0, h.key.Sub(h.since))
	default:
		*h = holding{}
	}
	res := HoldButtonResult{Element: shell}
	if held >= d && !h.done {
		h.done, res.Confirmed = true, true
		slog.Debug("8bit: hold confirmed", "label", text)
	}
	if st.Press || keyed {
		k.Tick(c, 15)
	}
	steps := min(8, int(held*8/d))
	fill := pixel.Accent
	if o.Danger {
		fill = pixel.Bad
	}
	shell.Size(k.D(pixel.Measure(text)+16), k.D(pixel.LineHeight+4)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		ink, dy := face(g, Normal, st, false)
		w := (g.W - 2) * steps / 8
		g.Rect(1, 1+dy, w, g.H-2-dy, fill)
		g.Text(8, 2+dy, text, ink, 1)
		g.Clip(1, 0, w, g.H, func() { g.Text(8, 2+dy, text, pixel.AccentInk, 1) })
		g.Finish(st)
	})
	return res
}
