package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type TextInputOptions struct {
	Placeholder string
	MaxLength   int
	// Width is in dots; zero means 120.
	Width    int
	Disabled bool
}

type TextInputResult struct {
	// Element is laid out; Input takes the focus.
	Element, Input              ui.Element
	Changed, Submitted, Focused bool
}

// TextInput is a one-line pixel field; Enter reports Submitted.
func TextInput(c *ui.Context, label string, value *string, o TextInputOptions) TextInputResult {
	ed, e := field(c, value, pixel.EditorOptions{Label: label, Placeholder: o.Placeholder, MaxLength: o.MaxLength, Disabled: o.Disabled}, o.Width, nil, nil, nil)
	return result(e, ed, label)
}

// field frames an editor in a well with extras.
func field(c *ui.Context, value *string, eo pixel.EditorOptions, width int, bad func() bool, icon *pixel.Sprite, after func()) (pixel.Edit, ui.Element) {
	k := pixel.K(c)
	var ed pixel.Edit
	e := ui.Row(c).Width(k.D(pixel.Or("Width", width, 120))).Height(k.D(pixel.LineHeight+6)).Padding(0, k.D(4)).Gap(k.D(3)).AlignItems(ui.Center)
	e.Children(func() {
		if icon != nil {
			ui.Box(c).Size(k.D(icon.W), k.D(icon.H)).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Sprite(0, 0, icon, 0, pixel.Faint, 1)
			})
		}
		ed = pixel.Editor(c, value, eo)
		ed.Element.Grow(1)
		if after != nil {
			after()
		}
	})
	focus := ed.Focused
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		edge := pixel.Edge
		if bad != nil && bad() {
			edge = pixel.Bad
		}
		well(k.Canvas(p, r), edge, focus)
	})
	return ed, e
}

// well is a field's sunken box; focus dots its foot.
func well(g pixel.Canvas, edge pixel.Slot, focus bool) {
	g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
	g.HLine(1, 0, g.W-2, edge)
	g.HLine(1, g.H-1, g.W-2, edge)
	g.VLine(0, 1, g.H-2, edge)
	g.VLine(g.W-1, 1, g.H-2, edge)
	if focus {
		g.Checker(0, g.H-2, g.W, 1, pixel.Accent)
	}
}

// result wraps an editor's events and logs submits.
func result(e ui.Element, ed pixel.Edit, label string) TextInputResult {
	if ed.Submitted {
		slog.Debug("8bit: field submitted", "label", label)
	}
	return TextInputResult{Element: e, Input: ed.Input, Changed: ed.Changed, Submitted: ed.Submitted, Focused: ed.Focused}
}
