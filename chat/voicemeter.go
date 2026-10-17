package chat

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type VoiceMeterOptions struct {
	// Bars is how many level bars; zero means 8.
	Bars int
}

// VoiceMeter lights bars for 0–1: green, amber, then red.
func VoiceMeter(c *ui.Context, text string, level float64, o VoiceMeterOptions) ui.Element {
	k := pixel.K(c)
	if math.IsNaN(level) || level < 0 || level > 1 {
		panic(fmt.Sprintf("chat: VoiceMeter level %v outside 0–1", level))
	}
	n := pixel.Or("VoiceMeter Bars", o.Bars, 8)
	if n < 2 {
		panic(fmt.Sprintf("chat: VoiceMeter needs two bars, got %d", n))
	}
	lit := int(math.Round(level * float64(n)))
	return ui.Box(c).Role(ui.RoleMeter).Label(text).Description(fmt.Sprintf("%d%%", int(math.Round(level*100)))).Size(k.D(3*n-1), k.D(12)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i := range n {
			h := 3 + i*9/(n-1)
			s := pixel.Well
			switch {
			case i >= lit:
			case i == n-1:
				s = pixel.Bad
			case i >= n*3/4:
				s = pixel.Warn
			default:
				s = pixel.Good
			}
			g.Rect(3*i, 12-h, 2, h, s)
		}
	})
}
