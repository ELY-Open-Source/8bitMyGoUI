package agent

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// FileDiff is a changed file: lines added and removed.
type FileDiff struct {
	Path           string
	Added, Removed int
}

type DiffScrollOptions struct {
	// Width is in dots; zero means 220.
	Width int
}

// DiffScroll unrolls a parchment of changed files, lines counted.
func DiffScroll(c *ui.Context, text string, files []FileDiff, o DiffScrollOptions) *ui.Element {
	k := pixel.K(c)
	if len(files) == 0 {
		panic("agent: DiffScroll has no files")
	}
	add, del := 0, 0
	said := make([]string, 0, len(files)+1)
	for _, f := range files {
		if f.Path == "" || f.Added < 0 || f.Removed < 0 || f.Added > math.MaxInt-add || f.Removed > math.MaxInt-del {
			panic(fmt.Sprintf("agent: DiffScroll file %q +%d -%d", f.Path, f.Added, f.Removed))
		}
		add, del = add+f.Added, del+f.Removed
		said = append(said, f.Path+" +"+strconv.Itoa(f.Added)+" -"+strconv.Itoa(f.Removed))
	}
	w := pixel.Or("DiffScroll Width", o.Width, 220)
	total := fmt.Sprintf(k.Msg("files_changed"), len(files))
	said = append(said, total+" +"+strconv.Itoa(add)+" -"+strconv.Itoa(del))
	// Count columns fit the totals, the widest of them.
	aw, dw := pixel.MicroWidth("+"+strconv.Itoa(add)), pixel.MicroWidth("-"+strconv.Itoa(del))
	bx := w - 8 - 24
	dx := bx - 6 - dw
	ax := dx - 4 - aw
	if ax-12 < 24 || pixel.Measure(total) > ax-12 {
		panic(fmt.Sprintf("agent: DiffScroll Width %d is too narrow", w))
	}
	h := 6 + 4 + (len(files)+1)*pixel.LineHeight + 4 + 6
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(strings.Join(said, "; ")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(4, 6, w-8, h-12, pixel.Lit)
		g.VLine(4, 6, h-12, pixel.Shade)
		g.VLine(w-5, 6, h-12, pixel.Shade)
		rod(g, 0, w)
		rod(g, h-6, w)
		row := func(y int, name string, a, d int) {
			g.Text(8, y, name, pixel.Ink, 1)
			plus, minus := "+"+strconv.Itoa(a), "-"+strconv.Itoa(d)
			g.Micro(ax+aw-pixel.MicroWidth(plus), y+6, plus, ink(a, pixel.Good), 1)
			g.Micro(dx+dw-pixel.MicroWidth(minus), y+6, minus, ink(d, pixel.Bad), 1)
			blocks(g, bx, y+6, a, d)
		}
		y := 10
		for _, f := range files {
			row(y, tail(f.Path, ax-12), f.Added, f.Removed)
			y += pixel.LineHeight
		}
		g.Checker(8, y+1, w-16, 1, pixel.Shade)
		row(y+4, total, add, del)
	})
}

// ink is s for a count, Faint for none.
func ink(n int, s pixel.Slot) pixel.Slot {
	if n == 0 {
		return pixel.Faint
	}
	return s
}

// rod is a scroll's wooden roller, capped past the paper.
func rod(g pixel.Canvas, y, w int) {
	g.Rect(1, y+1, w-2, 4, pixel.Shade)
	g.HLine(1, y, w-2, pixel.Edge)
	g.HLine(1, y+5, w-2, pixel.Edge)
	g.Rect(0, y, 3, 6, pixel.Edge)
	g.Rect(w-3, y, 3, 6, pixel.Edge)
}

// blocks shares five squares between added and removed.
func blocks(g pixel.Canvas, x, y, a, d int) {
	good := 0
	if a > 0 || d > 0 {
		good = int(math.Round(5 * float64(a) / (float64(a) + float64(d))))
	}
	for i := range 5 {
		ink := pixel.Bad
		switch {
		case a == 0 && d == 0:
			ink = pixel.Shade
		case i < good:
			ink = pixel.Good
		}
		g.Rect(x+i*5, y, 4, 4, ink)
	}
}

// tail cuts s from the left, so the name shows.
func tail(s string, width int) string {
	if pixel.Measure(s) <= width {
		return s
	}
	rs := []rune(s)
	for len(rs) > 0 && pixel.Measure("…"+string(rs)) > width {
		rs = rs[1:]
	}
	return "…" + string(rs)
}
