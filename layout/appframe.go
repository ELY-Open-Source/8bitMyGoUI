package layout

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type AppFrameOptions struct {
	// SidebarWidth is zero for 96 dots.
	SidebarWidth int
	// FoldBelow stacks the sidebar over main; zero means 320.
	FoldBelow int
}

// AppFrame arranges header, sidebar, main and footer; nil ones vanish.
func AppFrame(c *ui.Context, o AppFrameOptions, header, sidebar, main, footer func()) *ui.Element {
	k := pixel.K(c)
	sw := orDef("AppFrame SidebarWidth", o.SidebarWidth, 96)
	fold := orDef("AppFrame FoldBelow", o.FoldBelow, 320)
	if main == nil {
		panic("layout: AppFrame needs main")
	}
	root := ui.Column(c).Fill()
	w, check := pixel.WrapWidth(c, root)
	folded := w > 0 && w < fold
	strip := func(fn func(), below bool) {
		ui.Box(c).FillWidth().Padding(k.D(4), k.D(6)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Rect(0, 0, g.W, g.H, pixel.Face)
			if below {
				g.HLine(0, g.H-1, g.W, pixel.Edge)
			} else {
				g.HLine(0, 0, g.W, pixel.Edge)
			}
		}).Children(fn)
	}
	return root.Draw(func(p *ui.Painter, r ui.Rect) {
		check(r)
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Back)
	}).Children(func() {
		if header != nil {
			strip(header, true)
		}
		var body *ui.Element
		if folded {
			body = ui.Column(c)
		} else {
			body = ui.Row(c)
		}
		body.Grow(1).AlignItems(ui.Stretch).Children(func() {
			if sidebar != nil {
				side := ui.Column(c).Padding(k.D(4)).Gap(k.D(2)).Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					g.Rect(0, 0, g.W, g.H, pixel.Well)
					if folded {
						g.HLine(0, g.H-1, g.W, pixel.Edge)
					} else {
						g.VLine(g.W-1, 0, g.H, pixel.Edge)
					}
				})
				if !folded {
					side.Width(k.D(sw))
				}
				side.Children(sidebar)
			}
			ui.Column(c).Grow(1).Padding(k.D(6)).Gap(k.D(4)).Clip().Children(main)
		})
		if footer != nil {
			strip(footer, false)
		}
	})
}
