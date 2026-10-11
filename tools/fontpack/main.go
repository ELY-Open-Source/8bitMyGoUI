// Command fontpack packs a BDF bitmap font for package pixel.
package main

import (
	"bufio"
	"cmp"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type glyph struct {
	r          rune
	adv        int
	w, h, x, y int
	rows       []byte
}

func main() {
	in := flag.String("in", "", "BDF file")
	out := flag.String("out", "pixel/font12.bin", "packed font")
	flag.Parse()
	f, err := os.Open(*in)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	asc, desc, gs, err := parse(f)
	if err != nil {
		log.Fatalf("fontpack: %s: %v", *in, err)
	}
	data, err := pack(asc, desc, gs)
	if err != nil {
		log.Fatalf("fontpack: %v", err)
	}
	if err := os.WriteFile(*out, data, 0o644); err != nil {
		log.Fatal(err)
	}
	log.Printf("fontpack: %d glyphs, %d bytes, %s", len(gs), len(data), *out)
}

func parse(r io.Reader) (asc, desc int, gs []glyph, err error) {
	asc, desc = -1, -1
	sc := bufio.NewScanner(r)
	var g *glyph
	inBitmap := false
	line := 0
	for sc.Scan() {
		line++
		f := strings.Fields(sc.Text())
		if len(f) == 0 {
			continue
		}
		bad := func(e error) (int, int, []glyph, error) { return 0, 0, nil, fmt.Errorf("line %d: %w", line, e) }
		switch {
		case inBitmap && f[0] == "ENDCHAR":
			if len(g.rows) != g.h*((g.w+7)/8) {
				return bad(fmt.Errorf("glyph %U has %d bitmap bytes, want %d", g.r, len(g.rows), g.h*((g.w+7)/8)))
			}
			if g.r >= 0 {
				gs = append(gs, *g)
			}
			g, inBitmap = nil, false
		case inBitmap:
			b, err := hex.DecodeString(f[0])
			if err != nil || len(b) != (g.w+7)/8 {
				return bad(fmt.Errorf("glyph %U row %q", g.r, f[0]))
			}
			g.rows = append(g.rows, b...)
		case f[0] == "FONT_ASCENT":
			asc, err = number(f)
		case f[0] == "FONT_DESCENT":
			desc, err = number(f)
		case f[0] == "STARTCHAR":
			g = &glyph{r: -1}
		case g == nil:
		case f[0] == "ENCODING":
			var e int
			e, err = number(f)
			g.r = rune(e)
		case f[0] == "DWIDTH":
			g.adv, err = number(f)
		case f[0] == "BBX":
			if len(f) != 5 {
				return bad(errors.New("BBX needs four numbers"))
			}
			v := make([]int, 4)
			for i := range v {
				if v[i], err = strconv.Atoi(f[i+1]); err != nil {
					return bad(err)
				}
			}
			g.w, g.h, g.x, g.y = v[0], v[1], v[2], v[3]
			if g.w < 0 || g.h < 0 {
				return bad(fmt.Errorf("glyph %U has a negative box", g.r))
			}
		case f[0] == "BITMAP":
			inBitmap = true
		}
		if err != nil {
			return bad(err)
		}
	}
	if err := sc.Err(); err != nil {
		return 0, 0, nil, err
	}
	if g != nil {
		return 0, 0, nil, fmt.Errorf("glyph %U has no ENDCHAR", g.r)
	}
	if asc < 0 || desc < 0 {
		return 0, 0, nil, errors.New("FONT_ASCENT or FONT_DESCENT missing")
	}
	return asc, desc, gs, nil
}

func number(f []string) (int, error) {
	if len(f) < 2 {
		return 0, fmt.Errorf("%s needs a number", f[0])
	}
	return strconv.Atoi(f[1])
}

// Layout: magic, ascent, descent, count, index, glyphs.
func pack(asc, desc int, gs []glyph) ([]byte, error) {
	slices.SortFunc(gs, func(a, b glyph) int { return cmp.Compare(a.r, b.r) })
	if asc > 255 || desc > 255 {
		return nil, fmt.Errorf("ascent %d or descent %d over 255", asc, desc)
	}
	out := []byte{'8', 'B', 'F', '1', byte(asc), byte(desc)}
	out = binary.LittleEndian.AppendUint32(out, uint32(len(gs)))
	index := len(out)
	out = append(out, make([]byte, 8*len(gs))...)
	for i, g := range gs {
		if i > 0 && gs[i-1].r == g.r {
			return nil, fmt.Errorf("glyph %U twice", g.r)
		}
		if g.adv < 0 || g.adv > 255 || g.w < 0 || g.w > 255 || g.h < 0 || g.h > 255 || g.x < -128 || g.x > 127 || g.y < -128 || g.y > 127 {
			return nil, fmt.Errorf("glyph %U metrics out of range", g.r)
		}
		binary.LittleEndian.PutUint32(out[index+8*i:], uint32(g.r))
		binary.LittleEndian.PutUint32(out[index+8*i+4:], uint32(len(out)))
		out = append(out, byte(g.adv), byte(g.w), byte(g.h), byte(int8(g.x)), byte(int8(g.y)))
		out = append(out, g.rows...)
	}
	return out, nil
}
