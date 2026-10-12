package files

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PathBarOptions struct {
	Disabled bool
}

type PathBarResult struct {
	Element *ui.Element
	// Clicked is the folder clicked, else -1.
	Clicked int
}

// PathBar sets a path's folders in a well, ▸ between.
func PathBar(c *ui.Context, text string, parts []string, o PathBarOptions) PathBarResult {
	k := pixel.K(c)
	if len(parts) == 0 || slices.Contains(parts, "") {
		panic(fmt.Sprintf("files: PathBar parts %q", parts))
	}
	res := PathBarResult{Clicked: -1}
	res.Element = layout.Panel(c, layout.PanelOptions{Style: pixel.Sunken, Fill: pixel.Well, Padding: 1}, func() {
		// One Tab stop; arrows walk the folders.
		ui.Row(c).Gap(k.D(2)).AlignItems(ui.Center).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).Label(text).Children(func() {
			ui.Box(c).Size(k.D(icon.FolderOpen.W), k.D(icon.FolderOpen.H)).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Sprite(0, 0, icon.FolderOpen, 0, pixel.Warn, 1)
			})
			for i, part := range parts {
				if i > 0 {
					ui.Box(c).Size(k.D(pixel.Measure("▸")), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
						k.Canvas(p, r).Text(0, 0, "▸", pixel.Faint, 1)
					})
				}
				// The last folder is where you are; it stays put.
				if i == len(parts)-1 {
					ui.Box(c).Size(k.D(pixel.Measure(part)+4), k.D(pixel.LineHeight)).Role(ui.RoleText).Label(part).Draw(func(p *ui.Painter, r ui.Rect) {
						k.Canvas(p, r).Text(2, 0, part, pixel.Accent, 1)
					})
					break
				}
				shell, base, st := pixel.Shell(c, func() *ui.Element {
					return ui.ButtonBase(c).Label(part).Disabled(o.Disabled)
				})
				if base.Clicked() {
					res.Clicked = i
				}
				shell.Size(k.D(pixel.Measure(part)+4), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					if st.Hover || st.Press {
						g.Rect(0, 0, g.W, g.H, pixel.Face)
					}
					g.Text(2, 0, part, pixel.Ink, 1)
					g.Finish(st)
				})
			}
		})
	})
	if res.Clicked >= 0 {
		slog.Debug("8bit: path bar", "label", text, "part", res.Clicked)
	}
	return res
}
