package pixel

import "fmt"

var bayer = [4][4]int{{0, 8, 2, 10}, {12, 4, 14, 6}, {3, 11, 1, 9}, {15, 7, 13, 5}}

// Bayer is the 4×4 ordered threshold, 0–15, at x, y.
func Bayer(x, y int) int { return bayer[y&3][x&3] }

// Dither fills b below the level threshold, a elsewhere.
func (g Canvas) Dither(x, y, w, h int, a, b Slot, level int) {
	checkLevel(level)
	if level <= 8 {
		g.Rect(x, y, w, h, a)
		g.cells(x, y, w, h, b, func(t int) bool { return t < level })
		return
	}
	g.Rect(x, y, w, h, b)
	g.cells(x, y, w, h, a, func(t int) bool { return t >= level })
}

// Veil dithers s over what is already there.
func (g Canvas) Veil(x, y, w, h int, s Slot, level int) {
	checkLevel(level)
	g.cells(x, y, w, h, s, func(t int) bool { return t < level })
}

// Checker is a half Veil: every other dot.
func (g Canvas) Checker(x, y, w, h int, s Slot) { g.Veil(x, y, w, h, s, 8) }

func (g Canvas) cells(x, y, w, h int, s Slot, on func(t int) bool) {
	for j := range h {
		row := bayer[(g.gy+y+j)&3]
		g.bits(x, y+j, w, 1, s, func(i int) bool { return on(row[(g.gx+x+i)&3]) })
	}
}

func checkLevel(level int) {
	if level < 0 || level > 16 {
		panic(fmt.Sprintf("pixel: dither level %d outside 0–16", level))
	}
}
