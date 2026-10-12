package game

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TimerHUDOptions struct{}

// TimerHUD shows TIME in seconds; under 30 it blinks red.
func TimerHUD(c *ui.Context, left time.Duration, o TimerHUDOptions) *ui.Element {
	k := pixel.K(c)
	if left < 0 {
		panic(fmt.Sprintf("game: TimerHUD %v left", left))
	}
	secs := int(math.Ceil(left.Seconds()))
	name, n := k.Msg("time"), fmt.Sprintf("%03d", secs)
	hurry := left < 30*time.Second
	top := tagH(name) + 2
	w := max(tagW(name), 2*pixel.Measure(n))
	return ui.Box(c).Role(ui.RoleText).Label(name+" "+n).Size(k.D(w), k.D(top+2*pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		tag(g, 0, 0, name, pixel.Faint)
		ink := pixel.Ink
		if hurry {
			ink = pixel.Bad
			// It blinks at 2 Hz; reduced, it holds steady.
			if !k.Reduce() && g.Busy(4)%2 == 1 {
				return
			}
		}
		g.Text(0, top-2, n, ink, 2)
	})
}
