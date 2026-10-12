package pixel

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
)

type EditorOptions struct {
	Multiline bool
	MaxLength int
	// Mask hides each rune behind Mask and turns copy off.
	Mask        rune
	Placeholder string
	Label       string
	Disabled    bool
	// EnterSends submits a multi-line field on Enter; Shift+Enter breaks lines.
	EnterSends bool
}

type Edit struct {
	// Element is laid out; Input takes the focus.
	Element, Input              *ui.Element
	Changed, Submitted, Focused bool
	// Line is the caret's line, from 0.
	Line int
}

type editKey struct{}

type editLine struct{ start, end int }

type editEnv struct {
	multi  bool
	sends  bool
	maxLen int
	mask   rune
	width  int
	dot    float32
	read   func() string
	write  func(string)
}

// Editor edits *value in the pixel font.
func Editor(c *ui.Context, value *string, o EditorOptions) Edit {
	k := K(c)
	if o.MaxLength < 0 {
		panic(fmt.Sprintf("pixel: MaxLength %d", o.MaxLength))
	}
	if o.EnterSends && !o.Multiline {
		panic("pixel: EnterSends needs a multi-line field")
	}
	shell, base, _ := Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleTextField).Label(o.Label).Disabled(o.Disabled)
	})
	s := ui.Local(base, editKey{}, func() editState { return editState{} })
	if string(s.text) != *value {
		s.reset(*value)
	}
	if o.Mask == 0 {
		base.Description(*value)
	}
	width, check := WrapWidth(c, shell)
	env := editEnv{multi: o.Multiline, sends: o.EnterSends, maxLen: o.MaxLength, mask: o.Mask, width: width, dot: k.dot, read: c.ReadClipboard, write: c.WriteClipboard}
	shown := s.view(o.Mask, true)
	caret := s.caret + s.composeCaret
	lines := editLayout(shown, width, o.Multiline)
	cx, cl := caretAt(lines, shown, caret)
	if o.Multiline {
		shell.Height(k.D(len(lines) * LineHeight))
	} else {
		shell.Height(k.D(LineHeight))
		s.keepVisible(cx, Measure(string(shown)), width)
	}
	scroll, sa, sb := s.scroll, s.caret, s.anchor
	sa, sb = min(sa, sb), max(sa, sb)
	ca, cb := s.caret, s.caret+len(s.compose)
	composing := cb > ca
	base.TextCaret(ui.Rect{X: k.D(cx - scroll), Y: k.D(cl * LineHeight), W: k.dot, H: k.D(LineHeight)})
	base.HandleInput(func(ev ui.InputEvent) bool {
		// The caller may have changed the value since this frame.
		if string(s.text) != *value {
			s.reset(*value)
		}
		took := s.handle(ev, env)
		if s.changed {
			*value = string(s.text)
		}
		return took
	})
	res := Edit{Element: shell, Input: base, Changed: s.changed, Submitted: s.submitted, Focused: base.Focused(), Line: cl}
	if res.Focused {
		frameOf(c).typing = true
	}
	if s.submitted {
		slog.Debug("8bit: editor submit", "label", o.Label)
	}
	s.changed, s.submitted = false, false
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		check(r)
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() {
			if len(shown) == 0 {
				g.Text(0, 0, o.Placeholder, Faint, 1)
			}
			ink := Ink
			if o.Disabled {
				ink = Faint
			}
			for li, l := range lines {
				y, x0 := li*LineHeight, -scroll
				g.Text(x0, y, string(shown[l.start:l.end]), ink, 1)
				if a, b := max(sa, l.start), min(sb, l.end); a < b && !composing {
					xa := x0 + Measure(string(shown[l.start:a]))
					g.Rect(xa, y, Measure(string(shown[a:b])), LineHeight, Accent)
					g.Text(xa, y, string(shown[a:b]), AccentInk, 1)
				}
				if a, b := max(ca, l.start), min(cb, l.end); a < b {
					xa := x0 + Measure(string(shown[l.start:a]))
					for i := range Measure(string(shown[a:b])) {
						if i%2 == 0 {
							g.Dot(xa+i, y+LineHeight-1, Accent)
						}
					}
				}
			}
			if res.Focused && !o.Disabled && g.Step(2)%2 == 0 {
				g.VLine(cx-scroll, cl*LineHeight+1, LineHeight-2, Ink)
			}
		})
	})
	return res
}

