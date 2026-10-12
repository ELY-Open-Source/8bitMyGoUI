package feedback

import (
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LevelUpOptions struct{}

type levelKey struct{}

// rise is the last level seen, and when it rose.
type rise struct {
	level int
	at    time.Time
}

// LevelUp slides in LEVEL UP! with sparkles when level rises.
func LevelUp(c *ui.Context, level int, o LevelUpOptions) *ui.Element {
	k := pixel.K(c)
	e := ui.Box(c)
	s := ui.Local(e, levelKey{}, func() rise { return rise{level: level} })
	if level > s.level {
		s.at = k.Now()
	}
	s.level = level
	// Ten steps a second: six slide in, twenty stay.
	step := since(c, k, s.at, 10, 26)
	if step >= 26 {
		return e
	}
	slide := 5 - min(step, 5)
	if k.Reduce() {
		slide = 0
	}
	title := k.Msg("level_up") + " LV " + strconv.Itoa(level)
	w := pixel.Measure(title)*2 + 24
	return e.Role(ui.RoleStatus).Label(title).Size(k.D(w), k.D(40)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		x := -slide * w / 5
		g.Clip(0, 0, g.W, g.H, func() {
			g.Frame(x, 4, w, 32, pixel.Window, pixel.Accent)
			g.Text(x+13, 5, title, pixel.Shade, 2)
			g.Text(x+12, 4, title, pixel.AccentInk, 2)
		})
		if slide > 0 {
			return
		}
		for i, at := range [][2]int{{2, 0}, {w - 14, 2}, {w / 2, 28}, {8, 26}} {
			if (g.Step(8)+i)%3 != 0 {
				g.Sprite(at[0], at[1], icon.Sparkle, 0, pixel.Warn, 1)
			}
		}
	})
}
