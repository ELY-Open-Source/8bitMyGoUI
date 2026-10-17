package chart

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type RadarChartOptions struct {
	// Max is the web's rim; zero takes the largest value.
	Max float64
	// Size is the web's width in dots; zero means 80.
	Size int
	Slot pixel.Slot
}

// RadarChart draws a stat polygon on a web of axes.
func RadarChart(c *ui.Context, text string, axes []string, values []float64, o RadarChartOptions) ui.Element {
	k := pixel.K(c)
	if len(axes) != len(values) || len(axes) < 3 {
		panic(fmt.Sprintf("chart: RadarChart has %d axes for %d values", len(axes), len(values)))
	}
	span([]float64{o.Max})
	lo, top := span(values)
	if lo < 0 || o.Max < 0 || o.Max > 0 && top > o.Max {
		panic(fmt.Sprintf("chart: RadarChart values %v–%v outside 0–%v", lo, top, o.Max))
	}
	if o.Max > 0 {
		top = o.Max
	}
	// All zeros: any rim works; the polygon sits mid-web.
	if top == 0 {
		top = 1
	}
	n := pixel.Or("RadarChart Size", o.Size, 80)
	sm, ink := small(axes...), pick(o.Slot, 0)
	rad := web(axes, sm, n)
	// point is where value v sits on axis i.
	point := func(i int, v float64) (int, int) {
		sin, cos := math.Sincos(2 * math.Pi * float64(i) / float64(len(axes)))
		reach := v / top * float64(rad)
		return n/2 + int(math.Round(reach*sin)), n/2 - int(math.Round(reach*cos))
	}
	return ui.Box(c).Size(k.D(n), k.D(n)).Role(ui.RoleImage).Label(text).Description(read(axes, values)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for _, ring := range []float64{1.0 / 3, 2.0 / 3, 1} {
			for i := range axes {
				x0, y0 := point(i, ring*top)
				x1, y1 := point((i+1)%len(axes), ring*top)
				g.Line(x0, y0, x1, y1, pixel.Shade)
			}
		}
		for i, a := range axes {
			x, y := point(i, top)
			g.Line(n/2, n/2, x, y, pixel.Shade)
			lx, ly := tag(i, len(axes), rad, markW(a, sm), inkH(sm), n)
			mark(g, lx, ly, a, sm, pixel.Faint)
		}
		for i, v := range values {
			x0, y0 := point(i, v)
			x1, y1 := point((i+1)%len(values), values[(i+1)%len(values)])
			g.Line(x0, y0, x1, y1, ink)
		}
		for i, v := range values {
			x, y := point(i, v)
			g.Rect(x-1, y-1, 3, 3, ink)
		}
	})
}

// gap is the room between the rim and a label.
const gap = 3

// lean is -1, 0 or 1: a label's side.
func lean(v float64) float64 {
	switch {
	case v > 0.3:
		return 1
	case v < -0.3:
		return -1
	}
	return 0
}

// web is the largest rim keeping every label inside.
func web(axes []string, sm bool, n int) int {
	rad := float64(n/2 - 1)
	for i, a := range axes {
		sin, cos := math.Sincos(2 * math.Pi * float64(i) / float64(len(axes)))
		w, h := float64(markW(a, sm)), float64(inkH(sm))
		// A leaning label reaches its whole size out; others half.
		ex, ey := w/2*(1+math.Abs(lean(sin))), h/2*(1+math.Abs(lean(cos)))
		if ex > float64(n)/2 || ey > float64(n)/2 {
			rad = 0
		}
		if s := math.Abs(sin); s > 1e-9 {
			rad = min(rad, (float64(n/2)-ex)/s-gap)
		}
		if c := math.Abs(cos); c > 1e-9 {
			rad = min(rad, (float64(n/2)-ey)/c-gap)
		}
	}
	if rad < 4 {
		panic(fmt.Sprintf("chart: RadarChart Size %d leaves no room inside its labels", n))
	}
	return int(rad)
}

// tag is the top left of axis i's label.
func tag(i, count, rad, w, h, n int) (int, int) {
	sin, cos := math.Sincos(2 * math.Pi * float64(i) / float64(count))
	ax, ay := float64(n)/2+float64(rad+gap)*sin, float64(n)/2-float64(rad+gap)*cos
	x := ax - float64(w)/2*(1-lean(sin))
	y := ay - float64(h)/2*(1+lean(cos))
	return max(0, min(n-w, int(math.Round(x)))), max(0, min(n-h, int(math.Round(y))))
}
