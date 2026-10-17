package input

import (
	"log/slog"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type EditableTextOptions struct {
	Placeholder string
	// Width is in dots; zero means 120.
	Width    int
	Disabled bool
}

type EditableTextResult struct {
	Element ui.Element
	Changed bool
}

type editableKey struct{}

type editable struct {
	on, refocus bool
	draft       string
}

// EditableText edits on double click or Enter; Escape reverts.
func EditableText(c *ui.Context, text string, value *string, o EditableTextOptions) EditableTextResult {
	k := pixel.K(c)
	w := pixel.Or("EditableText Width", o.Width, 120)
	box := ui.Box(c)
	s := ui.Local(box, editableKey{}, func() editable { return editable{} })
	res := EditableTextResult{Element: box}
	if o.Disabled {
		s.on = false
	}
	box.Children(func() {
		if s.on {
			ed, e := field(c, &s.draft, pixel.EditorOptions{Label: text, Placeholder: o.Placeholder}, w, nil, nil, nil)
			pixel.AutoFocus(c, ed.Input)
			// Keys reach the wrapper; the editor's state gates them.
			on := !ed.Input.IsDisabled()
			enter, escape := on && pixel.Shortcut(e, 0, ui.KeyEnter), on && pixel.Shortcut(e, 0, ui.KeyEscape)
			if enter {
				res.Changed = *value != s.draft
				*value = s.draft
			}
			if enter || escape {
				s.on, s.refocus = false, true
			}
			return
		}
		shell, base, st := pixel.Shell(c, func() ui.Element {
			return ui.Box(c).Focusable().Role(ui.RoleButton).Label(text).Description(*value).Disabled(o.Disabled)
		})
		if s.refocus {
			base.Focus()
			s.refocus = false
		}
		enter, double := pixel.Shortcut(base, 0, ui.KeyEnter), base.DoubleClicked()
		if enter || double {
			s.on, s.draft = true, *value
		}
		shell.Size(k.D(w), k.D(pixel.LineHeight+6)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			if st.Hover {
				g.Rect(0, 0, g.W, g.H, pixel.Lit)
			}
			shown, ink := *value, pixel.Ink
			if shown == "" {
				shown, ink = o.Placeholder, pixel.Faint
			}
			g.Text(4, 3, shown, ink, 1)
			g.Veil(0, g.H-1, g.W, 1, pixel.Faint, 8)
			g.Finish(st)
		})
	})
	if res.Changed {
		slog.Debug("8bit: text edited", "label", text, "runes", utf8.RuneCountInString(*value))
	}
	return res
}
