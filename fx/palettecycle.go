package fx

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PaletteCycleOptions struct {
	// FPS is 4–15 steps a second; zero means 8.
	FPS int
	// Slot inks the sprite's # dots.
	Slot  pixel.Slot
	Scale int
	Label string
}

// PaletteCycle draws sp, rotating band's slots one place a step.
func PaletteCycle(c *ui.Context, sp *pixel.Sprite, band []pixel.Slot, o PaletteCycleOptions) ui.Element {
	k := pixel.K(c)
	fps, sc := pixel.Or("PaletteCycle FPS", o.FPS, 8), scale("PaletteCycle", o.Scale, 1)
	if len(band) < 2 {
		panic(fmt.Sprintf("fx: PaletteCycle band %v needs two slots", band))
	}
	return box(c, k, o.Label, sp.W*sc, sp.H*sc).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Cycled(band, g.Step(fps)).Sprite(0, 0, sp, 0, o.Slot, sc)
	})
}
