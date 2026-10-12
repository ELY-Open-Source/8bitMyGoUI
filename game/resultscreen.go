package game

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ResultScreenOptions struct{}

// ResultScreen counts time and score up, then shows the rank.
func ResultScreen(c *ui.Context, text string, took time.Duration, score int, rank string, o ResultScreenOptions) *ui.Element {
	k := pixel.K(c)
	if took < 0 || score < 0 || rank == "" {
		panic(fmt.Sprintf("game: ResultScreen %v, %d, rank %q", took, score, rank))
	}
	rows := [][2]string{{k.Msg("time"), ""}, {k.Msg("score"), ""}, {k.Msg("rank"), rank}}
	label := rows[0][0] + " " + lap(took) + ", " + rows[1][0] + " " + fmt.Sprint(score) + ", " + rows[2][0] + " " + rank
	e := ui.Box(c).Role(ui.RoleGroup).Label(text).Description(label)
	// Counts run 1.5 s from showing; reduced, they end.
	step := 15
	if !k.Reduce() {
		step = k.Since(c, start(k, e, 0).at, 10, 15)
	}
	rows[0][1] = lap(upTo(took, step, 15))
	rows[1][1] = fmt.Sprint(upTo(score, step, 15))
	name := 0
	for _, r := range rows {
		name = max(name, pixel.Measure(r[0]))
	}
	w := name + 12 + max(pixel.Measure(lap(took)), pixel.Measure(fmt.Sprint(score)), 3*pixel.Measure(rank))
	return e.Size(k.D(w), k.D(2*pixel.LineHeight+3*pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, row := range rows[:2] {
			g.Text(0, i*pixel.LineHeight, row[0], pixel.Faint, 1)
			g.Text(w-pixel.Measure(row[1]), i*pixel.LineHeight, row[1], pixel.Ink, 1)
		}
		g.Text(0, 2*pixel.LineHeight+pixel.LineHeight, rows[2][0], pixel.Faint, 1)
		if step == 15 {
			g.Text(w-3*pixel.Measure(rank), 2*pixel.LineHeight, rank, pixel.Warn, 3)
		}
	})
}

// lap prints a duration as M:SS.CC.
func lap(d time.Duration) string {
	cs := int(d / (10 * time.Millisecond))
	return fmt.Sprintf("%d:%02d.%02d", cs/6000, cs/100%60, cs%100)
}

// upTo is step/of of v, never forming the full product.
func upTo[T ~int | ~int64](v T, step, of int) T {
	return v/T(of)*T(step) + v%T(of)*T(step)/T(of)
}
