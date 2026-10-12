package tabletop

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DiceOptions struct {
	// Cue rolls it: a new cue tumbles, then lands it.
	Cue int
}

// hop is how high a tumbling die flies, each frame.
var hop = [...]int{0, 3, 5, 6, 6, 5, 3, 1}

// Dice is a d6 showing value; new Cues tumble it.
func Dice(c *ui.Context, text string, value int, o DiceOptions) *ui.Element {
	k := pixel.K(c)
	if value < 1 || value > 6 {
		panic(fmt.Sprintf("tabletop: Dice value %d", value))
	}
	e := ui.Box(c).Size(k.D(20), k.D(26)).Role(ui.RoleImage).Label(text).Description(strconv.Itoa(value))
	// A roll tumbles eight frames; reduced, it lands at once.
	f := max(0, k.Shot(c, e, o.Cue, 12, len(hop)))
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		face, y, h := value, 8, 16
		if f < len(hop) {
			// Turning, it shows other faces, squashed every other frame.
			face, y = (value+f*5)%6+1, 8-hop[f]
			if f%2 == 1 {
				y, h = y+2, 12
			}
		}
		g.Frame(2, y, 16, h, pixel.Raised, pixel.Lit)
		ink := pixel.Ink
		if face == 1 {
			ink = pixel.Bad
		}
		if h == 16 {
			pips(g, 2, y, 16, face, ink)
		}
	})
}
