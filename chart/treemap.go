package chart

import (
	"sort"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type TreemapOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
}

// box is a laid-out rectangle in dots.
type box struct{ x, y, w, h int }

// Treemap tiles slices by value in squarified rectangles.
func Treemap(c *ui.Context, text string, slices []Slice, o TreemapOptions) ui.Element {
	k := pixel.K(c)
	w, h := size("Treemap", o.Width, o.Height, 160, 80)
	total, _ := shares(slices)
	// Largest first keeps rows square; empty slices take no room.
	var order []int
	for i, s := range slices {
		if s.Value > 0 {
			order = append(order, i)
		}
	}
	sort.SliceStable(order, func(a, b int) bool { return slices[order[a]].Value > slices[order[b]].Value })
	areas := make([]float64, len(order))
	for i, j := range order {
		areas[i] = slices[j].Value / total * float64(w*h)
	}
	placed := squarify(areas, box{0, 0, w, h})
	boxes := make([]box, len(slices))
	for i, j := range order {
		boxes[j] = placed[i]
	}
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(readSlices(slices)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, b := range boxes {
			if b.w < 2 || b.h < 2 {
				continue
			}
			g.Rect(b.x, b.y, b.w-1, b.h-1, pick(slices[i].Slot, i))
			name := slices[i].Name
			if sm := small(name); markW(name, sm)+4 <= b.w && inkH(sm)+4 <= b.h {
				mark(g, b.x+2, b.y+2, name, sm, pixel.AccentInk)
			}
		}
	})
}

// squarify lays areas into b in rows kept near square.
func squarify(areas []float64, b box) []box {
	out := make([]box, 0, len(areas))
	for len(areas) > 0 {
		short := float64(min(b.w, b.h))
		n, worst := 1, aspect(areas[:1], short)
		for n < len(areas) {
			if next := aspect(areas[:n+1], short); next <= worst {
				n, worst = n+1, next
				continue
			}
			break
		}
		sum := 0.0
		for _, a := range areas[:n] {
			sum += a
		}
		// A row spans the short side; depth gives its area.
		if b.w >= b.h {
			depth := min(b.w, int(sum/float64(b.h)+0.5))
			if n == len(areas) {
				depth = b.w
			}
			y := b.y
			for i, a := range areas[:n] {
				hh := int(a/sum*float64(b.h) + 0.5)
				if i == n-1 {
					hh = b.y + b.h - y
				}
				out, y = append(out, box{b.x, y, depth, hh}), y+hh
			}
			b = box{b.x + depth, b.y, b.w - depth, b.h}
		} else {
			depth := min(b.h, int(sum/float64(b.w)+0.5))
			if n == len(areas) {
				depth = b.h
			}
			x := b.x
			for i, a := range areas[:n] {
				ww := int(a/sum*float64(b.w) + 0.5)
				if i == n-1 {
					ww = b.x + b.w - x
				}
				out, x = append(out, box{x, b.y, ww, depth}), x+ww
			}
			b = box{b.x, b.y + depth, b.w, b.h - depth}
		}
		areas = areas[n:]
	}
	return out
}

// aspect is a row's worst ratio laid along side.
func aspect(row []float64, side float64) float64 {
	sum, lo, hi := 0.0, row[0], row[0]
	for _, a := range row {
		sum, lo, hi = sum+a, min(lo, a), max(hi, a)
	}
	s2 := side * side
	return max(s2*hi/(sum*sum), sum*sum/(s2*lo))
}
