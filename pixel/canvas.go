package pixel

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
)

// Canvas draws dots for one Draw call; coordinates are dots.
type Canvas struct {
	p      painter
	k      Kit
	gx, gy int
	W, H   int
	// remap swaps slots as they draw; nil keeps them.
	remap *[Slots]Slot
}

// painter is what Canvas needs of *ui.Painter.
type painter interface {
	Fill(r ui.Rect, c ui.Color, radius float32)
	Clip(r ui.Rect, radius float32, fn func())
	After(d time.Duration)
	Now() time.Time
}

func (k Kit) Canvas(p *ui.Painter, r ui.Rect) Canvas {
	gx := int(math.Round(float64(r.X / k.dot)))
	gy := int(math.Round(float64(r.Y / k.dot)))
	return Canvas{p: p, k: k, gx: gx, gy: gy, W: int(r.W / k.dot), H: int(r.H / k.dot)}
}

func (g Canvas) now() time.Time {
	if g.k.clock != nil {
		return g.k.clock()
	}
	return g.p.Now()
}

func (g Canvas) Rect(x, y, w, h int, s Slot) {
	if w <= 0 || h <= 0 {
		return
	}
	if g.remap != nil {
		s = g.remap[s]
	}
	d := g.k.dot
	g.p.Fill(ui.Rect{X: float32(g.gx+x) * d, Y: float32(g.gy+y) * d, W: float32(w) * d, H: float32(h) * d}, g.k.C(s), 0)
}

// Cycled draws band's slots where Cycle(band, step) moves them.
func (g Canvas) Cycled(band []Slot, step int) Canvas {
	var m [Slots]Slot
	for i := range m {
		m[i] = Slot(i)
	}
	to := Cycle(band, step)
	for i, s := range band {
		if s >= Slots || slices.Index(band, s) != i {
			panic(fmt.Sprintf("pixel: Cycled band %v repeats or overflows", band))
		}
		m[s] = to[i]
	}
	if g.remap != nil {
		for i := range m {
			m[i] = g.remap[m[i]]
		}
	}
	g.remap = &m
	return g
}

// Clip limits fn's drawing to a dot rectangle.
func (g Canvas) Clip(x, y, w, h int, fn func()) {
	d := g.k.dot
	g.p.Clip(ui.Rect{X: float32(g.gx+x) * d, Y: float32(g.gy+y) * d, W: float32(w) * d, H: float32(h) * d}, 0, fn)
}

func (g Canvas) Dot(x, y int, s Slot)      { g.Rect(x, y, 1, 1, s) }
func (g Canvas) HLine(x, y, w int, s Slot) { g.Rect(x, y, w, 1, s) }
func (g Canvas) VLine(x, y, h int, s Slot) { g.Rect(x, y, 1, h, s) }

// Line draws a Bresenham line, both ends included.
func (g Canvas) Line(x0, y0, x1, y1 int, s Slot) {
	dx, dy := abs(x1-x0), -abs(y1-y0)
	sx, sy := sign(x1-x0), sign(y1-y0)
	e := dx + dy
	var pts [][2]int
	for {
		pts = append(pts, [2]int{x0, y0})
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * e
		if e2 >= dy {
			e += dy
			x0 += sx
		}
		if e2 <= dx {
			e += dx
			y0 += sy
		}
	}
	g.points(pts, s)
}

// points merges dots in a row or column into rects.
func (g Canvas) points(pts [][2]int, s Slot) {
	if len(pts) == 0 {
		return
	}
	x0, y0, x1, y1 := pts[0][0], pts[0][1], pts[0][0], pts[0][1]
	for _, p := range pts[1:] {
		row := y0 == y1 && p[1] == y0 && (p[0] == x1+1 || p[0] == x0-1)
		col := x0 == x1 && p[0] == x0 && (p[1] == y1+1 || p[1] == y0-1)
		if row || col {
			x0, y0, x1, y1 = min(x0, p[0]), min(y0, p[1]), max(x1, p[0]), max(y1, p[1])
			continue
		}
		g.Rect(x0, y0, x1-x0+1, y1-y0+1, s)
		x0, y0, x1, y1 = p[0], p[1], p[0], p[1]
	}
	g.Rect(x0, y0, x1-x0+1, y1-y0+1, s)
}

// Circle draws a midpoint circle outline of radius r.
func (g Canvas) Circle(cx, cy, r int, s Slot) {
	rows := map[int][]int{}
	octants(r, func(x, y int) {
		for _, p := range [][2]int{{x, y}, {y, x}, {-x, y}, {-y, x}, {x, -y}, {y, -x}, {-x, -y}, {-y, -x}} {
			rows[cy+p[1]] = append(rows[cy+p[1]], cx+p[0])
		}
	})
	for y, xs := range rows {
		slices.Sort(xs)
		xs = slices.Compact(xs)
		start := 0
		for i := 1; i <= len(xs); i++ {
			if i == len(xs) || xs[i] != xs[i-1]+1 {
				g.HLine(xs[start], y, xs[i-1]-xs[start]+1, s)
				start = i
			}
		}
	}
}

