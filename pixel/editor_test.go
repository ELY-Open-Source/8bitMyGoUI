package pixel_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type editHarness struct {
	t                *testing.T
	tt               *ui.Tester
	value            string
	o                pixel.EditorOptions
	last             pixel.Edit
	changes, submits int
	at               time.Time
}

func newEdit(t *testing.T, v string, width float32, o pixel.EditorOptions) *editHarness {
	t.Helper()
	h := &editHarness{t: t, value: v, o: o}
	if h.o.Label == "" {
		h.o.Label = "field"
	}
	h.tt = ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SummerDay, Dot: 2, Now: func() time.Time {
			if h.at.IsZero() {
				return time.Now()
			}
			return h.at
		}})
		ui.Column(c).Width(width).Padding(4).Children(func() {
			h.last = pixel.Editor(c, &h.value, h.o)
			if h.last.Changed {
				h.changes++
			}
			if h.last.Submitted {
				h.submits++
			}
		})
	}, 300, 200)
	return h
}

func (h *editHarness) focus() {
	h.t.Helper()
	if err := h.tt.Click(h.o.Label); err != nil {
		h.t.Fatal(err)
	}
}

func (h *editHarness) want(v string) {
	h.t.Helper()
	if h.value != v {
		h.t.Fatalf("value %q, want %q", h.value, v)
	}
}

func TestEditorTypingAndKeys(t *testing.T) {
	h := newEdit(t, "", 200, pixel.EditorOptions{})
	h.focus()
	h.tt.Type("hello")
	h.want("hello")
	if h.changes == 0 {
		t.Fatal("typing reported no change")
	}
	h.tt.Key(0, ui.KeyLeft)
	h.tt.Key(0, ui.KeyLeft)
	h.tt.Type("X")
	h.want("helXlo")
	h.tt.Key(0, ui.KeyBackspace)
	h.tt.Key(0, ui.KeyDelete)
	h.want("helo")
	h.tt.Key(0, ui.KeyHome)
	h.tt.Type("<")
	h.tt.Key(0, ui.KeyEnd)
	h.tt.Type(">")
	h.want("<helo>")
	h.tt.Key(ui.Shift, ui.KeyLeft)
	h.tt.Key(ui.Shift, ui.KeyLeft)
	h.tt.Type("!")
	h.want("<hel!")
	h.tt.Key(0, ui.KeyEnter)
	if h.submits != 1 {
		t.Fatalf("Enter submitted %d times", h.submits)
	}
	h.want("<hel!")
}

func TestEditorWordsClipboardUndo(t *testing.T) {
	h := newEdit(t, "foo bar baz", 200, pixel.EditorOptions{})
	h.focus()
	h.tt.Key(0, ui.KeyEnd)
	h.tt.Key(ui.Alt, ui.KeyLeft)
	h.tt.Type("X")
	h.want("foo bar Xbaz")
	h.tt.Key(ui.Cmd, ui.KeyZ)
	h.want("foo bar baz")
	h.tt.Key(ui.Cmd|ui.Shift, ui.KeyZ)
	h.want("foo bar Xbaz")
	h.tt.Key(ui.Alt, ui.KeyBackspace)
	h.want("foo bar baz")
	h.tt.Key(ui.Cmd, ui.KeyA)
	h.tt.Key(ui.Cmd, ui.KeyC)
	if h.tt.Clipboard() != "foo bar baz" {
		t.Fatalf("clipboard %q", h.tt.Clipboard())
	}
	h.tt.Key(0, ui.KeyRight)
	h.tt.Key(ui.Cmd, ui.KeyV)
	h.want("foo bar bazfoo bar baz")
	h.tt.Command("selectAll")
	h.tt.Command("cut")
	h.want("")
	if h.tt.Clipboard() != "foo bar bazfoo bar baz" {
		t.Fatalf("cut put %q on the clipboard", h.tt.Clipboard())
	}
	h.tt.SetClipboard("a\tb\nc")
	h.tt.Command("paste")
	h.want("a b c")
}

