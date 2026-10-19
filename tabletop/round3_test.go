package tabletop_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/tabletop"
)

func TestBingoCardDrawsInsideItsGrid(t *testing.T) {
	var numbers [25]int
	for i := range numbers {
		numbers[i] = i + 1
	}
	var m [25]bool
	tt := view(400, 400, newWorld(), func(c *ui.Context) { tabletop.BingoCard(c, "B", numbers, &m, tabletop.BingoCardOptions{}) })
	x, y, w, h := dots(find(t, tt, "B"))
	// The letters sit nine dots above the grid.
	if n := spill(tt.Image(), x, y-9, w, h+9); n > 0 {
		t.Fatalf("the card spills %d dots past its grid", n)
	}
}

func TestCardHandPlaysTheTopOfAnOverlap(t *testing.T) {
	cards := []tabletop.Card{{Rank: 2, Suit: tabletop.Clubs}, {Rank: 3, Suit: tabletop.Clubs}}
	played := -1
	tt := view(400, 300, newWorld(), func(c *ui.Context) {
		if i := tabletop.CardHand(c, "Hand", cards, tabletop.CardHandOptions{}).Played; i >= 0 {
			played = i
		}
	})
	r := find(t, tt, "Hand")
	// Dots 18–35 hold both cards; the second is on top.
	tt.ClickAt(r.X+2*24, r.Y+2*30)
	if played != 1 {
		t.Fatalf("a click where the cards overlap played %d", played)
	}
}
