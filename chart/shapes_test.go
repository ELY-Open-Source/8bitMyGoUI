package chart_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/chart"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var loot = []chart.Slice{{Name: "Gold", Value: 3, Slot: pixel.Warn}, {Name: "Gems", Value: 1, Slot: pixel.Info}}

func TestPieAndDonutShareBySlice(t *testing.T) {
	tt := view(300, 200, newClock(), func(c *ui.Context) {
		ui.Row(c).Children(func() {
			chart.PieChart(c, "Pie", loot, chart.PieChartOptions{})
			chart.DonutChart(c, "Donut", loot, chart.DonutChartOptions{})
		})
	})
	img := tt.Image()
	gold, gems := count(t, img, 0, 0, 64, 64, pixel.Warn), count(t, img, 0, 0, 64, 64, pixel.Info)
	if share := float64(gold) / float64(gold+gems); share < 0.72 || share > 0.78 {
		t.Fatalf("pie dots %d and %d for 3:1", gold, gems)
	}
	// The donut's middle is hollow.
	if at(t, img, 64+32, 26) == pixel.Warn || count(t, img, 64, 0, 64, 64, pixel.Warn) == 0 {
		t.Fatal("the donut has no hole, or no ring")
	}
	keeps(t, img)
	panics(t, "slices of nothing", func() {
		view(100, 100, newClock(), func(c *ui.Context) { chart.PieChart(c, "x", []chart.Slice{{Name: "a"}}, chart.PieChartOptions{}) })
	})
}

func TestPolarAreaReachesByValue(t *testing.T) {
	tt := view(200, 200, newClock(), func(c *ui.Context) {
		chart.PolarArea(c, "Stats", []chart.Slice{{Name: "A", Value: 4, Slot: pixel.Warn}, {Name: "B", Value: 2, Slot: pixel.Info}}, chart.PolarAreaOptions{})
	})
	img := tt.Image()
	a, b := count(t, img, 0, 0, 64, 64, pixel.Warn), count(t, img, 0, 0, 64, 64, pixel.Info)
	// Twice the reach is four times the area.
	if a < 3*b || a > 5*b {
		t.Fatalf("polar dots %d and %d for reaches 2:1", a, b)
	}
	tt = view(200, 200, newClock(), func(c *ui.Context) {
		chart.PolarArea(c, "Stats", []chart.Slice{{Name: "A", Value: 4, Slot: pixel.Warn}, {Name: "B", Slot: pixel.Info}}, chart.PolarAreaOptions{})
	})
	if n := count(t, tt.Image(), 0, 0, 64, 64, pixel.Info); n != 0 {
		t.Fatalf("a zero slice drew %d dots", n)
	}
}

func TestProgressRingFillsClockwise(t *testing.T) {
	tt := view(200, 100, newClock(), func(c *ui.Context) {
		chart.ProgressRing(c, "Load", 0.25, chart.ProgressRingOptions{Slot: pixel.Good})
	})
	img := tt.Image()
	if right, left := count(t, img, 20, 0, 20, 20, pixel.Good), count(t, img, 0, 0, 20, 40, pixel.Good); right == 0 || left != 0 {
		t.Fatalf("a quarter lit %d dots top right, %d on the left", right, left)
	}
	panics(t, "a value of 2", func() {
		view(100, 100, newClock(), func(c *ui.Context) { chart.ProgressRing(c, "x", 2, chart.ProgressRingOptions{}) })
	})
}

func TestGaugeNeedlePoints(t *testing.T) {
	for _, tc := range []struct {
		v    float64
		x, y int
	}{{0, 6, 32}, {100, 58, 32}, {50, 32, 6}} {
		tt := view(200, 100, newClock(), func(c *ui.Context) { chart.Gauge(c, "Speed", tc.v, chart.GaugeOptions{}) })
		if s := at(t, tt.Image(), tc.x, tc.y); s != pixel.Ink {
			t.Fatalf("value %v: no needle at %d,%d (%v)", tc.v, tc.x, tc.y, s)
		}
	}
}

func TestRadarChartReachesTheRim(t *testing.T) {
	for _, tc := range []struct {
		axes []string
		rim  int
	}{
		// 11-dot labels on the 60° axes leave a 30 rim.
		{[]string{"STR", "DEX", "INT", "VIT", "AGI", "LUK"}, 30},
		// Pixel-font labels, 24 wide, shrink it to 15.
		{[]string{"力量", "敏捷", "智力", "体力", "速度", "运气"}, 15},
	} {
		tt := view(200, 200, newClock(), func(c *ui.Context) {
			chart.RadarChart(c, "Hero", tc.axes, []float64{10, 5, 5, 5, 5, 5}, chart.RadarChartOptions{Slot: pixel.Magic})
		})
		if top, _ := column(t, tt.Image(), 40, 80, pixel.Magic); top != 40-tc.rim-1 {
			t.Fatalf("%s: the largest value's marker tops row %d", tc.axes[0], top)
		}
		keeps(t, tt.Image())
	}
	panics(t, "a value above Max", func() {
		view(200, 200, newClock(), func(c *ui.Context) {
			chart.RadarChart(c, "x", []string{"A", "B", "C"}, []float64{1, 2, 12}, chart.RadarChartOptions{Max: 10})
		})
	})
	panics(t, "labels too wide for the size", func() {
		view(200, 200, newClock(), func(c *ui.Context) {
			chart.RadarChart(c, "x", []string{"STRENGTH", "B", "C"}, []float64{1, 2, 3}, chart.RadarChartOptions{Size: 30})
		})
	})
}

