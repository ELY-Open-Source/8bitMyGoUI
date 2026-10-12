package data

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Change is a stat before and after.
type Change struct {
	Label    string
	Old, New int
}

type StatCompareOptions struct{}

// StatCompare lists stats old → new with ▲▼ deltas.
func StatCompare(c *ui.Context, changes []Change, o StatCompareOptions) *ui.Element {
	k := pixel.K(c)
	lw, ow, nw, dw := 0, 0, 0, 0
	deltas := make([]string, len(changes))
	for i, ch := range changes {
		switch d := ch.New - ch.Old; {
		case d > 0:
			deltas[i] = "▲" + strconv.Itoa(d)
		case d < 0:
			deltas[i] = "▼" + strconv.Itoa(-d)
		}
		lw = max(lw, pixel.Measure(ch.Label))
		ow = max(ow, pixel.Measure(strconv.Itoa(ch.Old)))
		nw = max(nw, pixel.Measure(strconv.Itoa(ch.New)))
		dw = max(dw, pixel.Measure(deltas[i]))
	}
	arrow := lw + 8 + ow + 4
	newAt := arrow + pixel.Measure("→") + 4
	deltaAt := newAt + nw + 6
	return ui.Column(c).Children(func() {
		for i, ch := range changes {
			delta, ink := deltas[i], pixel.Good
			if ch.New < ch.Old {
				ink = pixel.Bad
			}
			ui.Box(c).Size(k.D(deltaAt+dw), k.D(pixel.LineHeight)).Role(ui.RoleText).
				Label(fmt.Sprintf("%s %d → %d %s", ch.Label, ch.Old, ch.New, delta)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				old, now := strconv.Itoa(ch.Old), strconv.Itoa(ch.New)
				g.Text(0, 0, ch.Label, pixel.Faint, 1)
				g.Text(arrow-4-pixel.Measure(old), 0, old, pixel.Ink, 1)
				g.Text(arrow, 0, "→", pixel.Faint, 1)
				g.Text(newAt+nw-pixel.Measure(now), 0, now, pixel.Ink, 1)
				g.Text(deltaAt, 0, delta, ink, 1)
			})
		}
	})
}
