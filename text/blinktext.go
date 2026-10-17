package text

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type BlinkTextOptions struct {
	Slot  pixel.Slot
	Scale int
	// FPS steps a four-step blink, 4–15; zero means 4.
	FPS int
}

// BlinkText blinks; under Reduce Motion it stays shown.
func BlinkText(c *ui.Context, s string, o BlinkTextOptions) ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	fps := pixel.Or("BlinkText FPS", o.FPS, 4)
	return block(c, k, s, pixel.Measure(s)*sc, pixel.LineHeight*sc, func(g pixel.Canvas) {
		if g.Step(fps)%4 < 2 {
			g.Text(0, 0, s, o.Slot, sc)
		}
	})
}
