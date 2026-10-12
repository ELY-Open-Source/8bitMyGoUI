package chart

import (
	"fmt"
	"strings"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Point sits on two value axes; Size sizes bubbles.
type Point struct {
	X, Y, Size float64
}

// field is a plot with a value axis across too.
type field struct {
	plot
	xfrom, xstep float64
	xn           int
}

// fitPoints fits both axes around points, marks pad dots in.
func fitPoints(points []Point, w, h, pad int) field {
	if len(points) == 0 {
		panic("chart: no points to plot")
	}
	xs, ys := make([]float64, len(points)), make([]float64, len(points))
	for i, pt := range points {
		xs[i], ys[i] = pt.X, pt.Y
	}
	xlo, xhi := span(xs)
	ylo, yhi := span(ys)
	f := field{plot: fit(ylo, yhi, w, h, []string{"0"})}
	f.pad = pad
	if f.w-1-2*pad < 2 || f.h-1-2*pad < 2 {
		panic(fmt.Sprintf("chart: a %d×%d chart leaves no room inside %d-dot margins", w, h, pad))
	}
	f.xfrom, f.xstep, f.xn = ticks(xlo, xhi, 5)
	return f
}

// ax is the x of v.
func (f field) ax(v float64) int {
	top := f.xfrom + float64(f.xn-1)*f.xstep
	return f.x + f.pad + int((v-f.xfrom)/(top-f.xfrom)*float64(f.w-1-2*f.pad)+0.5)
}

// grid draws both axes with their ticks.
func (f field) grid(g pixel.Canvas) {
	f.plot.grid(g, nil)
	for i := range f.xn {
		v := f.xfrom + float64(i)*f.xstep
		x, s := f.ax(v), tick(v, f.xstep)
		g.Dot(x, f.y+f.h+1, pixel.Edge)
		tiny(g, max(0, min(g.W-tinyW(s), x-tinyW(s)/2)), f.y+f.h+3, s, pixel.Faint)
	}
}

// where lists points as "(x, y)", sizes after for bubbles.
func where(points []Point, sized bool) string {
	parts := make([]string, len(points))
	for i, p := range points {
		parts[i] = "(" + num(p.X) + ", " + num(p.Y) + ")"
		if sized {
			parts[i] += " " + num(p.Size)
		}
	}
	return strings.Join(parts, " ")
}
