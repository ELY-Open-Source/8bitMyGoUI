package text

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// segments lists a..g as rects in a 7×11 cell.
var segments = [7][4]int{{1, 0, 5, 1}, {6, 1, 1, 4}, {6, 6, 1, 4}, {1, 10, 5, 1}, {0, 6, 1, 4}, {0, 1, 1, 4}, {1, 5, 5, 1}}

// lit holds each rune's segments, bit 0 for a.
var lit = map[rune]uint8{
	'0': 0x3F, '1': 0x06, '2': 0x5B, '3': 0x4F, '4': 0x66, '5': 0x6D,
	'6': 0x7D, '7': 0x07, '8': 0x7F, '9': 0x6F, '-': 0x40, ' ': 0,
}

type SevenSegmentOptions struct {
	Scale int
	// Slot lights segments; zero means Bad.
	Slot pixel.Slot
}

// SevenSegment shows digits, minus and colon with ghost segments.
func SevenSegment(c *ui.Context, s string, o SevenSegmentOptions) *ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	ink := slotOr(o.Slot, pixel.Bad)
	w := 0
	for _, r := range s {
		w += cellW(r) + 2
	}
	return block(c, k, s, max(0, w-2)*sc, 11*sc, func(g pixel.Canvas) {
		x := 0
		for _, r := range s {
			if r == ':' {
				g.Rect(x+sc, 3*sc, sc, sc, ink)
				g.Rect(x+sc, 7*sc, sc, sc, ink)
			} else {
				for i, sg := range segments {
					if lit[r]&(1<<i) != 0 {
						g.Rect(x+sg[0]*sc, sg[1]*sc, sg[2]*sc, sg[3]*sc, ink)
					} else {
						g.Checker(x+sg[0]*sc, sg[1]*sc, sg[2]*sc, sg[3]*sc, pixel.Faint)
					}
				}
			}
			x += (cellW(r) + 2) * sc
		}
	})
}

// cellW is a rune's cell width; other runes panic.
func cellW(r rune) int {
	if r == ':' {
		return 3
	}
	if _, ok := lit[r]; !ok {
		panic(fmt.Sprintf("text: SevenSegment cannot show %q", r))
	}
	return 7
}
