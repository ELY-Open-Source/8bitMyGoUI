package tabletop

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Family is a mahjong suit, or the winds or dragons.
type Family int

const (
	Dots Family = iota
	Bamboo
	Characters
	Winds
	Dragons
)

var families = [...]struct {
	msg   string
	ranks int
}{{"dots_suit", 9}, {"bamboo", 9}, {"characters", 9}, {"", 4}, {"", 3}}

// Tile: suits rank 1–9, winds 1–4, dragons 1–3.
type Tile struct {
	Family Family
	Rank   int
}

type MahjongTileOptions struct {
	// Back shows the tile face down.
	Back bool
}

// spots are where 1–9 dots or sticks stand, on 3×3.
var spots = [...][][2]int{{}, {{1, 1}}, {{1, 0}, {1, 2}}, {{0, 0}, {1, 1}, {2, 2}}, {{0, 0}, {2, 0}, {0, 2}, {2, 2}},
	{{0, 0}, {2, 0}, {1, 1}, {0, 2}, {2, 2}}, {{0, 0}, {2, 0}, {0, 1}, {2, 1}, {0, 2}, {2, 2}},
	{{0, 0}, {2, 0}, {0, 1}, {1, 1}, {2, 1}, {0, 2}, {2, 2}}, {{0, 0}, {1, 0}, {2, 0}, {0, 1}, {2, 1}, {0, 2}, {1, 2}, {2, 2}},
	{{0, 0}, {1, 0}, {2, 0}, {0, 1}, {1, 1}, {2, 1}, {0, 2}, {1, 2}, {2, 2}}}

// Tile faces are art: their hanzi stay, whatever the locale.
var (
	numerals = []string{"一", "二", "三", "四", "五", "六", "七", "八", "九"}
	winds    = []string{"east", "south", "west", "north"}
	dragons  = []string{"red_dragon", "green_dragon", "white_dragon"}
)

const tileW, tileH = 26, 34

// MahjongTile is a suit or honor tile, or its back.
func MahjongTile(c *ui.Context, t Tile, o MahjongTileOptions) ui.Element {
	k := pixel.K(c)
	if t.Family < Dots || t.Family > Dragons || t.Rank < 1 || t.Rank > families[t.Family].ranks {
		panic(fmt.Sprintf("tabletop: MahjongTile %+v", t))
	}
	name := k.Msg("tile_back")
	switch {
	case o.Back:
	case t.Family == Winds:
		name = k.Msg(winds[t.Rank-1])
	case t.Family == Dragons:
		name = k.Msg(dragons[t.Rank-1])
	default:
		name = fmt.Sprintf(k.Msg("tile_name"), t.Rank, k.Msg(families[t.Family].msg))
	}
	return ui.Box(c).Size(k.D(tileW), k.D(tileH)).Role(ui.RoleImage).Label(name).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		// The green back shows as a rim below, right.
		g.Frame(2, 2, tileW-2, tileH-2, pixel.Flat, pixel.Good)
		face := pixel.Lit
		if o.Back {
			face = pixel.Good
		}
		g.Frame(0, 0, tileW-2, tileH-2, pixel.Flat, face)
		if !o.Back {
			tileFace(g, t)
		}
	})
}

// tileFace paints a tile's face in its 24×32 front.
func tileFace(g pixel.Canvas, t Tile) {
	big := func(s string, ink pixel.Slot) { g.Text((24-2*pixel.Measure(s))/2, 2, s, ink, 2) }
	switch t.Family {
	case Dots:
		for i, s := range spots[t.Rank] {
			ink := []pixel.Slot{pixel.Info, pixel.Good, pixel.Bad}[i%3]
			r := 2
			if t.Rank == 1 {
				r = 6
			}
			g.Circle(5+s[0]*7, 7+s[1]*9, r, ink)
			g.Dot(5+s[0]*7, 7+s[1]*9, ink)
		}
	case Bamboo:
		for _, s := range spots[t.Rank] {
			x, y := 4+s[0]*7, 3+s[1]*9
			g.Rect(x, y, 2, 8, pixel.Good)
			g.Rect(x, y+3, 2, 1, pixel.Lit)
		}
	case Characters:
		n := numerals[t.Rank-1]
		g.Text((24-pixel.Measure(n))/2, 1, n, pixel.Ink, 1)
		g.Text((24-pixel.Measure("万"))/2, 15, "万", pixel.Bad, 1)
	case Winds:
		big([]string{"东", "南", "西", "北"}[t.Rank-1], pixel.Ink)
	case Dragons:
		switch t.Rank {
		case 1:
			big("中", pixel.Bad)
		case 2:
			big("发", pixel.Good)
		default:
			g.Rect(5, 6, 14, 20, pixel.Info)
			g.Rect(7, 8, 10, 16, pixel.Lit)
		}
	}
}
