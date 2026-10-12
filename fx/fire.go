package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type FireOptions struct {
	// Scale multiplies the 16×20 campfire; zero means 2.
	Scale int
	Label string
}

// tongues are each flicker frame's three tongue heights.
var tongues = [3][3]int{{3, 6, 2}, {4, 5, 3}, {2, 7, 4}}

// bases are the flame's half widths from row 6 down.
var bases = [8]int{3, 4, 5, 5, 6, 6, 5, 4}

// flames are flicker frames: a Bad rim, a Warn body.
var flames = func() (out [3]*pixel.Sprite) {
	for f := range out {
		in := func(x, y int) bool {
			if y < 0 || y > 13 || x < 0 || x > 12 {
				return false
			}
			if y >= 6 && abs(x-6) <= bases[y-6] {
				return true
			}
			for i, h := range tongues[f] {
				if tip := 6 - h; y >= tip && y <= 6 && abs(x-3-3*i) <= (y-tip)/2 {
					return true
				}
			}
			return false
		}
		out[f] = pixel.SpriteOf(13, 14, func(x, y int) (pixel.Slot, bool) {
			switch {
			case !in(x, y):
				return 0, false
			case !in(x-1, y) || !in(x+1, y) || !in(x, y-1) || !in(x, y+1):
				return pixel.Bad, true
			}
			return pixel.Warn, true
		})
	}
	return out
}()

// logs cross under the fire; s marks the cut ends.
var logs = pixel.MustSprite(`
ss............ss
s###........###s
..####....####..
....########....
..####....####..
s###........###s
`)

// Fire burns a campfire: flicker frames, heat rising through it.
func Fire(c *ui.Context, o FireOptions) *ui.Element {
	k := pixel.K(c)
	sc := scale("Fire", o.Scale, 2)
	return box(c, k, o.Label, 16*sc, 20*sc).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		fire(g, k, sc, g.Step(10))
	})
}

// fire lights a third of every fourth body row, rising.
func fire(g pixel.Canvas, k pixel.Kit, sc, step int) {
	px := func(x, y, w, h int, s pixel.Slot) { g.Rect(x*sc, y*sc, w*sc, h*sc, s) }
	wood := g
	if k.Palette().Night {
		wood = g.Cycled([]pixel.Slot{pixel.Shade, pixel.Lit}, 1)
	}
	wood.Sprite(0, 14*sc, logs, 0, pick(k, pixel.Faint, pixel.Face), sc)
	sp, hot := flames[wrap(step, 3)], nightLit(g, k)
	hot.Sprite(sc, 2*sc, sp, 0, pixel.Ink, sc)
	// Body rows index the band; turning it lifts the heat.
	heat := pixel.Cycle([]pixel.Slot{pixel.Lit, pixel.Warn, pixel.Warn, pixel.Warn}, step)
	for y := range sp.H {
		for x := range sp.W {
			if s, ok := sp.At(0, x, y, pixel.Ink); ok && s == pixel.Warn && rnd(16*y+x, 200)%3 == 0 && heat[y%4] == pixel.Lit {
				hot.Rect((1+x)*sc, (2+y)*sc, sc, sc, pixel.Lit)
			}
		}
	}
	// Sparks rise beside the tips.
	for i := range 2 {
		if t := wrap(step+5*i, 8); t < 5 {
			px(4+8*i+sway(step+3*i, 8, 1), 5-t, 1, 1, pixel.Warn)
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
