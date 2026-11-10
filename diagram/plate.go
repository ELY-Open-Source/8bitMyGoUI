package diagram

import (
	"math"
	"slices"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// plate holds a slot per dot, 0 for none.
type plate struct {
	w, h int
	dots []uint8
}

// ink is a slot and its dither level, 1–16.
type ink struct {
	slot  pixel.Slot
	level int
}

func (p *plate) set(x, y int, k ink) {
	if x >= 0 && y >= 0 && x < p.w && y < p.h && pixel.Bayer(x, y) < k.level {
		p.dots[y*p.w+x] = uint8(k.slot) + 1
	}
}

// block inks a rectangle of dots, cut to the plate.
func (p *plate) block(x, y, w, h int, k ink) {
	for yy := max(y, 0); yy < min(y+h, p.h); yy++ {
		for xx := max(x, 0); xx < min(x+w, p.w); xx++ {
			p.set(xx, yy, k)
		}
	}
}

// fill inks the dots whose centers the outline winds around.
func (p *plate) fill(o outline, k ink) {
	type cross struct {
		x   float64
		dir int
	}
	top, bottom := math.Inf(1), math.Inf(-1)
	for _, line := range o {
		for _, q := range line {
			top, bottom = min(top, q.y), max(bottom, q.y)
		}
	}
	var hits []cross
	for y := max(0, int(math.Floor(top))); y < p.h && float64(y) <= bottom; y++ {
		hits = hits[:0]
		yc := float64(y) + 0.5
		for _, line := range o {
			for i := range line {
				a, b := line[i], line[(i+1)%len(line)]
				if (a.y <= yc) == (b.y <= yc) {
					continue
				}
				dir := 1
				if b.y < a.y {
					dir = -1
				}
				hits = append(hits, cross{a.x + (yc-a.y)*(b.x-a.x)/(b.y-a.y), dir})
			}
		}
		slices.SortFunc(hits, func(a, b cross) int {
			if a.x < b.x {
				return -1
			}
			return 1
		})
		wind := 0
		for i, hit := range hits {
			if wind != 0 {
				for x := int(math.Ceil(max(hits[i-1].x, 0.5) - 0.5)); x < p.w && float64(x)+0.5 <= hit.x; x++ {
					p.set(x, y, k)
				}
			}
			wind += hit.dir
		}
	}
}

// dot floors a point; a hair under rounds up.
func dot(v float64) int { return int(math.Floor(v + 1e-6)) }

// stroke walks each polyline with a square brush.
func (p *plate) stroke(o outline, width int, dash []float64, offset float64, k ink) {
	total := 0.0
	for _, d := range dash {
		total += d
	}
	for _, line := range o {
		gone := offset
		for i := 1; i < len(line); i++ {
			a, b := line[i-1], line[i]
			x, y, x1, y1 := dot(a.x), dot(a.y), dot(b.x), dot(b.y)
			dx, dy := abs(x1-x), -abs(y1-y)
			sx, sy, n := sign(x1-x), sign(y1-y), max(abs(x1-x), abs(y1-y))
			// A dashed line's brush is butt-ended: it widens across only.
			wide, tall := width, width
			if total > 0 && dx >= -dy {
				wide = 1
			} else if total > 0 {
				tall = 1
			}
			mx, my := 0, 0
			for s, lean := 0, dx+dy; ; s++ {
				left, top := x-wide/2, y-tall/2
				switch {
				case total > 0 && !inked(dash, total, gone+float64(s)/float64(max(n, 1))*dist(a, b)):
				case total > 0 || s == 0 && i == 1:
					p.block(left, top, wide, tall, k)
				default:
					// A solid brush that moved inks its leading edges.
					if mx != 0 {
						p.block(left+(wide-1)*(mx+1)/2, top, 1, tall, k)
					}
					if my != 0 {
						p.block(left, top+(tall-1)*(my+1)/2, wide, 1, k)
					}
				}
				if x == x1 && y == y1 {
					break
				}
				twice := 2 * lean
				mx, my = 0, 0
				if twice >= dy {
					lean, x, mx = lean+dy, x+sx, sx
				}
				if twice <= dx {
					lean, y, my = lean+dx, y+sy, sy
				}
			}
			gone += dist(a, b)
		}
	}
}

func abs(v int) int { return max(v, -v) }

func sign(v int) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}

// inked reports whether a dash covers the length at.
func inked(dash []float64, total, at float64) bool {
	if len(dash)%2 == 1 {
		total *= 2
	}
	at = math.Mod(math.Mod(at, total)+total, total)
	for i := 0; ; i++ {
		if d := dash[i%len(dash)]; at < d {
			return i%2 == 0
		} else {
			at -= d
		}
	}
}

// cell is a rectangle of one slot, in dots.
type cell struct {
	x, y, w, h int
	slot       pixel.Slot
}

// cells covers the plate's inked dots with rectangles.
func (p *plate) cells() []cell {
	var out []cell
	open := map[[3]int]int{}
	for y := range p.h {
		next := map[[3]int]int{}
		for x := 0; x < p.w; {
			v := p.dots[y*p.w+x]
			from := x
			for x < p.w && p.dots[y*p.w+x] == v {
				x++
			}
			if v == 0 {
				continue
			}
			key := [3]int{from, x, int(v)}
			if i, ok := open[key]; ok {
				out[i].h++
				next[key] = i
			} else {
				next[key] = len(out)
				out = append(out, cell{from, y, x - from, 1, pixel.Slot(v - 1)})
			}
		}
		open = next
	}
	return out
}