func TestHeatmapsShadeByLevel(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.Heatmap(c, "Heat", []string{"A"}, []string{"X", "Y"}, [][]float64{{0, 10}}, chart.HeatmapOptions{})
		chart.CalendarHeatmap(c, "Raids", 2026, func(d time.Time) int {
			if d.Month() == time.January {
				return 4
			}
			return 0
		}, chart.CalendarHeatmapOptions{})
	})
	img := tt.Image()
	r := find(t, tt, "Heat")
	ox, oy := int(r.X)/2, int(r.Y)/2
	// Labels take six dots left, eight on top.
	low, high := count(t, img, ox+6, oy+8, 9, 9, pixel.Bad), count(t, img, ox+16, oy+8, 9, 9, pixel.Bad)
	if low != 0 || high != 81 {
		t.Fatalf("cells of 0 and 10 hold %d and %d Bad dots", low, high)
	}
	cal := find(t, tt, "Raids 2026")
	cx, cy := int(cal.X)/2, int(cal.Y)/2
	if jan, jun := count(t, img, cx, cy, 16, 28, pixel.Good), count(t, img, cx+100, cy, 16, 28, pixel.Good); jan == 0 || jun != 0 {
		t.Fatalf("January lit %d dots, June %d", jan, jun)
	}
	panics(t, "a level of 5", func() {
		view(400, 100, newClock(), func(c *ui.Context) {
			chart.CalendarHeatmap(c, "x", 2026, func(time.Time) int { return 5 }, chart.CalendarHeatmapOptions{})
		})
	})
}

func TestTreemapAreasFollowValues(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.Treemap(c, "Bag", []chart.Slice{{Name: "a", Value: 1, Slot: pixel.Info}, {Name: "b", Value: 3, Slot: pixel.Warn}}, chart.TreemapOptions{})
	})
	img := tt.Image()
	a, b := count(t, img, 0, 0, 160, 80, pixel.Info), count(t, img, 0, 0, 160, 80, pixel.Warn)
	if b < 3*a-200 || b > 3*a+200 {
		t.Fatalf("treemap dots %d and %d for 1:3", a, b)
	}
	keeps(t, img)
}

func TestBulletFunnelCandlesBoxes(t *testing.T) {
	tt := view(400, 400, newClock(), func(c *ui.Context) {
		chart.BulletChart(c, "Goal", 60, 80, []float64{50, 100}, chart.BulletChartOptions{Width: 100})
		chart.Funnel(c, "Sales", []chart.Slice{{Name: "SEEN", Value: 100, Slot: pixel.Info}, {Name: "BOUGHT", Value: 25, Slot: pixel.Info}}, chart.FunnelOptions{Width: 100})
		chart.Candlestick(c, "Gold", []chart.Candle{{Label: "A", Open: 10, High: 20, Low: 5, Close: 15}, {Label: "B", Open: 15, High: 16, Low: 2, Close: 4}}, chart.CandlestickOptions{})
	})
	img := tt.Image()
	if bar := count(t, img, 0, 7, 100, 1, pixel.Ink); bar < 58 || bar > 62 {
		t.Fatalf("a 60 of 100 bar is %d dots", bar)
	}
	if mark := at(t, img, 79, 7); mark != pixel.Accent {
		t.Fatalf("no target mark at 80: %v", mark)
	}
	f := find(t, tt, "Sales")
	fy := int(f.Y) / 2
	if top, bottom := count(t, img, 0, fy+2, 100, 1, pixel.Info), count(t, img, 0, fy+18, 100, 1, pixel.Info); top < 3*bottom {
		t.Fatalf("funnel rows %d and %d dots", top, bottom)
	}
	cs := find(t, tt, "Gold")
	cx, cy := int(cs.X)/2, int(cs.Y)/2
	if count(t, img, cx, cy, 80, 80, pixel.Good) == 0 || count(t, img, cx+80, cy, 80, 80, pixel.Bad) == 0 {
		t.Fatal("candles are not inked by direction")
	}
	keeps(t, img)
	panics(t, "a candle closing above its high", func() {
		view(200, 100, newClock(), func(c *ui.Context) {
			chart.Candlestick(c, "x", []chart.Candle{{Open: 1, High: 2, Low: 0, Close: 3}}, chart.CandlestickOptions{})
		})
	})
}

func TestBoxPlotMedianAndWhiskers(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.BoxPlot(c, "Rolls", []chart.Series{{Name: "D6", Values: []float64{1, 2, 3, 4, 5, 6, 7}}}, chart.BoxPlotOptions{})
	})
	img := tt.Image()
	if count(t, img, 0, 0, 160, 80, pixel.Accent) == 0 || count(t, img, 0, 0, 160, 80, pixel.Edge) < 40 {
		t.Fatal("no box or no whiskers")
	}
	keeps(t, img)
}

