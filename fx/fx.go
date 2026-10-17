// Package fx draws weather, scenes and transitions in stepped frames.
package fx

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// box sizes an effect; a label makes it an image.
func box(c *ui.Context, k pixel.Kit, label string, w, h int) ui.Element {
	e := ui.Box(c).Size(k.D(w), k.D(h))
	if label != "" {
		return e.Label(label).Role(ui.RoleImage)
	}
	return e.Role(ui.RoleNone)
}

// size is w×h; zeros take dw×dh.
func size(name string, w, h, dw, dh int) (int, int) {
	return pixel.Or(name+" Width", w, dw), pixel.Or(name+" Height", h, dh)
}

// scale is s, or def for zero; above 4 panics.
func scale(name string, s, def int) int {
	s = pixel.Or(name+" Scale", s, def)
	if s > 4 {
		panic(fmt.Sprintf("fx: %s Scale %d outside 1–4", name, s))
	}
	return s
}

// rnd is a fixed pseudo-random number for i and salt.
func rnd(i, salt int) int {
	x := uint64(i)*0x9E3779B97F4A7C15 + uint64(salt)*0xD1B54A32D192ED03
	x = (x ^ x>>30) * 0xBF58476D1CE4E5B9
	x = (x ^ x>>27) * 0x94D049BB133111EB
	return int((x ^ x>>31) >> 1)
}

// wrap is v modulo n, never negative.
func wrap(v, n int) int { return (v%n + n) % n }

// sway swings 0 to a, back through -a, over p.
func sway(t, p, a int) int {
	t = wrap(t, p)
	q := p / 4
	switch {
	case t < q:
		return a * t / q
	case t < 3*q:
		return a - 2*a*(t-q)/(2*q)
	default:
		return -a + a*(t-3*q)/q
	}
}

// pick is day under day palettes, night under night ones.
func pick(k pixel.Kit, day, night pixel.Slot) pixel.Slot {
	if k.Palette().Night {
		return night
	}
	return day
}

// light is the palette's brightest slot.
func light(k pixel.Kit) pixel.Slot { return pick(k, pixel.Lit, pixel.Ink) }

// dark is the palette's darkest slot.
func dark(k pixel.Kit) pixel.Slot { return pick(k, pixel.Ink, pixel.Well) }

// slotOr is s, or def for the zero slot.
func slotOr(s, def pixel.Slot) pixel.Slot {
	if s == pixel.Ink {
		return def
	}
	return s
}

// checkHour panics outside 0–23.
func checkHour(name string, hour int) {
	if hour < 0 || hour > 23 {
		panic(fmt.Sprintf("fx: %s hour %d outside 0–23", name, hour))
	}
}

// cover draws over e's content, clipped to e.
func cover(k pixel.Kit, e ui.Element, draw func(g pixel.Canvas), content func()) ui.Element {
	return e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { draw(g) })
	}).Children(content)
}
