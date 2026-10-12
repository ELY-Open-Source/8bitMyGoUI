package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ToggleOptions struct {
	Size     Size
	Icon     *pixel.Sprite
	Disabled bool
	// Label names it for assistive technology; zero uses text.
	Label string
}

type ToggleResult struct {
	Element *ui.Element
	Changed bool
}

// Toggle is a button that stays sunken while *on.
func Toggle(c *ui.Context, text string, on *bool, o ToggleOptions) ToggleResult {
	k := pixel.K(c)
	h, pad, top, sc := o.Size.metrics()
	// The base turns before it learns it is off.
	v := *on
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.ToggleBase(c, &v).Label(label(o.Label, text)).Disabled(o.Disabled)
	})
	changed := v != *on && !base.IsDisabled()
	if changed {
		*on = v
	}
	lead := 0
	if o.Icon != nil {
		lead = o.Icon.W
		if text != "" {
			lead += 3
		}
	}
	shell.Size(k.D(pixel.Measure(text)*sc+lead+2*pad), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		ink, dy := face(g, Normal, st, *on)
		if o.Icon != nil {
			g.Sprite(pad, (h-o.Icon.H)/2+dy, o.Icon, 0, ink, 1)
		}
		g.Text(pad+lead, top+dy, text, ink, sc)
		g.Finish(st)
	})
	res := ToggleResult{Element: shell, Changed: changed}
	if res.Changed {
		slog.Debug("8bit: toggle", "label", label(o.Label, text), "on", *on)
	}
	return res
}
