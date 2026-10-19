package project_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/project"
)

func TestPomodoroWorksRestsAndBanks(t *testing.T) {
	wd := newWorld()
	changed, rang := 0, 0
	tt := view(500, 200, wd, func(c *ui.Context) {
		r := project.Pomodoro(c, "Tomato", project.PomodoroOptions{Work: time.Minute, Rest: 30 * time.Second})
		if r.Changed {
			changed++
		}
		if r.Rang {
			rang++
		}
	})
	click(t, tt, "START")
	later(tt, wd, 20*time.Second)
	click(t, tt, "PAUSE")
	// Paused, the time holds; started again, it runs on.
	later(tt, wd, time.Hour)
	if !tt.HasText("00:40") || changed != 2 {
		t.Fatalf("paused at 20 s: texts %q, changed %d", tt.Texts(), changed)
	}
	x, y, _, _ := dots(find(t, tt, "FOCUS"))
	awake := badge(tt, x, y, 32)
	click(t, tt, "START")
	later(tt, wd, 40*time.Second)
	if rang != 1 || !tt.HasText("00:30") {
		t.Fatalf("work ran out: rang %d, texts %q", rang, tt.Texts())
	}
	// At rest the tomato shuts its eyes.
	if x, y, _, _ = dots(find(t, tt, "REST")); badge(tt, x, y, 32) == awake {
		t.Fatal("the tomato looks the same at rest")
	}
	click(t, tt, "RESET")
	if !tt.HasText("01:00") || changed != 4 {
		t.Fatalf("RESET: texts %q, changed %d", tt.Texts(), changed)
	}
	find(t, tt, "FOCUS")
}

func TestTimeTrackerStartsStopsAndTotals(t *testing.T) {
	wd := newWorld()
	entries := []project.TimeEntry{{Task: "MAP", Start: wd.clk.t.Add(-time.Hour), Stop: wd.clk.t.Add(-30 * time.Minute)}}
	tt := view(600, 300, wd, func(c *ui.Context) {
		project.TimeTracker(c, "Time", "GRIND", &entries, project.TimeTrackerOptions{})
	})
	click(t, tt, "START")
	later(tt, wd, 90*time.Second)
	if len(entries) != 2 || !entries[1].Stop.IsZero() || !tt.HasText("01:30") {
		t.Fatalf("90 s after START: %+v, texts %q", entries, tt.Texts())
	}
	click(t, tt, "STOP")
	later(tt, wd, time.Minute)
	if entries[1].Stop.Sub(entries[1].Start) != 90*time.Second || !tt.HasText("31:30") {
		t.Fatalf("STOP: %+v, texts %q", entries, tt.Texts())
	}
	mustPanic(t, "two running entries", func() {
		es := []project.TimeEntry{{Task: "A", Start: wd.clk.t}, {Task: "B", Start: wd.clk.t}}
		view(600, 300, newWorld(), func(c *ui.Context) { project.TimeTracker(c, "T", "A", &es, project.TimeTrackerOptions{}) })
	})
}

func TestSprintPanelCountsDays(t *testing.T) {
	wd := newWorld()
	s := project.Sprint{Name: "S7", Start: wd.clk.t.AddDate(0, 0, -9), End: wd.clk.t.Add(3*24*time.Hour + time.Hour), Done: 18, Total: 30}
	tt := view(600, 300, wd, func(c *ui.Context) { project.SprintPanel(c, s, project.SprintPanelOptions{}) })
	if !tt.HasText("4 DAYS LEFT") || !tt.HasText("18/30 POINTS") {
		t.Fatalf("texts %q", tt.Texts())
	}
	later(tt, wd, 3*24*time.Hour)
	if !tt.HasText("LAST DAY") {
		t.Fatalf("an hour left: %q", tt.Texts())
	}
	later(tt, wd, 2*time.Hour)
	if !tt.HasText("OVER") {
		t.Fatalf("past the end: %q", tt.Texts())
	}
	mustPanic(t, "more done than the total", func() {
		bad := s
		bad.Done = 31
		view(600, 300, newWorld(), func(c *ui.Context) { project.SprintPanel(c, bad, project.SprintPanelOptions{}) })
	})
}

