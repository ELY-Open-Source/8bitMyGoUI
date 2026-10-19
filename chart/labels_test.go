package chart_test

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/chart"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var zhDays = []string{"周一", "周二", "周三"}

func TestFontLabelsClearTheAxis(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.BarChart(c, "金币", zhDays, []float64{3, 9, 4}, chart.BarChartOptions{})
	})
	img := tt.Image()
	// With 15-row font labels the axis is row 63.
	if n := count(t, img, 10, 63, 149, 1, pixel.Edge); n != 149 {
		t.Fatalf("the axis row holds %d Edge dots of 149", n)
	}
	if count(t, img, 0, 64, 160, 16, pixel.Faint) == 0 || count(t, img, 11, 3, 148, 60, pixel.Faint) != 0 {
		t.Fatal("labels are not all below the axis")
	}
}

func TestHeatmapAndCrosshairGrowForFontLabels(t *testing.T) {
	tt := view(400, 300, newClock(), func(c *ui.Context) {
		chart.Heatmap(c, "热", []string{"周一", "周二"}, []string{"早", "晚上"}, [][]float64{{1, 2}, {3, 4}}, chart.HeatmapOptions{})
		chart.Crosshair(c, "血", zhDays, []chart.Series{{Name: "A", Values: []float64{1, 2, 3}}, {Name: "B", Values: []float64{3, 2, 1}}}, chart.CrosshairOptions{})
	})
	// Rows are 14 dots under a 15-dot label head.
	if r := find(t, tt, "热"); r.H != 2*(15+2*14) {
		t.Fatalf("heatmap is %v px tall", r.H)
	}
	r := find(t, tt, "血")
	tt.Move(r.X+2*80, r.Y+2*30)
	// The readout holds a 14-row label and two 7-row values.
	x, last := int(r.X)/2+95, -1
	top, _ := column(t, tt.Image(), x, 150, pixel.Lit)
	for y := range 150 {
		if at(t, tt.Image(), x, y) == pixel.Lit {
			last = y
		}
	}
	if top < 0 || last-top+1 < 3+14+2*7-2 {
		t.Fatalf("readout spans rows %d–%d", top, last)
	}
}

func TestTreemapSkipsEmptySlices(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.Treemap(c, "Bag", []chart.Slice{{Name: "a", Value: 1, Slot: pixel.Info}, {Name: "none", Slot: pixel.Magic},
			{Name: "b", Value: 3, Slot: pixel.Warn}, {Name: "zero", Slot: pixel.Magic}}, chart.TreemapOptions{})
	})
	img := tt.Image()
	a, none, b := count(t, img, 0, 0, 160, 80, pixel.Info), count(t, img, 0, 0, 160, 80, pixel.Magic), count(t, img, 0, 0, 160, 80, pixel.Warn)
	if none != 0 || a == 0 || b < 3*a-200 || b > 3*a+200 {
		t.Fatalf("dots: a %d, empty %d, b %d", a, none, b)
	}
}

