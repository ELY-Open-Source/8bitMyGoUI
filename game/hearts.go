package game

import (
	"fmt"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type HeartsOptions struct {
	// PerRow wraps the hearts; zero means 10.
	PerRow int
}

// Hearts shows hp of full hearts, in quarters.
func Hearts(c *ui.Context, text string, hp, full float64, o HeartsOptions) ui.Element {
	k := pixel.K(c)
	if full < 1 || full != math.Trunc(full) || !(hp >= 0 && hp <= full) || hp*4 != math.Trunc(hp*4) {
		panic(fmt.Sprintf("game: Hearts %v of %v; hearts are whole, health in quarters", hp, full))
	}
	per := pixel.Or("Hearts PerRow", o.PerRow, 10)
	n := int(full)
	cols, rows := min(n, per), (n+per-1)/per
	shown := strconv.FormatFloat(hp, 'f', -1, 64) + "/" + strconv.Itoa(n)
	return ui.Box(c).Role(ui.RoleMeter).Label(text).Description(shown).Size(k.D(cols*11+1), k.D(rows*10+2)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i := range n {
			x, y := i%per*11, i/per*10
			// Quarters fill from the bottom left, clockwise.
			q := min(4, max(0, int(hp*4)-4*i))
			for j, part := range [][2]int{{0, 6}, {0, 0}, {6, 0}, {6, 6}} {
				if j < q {
					g.Clip(x+part[0], y+part[1], 6, 6, func() { g.Sprite(x, y, icon.Heart, 0, pixel.Bad, 1) })
				}
			}
			g.Sprite(x, y, icon.HeartEmpty, 0, pixel.Edge, 1)
		}
	})
}
