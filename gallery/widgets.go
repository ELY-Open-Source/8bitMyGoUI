package gallery

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// press is a clickable box; w ≤ 0 stretches.
func press(c *ui.Context, w, h int, label string, draw func(g pixel.Canvas, st pixel.State)) bool {
	k := pixel.K(c)
	shell, base, st := pixel.Shell(c, func() *ui.Element { return ui.ButtonBase(c).Label(label) })
	if w > 0 {
		shell.Width(k.D(w))
	}
	shell.Height(k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		draw(g, st)
		if st.Focus {
			g.FocusMark(0, 0, g.W, g.H)
		}
	})
	return base.Clicked()
}

// button is a raised pixel button with a centered label.
func button(c *ui.Context, label string) bool {
	return press(c, pixel.Measure(label)+12, 16, label, func(g pixel.Canvas, st pixel.State) {
		y := 0
		if st.Press {
			y = 1
		}
		g.Frame(0, y, g.W, g.H-1, pixel.Raised, pixel.Face)
		g.Text((g.W-pixel.Measure(label))/2, y, label, pixel.Ink, 1)
	})
}

// panel is a framed box with padding for children.
func panel(c *ui.Context, style pixel.FrameStyle, fill pixel.Slot, pad int) *ui.Element {
	k := pixel.K(c)
	return ui.Column(c).Padding(k.D(pad)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, style, fill)
	})
}

// shadowText draws text with a one-dot shadow.
func shadowText(g pixel.Canvas, x, y int, s string, ink pixel.Slot, scale int) {
	g.Text(x+scale, y+scale, s, pixel.Shade, scale)
	g.Text(x, y, s, ink, scale)
}

// label is non-wrapping pixel text.
func label(c *ui.Context, s string, ink pixel.Slot) *ui.Element {
	return pixel.Text(c, s, pixel.TextOptions{Slot: ink})
}
