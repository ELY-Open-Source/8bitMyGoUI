package text

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type HeadingOptions struct {
	// Level 1–3 draws at ×3, ×2, ×1; zero means 1.
	Level int
	Slot  pixel.Slot
}

// Heading is a title with a 1-dot Shade shadow.
func Heading(c *ui.Context, s string, o HeadingOptions) *ui.Element {
	level := pixel.Or("Heading Level", o.Level, 1)
	if level > 3 {
		panic(fmt.Sprintf("text: Heading Level %d outside 1–3", level))
	}
	k, sc := pixel.K(c), 4-level
	return block(c, k, s, pixel.Measure(s)*sc+1, pixel.LineHeight*sc+1, func(g pixel.Canvas) {
		g.Text(1, 1, s, pixel.Shade, sc)
		g.Text(0, 0, s, o.Slot, sc)
	})
}
