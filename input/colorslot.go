package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type ColorSlotOptions struct {
	Disabled bool
}

type ColorSlotResult struct {
	Element ui.Element
	Changed bool
}

// ColorSlot picks one of 16 palette slots from a grid.
func ColorSlot(c *ui.Context, text string, value *pixel.Slot, o ColorSlotOptions) ColorSlotResult {
	k := pixel.K(c)
	cur := int(*value)
	shell, base, st, _ := pixel.Pad(c, text, 8, pixel.Slots, 12, 12, &cur, o.Disabled)
	res := ColorSlotResult{Element: shell, Changed: pixel.Slot(cur) != *value}
	*value = pixel.Slot(cur)
	base.Description(value.String())
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i := range pixel.Slots {
			x, y := i%8*12, i/8*12
			if pixel.Slot(i) == *value {
				g.Rect(x, y, 12, 12, pixel.Ink)
				g.Rect(x+1, y+1, 10, 10, pixel.Lit)
			}
			g.Rect(x+2, y+2, 8, 8, pixel.Edge)
			g.Rect(x+3, y+3, 6, 6, pixel.Slot(i))
		}
		g.Finish(st)
	})
	if res.Changed {
		slog.Debug("8bit: color slot", "label", text, "slot", *value)
	}
	return res
}
