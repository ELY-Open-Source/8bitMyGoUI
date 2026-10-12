package navigation

import (
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var star5 = pixel.MustSprite(`
	..#..
	.###.
	#####
	.###.
	.#.#.
`)

type Stage struct {
	Name string
	// Locked stages cannot be entered.
	Locked bool
	// Cleared stages show Stars, 0 to 3.
	Cleared bool
	Stars   int
}

type LevelSelectOptions struct {
	// Columns is the grid's width; zero means 5.
	Columns  int
	Disabled bool
}

type LevelSelectResult struct {
	Element *ui.Element
	// Chosen is Enter or a click on an open stage.
	Changed, Chosen bool
}

// LevelSelect is a grid of locked, open or starred stages.
func LevelSelect(c *ui.Context, text string, stages []Stage, selected *int, o LevelSelectOptions) LevelSelectResult {
	k := pixel.K(c)
	cols := pixel.Or("LevelSelect Columns", o.Columns, 5)
	old := *selected
	shell, base, st, picked := pixel.Pad(c, text, cols, len(stages), 30, 30, selected, o.Disabled)
	res := LevelSelectResult{Element: shell, Changed: *selected != old, Chosen: picked >= 0 && !stages[picked].Locked}
	base.Description(stages[*selected].Name)
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, s := range stages {
			x, y := i%cols*30, i/cols*30
			fill := pixel.Face
			if s.Locked {
				fill = pixel.Shade
			}
			if i == *selected {
				g.Cursor(x, y, 30, 30)
			}
			g.Frame(x+2, y+2, 26, 26, pixel.Raised, fill)
			if s.Locked {
				g.Sprite(x+9, y+8, icon.Lock, 0, pixel.Faint, 1)
				continue
			}
			n := strconv.Itoa(i + 1)
			g.Micro(x+15-pixel.MicroWidth(n), y+5, n, pixel.Ink, 2)
			if s.Cleared {
				for j := range 3 {
					ink := pixel.Shade
					if j < s.Stars {
						ink = pixel.Warn
					}
					g.Sprite(x+5+j*7, y+20, star5, 0, ink, 1)
				}
			}
		}
		g.Finish(st)
	})
	if res.Chosen {
		slog.Debug("8bit: stage chosen", "label", text, "stage", stages[*selected].Name)
	}
	return res
}
