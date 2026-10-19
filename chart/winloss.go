package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type WinLossOptions struct{}

// WinLoss rows blocks: wins up, losses down, draws dashed.
func WinLoss(c *ui.Context, text string, results []int, o WinLossOptions) ui.Element {
	k := pixel.K(c)
	return ui.Box(c).Size(k.D(len(results)*5), k.D(pixel.LineHeight)).Role(ui.RoleImage).Label(text).
		Description(record(results)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		mid := pixel.LineHeight / 2
		for i, res := range results {
			x := i * 5
			switch res {
			case 1:
				g.Rect(x, mid-6, 4, 5, pixel.Good)
			case -1:
				g.Rect(x, mid+1, 4, 5, pixel.Bad)
			default:
				g.Rect(x, mid-1, 4, 2, pixel.Faint)
			}
		}
	})
}

// record is "wins-losses-draws: " then each result, checked.
func record(results []int) string {
	var n [3]int
	all := make([]float64, len(results))
	for i, r := range results {
		if r < -1 || r > 1 {
			panic(fmt.Sprintf("chart: WinLoss result %d is not -1, 0 or 1", r))
		}
		n[1-r]++
		all[i] = float64(r)
	}
	return fmt.Sprintf("%d-%d-%d: %s", n[0], n[2], n[1], nums(all))
}
