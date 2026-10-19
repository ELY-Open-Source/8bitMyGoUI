package tabletop_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/tabletop"
)

func TestGoBoardShrinksUnderItsCursor(t *testing.T) {
	board := make([]tabletop.Stone, 361)
	tt := view(600, 600, newWorld(), func(c *ui.Context) { tabletop.GoBoard(c, "Goban", board, tabletop.GoBoardOptions{}) })
	board = make([]tabletop.Stone, 81)
	tt.Frame()
	find(t, tt, "Goban")
}

func TestCourtCardsShowTheirSuit(t *testing.T) {
	shot := func(cd tabletop.Card) []byte {
		return view(200, 200, newWorld(), func(c *ui.Context) { tabletop.PlayingCard(c, cd, tabletop.PlayingCardOptions{}) }).Image().Pix
	}
	for _, rank := range []int{11, 12, 13} {
		if bytes.Equal(shot(tabletop.Card{Rank: rank, Suit: tabletop.Hearts}), shot(tabletop.Card{Rank: rank, Suit: tabletop.Diamonds})) ||
			bytes.Equal(shot(tabletop.Card{Rank: rank, Suit: tabletop.Spades}), shot(tabletop.Card{Rank: rank, Suit: tabletop.Clubs})) {
			t.Errorf("rank %d hides its suit", rank)
		}
	}
}

func TestCardHandHitsTheCardDrawn(t *testing.T) {
	cards := []tabletop.Card{{Rank: 2, Suit: tabletop.Clubs}, {Rank: 3, Suit: tabletop.Clubs}}
	wd, played := newWorld(), -1
	tt := view(400, 300, wd, func(c *ui.Context) {
		if i := tabletop.CardHand(c, "Hand", cards, tabletop.CardHandOptions{}).Played; i >= 0 {
			played = i
		}
	})
	r := find(t, tt, "Hand")
	// Above the cards is felt, not a card.
	tt.ClickAt(r.X+44, r.Y+2)
	if played != -1 {
		t.Fatalf("a click above the cards played %d", played)
	}
	// Lifted, the first card shows above; clicking there takes it.
	tt.Move(r.X+10, r.Y+40)
	later(tt, wd, 1e9)
	tt.ClickAt(r.X+44, r.Y+8)
	if played != 0 {
		t.Fatalf("a click on the lifted first card played %d", played)
	}
}

func TestDeckAndBingoStayInTheirCells(t *testing.T) {
	mustPanic(t, "deck", func() {
		view(300, 300, newWorld(), func(c *ui.Context) { tabletop.Deck(c, "D", 1000, tabletop.DeckOptions{}) })
	})
	var numbers [25]int
	for i := range numbers {
		numbers[i] = i + 1
	}
	card := func(l pixel.Locale) []byte {
		var m [25]bool
		return ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Locale: l})
			tabletop.BingoCard(c, "B", numbers, &m, tabletop.BingoCardOptions{})
		}, 400, 400).Image().Pix
	}
	// FREE is a star, alike in every locale.
	if !bytes.Equal(card(pixel.En), card(pixel.Zh)) {
		t.Fatal("the FREE cell differs by locale")
	}
	var m [25]bool
	tt := view(400, 400, newWorld(), func(c *ui.Context) { tabletop.BingoCard(c, "B", numbers, &m, tabletop.BingoCardOptions{}) })
	x, y, _, _ := dots(find(t, tt, "B"))
	if count(tt.Image(), x+36, y+32, 18, 16, pixel.Accent) == 0 {
		t.Fatal("no star on FREE")
	}
}

func TestMeepleStandsInTheItemSlot(t *testing.T) {
	tt := view(200, 200, newWorld(), func(c *ui.Context) { tabletop.Meeple(c, "M", tabletop.MeepleOptions{}) })
	x, y, w, h := dots(find(t, tt, "M"))
	if w != data.SlotSize || h != data.SlotSize || count(tt.Image(), x, y, w, h, pixel.Well) == 0 {
		t.Fatalf("a meeple slot is %d×%d, or not an item slot", w, h)
	}
}

func TestHoldsAndRestartsAreLogged(t *testing.T) {
	var board [64]tabletop.Piece
	board[0] = tabletop.Piece{Kind: tabletop.Rook}
	restarts := 0
	log := logged(func() {
		tt := view(600, 700, newWorld(), func(c *ui.Context) {
			tabletop.Chessboard(c, "B", &board, tabletop.ChessboardOptions{})
			if tabletop.TicTacToe(c, "Grid", tabletop.TicTacToeOptions{}).Restarted {
				restarts++
			}
		})
		cell(t, tt, "CHESSBOARD", 8, 16, 0)
		for _, i := range []int{0, 8, 1, 4} {
			cell(t, tt, "Grid", 3, 22, i)
		}
	})
	if !strings.Contains(log, "chess hold") || !strings.Contains(log, "tic-tac-toe restart") || restarts != 1 {
		t.Fatalf("%d restarts; log:\n%s", restarts, log)
	}
}
