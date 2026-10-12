package display

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SpriteOptions struct {
	Slot  pixel.Slot
	Scale int
	// FPS plays the frames; zero shows the first.
	FPS    int
	Mirror bool
	Label  string
}

// Sprite plays a multi-slot sprite, mirrored when asked.
func Sprite(c *ui.Context, sp *pixel.Sprite, o SpriteOptions) *ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	fps := pixel.Or("Sprite FPS", o.FPS, 0)
	if o.Mirror {
		sp = sp.Mirror()
	}
	return art(c, k, o.Label, sp.W*sc, sp.H*sc, func(g pixel.Canvas) {
		frame := 0
		if fps > 0 {
			frame = g.Step(fps)
		}
		g.Sprite(0, 0, sp, frame, o.Slot, sc)
	})
}
