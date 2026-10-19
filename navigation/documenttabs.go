package navigation

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var cross = pixel.MustSprite(`
	#...#
	.#.#.
	..#..
	.#.#.
	#...#
`)

type DocumentTabsOptions struct {
	Disabled bool
}

type DocumentTabsResult struct {
	Element ui.Element
	Changed bool
	// Closed is the tab × or Delete closed, or -1.
	Closed int
}

// DocumentTabs are closable tabs; × or Delete reports Closed.
func DocumentTabs(c *ui.Context, text string, titles []string, selected *int, o DocumentTabsOptions) DocumentTabsResult {
	k := pixel.K(c)
	seen := map[string]bool{}
	for _, t := range titles {
		if seen[t] {
			panic(fmt.Sprintf("navigation: DocumentTabs has %q twice", t))
		}
		seen[t] = true
	}
	res := DocumentTabsResult{Closed: -1}
	// The base turns before it learns it is off.
	v := *selected
	p := ui.TabsBase(c, &v, len(titles))
	p.List.Label(text).Children(func() {
		for i, t := range titles {
			// Keyed by title, so focus follows its tab.
			ui.Box(c).Key(t).Children(func() {
				var x ui.Element
				shell, base, st := pixel.Shell(c, func() ui.Element {
					tab := p.Tab(i).Label(t).Disabled(o.Disabled)
					// The × is no Tab stop; Delete closes by keyboard.
					tab.Children(func() {
						x = ui.Box(c).Absolute().Right(k.D(4)).Top(k.D(6)).Size(k.D(7), k.D(7)).Role(ui.RoleButton).Label(k.Msg("close") + " " + t)
					})
					return tab
				})
				del := pixel.Shortcut(base, 0, ui.KeyDelete)
				if x.Clicked() || del {
					res.Closed = i
				}
				hot := x.Hovered()
				shell.Size(k.D(pixel.Measure(t)+24), k.D(pixel.LineHeight+6)).Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					folder(g, t, i == *selected, false, st)
					ink := pixel.Faint
					if hot {
						ink = pixel.Bad
					}
					g.Sprite(g.W-10, 7, cross, 0, ink, 1)
				})
			})
		}
	})
	res.Element, res.Changed = p.List, keep(c, v, selected, o.Disabled)
	if res.Closed >= 0 {
		slog.Debug("8bit: tab closed", "label", text, "tab", titles[res.Closed])
	}
	return res
}
