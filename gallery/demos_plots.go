package gallery

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/chart"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var plotEntries = []Entry{
	{231, "BulletChart", "子弹图", "chart", "Progress against a goal, over bands of poor, fine and great.", demoBulletChart},
	{232, "Funnel", "漏斗图", "chart", "From SEEN to BOUGHT at the item shop. Most just window-shop.", demoFunnel},
	{233, "Waterfall", "瀑布图", "chart", "Your wallet through a day: green steps up, red steps down.", demoWaterfall},
	{234, "Treemap", "矩形树图", "chart", "The bag by weight, tiled into near squares.", demoTreemap},
	{235, "Candlestick", "K线图", "chart", "Potion prices. Green days rose, red days fell.", demoCandlestick},
	{236, "BoxPlot", "箱线图", "chart", "Rolls by die: the middle half in the box, the ends on whiskers.", demoBoxPlot},
	{237, "StepChart", "阶梯图", "chart", "Hero level by week. It only ever steps.", demoStepChart},
	{238, "ProgressRing", "进度环", "chart", "Clockwise from twelve, a pixel arc at a time.", demoProgressRing},
	{239, "Legend", "图例", "chart", "Click a hero to hide or show their line.", demoLegend},
	{240, "Axis", "坐标轴", "chart", "Ticks and micro labels, across or down.", demoAxis},
	{241, "Crosshair", "十字线", "chart", "Hover the chart; the crosshair reads every line.", demoCrosshair},
	{242, "PolarArea", "极坐标图", "chart", "Equal turns, each reaching out by its value.", demoPolarArea},
	{243, "Pyramid", "金字塔图", "chart", "Heroes and monsters back to back, by level.", demoPyramid},
	{244, "SpectrumBars", "频谱条", "chart", "Bars dance in segments; peak dots fall behind them.", demoSpectrumBars},
	{245, "WinLoss", "胜负图", "chart", "Wins up, losses down, draws a dash. PLAY for one more.", demoWinLoss},
}

func demoBulletChart(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		pixel.Text(c, "Quest XP: 62 of 80", pixel.TextOptions{Slot: pixel.Faint})
		chart.BulletChart(c, "Quest XP", 62, 80, []float64{40, 70, 100}, chart.BulletChartOptions{Width: 200})
	})
}

func demoFunnel(c *ui.Context, log func(string)) {
	chart.Funnel(c, "Item shop", []chart.Slice{{Name: "SEEN", Value: 120}, {Name: "CLICKED", Value: 64}, {Name: "IN CART", Value: 30},
		{Name: "BOUGHT", Value: 12}}, chart.FunnelOptions{Width: 200})
}

func demoWaterfall(c *ui.Context, log func(string)) {
	chart.Waterfall(c, "Wallet", []string{"QUEST", "INN", "SHOP", "LOOT", "TAX"}, []float64{40, -8, -25, 18, -5}, chart.WaterfallOptions{Width: 200, Height: 90})
}

func demoTreemap(c *ui.Context, log func(string)) {
	chart.Treemap(c, "Bag weight", []chart.Slice{{Name: "ARMOR", Value: 10}, {Name: "SWORD", Value: 8}, {Name: "SHIELD", Value: 6},
		{Name: "HERBS", Value: 2}, {Name: "KEYS", Value: 1}, {Name: "MAP", Value: 1}}, chart.TreemapOptions{Width: 200, Height: 90})
}

func demoCandlestick(c *ui.Context, log func(string)) {
	chart.Candlestick(c, "Potion price", []chart.Candle{
		{Label: "MON", Open: 10, High: 14, Low: 9, Close: 13}, {Label: "TUE", Open: 13, High: 15, Low: 8, Close: 9},
		{Label: "WED", Open: 9, High: 12, Low: 7, Close: 11}, {Label: "THU", Open: 11, High: 18, Low: 11, Close: 17},
		{Label: "FRI", Open: 17, High: 17, Low: 12, Close: 13},
	}, chart.CandlestickOptions{Width: 200, Height: 90})
}

func demoBoxPlot(c *ui.Context, log func(string)) {
	chart.BoxPlot(c, "Rolls", []chart.Series{
		{Name: "D6", Values: []float64{1, 2, 3, 3, 4, 5, 6}},
		{Name: "D8", Values: []float64{1, 3, 4, 5, 5, 7, 8}},
		{Name: "D12", Values: []float64{2, 4, 6, 7, 9, 11, 12}},
	}, chart.BoxPlotOptions{Width: 200, Height: 90})
}

