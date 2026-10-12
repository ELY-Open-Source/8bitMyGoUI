package game

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Fighter is a portrait in the turn order.
type Fighter struct {
	Name     string
	Portrait *pixel.Sprite
	Foe      bool
}

type TurnOrderOptions struct{}

// TurnOrder rows portraits in order, the active raised.
func TurnOrder(c *ui.Context, text string, order []Fighter, active int, o TurnOrderOptions) *ui.Element {
	k := pixel.K(c)
	if active < 0 || active >= len(order) {
		panic(fmt.Sprintf("game: TurnOrder active %d of %d", active, len(order)))
	}
	// Assistive technology hears who acts, from the active one on.
	names := make([]string, len(order))
	for i := range order {
		f := order[(active+i)%len(order)]
		if f.Portrait == nil {
			panic(fmt.Sprintf("game: fighter %q has no portrait", f.Name))
		}
		names[i] = f.Name
	}
	return ui.Box(c).Role(ui.RoleImage).Label(text).Description(strings.Join(names, ", ")).Size(k.D(len(order)*19), k.D(22)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, f := range order {
			x, y, edge := i*19, 5, pixel.Good
			if f.Foe {
				edge = pixel.Bad
			}
			if i == active {
				y, edge = 0, pixel.Accent
			}
			g.Rect(x, y, 18, 17, edge)
			g.Rect(x+1, y+1, 16, 15, pixel.Well)
			g.Sprite(x+3, y+3, f.Portrait, 0, pixel.Ink, 1)
		}
	})
}
