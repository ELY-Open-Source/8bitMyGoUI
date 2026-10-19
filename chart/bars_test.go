package chart_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/chart"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var days = []string{"MON", "TUE", "WED", "THU"}

func TestBarChartHeightsAndHover(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.BarChart(c, "Gold", days, []float64{10, 40, 20, 0}, chart.BarChartOptions{})
	})
	img := tt.Image()
	// Bars stand mid-category, past the 13-dot label gutter.
	var heights []int
	for i := range 4 {
		x := 13 + (2*i+1)*146/8
		_, n := column(t, img, x, 80, pixel.Accent)
		heights = append(heights, n)
	}
	// Topped at tick 40, that bar spans the 66-dot plot.
	if heights[1] != 66 || heights[2] != 33 || heights[0] < 16 || heights[0] > 17 || heights[3] != 1 {
		t.Fatalf("bar heights %v for 10, 40, 20, 0", heights)
	}
	keeps(t, img)
	tt.Move(2*(13+3*146/8), 2*40)
	if n := count(t, tt.Image(), 0, 0, 160, 80, pixel.Lit); n == 0 {
		t.Fatal("hovering a bar showed no readout")
	}
	// The value gutter left of the plot reads nothing.
	tt.Move(2*8, 2*40)
	if n := count(t, tt.Image(), 0, 0, 160, 80, pixel.Lit); n != 0 {
		t.Fatal("hovering the gutter showed a readout")
	}
	panics(t, "a short label list", func() {
		view(100, 100, newClock(), func(c *ui.Context) { chart.BarChart(c, "x", days[:1], []float64{1, 2}, chart.BarChartOptions{}) })
	})
}

func TestStackedBarStacksSeriesInOrder(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.StackedBar(c, "Loot", days[:2], []chart.Series{
			{Name: "Gold", Values: []float64{10, 5}, Slot: pixel.Warn},
			{Name: "Gems", Values: []float64{10, 15}, Slot: pixel.Info},
		}, chart.StackedBarOptions{})
	})
	img := tt.Image()
	x := 9 + 150/4
	warnTop, warn := column(t, img, x, 80, pixel.Warn)
	infoTop, info := column(t, img, x, 80, pixel.Info)
	if warn == 0 || info == 0 || infoTop >= warnTop || warn < info-3 || warn > info+3 {
		t.Fatalf("first stack: gold %d from %d, gems %d from %d", warn, warnTop, info, infoTop)
	}
	keeps(t, img)
	panics(t, "a negative stack", func() {
		view(100, 100, newClock(), func(c *ui.Context) {
			chart.StackedBar(c, "x", days[:1], []chart.Series{{Name: "a", Values: []float64{-1}}}, chart.StackedBarOptions{})
		})
	})
}

func TestHorizontalBarLongestReachesTheEdge(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.HorizontalBar(c, "Kills", []string{"Aria", "Brom"}, []float64{50, 100}, chart.HorizontalBarOptions{})
	})
	img := tt.Image()
	a, b := count(t, img, 0, 0, 160, 16, pixel.Accent), count(t, img, 0, 16, 160, 16, pixel.Accent)
	if a == 0 || b < 2*a-10 || b > 2*a+10 {
		t.Fatalf("bar dots %d and %d for 50 and 100", a, b)
	}
	keeps(t, img)
}

func TestHistogramBinsValues(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.Histogram(c, "Damage", []float64{1, 1, 1, 2, 2, 9}, chart.HistogramOptions{Bins: 2})
	})
	img := tt.Image()
	_, low := column(t, img, 40, 80, pixel.Accent)
	_, high := column(t, img, 120, 80, pixel.Accent)
	if low == 0 || high == 0 || low < 4*high-6 {
		t.Fatalf("bins of 5 and 1 stand %d and %d dots", low, high)
	}
}

func TestWaterfallInksStepsByDirection(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.Waterfall(c, "Gold", []string{"QUEST", "SHOP"}, []float64{30, -10}, chart.WaterfallOptions{})
	})
	img := tt.Image()
	if count(t, img, 0, 0, 80, 80, pixel.Good) == 0 || count(t, img, 80, 0, 80, 80, pixel.Bad) == 0 {
		t.Fatal("the up step is not Good, or the down step not Bad")
	}
	keeps(t, img)
}
