package tabletop_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/tabletop"
)

func TestDiceTumbleThenLand(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = reduce
		value, cue := 3, 0
		tt := view(200, 200, wd, func(c *ui.Context) { tabletop.Dice(c, "Die", value, tabletop.DiceOptions{Cue: cue}) })
		landed := func(v int) []byte {
			return view(200, 200, newWorld(), func(c *ui.Context) { tabletop.Dice(c, "Die", v, tabletop.DiceOptions{}) }).Image().Pix
		}
		if !bytes.Equal(tt.Image().Pix, landed(3)) {
			t.Fatal("a die shown first tumbled")
		}
		value, cue = 5, 1
		tt.Frame()
		if tumbling := !bytes.Equal(tt.Image().Pix, landed(5)); tumbling == reduce {
			t.Errorf("reduce %v: tumbling is %v", reduce, tumbling)
		}
		later(tt, wd, time.Second)
		if !bytes.Equal(tt.Image().Pix, landed(5)) {
			t.Errorf("reduce %v: it never landed on 5", reduce)
		}
	}
}

func TestDiceTrayAddsUp(t *testing.T) {
	tt := view(400, 200, newWorld(), func(c *ui.Context) { tabletop.DiceTray(c, "Dice", []int{6, 6, 1}, tabletop.DiceTrayOptions{}) })
	find(t, tt, "TOTAL 13")
	find(t, tt, "Dice 3")
}

func TestCardsReadAndInk(t *testing.T) {
	tt := view(400, 400, newWorld(), func(c *ui.Context) {
		tabletop.PlayingCard(c, tabletop.Card{Rank: 1, Suit: tabletop.Spades}, tabletop.PlayingCardOptions{})
		tabletop.PlayingCard(c, tabletop.Card{Rank: 12, Suit: tabletop.Hearts}, tabletop.PlayingCardOptions{})
		tabletop.PlayingCard(c, tabletop.Card{Rank: 9, Suit: tabletop.Clubs}, tabletop.PlayingCardOptions{Back: true})
	})
	img := tt.Image()
	x, y, w, h := dots(find(t, tt, "ACE OF SPADES"))
	if count(img, x, y, w, h, pixel.Bad) > 0 {
		t.Error("a spade is red")
	}
	x, y, w, h = dots(find(t, tt, "QUEEN OF HEARTS"))
	if count(img, x, y, w, h, pixel.Bad) == 0 {
		t.Error("a heart is not red")
	}
	x, y, w, h = dots(find(t, tt, "CARD BACK"))
	if count(img, x, y, w, h, pixel.Accent) == 0 {
		t.Error("the back shows no pattern")
	}
}

func TestCardHandLiftsAndPlays(t *testing.T) {
	cards := []tabletop.Card{{Rank: 2, Suit: tabletop.Clubs}, {Rank: 3, Suit: tabletop.Clubs}, {Rank: 4, Suit: tabletop.Clubs}}
	wd, played, off := newWorld(), -1, false
	tt := view(400, 300, wd, func(c *ui.Context) {
		if i := tabletop.CardHand(c, "Hand", cards, tabletop.CardHandOptions{Disabled: off}).Played; i >= 0 {
			played = i
		}
	})
	x, y, _, _ := dots(find(t, tt, "Hand"))
	// top is the first row of a card's outline.
	top := func(i int) int {
		for j := y; j < y+60; j++ {
			if at(tt.Image(), x+i*18+4, j) == pixel.Edge {
				return j - y
			}
		}
		return -1
	}
	rest := top(1)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	later(tt, wd, time.Second)
	if top(1) != rest-6 || top(0) != rest {
		t.Fatalf("the chosen card rose to %d from %d", top(1), rest)
	}
	tt.Key(0, ui.KeyEnter)
	if played != 1 {
		t.Fatalf("Enter played %d", played)
	}
	r := find(t, tt, "Hand")
	tt.ClickAt(r.X+r.W-4, r.Y+r.H/2)
	if played != 2 {
		t.Fatalf("a click on the last card played %d", played)
	}
	off, played = true, -1
	tt.Key(0, ui.KeyEnter)
	if played != -1 {
		t.Fatal("a disabled hand played")
	}
}

