package chart

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Candle is a period's open, high, low and close.
type Candle struct {
	Label                  string
	Open, High, Low, Close float64
}

type CandlestickOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
}

// Candlestick draws wicks and bodies: up Good, down Bad.
func Candlestick(c *ui.Context, text string, candles []Candle, o CandlestickOptions) *ui.Element {
	k := pixel.K(c)
	if len(candles) == 0 {
		panic("chart: Candlestick has no candles")
	}
	w, h := size("Candlestick", o.Width, o.Height, 160, 80)
	labels, lows, highs := make([]string, len(candles)), make([]float64, len(candles)), make([]float64, len(candles))
	for i, cd := range candles {
		span([]float64{cd.Open, cd.High, cd.Low, cd.Close})
		if cd.Low > min(cd.Open, cd.Close) || cd.High < max(cd.Open, cd.Close) {
			panic(fmt.Sprintf("chart: candle %q has open %v and close %v outside %v–%v", cd.Label, cd.Open, cd.Close, cd.Low, cd.High))
		}
		labels[i], lows[i], highs[i] = cd.Label, cd.Low, cd.High
	}
	lo, _ := span(lows)
	_, hi := span(highs)
	p := fit(lo, hi, w, h, labels)
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(ohlc(candles)).Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		p.grid(g, labels)
		bw := max(1, p.w/len(candles)/2)
		for i, cd := range candles {
			x, ink := p.slot(i, len(candles)), pixel.Good
			if cd.Close < cd.Open {
				ink = pixel.Bad
			}
			g.VLine(x, p.at(cd.High), p.at(cd.Low)-p.at(cd.High)+1, pixel.Edge)
			top, bottom := p.at(max(cd.Open, cd.Close)), p.at(min(cd.Open, cd.Close))
			g.Rect(x-bw/2, top, max(1, bw), max(1, bottom-top), ink)
		}
	})
}

// ohlc lists each candle's open, high, low and close.
func ohlc(candles []Candle) string {
	parts := make([]string, len(candles))
	for i, cd := range candles {
		parts[i] = cd.Label + " " + nums([]float64{cd.Open, cd.High, cd.Low, cd.Close})
	}
	return strings.Join(parts, "; ")
}
