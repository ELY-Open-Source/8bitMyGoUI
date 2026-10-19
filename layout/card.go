package layout

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type CardOptions struct {
	Title string
	Icon  *pixel.Sprite
	// Footer builds a strip under the body.
	Footer func()
}

// Card is a raised panel: header strip, body, footer.
func Card(c *ui.Context, o CardOptions, body func()) ui.Element {
	k := pixel.K(c)
	e := ui.Column(c).Padding(k.D(2)).Draw(func(p *ui.Painter, r ui.Rect) {
		box(k.Canvas(p, r), pixel.Raised, pixel.Face)
	})
	if o.Title != "" {
		e.Role(ui.RoleGroup).Label(o.Title)
	}
	return e.Children(func() {
		if o.Title != "" || o.Icon != nil {
			cardHeader(c, k, o)
		}
		ui.Column(c).Padding(k.D(6)).Gap(k.D(4)).Children(body)
		if o.Footer != nil {
			ui.Row(c).Padding(k.D(5), k.D(6), k.D(4), k.D(6)).Gap(k.D(4)).AlignItems(ui.Center).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.HLine(0, 0, g.W, pixel.Shade)
				g.HLine(0, 1, g.W, pixel.Lit)
			}).Children(o.Footer)
		}
	})
}

func cardHeader(c *ui.Context, k pixel.Kit, o CardOptions) {
	ui.Box(c).FillWidth().Height(k.D(pixel.LineHeight + 4)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Accent)
		x := 4
		if o.Icon != nil {
			g.Sprite(x, (g.H-o.Icon.H)/2, o.Icon, 0, pixel.AccentInk, 1)
			x += o.Icon.W + 4
		}
		g.Clip(0, 0, g.W-2, g.H, func() { g.Text(x, 2, o.Title, pixel.AccentInk, 1) })
	})
}
