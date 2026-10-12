package agent

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CostCoinsOptions struct {
	// Symbol leads the amount; zero means "$".
	Symbol string
}

// CostCoins stacks cost as gold dollars, silver dimes, bronze cents.
func CostCoins(c *ui.Context, text string, cents int, o CostCoinsOptions) *ui.Element {
	k := pixel.K(c)
	if cents < 0 {
		panic(fmt.Sprintf("agent: CostCoins of %d cents", cents))
	}
	sym := o.Symbol
	if sym == "" {
		sym = "$"
	}
	amount := fmt.Sprintf("%s%d.%02d", sym, cents/100, cents%100)
	counts := [3]int{cents / 100, cents / 10 % 10, cents % 10}
	var cols [3]int
	stacks := 8
	for i, n := range counts {
		cols[i] = max(14, pixel.MicroWidth(strconv.Itoa(n))+2)
		stacks += cols[i]
	}
	w := max(stacks, pixel.Measure(amount))
	e := ui.Box(c).Size(k.D(w), k.D(66)).Role(ui.RoleImage).Label(text).Description(amount)
	// A change drops each stack's top coin in.
	drop := []int{-6, -3, -1, 0}[k.Shot(c, e, cents, 12, 3)]
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		x := (w - stacks) / 2
		for i, n := range counts {
			cx := x + (cols[i]-12)/2
			g.HLine(cx, 42, 12, pixel.Shade)
			top := min(n, 10)
			for j := range top {
				y := 38 - 3*j
				if j == top-1 {
					y += drop
				}
				coin(g, cx, y, i)
			}
			s := strconv.Itoa(n)
			g.Micro(x+(cols[i]-pixel.MicroWidth(s))/2, 45, s, pixel.Ink, 1)
			x += cols[i] + 4
		}
		g.Text((w-pixel.Measure(amount))/2, 50, amount, pixel.Ink, 1)
	})
}

// coin draws one edge-on: 0 gold, 1 silver, 2 bronze.
func coin(g pixel.Canvas, x, y, metal int) {
	fill := pixel.Warn
	if metal == 1 {
		fill = pixel.Face
	}
	g.HLine(x+1, y, 10, pixel.Edge)
	g.Rect(x, y+1, 12, 2, pixel.Edge)
	g.Rect(x+1, y+1, 10, 2, fill)
	if metal == 2 {
		g.Checker(x+1, y+1, 10, 2, pixel.Bad)
	}
	g.Dot(x+2, y+1, pixel.Lit)
	g.HLine(x+1, y+3, 10, pixel.Edge)
}
