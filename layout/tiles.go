package layout

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TilesOptions struct {
	// W is zero for 48; H is zero for W.
	W, H int
	// Gap is zero for 4.
	Gap int
}

// Tiles wraps n equal tiles into rows.
func Tiles(c *ui.Context, n int, o TilesOptions, tile func(i int)) *ui.Element {
	k := pixel.K(c)
	w := orDef("Tiles W", o.W, 48)
	h := orDef("Tiles H", o.H, w)
	orDef("Tiles n", n, 0)
	return ui.Row(c).Wrap().Gap(k.D(orDef("Tiles Gap", o.Gap, 4))).Children(func() {
		for i := range n {
			ui.Box(c).Size(k.D(w), k.D(h)).Clip().Children(func() { tile(i) })
		}
	})
}
