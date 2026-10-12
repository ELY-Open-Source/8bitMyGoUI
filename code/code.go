// Package code shows, runs and debugs programs and their data.
package code

import (
	"go/scanner"
	"go/token"
	"strings"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var (
	builtins = set("any bool byte comparable complex64 complex128 error float32 float64 int int8 int16 int32 int64 rune string " +
		"uint uint8 uint16 uint32 uint64 uintptr append cap clear close complex copy delete imag len make max min new panic print println real recover")
	constants = set("true false nil iota")
)

func set(words string) map[string]bool {
	m := map[string]bool{}
	for _, w := range strings.Fields(words) {
		m[w] = true
	}
	return m
}

// tidy drops carriage returns and turns tabs into four spaces.
func tidy(src string) string {
	src = strings.ReplaceAll(src, "\r", "")
	return strings.TrimRight(strings.ReplaceAll(src, "\t", "    "), "\n")
}

// paint is the slot of every byte of Go-like source.
func paint(src string) []pixel.Slot {
	slots := make([]pixel.Slot, len(src))
	fs := token.NewFileSet()
	f := fs.AddFile("", -1, len(src))
	var s scanner.Scanner
	s.Init(f, []byte(src), nil, scanner.ScanComments)
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			return slots
		}
		n := len(lit)
		if n == 0 {
			n = len(tok.String())
		}
		var slot pixel.Slot
		switch {
		case tok == token.COMMENT:
			slot = pixel.Faint
		case tok == token.STRING || tok == token.CHAR:
			slot = pixel.Good
		case tok == token.INT || tok == token.FLOAT || tok == token.IMAG:
			slot = pixel.Warn
		case tok.IsKeyword():
			slot = pixel.Accent
		case tok == token.IDENT && builtins[lit]:
			slot = pixel.Info
		case tok == token.IDENT && constants[lit]:
			slot = pixel.Magic
		default:
			continue
		}
		off := f.Offset(pos)
		for i := off; i < min(off+n, len(src)); i++ {
			slots[i] = slot
		}
	}
}

// colored draws s in runs by slot; returns the end.
func colored(g pixel.Canvas, x, y int, s string, slots []pixel.Slot) int {
	for i := 0; i < len(s); {
		j := i
		for j < len(s) && slots[j] == slots[i] {
			j++
		}
		g.Text(x, y, s[i:j], slots[i], 1)
		x += pixel.Measure(s[i:j])
		i = j
	}
	return x
}

// dark is a screen's body, light ink and dim ink.
func dark(k pixel.Kit) (body, ink, dim pixel.Slot) {
	if k.Palette().Night {
		return pixel.Well, pixel.Ink, pixel.Faint
	}
	return pixel.Ink, pixel.Lit, pixel.Shade
}
