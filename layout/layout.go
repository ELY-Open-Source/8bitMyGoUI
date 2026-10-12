// Package layout arranges screens, windows and panes on dots.
package layout

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// box frames the canvas unless it is too small.
func box(g pixel.Canvas, style pixel.FrameStyle, fill pixel.Slot) {
	if g.W >= 3 && g.H >= 3 {
		g.Frame(0, 0, g.W, g.H, style, fill)
	}
}

// rule draws n dots on an axis in on/off runs.
func rule(g pixel.Canvas, x, y, n int, vertical bool, on, off int, s pixel.Slot) {
	for i := 0; i < n; i += on + off {
		run := min(on, n-i)
		if vertical {
			g.VLine(x, y+i, run, s)
		} else {
			g.HLine(x+i, y, run, s)
		}
	}
}

type dragKey struct{}

// drag is enabled e's pointer travel this frame, in dots.
func drag(k pixel.Kit, e *ui.Element) (dx, dy int, ok bool) {
	rest := ui.Local(e, dragKey{}, func() [2]float32 { return [2]float32{} })
	fx, fy, ok := e.Dragged()
	if !ok || e.IsDisabled() {
		*rest = [2]float32{}
		return 0, 0, false
	}
	x, y := rest[0]+fx/k.Dot(), rest[1]+fy/k.Dot()
	dx, dy = int(x), int(y)
	*rest = [2]float32{x - float32(dx), y - float32(dy)}
	return dx, dy, true
}

// label is o's label, or the built-in copy for key.
func label(k pixel.Kit, s, key string) string {
	if s != "" {
		return s
	}
	return k.Msg(key)
}
