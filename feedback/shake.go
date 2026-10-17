package feedback

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ShakeOptions struct {
	// Cue shakes the child each time it changes.
	Cue int
}

// shakes are the child's offsets in dots, a step each.
var shakes = [...]int{2, -2, 2, -2, 1, -1, 0}

// Shake jolts its child side to side briefly.
func Shake(c *ui.Context, o ShakeOptions, content func()) ui.Element {
	k := pixel.K(c)
	e := ui.Box(c)
	dx := shakes[k.Shot(c, e, o.Cue, 15, len(shakes)-1)]
	return e.Padding(0, k.D(2-dx), 0, k.D(2+dx)).Children(content)
}
