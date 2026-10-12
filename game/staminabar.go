package game

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type StaminaBarOptions struct {
	// Segments are the bar's blocks; zero means 10.
	Segments int
}

// StaminaBar blocks stamina out; the part-full block dithers.
func StaminaBar(c *ui.Context, text string, value float64, o StaminaBarOptions) *ui.Element {
	k := pixel.K(c)
	if !(value >= 0 && value <= 1) {
		panic(fmt.Sprintf("game: StaminaBar value %v outside 0–1", value))
	}
	n := pixel.Or("StaminaBar Segments", o.Segments, 10)
	top := headH(text)
	w := max(n*6, tagW(text))
	return ui.Box(c).Role(ui.RoleMeter).Label(text).Description(strconv.Itoa(int(value*100))+"%").Size(k.D(w), k.D(top+6)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		tag(g, 0, 0, text, pixel.Ink)
		full := value * float64(n)
		for i := range n {
			x := i * 6
			switch part := full - float64(i); {
			case part >= 1:
				g.Rect(x, top, 5, 6, pixel.Good)
			case part > 0:
				g.Rect(x, top, 5, 6, pixel.Well)
				g.Veil(x, top, 5, 6, pixel.Good, max(1, int(part*16)))
			default:
				g.Checker(x, top, 5, 6, pixel.Shade)
			}
		}
	})
}
