package chart_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/chart"
)

func TestSpectrumPeaksStayWhenTheClockStepsBack(t *testing.T) {
	clk, levels := newClock(), []float64{1}
	tt := view(200, 100, clk, func(c *ui.Context) { chart.SpectrumBars(c, "EQ", levels, chart.SpectrumBarsOptions{}) })
	levels[0] = 0
	clk.t = clk.t.Add(200 * time.Millisecond)
	tt.Frame()
	before := bytes.Clone(tt.Image().Pix)
	clk.t = clk.t.Add(-time.Hour)
	tt.Frame()
	// A peak flung an hour high would leave the box.
	if !bytes.Equal(before, tt.Image().Pix) {
		t.Fatal("the falling peak moved as the clock stepped back")
	}
}

func TestSpectrumPeaksFallInWholeFrames(t *testing.T) {
	clk, levels := newClock(), []float64{1}
	tt := view(200, 400, clk, func(c *ui.Context) {
		chart.SpectrumBars(c, "Peaks", levels, chart.SpectrumBarsOptions{Segments: 100})
	})
	levels[0] = 0
	tt.Frame()
	seen := map[string]bool{}
	for range 7 {
		seen[string(tt.Image().Pix)] = true
		clk.t = clk.t.Add(10 * time.Millisecond)
		tt.Frame()
	}
	// S5: 60 ms holds two pictures at most.
	if len(seen) > 2 {
		t.Fatalf("%d pictures in 60 ms", len(seen))
	}
}
