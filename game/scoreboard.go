package game

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ScoreBoardOptions struct {
	// Digits pads both scores; zero means 6.
	Digits int
}

// ScoreBoard shows SCORE and HI-SCORE in big digits.
func ScoreBoard(c *ui.Context, score, hi int, o ScoreBoardOptions) ui.Element {
	k := pixel.K(c)
	if score < 0 || hi < 0 {
		panic(fmt.Sprintf("game: ScoreBoard %d and %d", score, hi))
	}
	digits := pixel.Or("ScoreBoard Digits", o.Digits, 6)
	// A score past the best becomes the best, lit.
	best, beat := max(hi, score), score > hi
	cols := []struct {
		name, n string
		ink     pixel.Slot
	}{{k.Msg("score"), fmt.Sprintf("%0*d", digits, score), pixel.Faint}, {k.Msg("hi_score"), fmt.Sprintf("%0*d", digits, best), pixel.Faint}}
	if beat {
		cols[1].ink = pixel.Warn
	}
	w, top := 0, 0
	for _, col := range cols {
		w, top = max(w, tagW(col.name), 2*pixel.Measure(col.n)), max(top, tagH(col.name)+2)
	}
	label := fmt.Sprintf("%s %d, %s %d", cols[0].name, score, cols[1].name, best)
	return ui.Box(c).Role(ui.RoleText).Label(label).Size(k.D(2*w+16), k.D(top+2*pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, col := range cols {
			x := i * (w + 16)
			tag(g, x, 0, col.name, col.ink)
			g.Text(x, top-2, col.n, pixel.Ink, 2)
		}
	})
}
