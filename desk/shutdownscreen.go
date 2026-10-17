package desk

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ShutdownScreenOptions struct {
	// Width and Height are dots; zeros mean 320 by 200.
	Width, Height int
}

// ShutdownScreen says it is safe to turn off the computer.
func ShutdownScreen(c *ui.Context, o ShutdownScreenOptions) ui.Element {
	k := pixel.K(c)
	w, h := pixel.Or("ShutdownScreen Width", o.Width, 320), pixel.Or("ShutdownScreen Height", o.Height, 200)
	s := k.Msg("safe_to_turn_off")
	lines := pixel.Wrap(s, (w-16)/2, 0)
	if len(lines)*2*pixel.LineHeight > h-16 {
		panic(fmt.Sprintf("desk: ShutdownScreen %d×%d is too small for its words", w, h))
	}
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleText).Label(s).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, w, h, pixel.Well)
		top := (h - len(lines)*2*pixel.LineHeight) / 2
		for i, l := range lines {
			g.Text((w-2*pixel.Measure(l))/2, top+i*2*pixel.LineHeight, l, pixel.Warn, 2)
		}
	})
}