// view is the shown runes: masked, with any composition.
func (s *editState) view(mask rune, composing bool) []rune {
	out := make([]rune, 0, len(s.text)+len(s.compose))
	for _, r := range s.text {
		if mask != 0 && r != '\n' {
			r = mask
		}
		out = append(out, r)
	}
	if composing {
		c := s.compose
		if mask != 0 {
			c = slices.Repeat([]rune{mask}, len(c))
		}
		out = slices.Insert(out, s.caret, c...)
	}
	return out
}

func (s *editState) keepVisible(cx, total, width int) {
	if width <= 0 {
		s.scroll = 0
		return
	}
	if cx-s.scroll > width-1 {
		s.scroll = cx - width + 1
	}
	s.scroll = max(0, min(s.scroll, cx, total-width+1))
}

// editLayout wraps a dot short, leaving the caret its column.
func editLayout(rs []rune, width int, multi bool) []editLine {
	if !multi {
		return []editLine{{0, len(rs)}}
	}
	width = max(0, width-1)
	var out []editLine
	start := 0
	for i := 0; i <= len(rs); i++ {
		if i == len(rs) || rs[i] == '\n' {
			for _, r := range wrapRanges(rs[start:i], width) {
				out = append(out, editLine{start + r[0], start + r[1]})
			}
			start = i + 1
		}
	}
	return out
}

// caretAt is the caret's x in dots and its line.
func caretAt(lines []editLine, rs []rune, idx int) (int, int) {
	li := 0
	for i, l := range lines {
		if l.start <= idx {
			li = i
		}
	}
	l := lines[li]
	return Measure(string(rs[l.start:min(idx, l.end)])), li
}

// hitLine is the rune boundary nearest x on line li.
func hitLine(lines []editLine, rs []rune, li, x int) int {
	l := lines[li]
	pos := 0
	for i := l.start; i < l.end; i++ {
		adv := glyphOf(rs[i]).adv
		if x < pos+adv/2 {
			return i
		}
		pos += adv
	}
	return lineEnd(lines, li)
}

// lineEnd keeps a wrapped line's caret on that line.
func lineEnd(lines []editLine, li int) int {
	l := lines[li]
	if li+1 < len(lines) && lines[li+1].start == l.end && l.end > l.start {
		return l.end - 1
	}
	return l.end
}

func (s *editState) hit(e editEnv, x, y float32) int {
	rs := s.view(e.mask, false)
	lines := editLayout(rs, e.width, e.multi)
	li := max(0, min(int(y/e.dot)/LineHeight, len(lines)-1))
	return hitLine(lines, rs, li, int(x/e.dot)+s.scroll)
}

func (s *editState) handle(ev ui.InputEvent, e editEnv) bool {
	switch ev.Kind {
	case ui.InputText:
		s.compose, s.composeCaret = nil, 0
		s.replace(clean(ev.Text, e.multi), e.maxLen)
		return true
	case ui.InputCompose:
		if ev.Text == "" {
			s.compose, s.composeCaret = nil, 0
			return true
		}
		if a, b := s.span(); len(s.compose) == 0 && a != b {
			s.erase(a, b)
		}
		s.compose = []rune(ev.Text)
		s.composeCaret = max(0, min(ev.Caret, len(s.compose)))
		return true
	case ui.InputCommand:
		return s.command(ev.Text, e)
	case ui.InputKeyDown:
		return s.key(ev, e)
	case ui.InputPointerDown:
		if ev.Button != 0 {
			return false
		}
		s.compose, s.composeCaret = nil, 0
		s.press(ev, e)
		return true
	case ui.InputPointerMove:
		if !s.dragging || ev.Button != 0 {
			return false
		}
		s.move(s.hit(e, ev.X, ev.Y), true)
		return true
	case ui.InputPointerUp:
		was := s.dragging
		s.dragging = false
		return was
	}
	return false
}

