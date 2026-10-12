// Package tabletop is dice, cards, boards and tokens.
package tabletop

import (
	"fmt"
	"strconv"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Suit is a playing card's suit.
type Suit int

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

var suits = [...]struct {
	msg string
	sp  *pixel.Sprite
}{{"spades", icon.Spade}, {"hearts", icon.Heart}, {"diamonds", icon.Diamond}, {"clubs", icon.Club}}

// Card's Rank runs from ace, 1, to king, 13.
type Card struct {
	Rank int
	Suit Suit
}

// check fails on a rank or suit no deck has.
func (cd Card) check(name string) {
	if cd.Rank < 1 || cd.Rank > 13 || cd.Suit < Spades || cd.Suit > Clubs {
		panic(fmt.Sprintf("tabletop: %s card %+v", name, cd))
	}
}

// ink is a card's color: hearts and diamonds are red.
func (cd Card) ink() pixel.Slot {
	if cd.Suit == Hearts || cd.Suit == Diamonds {
		return pixel.Bad
	}
	return pixel.Ink
}

// index is the printed rank: A, 2–10, J, Q, K.
func (cd Card) index() string {
	return [...]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}[cd.Rank-1]
}

// name reads a card out, as ACE OF SPADES.
func (cd Card) name(k pixel.Kit) string {
	rank := strconv.Itoa(cd.Rank)
	switch cd.Rank {
	case 1:
		rank = k.Msg("ace")
	case 11, 12, 13:
		rank = k.Msg([...]string{"jack", "queen", "king"}[cd.Rank-11])
	}
	return fmt.Sprintf(k.Msg("card_name"), rank, k.Msg(suits[cd.Suit].msg))
}

const cardW, cardH = 36, 50

// drawCard paints a card at x, y: face or back.
func drawCard(g pixel.Canvas, x, y int, cd Card, back bool) {
	g.Frame(x, y, cardW, cardH, pixel.Flat, pixel.Lit)
	if back {
		g.Rect(x+3, y+3, cardW-6, cardH-6, pixel.Accent)
		for j := 0; j < cardH-10; j += 4 {
			for i := j / 4 % 2 * 2; i < cardW-10; i += 4 {
				g.Rect(x+5+i, y+5+j, 2, 2, pixel.AccentLit)
			}
		}
		return
	}
	ink, idx := cd.ink(), cd.index()
	g.Micro(x+3, y+3, idx, ink, 1)
	g.Micro(x+cardW-3-pixel.MicroWidth(idx), y+cardH-8, idx, ink, 1)
	// Every card shows its suit in a corner.
	g.Sprite(x+cardW-14, y+2, suits[cd.Suit].sp, 0, ink, 1)
	// Court cards show their piece, the rest a big pip.
	sp := suits[cd.Suit].sp
	if cd.Rank > 10 {
		sp = [...]*pixel.Sprite{icon.Knight, icon.Queen, icon.King}[cd.Rank-11]
	}
	g.Sprite(x+(cardW-24)/2, y+(cardH-24)/2+2, sp, 0, ink, 2)
}

// pips places 0–6 pips on a face of side s.
func pips(g pixel.Canvas, x, y, s, n int, ink pixel.Slot) {
	if n < 0 || n > 6 {
		panic(fmt.Sprintf("tabletop: %d pips", n))
	}
	spots := [...][][2]int{{}, {{1, 1}}, {{0, 0}, {2, 2}}, {{0, 0}, {1, 1}, {2, 2}}, {{0, 0}, {2, 0}, {0, 2}, {2, 2}},
		{{0, 0}, {2, 0}, {1, 1}, {0, 2}, {2, 2}}, {{0, 0}, {2, 0}, {0, 1}, {2, 1}, {0, 2}, {2, 2}}}[n]
	d := max(2, s/7)
	for _, p := range spots {
		g.Rect(x+s*(p[0]+1)/4-d/2, y+s*(p[1]+1)/4-d/2, d, d, ink)
	}
}

// shades are Lit and Ink, lighter first: White looks white.
func shades(k pixel.Kit) (light, dark pixel.Slot) {
	lum := func(s pixel.Slot) int {
		c := k.C(s)
		return 299*int(c.R) + 587*int(c.G) + 114*int(c.B)
	}
	if lum(pixel.Lit) >= lum(pixel.Ink) {
		return pixel.Lit, pixel.Ink
	}
	return pixel.Ink, pixel.Lit
}

// players color tokens, one slot a player.
var players = []pixel.Slot{pixel.Bad, pixel.Info, pixel.Good, pixel.Warn, pixel.Magic, pixel.Ink}
