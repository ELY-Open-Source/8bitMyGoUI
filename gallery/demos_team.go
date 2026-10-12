package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/project"
)

var teamEntries = []Entry{
	{429, "LabelPicker", "标签选择", "project", "Labels in slots. Toggle the ones that fit.", demoLabelPicker},
	{430, "TimeTracker", "计时记录", "project", "START, STOP, and where the hours went.", demoTimeTracker},
	{431, "Roadmap", "路线图", "project", "Lanes of work laid over the months. The dashed line is today.", demoRoadmap},
	{432, "SprintPanel", "迭代", "project", "Points done and days left. No pressure.", demoSprintPanel},
	{433, "IssueCard", "问题卡", "project", "An issue at a glance: number, title, labels, owner. Click to open.", demoIssueCard},
	{434, "Checklist", "清单", "project", "Small steps under a bar. Finish them all for a fanfare.", demoChecklist},
	{435, "WorkloadBar", "工作量", "project", "Hours per party member. Gold is busy; red is too much.", demoWorkloadBar},
}

var tags = []project.Label{
	{Name: "BUG", Slot: pixel.Bad}, {Name: "FEATURE", Slot: pixel.Good}, {Name: "DOCS", Slot: pixel.Info},
	{Name: "UI", Slot: pixel.Magic}, {Name: "SPEED", Slot: pixel.Warn}, {Name: "LORE", Slot: pixel.Accent},
}

type teamDemo struct {
	on      []bool
	entries []project.TimeEntry
	steps   []bool
}

func teamState(c *ui.Context) *teamDemo {
	now := pixel.K(c).Now()
	return ui.Local(c.Root(), demoKey("team"), func() teamDemo {
		return teamDemo{
			on: []bool{true, false, false, true, false, false},
			entries: []project.TimeEntry{
				{Task: "GRIND XP", Start: now.Add(-3 * time.Hour), Stop: now.Add(-2*time.Hour - 35*time.Minute)},
				{Task: "MAP THE CAVE", Start: now.Add(-90 * time.Minute), Stop: now.Add(-48 * time.Minute)},
			},
			steps: []bool{true, true, false, false},
		}
	})
}

func demoLabelPicker(c *ui.Context, log func(string)) {
	s := teamState(c)
	if project.LabelPicker(c, "Labels", tags, s.on, project.LabelPickerOptions{}).Changed {
		n := 0
		for _, on := range s.on {
			if on {
				n++
			}
		}
		log(fmt.Sprintf("%d labels on", n))
	}
}

func demoTimeTracker(c *ui.Context, log func(string)) {
	s := teamState(c)
	r := project.TimeTracker(c, "Time", "GRIND XP", &s.entries, project.TimeTrackerOptions{})
	if r.Started {
		log("started")
	}
	if r.Stopped {
		log("stopped")
	}
}

func demoRoadmap(c *ui.Context, log func(string)) {
	now := pixel.K(c).Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0)
	m := func(month, day int) time.Time { return from.AddDate(0, month, day) }
	project.Roadmap(c, "Plan", from, []project.Lane{
		{Name: "ENGINE", Bars: []project.Bar{{Title: "PALETTES", Start: m(0, 0), End: m(1, 10)}, {Title: "SPRITES", Start: m(1, 14), End: m(3, 0)}}},
		{Name: "GAME", Bars: []project.Bar{{Title: "LEVEL 1", Start: m(1, 0), End: m(4, 0)}}},
		{Name: "SOUND", Bars: []project.Bar{{Title: "CHIPTUNE", Start: m(3, 10), End: m(5, 20)}}},
	}, project.RoadmapOptions{})
}

func demoSprintPanel(c *ui.Context, log func(string)) {
	now := pixel.K(c).Now()
	project.SprintPanel(c, project.Sprint{Name: "SPRINT 7: THE DUNGEON", Start: now.AddDate(0, 0, -9), End: now.AddDate(0, 0, 4), Done: 18, Total: 30}, project.SprintPanelOptions{})
}

func demoIssueCard(c *ui.Context, log func(string)) {
	is := project.Issue{ID: "#42", Title: "Save file vanishes after the boss", Labels: tags[:1], Assignee: "Brom", Priority: project.Critical}
	if project.IssueCard(c, is, project.IssueCardOptions{}).Clicked {
		log("open " + is.ID)
	}
}

func demoChecklist(c *ui.Context, log func(string)) {
	s := teamState(c)
	if project.Checklist(c, "Before the quest", []string{"Sharpen the sword", "Pack potions", "Tell mom", "Leave the village"}, s.steps, project.ChecklistOptions{}).Changed {
		log(fmt.Sprintf("steps %v", s.steps))
	}
}

func demoWorkloadBar(c *ui.Context, log func(string)) {
	project.WorkloadBar(c, "Load", []project.Load{{Name: "Aria", Hours: 32}, {Name: "Brom", Hours: 44}, {Name: "Cass", Hours: 18}, {Name: "Dax", Hours: 39}}, project.WorkloadBarOptions{})
}
