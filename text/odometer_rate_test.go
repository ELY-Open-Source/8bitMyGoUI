package text_test

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

func TestOdometerHoldsMaxIntAndRollsWithinS5(t *testing.T) {
	clk := newClock()
	view(400, 100, clk, func(c *ui.Context) { text.Odometer(c, math.MaxInt, text.OdometerOptions{}) })
	value := 0
	tt := view(400, 100, clk, func(c *ui.Context) { text.Odometer(c, value, text.OdometerOptions{Digits: 3}) })
	value = 999
	tt.Frame()
	// S5: a second of rolling changes at most 15 times.
	changes, last := 0, append([]byte(nil), tt.Image().Pix...)
	for range 30 {
		clk.t = clk.t.Add(time.Second / 30)
		tt.Frame()
		if !bytes.Equal(last, tt.Image().Pix) {
			changes++
			last = append(last[:0], tt.Image().Pix...)
		}
	}
	if changes > 15 || changes == 0 {
		t.Fatalf("the odometer changed %d times in a second", changes)
	}
}

func TestOdometerShowsBigIntsExactly(t *testing.T) {
	clk := newClock()
	cell := func(v, digits int) []byte {
		tt := view(400, 100, clk, func(c *ui.Context) { text.Odometer(c, v, text.OdometerOptions{Digits: digits}) })
		img := tt.Image()
		// The last digit's 10-dot cell, at two pixels a dot.
		var out []byte
		for y := 0; y < 2*18; y++ {
			x0 := 2 * 10 * (digits - 1)
			out = append(out, img.Pix[y*img.Stride+4*x0:y*img.Stride+4*(x0+20)]...)
		}
		return out
	}
	// 2^53+1 ends in 3; a float rounds it away.
	if !bytes.Equal(cell(9007199254740993, 16), cell(3, 1)) {
		t.Fatal("9007199254740993 does not end in a drawn 3")
	}
}

func TestOdometerRollsHomeOnBigInts(t *testing.T) {
	clk := newClock()
	value := 0
	tt := view(400, 100, clk, func(c *ui.Context) { text.Odometer(c, value, text.OdometerOptions{Digits: 16}) })
	value = 9007199254740993
	for range 400 {
		clk.add(time.Second / 15)
		tt.Frame()
	}
	// Floats stall short of 2^53+1; the roll must still land.
	ref := view(400, 100, clk, func(c *ui.Context) { text.Odometer(c, 9007199254740993, text.OdometerOptions{Digits: 16}) })
	if !bytes.Equal(tt.Image().Pix, ref.Image().Pix) {
		t.Fatal("the roll never landed on 9007199254740993")
	}
}
