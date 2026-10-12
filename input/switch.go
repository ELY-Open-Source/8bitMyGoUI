package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SwitchOptions struct {
	Disabled bool
}

type SwitchResult struct {
	Element *ui.Element
	Changed bool
}

// Switch is a lever on an OFF/ON plate.
func Switch(c *ui.Context, text string, on *bool, o SwitchOptions) SwitchResult {
	k := pixel.K(c)
	shell, base, st := pixel.Shell(c, func() *ui.Element { return ui.SwitchBase(c, on).Label(text).Disabled(o.Disabled) })
	shell.Size(k.D(32+pixel.Measure(text)), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 1, 27, 14, pixel.Raised, pixel.Shade)
		g.Micro(3, 2, "OFF", pixel.Faint, 1)
		g.Micro(17, 2, "ON", pixel.Faint, 1)
		x, lx, lit := 6, 3, "OFF"
		if *on {
			x, lx, lit = 19, 17, "ON"
		}
		g.Micro(lx, 2, lit, pixel.Ink, 1)
		g.Rect(x, 8, 3, 5, pixel.Edge)
		g.Rect(x-1, 8, 5, 3, pixel.Edge)
		g.Rect(x, 9, 3, 1, pixel.Lit)
		ink := pixel.Ink
		if st.Disabled {
			ink = pixel.Faint
		}
		g.Text(32, 0, text, ink, 1)
		g.Finish(st)
	})
	res := SwitchResult{Element: shell, Changed: base.Changed()}
	if res.Changed {
		slog.Debug("8bit: switch", "label", text, "on", *on)
	}
	return res
}
