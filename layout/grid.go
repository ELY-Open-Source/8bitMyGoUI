package layout

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type GridLines uint8

const (
	NoLines GridLines = iota
	SolidLines
	DottedLines
)

type GridOptions struct {
	// Columns is zero for 2.
	Columns int
	// Pad is around each cell's content; zero means 4.
	Pad   int
	Lines GridLines
}

// Grid lays n cells in equal columns, lines between them.
func Grid(c *ui.Context, n int, o GridOptions, cell func(i int)) *ui.Element {
	k := pixel.K(c)
	cols := orDef("Grid Columns", o.Columns, 2)
	pad := orDef("Grid Pad", o.Pad, 4)
	orDef("Grid n", n, 0)
	if o.Lines > DottedLines {
		panic(fmt.Sprintf("layout: unknown grid lines %d", o.Lines))
	}
	grid := ui.Grid(c).FillWidth()
	w, check := pixel.WrapWidth(c, grid)
	tracks := make([]ui.Track, cols)
	for i, d := range shares(w, cols) {
		tracks[i] = ui.Fixed(k.D(d))
	}
	return grid.ColumnTracks(tracks...).Draw(func(p *ui.Painter, r ui.Rect) { check(r) }).Children(func() {
		for i := range n {
			right, below := edges(i, n, cols)
			ui.Box(c).Padding(k.D(pad)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				line := func(x, y, w, h int) {
					if o.Lines == DottedLines {
						g.Checker(x, y, w, h, pixel.Edge)
					} else if o.Lines == SolidLines {
						g.Rect(x, y, w, h, pixel.Edge)
					}
				}
				if right {
					line(g.W-1, 0, 1, g.H)
				}
				if below {
					line(0, g.H-1, g.W, 1)
				}
			}).Children(func() { cell(i) })
		}
	})
}

// edges reports neighbors right of and below cell i.
func edges(i, n, cols int) (right, below bool) {
	return i%cols != cols-1 && i != n-1, i < (n-1)/cols*cols
}
