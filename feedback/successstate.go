package feedback

import (
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

type SuccessStateOptions struct {
	// Stars pop in one by one; zero means 3.
	Stars int
}

// SuccessState is STAGE CLEAR! with stars popping in steps.
func SuccessState(c *ui.Context, message string, o SuccessStateOptions) *ui.Element {
	k := pixel.K(c)
	n := pixel.Or("SuccessState Stars", o.Stars, 3)
	col := ui.Column(c).Gap(k.D(8)).Padding(k.D(16)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(k.Msg("stage_clear"))
	start := ui.Local(col, startKey{}, func() time.Time { return k.Now() })
	// Stars come four a second; Reduce Motion shows them all.
	shown := n + 1
	if !k.Reduce() {
		shown = min(n+1, max(0, int(k.Now().Sub(*start)*4/time.Second)))
		if shown <= n {
			k.Tick(c, 8)
		}
	}
	return col.Children(func() {
		text.ShadowText(c, k.Msg("stage_clear"), text.ShadowTextOptions{Slot: pixel.Good, Scale: 3})
		w := icon.Star.W*2 + 4
		ui.Box(c).Size(k.D(n*w), k.D(icon.Star.H*2+4)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range min(shown, n) {
				if i == shown-1 {
					for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
						g.Sprite(i*w+2+d[0], 2+d[1], icon.Star, 0, pixel.Lit, 2)
					}
				}
				g.Sprite(i*w+2, 2, icon.Star, 0, pixel.Warn, 2)
			}
		})
		pixel.Text(c, message, pixel.TextOptions{})
	})
}
