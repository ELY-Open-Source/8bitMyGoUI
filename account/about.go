package account

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

// Credit is a role and who filled it.
type Credit struct {
	Role, Name string
}

type AboutOptions struct {
	Version string
	Credits []Credit
}

// About shows the title as a logo, version, then credits.
func About(c *ui.Context, title string, o AboutOptions) *ui.Element {
	k := pixel.K(c)
	if o.Version == "" {
		panic("account: About needs a version")
	}
	return ui.Column(c).Gap(k.D(4)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(title).Children(func() {
		text.OutlineText(c, title, text.OutlineTextOptions{Slot: pixel.Accent, Scale: 3})
		pixel.Text(c, k.Msg("version")+" "+o.Version, pixel.TextOptions{Slot: pixel.Faint})
		if len(o.Credits) == 0 {
			return
		}
		ui.Box(c).Size(k.D(80), k.D(5)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for x := 0; x < g.W; x += 4 {
				g.Rect(x, 2, 2, 1, pixel.Shade)
			}
		})
		pixel.Text(c, k.Msg("credits"), pixel.TextOptions{Slot: pixel.Accent})
		for _, cr := range o.Credits {
			pixel.Text(c, cr.Role, pixel.TextOptions{Slot: pixel.Faint})
			pixel.Text(c, cr.Name, pixel.TextOptions{})
		}
	})
}
