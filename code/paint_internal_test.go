package code

import (
	"strings"
	"testing"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestPaintColorsGoTokens(t *testing.T) {
	src := "package main\n// hi\nfunc f() int { return len(\"x\") + 42 + nil }\nvar s = `a\nb` /* c\nd */"
	slots := paint(src)
	at := func(word string, nth int) pixel.Slot {
		i := -1
		for range nth + 1 {
			i += 1 + strings.Index(src[i+1:], word)
		}
		return slots[i]
	}
	for _, tc := range []struct {
		word string
		nth  int
		want pixel.Slot
	}{
		{"package", 0, pixel.Accent}, {"main", 0, pixel.Ink}, {"// hi", 0, pixel.Faint}, {"func", 0, pixel.Accent},
		{"f()", 0, pixel.Ink}, {"int", 0, pixel.Info}, {"len", 0, pixel.Info}, {`"x"`, 0, pixel.Good},
		{"42", 0, pixel.Warn}, {"nil", 0, pixel.Magic}, {"+", 0, pixel.Ink}, {"b`", 0, pixel.Good}, {"d */", 0, pixel.Faint},
	} {
		if got := at(tc.word, tc.nth); got != tc.want {
			t.Errorf("%q is %v, want %v", tc.word, got, tc.want)
		}
	}
	if tidy("a\r\n\tb\n\n") != "a\n    b" {
		t.Fatalf("tidy = %q", tidy("a\r\n\tb\n\n"))
	}
}
