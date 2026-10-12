package project_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/project"
)

func TestPriorityIconsRiseThenBurn(t *testing.T) {
	levels := []project.Level{project.Low, project.Medium, project.High, project.Critical}
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		ui.Row(c).Children(func() {
			for _, l := range levels {
				project.Priority(c, l, project.PriorityOptions{})
			}
		})
	})
	img := tt.Image()
	// LOW lights one bar, MEDIUM two, HIGH three.
	for _, tc := range []struct {
		name string
		slot pixel.Slot
		dots int
	}{{"LOW", pixel.Good, 12}, {"MEDIUM", pixel.Warn, 33}, {"HIGH", pixel.Bad, 63}} {
		x, y, _, _ := dots(find(t, tt, tc.name))
		if n := count(img, x, y, 12, 12, tc.slot); n != tc.dots {
			t.Errorf("%s has %d dots of %v, want %d", tc.name, n, tc.slot, tc.dots)
		}
	}
	// CRITICAL is the flame: gold at its heart.
	x, y, _, _ := dots(find(t, tt, "CRITICAL"))
	if count(img, x, y, 12, 12, pixel.Warn) == 0 {
		t.Fatal("CRITICAL is not on fire")
	}
	mustPanic(t, "a fifth level", func() {
		view(400, 200, newWorld(), func(c *ui.Context) { project.Priority(c, 4, project.PriorityOptions{}) })
	})
}

func TestIssueCardOpens(t *testing.T) {
	clicked := 0
	is := project.Issue{ID: "#42", Title: "Save lost", Labels: []project.Label{{Name: "BUG", Slot: pixel.Bad}}, Assignee: "Brom", Priority: project.Critical}
	tt := view(400, 300, newWorld(), func(c *ui.Context) {
		if project.IssueCard(c, is, project.IssueCardOptions{}).Clicked {
			clicked++
		}
	})
	click(t, tt, "#42 Save lost")
	if clicked != 1 {
		t.Fatalf("a click opened %d times", clicked)
	}
	x, y, _, h := dots(find(t, tt, "#42 Save lost"))
	// BUG's swatch is Bad on the labels row.
	if count(tt.Image(), x+4, y+h-24, 8, 16, pixel.Bad) == 0 {
		t.Fatal("the BUG label has no Bad swatch")
	}
}

func TestWorkloadBarWarnsAndOverflows(t *testing.T) {
	loads := []project.Load{{Name: "Aria", Hours: 30}, {Name: "Brom", Hours: 36}, {Name: "Cass", Hours: 44}}
	tt := view(700, 300, newWorld(), func(c *ui.Context) {
		project.WorkloadBar(c, "Load", loads, project.WorkloadBarOptions{})
	})
	// Meters start past the face, the names and gaps.
	nameW := max(pixel.Measure("Aria"), pixel.Measure("Brom"), pixel.Measure("Cass"))
	for _, tc := range []struct {
		name string
		slot pixel.Slot
	}{{"Aria", pixel.Good}, {"Brom", pixel.Warn}, {"Cass", pixel.Bad}} {
		x, y, _, _ := dots(find(t, tt, tc.name))
		if count(tt.Image(), x+24+nameW, y, 100, 16, tc.slot) == 0 {
			t.Errorf("%s's meter shows no %v", tc.name, tc.slot)
		}
	}
	// Six dots of meter is too few to read.
	narrow := 16 + nameW + pixel.Measure("44/40 H") + 12 + 6
	mustPanic(t, "a six-dot meter", func() {
		view(700, 300, newWorld(), func(c *ui.Context) {
			project.WorkloadBar(c, "L", loads, project.WorkloadBarOptions{Width: narrow})
		})
	})
}

// offPalette counts dots no palette color matches.
func offPalette(img *image.RGBA) int {
	b := img.Bounds()
	return count(img, 0, 0, b.Dx()/2, b.Dy()/2, none)
}

func TestProjectStaysOnThePaletteWhenOff(t *testing.T) {
	builds := map[string]func(c *ui.Context, off bool){
		"board": func(c *ui.Context, off bool) {
			cols := [3][]string{{"A"}, {"B"}, {}}
			project.QuestBoard(c, "X", &cols, project.QuestBoardOptions{Disabled: off})
		},
		"task": func(c *ui.Context, off bool) {
			tk := project.Task{Title: "A", Assignee: "Brom"}
			project.TaskItem(c, &tk, project.TaskItemOptions{Disabled: off})
		},
		"list": func(c *ui.Context, off bool) {
			project.TaskList(c, "X", []project.Task{{Title: "A"}}, project.TaskListOptions{Disabled: off})
		},
		"tomato": func(c *ui.Context, off bool) {
			project.Pomodoro(c, "X", project.PomodoroOptions{Disabled: off})
		},
		"party": func(c *ui.Context, off bool) {
			i := 0
			project.AssigneePicker(c, "X", []string{"Aria"}, &i, project.AssigneePickerOptions{Disabled: off})
		},
		"labels": func(c *ui.Context, off bool) {
			project.LabelPicker(c, "X", []project.Label{{Name: "A", Slot: pixel.Bad}}, []bool{true}, project.LabelPickerOptions{Disabled: off})
		},
		"time": func(c *ui.Context, off bool) {
			es := []project.TimeEntry{{Task: "A", Start: time.Unix(0, 0)}}
			project.TimeTracker(c, "X", "A", &es, project.TimeTrackerOptions{Disabled: off})
		},
		"issue": func(c *ui.Context, off bool) {
			project.IssueCard(c, project.Issue{ID: "#1", Title: "A", Assignee: "Brom"}, project.IssueCardOptions{Disabled: off})
		},
		"checklist": func(c *ui.Context, off bool) {
			project.Checklist(c, "X", []string{"A"}, []bool{true}, project.ChecklistOptions{Disabled: off})
		},
	}
	for name, build := range builds {
		for _, inert := range []bool{false, true} {
			wd := newWorld()
			wd.inert = inert
			tt := view(700, 400, wd, func(c *ui.Context) { build(c, !inert) })
			if n := offPalette(tt.Image()); n != 0 {
				t.Errorf("%s, inert %v: %d dots off the palette", name, inert, n)
			}
		}
	}
}