// Disc fills a midpoint circle of radius r.
func (g Canvas) Disc(cx, cy, r int, s Slot) {
	half := make([]int, r+1)
	octants(r, func(x, y int) {
		half[y] = max(half[y], x)
		half[x] = max(half[x], y)
	})
	for dy := -r; dy <= r; dy++ {
		h := half[abs(dy)]
		g.HLine(cx-h, cy+dy, 2*h+1, s)
	}
}

func octants(r int, plot func(x, y int)) {
	if r < 0 {
		panic(fmt.Sprintf("pixel: negative radius %d", r))
	}
	x, y, e := r, 0, 1-r
	for x >= y {
		plot(x, y)
		y++
		if e < 0 {
			e += 2*y + 1
		} else {
			x--
			e += 2*(y-x) + 1
		}
	}
}

// Bevel lights the top-left edge and darkens the bottom-right.
func (g Canvas) Bevel(x, y, w, h int, light, dark Slot) {
	g.HLine(x, y, w-1, light)
	g.VLine(x, y+1, h-2, light)
	g.HLine(x+1, y+h-1, w-1, dark)
	g.VLine(x+w-1, y, h-1, dark)
}

type FrameStyle uint8

const (
	Raised FrameStyle = iota
	Sunken
	Flat
	Double
	Window
	Notched
)

// Frame draws an outlined, filled box in a style.
func (g Canvas) Frame(x, y, w, h int, style FrameStyle, fill Slot) {
	if w < 3 || h < 3 {
		panic(fmt.Sprintf("pixel: frame %d×%d smaller than 3×3", w, h))
	}
	notch := style == Window || style == Notched
	g.outline(x, y, w, h, Edge, notch)
	g.Rect(x+1, y+1, w-2, h-2, fill)
	switch style {
	case Raised:
		g.Bevel(x+1, y+1, w-2, h-2, Lit, Shade)
	case Sunken:
		g.Bevel(x+1, y+1, w-2, h-2, Shade, Lit)
	case Double:
		if w >= 5 && h >= 5 {
			g.outline(x+2, y+2, w-4, h-4, Edge, false)
		}
	case Window:
		g.outline(x+1, y+1, w-2, h-2, Lit, false)
		if w >= 5 && h >= 5 {
			g.outline(x+2, y+2, w-4, h-4, Shade, false)
		}
	case Flat, Notched:
	default:
		panic(fmt.Sprintf("pixel: unknown frame style %d", style))
	}
}

func (g Canvas) outline(x, y, w, h int, s Slot, notch bool) {
	n := 0
	if notch {
		n = 1
	}
	g.HLine(x+n, y, w-2*n, s)
	g.HLine(x+n, y+h-1, w-2*n, s)
	g.VLine(x, y+1, h-2, s)
	g.VLine(x+w-1, y+1, h-2, s)
}

// Text draws one line with its top-left at x, y.
func (g Canvas) Text(x, y int, s string, ink Slot, scale int) {
	checkScale(scale)
	for _, r := range s {
		gl := glyphOf(r)
		top := y + gl.top()*scale
		for row := range gl.h {
			g.bits(x+gl.x*scale, top+row*scale, gl.w, scale, ink, func(col int) bool { return gl.on(col, row) })
		}
		x += gl.adv * scale
	}
}

// TextBlend draws s, glyph dots dithered from top to bottom.
func (g Canvas) TextBlend(x, y int, s string, top, bottom Slot, scale int) {
	checkScale(scale)
	span := LineHeight * scale
	for _, r := range s {
		gl := glyphOf(r)
		gx := x + gl.x*scale
		for j := range gl.h * scale {
			yy := y + gl.top()*scale + j
			level := min(16, (yy-y)*17/span)
			row := bayer[(g.gy+yy)&3]
			for _, part := range [2]struct {
				s    Slot
				down bool
			}{{top, false}, {bottom, true}} {
				runs(gl.w*scale, func(i int) bool {
					return gl.on(i/scale, j/scale) && (row[(g.gx+gx+i)&3] < level) == part.down
				}, func(a, b int) { g.Rect(gx+a, yy, b-a, 1, part.s) })
			}
		}
		x += gl.adv * scale
	}
}

// Micro draws s in the 3×5 micro font.
func (g Canvas) Micro(x, y int, s string, ink Slot, scale int) {
	checkScale(scale)
	for _, r := range s {
		gl := microGlyph(r)
		for row, bits := range gl {
			g.bits(x, y+row*scale, 3, scale, ink, func(col int) bool { return bits&(4>>col) != 0 })
		}
		x += MicroAdvance * scale
	}
}

// bits draws one row's set dots as merged runs.
func (g Canvas) bits(x, y, w, scale int, s Slot, on func(int) bool) {
	runs(w, on, func(i, j int) { g.Rect(x+i*scale, y, (j-i)*scale, scale, s) })
}

// runs emits each run [i, j) of set dots.
func runs(w int, on func(int) bool, emit func(i, j int)) {
	for i := 0; i < w; {
		if !on(i) {
			i++
			continue
		}
		j := i
		for j < w && on(j) {
			j++
		}
		emit(i, j)
		i = j
	}
}

func checkScale(scale int) {
	if scale < 1 || scale > 4 {
		panic(fmt.Sprintf("pixel: scale %d outside 1–4", scale))
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func sign(v int) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}
