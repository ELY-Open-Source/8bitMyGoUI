package feedback

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type HourglassOptions struct{}

// glass is the hourglass, by hand; # marks its walls.
var glass = []string{
	"############",
	".#........#.",
	".#........#.",
	"..#......#..",
	"...#....#...",
	"....#..#....",
	".....##.....",
	".....##.....",
	"....#..#....",
	"...#....#...",
	"..#......#..",
	".#........#.",
	".#........#.",
	"############",
}

// Hourglass drains in steps, then flips.
func Hourglass(c *ui.Context, text string, o HourglassOptions) *ui.Element {
	k := pixel.K(c)
	return ui.Box(c).Role(ui.RoleProgress).Label(text).Size(k.D(14), k.D(14)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		f := g.Busy(4) % 8
		// Frame 6 lies on its side, mid-flip; 7 rests.
		side := f == 6
		dot := func(x, y int, s pixel.Slot) {
			if side {
				g.Dot(y, 12-x, s)
			} else {
				g.Dot(x+1, y, s)
			}
		}
		top, bottom := 4-min(f, 4), min(f, 4)
		if f >= 6 {
			top, bottom = 4, 0
		}
		for y, row := range glass {
			l, r := strings.Index(row, "#"), strings.LastIndex(row, "#")
			for x := range len(row) {
				switch {
				case row[x] == '#':
					dot(x, y, pixel.Edge)
				case x > l && x < r && (y >= 6-top && y <= 5 || y >= 13-bottom && y <= 12):
					dot(x, y, pixel.Warn)
				}
			}
		}
		if f > 0 && f < 5 {
			for y := 6; y < 13-bottom; y++ {
				dot(5, y, pixel.Warn)
			}
		}
	})
}
