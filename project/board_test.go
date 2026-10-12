package project_test

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/project"
)

func TestQuestBoardMovesCardsByKeysAndDrags(t *testing.T) {
	for _, off := range []bool{false, true} {
		cols, moved := [3][]string{{"A", "B"}, {"C"}, {}}, 0
		wd := newWorld()
		tt := view(700, 400, wd, func(c *ui.Context) {
			if project.QuestBoard(c, "Board", &cols, project.QuestBoardOptions{Disabled: off}).Moved {
				moved++
			}
		})
		tt.Key(0, ui.KeyTab)
		// Alt carries A to DOING; Down, Alt takes C on.
		tt.Key(ui.Alt, ui.KeyRight)
		tt.Key(0, ui.KeyDown)
		tt.Key(ui.Alt, ui.KeyRight)
		want := [3][]string{{"B"}, {"A"}, {"C"}}
		if off {
			want = [3][]string{{"A", "B"}, {"C"}, {}}
		}
		if !slices.Equal(cols[0], want[0]) || !slices.Equal(cols[1], want[1]) || !slices.Equal(cols[2], want[2]) || off != (moved == 0) {
			t.Fatalf("off %v: keys left %q, moved %d", off, cols, moved)
		}
		if off {
			continue
		}
		// The last card can go no lower.
		tt.Key(ui.Alt, ui.KeyDown)
		if moved != 2 {
			t.Fatalf("Alt+Down on the last card moved it: %d", moved)
		}
		// Columns are 100 dots apart; cards 20 down, 38 tall.
		r := find(t, tt, "Board")
		drag := func(x0, y0, x1, y1 float32) {
			tt.Press(r.X+2*x0, r.Y+2*y0)
			tt.Move(r.X+2*x1, r.Y+2*y1)
			tt.Release(r.X+2*x1, r.Y+2*y1)
		}
		drag(150, 60, 40, 30)
		drag(40, 30, 240, 68)
		if len(cols[0]) != 0 || !slices.Equal(cols[1], []string{"A"}) || !slices.Equal(cols[2], []string{"C", "B"}) {
			t.Fatalf("a press on nothing, then B dropped under C, left %q", cols)
		}
		// Left walks to A; Alt+Left takes it home.
		tt.Key(0, ui.KeyLeft)
		tt.Key(ui.Alt, ui.KeyLeft)
		if !slices.Equal(cols[0], []string{"A"}) || len(cols[1]) != 0 {
			t.Fatalf("Left, Alt+Left left %q", cols)
		}
	}
	mustPanic(t, "a card with no title", func() {
		cols := [3][]string{{""}, {}, {}}
		view(700, 400, newWorld(), func(c *ui.Context) { project.QuestBoard(c, "B", &cols, project.QuestBoardOptions{}) })
	})
}

func TestTaskItemTicksAndShowsWho(t *testing.T) {
	task := project.Task{Title: "Slay", Priority: project.High, Assignee: "Brom"}
	changed := false
	tt := view(500, 200, newWorld(), func(c *ui.Context) {
		changed = project.TaskItem(c, &task, project.TaskItemOptions{}).Changed || changed
	})
	for _, s := range []string{"HIGH", "Brom"} {
		find(t, tt, s)
	}
	click(t, tt, "Slay")
	if !task.Done || !changed {
		t.Fatal("clicking the task did not tick it")
	}
	task.Assignee = ""
	tt.Frame()
	find(t, tt, "NOBODY")
}

func TestTaskListFiltersAndKeepsRows(t *testing.T) {
	tasks := []project.Task{
		{Title: "Low open", Priority: project.Low},
		{Title: "High done", Priority: project.High, Done: true},
		{Title: "Critical open", Priority: project.Critical},
	}
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		project.TaskList(c, "Tasks", tasks, project.TaskListOptions{})
	})
	shown := func() (out []string) {
		for _, tk := range tasks {
			if _, ok := tt.Find(tk.Title); ok {
				out = append(out, tk.Title)
			}
		}
		return out
	}
	click(t, tt, "OPEN")
	if got := shown(); !slices.Equal(got, []string{"Low open", "Critical open"}) {
		t.Fatalf("OPEN shows %q", got)
	}
	click(t, tt, "URGENT")
	if got := shown(); !slices.Equal(got, []string{"Critical open"}) {
		t.Fatalf("OPEN and URGENT show %q", got)
	}
	// The row that stayed still ticks its own task.
	click(t, tt, "Critical open")
	if !tasks[2].Done || tasks[0].Done {
		t.Fatalf("the tick landed on %+v", tasks)
	}
	click(t, tt, "OPEN")
	click(t, tt, "URGENT")
	click(t, tt, "DONE")
	if got := shown(); !slices.Equal(got, []string{"High done", "Critical open"}) {
		t.Fatalf("DONE shows %q", got)
	}
	click(t, tt, "DONE")
	click(t, tt, "URGENT")
	tasks[2].Done = false
	tt.Frame()
	if got := shown(); !slices.Equal(got, []string{"High done", "Critical open"}) {
		t.Fatalf("URGENT shows %q", got)
	}
}

