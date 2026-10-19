package layout

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type DividerStyle uint8

const (
	Solid DividerStyle = iota
	Dotted
	Dashed
	Double
	Ornament
)

type DividerOptions struct {
	Vertical bool
	Style    DividerStyle
	// Slot is the ink; zero means Edge.
	Slot pixel.Slot
}

var diamond = pixel.MustSprite(`
	..#..
	.###.
	#####
	.###.
	..#..
`)

// Divider is a rule across its parent.
func Divider(c *ui.Context, o DividerOptions) ui.Element {
	k := pixel.K(c)
	ink := o.Slot
	if ink == pixel.Ink {
		ink = pixel.Edge
	}
	thick := [...]int{Solid: 1, Dotted: 1, Dashed: 1, Double: 3, Ornament: 5}
	if int(o.Style) >= len(thick) {
		panic(fmt.Sprintf("layout: unknown divider style %d", o.Style))
	}
	t := thick[o.Style]
	e := ui.Box(c).Role(ui.RoleNone)
	if o.Vertical {
		e.FillHeight().Width(k.D(t))
	} else {
		e.FillWidth().Height(k.D(t))
	}
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		n := g.W
		if o.Vertical {
			n = g.H
		}
		line := func(at, from, n int) {
			x, y := from, at
			if o.Vertical {
				x, y = at, from
			}
			switch o.Style {
			case Dotted:
				w, h := n, 1
				if o.Vertical {
					w, h = 1, n
				}
				g.Checker(x, y, w, h, ink)
			case Dashed:
				rule(g, x, y, n, o.Vertical, 4, 2, ink)
			default:
				rule(g, x, y, n, o.Vertical, max(n, 1), 0, ink)
			}
		}
		switch o.Style {
		case Double:
			line(0, 0, n)
			line(2, 0, n)
		case Ornament:
			mid := n / 2
			line(2, 0, mid-4)
			line(2, mid+5, n-mid-5)
			if o.Vertical {
				g.Sprite(0, mid-2, diamond, 0, ink, 1)
			} else {
				g.Sprite(mid-2, 0, diamond, 0, ink, 1)
			}
		default:
			line(0, 0, n)
		}
	})
}
