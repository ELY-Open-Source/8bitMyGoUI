package game

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ATBGaugeOptions struct {
	// Width is in dots; zero means 64.
	Width int
}

// ATBGauge fills as a turn nears; full, it flashes READY.
func ATBGauge(c *ui.Context, text string, value float64, o ATBGaugeOptions) ui.Element {
	k := pixel.K(c)
	if !(value >= 0 && value <= 1) {
		panic(fmt.Sprintf("game: ATBGauge value %v outside 0–1", value))
	}
	w := pixel.Or("ATBGauge Width", o.Width, 64)
	ready := k.Msg("ready")
	fits("ATBGauge", w, text, ready)
	state := strconv.Itoa(int(value*100)) + "%"
	if value == 1 {
		state = ready
	}
	// Both labels sit on one line, tall enough for either.
	top := max(headH(text), headH(ready))
	return ui.Box(c).Role(ui.RoleMeter).Label(text).Description(state).Size(k.D(w), k.D(top+8)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		tag(g, 0, top-2-tagH(text), text, pixel.Ink)
		if value < 1 {
			meter(g, 0, top, w, 8, value, pixel.Info)
			return
		}
		// Full, READY flashes at 4 Hz; reduced, it holds.
		lit := g.Step(8)%2 == 0
		fill := pixel.Warn
		if !lit {
			fill = pixel.Lit
		}
		meter(g, 0, top, w, 8, 1, fill)
		if lit {
			tag(g, w-tagW(ready), top-2-tagH(ready), ready, pixel.Warn)
		}
	})
}
