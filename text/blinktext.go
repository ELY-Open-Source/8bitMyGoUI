package text

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type BlinkTextOptions struct {
	Slot  pixel.Slot
	Scale int
	// FPS is blinks per second, times two; zero means 2.
	FPS int
}

// BlinkText blinks; under Reduce Motion it stays shown.
func BlinkText(c *ui.Context, s string, o BlinkTextOptions) *ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	fps := orDef("BlinkText FPS", o.FPS, 2)
	return block(c, k, s, pixel.Measure(s)*sc, pixel.LineHeight*sc, func(g pixel.Canvas) {
		if g.Step(fps)%2 == 0 {
			g.Text(0, 0, s, o.Slot, sc)
		}
	})
}
