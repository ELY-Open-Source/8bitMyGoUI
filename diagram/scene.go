package diagram

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/egoist/mygo/ui"
)

// scene is a diagram drawn on dots.
type scene struct {
	plate
	title string
}

// frame is an open element.
type frame struct {
	name string
	st   style
	geo  map[string]string
}

// word is text in one ink and spacing.
type word struct {
	text    string
	k       ink
	scale   int
	gap     int
	shift   int
	under   bool
	struck  bool
	outline bool
}

// line is a line of text being read.
type line struct {
	x, y   float32
	anchor string
	words  []word
}

// walker draws Pele's SVG onto a plate.
type walker struct {
	pal   palette
	slots *[pixel.Slots]ui.Color
	plate plate
	title string
	stack []frame
	lines []line
}

// parse reads pele-go's SVG onto dots in the palette's slots.
func parse(svg string, pal palette, slots *[pixel.Slots]ui.Color) (*scene, error) {
	w := &walker{pal: pal, slots: slots}
	root := style{fill: ui.RGB(0, 0, 0), fillA: 1, strokeA: 1, alpha: 1, width: 1, size: glyph, root: glyph,
		anchor: "start", t: identity}
	dec := xml.NewDecoder(strings.NewReader(svg))
	for {
		tok, err := dec.Token()
		if err == io.EOF && w.plate.w == 0 {
			return nil, fmt.Errorf("diagram: svg has no size")
		}
		if err == io.EOF {
			return &scene{w.plate, w.title}, nil
		}
		if err == nil {
			switch t := tok.(type) {
			case xml.StartElement:
				err = w.open(t, root)
			case xml.EndElement:
				err = w.close()
			case xml.CharData:
				w.chars(string(t))
			}
		}
		if err != nil {
			return nil, fmt.Errorf("diagram: %w", err)
		}
	}
}

func (w *walker) top() *frame { return &w.stack[len(w.stack)-1] }

func (w *walker) open(el xml.StartElement, root style) error {
	st, outer := root, len(w.stack) == 0
	if !outer {
		st = w.top().st
	}
	from := st
	geo, css, decls := map[string]string{}, "", [][2]string{}
	// add keeps a value; CSS drops a bad corner.
	add := func(name, value string, styled bool) error {
		if name != "fill" && name != "stroke" {
			var err error
			if value, err = w.plain(value); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
		}
		_, round := st.corner(value, 0)
		switch corner := name == "rx" || name == "ry"; {
		case styled && corner && !round:
		case corner || geometry[name] && !styled:
			geo[name] = value
		default:
			decls = append(decls, [2]string{name, value})
		}
		return nil
	}
	for _, a := range el.Attr {
		n := a.Name.Local
		if n == "style" {
			css = a.Value
		} else if !ignored[n] {
			if err := add(n, a.Value, false); err != nil {
				return err
			}
		}
	}
	// The outer style holds Pele's tokens alone.
	if !outer {
		err := declarations(css, func(name, value string) error { return add(name, value, true) })
		if err != nil {
			return err
		}
	}
	// Font size goes first: em lengths read it.
	for _, sized := range []bool{true, false} {
		for _, d := range decls {
			if (d[0] == "font-size") != sized {
				continue
			}
			if err := w.set(&st, &from, d[0], d[1]); err != nil {
				return err
			}
		}
	}
	w.stack = append(w.stack, frame{el.Name.Local, st, geo})
	switch el.Name.Local {
	case "svg":
		return w.svg(outer)
	case "g", "title", "desc":
		return nil
	case "rect", "circle", "ellipse", "polygon", "path":
		return w.shape()
	case "text", "tspan":
		return w.breakLine(el.Name.Local == "text")
	}
	return fmt.Errorf("unsupported element <%s>", el.Name.Local)
}

// svg sizes the plate, or places an icon's own box.
func (w *walker) svg(outer bool) error {
	f := w.top()
	n, err := f.numbers("x", "y", "width", "height")
	if err != nil {
		return err
	}
	if outer {
		if n[2] < 1 || n[3] < 1 || n[2] > limit || n[3] > limit {
			return fmt.Errorf("svg sized %g by %g, want 1 to %d dots a side", n[2], n[3], limit)
		}
		w.plate.w, w.plate.h = int(math.Ceil(float64(n[2]))), int(math.Ceil(float64(n[3])))
		w.plate.dots = make([]uint8, w.plate.w*w.plate.h)
		return nil
	}
	box, ok := lengths(f.geo["viewBox"])
	if !ok || len(box) != 4 || box[2] <= 0 || box[3] <= 0 {
		return fmt.Errorf("unsupported viewBox %q", f.geo["viewBox"])
	}
	f.st.t = translation(-float64(box[0]), -float64(box[1])).then(scaling(float64(n[2]/box[2]), float64(n[3]/box[3]))).
		then(translation(float64(n[0]), float64(n[1]))).then(f.st.t)
	return nil
}

// limit is the longest side of a plate, in dots.
const limit = 4096

// dots rounds a length; past two plates nothing more shows.
func dots(v float64) int { return int(math.Round(max(-2*limit, min(v, 2*limit)))) }

