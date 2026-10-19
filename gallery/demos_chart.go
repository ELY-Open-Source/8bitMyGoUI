package gallery

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/chart"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var chartEntries = []Entry{
	{216, "BarChart", "柱状图", "chart", "Gold earned per day. Point at a bar to read it.", demoBarChart},
	{217, "StackedBar", "堆叠柱图", "chart", "Gold under gems, stacked day by day.", demoStackedBar},
	{218, "HorizontalBar", "条形图", "chart", "Slimes slain, by hero. Brom is not sharing.", demoHorizontalBar},
	{219, "LineChart", "折线图", "chart", "HP and MP down a dungeon, joined dot by dot.", demoLineChart},
	{220, "AreaChart", "面积图", "chart", "Gold over a week; the ground under the line dithers.", demoAreaChart},
	{221, "Sparkline", "迷你图", "chart", "A chart one text line tall. It fits in a sentence.", demoSparkline},
	{222, "PieChart", "饼图", "chart", "Where the loot went, sector by sector, dot by dot.", demoPieChart},
	{223, "DonutChart", "环形图", "chart", "A pie with a hole; the total waits inside.", demoDonutChart},
	{224, "RadarChart", "雷达图", "chart", "Six stats on a hexagon web. LEVEL UP and watch it grow.", demoRadarChart},
	{225, "Histogram", "直方图", "chart", "Damage rolls in bins. ROLL throws ten more.", demoHistogram},
	{226, "ScatterPlot", "散点图", "chart", "Every hit by distance and damage, one 2×2 dot each.", demoScatterPlot},
	{227, "BubbleChart", "气泡图", "chart", "Towns by trade and danger, sized by how many live there.", demoBubbleChart},
	{228, "Heatmap", "密度图", "chart", "Raids by day and hour. Hot cells dither red.", demoHeatmap},
	{229, "CalendarHeatmap", "日历密度图", "chart", "A year of play, a square a day. A commit graph for slimes.", demoCalendarHeatmap},
	{230, "Gauge", "仪表", "chart", "Push the throttle; the needle follows.", demoGauge},
}

type chartDemo struct {
	stats    []float64
	rolls    []float64
	dice     *rand.Rand
	throttle float64
	hidden   []bool
	since    time.Time
	play     bool
	matches  []int
}

func chartState(c *ui.Context) *chartDemo {
	return ui.Local(c.Root(), demoKey("chart"), func() chartDemo {
		s := chartDemo{stats: []float64{4, 3, 2, 4, 3, 2}, dice: rand.New(rand.NewPCG(1, 2)), throttle: 40,
			hidden: []bool{false, false, false}, play: true, matches: []int{1, 1, -1, 1, 0, -1, 1, 1}}
		for range 30 {
			s.rolls = append(s.rolls, roll(s.dice))
		}
		return s
	})
}

// roll throws two dice.
func roll(r *rand.Rand) float64 { return float64(r.IntN(6) + r.IntN(6) + 2) }

var week = []string{"MON", "TUE", "WED", "THU", "FRI"}

// heroes are three HP series with fixed colors.
var heroes = []chart.Series{
	{Name: "Aria", Values: []float64{30, 26, 28, 18, 24}, Slot: pixel.Good},
	{Name: "Brom", Values: []float64{40, 34, 22, 26, 12}, Slot: pixel.Warn},
	{Name: "Cid", Values: []float64{18, 20, 16, 22, 20}, Slot: pixel.Info},
}

func demoBarChart(c *ui.Context, log func(string)) {
	chart.BarChart(c, "Gold per day", week, []float64{120, 340, 90, 260, 410}, chart.BarChartOptions{Width: 200, Height: 90})
}

func demoStackedBar(c *ui.Context, log func(string)) {
	chart.StackedBar(c, "Loot per day", week, []chart.Series{
		{Name: "Gold", Values: []float64{12, 30, 8, 22, 18}, Slot: pixel.Warn},
		{Name: "Gems", Values: []float64{4, 6, 12, 2, 9}, Slot: pixel.Magic},
	}, chart.StackedBarOptions{Width: 200, Height: 90})
}

func demoHorizontalBar(c *ui.Context, log func(string)) {
	chart.HorizontalBar(c, "Slimes slain", []string{"Aria", "Brom", "Cid", "Dax"}, []float64{48, 112, 37, 64}, chart.HorizontalBarOptions{Width: 200})
}

func demoLineChart(c *ui.Context, log func(string)) {
	floors := []string{"B1", "B2", "B3", "B4", "B5", "B6"}
	chart.LineChart(c, "Dungeon run", floors, []chart.Series{
		{Name: "HP", Values: []float64{50, 44, 47, 30, 36, 12}, Slot: pixel.Good},
		{Name: "MP", Values: []float64{20, 16, 8, 12, 4, 6}, Slot: pixel.Info},
	}, chart.LineChartOptions{Width: 200, Height: 90})
}

func demoAreaChart(c *ui.Context, log func(string)) {
	chart.AreaChart(c, "Gold", week, []float64{120, 180, 150, 260, 310}, chart.AreaChartOptions{Width: 200, Height: 90, Slot: pixel.Warn})
}