func TestLegendTogglesSeries(t *testing.T) {
	hidden, changes := []bool{false, false}, 0
	disabled, inert := false, false
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Inert: inert})
		if chart.Legend(c, "Series", hp, hidden, chart.LegendOptions{Disabled: disabled}).Changed {
			changes++
		}
	}, 300, 60)
	if err := tt.Click("MP"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeySpace)
	if !hidden[0] || !hidden[1] || changes != 2 {
		t.Fatalf("click MP, Left, Space left %v", hidden)
	}
	keeps(t, tt.Image())
	for name, set := range map[string]func(bool){"disabled": func(b bool) { disabled = b }, "inert": func(b bool) { inert = b }} {
		set(true)
		tt.Frame()
		tt.Click("HP")
		tt.Key(0, ui.KeySpace)
		if !hidden[0] || changes != 2 {
			t.Fatalf("a %s legend toggled", name)
		}
		set(false)
	}
}

func TestScatterAndBubblesPlacePoints(t *testing.T) {
	tt := view(400, 400, newClock(), func(c *ui.Context) {
		chart.ScatterPlot(c, "Hits", []chart.Point{{X: 1, Y: 1}, {X: 5, Y: 5}, {X: 9, Y: 9}}, chart.ScatterPlotOptions{Slot: pixel.Magic})
		chart.BubbleChart(c, "Towns", []chart.Point{{X: 0, Y: 5, Size: 1}, {X: 10, Y: 5, Size: 4}}, chart.BubbleChartOptions{Slot: pixel.Info})
	})
	img := tt.Image()
	if n := count(t, img, 0, 0, 160, 80, pixel.Magic); n != 12 {
		t.Fatalf("three points drew %d dots", n)
	}
	// Each point sits up and right of the one before.
	var xs, ys []int
	for y := range 80 {
		for x := range 160 {
			if at(t, img, x, y) == pixel.Magic && (len(xs) == 0 || abs(xs[len(xs)-1]-x) > 2) {
				xs, ys = append(xs, x), append(ys, y)
			}
		}
	}
	if len(xs) != 3 || !(xs[0] > xs[1] && xs[1] > xs[2] && ys[0] < ys[1] && ys[1] < ys[2]) {
		t.Fatalf("points at x %v, y %v", xs, ys)
	}
	small, big := count(t, img, 0, 80, 80, 80, pixel.Info), count(t, img, 80, 80, 80, 80, pixel.Info)
	// Area follows size: 2–4× the dots for 4× the size.
	if small == 0 || big < 2*small || big > 4*small {
		t.Fatalf("bubbles of size 1 and 4 hold %d and %d dots", small, big)
	}
	keeps(t, img)
	panics(t, "a negative bubble", func() {
		view(200, 100, newClock(), func(c *ui.Context) {
			chart.BubbleChart(c, "x", []chart.Point{{Size: -1}}, chart.BubbleChartOptions{})
		})
	})
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func TestPyramidAndWinLoss(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.Pyramid(c, "Party", []string{"KIDS"}, chart.Series{Name: "BOYS", Values: []float64{10}, Slot: pixel.Info},
			chart.Series{Name: "GIRLS", Values: []float64{5}, Slot: pixel.Warn}, chart.PyramidOptions{Width: 100})
		chart.WinLoss(c, "Season", []int{1, -1, 0}, chart.WinLossOptions{})
	})
	img := tt.Image()
	left, right := count(t, img, 0, 20, 50, 4, pixel.Info), count(t, img, 50, 20, 50, 4, pixel.Warn)
	if left == 0 || right == 0 || left < 2*right-8 {
		t.Fatalf("pyramid bars %d and %d for 10 and 5", left, right)
	}
	w := find(t, tt, "Season")
	wx, wy := int(w.X)/2, int(w.Y)/2
	if at(t, img, wx+1, wy+4) != pixel.Good || at(t, img, wx+6, wy+11) != pixel.Bad || at(t, img, wx+11, wy+8) != pixel.Faint {
		t.Fatal("win, loss and draw blocks are misplaced")
	}
	panics(t, "a result of 2", func() {
		view(100, 50, newClock(), func(c *ui.Context) { chart.WinLoss(c, "x", []int{2}, chart.WinLossOptions{}) })
	})
}

func TestAxisMarksTicksOrLabels(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.Axis(c, chart.AxisOptions{})
		chart.Axis(c, chart.AxisOptions{Labels: []string{"LOW", "HIGH"}, Vertical: true, Length: 40})
		chart.Axis(c, chart.AxisOptions{Max: 250})
	})
	if !tt.HasText("0–100") || !tt.HasText("LOW–HIGH") || !tt.HasText("0–250") {
		t.Fatalf("texts %q", tt.Texts())
	}
	keeps(t, tt.Image())
	panics(t, "one mark", func() {
		view(100, 50, newClock(), func(c *ui.Context) { chart.Axis(c, chart.AxisOptions{Labels: []string{"A"}}) })
	})
}
