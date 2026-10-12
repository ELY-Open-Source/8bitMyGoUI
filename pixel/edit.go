package pixel

import (
	"slices"
	"strings"
	"time"
	"unicode"
)

type snapshot struct {
	text          []rune
	caret, anchor int
}

// editState is one Editor's text, caret and history.
type editState struct {
	text          []rune
	caret, anchor int
	compose       []rune
	composeCaret  int
	undo, redo    []snapshot
	changed       bool
	submitted     bool
	// kept is the value the last frame ended with.
	kept string
	// clip is a copy waiting for an enabled build.
	clip           string
	clipped        bool
	scroll         int
	dragging       bool
	pressAt        time.Time
	pressX, pressY float32
	clicks         int
}

// reset takes v as the text and the value kept.
func (s *editState) reset(v string) {
	*s = editState{text: []rune(v), kept: v}
	s.caret = len(s.text)
	s.anchor = s.caret
}

func (s *editState) span() (int, int) { return min(s.caret, s.anchor), max(s.caret, s.anchor) }

func (s *editState) save() {
	s.undo = append(s.undo, snapshot{slices.Clone(s.text), s.caret, s.anchor})
	if len(s.undo) > 100 {
		s.undo = s.undo[1:]
	}
	s.redo = nil
}

// replace puts r over the selection, within maxLen runes.
func (s *editState) replace(r []rune, maxLen int) {
	a, b := s.span()
	if maxLen > 0 {
		r = r[:min(len(r), max(0, maxLen-len(s.text)+b-a))]
	}
	if a == b && len(r) == 0 {
		return
	}
	s.save()
	s.text = slices.Concat(s.text[:a], r, s.text[b:])
	s.caret = a + len(r)
	s.anchor = s.caret
	s.changed = true
}

func (s *editState) erase(a, b int) {
	a, b = max(0, min(a, b)), min(len(s.text), max(a, b))
	if a == b {
		return
	}
	s.save()
	s.text = slices.Delete(slices.Clone(s.text), a, b)
	s.caret, s.anchor = a, a
	s.changed = true
}

// eraseOr erases the selection, or else a..b.
func (s *editState) eraseOr(a, b int) {
	if x, y := s.span(); x != y {
		a, b = x, y
	}
	s.erase(a, b)
}

func (s *editState) move(to int, extend bool) {
	s.caret = max(0, min(to, len(s.text)))
	if !extend {
		s.anchor = s.caret
	}
}

// restore pops a snapshot from one history onto the other.
func (s *editState) restore(from, to *[]snapshot) {
	if len(*from) == 0 {
		return
	}
	*to = append(*to, snapshot{slices.Clone(s.text), s.caret, s.anchor})
	last := (*from)[len(*from)-1]
	*from = (*from)[:len(*from)-1]
	s.text, s.caret, s.anchor = last.text, last.caret, last.anchor
	s.changed = true
}

func (s *editState) selected() string {
	a, b := s.span()
	return string(s.text[a:b])
}

// clean drops control runes; one line turns breaks into spaces.
func clean(t string, multiline bool) []rune {
	var out []rune
	for _, r := range strings.ReplaceAll(t, "\r\n", "\n") {
		switch {
		case r == '\n' && multiline:
			out = append(out, r)
		case r == '\n' || r == '\t' || r == '\r':
			out = append(out, ' ')
		case !unicode.IsControl(r):
			out = append(out, r)
		}
	}
	return out
}

func class(r rune) int {
	switch {
	case unicode.IsSpace(r):
		return 0
	case wide(r):
		return 2
	case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_':
		return 1
	}
	return 3
}

func wordRight(rs []rune, i int) int {
	for i < len(rs) && class(rs[i]) == 0 {
		i++
	}
	if i < len(rs) {
		c := class(rs[i])
		for i < len(rs) && class(rs[i]) == c {
			i++
		}
	}
	return i
}

func wordLeft(rs []rune, i int) int {
	for i > 0 && class(rs[i-1]) == 0 {
		i--
	}
	if i > 0 {
		c := class(rs[i-1])
		for i > 0 && class(rs[i-1]) == c {
			i--
		}
	}
	return i
}

// wordAt is the run of one class around rune i.
func wordAt(rs []rune, i int) (int, int) {
	if len(rs) == 0 {
		return 0, 0
	}
	i = min(i, len(rs)-1)
	c := class(rs[i])
	a, b := i, i+1
	for a > 0 && class(rs[a-1]) == c {
		a--
	}
	for b < len(rs) && class(rs[b]) == c {
		b++
	}
	return a, b
}
