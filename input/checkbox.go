package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var (
	check = pixel.MustSprite(`
		.....#
		....##
		#..##.
		####..
		.##...
	`)
	dash = pixel.MustSprite(`
		####
	`)
)

type CheckboxOptions struct {
	Disabled bool
}

type CheckboxResult struct {
	Element *ui.Element
	Changed bool
}

// Checkbox is a labeled 8×8 box with a hand-placed check.
func Checkbox(c *ui.Context, text string, checked *bool, o CheckboxOptions) CheckboxResult {
	k := pixel.K(c)
	shell, base, st := pixel.Shell(c, func() *ui.Element { return ui.CheckboxBase(c, checked).Label(text).Disabled(o.Disabled) })
	mark := check
	if !*checked {
		mark = nil
	}
	shell.Size(k.D(12+pixel.Measure(text)), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		box(k.Canvas(p, r), text, mark, st)
	})
	res := CheckboxResult{Element: shell, Changed: base.Changed()}
	if res.Changed {
		slog.Debug("8bit: checkbox", "label", text, "checked", *checked)
	}
	return res
}

// box draws the 8×8 box, its mark and the label.
func box(g pixel.Canvas, text string, mark *pixel.Sprite, st pixel.State) {
	y := (pixel.LineHeight - 8) / 2
	fill := pixel.Well
	if st.Hover {
		fill = pixel.Lit
	}
	g.Frame(0, y, 8, 8, pixel.Sunken, fill)
	if mark != nil {
		g.Sprite((8-mark.W)/2, y+(8-mark.H)/2, mark, 0, pixel.Accent, 1)
	}
	ink := pixel.Ink
	if st.Disabled {
		ink = pixel.Faint
	}
	g.Text(12, 0, text, ink, 1)
	finish(g, st)
}