func TestDeckDrawsUntilEmpty(t *testing.T) {
	left, drawn := 2, 0
	tt := view(300, 300, newWorld(), func(c *ui.Context) {
		if tabletop.Deck(c, "Deck", left, tabletop.DeckOptions{}).Drawn {
			left--
			drawn++
		}
	})
	// The count's badge is the only AccentInk on the pile.
	x, y, w, h := dots(find(t, tt, "Deck"))
	if count(tt.Image(), x, y, w, h, pixel.AccentInk) == 0 {
		t.Fatal("no count badge")
	}
	click(t, tt, "Deck")
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyEnter)
	click(t, tt, "Deck")
	if left != 0 || drawn != 2 {
		t.Fatalf("drew %d, %d left", drawn, left)
	}
}

func TestDominoAndMahjongRead(t *testing.T) {
	tt := view(600, 600, newWorld(), func(c *ui.Context) {
		tabletop.Domino(c, 3, 5, tabletop.DominoOptions{Vertical: true})
		tabletop.MahjongTile(c, tabletop.Tile{Family: tabletop.Dots, Rank: 3}, tabletop.MahjongTileOptions{})
		tabletop.MahjongTile(c, tabletop.Tile{Family: tabletop.Winds, Rank: 1}, tabletop.MahjongTileOptions{})
		tabletop.MahjongTile(c, tabletop.Tile{Family: tabletop.Dragons, Rank: 1}, tabletop.MahjongTileOptions{})
		tabletop.MahjongTile(c, tabletop.Tile{Family: tabletop.Bamboo, Rank: 2}, tabletop.MahjongTileOptions{Back: true})
	})
	if _, _, w, h := dots(find(t, tt, "3 | 5")); w != 20 || h != 39 {
		t.Fatalf("a standing domino is %d×%d", w, h)
	}
	for _, label := range []string{"3 DOTS", "EAST WIND", "TILE BACK"} {
		find(t, tt, label)
	}
	x, y, w, h := dots(find(t, tt, "RED DRAGON"))
	if count(tt.Image(), x, y, w, h, pixel.Bad) == 0 {
		t.Fatal("the red dragon is not red")
	}
}

func TestDiceShowPips(t *testing.T) {
	for v := 2; v <= 6; v++ {
		tt := view(200, 200, newWorld(), func(c *ui.Context) { tabletop.Dice(c, "Die", v, tabletop.DiceOptions{}) })
		x, y, w, h := dots(find(t, tt, "Die"))
		// Each pip is two dots square, in Ink.
		if n := count(tt.Image(), x, y, w, h, pixel.Ink); n != 4*v {
			t.Errorf("a %d shows %d Ink dots", v, n)
		}
	}
}

func TestWhitePiecesLookLight(t *testing.T) {
	var board [64]tabletop.Piece
	board[1] = tabletop.Piece{Kind: tabletop.Pawn, White: true}
	board[3] = tabletop.Piece{Kind: tabletop.Pawn}
	tt := view(600, 400, newWorld(), func(c *ui.Context) { tabletop.Chessboard(c, "B", &board, tabletop.ChessboardOptions{}) })
	x, y, _, _ := dots(find(t, tt, "CHESSBOARD"))
	// On dark squares, a white pawn takes the lighter slot.
	white, black := count(tt.Image(), x+16, y, 16, 16, pixel.Lit), count(tt.Image(), x+48, y, 16, 16, pixel.Lit)
	if white <= black {
		t.Fatalf("the white pawn has %d Lit dots, the black %d", white, black)
	}
}

func TestEveryTileAndDominoLooksItsOwn(t *testing.T) {
	seen := map[string]tabletop.Tile{}
	for f, ranks := range []int{9, 9, 9, 4, 3} {
		for r := 1; r <= ranks; r++ {
			tile := tabletop.Tile{Family: tabletop.Family(f), Rank: r}
			pix := string(view(100, 100, newWorld(), func(c *ui.Context) { tabletop.MahjongTile(c, tile, tabletop.MahjongTileOptions{}) }).Image().Pix)
			if other, ok := seen[pix]; ok {
				t.Errorf("%+v looks like %+v", tile, other)
			}
			seen[pix] = tile
		}
	}
	ink := func(a, b int) int {
		tt := view(200, 100, newWorld(), func(c *ui.Context) { tabletop.Domino(c, a, b, tabletop.DominoOptions{}) })
		x, y, w, h := dots(find(t, tt, fmt.Sprintf("%d | %d", a, b)))
		return count(tt.Image(), x, y, w, h, pixel.Ink)
	}
	// Two dots square a pip, as on dice.
	if blank, two, six := ink(0, 0), ink(1, 1), ink(6, 6); blank != 0 || two != 8 || six != 48 {
		t.Errorf("0|0, 1|1 and 6|6 show %d, %d and %d Ink dots", blank, two, six)
	}
}
