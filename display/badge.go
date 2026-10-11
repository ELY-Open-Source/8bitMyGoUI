package display

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type BadgeOptions struct {
	// Text replaces the count.
	Text string
	// Slot is the fill; zero means Bad.
	Slot pixel.Slot
}

// Badge shows n, capped at 99+, or text.
func Badge(c *ui.Context, n int, o BadgeOptions) *ui.Element {
	k := pixel.K(c)
	s := o.Text
	if s == "" {
		s = itoa(n)
		if n > 99 {
			s = "99+"
		}
	}
	fill := slotOr(o.Slot, pixel.Bad)
	if !pixel.MicroOK(s) {
		w := pixel.Measure(s) + 6
		return art(c, k, s, w, pixel.LineHeight+1, func(g pixel.Canvas) {
			g.Frame(0, 0, w, pixel.LineHeight+1, pixel.Notched, fill)
			g.Text(3, 0, s, pixel.AccentInk, 1)
		})
	}
	w := max(pixel.MicroWidth(s)+6, 9)
	return art(c, k, s, w, 9, func(g pixel.Canvas) {
		g.Frame(0, 0, w, 9, pixel.Notched, fill)
		g.Micro((w-pixel.MicroWidth(s))/2, 2, s, pixel.AccentInk, 1)
	})
}
