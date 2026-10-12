package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type BreadcrumbsOptions struct {
	Disabled bool
}

type BreadcrumbsResult struct {
	Element *ui.Element
	// Clicked is the part clicked, or -1.
	Clicked int
}

// Breadcrumbs is a path whose parts, but the last, link.
func Breadcrumbs(c *ui.Context, text string, parts []string, o BreadcrumbsOptions) BreadcrumbsResult {
	k := pixel.K(c)
	res := BreadcrumbsResult{Clicked: -1}
	res.Element = ui.Row(c).Gap(k.D(3)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(text).Children(func() {
		for i, part := range parts {
			if i == len(parts)-1 {
				pixel.Text(c, part, pixel.TextOptions{})
				break
			}
			shell, base, st := pixel.Shell(c, func() *ui.Element {
				return ui.ButtonBase(c).Role(ui.RoleLink).Label(part).Disabled(o.Disabled)
			})
			if base.Clicked() {
				res.Clicked = i
			}
			shell.Size(k.D(pixel.Measure(part)+2), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				ink := pixel.Accent
				if st.Hover {
					ink = pixel.AccentLit
				}
				g.Text(1, 0, part, ink, 1)
				g.HLine(1, g.H-2, g.W-2, ink)
				g.Finish(st)
			})
			ui.Box(c).Size(k.D(pointer.W), k.D(pointer.H)).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Sprite(0, 0, pointer, 0, pixel.Faint, 1)
			})
		}
	})
	if res.Clicked >= 0 {
		slog.Debug("8bit: breadcrumb", "part", parts[res.Clicked])
	}
	return res
}
