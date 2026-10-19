// Package shop holds prices, carts, rewards and trading components.
package shop

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Currency is what a price is counted in.
type Currency int

const (
	Gold Currency = iota
	Gems
)

// coin is a currency's icon and word; unknown ones panic.
func coin(k pixel.Kit, cur Currency) (*pixel.Sprite, string) {
	switch cur {
	case Gold:
		return icon.Coin, k.Msg("coins")
	case Gems:
		return icon.Gem, k.Msg("gems")
	}
	panic(fmt.Sprintf("shop: unknown currency %d", cur))
}

// amount checks n is not negative and prints it.
func amount(name string, n int) string {
	if n < 0 {
		panic(fmt.Sprintf("shop: %s amount %d is negative", name, n))
	}
	return strconv.Itoa(n)
}

// money draws a currency icon and amount, named for readers.
func money(c *ui.Context, k pixel.Kit, cur Currency, n int, slot pixel.Slot) ui.Element {
	sp, word := coin(k, cur)
	s := amount("money", n)
	return ui.Box(c).Size(k.D(14+pixel.Measure(s)), k.D(pixel.LineHeight)).Role(ui.RoleText).Label(s + " " + word).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Sprite(0, 2, sp, 0, pixel.Warn, 1)
		g.Text(14, 0, s, slot, 1)
	})
}