func TestChartsRejectBadData(t *testing.T) {
	nan := math.NaN()
	for name, build := range map[string]func(c *ui.Context){
		"a widening funnel": func(c *ui.Context) {
			chart.Funnel(c, "x", []chart.Slice{{Name: "A", Value: 9}, {Name: "B", Value: 3}, {Name: "C", Value: 5}}, chart.FunnelOptions{})
		},
		"a NaN funnel stage": func(c *ui.Context) {
			chart.Funnel(c, "x", []chart.Slice{{Name: "A", Value: nan}}, chart.FunnelOptions{})
		},
		"no stacked labels":  func(c *ui.Context) { chart.StackedBar(c, "x", nil, nil, chart.StackedBarOptions{}) },
		"no box groups":      func(c *ui.Context) { chart.BoxPlot(c, "x", nil, chart.BoxPlotOptions{}) },
		"an axis of nothing": func(c *ui.Context) { chart.Axis(c, chart.AxisOptions{Min: 5, Max: 5}) },
		"a NaN bullet":       func(c *ui.Context) { chart.BulletChart(c, "x", nan, 1, []float64{2}, chart.BulletChartOptions{}) },
		"a NaN bubble": func(c *ui.Context) {
			chart.BubbleChart(c, "x", []chart.Point{{X: 1, Y: 1, Size: nan}}, chart.BubbleChartOptions{})
		},
		"a NaN candle": func(c *ui.Context) {
			chart.Candlestick(c, "x", []chart.Candle{{Label: "A", Open: nan, High: 2, Low: 0, Close: 1}}, chart.CandlestickOptions{})
		},
		"a negative radar": func(c *ui.Context) {
			chart.RadarChart(c, "x", []string{"A", "B", "C"}, []float64{1, -2, 3}, chart.RadarChartOptions{})
		},
		"a NaN radar Max": func(c *ui.Context) {
			chart.RadarChart(c, "x", []string{"A", "B", "C"}, []float64{1, 2, 3}, chart.RadarChartOptions{Max: nan})
		},
		"a negative radar Max": func(c *ui.Context) {
			chart.RadarChart(c, "x", []string{"A", "B", "C"}, []float64{1, 2, 3}, chart.RadarChartOptions{Max: -1})
		},
		"bars with no room": func(c *ui.Context) {
			chart.BarChart(c, "x", []string{"A"}, []float64{1}, chart.BarChartOptions{Width: 10})
		},
		"bubbles with no room": func(c *ui.Context) {
			chart.BubbleChart(c, "x", []chart.Point{{Size: 1}}, chart.BubbleChartOptions{Height: 20})
		},
		"an endless gauge": func(c *ui.Context) {
			chart.Gauge(c, "x", math.Inf(1), chart.GaugeOptions{Max: math.Inf(1)})
		},
	} {
		panics(t, name, func() { view(400, 200, newClock(), build) })
	}
}

func TestLabelsStayInsideTheChart(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: newClock().now, Locale: pixel.Zh})
		ui.Column(c).Children(func() {
			chart.ScatterPlot(c, "Hits", []chart.Point{{X: 100, Y: 1}, {X: 900, Y: 9}}, chart.ScatterPlotOptions{})
			chart.CalendarHeatmap(c, "Days", 2026, func(time.Time) int { return 0 }, chart.CalendarHeatmapOptions{})
			chart.BarChart(c, "Days", []string{"A", "WEDNESDAY"}, []float64{1, 2}, chart.BarChartOptions{Width: 60})
		})
	}, 800, 600)
	img := tt.Image()
	if n := count(t, img, 160, 0, 240, 80, pixel.Faint); n != 0 {
		t.Fatalf("%d scatter label dots past the right edge", n)
	}
	cal := find(t, tt, "Days 2026")
	// The year spans 53 weeks of 5 dots.
	if n := count(t, img, 265, int(cal.Y)/2, 135, 20, pixel.Faint); n != 0 {
		t.Fatalf("%d month label dots past the right edge", n)
	}
	bars := find(t, tt, "Days")
	if n := count(t, img, 60, int(bars.Y)/2, 340, 80, pixel.Faint); n != 0 {
		t.Fatalf("%d bar label dots past the right edge", n)
	}
}

func TestBulletTargetAtTheEndStaysInside(t *testing.T) {
	tt := view(400, 100, newClock(), func(c *ui.Context) {
		ui.Row(c).Children(func() {
			chart.BulletChart(c, "Goal", 50, 100, []float64{100}, chart.BulletChartOptions{Width: 100})
		})
	})
	if n := count(t, tt.Image(), 100, 0, 100, 14, pixel.Accent); n != 0 {
		t.Fatalf("%d target dots past the right edge", n)
	}
	if n := count(t, tt.Image(), 98, 2, 2, 10, pixel.Accent); n != 20 {
		t.Fatalf("the end target holds %d dots of 20", n)
	}
}

