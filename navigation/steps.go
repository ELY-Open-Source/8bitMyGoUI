package navigation

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// waving is a flag in two frames of wind.
var waving = pixel.MustSprite(`
	.#..........
	.#######....
	.########...
	.#########..
	.########...
	.#######....
	.#..........
	.#..........
	.#..........
	.#..........

	.#..........
	.##.#####...
	.#########..
	.#########..
	.##########.
	.##.#####...
	.#..........
	.#..........
	.#..........
	.#..........
`)

type StepsOptions struct{}

// Steps are flags on a path: done, waving, then locked.
func Steps(c *ui.Context, labels []string, current int, o StepsOptions) ui.Element {
	k := pixel.K(c)
	if current < 0 || current >= len(labels) {
		panic(fmt.Sprintf("navigation: step %d of %d", current, len(labels)))
	}
	ws, total := make([]int, len(labels)), 0
	for i, l := range labels {
		ws[i] = max(pixel.Measure(l), 12) + 16
		total += ws[i]
	}
	return ui.Box(c).Role(ui.RoleList).Label(strings.Join(labels, ", ")).Description(fmt.Sprintf("%d/%d", current+1, len(labels))).
		Size(k.D(total), k.D(14+pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		x := 0
		for i, l := range labels {
			mid := x + ws[i]/2
			if i > 0 {
				for d := x - ws[i-1]/2 + 10; d < mid-8; d += 3 {
					g.Dot(d, 6, pixel.Faint)
				}
			}
			switch {
			case i < current:
				g.Sprite(mid-6, 0, waving, 0, pixel.Good, 1)
			case i == current:
				g.Sprite(mid-6, 0, waving, g.Step(4)%2, pixel.Accent, 1)
			default:
				g.Sprite(mid-6, 0, icon.Lock, 0, pixel.Faint, 1)
			}
			ink := pixel.Faint
			if i <= current {
				ink = pixel.Ink
			}
			g.Text(mid-pixel.Measure(l)/2, 14, l, ink, 1)
			x += ws[i]
		}
	})
}
