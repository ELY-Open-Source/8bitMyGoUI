package tabletop

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Stone is a Go point: Empty, Black or White.
type Stone int8

const (
	Empty Stone = iota
	Black
	White
)

type GoBoardOptions struct {
	Disabled bool
}

type GoBoardResult struct {
	Element ui.Element
	// Placed is the point played this frame, else -1.
	Placed int
}

type goKey struct{}

const point = 10

// GoBoard places stones in turn, Black first, 9–19 lines.
func GoBoard(c *ui.Context, text string, board []Stone, o GoBoardOptions) GoBoardResult {
	k := pixel.K(c)
	size := map[int]int{81: 9, 169: 13, 361: 19}[len(board)]
	if size == 0 {
		panic(fmt.Sprintf("tabletop: GoBoard of %d points is not 9, 13 or 19 lines", len(board)))
	}
	blacks, whites := 0, 0
	for i, s := range board {
		switch s {
		case Black:
			blacks++
		case White:
			whites++
		case Empty:
		default:
			panic(fmt.Sprintf("tabletop: GoBoard point %d holds %d", i, s))
		}
	}
	turn := Black
	if blacks > whites {
		turn = White
	}
	col := ui.Column(c).Gap(k.D(4))
	cur := ui.Local(col, goKey{}, func() int { return len(board) / 2 })
	// A smaller board pulls the cursor in.
	*cur = min(*cur, len(board)-1)
	res := GoBoardResult{Placed: -1}
	res.Element = col.Children(func() {
		side := k.Msg("black")
		if turn == White {
			side = k.Msg("white")
		}
		next := fmt.Sprintf(k.Msg("to_play"), side)
		pixel.Text(c, next, pixel.TextOptions{})
		shell, base, st, picked := pixel.Pad(c, text, size, len(board), point, point, cur, o.Disabled)
		if picked >= 0 && board[picked] == Empty {
			board[picked], res.Placed = turn, picked
		}
		base.Description(goPoint(size, *cur) + " " + stoneName(k, board[*cur]))
		at, shown := *cur, append([]Stone(nil), board...)
		light, dark := shades(k)
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Rect(0, 0, g.W, g.H, pixel.Face)
			end := (size-1)*point + 1
			for i := range size {
				g.HLine(point/2, point/2+i*point, end, pixel.Ink)
				g.VLine(point/2+i*point, point/2, end, pixel.Ink)
			}
			// Star points sit three or four lines in, and mid-board.
			edge := 2
			if size > 9 {
				edge = 3
			}
			for _, x := range []int{edge, size / 2, size - 1 - edge} {
				for _, y := range []int{edge, size / 2, size - 1 - edge} {
					g.Rect(point/2+x*point-1, point/2+y*point-1, 3, 3, pixel.Ink)
				}
			}
			for i, s := range shown {
				cx, cy := point/2+i%size*point, point/2+i/size*point
				switch s {
				case Black:
					g.Disc(cx, cy, 4, dark)
				case White:
					g.Disc(cx, cy, 4, light)
					g.Circle(cx, cy, 4, dark)
				}
			}
			if st.Focus {
				x, y := at%size*point, at/size*point
				g.Rect(x, y, point, 1, pixel.Accent)
				g.Rect(x, y+point-1, point, 1, pixel.Accent)
				g.Rect(x, y, 1, point, pixel.Accent)
				g.Rect(x+point-1, y, 1, point, pixel.Accent)
			}
			g.Finish(st)
		})
	})
	if res.Placed >= 0 {
		slog.Debug("8bit: go stone placed", "label", text, "point", goPoint(size, res.Placed))
	}
	return res
}

// goPoint names a point as Go does: columns skip I.
func goPoint(size, i int) string {
	return fmt.Sprintf("%c%d", "ABCDEFGHJKLMNOPQRST"[i%size], size-i/size)
}

// stoneName reads a stone out, or the point's emptiness.
func stoneName(k pixel.Kit, s Stone) string {
	switch s {
	case Black:
		return fmt.Sprintf(k.Msg("stone"), k.Msg("black"))
	case White:
		return fmt.Sprintf(k.Msg("stone"), k.Msg("white"))
	}
	return k.Msg("empty")
}
