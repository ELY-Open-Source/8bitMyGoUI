package chart

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ProgressRingOptions struct {
	// Size is the ring's width in dots; zero means 40.
	Size int
	Slot pixel.Slot
}

// ProgressRing fills a ring clockwise from twelve, in pixel arcs.
func ProgressRing(c *ui.Context, text string, value float64, o ProgressRingOptions) ui.Element {
	k := pixel.K(c)
	if !(value >= 0 && value <= 1) {
		panic(fmt.Sprintf("chart: ProgressRing value %v outside 0–1", value))
	}
	n := pixel.Or("ProgressRing Size", o.Size, 40)
	pct := strconv.Itoa(int(value*100)) + "%"
	ink := pick(o.Slot, 0)
	return ui.Box(c).Size(k.D(n), k.D(n)).Role(ui.RoleProgress).Label(text).Description(pct).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		c, rad := n/2, n/2-1
		sector(g, c, c, rad-4, rad, 0, 1, pixel.Well)
		sector(g, c, c, rad-4, rad, 0, value, ink)
		tiny(g, c-tinyW(pct)/2+1, c-2, pct, pixel.Ink)
	})
}
