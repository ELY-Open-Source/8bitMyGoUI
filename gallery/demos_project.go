package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/project"
)

var projectEntries = []Entry{
	{421, "QuestBoard", "任务板", "project", "Quests in three columns. Drag a card across, or hold Alt with the arrows.", demoQuestBoard},
	{422, "TaskItem", "任务项", "project", "A box to tick, a title, how urgent, and who.", demoTaskItem},
	{423, "TaskList", "任务列表", "project", "Tasks with OPEN, DONE and URGENT chips to narrow them down.", demoTaskList},
	{424, "Milestones", "里程碑", "project", "Flags along the road. The little hero stands at today.", demoMilestones},
	{425, "Burndown", "冲刺进度", "project", "The straight line is the plan. The other one is what happened.", demoBurndown},
	{426, "Pomodoro", "番茄钟", "project", "Twenty-five minutes of work, five of rest. The tomato naps on the break.", demoPomodoro},
	{427, "Priority", "优先级", "project", "LOW, MEDIUM, HIGH and, for the worst, CRITICAL: it is on fire.", demoPriority},
	{428, "AssigneePicker", "指派", "project", "The party lines up. Pick who takes the quest, or nobody.", demoAssigneePicker},
}

type projectDemo struct {
	board             [3][]string
	task              project.Task
	tasks             []project.Task
	assignee, changes int
}

func projectState(c *ui.Context) *projectDemo {
	return ui.Local(c.Root(), demoKey("project"), func() projectDemo {
		return projectDemo{
			board: [3][]string{{"SLAY DRAGON", "FIND THE KEY", "SAVE VILLAGE"}, {"GRIND XP"}, {"TALK TO OLD MAN"}},
			task:  project.Task{Title: "Fix the save bug", Priority: project.High, Assignee: "Brom"},
			tasks: []project.Task{
				{Title: "Feed the slime", Priority: project.Low, Assignee: "Aria"},
				{Title: "Fix the save bug", Priority: project.High, Assignee: "Brom", Done: true},
				{Title: "Ship it", Priority: project.Critical},
				{Title: "Draw the map", Priority: project.Medium, Assignee: "Cass"},
			},
			assignee: -1,
		}
	})
}

func demoQuestBoard(c *ui.Context, log func(string)) {
	s := projectState(c)
	if project.QuestBoard(c, "Quests", &s.board, project.QuestBoardOptions{}).Moved {
		log(fmt.Sprintf("moved: %d to do, %d doing, %d done", len(s.board[0]), len(s.board[1]), len(s.board[2])))
	}
}

func demoTaskItem(c *ui.Context, log func(string)) {
	s := projectState(c)
	if project.TaskItem(c, &s.task, project.TaskItemOptions{}).Changed {
		log(fmt.Sprintf("done: %v", s.task.Done))
	}
}

func demoTaskList(c *ui.Context, log func(string)) {
	s := projectState(c)
	if project.TaskList(c, "Tasks", s.tasks, project.TaskListOptions{}).Changed {
		s.changes++
		log(fmt.Sprintf("change %d", s.changes))
	}
}

func demoMilestones(c *ui.Context, log func(string)) {
	now := pixel.K(c).Now()
	d := func(days int) time.Time { return now.AddDate(0, 0, days) }
	project.Milestones(c, "Road", []project.Milestone{
		{Title: "ALPHA", Date: d(-40), Done: true}, {Title: "BETA", Date: d(-12), Done: true},
		{Title: "DEMO", Date: d(9)}, {Title: "LAUNCH", Date: d(30)}, {Title: "DLC", Date: d(60)},
	}, project.MilestonesOptions{})
}

func demoBurndown(c *ui.Context, log func(string)) {
	project.Burndown(c, "Sprint 7", 40, 10, []float64{40, 37, 35, 30, 28, 22}, project.BurndownOptions{})
}

func demoPomodoro(c *ui.Context, log func(string)) {
	r := project.Pomodoro(c, "Tomato", project.PomodoroOptions{})
	if r.Changed {
		log("button")
	}
	if r.Rang {
		log("ring")
	}
}

func demoPriority(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		for _, l := range []project.Level{project.Low, project.Medium, project.High, project.Critical} {
			ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
				project.Priority(c, l, project.PriorityOptions{Scale: 2})
				pixel.Text(c, k.Msg([]string{"low", "medium", "high", "critical"}[l]), pixel.TextOptions{})
			})
		}
	})
}

func demoAssigneePicker(c *ui.Context, log func(string)) {
	s := projectState(c)
	if project.AssigneePicker(c, "Party", party[:4], &s.assignee, project.AssigneePickerOptions{}).Changed {
		who := "nobody"
		if s.assignee >= 0 {
			who = party[s.assignee]
		}
		log("assigned to " + who)
	}
}
