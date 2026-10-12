package media

import (
	"fmt"
	"math"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type VisualizerOptions struct {
	// Height is in dots; zero means 32.
	Height int
}

// Visualizer mirrors each level, 0 to 1, about the middle.
func Visualizer(c *ui.Context, text string, levels []float64, o VisualizerOptions) *ui.Element {
	k := pixel.K(c)
	if len(levels) == 0 {
		panic("media: Visualizer has no levels")
	}
	parts := make([]string, len(levels))
	for i, l := range levels {
		if !(l >= 0 && l <= 1) {
			panic(fmt.Sprintf("media: Visualizer level %v outside 0–1", l))
		}
		parts[i] = fmt.Sprintf("%.2f", l)
	}
	h := pixel.Or("Visualizer Height", o.Height, 32)
	mid := h / 2
	return art(c, k, text, 4*len(levels)-1, h).Description(strings.Join(parts, " ")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, l := range levels {
			n := int(math.Round(l * float64(mid)))
			if n == 0 {
				continue
			}
			// Magic at the root, Accent above, an AccentLit tip.
			x, deep := 4*i, n/2
			g.Rect(x, mid-deep, 3, 2*deep, pixel.Magic)
			g.Rect(x, mid-n+1, 3, n-1-deep, pixel.Accent)
			g.Rect(x, mid+deep, 3, n-1-deep, pixel.Accent)
			g.Rect(x, mid-n, 3, 1, pixel.AccentLit)
			g.Rect(x, mid+n-1, 3, 1, pixel.AccentLit)
		}
	})
}
