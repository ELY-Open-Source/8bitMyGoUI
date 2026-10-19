package tabletop_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/tabletop"
)

func TestBadInputPanics(t *testing.T) {
	n := 0
	for name, build := range map[string]func(c *ui.Context){
		"die 0":   func(c *ui.Context) { tabletop.Dice(c, "D", 0, tabletop.DiceOptions{}) },
		"die 7":   func(c *ui.Context) { tabletop.Dice(c, "D", 7, tabletop.DiceOptions{}) },
		"no dice": func(c *ui.Context) { tabletop.DiceTray(c, "T", nil, tabletop.DiceTrayOptions{}) },
		"rank 14": func(c *ui.Context) { tabletop.PlayingCard(c, tabletop.Card{Rank: 14}, tabletop.PlayingCardOptions{}) },
		"suit 4": func(c *ui.Context) {
			tabletop.PlayingCard(c, tabletop.Card{Rank: 1, Suit: 4}, tabletop.PlayingCardOptions{})
		},
		"twin cards": func(c *ui.Context) {
			tabletop.CardHand(c, "H", []tabletop.Card{{Rank: 1}, {Rank: 1}}, tabletop.CardHandOptions{})
		},
		"minus deck": func(c *ui.Context) { tabletop.Deck(c, "D", -1, tabletop.DeckOptions{}) },
		"kind 9": func(c *ui.Context) {
			var b [64]tabletop.Piece
			b[0].Kind = 9
			tabletop.Chessboard(c, "B", &b, tabletop.ChessboardOptions{})
		},
		"go 10x10": func(c *ui.Context) { tabletop.GoBoard(c, "G", make([]tabletop.Stone, 100), tabletop.GoBoardOptions{}) },
		"stone 3": func(c *ui.Context) {
			tabletop.GoBoard(c, "G", append(make([]tabletop.Stone, 80), 3), tabletop.GoBoardOptions{})
		},
		"dots 10": func(c *ui.Context) { tabletop.MahjongTile(c, tabletop.Tile{Rank: 10}, tabletop.MahjongTileOptions{}) },
		"wind 5": func(c *ui.Context) {
			tabletop.MahjongTile(c, tabletop.Tile{Family: tabletop.Winds, Rank: 5}, tabletop.MahjongTileOptions{})
		},
		"domino 7":  func(c *ui.Context) { tabletop.Domino(c, 7, 0, tabletop.DominoOptions{}) },
		"one wedge": func(c *ui.Context) { tabletop.Roulette(c, "R", []string{"A"}, 0, tabletop.RouletteOptions{}) },
		"past end":  func(c *ui.Context) { tabletop.Roulette(c, "R", []string{"A", "B"}, 2, tabletop.RouletteOptions{}) },
		"blank":     func(c *ui.Context) { tabletop.Roulette(c, "R", []string{"A", ""}, 0, tabletop.RouletteOptions{}) },
		"no players": func(c *ui.Context) {
			tabletop.ScoreTrack(c, "S", nil, tabletop.ScoreTrackOptions{})
		},
		"7 players":   func(c *ui.Context) { tabletop.ScoreTrack(c, "S", make([]int, 7), tabletop.ScoreTrackOptions{}) },
		"minus score": func(c *ui.Context) { tabletop.ScoreTrack(c, "S", []int{-1}, tabletop.ScoreTrackOptions{}) },
		"player 6":    func(c *ui.Context) { tabletop.Meeple(c, "M", tabletop.MeepleOptions{Player: 6}) },
	} {
		n++
		// The panic must be the package's own, not a fault.
		func() {
			defer func() {
				if r := fmt.Sprint(recover()); !strings.HasPrefix(r, "tabletop:") {
					t.Errorf("%s panicked with %q", name, r)
				}
			}()
			view(600, 300, newWorld(), build)
		}()
	}
	if n != 20 {
		t.Fatalf("%d cases", n)
	}
}

func TestExtremesStayInside(t *testing.T) {
	long := strings.Repeat("W", 40)
	for _, tc := range []struct {
		label string
		build func(c *ui.Context)
	}{
		{"R", func(c *ui.Context) { tabletop.Roulette(c, "R", []string{long, "B"}, 0, tabletop.RouletteOptions{}) }},
		{"S", func(c *ui.Context) {
			tabletop.ScoreTrack(c, "S", []int{9, 9, 9, 9, 9, 9}, tabletop.ScoreTrackOptions{Spaces: 10})
		}},
		{"T", func(c *ui.Context) {
			tabletop.DiceTray(c, "T", []int{1, 2, 3, 4, 5, 6, 6, 6}, tabletop.DiceTrayOptions{})
		}},
	} {
		tt := view(1400, 700, newWorld(), tc.build)
		x, y, w, h := dots(find(t, tt, tc.label))
		if n := spill(tt.Image(), x, y, w, h); n > 0 {
			t.Errorf("%s spills %d dots", tc.label, n)
		}
	}
}
