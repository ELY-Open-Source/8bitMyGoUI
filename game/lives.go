package game

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LivesOptions struct {
	// Sprite is the life icon; nil means Hero.
	Sprite *pixel.Sprite
}

// Lives shows a life icon and ×count.
func Lives(c *ui.Context, text string, n int, o LivesOptions) ui.Element {
	k := pixel.K(c)
	if n < 0 {
		panic(fmt.Sprintf("game: Lives %d", n))
	}
	sp := o.Sprite
	if sp == nil {
		sp = Hero
	}
	count := "×" + strconv.Itoa(n)
	return ui.Box(c).Role(ui.RoleImage).Label(text).Description(count).Size(k.D(sp.W+2+pixel.Measure(count)), k.D(max(sp.H, pixel.LineHeight))).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Sprite(0, (g.H-sp.H)/2, sp, 0, pixel.Ink, 1)
		g.Text(sp.W+2, (g.H-pixel.LineHeight)/2, count, pixel.Ink, 1)
	})
}
