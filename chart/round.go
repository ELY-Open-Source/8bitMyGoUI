package chart

import (
	"fmt"
	"math"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// turn is a dot's angle, 0–1 clockwise from twelve.
func turn(dx, dy int) float64 {
	a := math.Atan2(float64(dx)+0.5, -float64(dy)-0.5) / (2 * math.Pi)
	if a < 0 {
		a++
	}
	return a
}

// sector fills dots between radii r0–r1 and turns a0–a1.
func sector(g pixel.Canvas, cx, cy, r0, r1 int, a0, a1 float64, s pixel.Slot) {
	for dy := -r1; dy <= r1; dy++ {
		for dx := -r1; dx <= r1; dx++ {
			d := dx*dx + dy*dy
			if d > r1*r1+r1 || d < r0*r0 {
				continue
			}
			if a := turn(dx, dy); a >= a0 && a < a1 {
				g.Dot(cx+dx, cy+dy, s)
			}
		}
	}
}

// shares are slices' fractions of their total, checked.
func shares(slices []Slice) (total float64, cuts []float64) {
	for _, s := range slices {
		if s.Value < 0 || math.IsNaN(s.Value) || math.IsInf(s.Value, 0) {
			panic(fmt.Sprintf("chart: slice %q has value %v", s.Name, s.Value))
		}
		total += s.Value
	}
	if total == 0 {
		panic("chart: slices sum to zero")
	}
	at := 0.0
	cuts = []float64{0}
	for _, s := range slices {
		at += s.Value / total
		cuts = append(cuts, at)
	}
	cuts[len(cuts)-1] = 1
	return total, cuts
}

// readSlices lists slices for assistive technology.
func readSlices(slices []Slice) string {
	labels, values := make([]string, len(slices)), make([]float64, len(slices))
	for i, s := range slices {
		labels[i], values[i] = s.Name, s.Value
	}
	return read(labels, values)
}
