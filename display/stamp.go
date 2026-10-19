package display

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type StampOptions struct {
	// Slot is the ink; zero means Bad.
	Slot pixel.Slot
	Seed int
}

// Stamp is a rubber-stamp mark whose border wears unevenly.
func Stamp(c *ui.Context, s string, o StampOptions) ui.Element {
	k := pixel.K(c)
	ink := slotOr(o.Slot, pixel.Bad)
	w, h := pixel.Measure(s)+12, pixel.LineHeight+8
	return art(c, k, s, w, h, func(g pixel.Canvas) {
		ring := func(in int) {
			for i := in; i < w-in; i++ {
				worn(g, o.Seed, i, in, ink)
				worn(g, o.Seed, i, h-1-in, ink)
			}
			for j := in + 1; j < h-1-in; j++ {
				worn(g, o.Seed, in, j, ink)
				worn(g, o.Seed, w-1-in, j, ink)
			}
		}
		ring(0)
		ring(2)
		g.Text(6, 4, s, ink, 1)
	})
}

// worn inks a border dot unless seeded wear skips it.
func worn(g pixel.Canvas, seed, x, y int, ink pixel.Slot) {
	if hash(seed, x, y)%7 != 0 {
		g.Dot(x, y, ink)
	}
}

func hash(seed, x, y int) uint32 {
	n := uint32(seed*73856093 ^ x*19349663 ^ y*83492791)
	n ^= n >> 13
	n *= 0x5bd1e995
	return n ^ n>>15
}
