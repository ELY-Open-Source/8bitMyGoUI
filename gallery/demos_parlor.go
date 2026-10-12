package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/tabletop"
)

var parlorEntries = []Entry{
	{494, "Domino", "多米诺", "tabletop", "Two halves, zero to six pips. Lay them flat or on end.", demoDomino},
	{495, "Roulette", "轮盘", "tabletop", "SPIN and the wheel runs a lap before it stops. The house is green.", demoRoulette},
	{496, "CoinFlip", "抛硬币", "tabletop", "Heads or tails, turned in the air frame by frame.", demoCoinFlip},
	{497, "BingoCard", "宾果卡", "tabletop", "Mark your numbers; five in a line shouts BINGO!", demoBingoCard},
	{498, "TicTacToe", "井字棋", "tabletop", "You are X. The computer is a decent O, not a genius.", demoTicTacToe},
	{499, "ScoreTrack", "计分轨道", "tabletop", "A winding track of spaces; each player's meeple marks the score.", demoScoreTrack},
	{500, "Meeple", "棋子", "tabletop", "A player's token in its slot. Number five hundred: thank you for scrolling!", demoMeeple},
}

type parlorDemo struct {
	spins, result, flips int
	heads                bool
	marked               [25]bool
	scores               []int
	turn                 int
}

// prizes ring the demo wheel; the house takes the first.
var prizes = []string{"HOUSE", "50", "100", "1UP", "200", "LOSE", "500", "x2", "300", "BOMB", "1000", "?"}

func parlorState(c *ui.Context) *parlorDemo {
	return ui.Local(c.Root(), demoKey("parlor"), func() parlorDemo {
		return parlorDemo{result: 3, heads: true, scores: []int{12, 7, 12}}
	})
}

func demoDomino(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		tabletop.Domino(c, 6, 6, tabletop.DominoOptions{})
		tabletop.Domino(c, 3, 5, tabletop.DominoOptions{})
		tabletop.Domino(c, 0, 1, tabletop.DominoOptions{Vertical: true})
		tabletop.Domino(c, 2, 4, tabletop.DominoOptions{Vertical: true})
	})
}

func demoRoulette(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := parlorState(c)
	rng := tableState(c).rng
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		tabletop.Roulette(c, "Wheel", prizes, s.result, tabletop.RouletteOptions{Cue: s.spins})
		if input.Button(c, "SPIN", input.ButtonOptions{Size: input.Small}).Clicked {
			s.result, s.spins = rng.IntN(len(prizes)), s.spins+1
			log("landed on " + prizes[s.result])
		}
	})
}

func demoCoinFlip(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := parlorState(c)
	rng := tableState(c).rng
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		tabletop.CoinFlip(c, "Coin", s.heads, tabletop.CoinFlipOptions{Cue: s.flips})
		if input.Button(c, "FLIP", input.ButtonOptions{Size: input.Small}).Clicked {
			s.heads, s.flips = rng.IntN(2) == 0, s.flips+1
			log(fmt.Sprintf("heads %v", s.heads))
		}
	})
}

func demoBingoCard(c *ui.Context, log func(string)) {
	s := parlorState(c)
	numbers := [25]int{}
	for i := range numbers {
		numbers[i] = i%5*15 + 1 + (i/5*7)%15
	}
	r := tabletop.BingoCard(c, "Card", numbers, &s.marked, tabletop.BingoCardOptions{})
	if r.Changed {
		log("marked")
	}
	if r.Bingo {
		log("BINGO!")
	}
}

func demoTicTacToe(c *ui.Context, log func(string)) {
	r := tabletop.TicTacToe(c, "Grid", tabletop.TicTacToeOptions{})
	if r.Played >= 0 {
		log(fmt.Sprintf("X takes %d", r.Played))
	}
	if r.Over {
		log([]string{"a draw", "you won", "the computer won"}[r.Winner])
	}
	if r.Restarted {
		log("a new game")
	}
}

func demoScoreTrack(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := parlorState(c)
	rng := tableState(c).rng
	ui.Column(c).Gap(k.D(6)).Children(func() {
		tabletop.ScoreTrack(c, "Scores", s.scores, tabletop.ScoreTrackOptions{})
		if input.Button(c, "TAKE A TURN", input.ButtonOptions{Size: input.Small}).Clicked {
			gain := rng.IntN(6) + 1
			s.scores[s.turn] += gain
			log(fmt.Sprintf("player %d +%d", s.turn+1, gain))
			s.turn = (s.turn + 1) % len(s.scores)
		}
	})
}

func demoMeeple(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(4)).Children(func() {
		for i := range 6 {
			tabletop.Meeple(c, fmt.Sprintf("Player %d", i+1), tabletop.MeepleOptions{Player: i})
		}
		tabletop.Meeple(c, "Open seat", tabletop.MeepleOptions{Empty: true})
	})
}
