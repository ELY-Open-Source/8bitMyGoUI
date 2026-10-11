package main

import (
	"encoding/binary"
	"strings"
	"testing"
)

const sample = `STARTFONT 2.1
FONT_ASCENT 13
FONT_DESCENT 3
CHARS 4
STARTCHAR B
ENCODING 66
DWIDTH 7 0
BBX 5 2 1 0
BITMAP
F8
88
ENDCHAR
STARTCHAR unencoded
ENCODING -1
DWIDTH 3 0
BBX 1 1 0 0
BITMAP
80
ENDCHAR
STARTCHAR A
ENCODING 65
DWIDTH 8 0
BBX 9 1 0 -2
BITMAP
FF80
ENDCHAR
STARTCHAR space
ENCODING 32
DWIDTH 6 0
BBX 0 0 0 0
BITMAP
ENDCHAR
ENDFONT
`

func TestParseAndPack(t *testing.T) {
	asc, desc, gs, err := parse(strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}
	if asc != 13 || desc != 3 || len(gs) != 3 {
		t.Fatalf("ascent %d descent %d glyphs %d", asc, desc, len(gs))
	}
	b, err := pack(asc, desc, gs)
	if err != nil {
		t.Fatal(err)
	}
	if string(b[:4]) != "8BF1" || b[4] != 13 || b[5] != 3 || binary.LittleEndian.Uint32(b[6:]) != 3 {
		t.Fatalf("header % X", b[:10])
	}
	runes := []rune{' ', 'A', 'B'}
	records := [][]byte{
		{6, 0, 0, 0, 0},
		{8, 9, 1, 0, 0xFE, 0xFF, 0x80},
		{7, 5, 2, 1, 0, 0xF8, 0x88},
	}
	for i, r := range runes {
		e := b[10+8*i:]
		if got := rune(binary.LittleEndian.Uint32(e)); got != r {
			t.Fatalf("index %d rune %U, want %U", i, got, r)
		}
		off := int(binary.LittleEndian.Uint32(e[4:]))
		if got := b[off : off+len(records[i])]; string(got) != string(records[i]) {
			t.Fatalf("%U record % X, want % X", r, got, records[i])
		}
	}
}

func TestParseErrors(t *testing.T) {
	cases := map[string]string{
		"no ascent":  strings.Replace(sample, "FONT_ASCENT 13\n", "", 1),
		"short rows": strings.Replace(sample, "F8\n88\n", "F8\n", 1),
		"bad hex":    strings.Replace(sample, "F8\n", "ZZ\n", 1),
		"wide row":   strings.Replace(sample, "F8\n", "F800\n", 1),
		"bad bbx":    strings.Replace(sample, "BBX 5 2 1 0", "BBX 5 2", 1),
		"no ENDCHAR": strings.TrimSuffix(sample, "ENDCHAR\nENDFONT\n"),
		"negative":   strings.Replace(sample, "BBX 0 0 0 0", "BBX -1 -1 0 0", 1),
	}
	for name, src := range cases {
		if _, _, _, err := parse(strings.NewReader(src)); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
}

func TestPackRejectsDuplicates(t *testing.T) {
	gs := []glyph{{r: 'A', adv: 1}, {r: 'A', adv: 2}}
	if _, err := pack(13, 3, gs); err == nil {
		t.Fatal("duplicate rune packed")
	}
	if _, err := pack(13, 3, []glyph{{r: 'A', adv: 300}}); err == nil {
		t.Fatal("advance 300 packed")
	}
}
