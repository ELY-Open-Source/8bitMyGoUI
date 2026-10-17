// Package display shows things: pictures, faces, badges and prizes.
package display

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// art is a labeled w×h dot image that fn draws.
func art(c *ui.Context, k pixel.Kit, label string, w, h int, fn func(g pixel.Canvas)) ui.Element {
	e := ui.Box(c).Size(k.D(w), k.D(h))
	if label != "" {
		e.Label(label).Role(ui.RoleImage)
	} else {
		e.Role(ui.RoleNone)
	}
	if fn != nil {
		e.Draw(func(p *ui.Painter, r ui.Rect) { fn(k.Canvas(p, r)) })
	}
	return e
}

// scale is s, or 1 for zero; outside 1–4 panics.
func scale(s int) int {
	s = pixel.Or("Scale", s, 1)
	if s > 4 {
		panic(fmt.Sprintf("display: Scale %d outside 1–4", s))
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

func itoa(n int) string { return fmt.Sprint(n) }
