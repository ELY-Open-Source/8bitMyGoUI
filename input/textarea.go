package input

import (
	"fmt"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TextAreaOptions struct {
	Placeholder string
	// MaxLength caps runes and shows a count.
	MaxLength int
	// Width is zero for 160 dots, Lines zero for 4.
	Width, Lines int
	Disabled     bool
	// EnterSends submits on Enter; Shift+Enter breaks lines.
	EnterSends bool
}

type areaKey struct{}

type areaState struct {
	scroll ui.ScrollState
	line   int
	text   string
	// held is last frame's offset, kept while disabled.
	held float32
}

type TextAreaResult = TextInputResult

// TextArea wraps lines in its own scroll, following the caret.
func TextArea(c *ui.Context, text string, value *string, o TextAreaOptions) TextAreaResult {
	k := pixel.K(c)
	lines := pixel.Or("TextArea Lines", o.Lines, 4)
	col := ui.Column(c).Width(k.D(pixel.Or("TextArea Width", o.Width, 160))).Gap(k.D(2))
	st := ui.Local(col, areaKey{}, func() areaState { return areaState{line: -1} })
	// Disabled or inert, the view holds still.
	if o.Disabled || pixel.Current(c).Inert {
		st.scroll.Y = st.held
	}
	var ed pixel.Edit
	col.Children(func() {
		view := ui.Scroll(c).TrackScroll(&st.scroll).Height(k.D(lines*pixel.LineHeight+6)).Padding(k.D(3), k.D(4))
		view.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
		}).Children(func() {
			ed = pixel.Editor(c, value, pixel.EditorOptions{Multiline: true, Label: text, Placeholder: o.Placeholder, MaxLength: o.MaxLength, Disabled: o.Disabled, EnterSends: o.EnterSends})
		})
		// Keep the caret's line in view after it moves.
		if ed.Line != st.line || *value != st.text {
			top, bottom := k.D(ed.Line*pixel.LineHeight), k.D((ed.Line+1)*pixel.LineHeight)
			if y := max(min(st.scroll.Y, top), bottom-k.D(lines*pixel.LineHeight)); y != st.scroll.Y {
				st.scroll.Y = y
				view.TrackScroll(&st.scroll)
			}
			st.line, st.text = ed.Line, *value
		}
		st.held = st.scroll.Y
		if o.MaxLength > 0 {
			count := fmt.Sprintf("%d/%d", utf8.RuneCountInString(*value), o.MaxLength)
			ui.Row(c).Justify(ui.End).Children(func() { pixel.Text(c, count, pixel.TextOptions{Slot: pixel.Faint}) })
		}
	})
	return result(col, ed, text)
}
