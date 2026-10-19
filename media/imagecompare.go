package media

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type ImageCompareOptions struct {
	// Scale is dots per picture dot, 1–4; zero means 2.
	Scale    int
	Disabled bool
}

type ImageCompareResult struct {
	Element ui.Element
	Changed bool
}

// ImageCompare shows before left of *split, after right.
func ImageCompare(c *ui.Context, text string, before, after *pixel.Sprite, split *float64, o ImageCompareOptions) ImageCompareResult {
	k := pixel.K(c)
	sameSize("ImageCompare", []*pixel.Sprite{before, after})
	if !(*split >= 0 && *split <= 1) {
		panic(fmt.Sprintf("media: ImageCompare split %v outside 0–1", *split))
	}
	s := scale("ImageCompare", o.Scale, 2)
	w, h := before.W*s, before.H*s
	shell, st, changed := input.SlideBase(c, text, split, input.RailOptions{Disabled: o.Disabled})
	if changed {
		slog.Debug("8bit: image compare", "label", text, "split", *split)
	}
	at := *split
	shell.Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		x := int(math.Round(at * float64(w)))
		g.Sprite(0, 0, before, 0, pixel.Ink, s)
		g.Clip(x, 0, w-x, h, func() { g.Sprite(0, 0, after, 0, pixel.Ink, s) })
		g.VLine(min(x, w-1), 0, h, pixel.Edge)
		// The grip stays whole at either edge.
		gx := max(0, min(w-7, x-3))
		g.Frame(gx, h/2-5, 7, 11, pixel.Raised, pixel.Face)
		for _, d := range []int{-2, 0, 2} {
			g.HLine(gx+2, h/2+d, 3, pixel.Shade)
		}
		g.Finish(st)
	})
	return ImageCompareResult{Element: shell, Changed: changed}
}
