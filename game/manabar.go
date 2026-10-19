package game

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type ManaBarOptions struct {
	// Width is in dots; zero means 96.
	Width int
	// Tick marks every Tick points; zero means 10.
	Tick int
}

// ManaBar shows mp of full, ticked every few points.
func ManaBar(c *ui.Context, text string, mp, full int, o ManaBarOptions) ui.Element {
	k := pixel.K(c)
	amount("ManaBar mp", mp, full)
	w := pixel.Or("ManaBar Width", o.Width, 96)
	fits("ManaBar", w, text, ratio(full, full))
	step := pixel.Or("ManaBar Tick", o.Tick, 10)
	top := headH(text)
	return ui.Box(c).Role(ui.RoleMeter).Label(text).Description(ratio(mp, full)).Size(k.D(w), k.D(top+8)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		head(g, w, text, ratio(mp, full))
		meter(g, 0, top, w, 8, float64(mp)/float64(full), pixel.Info)
		// Ticks under two dots apart are left out.
		if float64(step)/float64(full)*float64(w-4) < 2 {
			return
		}
		for v := step; v < full; v += step {
			g.VLine(2+min(w-5, int(float64(v)/float64(full)*float64(w-4))), top+2, 4, pixel.Edge)
			if v > full-step {
				break
			}
		}
	})
}
