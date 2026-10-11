package layout

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MasonryOptions struct {
	// Columns is zero for 3.
	Columns int
	// Gap is zero for 4.
	Gap int
}

// Masonry puts each item into the shortest column.
func Masonry(c *ui.Context, heights []int, o MasonryOptions, item func(i int)) *ui.Element {
	k := pixel.K(c)
	gap := orDef("Masonry Gap", o.Gap, 4)
	n := orDef("Masonry Columns", o.Columns, 3)
	cols := lanes(heights, n, gap)
	row := ui.Row(c).FillWidth().Gap(k.D(gap)).AlignItems(ui.Start)
	w, check := pixel.WrapWidth(c, row)
	widths := shares(w-gap*(n-1), n)
	return row.Draw(func(p *ui.Painter, r ui.Rect) { check(r) }).Children(func() {
		for j, lane := range cols {
			ui.Column(c).Width(k.D(widths[j])).Gap(k.D(gap)).Children(func() {
				for _, i := range lane {
					ui.Box(c).FillWidth().Height(k.D(heights[i])).Clip().Children(func() { item(i) })
				}
			})
		}
	})
}

// shares splits total into n whole parts, the first larger.
func shares(total, n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = max(0, total/n)
		if i < total%n {
			out[i]++
		}
	}
	return out
}

// lanes deals items into n columns, shortest first.
func lanes(heights []int, n, gap int) [][]int {
	out, tall := make([][]int, n), make([]int, n)
	for i, h := range heights {
		if h <= 0 {
			panic(fmt.Sprintf("layout: Masonry height %d at %d", h, i))
		}
		j := 0
		for k := range tall {
			if tall[k] < tall[j] {
				j = k
			}
		}
		out[j] = append(out[j], i)
		tall[j] += h + gap
	}
	return out
}
