package diagram

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// pt is a point in dots.
type pt struct{ x, y float64 }

// affine maps x, y to a·x+c·y+e, b·x+d·y+f.
type affine struct{ a, b, c, d, e, f float64 }

var identity = affine{a: 1, d: 1}

func (m affine) at(p pt) pt { return pt{m.a*p.x + m.c*p.y + m.e, m.b*p.x + m.d*p.y + m.f} }

// then applies m first, then n.
func (m affine) then(n affine) affine {
	return affine{n.a*m.a + n.c*m.b, n.b*m.a + n.d*m.b, n.a*m.c + n.c*m.d, n.b*m.c + n.d*m.d,
		n.a*m.e + n.c*m.f + n.e, n.b*m.e + n.d*m.f + n.f}
}

func translation(x, y float64) affine { return affine{a: 1, d: 1, e: x, f: y} }

func scaling(x, y float64) affine { return affine{a: x, d: y} }

// rotation turns clockwise by degrees, as SVG's rotate.
func rotation(degrees float64) affine {
	s, c := math.Sincos(degrees * math.Pi / 180)
	return affine{a: c, b: s, c: -s, d: c}
}

// outline is polylines; a closed one repeats its start.
type outline [][]pt

// reader walks SVG path data.
type reader struct {
	d  string
	at int
}

func (r *reader) skip() {
	for r.at < len(r.d) && strings.ContainsRune(" ,\t\r\n", rune(r.d[r.at])) {
		r.at++
	}
}

// number reads one SVG number; numbers may touch.
func (r *reader) number() (float64, error) {
	r.skip()
	from, digits := r.at, func() {
		for r.at < len(r.d) && r.d[r.at] >= '0' && r.d[r.at] <= '9' {
			r.at++
		}
	}
	if r.at < len(r.d) && (r.d[r.at] == '+' || r.d[r.at] == '-') {
		r.at++
	}
	digits()
	if r.at < len(r.d) && r.d[r.at] == '.' {
		r.at++
		digits()
	}
	if r.at+1 < len(r.d) && (r.d[r.at] == 'e' || r.d[r.at] == 'E') && strings.ContainsRune("+-0123456789", rune(r.d[r.at+1])) {
		r.at += 2
		digits()
	}
	v, err := strconv.ParseFloat(r.d[from:r.at], 64)
	if err != nil || math.IsInf(v, 0) {
		return 0, fmt.Errorf("path data: no number at offset %d", from)
	}
	return v, nil
}

// flag reads an arc flag, one digit.
func (r *reader) flag() (bool, error) {
	r.skip()
	if r.at >= len(r.d) || r.d[r.at] != '0' && r.d[r.at] != '1' {
		return false, fmt.Errorf("path data: no arc flag at offset %d", r.at)
	}
	r.at++
	return r.d[r.at-1] == '1', nil
}

