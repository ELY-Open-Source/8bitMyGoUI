package game

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type EndingCardOptions struct {
	// Width and Height are dots; zeros mean 200 by 80.
	Width, Height int
	// Cue fades it in again each time it changes.
	Cue int
}

// EndingCard dithers THE END in, step by step.
func EndingCard(c *ui.Context, o EndingCardOptions) *ui.Element {
	k := pixel.K(c)
	w, h := pixel.Or("EndingCard Width", o.Width, 200), pixel.Or("EndingCard Height", o.Height, 80)
	s := k.Msg("the_end")
	if w < 3*pixel.Measure(s) || h < 3*pixel.LineHeight {
		panic(fmt.Sprintf("game: EndingCard %d×%d is too small for %q", w, h, s))
	}
	e := ui.Box(c).Role(ui.RoleText).Label(s).Size(k.D(w), k.D(h))
	level := 16
	// It fades in over 2 s; reduced, at once.
	if !k.Reduce() {
		level = k.Since(c, start(k, e, o.Cue).at, 8, 16)
	}
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, w, h, pixel.Ink)
		g.TextVeil((w-3*pixel.Measure(s))/2, (h-3*pixel.LineHeight)/2, s, pixel.Back, level, 3)
	})
}
