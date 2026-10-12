package chart

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CrosshairOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
}

// Crosshair is a line chart; hover crosses the nearest point.
func Crosshair(c *ui.Context, text string, labels []string, series []Series, o CrosshairOptions) *ui.Element {
	k := pixel.K(c)
	w, h := size("Crosshair", o.Width, o.Height, 160, 80)
	p := fitSeries(labels, series, w, h)
	e := ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(readSeries(labels, series))
	at, py := -1, 0
	if x, y, over := e.PointerPosition(); over && len(labels) > 0 {
		at = max(0, min(len(labels)-1, (int(x/k.Dot())-p.x)*len(labels)/max(1, p.w)))
		py = max(p.y, min(p.y+p.h-1, int(y/k.Dot())))
	}
	return e.Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		p.grid(g, labels)
		lines(g, p, len(labels), series)
		if at < 0 {
			return
		}
		x := p.slot(at, len(labels))
		g.Checker(x, p.y, 1, p.h, pixel.Ink)
		g.Checker(p.x, py, p.w, 1, pixel.Ink)
		// The readout lists the label, then each series' value.
		rows := []string{labels[at]}
		for _, s := range series {
			rows = append(rows, num(s.Values[at]))
		}
		bw, bh := 4, 3
		for _, row := range rows {
			bw, bh = max(bw, tinyW(row)+4), bh+inkH(small(row))+2
		}
		bx := min(w-bw, x+3)
		g.Frame(bx, p.y, bw, bh, pixel.Flat, pixel.Lit)
		y := p.y + 2
		for i, row := range rows {
			ink := pixel.Ink
			if i > 0 {
				ink = pick(series[i-1].Slot, i-1)
			}
			sm := small(row)
			mark(g, bx+2, y, row, sm, ink)
			y += inkH(sm) + 2
		}
	})
}
