package text

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LinkOptions struct {
	Disabled bool
}

type LinkResult struct {
	Element ui.Element
	Clicked bool
}

// Link is underlined text opening url on click or Enter.
func Link(c *ui.Context, label, url string, o LinkOptions) LinkResult {
	k := pixel.K(c)
	shell, base, st := pixel.Shell(c, func() ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleLink).Label(label).Description(url).Disabled(o.Disabled).Cursor(ui.CursorPointer)
	})
	w := pixel.Measure(label)
	shell.Size(k.D(w+2), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		ink := pixel.Info
		switch {
		case st.Disabled:
			ink = pixel.Faint
		case st.Press:
			g.Rect(0, 0, g.W, g.H, pixel.Shade)
		case st.Hover:
			g.Rect(0, 0, g.W, g.H, pixel.Lit)
		}
		g.Text(1, 0, label, ink, 1)
		if !st.Disabled {
			g.HLine(1, pixel.LineHeight-2, w, ink)
		}
		if st.Focus {
			g.FocusMark(0, 0, g.W, g.H)
		}
	})
	res := LinkResult{Element: shell, Clicked: base.Clicked()}
	if res.Clicked {
		slog.Debug("8bit: link open", "url", url)
		c.OpenURL(url)
	}
	return res
}
