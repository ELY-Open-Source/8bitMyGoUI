package game

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/feedback"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CoinCounterOptions struct {
	// Scale is 1–4; zero means 1.
	Scale int
}

// CoinCounter spins a coin beside a count that ticks up.
func CoinCounter(c *ui.Context, text string, coins int, o CoinCounterOptions) ui.Element {
	k := pixel.K(c)
	if coins < 0 {
		panic(fmt.Sprintf("game: CoinCounter %d", coins))
	}
	sc := pixel.Or("CoinCounter Scale", o.Scale, 1)
	return ui.Row(c).Gap(k.D(2 * sc)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(text).Children(func() {
		display.Sprite(c, coin, display.SpriteOptions{Scale: sc, FPS: 8})
		feedback.CountUp(c, coins, feedback.CountUpOptions{Scale: sc})
	})
}