// ink is the slot nearest c, dithered by its alpha.
func (w *walker) ink(c ui.Color, alpha float32) ink {
	best, at := math.MaxInt, pixel.Ink
	for s, have := range w.slots {
		dr, dg, db := int(c.R)-int(have.R), int(c.G)-int(have.G), int(c.B)-int(have.B)
		if d := dr*dr + dg*dg + db*db; d < best {
			best, at = d, pixel.Slot(s)
		}
	}
	return ink{at, int(math.Round(float64(alpha) * float64(c.A) / 255 * 16))}
}

// shape fills and strokes the element on top.
func (w *walker) shape() error {
	f := w.top()
	d, err := f.data()
	if err != nil {
		return err
	}
	st := &f.st
	o, err := flatten(d, st.t)
	if err != nil {
		return fmt.Errorf("<%s>: %w", f.name, err)
	}
	if k := w.ink(st.fill, st.fillA*st.alpha); k.level > 0 {
		w.plate.fill(o, k)
	}
	grow := math.Sqrt(math.Abs(st.t.a*st.t.d - st.t.b*st.t.c))
	if k := w.ink(st.stroke, st.strokeA*st.alpha); k.level > 0 && st.width > 0 {
		dash := make([]float64, len(st.dash))
		for i, v := range st.dash {
			dash[i] = float64(v) * grow
		}
		w.plate.stroke(o, max(1, dots(float64(st.width)*grow)), dash, float64(st.dashOff)*grow, k)
	}
	return nil
}

// breakLine starts a line at a text or placed tspan.
func (w *walker) breakLine(text bool) error {
	f := w.top()
	if text {
		w.lines = w.lines[:0]
	} else if f.geo["x"] == "" && f.geo["y"] == "" {
		return nil
	}
	n, err := f.numbers("x", "y")
	if err != nil {
		return err
	}
	w.lines = append(w.lines, line{x: n[0], y: n[1], anchor: f.st.anchor})
	return nil
}

// chars keeps a title, or adds words to the line.
func (w *walker) chars(s string) {
	if len(w.stack) == 0 {
		return
	}
	f := w.top()
	switch {
	case f.name == "title" && len(w.stack) == 2:
		w.title += s
	case f.st.size > 0 && (f.name == "text" || f.name == "tspan"):
		st := &f.st
		wd := word{k: w.ink(st.fill, st.fillA*st.alpha), scale: scale(st.size), gap: dots(float64(st.spacing)),
			shift: dots(float64(st.shift)), under: st.underline, struck: st.struck}
		if wd.k.level == 0 {
			wd.k, wd.outline = w.ink(st.stroke, st.strokeA*st.alpha), true
		}
		l := &w.lines[len(w.lines)-1]
		for _, r := range s {
			wd.text = string(r)
			each := wd
			if r == ' ' || r == nbsp {
				each.gap += dots(float64(st.word))
			}
			l.words = append(l.words, each)
		}
	}
}

// nbsp is a word separator, as the space is.
const nbsp = '\u00a0'

// scale is the dots a font size gives a glyph's.
func scale(size float32) int { return max(1, min(4, int(math.Round(float64(size)/glyph)))) }

// advance is a word's width in dots.
func (wd word) advance() int { return pixel.Measure(wd.text)*wd.scale + wd.gap }

func (w *walker) close() error {
	f := w.top()
	if f.name == "text" {
		w.text(&f.st)
	}
	w.stack = w.stack[:len(w.stack)-1]
	return nil
}

// text stamps the lines read, placed by their anchors.
func (w *walker) text(st *style) {
	for _, l := range w.lines {
		width := 0
		for _, wd := range l.words {
			width += wd.advance()
		}
		x := float64(l.x)
		switch l.anchor {
		case "middle":
			x -= float64(width) / 2
		case "end":
			x -= float64(width)
		}
		// A quarter turn keeps whole dots; other turns stay upright.
		m := st.t
		at := m.at(pt{x, float64(l.y)})
		if !quarter(m) {
			m = identity
			at = st.t.at(pt{float64(l.x), float64(l.y)})
			at.x += x - float64(l.x)
		}
		ox, oy, pen := int(math.Round(at.x)), int(math.Round(at.y)), 0
		put := func(dx, dy int, k ink) {
			w.plate.set(ox+int(math.Round(m.a*float64(dx)+m.c*float64(dy))), oy+int(math.Round(m.b*float64(dx)+m.d*float64(dy))), k)
		}
		for _, wd := range l.words {
			wd.stamp(pen, put)
			pen += wd.advance()
		}
	}
}

// quarter reports a transform that turns by right angles only.
func quarter(m affine) bool {
	unit := func(v float64) bool { return math.Abs(v) < 1e-9 || math.Abs(math.Abs(v)-1) < 1e-9 }
	return unit(m.a) && unit(m.b) && unit(m.c) && unit(m.d)
}

// stamp plots a word from pen; outlines ring its dots.
func (wd word) stamp(pen int, put func(dx, dy int, k ink)) {
	pixel.Dots(wd.text, func(x, y int) {
		for sy := range wd.scale {
			for sx := range wd.scale {
				dx, dy := pen+x*wd.scale+sx, y*wd.scale+sy-wd.shift
				if !wd.outline {
					put(dx, dy, wd.k)
					continue
				}
				for _, by := range [8][2]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}} {
					put(dx+by[0], dy+by[1], wd.k)
				}
			}
		}
	})
	for i := range wd.advance() {
		if wd.under {
			put(pen+i, 2*wd.scale-wd.shift, wd.k)
		}
		if wd.struck {
			put(pen+i, -4*wd.scale-wd.shift, wd.k)
		}
	}
}
