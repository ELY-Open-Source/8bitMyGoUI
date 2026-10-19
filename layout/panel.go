package layout

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type PanelOptions struct {
	Style pixel.FrameStyle
	// Fill is the body slot; zero means Face.
	Fill pixel.Slot
	// Padding is inside the border; zero means 4.
	Padding int
}

// Panel frames children in a style.
func Panel(c *ui.Context, o PanelOptions, children func()) ui.Element {
	k := pixel.K(c)
	fill := o.Fill
	if fill == pixel.Ink {
		fill = pixel.Face
	}
	pad := border(o.Style) + pixel.Or("Panel Padding", o.Padding, 4)
	return ui.Column(c).Padding(k.D(pad)).Draw(func(p *ui.Painter, r ui.Rect) {
		box(k.Canvas(p, r), o.Style, fill)
	}).Children(children)
}

// border is the dots a frame style takes per side.
func border(s pixel.FrameStyle) int {
	switch s {
	case pixel.Flat, pixel.Notched:
		return 1
	case pixel.Raised, pixel.Sunken:
		return 2
	case pixel.Double, pixel.Window:
		return 3
	}
	panic("layout: unknown frame style")
}
