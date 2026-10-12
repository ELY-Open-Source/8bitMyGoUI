package data

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CardGridOptions struct {
	// MinWidth is a card's least dots; zero means 96.
	MinWidth int
}

// CardGrid fits as many card columns as width allows.
func CardGrid(c *ui.Context, n int, o CardGridOptions, card func(i int)) *ui.Element {
	mw := pixel.Or("CardGrid MinWidth", o.MinWidth, 96)
	box := ui.Box(c).FillWidth()
	w, check := pixel.WrapWidth(c, box)
	return box.Draw(func(p *ui.Painter, r ui.Rect) { check(r) }).Children(func() {
		layout.Grid(c, n, layout.GridOptions{Columns: max(1, w/mw), Pad: 2}, card)
	})
}