func (s *editState) press(ev ui.InputEvent, e editEnv) {
	now := time.Now()
	if now.Sub(s.pressAt) < 500*time.Millisecond && abs32(ev.X-s.pressX) < 5 && abs32(ev.Y-s.pressY) < 5 {
		s.clicks++
	} else {
		s.clicks = 1
	}
	s.pressAt, s.pressX, s.pressY = now, ev.X, ev.Y
	i := s.hit(e, ev.X, ev.Y)
	switch clicks := max(ev.Clicks, s.clicks); {
	case clicks == 2:
		s.anchor, s.caret = wordAt(s.text, i)
	case clicks >= 3:
		s.anchor, s.caret = 0, len(s.text)
	default:
		s.move(i, ev.Mods&ui.Shift != 0)
		s.dragging = true
	}
}

func (s *editState) command(name string, e editEnv) bool {
	switch name {
	case "copy":
		if a, b := s.span(); e.mask == 0 && a != b {
			e.write(s.selected())
		}
	case "cut":
		if a, b := s.span(); e.mask == 0 && a != b {
			e.write(s.selected())
			s.erase(a, b)
		}
	case "paste":
		s.replace(clean(e.read(), e.multi), e.maxLen)
	case "selectAll":
		s.anchor, s.caret = 0, len(s.text)
	case "undo":
		s.restore(&s.undo, &s.redo)
	case "redo":
		s.restore(&s.redo, &s.undo)
	case "delete":
		a, b := s.span()
		s.erase(a, b)
	default:
		return false
	}
	return true
}

func abs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func (s *editState) key(ev ui.InputEvent, e editEnv) bool {
	m, shift := ev.Mods&^ui.Shift, ev.Mods&ui.Shift != 0
	if m == ui.Cmd {
		switch ev.Key {
		case ui.KeyA:
			return s.command("selectAll", e)
		case ui.KeyC:
			return s.command("copy", e)
		case ui.KeyX:
			return s.command("cut", e)
		case ui.KeyV:
			return s.command("paste", e)
		case ui.KeyZ:
			if shift {
				return s.command("redo", e)
			}
			return s.command("undo", e)
		case ui.KeyY:
			return s.command("redo", e)
		case ui.KeyEnter:
			s.submitted = true
			return true
		}
	}
	word, line := m&(ui.Alt|ui.Ctrl) != 0 && m&ui.Super == 0, m == ui.Super
	rs := s.view(e.mask, false)
	lines := editLayout(rs, e.width, e.multi)
	x, li := caretAt(lines, rs, s.caret)
	a, b := s.span()
	switch ev.Key {
	case ui.KeyLeft:
		switch {
		case a != b && !shift && m == 0:
			s.move(a, false)
		case word:
			s.move(wordLeft(s.text, s.caret), shift)
		case line:
			s.move(lines[li].start, shift)
		default:
			s.move(s.caret-1, shift)
		}
	case ui.KeyRight:
		switch {
		case a != b && !shift && m == 0:
			s.move(b, false)
		case word:
			s.move(wordRight(s.text, s.caret), shift)
		case line:
			s.move(lineEnd(lines, li), shift)
		default:
			s.move(s.caret+1, shift)
		}
	case ui.KeyUp, ui.KeyDown:
		to, end := li+1, len(s.text)
		if ev.Key == ui.KeyUp {
			to, end = li-1, 0
		}
		if !e.multi || line || to < 0 || to >= len(lines) {
			s.move(end, shift)
		} else {
			s.move(hitLine(lines, rs, to, x), shift)
		}
	case ui.KeyHome:
		to := lines[li].start
		if m&ui.Ctrl != 0 {
			to = 0
		}
		s.move(to, shift)
	case ui.KeyEnd:
		to := lineEnd(lines, li)
		if m&ui.Ctrl != 0 {
			to = len(s.text)
		}
		s.move(to, shift)
	case ui.KeyBackspace:
		switch {
		case word:
			s.eraseOr(wordLeft(s.text, s.caret), s.caret)
		case line:
			s.eraseOr(lines[li].start, s.caret)
		default:
			s.eraseOr(s.caret-1, s.caret)
		}
	case ui.KeyDelete:
		if word {
			s.eraseOr(s.caret, wordRight(s.text, s.caret))
		} else {
			s.eraseOr(s.caret, s.caret+1)
		}
	case ui.KeyEnter:
		brk := e.multi && m == 0 && (shift || !e.sends)
		if !brk {
			s.submitted = true
			return true
		}
		s.replace([]rune{'\n'}, e.maxLen)
	default:
		return false
	}
	return true
}
