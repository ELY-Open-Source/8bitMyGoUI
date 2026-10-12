package gallery

import (
	"strings"
	"testing"
)

func TestFilesDemosLogTheirActions(t *testing.T) {
	// Clicks that only focus a component may log nothing.
	quiet := map[string]bool{"FOLDERS": true}
	for _, tc := range []struct {
		path  string
		steps []string
	}{
		{"/c/406-fileexplorer", []string{"FOLDERS", "key:Up"}}, {"/c/407-filegrid", []string{"SAVE1.SAV", "key:Right", "key:Enter"}},
		{"/c/409-pathbar", []string{"GAMES", "BACK DOWN"}}, {"/c/410-memorycard", []string{"DELETE HERO QUEST"}},
		{"/c/412-uploadqueue", []string{"CANCEL BOSS.PNG"}}, {"/c/413-trash", []string{"EMPTY", "THROW AWAY"}},
		{"/c/414-recentfiles", []string{"Recent", "key:Down", "key:Enter"}}, {"/c/415-filepreview", []string{"SUNSET.PNG"}},
		{"/c/416-archiveviewer", []string{"LOOT.ZIP", "key:Down", "SIZE"}}, {"/c/418-sharepanel", []string{"COPY", "EDIT"}},
	} {
		a, tt := open(tc.path)
		slug := tc.path[3:]
		for _, step := range tc.steps {
			before := len(a.logs[slug])
			if k, ok := strings.CutPrefix(step, "key:"); ok {
				tt.Key(0, letters[k])
			} else {
				click(t, tt, step)
			}
			if len(a.logs[slug]) <= before && !quiet[step] {
				t.Errorf("%s: %s logged nothing: %q", slug, step, a.logs[slug])
			}
		}
	}
}