func demoStepChart(c *ui.Context, log func(string)) {
	chart.StepChart(c, "Level", []string{"W1", "W2", "W3", "W4", "W5", "W6"}, []float64{1, 3, 3, 4, 7, 8}, chart.StepChartOptions{Width: 200, Height: 90, Slot: pixel.Magic})
}

func demoProgressRing(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := chartState(c)
	// The pretend download takes four seconds.
	done := k.Since(c, s.since, 10, 40)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		chart.ProgressRing(c, "Download", float64(done)/40, chart.ProgressRingOptions{Size: 48})
		if input.Button(c, "DOWNLOAD", input.ButtonOptions{Disabled: done < 40}).Clicked {
			s.since = k.Now()
			log("downloading more RAM")
		}
	})
}

func demoLegend(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := chartState(c)
	var shown []chart.Series
	for i, h := range heroes {
		if !s.hidden[i] {
			shown = append(shown, h)
		}
	}
	ui.Column(c).Gap(k.D(4)).Children(func() {
		if chart.Legend(c, "Heroes", heroes, s.hidden, chart.LegendOptions{}).Changed {
			log(fmt.Sprintf("hidden %v", s.hidden))
		}
		chart.LineChart(c, "Hero HP", week, shown, chart.LineChartOptions{Width: 200, Height: 90})
	})
}

func demoAxis(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(12)).AlignItems(ui.Center).Children(func() {
		chart.Axis(c, chart.AxisOptions{Labels: []string{"LOW", "MID", "HIGH"}, Vertical: true, Length: 60})
		chart.Axis(c, chart.AxisOptions{Min: 0, Max: 250, Length: 160})
	})
}

func demoCrosshair(c *ui.Context, log func(string)) {
	chart.Crosshair(c, "Hero HP", week, heroes, chart.CrosshairOptions{Width: 200, Height: 90})
}

func demoPolarArea(c *ui.Context, log func(string)) {
	chart.PolarArea(c, "Rain", []chart.Slice{{Name: "SPRING", Value: 6, Slot: pixel.Good}, {Name: "SUMMER", Value: 3, Slot: pixel.Warn},
		{Name: "AUTUMN", Value: 8, Slot: pixel.Bad}, {Name: "WINTER", Value: 4, Slot: pixel.Info}}, chart.PolarAreaOptions{Size: 80})
}

func demoPyramid(c *ui.Context, log func(string)) {
	chart.Pyramid(c, "By level", []string{"LV 1-9", "LV 10-19", "LV 20-29", "LV 30+"},
		chart.Series{Name: "HEROES", Values: []float64{40, 26, 12, 4}, Slot: pixel.Info},
		chart.Series{Name: "MONSTERS", Values: []float64{60, 30, 18, 9}, Slot: pixel.Bad}, chart.PyramidOptions{Width: 220})
}

func demoSpectrumBars(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := chartState(c)
	n := 0
	if s.play && !k.Reduce() {
		n = k.Tick(c, 8)
	}
	levels := make([]float64, 12)
	for i := range levels {
		// A made-up song: a slow swell times a quick wobble.
		t, f := float64(n), float64(i)
		levels[i] = 0.5 + 0.45*math.Sin(t*0.7+f*1.3)*math.Sin(t*0.13+f*0.4)
	}
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		chart.SpectrumBars(c, "Music", levels, chart.SpectrumBarsOptions{Segments: 12})
		if input.Toggle(c, "PLAY", &s.play, input.ToggleOptions{Size: input.Small}).Changed {
			log(fmt.Sprintf("play %v", s.play))
		}
	})
}

func demoWinLoss(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := chartState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		chart.WinLoss(c, "Arena", s.matches, chart.WinLossOptions{})
		if input.Button(c, "PLAY", input.ButtonOptions{Size: input.Small}).Clicked {
			r := s.dice.IntN(3) - 1
			// The strip keeps the last twenty matches.
			s.matches = append(s.matches[max(0, len(s.matches)-19):], r)
			log([]string{"lost", "draw", "won"}[r+1])
		}
	})
}