func TestMilestonesFlagTheNextOne(t *testing.T) {
	wd := newWorld()
	d := func(days int) time.Time { return wd.clk.t.AddDate(0, 0, days) }
	ms := []project.Milestone{{Title: "A", Date: d(-10), Done: true}, {Title: "B", Date: d(10)}, {Title: "C", Date: d(30)}}
	tt := view(600, 200, wd, func(c *ui.Context) { project.Milestones(c, "Road", ms, project.MilestonesOptions{Width: 200}) })
	x, y, _, _ := dots(find(t, tt, "Road"))
	img := tt.Image()
	// Done flies Good, next flies Accent, later ones Faint.
	flags := []pixel.Slot{pixel.Good, pixel.Accent, pixel.Faint}
	for i, at := range []int{8, 100, 192} {
		if count(img, x+at, y+5, 8, 5, flags[i]) == 0 {
			t.Errorf("flag %d at %d is not %v", i, at, flags[i])
		}
	}
	mustPanic(t, "dates out of order", func() {
		bad := []project.Milestone{{Title: "A", Date: d(5)}, {Title: "B", Date: d(1)}}
		view(600, 200, newWorld(), func(c *ui.Context) { project.Milestones(c, "R", bad, project.MilestonesOptions{}) })
	})
}

func TestBurndownJudgesTheLastDay(t *testing.T) {
	for _, tc := range []struct {
		left []float64
		ink  pixel.Slot
	}{{[]float64{10, 9, 9}, pixel.Bad}, {[]float64{10, 7, 4}, pixel.Good}} {
		tt := view(600, 300, newWorld(), func(c *ui.Context) {
			project.Burndown(c, "Burn", 10, 5, tc.left, project.BurndownOptions{})
		})
		x, y, w, _ := dots(find(t, tt, "Burn"))
		// The verdict sits top right: BEHIND Bad, AHEAD Good.
		if count(tt.Image(), x+w-70, y+3, 66, 16, tc.ink) == 0 {
			t.Errorf("%v: no verdict in slot %v", tc.left, tc.ink)
		}
	}
	mustPanic(t, "more days of data than the sprint", func() {
		view(600, 300, newWorld(), func(c *ui.Context) { project.Burndown(c, "B", 10, 1, []float64{1, 2, 3}, project.BurndownOptions{}) })
	})
}

func TestRoadmapMarksTodayAndBars(t *testing.T) {
	wd := newWorld()
	from := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	lanes := []project.Lane{{Name: "ENG", Bars: []project.Bar{{Title: "A", Start: from, End: from.AddDate(0, 2, 0)}}}}
	tt := view(700, 200, wd, func(c *ui.Context) {
		project.Roadmap(c, "Plan", from, lanes, project.RoadmapOptions{Months: 4, MonthWidth: 40})
	})
	x, y, _, h := dots(find(t, tt, "Plan"))
	// Today, Sep 2, falls in September's column, by its days.
	nameW := pixel.Measure("ENG") + 8
	sep := from.AddDate(0, 1, 0)
	today := x + nameW + 40 + int(40*float64(wd.clk.t.Sub(sep))/float64(sep.AddDate(0, 1, 0).Sub(sep)))
	if count(tt.Image(), today, y+20, 1, h-20, pixel.Accent) == 0 {
		t.Fatal("no Accent line at today")
	}
	if count(tt.Image(), x+nameW+2, y+24, 3, 8, pixel.Info) == 0 {
		t.Fatal("the first bar has no Info cap")
	}
	mustPanic(t, "a bar that ends first", func() {
		bad := []project.Lane{{Name: "E", Bars: []project.Bar{{Title: "A", Start: from, End: from}}}}
		view(700, 200, newWorld(), func(c *ui.Context) { project.Roadmap(c, "P", from, bad, project.RoadmapOptions{}) })
	})
}

// badge is the slots of an n-dot square there.
func badge(tt *ui.Tester, x, y, n int) string {
	b := make([]byte, 0, n*n)
	for j := y; j < y+n; j++ {
		for i := x; i < x+n; i++ {
			b = append(b, byte(at(tt.Image(), i, j)))
		}
	}
	return string(b)
}

func TestRoadmapLinesUpWithMonths(t *testing.T) {
	wd := newWorld()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	wd.clk.t = time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	lanes := []project.Lane{{Name: "E", Bars: []project.Bar{{Title: "B", Start: wd.clk.t, End: wd.clk.t.AddDate(0, 0, 14)}}}}
	tt := view(900, 200, wd, func(c *ui.Context) {
		project.Roadmap(c, "Plan", from, lanes, project.RoadmapOptions{Months: 3, MonthWidth: 100})
	})
	x, y, _, _ := dots(find(t, tt, "Plan"))
	// February, today and B start 100 dots past the names.
	feb, img := x+pixel.Measure("E")+8+100, tt.Image()
	if at(img, feb, y+18) != pixel.Accent || at(img, feb-1, y+23) != pixel.Back || at(img, feb, y+23) == pixel.Back {
		t.Fatalf("at February's line: today %v; bar %v after %v", at(img, feb, y+18), at(img, feb, y+23), at(img, feb-1, y+23))
	}
}
