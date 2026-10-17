package game

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Award is an achievement; locked ones show as silhouettes.
type Award struct {
	Sprite *pixel.Sprite
	Name   string
	// How tells how it is earned.
	How      string
	Unlocked bool
}

type AchievementsOptions struct {
	// Columns is the grid's width; zero means 6.
	Columns  int
	Disabled bool
}

type AchievementsResult struct {
	Element ui.Element
	Changed bool
}

// tile is an award's cell, in dots.
const tile = 24

// Achievements grids awards; arrows read each one's name and how.
func Achievements(c *ui.Context, text string, awards []Award, selected *int, o AchievementsOptions) AchievementsResult {
	k := pixel.K(c)
	cols := pixel.Or("Achievements Columns", o.Columns, 6)
	if len(awards) == 0 || *selected < 0 || *selected >= len(awards) {
		panic(fmt.Sprintf("game: Achievements award %d of %d", *selected, len(awards)))
	}
	got := 0
	for _, a := range awards {
		if a.Sprite == nil {
			panic(fmt.Sprintf("game: award %q has no sprite", a.Name))
		}
		if a.Unlocked {
			got++
		}
	}
	old := *selected
	col := ui.Column(c).Gap(k.D(4))
	res := AchievementsResult{Element: col}
	col.Children(func() {
		pixel.Text(c, fmt.Sprintf(k.Msg("unlocked"), got, len(awards)), pixel.TextOptions{Slot: pixel.Accent})
		shell, base, st, _ := pixel.Pad(c, text, cols, len(awards), tile, tile, selected, o.Disabled)
		a := awards[*selected]
		name := a.Name
		if !a.Unlocked {
			name = k.Msg("unknown")
		}
		base.Description(name + ": " + a.How)
		sel, shown := *selected, append([]Award(nil), awards...)
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, a := range shown {
				x, y := i%cols*tile, i/cols*tile
				sx, sy := x+(tile-a.Sprite.W)/2, y+(tile-a.Sprite.H)/2
				if a.Unlocked {
					g.Frame(x+1, y+1, tile-2, tile-2, pixel.Raised, pixel.Lit)
					g.Sprite(sx, sy, a.Sprite, 0, pixel.Ink, 1)
				} else {
					g.Frame(x+1, y+1, tile-2, tile-2, pixel.Sunken, pixel.Well)
					g.Silhouette(sx, sy, a.Sprite, 0, pixel.Shade, 1)
				}
				if i == sel {
					g.Rect(x, y, tile, 1, pixel.Accent)
					g.Rect(x, y+tile-1, tile, 1, pixel.Accent)
					g.Rect(x, y, 1, tile, pixel.Accent)
					g.Rect(x+tile-1, y, 1, tile, pixel.Accent)
				}
			}
			g.Finish(st)
		})
		pixel.Text(c, name, pixel.TextOptions{})
		pixel.Text(c, a.How, pixel.TextOptions{Slot: pixel.Faint, Wrap: true})
	})
	res.Changed = *selected != old
	if res.Changed {
		slog.Debug("8bit: achievement", "label", text, "award", *selected)
	}
	return res
}