func TestEditorMultiline(t *testing.T) {
	h := newEdit(t, "ab\ncd", 200, pixel.EditorOptions{Multiline: true})
	h.focus()
	h.tt.Key(0, ui.KeyUp)
	h.tt.Type("X")
	h.want("abX\ncd")
	if h.last.Line != 0 {
		t.Fatalf("caret line %d on the first line", h.last.Line)
	}
	h.tt.Key(0, ui.KeyEnter)
	h.want("abX\n\ncd")
	if h.last.Line != 1 {
		t.Fatalf("caret line %d after a newline", h.last.Line)
	}
	h.tt.Key(ui.Cmd, ui.KeyEnter)
	if h.submits != 1 {
		t.Fatalf("Cmd+Enter submitted %d times", h.submits)
	}
	if r, _ := h.tt.Find("field"); r.H != 3*16*2 {
		t.Fatalf("three lines are %v high, want 96", r.H)
	}
}

func TestEditorEnterSends(t *testing.T) {
	h := newEdit(t, "hi", 200, pixel.EditorOptions{Multiline: true, EnterSends: true})
	h.focus()
	h.tt.Key(ui.Shift, ui.KeyEnter)
	h.tt.Type("yo")
	h.want("hi\nyo")
	h.tt.Key(0, ui.KeyEnter)
	if h.submits != 1 {
		t.Fatalf("Enter submitted %d times", h.submits)
	}
	h.want("hi\nyo")
	mustPanic(t, "one line", func() { newEdit(t, "", 200, pixel.EditorOptions{EnterSends: true}).tt.Frame() })
}

func TestEditorLimitsAndMask(t *testing.T) {
	h := newEdit(t, "", 200, pixel.EditorOptions{MaxLength: 5})
	h.focus()
	h.tt.Type("abcdefg")
	h.want("abcde")
	m := newEdit(t, "secret", 200, pixel.EditorOptions{Mask: '★'})
	m.focus()
	m.tt.Key(ui.Cmd, ui.KeyA)
	m.tt.Key(ui.Cmd, ui.KeyC)
	m.tt.Key(ui.Cmd, ui.KeyX)
	if m.tt.Clipboard() != "" {
		t.Fatalf("masked text reached the clipboard: %q", m.tt.Clipboard())
	}
	m.want("secret")
}

func TestEditorComposition(t *testing.T) {
	h := newEdit(t, "ab", 200, pixel.EditorOptions{})
	h.focus()
	h.tt.Compose("ni", 2)
	h.want("ab")
	if r, ok := h.tt.TextCaret(); !ok || r.H != 32 {
		t.Fatalf("caret %v active %v", r, ok)
	}
	h.tt.Compose("", 0)
	h.tt.Type("你")
	h.want("ab你")
	h.tt.Key(ui.Cmd, ui.KeyA)
	h.tt.Compose("x", 1)
	h.want("")
}

func TestEditorPointer(t *testing.T) {
	h := newEdit(t, "hello world", 300, pixel.EditorOptions{})
	r, _ := h.tt.Find("field")
	x := r.X + float32(pixel.Measure("hello ")+4)*2
	for range 2 {
		h.tt.Press(x, r.Y+8)
		h.tt.Release(x, r.Y+8)
	}
	h.tt.Type("there")
	h.want("hello there")
	h.tt.Press(r.X+1, r.Y+8)
	h.tt.Move(r.X+r.W-1, r.Y+8)
	h.tt.Release(r.X+r.W-1, r.Y+8)
	h.tt.Type("Z")
	h.want("Z")
}

func TestEditorExternalChangeAndDisabled(t *testing.T) {
	h := newEdit(t, "start", 200, pixel.EditorOptions{})
	h.focus()
	h.value = "reset"
	h.tt.Frame()
	h.tt.Type("!")
	h.want("reset!")
	d := newEdit(t, "frozen", 200, pixel.EditorOptions{Disabled: true})
	d.focus()
	d.tt.Type("x")
	d.tt.Key(0, ui.KeyBackspace)
	d.want("frozen")
	if d.tt.Focused("field") {
		t.Fatal("disabled editor took the focus")
	}
}

func TestEditorKeepsCaretInView(t *testing.T) {
	h := newEdit(t, "", 48, pixel.EditorOptions{})
	h.focus()
	h.tt.Type("a long line that cannot fit in a narrow field")
	field, _ := h.tt.Find("field")
	caret, _ := h.tt.TextCaret()
	if caret.X < field.X || caret.X+caret.W > field.X+field.W {
		t.Fatalf("caret %v outside field %v", caret, field)
	}
}

