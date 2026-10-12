package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MosaicOptions struct {
	// Cue pixelates the picture and back each time it changes.
	Cue int
	// Slot inks the sprite's # dots.
	Slot  pixel.Slot
	Scale int
	Label string
}

// blocks are block sides a step each: growing, then back.
var blocks = [...]int{1, 2, 3, 4, 6, 8, 6, 4, 3, 2, 1}

// Mosaic pixelates sp in growing blocks, then sharpens it again.
func Mosaic(c *ui.Context, sp *pixel.Sprite, o MosaicOptions) *ui.Element {
	k := pixel.K(c)
	sc := scale("Mosaic", o.Scale, 1)
	e := box(c, k, o.Label, sp.W*sc, sp.H*sc)
	b := blocks[k.Shot(c, e, o.Cue, 10, len(blocks)-1)]
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		if b == 1 {
			g.Sprite(0, 0, sp, 0, o.Slot, sc)
			return
		}
		for y := 0; y < sp.H; y += b {
			for x := 0; x < sp.W; x += b {
				if s, ok := majority(sp, o.Slot, x, y, b); ok {
					g.Rect(x*sc, y*sc, min(b, sp.W-x)*sc, min(b, sp.H-y)*sc, s)
				}
			}
		}
	})
}

// majority is a block's commonest slot; clear wins ties.
func majority(sp *pixel.Sprite, ink pixel.Slot, x0, y0, b int) (pixel.Slot, bool) {
	var n [pixel.Slots + 1]int
	for y := y0; y < min(y0+b, sp.H); y++ {
		for x := x0; x < min(x0+b, sp.W); x++ {
			if s, ok := sp.At(0, x, y, ink); ok {
				n[s]++
			} else {
				n[pixel.Slots]++
			}
		}
	}
	best := pixel.Slots
	for s := range pixel.Slots {
		if n[s] > n[best] {
			best = s
		}
	}
	return pixel.Slot(best), best != pixel.Slots
}
