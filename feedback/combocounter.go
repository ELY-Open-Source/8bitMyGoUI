package feedback

import (
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ComboCounterOptions struct{}

type comboKey struct{}

// ComboCounter shows ×n COMBO!, popping bigger as n rises.
func ComboCounter(c *ui.Context, count int, o ComboCounterOptions) *ui.Element {
	k := pixel.K(c)
	e := ui.Box(c)
	s := ui.Local(e, comboKey{}, func() rise { return rise{level: count} })
	if count > s.level {
		s.at = k.Now()
	}
	s.level = count
	if count <= 0 {
		return e
	}
	pop := !k.Reduce() && k.Since(c, s.at, 10, 2) < 2
	n, tag := "×"+strconv.Itoa(count), k.Msg("combo")
	w := pixel.Measure(n)*3 + 6 + pixel.Measure(tag)*2 + 2
	return e.Role(ui.RoleStatus).Label(n+" "+tag).Size(k.D(w), k.D(3*pixel.LineHeight+2)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		x := pixel.Measure(n)*3 + 6
		word := func(dx, dy int, ink pixel.Slot) {
			g.Text(1+dx, 1+dy, n, ink, 3)
			g.Text(x+dx, pixel.LineHeight+dy, tag, ink, 2)
		}
		word(1, 1, pixel.Shade)
		if pop {
			for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				word(d[0], d[1], pixel.Warn)
			}
		}
		word(0, 0, pixel.Warn)
	})
}
