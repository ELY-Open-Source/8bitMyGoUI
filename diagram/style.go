package diagram

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
)

// style is what an element inherits, as SVG does.
type style struct {
	fill, stroke          ui.Color
	fillA, strokeA, alpha float32
	width                 float32
	dash                  []float32
	dashOff               float32
	size, root            float32
	spacing, word, shift  float32
	underline, struck     bool
	anchor                string
	t                     affine
}

// ignored are attributes that draw nothing.
var ignored = map[string]bool{"class": true, "id": true, "role": true, "xmlns": true, "space": true,
	"data-id": true, "data-icon": true, "aria-roledescription": true, "aria-labelledby": true, "aria-describedby": true}

// geometry are attributes a shape reads itself.
var geometry = map[string]bool{"x": true, "y": true, "width": true, "height": true, "rx": true, "ry": true,
	"cx": true, "cy": true, "r": true, "d": true, "points": true, "viewBox": true}

// wide are the keywords CSS takes for every property.
var wide = map[string]bool{"inherit": true, "initial": true, "unset": true, "revert": true, "revert-layer": true}

// units are CSS's absolute lengths in px.
var units = map[string]float32{"": 1, "px": 1, "pt": 4.0 / 3, "pc": 16, "in": 96, "cm": 96 / 2.54, "mm": 96 / 25.4, "q": 96 / 101.6}

// sizes are CSS's font-size keywords in px.
var sizes = map[string]float32{"xx-small": 9, "x-small": 10, "small": 13, "medium": 16, "large": 18, "x-large": 24, "xx-large": 32}

// measure reads a CSS length; em is st's font size.
func (st *style) measure(v string) (float32, bool) {
	at := strings.LastIndexAny(v, "0123456789.") + 1
	n, ok := number(v[:at])
	if !ok {
		return 0, false
	}
	switch unit := v[at:]; unit {
	case "em":
		return float32(n) * st.size, true
	case "rem":
		return float32(n) * st.root, true
	default:
		k, ok := units[unit]
		return float32(n) * k, ok
	}
}

// fraction reads a number or percentage, held to 0–1.
func fraction(v string) (float32, bool) {
	s, pct := strings.CutSuffix(v, "%")
	f, ok := number(s)
	if !ok {
		return 0, false
	}
	if pct {
		f /= 100
	}
	return float32(max(0, min(f, 1))), true
}

// set applies one property to st; from is inherited.
func (w *walker) set(st, from *style, name, value string) error {
	value = strings.TrimSpace(value)
	v := strings.ToLower(value)
	// A CSS-wide keyword is outside the subset: drop it.
	if wide[v] {
		value, v = "", ""
	}
	switch name {
	case "fill", "stroke":
		c, ok, err := w.pal.color(v, w.pal.fg)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		if ok && name == "fill" {
			st.fill = c
		} else if ok {
			st.stroke = c
		}
	case "fill-opacity", "stroke-opacity", "opacity":
		k, ok := fraction(v)
		switch {
		case !ok:
		case name == "fill-opacity":
			st.fillA = k
		case name == "stroke-opacity":
			st.strokeA = k
		default:
			st.alpha = from.alpha * k
		}
	case "stroke-width":
		if n, ok := st.measure(v); ok && n >= 0 {
			st.width = n
		}
	case "stroke-dasharray":
		st.dashes(v)
	case "stroke-dashoffset":
		if n, ok := st.measure(v); ok {
			st.dashOff = n
		}
	case "stroke-linecap", "stroke-linejoin":
		// A square brush of dots has no caps or joins.
	case "font-size":
		if n, ok := fontSize(v, from); ok {
			st.size = n
		}
	case "font-weight", "font-style", "font-family":
		// The pixel font has one face.
	case "text-decoration":
		under, struck := strings.Contains(v, "underline"), strings.Contains(v, "line-through")
		if under || struck || v == "none" || v == "overline" {
			st.underline, st.struck = under, struck
		}
	case "letter-spacing", "word-spacing":
		n, ok := st.measure(v)
		switch {
		case !ok && v != "normal":
		case name == "word-spacing":
			st.word = n
		default:
			st.spacing = n
		}
	case "baseline-shift":
		if k, ok := map[string]float32{"baseline": 0, "sub": -0.2, "super": 0.35}[v]; ok {
			st.shift = k * st.size
		}
	case "text-anchor":
		if v == "start" || v == "middle" || v == "end" {
			st.anchor = v
		}
	case "transform":
		t, err := transform(v)
		if err != nil {
			return err
		}
		st.t = t.then(from.t)
	default:
		return fmt.Errorf("unsupported attribute %q", name)
	}
	return nil
}

// dashes reads a dash array; a bad one is dropped.
func (st *style) dashes(v string) {
	var dash []float32
	if v != "none" {
		for _, f := range strings.FieldsFunc(v, func(r rune) bool { return r == ',' || r == ' ' }) {
			n, ok := st.measure(f)
			if !ok || n < 0 {
				return
			}
			dash = append(dash, n)
		}
	}
	st.dash = dash
}

// fontSize reads a font size against the inherited one.
func fontSize(v string, from *style) (float32, bool) {
	if n, ok := sizes[v]; ok {
		return n, true
	}
	switch {
	case v == "smaller":
		return from.size / 1.2, true
	case v == "larger":
		return from.size * 1.2, true
	case strings.HasSuffix(v, "%"):
		k, ok := number(v[:len(v)-1])
		return from.size * float32(k) / 100, ok && k >= 0
	}
	n, ok := from.measure(v)
	return n, ok && n >= 0
}

// transform reads a list of translate and rotate.
func transform(v string) (affine, error) {
	t := identity
	for _, call := range strings.SplitAfter(v, ")") {
		name, body, ok := strings.Cut(strings.TrimSpace(call), "(")
		if !ok {
			continue
		}
		f, ok := lengths(strings.TrimSuffix(body, ")"))
		n := make([]float64, len(f))
		for i, v := range f {
			n[i] = float64(v)
		}
		var op affine
		switch {
		case ok && name == "translate" && len(n) == 1:
			op = translation(n[0], 0)
		case ok && name == "translate" && len(n) == 2:
			op = translation(n[0], n[1])
		case ok && name == "rotate" && len(n) == 1:
			op = rotation(n[0])
		case ok && name == "rotate" && len(n) == 3:
			op = translation(-n[1], -n[2]).then(rotation(n[0])).then(translation(n[1], n[2]))
		default:
			return t, fmt.Errorf("unsupported transform %q", v)
		}
		t = op.then(t)
	}
	return t, nil
}

// declarations splits a style attribute.
func declarations(v string, each func(name, value string) error) error {
	for _, decl := range strings.Split(v, ";") {
		name, value, ok := strings.Cut(decl, ":")
		if !ok {
			continue
		}
		if err := each(strings.ToLower(strings.TrimSpace(name)), value); err != nil {
			return err
		}
	}
	return nil
}
