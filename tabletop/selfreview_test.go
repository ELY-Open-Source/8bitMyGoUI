package tabletop_test

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/tabletop"
)

func TestTossesTakeAClockStepBack(t *testing.T) {
	wd := newWorld()
	cue := 0
	tt := view(600, 400, wd, func(c *ui.Context) {
		tabletop.Dice(c, "Die", 4, tabletop.DiceOptions{Cue: cue})
		tabletop.CoinFlip(c, "Coin", true, tabletop.CoinFlipOptions{Cue: cue})
		tabletop.Roulette(c, "Wheel", []string{"A", "B", "C"}, 1, tabletop.RouletteOptions{Cue: cue})
	})
	cue = 1
	tt.Frame()
	later(tt, wd, -time.Second)
	later(tt, wd, 3*time.Second)
	find(t, tt, "HEADS")
}

func TestTicTacToeReportsEachMove(t *testing.T) {
	played := []int{}
	log := logged(func() {
		tt := view(400, 400, newWorld(), func(c *ui.Context) {
			if r := tabletop.TicTacToe(c, "Grid", tabletop.TicTacToeOptions{}); r.Played >= 0 {
				played = append(played, r.Played)
			}
		})
		cell(t, tt, "Grid", 3, 22, 0)
	})
	if len(played) != 1 || played[0] != 0 || !strings.Contains(log, "tic-tac-toe move") {
		t.Fatalf("played %v; log:\n%s", played, log)
	}
}

func TestTooWideNumbersPanic(t *testing.T) {
	for name, build := range map[string]func(c *ui.Context){
		"track": func(c *ui.Context) { tabletop.ScoreTrack(c, "S", []int{1}, tabletop.ScoreTrackOptions{Spaces: 1000}) },
		"bingo": func(c *ui.Context) {
			var n [25]int
			var m [25]bool
			for i := range n {
				n[i] = 1000 + i
			}
			tabletop.BingoCard(c, "B", n, &m, tabletop.BingoCardOptions{})
		},
	} {
		mustPanic(t, name, func() { view(600, 400, newWorld(), build) })
	}
}
