package pixel

import (
	_ "embed"
	"encoding/binary"
	"strings"
	"unicode"
)

//go:embed font12.bin
var fontData []byte

// Text line height in dots at ×1.
const LineHeight = 16

type font struct {
	ascent, descent, n int
	index, data        []byte
}

type glyph struct {
	adv, w, h, x, y, stride int
	rows                    []byte
}

var font12 = loadFont(fontData)

// tofu stands in for runes the font lacks.
var tofu = glyph{adv: 8, w: 6, h: 9, x: 1, stride: 1, rows: []byte{0xFC, 0x84, 0x84, 0x84, 0x84, 0x84, 0x84, 0x84, 0xFC}}

func loadFont(b []byte) *font {
	if len(b) < 10 || string(b[:4]) != "8BF1" {
		panic("pixel: font12.bin has a bad header")
	}
	n := int(binary.LittleEndian.Uint32(b[6:10]))
	if len(b) < 10+8*n {
		panic("pixel: font12.bin index truncated")
	}
	return &font{ascent: int(b[4]), descent: int(b[5]), n: n, index: b[10 : 10+8*n], data: b}
}

func (f *font) find(r rune) (glyph, bool) {
	lo, hi := 0, f.n
	for lo < hi {
		m := int(uint(lo+hi) >> 1)
		if rune(binary.LittleEndian.Uint32(f.index[8*m:])) < r {
			lo = m + 1
		} else {
			hi = m
		}
	}
	if lo == f.n || rune(binary.LittleEndian.Uint32(f.index[8*lo:])) != r {
		return glyph{}, false
	}
	off := int(binary.LittleEndian.Uint32(f.index[8*lo+4:]))
	d := f.data
	if off+5 > len(d) {
		panic("pixel: font12.bin glyph truncated")
	}
	g := glyph{adv: int(d[off]), w: int(d[off+1]), h: int(d[off+2]), x: int(int8(d[off+3])), y: int(int8(d[off+4]))}
	g.stride = (g.w + 7) / 8
	end := off + 5 + g.h*g.stride
	if end > len(d) {
		panic("pixel: font12.bin glyph truncated")
	}
	g.rows = d[off+5 : end]
	return g, true
}

func glyphOf(r rune) glyph {
	if g, ok := font12.find(r); ok {
		return g
	}
	return tofu
}

func (g glyph) on(x, y int) bool { return g.rows[y*g.stride+x/8]&(0x80>>(x%8)) != 0 }

// top is the bitmap's first row below the line top.
func (g glyph) top() int { return font12.ascent - g.y - g.h }

// Measure returns one line's width in dots at ×1.
func Measure(s string) int {
	w := 0
	for _, r := range s {
		w += glyphOf(r).adv
	}
	return w
}

// Lines wraps s to width dots; zero width keeps lines.
func Lines(s string, width int) []string {
	var out []string
	for _, p := range strings.Split(s, "\n") {
		out = append(out, wrap([]rune(p), width)...)
	}
	return out
}

func wrap(rs []rune, width int) []string {
	if width <= 0 || Measure(string(rs)) <= width {
		return []string{string(rs)}
	}
	var lines []string
	for start := 0; start < len(rs); {
		w, brk, i := 0, -1, start
		for ; i < len(rs); i++ {
			if i > start && breakable(rs, i) {
				brk = i
			}
			w += glyphOf(rs[i]).adv
			if w > width && i > start {
				break
			}
		}
		if i == len(rs) {
			lines = append(lines, string(rs[start:]))
			break
		}
		if brk <= start {
			brk = hardBreak(rs, start, i)
		}
		lines = append(lines, strings.TrimRight(string(rs[start:brk]), " "))
		start = skip(rs, brk)
	}
	return lines
}

// breakable reports whether a line may break before rs[i].
func breakable(rs []rune, i int) bool {
	if n := skip(rs, i); n < len(rs) && closer(rs[n]) || opener(rs[back(rs, 0, i)]) {
		return false
	}
	return rs[i] == ' ' || rs[i-1] == ' ' || wide(rs[i]) || wide(rs[i-1])
}

// hardBreak splits an overlong run, keeping punctuation attached.
func hardBreak(rs []rune, start, i int) int {
	q, hangs := i, false
	for q > start && (rs[q-1] == ' ' || opener(rs[q-1])) {
		hangs = hangs || opener(rs[q-1])
		q--
	}
	switch {
	case hangs && q > start:
		return q
	case hangs:
		for i < len(rs) && (rs[i] == ' ' || opener(rs[i])) {
			i++
		}
		i = min(i+1, len(rs))
	}
	for n := skip(rs, i); n < len(rs) && closer(rs[n]); n = skip(rs, i) {
		i = n + 1
	}
	return i
}

func skip(rs []rune, i int) int {
	for i < len(rs) && rs[i] == ' ' {
		i++
	}
	return i
}

// back returns the last non-space index before i.
func back(rs []rune, start, i int) int {
	p := i - 1
	for p > start && rs[p] == ' ' {
		p--
	}
	return p
}

func closer(r rune) bool {
	return unicode.In(r, unicode.Pe, unicode.Pf) || strings.ContainsRune(stops, r)
}

func opener(r rune) bool { return unicode.In(r, unicode.Ps, unicode.Pi) }

// Sentence ends that must not start a line.
const stops = ",.;:!?，。、；：？！…‥·・"

func wide(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana, unicode.Hangul) ||
		r >= 0x3000 && r <= 0x303F || r >= 0xFF00 && r <= 0xFFEF
}
