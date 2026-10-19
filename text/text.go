// Package text draws words in the pixel font.
package text

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// block is a w×h dot text box that fn draws.
func block(c *ui.Context, k pixel.Kit, s string, w, h int, fn func(g pixel.Canvas)) ui.Element {
	return ui.Box(c).Size(k.D(w), k.D(h)).Label(s).Role(ui.RoleText).Draw(func(p *ui.Painter, r ui.Rect) {
		fn(k.Canvas(p, r))
	})
}

// scale is s, or 1 for zero; outside 1–4 panics.
func scale(s int) int {
	s = pixel.Or("Scale", s, 1)
	if s > 4 {
		panic(fmt.Sprintf("text: Scale %d outside 1–4", s))
	}
	return s
}

// slotOr is s, or def for the zero slot.
func slotOr(s, def pixel.Slot) pixel.Slot {
	if s == pixel.Ink {
		return def
	}
	return s
}

// each walks s rune by rune, with x in dots.
func each(s string, sc int, fn func(i int, r rune, x int)) {
	i, x := 0, 0
	for _, r := range s {
		fn(i, r, x)
		x += pixel.Measure(string(r)) * sc
		i++
	}
}

// noise is a seeded hash for jitter and scrambles.
func noise(seed, i, step int) uint32 {
	x := uint32(seed*73856093 ^ i*19349663 ^ step*83492791)
	x ^= x >> 13
	x *= 0x5bd1e995
	return x ^ x>>15
}
