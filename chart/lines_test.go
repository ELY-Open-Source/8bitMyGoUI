package chart_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/chart"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var hp = []chart.Series{{Name: "HP", Values: []float64{10, 30, 20, 40}, Slot: pixel.Good}, {Name: "MP", Values: []float64{5, 5, 15, 10}, Slot: pixel.Info}}

func TestLineChartJoinsAndMarksPoints(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.LineChart(c, "Stats", days, hp, chart.LineChartOptions{})
	})
	img := tt.Image()
	good, info := count(t, img, 0, 0, 160, 80, pixel.Good), count(t, img, 0, 0, 160, 80, pixel.Info)
	if good < 100 || info < 100 {
		t.Fatalf("line dots: HP %d, MP %d", good, info)
	}
	// One-dot lines leave whole 3×3 blocks only at the markers.
	if hb, mb := blocks(t, img, pixel.Good), blocks(t, img, pixel.Info); hb != 4 || mb != 4 {
		t.Fatalf("markers: HP %d, MP %d", hb, mb)
	}
	keeps(t, img)
	panics(t, "a series too short", func() {
		view(100, 100, newClock(), func(c *ui.Context) {
			chart.LineChart(c, "x", days, []chart.Series{{Name: "a", Values: []float64{1}}}, chart.LineChartOptions{})
		})
	})
}

func TestAreaChartDithersBelowTheLine(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.AreaChart(c, "Gold", days, []float64{20, 20, 20, 20}, chart.AreaChartOptions{Slot: pixel.Warn})
	})
	img := tt.Image()
	// Under a flat line, the area mixes Warn and ground.
	lit, ground := count(t, img, 60, 50, 40, 10, pixel.Warn), count(t, img, 60, 50, 40, 10, pixel.Back)
	if lit == 0 || ground == 0 {
		t.Fatalf("area under the line: %d Warn, %d ground", lit, ground)
	}
	keeps(t, img)
}

func TestAreaChartFillsBelowZeroToo(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.AreaChart(c, "Debt", days, []float64{-20, -20, -20, -20}, chart.AreaChartOptions{Slot: pixel.Bad})
	})
	img := tt.Image()
	// Zero tops the plot; the area hangs below it.
	lit, ground := count(t, img, 60, 10, 40, 10, pixel.Bad), count(t, img, 60, 10, 40, 10, pixel.Back)
	if lit == 0 || ground == 0 {
		t.Fatalf("area above the line: %d Bad, %d ground", lit, ground)
	}
}

func TestStepChartHoldsThenSteps(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.StepChart(c, "Level", days, []float64{1, 1, 3, 3}, chart.StepChartOptions{Slot: pixel.Magic})
	})
	img := tt.Image()
	// A flat run of a whole share holds one row.
	rows := map[int]bool{}
	for x := 20; x < 40; x++ {
		for y := range 80 {
			if at(t, img, x, y) == pixel.Magic {
				rows[y] = true
			}
		}
	}
	if len(rows) != 1 {
		t.Fatalf("the first hold spans rows %v", rows)
	}
}

func TestSparklineMarksItsEnd(t *testing.T) {
	tt := view(200, 60, newClock(), func(c *ui.Context) {
		chart.Sparkline(c, "Ping", []float64{3, 9, 4, 12}, chart.SparklineOptions{Slot: pixel.Info})
	})
	img := tt.Image()
	if count(t, img, 0, 0, 48, 16, pixel.Info) < 30 || count(t, img, 44, 0, 4, 16, pixel.Accent) == 0 {
		t.Fatal("no line, or no mark at its end")
	}
	panics(t, "one value", func() {
		view(100, 60, newClock(), func(c *ui.Context) { chart.Sparkline(c, "x", []float64{1}, chart.SparklineOptions{}) })
	})
}

func TestCrosshairCrossesOnHover(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		chart.Crosshair(c, "Stats", days, hp, chart.CrosshairOptions{})
	})
	if n := count(t, tt.Image(), 0, 0, 160, 80, pixel.Lit); n != 0 {
		t.Fatal("a readout shows before hover")
	}
	tt.Move(2*90, 2*30)
	img := tt.Image()
	if count(t, img, 0, 0, 160, 80, pixel.Lit) == 0 || count(t, img, 0, 30, 160, 1, pixel.Ink) < 20 {
		t.Fatal("hover drew no crosshair or readout")
	}
	keeps(t, img)
}

func TestSpectrumPeaksFallAndRest(t *testing.T) {
	clk := newClock()
	levels := []float64{1, 0.5}
	tt := view(200, 100, clk, func(c *ui.Context) {
		chart.SpectrumBars(c, "EQ", levels, chart.SpectrumBarsOptions{})
	})
	img := tt.Image()
	if lit := count(t, img, 0, 0, 5, 30, pixel.Shade); lit != 0 {
		t.Fatalf("a full bar shows %d dark dots", lit)
	}
	// Five idle seconds do not count toward the fall.
	clk.t = clk.t.Add(5 * time.Second)
	levels[0] = 0
	tt.Frame()
	if top, _ := column(t, tt.Image(), 2, 30, pixel.Ink); top != 0 {
		t.Fatalf("at the drop the peak sits at row %d", top)
	}
	clk.t = clk.t.Add(300 * time.Millisecond)
	tt.Frame()
	// After 0.3 s the peak has fallen three blocks.
	top, _ := column(t, tt.Image(), 2, 30, pixel.Ink)
	if top != 9 {
		t.Fatalf("the peak sits at row %d", top)
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	if top, _ := column(t, tt.Image(), 2, 30, pixel.Ink); top != -1 {
		t.Fatalf("reduced, the peak still floats at %d", top)
	}
}

// blocks counts 3×3 squares of s on a 160×80 chart.
func blocks(t *testing.T, img *image.RGBA, s pixel.Slot) int {
	t.Helper()
	n := 0
	for y := range 78 {
		for x := range 158 {
			if count(t, img, x, y, 3, 3, s) == 9 {
				n++
			}
		}
	}
	return n
}