func demoSparkline(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(2)).Children(func() {
		for _, r := range []struct {
			name, now string
			values    []float64
			slot      pixel.Slot
		}{
			{"Ping", "42 ms", []float64{40, 44, 39, 61, 45, 42, 38, 42}, pixel.Info},
			{"FPS", "60", []float64{60, 60, 58, 41, 59, 60, 60, 60}, pixel.Good},
			{"Gold", "1.2K", []float64{300, 420, 380, 650, 700, 980, 1100, 1200}, pixel.Warn},
		} {
			ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Width(k.D(36)).Children(func() { pixel.Text(c, r.name, pixel.TextOptions{Slot: pixel.Faint}) })
				chart.Sparkline(c, r.name, r.values, chart.SparklineOptions{Slot: r.slot})
				pixel.Text(c, r.now, pixel.TextOptions{})
			})
		}
	})
}

func demoPieChart(c *ui.Context, log func(string)) {
	chart.PieChart(c, "Loot", []chart.Slice{{Name: "Gold", Value: 52, Slot: pixel.Warn}, {Name: "Gems", Value: 23, Slot: pixel.Magic},
		{Name: "Potions", Value: 15, Slot: pixel.Bad}, {Name: "Junk", Value: 10, Slot: pixel.Faint}}, chart.PieChartOptions{Size: 80})
}

func demoDonutChart(c *ui.Context, log func(string)) {
	chart.DonutChart(c, "Party gold", []chart.Slice{{Name: "Aria", Value: 120}, {Name: "Brom", Value: 340}, {Name: "Cid", Value: 90}}, chart.DonutChartOptions{Size: 80})
}

func demoRadarChart(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := chartState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		chart.RadarChart(c, "Hero", []string{"STR", "DEX", "INT", "VIT", "AGI", "LUK"}, s.stats, chart.RadarChartOptions{Max: 10, Size: 96, Slot: pixel.Accent})
		maxed := slices.Min(s.stats) >= 10
		if input.Button(c, "LEVEL UP", input.ButtonOptions{Disabled: maxed}).Clicked {
			for i := range s.stats {
				s.stats[i] = min(10, s.stats[i]+float64(1+i%2))
			}
			log(fmt.Sprintf("STR %g, DEX %g", s.stats[0], s.stats[1]))
		}
	})
}

func demoHistogram(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := chartState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		chart.Histogram(c, "2d6 rolls", s.rolls, chart.HistogramOptions{Bins: 11, Width: 200, Height: 90})
		if input.Button(c, "ROLL", input.ButtonOptions{}).Clicked {
			for range 10 {
				s.rolls = append(s.rolls, roll(s.dice))
			}
			log(fmt.Sprintf("%d rolls", len(s.rolls)))
		}
	})
}

func demoScatterPlot(c *ui.Context, log func(string)) {
	hits := []chart.Point{{X: 1, Y: 12}, {X: 2, Y: 10}, {X: 2, Y: 14}, {X: 3, Y: 9}, {X: 4, Y: 8}, {X: 5, Y: 8}, {X: 5, Y: 3},
		{X: 6, Y: 6}, {X: 7, Y: 5}, {X: 8, Y: 4}, {X: 8, Y: 9}, {X: 9, Y: 2}}
	chart.ScatterPlot(c, "Hits", hits, chart.ScatterPlotOptions{Width: 200, Height: 90, Slot: pixel.Bad})
}

func demoBubbleChart(c *ui.Context, log func(string)) {
	towns := []chart.Point{{X: 2, Y: 8, Size: 120}, {X: 5, Y: 3, Size: 900}, {X: 7, Y: 6, Size: 300}, {X: 9, Y: 9, Size: 40}}
	chart.BubbleChart(c, "Towns", towns, chart.BubbleChartOptions{Width: 200, Height: 90, Slot: pixel.Info})
}

func demoHeatmap(c *ui.Context, log func(string)) {
	chart.Heatmap(c, "Raids", []string{"MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"}, []string{"AM", "NOON", "PM", "NIGHT"}, [][]float64{
		{1, 2, 3, 6}, {0, 2, 4, 7}, {1, 1, 3, 5}, {0, 3, 5, 8}, {2, 4, 7, 12}, {5, 8, 10, 14}, {6, 7, 9, 4},
	}, chart.HeatmapOptions{Cell: 12})
}

func demoCalendarHeatmap(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	now := k.Now()
	// Tomorrow's square lights up at midnight.
	k.At(c, time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()))
	// Days come as UTC dates; today is the local date.
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	chart.CalendarHeatmap(c, "Play days", now.Year(), func(d time.Time) int {
		if today.Before(d) {
			return 0
		}
		// Weekends run hotter; the rest is a made-up habit.
		n := (d.YearDay()*d.YearDay()*7 + d.YearDay()*3) % 10
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			n += 4
		}
		return min(4, n/3)
	}, chart.CalendarHeatmapOptions{})
}

func demoGauge(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := chartState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		chart.Gauge(c, "Speed", s.throttle, chart.GaugeOptions{Size: 80})
		if input.Slider(c, "Throttle", &s.throttle, input.SliderOptions{Max: 100}).Changed {
			log(fmt.Sprintf("speed %.0f", s.throttle))
		}
	})
}
