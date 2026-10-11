package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type RadioOptions struct {
	Disabled bool
}

type RadioResult struct {
	Element *ui.Element
	Changed bool
}

// Radio sets *selected to value when chosen.
func Radio[T comparable](c *ui.Context, text string, value T, selected *T, o RadioOptions) RadioResult {
	k := pixel.K(c)
	was := *selected == value
	on := was
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.CheckboxBase(c, &on).Role(ui.RoleRadio).Label(text).Disabled(o.Disabled)
	})
	res := RadioResult{Element: shell}
	if base.Changed() && !was {
		*selected, res.Changed = value, true
		slog.Debug("8bit: radio", "label", text)
	}
	chosen := *selected == value
	shell.Size(k.D(12+pixel.Measure(text)), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		dot(k.Canvas(p, r), text, chosen, st)
	})
	return res
}

// dot draws the radio, filled when on, and the label.
func dot(g pixel.Canvas, text string, on bool, st pixel.State) {
	y := (pixel.LineHeight - 8) / 2
	fill := pixel.Well
	if st.Hover {
		fill = pixel.Lit
	}
	g.Rect(2, y, 4, 8, pixel.Edge)
	g.Rect(0, y+2, 8, 4, pixel.Edge)
	g.Rect(1, y+1, 6, 6, pixel.Edge)
	g.Rect(2, y+1, 4, 6, fill)
	g.Rect(1, y+2, 6, 4, fill)
	if on {
		g.Rect(3, y+2, 2, 4, pixel.Accent)
		g.Rect(2, y+3, 4, 2, pixel.Accent)
	}
	ink := pixel.Ink
	if st.Disabled {
		ink = pixel.Faint
	}
	g.Text(12, 0, text, ink, 1)
	finish(g, st)
}
