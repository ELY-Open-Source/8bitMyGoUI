package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SidebarItem struct {
	Icon  *pixel.Sprite
	Label string
	// Badge is a short count or word at the right.
	Badge string
}

type SidebarSection struct {
	Title string
	Items []SidebarItem
}

type SidebarOptions struct {
	// Width is in dots; zero means 140.
	Width    int
	Disabled bool
}

type SidebarResult struct {
	Element *ui.Element
	Changed bool
}

// Sidebar lists sections of items; Up and Down choose.
func Sidebar(c *ui.Context, text string, sections []SidebarSection, selected *int, o SidebarOptions) SidebarResult {
	k := pixel.K(c)
	w := pixel.Or("Sidebar Width", o.Width, 140)
	res := SidebarResult{}
	n := 0
	res.Element = ui.RadioGroup(c, func() {
		for _, sec := range sections {
			if sec.Title != "" {
				pixel.Text(c, sec.Title, pixel.TextOptions{Slot: pixel.Faint})
			}
			for _, it := range sec.Items {
				i := n
				n++
				shell, base, st := pixel.Shell(c, func() *ui.Element {
					return ui.RadioBase(c, selected, i).Label(it.Label).Disabled(o.Disabled)
				})
				if base.Changed() {
					res.Changed = true
				}
				shell.Size(k.D(w), k.D(pixel.LineHeight+2)).Draw(func(p *ui.Painter, r ui.Rect) {
					entry(k.Canvas(p, r), it, i == *selected, st)
				})
			}
		}
	}).Gap(k.D(1)).Label(text)
	if res.Changed {
		slog.Debug("8bit: sidebar", "label", text, "selected", *selected)
	}
	return res
}

// entry draws a sidebar row: icon, label, badge.
func entry(g pixel.Canvas, it SidebarItem, on bool, st pixel.State) {
	ink := pixel.Ink
	switch {
	case on:
		g.Rect(0, 0, g.W, g.H, pixel.Accent)
		ink = pixel.AccentInk
	case st.Hover:
		g.Rect(0, 0, g.W, g.H, pixel.Lit)
	}
	x := 4
	if it.Icon != nil {
		g.Sprite(x, (g.H-it.Icon.H)/2, it.Icon, 0, ink, 1)
		x += it.Icon.W + 4
	}
	g.Text(x, 1, it.Label, ink, 1)
	if it.Badge != "" {
		bw := pixel.Measure(it.Badge) + 6
		g.Frame(g.W-bw-3, 2, bw, g.H-4, pixel.Notched, pixel.Bad)
		g.Text(g.W-bw, 1, it.Badge, pixel.AccentInk, 1)
	}
	g.Finish(st)
}