func TestChecklistCountsAndCheers(t *testing.T) {
	done := []bool{true, false}
	tt := view(500, 300, newWorld(), func(c *ui.Context) {
		project.Checklist(c, "Prep", []string{"Pack", "Go"}, done, project.ChecklistOptions{})
	})
	if !tt.HasText("1/2") || tt.HasText("QUEST COMPLETE!") {
		t.Fatalf("texts %q", tt.Texts())
	}
	click(t, tt, "Go")
	if !done[1] || !tt.HasText("2/2") || !tt.HasText("QUEST COMPLETE!") {
		t.Fatalf("after Go: %v, texts %q", done, tt.Texts())
	}
	mustPanic(t, "two items and one mark", func() {
		view(500, 300, newWorld(), func(c *ui.Context) {
			project.Checklist(c, "P", []string{"A", "B"}, []bool{false}, project.ChecklistOptions{})
		})
	})
}

func TestAssigneePickerSeatsAndEmpties(t *testing.T) {
	who := -1
	tt := view(600, 200, newWorld(), func(c *ui.Context) {
		project.AssigneePicker(c, "Party", []string{"Aria", "Brom"}, &who, project.AssigneePickerOptions{})
	})
	click(t, tt, "Brom")
	if who != 1 {
		t.Fatalf("Brom chose %d", who)
	}
	// The chosen seat wears Accent at its top edge.
	x, y, _, _ := dots(find(t, tt, "Brom"))
	if at(tt.Image(), x+20, y) != pixel.Accent {
		t.Fatal("Brom's seat is not marked")
	}
	click(t, tt, "NOBODY")
	if who != -1 {
		t.Fatalf("NOBODY chose %d", who)
	}
	mustPanic(t, "a choice past the party", func() {
		w := 2
		view(600, 200, newWorld(), func(c *ui.Context) {
			project.AssigneePicker(c, "P", []string{"A", "B"}, &w, project.AssigneePickerOptions{})
		})
	})
}

func TestLabelPickerTogglesSlots(t *testing.T) {
	labels := []project.Label{{Name: "BUG", Slot: pixel.Bad}, {Name: "UI", Slot: pixel.Magic}, {Name: "DOCS", Slot: pixel.Info}, {Name: "LORE"}}
	on := make([]bool, 4)
	tt := view(600, 200, newWorld(), func(c *ui.Context) {
		project.LabelPicker(c, "Labels", labels, on, project.LabelPickerOptions{})
	})
	tt.Key(0, ui.KeyTab)
	// Three slots a row: Down from BUG is LORE.
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	r := find(t, tt, "Labels")
	tt.ClickAt(r.X+2*(76+38), r.Y+2*11)
	if !slices.Equal(on, []bool{false, true, false, true}) {
		t.Fatalf("Down, Enter, then a click on UI gave %v", on)
	}
	mustPanic(t, "a nameless label", func() {
		view(600, 200, newWorld(), func(c *ui.Context) {
			project.LabelPicker(c, "L", []project.Label{{}}, []bool{false}, project.LabelPickerOptions{})
		})
	})
}

func TestPickersIgnoreKeysAsTheyTurnOff(t *testing.T) {
	for _, inert := range []bool{false, true} {
		who, off := -1, false
		wd := newWorld()
		tt := view(600, 200, wd, func(c *ui.Context) {
			project.AssigneePicker(c, "Party", []string{"Aria"}, &who, project.AssigneePickerOptions{Disabled: off && !inert})
		})
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyRight)
		// Space, sent while on, lands as the picker turns off.
		off, wd.inert = true, inert
		tt.Key(0, ui.KeySpace)
		if who != -1 {
			t.Errorf("inert %v: Space as the picker turned off chose %d", inert, who)
		}
	}
}

func TestLabelPickerSurvivesFewerLabels(t *testing.T) {
	labels := []project.Label{{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"}}
	on := make([]bool, 4)
	tt := view(600, 200, newWorld(), func(c *ui.Context) {
		project.LabelPicker(c, "Labels", labels, on, project.LabelPickerOptions{})
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	labels, on = labels[:1], on[:1]
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if !on[0] {
		t.Fatal("Enter on the one label left did not toggle it")
	}
}
