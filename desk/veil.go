package desk

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type wakeKey struct{}

// wake is input between frames; from is the first pointer.
type wake struct {
	woke, seen bool
	from       [2]float32
}

// cover fills the window; keys, clicks, long moves end it.
func cover(c *ui.Context, k pixel.Kit, text string, open *bool, moves bool, draw func(g pixel.Canvas)) (anchor, panel ui.Element, woke bool) {
	anchor = ui.Box(c).Role(ui.RoleNone)
	st := ui.Local(anchor, wakeKey{}, func() wake { return wake{} })
	was := *open
	if st.woke {
		*open = false
	}
	if !*open {
		*st = wake{}
	}
	panel = ui.DialogBase(c, open, func(back, panel ui.Element) {
		back.Draw(func(p *ui.Painter, r ui.Rect) { draw(k.Canvas(p, r)) })
		back.HandleInput(func(ev ui.InputEvent) bool {
			// Any button and the wheel end it, as keys do.
			if ev.Kind == ui.InputPointerDown || ev.Kind == ui.InputScroll {
				st.woke = true
			}
			if !moves || ev.Kind != ui.InputPointerMove {
				return false
			}
			// The first move marks the spot; eight dots away wakes.
			at := [2]float32{ev.X / k.Dot(), ev.Y / k.Dot()}
			if !st.seen {
				st.from, st.seen = at, true
			}
			if dx, dy := at[0]-st.from[0], at[1]-st.from[1]; dx*dx+dy*dy > 64 {
				st.woke = true
			}
			return false
		})
		panel.Focusable().FocusRing(false).Label(text).HandleInput(func(ev ui.InputEvent) bool {
			if ev.Kind == ui.InputKeyDown {
				st.woke = true
				return true
			}
			return false
		})
	})
	return anchor, panel, was && !*open
}