func TestEditorCaretBlinks(t *testing.T) {
	h := newEdit(t, "pixel", 200, pixel.EditorOptions{})
	h.at = time.Unix(1000, 0)
	h.focus()
	h.tt.Key(0, ui.KeyEnd)
	r, _ := h.tt.Find("field")
	x, y := int(r.X/2)+pixel.Measure("pixel"), int(r.Y/2)
	caret := func() int { return dots{t: t, img: h.tt.Image(), pal: pixel.SummerDay}.count(x, y+1, 1, 14, pixel.Ink) }
	if caret() != 14 {
		t.Fatal("caret not drawn after the text")
	}
	h.at = time.Unix(1000, 500_000_000)
	h.tt.Frame()
	if caret() != 0 {
		t.Fatal("caret did not blink off")
	}
	h.tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if caret() != 14 {
		t.Fatal("caret hidden under Reduce Motion")
	}
}

func TestMultilineCaretInsideField(t *testing.T) {
	h := newEdit(t, "ab", float32(pixel.Measure("ab")*2+8), pixel.EditorOptions{Multiline: true})
	h.focus()
	h.tt.Key(ui.Ctrl, ui.KeyEnd)
	field, _ := h.tt.Find("field")
	caret, _ := h.tt.TextCaret()
	if caret.X+caret.W > field.X+field.W || caret.Y+caret.H > field.Y+field.H {
		t.Fatalf("caret %v outside field %v", caret, field)
	}
}

func TestEditorDrawsCaretAndSelection(t *testing.T) {
	h := newEdit(t, "pixel", 200, pixel.EditorOptions{})
	h.focus()
	d := dots{t: t, img: h.tt.Image(), pal: pixel.SummerDay}
	r, _ := h.tt.Find("field")
	x0, y0, w := int(r.X/2), int(r.Y/2), int(r.W/2)
	if d.count(x0, y0, w, 16, pixel.Ink) == 0 {
		t.Fatal("no text drawn")
	}
	if d.count(x0, y0, w, 16, pixel.Accent) != 0 {
		t.Fatal("selection drawn without a selection")
	}
	h.tt.Key(ui.Cmd, ui.KeyA)
	d = dots{t: t, img: h.tt.Image(), pal: pixel.SummerDay}
	if d.count(x0, y0, pixel.Measure("pixel"), 16, pixel.Accent) == 0 || d.count(x0, y0, w, 16, pixel.AccentInk) == 0 {
		t.Fatal("selection not drawn in Accent with AccentInk text")
	}
}

func TestEditorOwnsTypingKeys(t *testing.T) {
	value, shortcuts := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SummerDay, Dot: 2})
		pixel.Editor(c, &value, pixel.EditorOptions{Label: "field"})
		if pixel.Typing(c) {
			return
		}
		for _, key := range []ui.Key{ui.KeyP, ui.KeySlash, ui.KeyBracketRight, ui.KeySpace} {
			if c.Shortcut(0, key) {
				shortcuts++
			}
		}
	}, 300, 100)
	if err := tt.Click("field"); err != nil {
		t.Fatal(err)
	}
	tt.TypeKey(0, ui.KeyP, "p")
	tt.TypeKey(0, ui.KeySlash, "/")
	tt.TypeKey(0, ui.KeyBracketRight, "]")
	tt.TypeKey(0, ui.KeySpace, " ")
	if value != "p/] " || shortcuts != 0 {
		t.Fatalf("value %q, shortcuts fired %d times", value, shortcuts)
	}
	tt.ClickAt(150, 90)
	tt.Key(0, ui.KeyP)
	if shortcuts != 1 {
		t.Fatalf("after leaving the field, P fired %d shortcuts", shortcuts)
	}
}

func TestEditorSpaceDoesNotScroll(t *testing.T) {
	value, scroll := "", ui.ScrollState{}
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SummerDay, Dot: 2})
		ui.Scroll(c).Height(100).TrackScroll(&scroll).Children(func() {
			pixel.Editor(c, &value, pixel.EditorOptions{Label: "field"})
			ui.Box(c).Height(1000)
		})
	}, 300, 100)
	if err := tt.Click("field"); err != nil {
		t.Fatal(err)
	}
	tt.TypeKey(0, ui.KeySpace, " ")
	if value != " " || scroll.Y != 0 {
		t.Fatalf("value %q, scrolled to %v", value, scroll.Y)
	}
}

func TestEditorTakesAValueChangedBetweenFrames(t *testing.T) {
	h := newEdit(t, "help", 200, pixel.EditorOptions{})
	h.focus()
	// The caller clears it; a key lands before a frame.
	h.value = ""
	h.tt.Type("x")
	h.want("x")
}
