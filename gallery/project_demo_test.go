package gallery

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestProjectDemosLogTheirActions(t *testing.T) {
	// Clicks that only focus a component may log nothing.
	quiet := map[string]bool{"Quests": true}
	for _, tc := range []struct {
		path  string
		steps []string
	}{
		{"/c/421-questboard", []string{"Quests", "alt:Right"}}, {"/c/422-taskitem", []string{"Fix the save bug"}},
		{"/c/423-tasklist", []string{"OPEN"}}, {"/c/426-pomodoro", []string{"START", "PAUSE"}},
		{"/c/428-assigneepicker", []string{"Brom"}}, {"/c/429-labelpicker", []string{"Labels"}},
		{"/c/430-timetracker", []string{"START", "STOP"}}, {"/c/433-issuecard", []string{"#42 Save file vanishes after the boss"}},
		{"/c/434-checklist", []string{"Tell mom"}},
	} {
		a, tt := open(tc.path)
		slug := tc.path[3:]
		for _, step := range tc.steps {
			before := len(a.logs[slug])
			if k, ok := strings.CutPrefix(step, "alt:"); ok {
				tt.Key(ui.Alt, letters[k])
			} else {
				click(t, tt, step)
			}
			if len(a.logs[slug]) <= before && !quiet[step] {
				t.Errorf("%s: %s logged nothing: %q", slug, step, a.logs[slug])
			}
		}
	}
}
