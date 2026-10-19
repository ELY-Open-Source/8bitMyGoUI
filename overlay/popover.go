package overlay

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type PopoverOptions struct {
	// Width is the panel's in dots; zero fits the content.
	Width int
}

// notch is the pointer's height in dots.
const notch = 4

// Popover opens a panel below anchor, notched toward it.
func Popover(c *ui.Context, text string, anchor ui.Element, open *bool, o PopoverOptions, content func()) ui.Element {
	k := pixel.K(c)
	w := pixel.Or("Popover Width", o.Width, 0)
	ab := anchor.Bounds()
	return ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		if w > 0 {
			panel.Width(k.D(w))
		}
		panel.Label(text).Padding(k.D(notch+5), k.D(7)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			// Flipped above its anchor, the notch points down.
			up := r.Y >= ab.Y
			top := notch
			if !up {
				top = 0
			}
			g.Frame(0, top, g.W, g.H-notch, pixel.Window, pixel.Face)
			x := max(notch+1, min(g.W-notch-2, int((ab.X+ab.W/2-r.X)/k.Dot())))
			for j := range notch {
				y := j
				if !up {
					y = g.H - 1 - j
				}
				g.HLine(x-j, y, 2*j+1, pixel.Face)
				g.Dot(x-j, y, pixel.Edge)
				g.Dot(x+j, y, pixel.Edge)
			}
			edge := notch
			if !up {
				edge = g.H - notch - 1
			}
			g.HLine(x-notch+1, edge, 2*notch-1, pixel.Face)
		}).Children(content)
	})
}
