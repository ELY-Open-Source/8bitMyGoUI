package feedback

import (
	"fmt"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type ProgressOptions struct {
	// Width is the bar in dots; zero means 120.
	Width int
	// Percent shows the number after the bar.
	Percent bool
}

// Progress fills chunk by chunk; below zero a block runs.
func Progress(c *ui.Context, text string, value float64, o ProgressOptions) ui.Element {
	k := pixel.K(c)
	if value > 1 || math.IsNaN(value) {
		panic(fmt.Sprintf("feedback: Progress value %v is past 1", value))
	}
	w := pixel.Or("Progress Width", o.Width, 120)
	n := (w - 3) / 5
	pct := strconv.Itoa(int(math.Round(value*100))) + "%"
	e := ui.Box(c).Role(ui.RoleProgress).Label(text)
	full := w
	if value < 0 {
		e.Description(k.Msg("loading"))
	} else {
		e.Description(pct)
		if o.Percent {
			full += 4 + pixel.MicroWidth("100%")
		}
	}
	return e.Size(k.D(full), k.D(10)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, w, 10, pixel.Sunken, pixel.Well)
		chunk := func(i int) { g.Rect(3+i*5, 3, 4, 4, pixel.Accent) }
		if value < 0 {
			at := pingpong(g.Busy(12), n-3)
			for i := range 3 {
				chunk(at + i)
			}
			return
		}
		for i := range int(value * float64(n)) {
			chunk(i)
		}
		if o.Percent {
			g.Micro(w+4, 2, pct, pixel.Ink, 1)
		}
	})
}
