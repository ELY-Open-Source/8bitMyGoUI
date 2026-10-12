package game

import (
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LetterboxOptions struct {
	// Bar is each bar's height in dots; zero means 16.
	Bar int
}

type barsKey struct{}

// bars is the last on, its change time and step.
type bars struct {
	on   bool
	at   time.Time
	from int
}

// Letterbox slides cutscene bars over content in four steps.
func Letterbox(c *ui.Context, on bool, o LetterboxOptions, content func()) *ui.Element {
	k := pixel.K(c)
	bar := pixel.Or("Letterbox Bar", o.Bar, 16)
	var cover *ui.Element
	e := layout.Layers(c, content, func() { cover = ui.Box(c).Fill().PassThrough() })
	// First shown, the bars are already where on puts them.
	s := ui.Local(e, barsKey{}, func() bars { return bars{on: on, from: 4 * boolInt(on)} })
	// shown is the step, 0–4; turns start where it is.
	shown := func() int {
		if k.Reduce() {
			return 4 * boolInt(s.on)
		}
		if n := k.Since(c, s.at, 8, 4); s.on {
			return min(4, s.from+n)
		} else {
			return max(0, s.from-n)
		}
	}
	if on != s.on {
		s.from, s.on, s.at = shown(), on, k.Now()
	}
	h := bar * shown() / 4
	cover.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, h, pixel.Ink)
		g.Rect(0, g.H-h, g.W, h, pixel.Ink)
	})
	return e
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
