package pixel

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	MicroAdvance = 4
	MicroHeight  = 5
)

// Hand-placed 3×5 glyphs; rows split by |.
var micro = microFont(map[rune]string{
	' ': "...|...|...|...|...",
	'0': "###|#.#|#.#|#.#|###",
	'1': ".#.|##.|.#.|.#.|###",
	'2': "###|..#|###|#..|###",
	'3': "###|..#|.##|..#|###",
	'4': "#.#|#.#|###|..#|..#",
	'5': "###|#..|###|..#|###",
	'6': "###|#..|###|#.#|###",
	'7': "###|..#|.#.|.#.|.#.",
	'8': "###|#.#|###|#.#|###",
	'9': "###|#.#|###|..#|###",
	'A': ".#.|#.#|###|#.#|#.#",
	'B': "##.|#.#|##.|#.#|##.",
	'C': ".##|#..|#..|#..|.##",
	'D': "##.|#.#|#.#|#.#|##.",
	'E': "###|#..|##.|#..|###",
	'F': "###|#..|##.|#..|#..",
	'G': ".##|#..|#.#|#.#|.##",
	'H': "#.#|#.#|###|#.#|#.#",
	'I': "###|.#.|.#.|.#.|###",
	'J': "..#|..#|..#|#.#|.#.",
	'K': "#.#|#.#|##.|#.#|#.#",
	'L': "#..|#..|#..|#..|###",
	'M': "#.#|###|###|#.#|#.#",
	'N': "##.|#.#|#.#|#.#|#.#",
	'O': ".#.|#.#|#.#|#.#|.#.",
	'P': "##.|#.#|##.|#..|#..",
	'Q': ".#.|#.#|#.#|##.|.##",
	'R': "##.|#.#|##.|#.#|#.#",
	'S': ".##|#..|.#.|..#|##.",
	'T': "###|.#.|.#.|.#.|.#.",
	'U': "#.#|#.#|#.#|#.#|###",
	'V': "#.#|#.#|#.#|#.#|.#.",
	'W': "#.#|#.#|###|###|#.#",
	'X': "#.#|#.#|.#.|#.#|#.#",
	'Y': "#.#|#.#|.#.|.#.|.#.",
	'Z': "###|..#|.#.|#..|###",
	'+': "...|.#.|###|.#.|...",
	'-': "...|...|###|...|...",
	'.': "...|...|...|...|.#.",
	':': "...|.#.|...|.#.|...",
	'/': "..#|..#|.#.|#..|#..",
	'%': "#.#|..#|.#.|#..|#.#",
	'!': ".#.|.#.|.#.|...|.#.",
})

func microFont(art map[rune]string) map[rune][MicroHeight]uint8 {
	out := map[rune][MicroHeight]uint8{}
	for r, a := range art {
		rows := strings.Split(a, "|")
		if len(rows) != MicroHeight {
			panic(fmt.Sprintf("pixel: micro %q needs %d rows", r, MicroHeight))
		}
		var g [MicroHeight]uint8
		for y, row := range rows {
			if len(row) != 3 || strings.Trim(row, ".#") != "" {
				panic(fmt.Sprintf("pixel: micro %q row %q", r, row))
			}
			for x := range 3 {
				if row[x] == '#' {
					g[y] |= 4 >> x
				}
			}
		}
		out[r] = g
	}
	return out
}

func microGlyph(r rune) [MicroHeight]uint8 {
	g, ok := micro[unicode.ToUpper(r)]
	if !ok {
		panic(fmt.Sprintf("pixel: micro font lacks %q", r))
	}
	return g
}

// MicroWidth returns the width of s in dots.
func MicroWidth(s string) int {
	n := 0
	for _, r := range s {
		microGlyph(r)
		n++
	}
	if n == 0 {
		return 0
	}
	return n*MicroAdvance - 1
}
