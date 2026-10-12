package gallery

import (
	"fmt"
	"math/rand/v2"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/tabletop"
)

var tabletopEntries = []Entry{
	{486, "Dice", "骰子", "tabletop", "A d6 with pips. ROLL and it tumbles before it lands.", demoDice},
	{487, "DiceTray", "骰盘", "tabletop", "Several dice on felt, and what they add up to.", demoDiceTray},
	{488, "PlayingCard", "扑克牌", "tabletop", "Rank and suit, court cards with their pieces, or the back.", demoPlayingCard},
	{489, "CardHand", "手牌", "tabletop", "A fanned hand. The card you point at or choose lifts; a click plays it.", demoCardHand},
	{490, "Deck", "牌堆", "tabletop", "A face-down pile with its count. Click to draw.", demoDeck},
	{491, "Chessboard", "棋盘", "tabletop", "Click a piece, then a square. The rules are on the honor system.", demoChessboard},
	{492, "GoBoard", "围棋盘", "tabletop", "Black, then White, a stone a click. 9, 13 or 19 lines.", demoGoBoard},
	{493, "MahjongTile", "麻将牌", "tabletop", "Dots, bamboo, characters, the four winds and three dragons.", demoMahjongTile},
}

type tableDemo struct {
	rng        *rand.Rand
	die, rolls int
	tray       []int
	hand       []tabletop.Card
	left       int
	drawn      tabletop.Card
	board      [64]tabletop.Piece
	stones     []tabletop.Stone
}

func tableState(c *ui.Context) *tableDemo {
	return ui.Local(c.Root(), demoKey("tabletop"), func() tableDemo {
		return tableDemo{rng: rand.New(rand.NewPCG(8, 16)), die: 6, tray: []int{3, 5, 2}, hand: deal(), left: 52,
			board: chessStart(), stones: make([]tabletop.Stone, 81)}
	})
}

// deal is the demo's five-card hand.
func deal() []tabletop.Card {
	return []tabletop.Card{{Rank: 1, Suit: tabletop.Spades}, {Rank: 10, Suit: tabletop.Hearts}, {Rank: 7, Suit: tabletop.Clubs},
		{Rank: 12, Suit: tabletop.Diamonds}, {Rank: 13, Suit: tabletop.Hearts}}
}

// chessStart is the opening position, Black at the top.
func chessStart() (b [64]tabletop.Piece) {
	back := []tabletop.Kind{tabletop.Rook, tabletop.Knight, tabletop.Bishop, tabletop.Queen, tabletop.King, tabletop.Bishop, tabletop.Knight, tabletop.Rook}
	for i, kind := range back {
		b[i], b[56+i] = tabletop.Piece{Kind: kind}, tabletop.Piece{Kind: kind, White: true}
		b[8+i], b[48+i] = tabletop.Piece{Kind: tabletop.Pawn}, tabletop.Piece{Kind: tabletop.Pawn, White: true}
	}
	return b
}

func demoDice(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := tableState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		tabletop.Dice(c, "Die", s.die, tabletop.DiceOptions{Cue: s.rolls})
		if input.Button(c, "ROLL", input.ButtonOptions{Size: input.Small}).Clicked {
			s.die, s.rolls = s.rng.IntN(6)+1, s.rolls+1
			log(fmt.Sprintf("rolled a %d", s.die))
		}
	})
}

func demoDiceTray(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := tableState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		tabletop.DiceTray(c, "Dice", s.tray, tabletop.DiceTrayOptions{Cue: s.rolls})
		if input.Button(c, "ROLL ALL", input.ButtonOptions{Size: input.Small}).Clicked {
			total := 0
			for i := range s.tray {
				s.tray[i] = s.rng.IntN(6) + 1
				total += s.tray[i]
			}
			s.rolls++
			log(fmt.Sprintf("rolled %d", total))
		}
	})
}

func demoPlayingCard(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(4)).Children(func() {
		for _, cd := range deal()[:4] {
			tabletop.PlayingCard(c, cd, tabletop.PlayingCardOptions{})
		}
		tabletop.PlayingCard(c, tabletop.Card{Rank: 1}, tabletop.PlayingCardOptions{Back: true})
	})
}

func demoCardHand(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := tableState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		if i := tabletop.CardHand(c, "Hand", s.hand, tabletop.CardHandOptions{}).Played; i >= 0 {
			log(fmt.Sprintf("played the %d of suit %d", s.hand[i].Rank, s.hand[i].Suit))
			s.hand = append(s.hand[:i:i], s.hand[i+1:]...)
		}
		if input.Button(c, "DEAL", input.ButtonOptions{Size: input.Small}).Clicked {
			s.hand = deal()
			log("a fresh hand")
		}
	})
}

func demoDeck(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := tableState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.End).Children(func() {
		if tabletop.Deck(c, "Draw pile", s.left, tabletop.DeckOptions{}).Drawn {
			s.left--
			s.drawn = tabletop.Card{Rank: s.rng.IntN(13) + 1, Suit: tabletop.Suit(s.rng.IntN(4))}
			log(fmt.Sprintf("drew a card, %d left", s.left))
		}
		if s.drawn.Rank > 0 {
			tabletop.PlayingCard(c, s.drawn, tabletop.PlayingCardOptions{})
		}
	})
}

func demoChessboard(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := tableState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		r := tabletop.Chessboard(c, "Board", &s.board, tabletop.ChessboardOptions{})
		if r.Changed {
			log(fmt.Sprintf("holding square %d", r.Held))
		}
		if r.Moved {
			log(fmt.Sprintf("moved from square %d to %d", r.From, r.To))
		}
		if input.Button(c, "NEW GAME", input.ButtonOptions{Size: input.Small}).Clicked {
			s.board = chessStart()
			log("pieces reset")
		}
	})
}

func demoGoBoard(c *ui.Context, log func(string)) {
	s := tableState(c)
	if i := tabletop.GoBoard(c, "Goban", s.stones, tabletop.GoBoardOptions{}).Placed; i >= 0 {
		log(fmt.Sprintf("stone at %d", i))
	}
}

func demoMahjongTile(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	tiles := []tabletop.Tile{
		{Family: tabletop.Dots, Rank: 1}, {Family: tabletop.Dots, Rank: 5}, {Family: tabletop.Dots, Rank: 9},
		{Family: tabletop.Bamboo, Rank: 3}, {Family: tabletop.Bamboo, Rank: 8}, {Family: tabletop.Characters, Rank: 7},
		{Family: tabletop.Winds, Rank: 1}, {Family: tabletop.Winds, Rank: 4}, {Family: tabletop.Dragons, Rank: 1},
		{Family: tabletop.Dragons, Rank: 2}, {Family: tabletop.Dragons, Rank: 3},
	}
	ui.Row(c).Gap(k.D(2)).Wrap().Width(k.D(300)).Children(func() {
		for _, t := range tiles {
			tabletop.MahjongTile(c, t, tabletop.MahjongTileOptions{})
		}
		tabletop.MahjongTile(c, tiles[0], tabletop.MahjongTileOptions{Back: true})
	})
}
