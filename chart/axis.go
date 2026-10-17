package chart

import (
	"fmt"
	"math"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// plot is a chart's value range and box, in dots.
type plot struct {
	x, y, w, h    int
	from, step    float64
	n             int
	lowest, upper float64
	// pad keeps marks this many dots inside the box.
	pad int
}

// labelH is the room x labels need under a plot.
func labelH(labels []string) int { return inkH(small(labels...)) + 3 }

// fit lays lo–hi in a w×h canvas, leaving label room.
func fit(lo, hi float64, w, h int, labels []string) plot {
	from, step, n := ticks(lo, hi, 5)
	gutter := 0
	for i := range n {
		gutter = max(gutter, tinyW(tick(from+float64(i)*step, step)))
	}
	p := plot{x: gutter + 4, y: 3, from: from, step: step, n: n}
	p.w, p.h = w-p.x-1, h-p.y-labelH(labels)-2
	if p.w < 4 || p.h < 4 {
		panic(fmt.Sprintf("chart: a %d×%d chart leaves no room for its plot", w, h))
	}
	p.lowest, p.upper = from, from+float64(n-1)*step
	return p
}

// at is the y of v.
func (p plot) at(v float64) int {
	run := float64(p.h - 1 - 2*p.pad)
	return p.y + p.h - 1 - p.pad - int(math.Round((v-p.lowest)/(p.upper-p.lowest)*run))
}

// slot is the middle x of category i of n.
func (p plot) slot(i, n int) int { return p.x + (2*i+1)*p.w/(2*n) }

// grid draws value ticks, the frame and x labels.
func (p plot) grid(g pixel.Canvas, labels []string) {
	for i := range p.n {
		v := p.from + float64(i)*p.step
		y := p.at(v)
		g.Checker(p.x, y, p.w, 1, pixel.Shade)
		s := tick(v, p.step)
		tiny(g, p.x-3-tinyW(s), y-2, s, pixel.Faint)
	}
	g.VLine(p.x-1, p.y, p.h, pixel.Edge)
	g.HLine(p.x-1, p.y+p.h, p.w+1, pixel.Edge)
	sm := small(labels...)
	for i, l := range labels {
		x := p.slot(i, len(labels))
		g.Dot(x, p.y+p.h+1, pixel.Edge)
		mark(g, max(0, min(g.W-markW(l, sm), x-markW(l, sm)/2)), p.y+p.h+3, l, sm, pixel.Faint)
	}
}

type AxisOptions struct {
	// Min and Max bound ticks, zeros 0–100; Labels replace them.
	Min, Max float64
	Labels   []string
	// Length is in dots; zero means 120.
	Length   int
	Vertical bool
}

// Axis draws ticks and micro labels, for any chart.
func Axis(c *ui.Context, o AxisOptions) ui.Element {
	k := pixel.K(c)
	n := pixel.Or("Axis Length", o.Length, 120)
	marks := o.Labels
	if marks == nil {
		lo, hi := o.Min, o.Max
		if lo == 0 && hi == 0 {
			hi = 100
		}
		if span([]float64{lo, hi}); !(lo < hi) {
			panic(fmt.Sprintf("chart: Axis needs Min below Max, or Labels; got %v–%v", lo, hi))
		}
		from, step, count := ticks(lo, hi, 5)
		for i := range count {
			marks = append(marks, tick(from+float64(i)*step, step))
		}
	}
	if len(marks) < 2 {
		panic("chart: Axis needs two marks or more")
	}
	sm, thick := small(marks...), 0
	for _, m := range marks {
		thick = max(thick, markW(m, sm))
	}
	w, h := n, labelH(marks)+3
	if o.Vertical {
		w, h = thick+5, n
	}
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(marks[0] + "–" + marks[len(marks)-1]).Description(strings.Join(marks, ", ")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, m := range marks {
			at := 1 + i*(n-3)/(len(marks)-1)
			if o.Vertical {
				y := n - 1 - at
				g.HLine(w-3, y, 2, pixel.Edge)
				mark(g, w-4-markW(m, sm), max(0, min(n-inkH(sm), y-inkH(sm)/2)), m, sm, pixel.Faint)
				continue
			}
			g.VLine(at, 1, 2, pixel.Edge)
			mark(g, max(0, min(n-markW(m, sm), at-markW(m, sm)/2)), 4, m, sm, pixel.Faint)
		}
		if o.Vertical {
			g.VLine(w-1, 0, n, pixel.Edge)
		} else {
			g.HLine(0, 0, n, pixel.Edge)
		}
	})
}
