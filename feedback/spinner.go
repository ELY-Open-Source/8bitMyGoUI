package feedback

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type SpinnerOptions struct {
	// Scale is 1–4; zero means 1.
	Scale int
}

// ring is where the spinner's eight beads sit, by hand.
var ring = [8][2]int{{5, 1}, {8, 2}, {9, 5}, {8, 8}, {5, 9}, {2, 8}, {1, 5}, {2, 2}}

// Spinner turns through eight frames, even under Reduce Motion.
func Spinner(c *ui.Context, text string, o SpinnerOptions) ui.Element {
	k := pixel.K(c)
	s := pixel.Or("Spinner Scale", o.Scale, 1)
	if s > 4 {
		panic("feedback: Spinner Scale past 4")
	}
	return ui.Box(c).Role(ui.RoleProgress).Label(text).Size(k.D(12*s), k.D(12*s)).Draw(func(p *ui.Painter, r ui.Rect) {
		beads(k.Canvas(p, r), 0, 0, s)
	})
}

// beads draws the spinner ring at x, y, scaled s.
func beads(g pixel.Canvas, x, y, s int) {
	head := g.Busy(10) % 8
	for i, at := range ring {
		ink := pixel.Shade
		switch (head - i + 8) % 8 {
		case 0:
			ink = pixel.Accent
		case 1:
			ink = pixel.AccentLit
		case 2:
			ink = pixel.Faint
		}
		g.Rect(x+at[0]*s, y+at[1]*s, 2*s, 2*s, ink)
	}
}
