package datetime

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type EventChipOptions struct {
	// Slot colors the chip's edge; zero means Accent.
	Slot pixel.Slot
}

// EventChip is a small chip with time and title.
func EventChip(c *ui.Context, e Event, o EventChipOptions) ui.Element {
	k := pixel.K(c)
	edge := o.Slot
	if edge == 0 {
		edge = pixel.Accent
	}
	at := e.Start.Format("15:04")
	w := pixel.Measure(at) + pixel.Measure(e.Title) + 14
	return ui.Box(c).Size(k.D(w), k.D(pixel.LineHeight+2)).Role(ui.RoleText).Label(at + " " + e.Title).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Flat, pixel.Lit)
		g.Rect(1, 1, 2, g.H-2, edge)
		g.Text(6, 1, at, pixel.Faint, 1)
		g.Text(10+pixel.Measure(at), 1, e.Title, pixel.Ink, 1)
	})
}
