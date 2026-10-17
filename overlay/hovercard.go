package overlay

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type HoverCardOptions struct {
	// Width is the card's in dots; zero means 160.
	Width int
}

type hoverKey struct{}

// HoverCard shows card below trigger while either is hovered.
func HoverCard(c *ui.Context, text string, o HoverCardOptions, trigger func(), card func()) ui.Element {
	k := pixel.K(c)
	w := pixel.Or("HoverCard Width", o.Width, 160)
	wrap := ui.Box(c)
	last := ui.Local(wrap, hoverKey{}, func() ui.Rect { return ui.Rect{} })
	b := wrap.Bounds()
	px, py, _ := wrap.PointerPosition()
	show := hint(c, wrap, 0) || last.Contains(b.X+px, b.Y+py)
	wrap.Children(trigger)
	*last = ui.Rect{}
	if !show {
		return wrap
	}
	ww, _ := c.Size()
	ui.Overlay(c, func() {
		e := ui.Column(c).Absolute().Left(max(0, min(b.X, ww-k.D(w)))).Top(b.Y + b.H).Width(k.D(w)).
			Padding(k.D(6)).Role(ui.RoleGroup).Label(text)
		*last = e.Bounds()
		e.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
		}).Children(card)
	})
	return wrap
}
