package project_test

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/project"
)

func TestExtremesFailFastOrHold(t *testing.T) {
	// At 40 dots bar and count fit; at 20, none.
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		project.Checklist(c, "C", []string{"A"}, []bool{false}, project.ChecklistOptions{Width: 40})
	})
	if box, n := find(t, tt, "C"), find(t, tt, "0/1"); n.X+n.W > box.X+box.W {
		t.Errorf("the count ends at %v, past the 40-dot list", n.X+n.W)
	}
	mustPanic(t, "a checklist with no room for a bar", func() {
		view(400, 200, newWorld(), func(c *ui.Context) {
			project.Checklist(c, "C", []string{"A"}, []bool{false}, project.ChecklistOptions{Width: 20})
		})
	})
	mustPanic(t, "a pomodoro of 292 years", func() {
		view(400, 200, newWorld(), func(c *ui.Context) {
			project.Pomodoro(c, "P", project.PomodoroOptions{Work: math.MaxInt64})
		})
	})
	wd := newWorld()
	far := project.Sprint{Name: "S", Start: wd.clk.t, End: wd.clk.t.AddDate(300, 0, 0), Total: 1}
	tt = view(600, 200, wd, func(c *ui.Context) { project.SprintPanel(c, far, project.SprintPanelOptions{}) })
	// 300 years away saturates to 106752 days, never below zero.
	if !tt.HasText("106752 DAYS LEFT") {
		t.Errorf("a sprint 300 years long: %q", tt.Texts())
	}
	tt = view(600, 300, newWorld(), func(c *ui.Context) {
		project.Burndown(c, "Burn", 1e308, 5, []float64{1e308, 1.5e308}, project.BurndownOptions{})
	})
	x, y, w, _ := dots(find(t, tt, "Burn"))
	if count(tt.Image(), x+w-70, y+3, 66, 16, pixel.Bad) == 0 {
		t.Error("1.5e308 left of 1e308 is not BEHIND")
	}
	tt = view(700, 200, newWorld(), func(c *ui.Context) {
		project.WorkloadBar(c, "Load", []project.Load{{Name: "A", Hours: 1.5e308}}, project.WorkloadBarOptions{Capacity: 1.5e308})
	})
	x, y, _, _ = dots(find(t, tt, "A"))
	if count(tt.Image(), x+24+pixel.Measure("A"), y, 120, 16, pixel.Warn) == 0 {
		t.Error("a full load of 1.5e308 hours shows no Warn fill")
	}
	wd = newWorld()
	es := []project.TimeEntry{
		{Task: "A", Start: time.Unix(0, 0).AddDate(-280, 0, 0), Stop: time.Unix(0, 0)},
		{Task: "B", Start: time.Unix(0, 0), Stop: wd.clk.t},
	}
	tt = view(600, 300, wd, func(c *ui.Context) { project.TimeTracker(c, "T", "A", &es, project.TimeTrackerOptions{}) })
	if tt.HasText("-") {
		t.Errorf("a total past int64 went negative: %q", tt.Texts())
	}
}

func TestNarrowAndFutureInputs(t *testing.T) {
	wd := newWorld()
	s := project.Sprint{Name: "S", Start: wd.clk.t, End: wd.clk.t.AddDate(0, 0, 9), Total: 1}
	mustPanic(t, "a sprint panel with no room for a bar", func() {
		view(600, 200, newWorld(), func(c *ui.Context) { project.SprintPanel(c, s, project.SprintPanelOptions{Width: 19}) })
	})
	mustPanic(t, "a roadmap of 121 months", func() {
		view(600, 200, newWorld(), func(c *ui.Context) {
			project.Roadmap(c, "R", wd.clk.t, []project.Lane{{Name: "E"}}, project.RoadmapOptions{Months: 121})
		})
	})
	// A start a minute ahead, as clocks drift, shows zero.
	es := []project.TimeEntry{{Task: "A", Start: wd.clk.t.Add(time.Minute)}}
	tt := view(600, 300, wd, func(c *ui.Context) { project.TimeTracker(c, "T", "A", &es, project.TimeTrackerOptions{}) })
	if tt.HasText("-") || !tt.HasText("00:00") {
		t.Errorf("a start ahead of now reads %q", tt.Texts())
	}
}

func TestStopAheadAndHugeSprints(t *testing.T) {
	wd := newWorld()
	es := []project.TimeEntry{{Task: "A", Start: wd.clk.t.Add(time.Minute)}}
	tt := view(600, 300, wd, func(c *ui.Context) { project.TimeTracker(c, "T", "A", &es, project.TimeTrackerOptions{}) })
	// STOP before the start ends the entry at its start.
	click(t, tt, "STOP")
	tt.Frame()
	if !es[0].Stop.Equal(es[0].Start) {
		t.Fatalf("STOP ahead of the start set %v for %v", es[0].Stop, es[0].Start)
	}
	tt = view(600, 300, newWorld(), func(c *ui.Context) {
		project.Burndown(c, "Burn", 10, math.MaxInt/2, []float64{10}, project.BurndownOptions{})
	})
	x, y, w, h := dots(find(t, tt, "Burn"))
	// The ideal line crosses the right half to zero.
	if count(tt.Image(), x+w/2, y+2, w/2-2, h-12, pixel.Faint) == 0 {
		t.Fatal("the ideal line never reaches the right half")
	}
}
