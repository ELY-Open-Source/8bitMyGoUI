package project

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type BurndownOptions struct {
	// Width and Height are dots; zeros mean 200 by 100.
	Width, Height int
}

// Burndown draws the ideal line and the work left.
func Burndown(c *ui.Context, text string, total float64, days int, left []float64, o BurndownOptions) ui.Element {
	k := pixel.K(c)
	if !(total > 0) || math.IsInf(total, 0) || days < 1 || len(left) > days+1 {
		panic(fmt.Sprintf("project: Burndown of %v over %d days, %d points", total, days, len(left)))
	}
	top := total
	for _, v := range left {
		if !(v >= 0) || math.IsInf(v, 0) {
			panic(fmt.Sprintf("project: Burndown work left %v", v))
		}
		top = max(top, v)
	}
	w, h := pixel.Or("Burndown Width", o.Width, 200), pixel.Or("Burndown Height", o.Height, 100)
	// The share comes first, so huge totals stay finite.
	ideal := func(d int) float64 { return total * (float64(days-d) / float64(days)) }
	verdict, ink := "", pixel.Good
	if n := len(left); n > 0 {
		verdict = k.Msg("ahead")
		if left[n-1] > ideal(n-1) {
			verdict, ink = k.Msg("behind"), pixel.Bad
		}
	}
	vals := make([]string, len(left))
	for i, v := range left {
		vals[i] = strconv.FormatFloat(v, 'g', 4, 64)
	}
	desc := fmt.Sprintf("%s %g → 0, %d; %s %s", k.Msg("ideal"), total, days, k.Msg("work_left"), strings.Join(vals, ", "))
	if verdict != "" {
		desc += "; " + verdict
	}
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(desc).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		// Plot inside a frame; day numbers sit under it.
		pw, ph := w-8, h-14
		// Floats, so huge day counts cannot overflow.
		px := func(d int) int { return 4 + int(float64(d)*float64(pw)/float64(days)) }
		py := func(v float64) int { return 4 + ph - int(v/top*float64(ph)) }
		g.Frame(0, 0, w, h-8, pixel.Sunken, pixel.Well)
		g.Line(px(0), py(total), px(days), py(0), pixel.Faint)
		for d := 1; d < len(left); d++ {
			g.Line(px(d-1), py(left[d-1]), px(d), py(left[d]), pixel.Accent)
		}
		for d, v := range left {
			g.Rect(px(d)-1, py(v)-1, 3, 3, pixel.Accent)
		}
		for _, d := range []int{0, days / 2, days} {
			s := strconv.Itoa(d)
			g.Micro(max(0, min(w-pixel.MicroWidth(s), px(d)-pixel.MicroWidth(s)/2)), h-6, s, pixel.Faint, 1)
		}
		if verdict != "" {
			g.Text(w-4-pixel.Measure(verdict), 3, verdict, ink, 1)
		}
	})
}
