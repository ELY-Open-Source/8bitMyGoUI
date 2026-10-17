package data

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

// Score is a leaderboard entry.
type Score struct {
	Name  string
	Score int
}

type LeaderboardOptions struct {
	// Player is your row, marked; -1 marks none.
	Player int
}

// Leaderboard ranks scores; the top three wear medals.
func Leaderboard(c *ui.Context, scores []Score, o LeaderboardOptions) ui.Element {
	k := pixel.K(c)
	if o.Player < -1 || o.Player >= len(scores) {
		panic(fmt.Sprintf("data: Leaderboard Player %d of %d", o.Player, len(scores)))
	}
	return ui.Column(c).AlignItems(ui.Stretch).Role(ui.RoleTable).Children(func() {
		for i, s := range scores {
			you, ink := i == o.Player, pixel.Ink
			if you {
				ink = pixel.AccentInk
			}
			ui.Row(c).Height(k.D(pixel.LineHeight+4)).Gap(k.D(6)).Padding(0, k.D(4)).AlignItems(ui.Center).Role(ui.RoleRow).
				Label(fmt.Sprintf("%d %s %d", i+1, s.Name, s.Score)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				if you {
					g.Rect(0, 0, g.W, g.H, pixel.Accent)
				} else if i%2 == 1 {
					g.Veil(0, 0, g.W, g.H, pixel.Shade, 4)
				}
			}).Children(func() {
				ui.Box(c).Width(k.D(14)).Shrink(0).AlignItems(ui.Center).Children(func() {
					if i < 3 {
						display.Medal(c, display.Rank(i), display.MedalOptions{})
					} else {
						text.MicroText(c, strconv.Itoa(i+1), text.MicroTextOptions{Slot: ink})
					}
				})
				pixel.Text(c, s.Name, pixel.TextOptions{Slot: ink})
				ui.Box(c).Grow(1)
				text.MicroText(c, fmt.Sprintf("%07d", s.Score), text.MicroTextOptions{Slot: ink})
			})
		}
	})
}