// numbers reads n numbers; relative ones start from cur.
func (r *reader) numbers(n int) ([]float64, error) {
	out := make([]float64, n)
	for i := range out {
		v, err := r.number()
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

// pen draws a path into an outline through a transform.
type pen struct {
	m          affine
	out        outline
	cur, start pt
	open       bool
}

func (p *pen) move(to pt) {
	p.cur, p.start, p.open = to, to, false
}

func (p *pen) line(to pt) {
	if !p.open {
		p.out, p.open = append(p.out, []pt{p.m.at(p.cur)}), true
	}
	last := len(p.out) - 1
	p.out[last], p.cur = append(p.out[last], p.m.at(to)), to
}

// curve samples f at about two dots a step.
func (p *pen) curve(reach float64, f func(t float64) pt) {
	n := max(2, min(64, int(math.Ceil(reach*math.Sqrt(math.Abs(p.m.a*p.m.d-p.m.b*p.m.c))/2))))
	for i := 1; i <= n; i++ {
		p.line(f(float64(i) / float64(n)))
	}
}

func dist(a, b pt) float64 { return math.Hypot(b.x-a.x, b.y-a.y) }

func (p *pen) cubic(c1, c2, to pt) {
	from := p.cur
	p.curve(dist(from, c1)+dist(c1, c2)+dist(c2, to), func(t float64) pt {
		u := 1 - t
		a, b, c, d := u*u*u, 3*u*u*t, 3*u*t*t, t*t*t
		return pt{a*from.x + b*c1.x + c*c2.x + d*to.x, a*from.y + b*c1.y + c*c2.y + d*to.y}
	})
}

// arc draws SVG's elliptical arc to to.
func (p *pen) arc(rx, ry, turn float64, large, sweep bool, to pt) {
	from := p.cur
	rx, ry = math.Abs(rx), math.Abs(ry)
	if rx == 0 || ry == 0 || from == to {
		p.line(to)
		return
	}
	sin, cos := math.Sincos(turn * math.Pi / 180)
	dx, dy := (from.x-to.x)/2, (from.y-to.y)/2
	x1, y1 := cos*dx+sin*dy, -sin*dx+cos*dy
	if grow := x1*x1/(rx*rx) + y1*y1/(ry*ry); grow > 1 {
		rx, ry = rx*math.Sqrt(grow), ry*math.Sqrt(grow)
	}
	k := math.Sqrt(max(0, rx*rx*ry*ry-rx*rx*y1*y1-ry*ry*x1*x1) / (rx*rx*y1*y1 + ry*ry*x1*x1))
	if large == sweep {
		k = -k
	}
	cx1, cy1 := k*rx*y1/ry, -k*ry*x1/rx
	cx, cy := cos*cx1-sin*cy1+(from.x+to.x)/2, sin*cx1+cos*cy1+(from.y+to.y)/2
	a0 := math.Atan2((y1-cy1)/ry, (x1-cx1)/rx)
	sweepBy := math.Atan2((-y1-cy1)/ry, (-x1-cx1)/rx) - a0
	if sweep && sweepBy < 0 {
		sweepBy += 2 * math.Pi
	} else if !sweep && sweepBy > 0 {
		sweepBy -= 2 * math.Pi
	}
	p.curve(math.Abs(sweepBy)*max(rx, ry), func(t float64) pt {
		s, c := math.Sincos(a0 + sweepBy*t)
		return pt{cx + cos*rx*c - sin*ry*s, cy + sin*rx*c + cos*ry*s}
	})
	p.cur = to
}

// arity counts a command's numbers; arcs read their own.
var arity = map[byte]int{'M': 2, 'L': 2, 'H': 1, 'V': 1, 'C': 6, 'S': 4, 'Q': 4, 'T': 2, 'A': 0, 'Z': 0}

// family is the curve a smooth command may mirror.
var family = map[byte]byte{'C': 'C', 'S': 'C', 'Q': 'Q', 'T': 'Q'}

// flatten turns SVG path data into polylines under m.
func flatten(d string, m affine) (outline, error) {
	r, p := reader{d: d}, pen{m: m}
	var cmd byte
	var hold pt
	var held byte
	for {
		r.skip()
		if r.at >= len(r.d) {
			return p.out, nil
		}
		if c := r.d[r.at]; c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' {
			cmd = c
			r.at++
		} else if cmd == 0 || cmd == 'Z' || cmd == 'z' {
			return nil, fmt.Errorf("path data: unexpected %q at offset %d", c, r.at)
		}
		base, op := pt{}, cmd&^0x20
		if cmd >= 'a' {
			base = p.cur
		}
		want, known := arity[op]
		if !known {
			return nil, fmt.Errorf("path data: unknown command %q", cmd)
		}
		n, err := r.numbers(want)
		if err != nil {
			return nil, err
		}
		at := func(i int) pt { return pt{base.x + n[i], base.y + n[i+1]} }
		mirror := pt{2*p.cur.x - hold.x, 2*p.cur.y - hold.y}
		switch op {
		case 'M':
			p.move(at(0))
			cmd--
		case 'L':
			p.line(at(0))
		case 'H':
			p.line(pt{base.x + n[0], p.cur.y})
		case 'V':
			p.line(pt{p.cur.x, base.y + n[0]})
		case 'C':
			hold = at(2)
			p.cubic(at(0), hold, at(4))
		case 'S':
			if held != 'C' {
				mirror = p.cur
			}
			hold = at(0)
			p.cubic(mirror, hold, at(2))
		case 'Q', 'T':
			to := at(want - 2)
			if hold = mirror; op == 'Q' {
				hold = at(0)
			} else if held != 'Q' {
				hold = p.cur
			}
			from := p.cur
			p.cubic(pt{from.x + 2*(hold.x-from.x)/3, from.y + 2*(hold.y-from.y)/3},
				pt{to.x + 2*(hold.x-to.x)/3, to.y + 2*(hold.y-to.y)/3}, to)
		case 'A':
			if err := r.arcTo(&p, base); err != nil {
				return nil, err
			}
		case 'Z':
			p.line(p.start)
			p.open = false
		}
		held = family[op]
	}
}

// arcTo reads one arc's arguments and draws it.
func (r *reader) arcTo(p *pen, base pt) error {
	n, err := r.numbers(3)
	if err != nil {
		return err
	}
	large, err := r.flag()
	if err != nil {
		return err
	}
	sweep, err := r.flag()
	if err != nil {
		return err
	}
	to, err := r.numbers(2)
	if err != nil {
		return err
	}
	p.arc(n[0], n[1], n[2], large, sweep, pt{base.x + to[0], base.y + to[1]})
	return nil
}
