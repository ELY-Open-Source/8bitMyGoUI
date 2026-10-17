package chart

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type GaugeOptions struct {
	// Min and Max bound the dial; zeros mean 0–100.
	Min, Max float64
	// Size is the dial's width in dots; zero means 64.
	Size int
}

// Gauge needles a half dial and prints the value.
func Gauge(c *ui.Context, text string, value float64, o GaugeOptions) ui.Element {
	k := pixel.K(c)
	span([]float64{o.Min, o.Max, value})
	lo, hi := o.Min, o.Max
	if lo == 0 && hi == 0 {
		hi = 100
	}
	if !(lo < hi && value >= lo && value <= hi) {
		panic(fmt.Sprintf("chart: Gauge value %v outside %v–%v", value, lo, hi))
	}
	n := pixel.Or("Gauge Size", o.Size, 64)
	shown := num(value)
	t := (value - lo) / (hi - lo)
	return ui.Box(c).Size(k.D(n), k.D(n/2+10)).Role(ui.RoleMeter).Label(text).Description(shown).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		c, rad := n/2, n/2-1
		// The dial runs from nine o'clock, over twelve, to three.
		sector(g, c, c, rad-5, rad, 0.75, 1, pixel.Good)
		sector(g, c, c, rad-5, rad, 0, 0.125, pixel.Warn)
		sector(g, c, c, rad-5, rad, 0.125, 0.25, pixel.Bad)
		a := math.Pi * (1 - t)
		g.Line(c, c, c+int(math.Round(float64(rad-3)*math.Cos(a))), c-int(math.Round(float64(rad-3)*math.Sin(a))), pixel.Ink)
		g.Disc(c, c, 2, pixel.Edge)
		tiny(g, c-tinyW(shown)/2+1, c+5, shown, pixel.Ink)
	})
}
