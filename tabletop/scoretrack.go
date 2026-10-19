package tabletop

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type ScoreTrackOptions struct {
	// Spaces on the track; zero means 30.
	Spaces int
}

const space = 20

// ScoreTrack winds a path; player meeples mark their scores.
func ScoreTrack(c *ui.Context, text string, scores []int, o ScoreTrackOptions) ui.Element {
	k := pixel.K(c)
	n := pixel.Or("ScoreTrack Spaces", o.Spaces, 30)
	// A space prints three digits at most.
	if len(scores) == 0 || len(scores) > len(players) || n > 999 {
		panic(fmt.Sprintf("tabletop: ScoreTrack for %d players on %d spaces", len(scores), n))
	}
	said := make([]string, len(scores))
	for i, s := range scores {
		if s < 0 {
			panic(fmt.Sprintf("tabletop: ScoreTrack score %d", s))
		}
		said[i] = fmt.Sprintf(k.Msg("player_score"), i+1, s)
	}
	cols := min(n, 10)
	rows := (n + cols - 1) / cols
	// at places space i; rows snake back and forth.
	at := func(i int) (int, int) {
		x, y := i%cols, i/cols
		if y%2 == 1 {
			x = cols - 1 - x
		}
		return x * space, y * space
	}
	return ui.Box(c).Size(k.D(cols*space), k.D(rows*space)).Role(ui.RoleImage).Label(text).Description(strings.Join(said, ", ")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i := range n {
			x, y := at(i)
			fill := pixel.Face
			if i%5 == 0 {
				fill = pixel.Lit
			}
			g.Frame(x, y, space, space, pixel.Flat, fill)
			g.Micro(x+2, y+2, strconv.Itoa(i), pixel.Faint, 1)
		}
		// A lone player stands as a meeple; crowds, small pawns.
		crowd := map[int]int{}
		for _, s := range scores {
			crowd[s%n]++
		}
		seen := map[int]int{}
		for i, s := range scores {
			x, y := at(s % n)
			if crowd[s%n] == 1 {
				g.Silhouette(x+3, y+7, icon.Meeple, 0, pixel.Edge, 1)
				g.Sprite(x+4, y+7, icon.Meeple, 0, players[i], 1)
				continue
			}
			j := seen[s%n]
			seen[s%n]++
			g.Rect(x+2+j%3*6, y+8+j/3*6, 5, 5, pixel.Edge)
			g.Rect(x+3+j%3*6, y+9+j/3*6, 3, 3, players[i])
		}
	})
}
