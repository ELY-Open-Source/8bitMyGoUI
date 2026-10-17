package display

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type RibbonOptions struct {
	// Slot is the band; zero means Accent.
	Slot pixel.Slot
}

const tail = 12

// Ribbon is a title banner with notched, folded tails.
func Ribbon(c *ui.Context, title string, o RibbonOptions) ui.Element {
	k := pixel.K(c)
	band := slotOr(o.Slot, pixel.Accent)
	tw := pixel.Measure(title)
	w, h := tw+2*tail+8, pixel.LineHeight+4
	return ui.Box(c).Size(k.D(w), k.D(h+4)).Role(ui.RoleText).Label(title).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for j := range h {
			cut := h/2 - abs(j-h/2)
			for _, x0 := range []int{cut, w - tail} {
				n := tail - cut
				g.HLine(x0, 4+j, n, band)
				g.Veil(x0, 4+j, n, 1, pixel.Shade, 8)
				if j == 0 || j == h-1 {
					g.HLine(x0, 4+j, n, pixel.Edge)
				}
			}
		}
		g.Rect(tail-3, h+1, 3, 3, pixel.Edge)
		g.Rect(w-tail, h+1, 3, 3, pixel.Edge)
		g.Frame(tail-3, 0, w-2*tail+6, h, pixel.Flat, band)
		g.HLine(tail-2, 1, w-2*tail+4, pixel.AccentLit)
		g.Text((w-tw)/2, 2, title, pixel.AccentInk, 1)
	})
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
