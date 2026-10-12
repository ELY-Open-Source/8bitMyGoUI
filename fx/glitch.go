package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type GlitchOptions struct {
	// Cue glitches the picture each time it changes.
	Cue int
	// Slot inks the sprite's # dots.
	Slot  pixel.Slot
	Scale int
	Label string
}

const glitchSteps = 6

// Glitch shifts row bands and swaps two slots, briefly.
func Glitch(c *ui.Context, sp *pixel.Sprite, o GlitchOptions) *ui.Element {
	k := pixel.K(c)
	sc := scale("Glitch", o.Scale, 1)
	// Three dots either side hold the rows shifted out.
	e := box(c, k, o.Label, (sp.W+6)*sc, sp.H*sc)
	step := k.Shot(c, e, o.Cue, 15, glitchSteps)
	used := slotsOf(sp, o.Slot)
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		if step == glitchSteps {
			g.Sprite(3*sc, 0, sp, 0, o.Slot, sc)
			return
		}
		seed := o.Cue*glitchSteps + step
		if len(used) > 1 {
			a := rnd(seed, 160) % len(used)
			b := (a + 1 + rnd(seed, 161)%(len(used)-1)) % len(used)
			g = g.Cycled([]pixel.Slot{used[a], used[b]}, 1)
		}
		for y := 0; y < sp.H; {
			h, dx := 1+rnd(seed+y, 162)%4, 0
			if rnd(seed+y, 163)%3 == 0 {
				dx = rnd(seed+y, 164)%7 - 3
			}
			g.Clip(0, y*sc, g.W, h*sc, func() { g.Sprite((3+dx)*sc, 0, sp, 0, o.Slot, sc) })
			y += h
		}
	})
}

// slotsOf lists the slots sp's first frame uses.
func slotsOf(sp *pixel.Sprite, ink pixel.Slot) []pixel.Slot {
	var seen [pixel.Slots]bool
	var out []pixel.Slot
	for y := range sp.H {
		for x := range sp.W {
			if s, ok := sp.At(0, x, y, ink); ok && !seen[s] {
				seen[s] = true
				out = append(out, s)
			}
		}
	}
	return out
}
