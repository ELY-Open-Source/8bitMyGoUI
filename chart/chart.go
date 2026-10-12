// Package chart plots numbers in dots: bars, lines, pies.
package chart

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Series is a named row of values.
type Series struct {
	Name   string
	Values []float64
	// Slot colors it; zero takes the next series slot.
	Slot pixel.Slot
}

// Slice is a named share of a whole.
type Slice struct {
	Name  string
	Value float64
	// Slot colors it; zero takes the next series slot.
	Slot pixel.Slot
}

// slots color series and slices in turn.
var slots = []pixel.Slot{pixel.Accent, pixel.Info, pixel.Good, pixel.Warn, pixel.Magic, pixel.Bad}

// pick is s, or the i-th series slot for zero.
func pick(s pixel.Slot, i int) pixel.Slot {
	if s != 0 {
		return s
	}
	return slots[i%len(slots)]
}

// size is w×h in dots, or the defaults for zeros.
func size(name string, w, h, dw, dh int) (int, int) {
	return pixel.Or(name+" Width", w, dw), pixel.Or(name+" Height", h, dh)
}

// num prints a value exactly, float noise rounded off.
func num(v float64) string {
	r, err := strconv.ParseFloat(strconv.FormatFloat(v, 'g', 12, 64), 64)
	if err != nil {
		panic(fmt.Sprintf("chart: %v does not print: %v", v, err))
	}
	return strconv.FormatFloat(r, 'f', -1, 64)
}

// tick prints a tick as precisely as its step needs.
func tick(v, step float64) string {
	if v == 0 {
		// Minus zero prints as zero.
		v = 0
	}
	unit, suffix := 1.0, ""
	switch a := math.Abs(v); {
	case a >= 1e6 && step >= 1e5:
		unit, suffix = 1e6, "M"
	case a >= 1e4 && step >= 1e2:
		unit, suffix = 1e3, "K"
	}
	d := max(0, -int(math.Floor(math.Log10(step/unit)+1e-9)))
	return strconv.FormatFloat(v/unit, 'f', d, 64) + suffix
}

// part is v's share of whole, zero of nothing.
func part(v, whole float64) float64 {
	if whole == 0 {
		return 0
	}
	return v / whole
}

// small is whether every label fits the micro font.
func small(labels ...string) bool {
	for _, l := range labels {
		if !pixel.MicroOK(l) {
			return false
		}
	}
	return true
}

// inkH is a label's ink height (font rows 3–14).
func inkH(small bool) int {
	if small {
		return pixel.MicroHeight
	}
	return 12
}

// mark writes s with its ink's top at y.
func mark(g pixel.Canvas, x, y int, s string, small bool, ink pixel.Slot) {
	if small {
		g.Micro(x, y, s, ink, 1)
		return
	}
	g.Text(x, y-3, s, ink, 1)
}

// markW is s's width as mark writes it.
func markW(s string, small bool) int {
	if small {
		return pixel.MicroWidth(s)
	}
	return pixel.Measure(s)
}

// tiny centers s on a micro line's middle, y+2.
func tiny(g pixel.Canvas, x, y int, s string, ink pixel.Slot) {
	sm := small(s)
	mark(g, x, y+2-inkH(sm)/2, s, sm, ink)
}

// tinyW is s's width as tiny writes it.
func tinyW(s string) int { return markW(s, small(s)) }

// ticks steps nicely from below lo to above hi.
func ticks(lo, hi float64, want int) (from, step float64, n int) {
	// Flat data widens a tenth each way; zeros from zero.
	if hi <= lo {
		if lo == 0 {
			hi = 1
		} else {
			m := math.Abs(lo) / 10
			lo, hi = lo-m, hi+m
		}
	}
	step = nice((hi - lo) / float64(want))
	from = math.Floor(lo/step) * step
	n = int(math.Ceil((hi-from)/step-1e-9)) + 1
	return from, step, n
}

// nice rounds raw up to a 1-2-5 step.
func nice(raw float64) float64 {
	mag := math.Pow(10, math.Floor(math.Log10(raw)))
	for _, m := range []float64{1, 2, 5} {
		if m*mag >= raw {
			return m * mag
		}
	}
	return 10 * mag
}

// span is the least and most of every value.
func span(rows ...[]float64) (lo, hi float64) {
	lo, hi = math.Inf(1), math.Inf(-1)
	for _, r := range rows {
		for _, v := range r {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				panic(fmt.Sprintf("chart: value %v cannot be plotted", v))
			}
			lo, hi = min(lo, v), max(hi, v)
		}
	}
	if lo > hi {
		return 0, 1
	}
	return lo, hi
}

// nums lists values: "3, 9, 4".
func nums(values []float64) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = num(v)
	}
	return strings.Join(parts, ", ")
}
