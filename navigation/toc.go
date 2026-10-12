package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var pointer = pixel.MustSprite(`
	#..
	##.
	###
	##.
	#..
`)

type TocItem struct {
	Title string
	// Level indents sub-sections; zero is the top.
	Level int
}

type TableOfContentsOptions struct {
	// Width is in dots; zero means 160.
	Width    int
	Disabled bool
}

type TableOfContentsResult struct {
	Element *ui.Element
	Changed bool
}

// TableOfContents lists sections; ▶ marks the current one.
func TableOfContents(c *ui.Context, text string, items []TocItem, current *int, o TableOfContentsOptions) TableOfContentsResult {
	k := pixel.K(c)
	w := pixel.Or("TableOfContents Width", o.Width, 160)
	res := TableOfContentsResult{}
	res.Element = ui.RadioGroup(c, func() {
		for i, it := range items {
			shell, base, st := pixel.Shell(c, func() *ui.Element {
				return ui.RadioBase(c, current, i).Label(it.Title).Disabled(o.Disabled)
			})
			if base.Changed() {
				res.Changed = true
			}
			shell.Size(k.D(w), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				x, ink := 8+it.Level*10, pixel.Faint
				if i == *current {
					g.Sprite(x-6, 5, pointer, 0, pixel.Accent, 1)
					ink = pixel.Accent
				} else if st.Hover {
					ink = pixel.Ink
				}
				g.Text(x, 0, it.Title, ink, 1)
				g.Finish(st)
			})
		}
	}).Gap(0).Label(text)
	if res.Changed {
		slog.Debug("8bit: toc", "label", text, "section", items[*current].Title)
	}
	return res
}
