package text

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MarqueeOptions struct {
	// Width is the window in dots; zero means 120.
	Width int
	Slot  pixel.Slot
}

type marqueeKey struct{}

type marqueeState struct{ offset, last int }

// Marquee scrolls s a dot a step; hover pauses it.
func Marquee(c *ui.Context, s string, o MarqueeOptions) *ui.Element {
	k := pixel.K(c)
	w := pixel.Or("Marquee Width", o.Width, 120)
	e := ui.Box(c).Size(k.D(w), k.D(pixel.LineHeight)).Clip().Label(s).Role(ui.RoleText)
	ms := ui.Local(e, marqueeKey{}, func() marqueeState { return marqueeState{last: -1} })
	if k.Reduce() || e.Hovered() {
		ms.last = -1
	} else {
		step := k.Tick(c, 20)
		if ms.last >= 0 {
			ms.offset += step - ms.last
		}
		ms.last = step
	}
	period := w + pixel.Measure(s)
	x := -(ms.offset % period)
	if k.Reduce() {
		x = 0
	}
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, w, pixel.LineHeight, func() {
			g.Text(x, 0, s, o.Slot, 1)
			g.Text(x+period, 0, s, o.Slot, 1)
		})
	})
}
