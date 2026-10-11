package pixel

import (
	"slices"
	"testing"
)

func TestEditModel(t *testing.T) {
	var s editState
	s.reset("hello")
	if s.caret != 5 || s.anchor != 5 {
		t.Fatalf("reset caret %d anchor %d", s.caret, s.anchor)
	}
	s.anchor = 1
	s.replace([]rune("EY"), 0)
	if string(s.text) != "hEY" || s.caret != 3 || !s.changed {
		t.Fatalf("replace gave %q caret %d", string(s.text), s.caret)
	}
	s.replace([]rune("1234"), 5)
	if string(s.text) != "hEY12" {
		t.Fatalf("maxLen kept %q", string(s.text))
	}
	s.erase(9, -3)
	if string(s.text) != "" || s.caret != 0 {
		t.Fatalf("erase clamps: %q caret %d", string(s.text), s.caret)
	}
	s.restore(&s.undo, &s.redo)
	s.restore(&s.undo, &s.redo)
	if string(s.text) != "hEY" || len(s.redo) != 2 {
		t.Fatalf("two undos gave %q, redo %d", string(s.text), len(s.redo))
	}
	s.restore(&s.redo, &s.undo)
	if string(s.text) != "hEY12" {
		t.Fatalf("redo gave %q", string(s.text))
	}
	for range 150 {
		s.replace([]rune("x"), 0)
	}
	if len(s.undo) != 100 || len(s.redo) != 0 {
		t.Fatalf("history %d undo %d redo, want 100 and 0", len(s.undo), len(s.redo))
	}
}

func TestClean(t *testing.T) {
	if got := string(clean("a\tb\r\nc\x07d", false)); got != "a b cd" {
		t.Fatalf("one line: %q", got)
	}
	if got := string(clean("a\r\nb\nc", true)); got != "a\nb\nc" {
		t.Fatalf("multiline: %q", got)
	}
}

func TestWords(t *testing.T) {
	rs := []rune("foo  bar, 像素x")
	if wordRight(rs, 0) != 3 || wordRight(rs, 3) != 8 || wordRight(rs, 8) != 9 || wordRight(rs, 10) != 12 {
		t.Fatalf("wordRight %d %d %d %d", wordRight(rs, 0), wordRight(rs, 3), wordRight(rs, 8), wordRight(rs, 10))
	}
	if wordLeft(rs, len(rs)) != 12 || wordLeft(rs, 8) != 5 || wordLeft(rs, 5) != 0 {
		t.Fatalf("wordLeft %d %d %d", wordLeft(rs, len(rs)), wordLeft(rs, 8), wordLeft(rs, 5))
	}
	if a, b := wordAt(rs, 6); a != 5 || b != 8 {
		t.Fatalf("wordAt 6 = %d..%d", a, b)
	}
	if a, b := wordAt(rs, 11); a != 10 || b != 12 {
		t.Fatalf("wordAt in CJK = %d..%d", a, b)
	}
	if a, b := wordAt(nil, 0); a != 0 || b != 0 {
		t.Fatal("wordAt on empty text")
	}
}

func TestEditLayoutAndCaret(t *testing.T) {
	rs := []rune("ab cd\nef")
	lines := editLayout(rs, Measure("ab "), true)
	want := []editLine{{0, 2}, {3, 5}, {6, 8}}
	if !slices.Equal(lines, want) {
		t.Fatalf("layout %v, want %v", lines, want)
	}
	if x, li := caretAt(lines, rs, 3); x != 0 || li != 1 {
		t.Fatalf("caret at wrap is %d on line %d", x, li)
	}
	if x, li := caretAt(lines, rs, 5); x != Measure("cd") || li != 1 {
		t.Fatalf("caret before newline is %d on line %d", x, li)
	}
	if lineEnd(lines, 0) != 2 || lineEnd(lines, 1) != 5 || lineEnd(lines, 2) != 8 {
		t.Fatalf("lineEnd %d %d %d", lineEnd(lines, 0), lineEnd(lines, 1), lineEnd(lines, 2))
	}
	if hitLine(lines, rs, 2, 0) != 6 || hitLine(lines, rs, 2, 999) != 8 || hitLine(lines, rs, 0, Measure("a")) != 1 {
		t.Fatal("hitLine misses rune boundaries")
	}
	gap := []rune("ab   cd")
	gl := editLayout(gap, Measure("ab "), true)
	if x, li := caretAt(gl, gap, 3); li != 0 || x != Measure("ab") {
		t.Fatalf("caret in the break gap at %d on line %d", x, li)
	}
	if one := editLayout(rs, 1, false); !slices.Equal(one, []editLine{{0, 8}}) {
		t.Fatalf("single line layout %v", one)
	}
}

func TestKeepVisible(t *testing.T) {
	var s editState
	s.keepVisible(100, 120, 40)
	if s.scroll != 61 {
		t.Fatalf("scroll %d, want 61", s.scroll)
	}
	s.keepVisible(10, 120, 40)
	if s.scroll != 10 {
		t.Fatalf("scroll back %d, want 10", s.scroll)
	}
	s.keepVisible(10, 30, 40)
	if s.scroll != 0 {
		t.Fatalf("short text scroll %d, want 0", s.scroll)
	}
}

func TestMaskedComposition(t *testing.T) {
	var s editState
	s.reset("ab")
	s.compose = []rune("ni")
	if got := string(s.view('*', true)); got != "****" {
		t.Fatalf("masked view with composition = %q", got)
	}
	if got := string(s.view(0, true)); got != "abni" {
		t.Fatalf("plain view with composition = %q", got)
	}
}

func TestMultilineLeavesRoomForCaret(t *testing.T) {
	rs := []rune("ab")
	if lines := editLayout(rs, Measure("ab"), true); len(lines) != 2 {
		t.Fatalf("text as wide as the field: %v, want a wrap before the caret column", lines)
	}
}
