package feedback

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type SaveIndicatorOptions struct{}

type saveKey struct{}

// spin is the floppy's width, frame by frame.
var spin = [4]int{12, 8, 2, 8}

// SaveIndicator spins a floppy while saving, then shows SAVED ✓.
func SaveIndicator(c *ui.Context, saving bool, o SaveIndicatorOptions) ui.Element {
	k := pixel.K(c)
	e := ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center)
	ever := ui.Local(e, saveKey{}, func() bool { return false })
	if saving {
		*ever = true
	}
	if !*ever {
		return e
	}
	word, ink := k.Msg("saved"), pixel.Good
	if saving {
		word, ink = k.Msg("saving")+"…", pixel.Faint
	}
	return e.Role(ui.RoleStatus).Label(word).Children(func() {
		ui.Box(c).Size(k.D(12), k.D(12)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			if !saving {
				g.Sprite(0, 0, icon.Check, 0, pixel.Good, 1)
				return
			}
			// A coin spin: columns squeeze in, crisp.
			w := spin[g.Busy(8)%4]
			x0 := (12 - w) / 2
			for x := range w {
				sx := x * 12 / w
				g.Clip(x0+x, 0, 1, 12, func() { g.Sprite(x0+x-sx, 0, icon.Save, 0, pixel.Info, 1) })
			}
		})
		pixel.Text(c, word, pixel.TextOptions{Slot: ink})
	})
}