func TestTinyRangesKeepTheirShape(t *testing.T) {
	for _, scale := range []float64{1, 1e-5} {
		tt := view(400, 400, newClock(), func(c *ui.Context) {
			chart.ScatterPlot(c, "Hits", []chart.Point{{X: 1 * scale, Y: 1 * scale}, {X: 9 * scale, Y: 9 * scale}}, chart.ScatterPlotOptions{Slot: pixel.Magic})
			chart.Histogram(c, "Rolls", []float64{1 * scale, 1 * scale, 1 * scale, 2 * scale, 2 * scale, 9 * scale}, chart.HistogramOptions{Bins: 2, Slot: pixel.Info})
		})
		img := tt.Image()
		if n := count(t, img, 0, 0, 160, 80, pixel.Magic); n != 8 {
			t.Fatalf("scale %v: two points drew %d dots", scale, n)
		}
		_, low := column(t, img, 40, 160, pixel.Info)
		_, high := column(t, img, 120, 160, pixel.Info)
		if high == 0 || low < 4*high-6 {
			t.Fatalf("scale %v: bins of 5 and 1 stand %d and %d dots", scale, low, high)
		}
	}
}

func TestZeroRadarKeepsItsWeb(t *testing.T) {
	tt := view(200, 200, newClock(), func(c *ui.Context) {
		chart.RadarChart(c, "Newbie", []string{"STR", "DEX", "INT", "VIT", "AGI", "LUK"}, make([]float64, 6), chart.RadarChartOptions{})
	})
	if n := count(t, tt.Image(), 0, 0, 80, 80, pixel.Shade); n < 200 {
		t.Fatalf("the web holds %d dots", n)
	}
}

func TestBubblesStayInsideAShortChart(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		ui.Column(c).Padding(pixel.K(c).D(20)).Children(func() {
			chart.BubbleChart(c, "Towns", []chart.Point{{Size: 1}, {X: 10, Y: 10, Size: 1}}, chart.BubbleChartOptions{Height: 40})
		})
	})
	img := tt.Image()
	if n := 200*20 - count(t, img, 0, 0, 200, 20, pixel.Back); n != 0 {
		t.Fatalf("%d bubble dots above the chart", n)
	}
	if n := 20*40 - count(t, img, 180, 20, 20, 40, pixel.Back); n != 0 {
		t.Fatalf("%d bubble dots right of the chart", n)
	}
}

func TestTinyMagnitudesKeepTheirProportions(t *testing.T) {
	bubbles := map[float64]int{}
	for _, scale := range []float64{1, 1e-12} {
		tt := view(800, 400, newClock(), func(c *ui.Context) {
			ui.Row(c).Children(func() {
				chart.BubbleChart(c, "Towns", []chart.Point{{X: 0, Y: 0, Size: 1 * scale}, {X: 9, Y: 9, Size: 4 * scale}}, chart.BubbleChartOptions{Slot: pixel.Info})
				chart.HorizontalBar(c, "Kills", []string{"A", "B"}, []float64{1 * scale, 2 * scale}, chart.HorizontalBarOptions{Slot: pixel.Warn})
			})
			chart.Pyramid(c, "Ages", []string{"A"}, chart.Series{Name: "L", Values: []float64{1 * scale}, Slot: pixel.Magic},
				chart.Series{Name: "R", Values: []float64{2 * scale}, Slot: pixel.Good}, chart.PyramidOptions{})
			chart.Funnel(c, "Shop", []chart.Slice{{Name: "IN", Value: 4 * scale, Slot: pixel.Bad}, {Name: "OUT", Value: 1 * scale, Slot: pixel.Bad}}, chart.FunnelOptions{})
		})
		img, kills := tt.Image(), find(t, tt, "Kills")
		kx, ky := int(kills.X)/2, int(kills.Y)/2
		bubbles[scale] = count(t, img, 0, 0, 160, 80, pixel.Info)
		// One dot row through each bar gives its length.
		for _, bars := range [][2]int{
			{count(t, img, kx, ky+4, 160, 1, pixel.Warn), count(t, img, kx, ky+20, 160, 1, pixel.Warn)},
			{count(t, img, 0, 100, 200, 1, pixel.Magic), count(t, img, 0, 100, 200, 1, pixel.Good)},
			{count(t, img, 0, 131, 160, 1, pixel.Bad) * 2, count(t, img, 0, 115, 160, 1, pixel.Bad)},
		} {
			if bars[0] < 10 || bars[1]-2*bars[0] > 4 || 2*bars[0]-bars[1] > 4 {
				t.Errorf("scale %v: bars %v are not 1:2", scale, bars)
			}
		}
	}
	if bubbles[1] != bubbles[1e-12] || bubbles[1] == 0 {
		t.Fatalf("bubble dots %v", bubbles)
	}
}
