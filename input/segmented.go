package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SegmentedOptions struct {
	Disabled bool
}

type SegmentedResult struct {
	Element *ui.Element
	Changed bool
}

// Segmented joins segments choosing one; arrows move.
func Segmented(c *ui.Context, text string, labels []string, selected *int, o SegmentedOptions) SegmentedResult {
	k := pixel.K(c)
	p := ui.SegmentedBase(c, selected, len(labels))
	w := 0
	for _, l := range labels {
		w += pixel.Measure(l) + 12
	}
	p.Track.Width(k.D(w)).Label(text).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		g.Frame(0, 0, g.W, g.H, pixel.Notched, pixel.Face)
	}).Children(func() {
		for i, l := range labels {
			shell, _, st := pixel.Shell(c, func() *ui.Element { return p.Segment(i).Label(l).Disabled(o.Disabled) })
			shell.Size(k.D(pixel.Measure(l)+12), k.D(pixel.LineHeight+4)).Draw(func(pt *ui.Painter, r ui.Rect) {
				g := k.Canvas(pt, r)
				ink := pixel.Ink
				switch {
				case *selected == i:
					g.Rect(1, 1, g.W-2, g.H-2, pixel.Accent)
					ink = pixel.AccentInk
				case st.Hover:
					g.Rect(1, 1, g.W-2, g.H-2, pixel.Lit)
				}
				if i > 0 {
					g.VLine(0, 1, g.H-2, pixel.Edge)
				}
				g.Text(6, 2, l, ink, 1)
				finish(g, st)
			})
		}
	})
	res := SegmentedResult{Element: p.Track, Changed: p.Track.Changed()}
	if res.Changed {
		slog.Debug("8bit: segmented", "label", text, "selected", *selected)
	}
	return res
}
