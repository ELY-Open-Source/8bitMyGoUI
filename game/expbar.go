package game

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ExpBarOptions struct {
	// Width is in dots; zero means 96.
	Width int
}

type riseKey struct{}

// rise is the last level, xp and fill time.
type rise struct {
	level, xp int
	at        time.Time
}

// ExpBar shows xp toward next and the level; ups flash.
func ExpBar(c *ui.Context, text string, level, xp, next int, o ExpBarOptions) *ui.Element {
	k := pixel.K(c)
	if level < 1 {
		panic(fmt.Sprintf("game: ExpBar level %d", level))
	}
	amount("ExpBar xp", xp, next)
	w := pixel.Or("ExpBar Width", o.Width, 96)
	lv := fmt.Sprintf(k.Msg("level"), level)
	fits("ExpBar", w, lv, ratio(next, next))
	e := ui.Box(c).Role(ui.RoleMeter).Label(text).Description(lv + " " + ratio(xp, next))
	up := ui.Local(e, riseKey{}, func() rise { return rise{level: level, xp: xp} })
	// It fills as xp reaches next or the level rises.
	if level > up.level || level == up.level && xp == next && up.xp < next {
		up.at = k.Now()
	}
	up.level, up.xp = level, xp
	// A fill flashes 0.6 s; reduced, it jumps past.
	flash := !k.Reduce() && k.Since(c, up.at, 10, 6) < 6
	top := headH(lv)
	return e.Size(k.D(w), k.D(top+8)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		head(g, w, lv, ratio(xp, next))
		switch {
		case !flash:
			meter(g, 0, top, w, 8, float64(xp)/float64(next), pixel.Warn)
		case g.Busy(10)%2 == 0:
			meter(g, 0, top, w, 8, 1, pixel.Lit)
		default:
			meter(g, 0, top, w, 8, 1, pixel.Warn)
		}
	})
}
