package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PyramidOptions struct {
	// Width is in dots, zero means 200; rows are lines.
	Width int
}

// Pyramid sets two series back to back by group.
func Pyramid(c *ui.Context, text string, groups []string, left, right Series, o PyramidOptions) *ui.Element {
	k := pixel.K(c)
	if len(left.Values) != len(groups) || len(right.Values) != len(groups) {
		panic(fmt.Sprintf("chart: Pyramid has %d groups for %d and %d values", len(groups), len(left.Values), len(right.Values)))
	}
	w := pixel.Or("Pyramid Width", o.Width, 200)
	lo, hi := span(left.Values, right.Values)
	if lo < 0 {
		panic(fmt.Sprintf("chart: Pyramid value %v is negative", lo))
	}
	mid := 0
	for _, g := range groups {
		mid = max(mid, pixel.Measure(g)+12)
	}
	lslot, rslot := pick(left.Slot, 1), pick(right.Slot, 0)
	head := pixel.LineHeight
	return ui.Box(c).Size(k.D(w), k.D(head+len(groups)*pixel.LineHeight)).Role(ui.RoleImage).Label(text).
		Description(left.Name + ": " + read(groups, left.Values) + "; " + right.Name + ": " + read(groups, right.Values)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		reach := (w - mid) / 2
		cx := reach + mid/2
		g.Text(0, 0, left.Name, lslot, 1)
		g.Text(w-pixel.Measure(right.Name), 0, right.Name, rslot, 1)
		for i, name := range groups {
			y := head + i*pixel.LineHeight
			g.Text(cx-pixel.Measure(name)/2, y, name, pixel.Faint, 1)
			lw, rw := int(part(left.Values[i], hi)*float64(reach)), int(part(right.Values[i], hi)*float64(reach))
			g.Rect(reach-lw, y+4, lw, pixel.LineHeight-8, lslot)
			g.Rect(reach+mid, y+4, rw, pixel.LineHeight-8, rslot)
		}
	})
}
