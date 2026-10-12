package display

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type Side uint8

const (
	Bottom Side = iota
	Top
	Left
	Right
)

type BubbleOptions struct {
	// Tail is the side the tail leaves from.
	Tail Side
	// At offsets the tail along its side; zero means 8.
	At int
	// Fill is the bubble's slot; zero means Lit.
	Fill pixel.Slot
}

const tailLen = 5

// Bubble is a speech bubble with a tail around children.
func Bubble(c *ui.Context, o BubbleOptions, children func()) *ui.Element {
	if o.Tail > Right {
		panic(fmt.Sprintf("display: unknown Side %d", o.Tail))
	}
	k := pixel.K(c)
	at, fill := pixel.Or("Bubble At", o.At, 8), slotOr(o.Fill, pixel.Lit)
	pad := [4]int{6, 8, 6, 8}
	pad[[...]int{Bottom: 2, Top: 0, Left: 3, Right: 1}[o.Tail]] += tailLen
	return ui.Column(c).Padding(k.D(pad[0]), k.D(pad[1]), k.D(pad[2]), k.D(pad[3])).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		x, y, w, h := 0, 0, g.W, g.H
		switch o.Tail {
		case Bottom:
			h -= tailLen
		case Top:
			y, h = tailLen, h-tailLen
		case Left:
			x, w = tailLen, w-tailLen
		case Right:
			w -= tailLen
		}
		if w < 6 || h < 6 {
			return
		}
		g.Frame(x, y, w, h, pixel.Notched, fill)
		for i := range tailLen {
			n := tailLen - i
			switch o.Tail {
			case Bottom:
				g.HLine(x+at-1, y+h-1+i, n+1, pixel.Edge)
				g.HLine(x+at, y+h-1+i, n-1, fill)
			case Top:
				g.HLine(x+at-1, y-i, n+1, pixel.Edge)
				g.HLine(x+at, y-i, n-1, fill)
			case Left:
				g.VLine(x-i, y+at-1, n+1, pixel.Edge)
				g.VLine(x-i, y+at, n-1, fill)
			case Right:
				g.VLine(x+w-1+i, y+at-1, n+1, pixel.Edge)
				g.VLine(x+w-1+i, y+at, n-1, fill)
			}
		}
	}).Children(children)
}
