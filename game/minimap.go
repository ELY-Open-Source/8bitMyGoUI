package game

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Spot is a map tile, counted from the top left.
type Spot struct{ X, Y int }

// Marker is a colored spot on a map.
type Marker struct {
	Spot
	Slot pixel.Slot
}

type MinimapOptions struct {
	// Cell is a tile's side in dots; zero means 3.
	Cell int
}

// tiles color map runes: wall, floor, water, grass, rock, void.
var tiles = map[rune]pixel.Slot{'#': pixel.Edge, '.': pixel.Well, '~': pixel.Info, '"': pixel.Good, '^': pixel.Faint, ' ': pixel.Back}

// Minimap draws tiles, fog where unseen, player and marks.
func Minimap(c *ui.Context, text string, rows []string, seen [][]bool, player Spot, marks []Marker, o MinimapOptions) *ui.Element {
	k := pixel.K(c)
	cell := pixel.Or("Minimap Cell", o.Cell, 3)
	if len(rows) == 0 {
		panic("game: Minimap has no rows")
	}
	grid := make([][]rune, len(rows))
	for y, row := range rows {
		grid[y] = []rune(row)
		if len(grid[y]) != len(grid[0]) {
			panic(fmt.Sprintf("game: Minimap row %d is %d tiles, not %d", y, len(grid[y]), len(grid[0])))
		}
		for _, r := range grid[y] {
			if _, ok := tiles[r]; !ok {
				panic(fmt.Sprintf("game: Minimap tile %q", r))
			}
		}
	}
	w, h := len(grid[0]), len(grid)
	inside := func(s Spot) bool { return s.X >= 0 && s.X < w && s.Y >= 0 && s.Y < h }
	for y := range seen {
		if len(seen) != h || len(seen[y]) != w {
			panic(fmt.Sprintf("game: Minimap fog is not %d×%d", w, h))
		}
	}
	for _, m := range append([]Marker{{Spot: player}}, marks...) {
		if !inside(m.Spot) {
			panic(fmt.Sprintf("game: Minimap spot %v is off the %d×%d map", m.Spot, w, h))
		}
	}
	at := strconv.Itoa(player.X) + ", " + strconv.Itoa(player.Y)
	return ui.Box(c).Role(ui.RoleImage).Label(text).Description(at).Size(k.D(w*cell+2), k.D(h*cell+2)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Flat, pixel.Back)
		for y, row := range grid {
			for x, t := range row {
				if seen != nil && !seen[y][x] {
					g.Dither(1+x*cell, 1+y*cell, cell, cell, pixel.Back, pixel.Edge, 12)
					continue
				}
				g.Rect(1+x*cell, 1+y*cell, cell, cell, tiles[t])
			}
		}
		for _, m := range marks {
			g.Rect(1+m.X*cell, 1+m.Y*cell, cell, cell, m.Slot)
		}
		// The player blinks at 2 Hz; reduced, it holds steady.
		if k.Reduce() || g.Busy(4)%2 == 0 {
			g.Rect(1+player.X*cell, 1+player.Y*cell, cell, cell, pixel.Accent)
		}
	})
}
