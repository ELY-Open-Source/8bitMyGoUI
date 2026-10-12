package gallery

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

var letters = map[string]ui.Key{"A": ui.KeyA, "Right": ui.KeyRight, "Up": ui.KeyUp, "Down": ui.KeyDown, "Enter": ui.KeyEnter, "Tab": ui.KeyTab}

func TestMediaDemosLogTheirActions(t *testing.T) {
	for _, tc := range []struct {
		path  string
		steps []string
	}{
		{"/c/346-audioplayer", []string{"PLAY", "NEXT"}}, {"/c/347-transport", []string{"PLAY", "STOP", "PREVIOUS"}},
		{"/c/348-scrubber", []string{"Seek"}}, {"/c/349-volume", []string{"MUTE"}}, {"/c/350-equalizer", []string{"1K", "key:Up"}},
		{"/c/352-playlist", []string{"Save Point Lullaby, Triangle Wave", "key:Enter"}}, {"/c/353-cassette", []string{"PLAY"}},
		{"/c/355-vinyl", []string{"PLAY"}}, {"/c/356-pianokeys", []string{"Piano", "key:A"}}, {"/c/357-stepsequencer", []string{"Beat", "PLAY"}},
		{"/c/358-tracker", []string{"Pattern", "key:A"}}, {"/c/360-videoframe", []string{"PLAY"}},
		{"/c/361-imagecompare", []string{"Sketch and painting", "key:Right"}}, {"/c/362-gallery", []string{"Views"}},
		{"/c/364-recordbutton", []string{"Record, 0:00"}}, {"/c/365-tuner", []string{"Radio"}},
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
			// Clicking a field only to focus it logs nothing.
			if len(a.logs[slug]) <= before && step != "1K" && step != "Pattern" && step != "Sketch and painting" {
				t.Errorf("%s: %s logged nothing: %q", slug, step, a.logs[slug])
			}
		}
	}
}

func TestCodeDemosLogTheirActions(t *testing.T) {
	quiet := map[string]bool{"COMMAND": true}
	for _, tc := range []struct {
		path  string
		steps []string
	}{
		{"/c/366-terminal", []string{"COMMAND", "type:help", "key:Enter", "key:Up"}}, {"/c/367-codeblock", []string{"COPY"}},
		{"/c/368-logviewer", []string{"WARN", "+ LINE"}}, {"/c/369-hexviewer", []string{"save.dat"}},
		{"/c/371-jsontree", []string{"save.json"}}, {"/c/373-crashscreen", []string{"CONTINUE?", "CRASH AGAIN"}},
		{"/c/376-breakpoints", []string{"jump.go"}}, {"/c/377-variables", []string{"Locals"}},
		{"/c/378-commandhistory", []string{"run", "key:Enter"}}, {"/c/380-debuggrid", []string{"ALIGNED?"}},
		{"/c/381-pixeleditor", []string{"Heart", "COPY"}}, {"/c/382-paletteeditor", []string{"R"}},
		{"/c/383-spritesheet", []string{"PLAY"}}, {"/c/384-glyphmap", []string{"Font"}}, {"/c/385-registerview", []string{"STEP"}},
	} {
		a, tt := open(tc.path)
		slug := tc.path[3:]
		for _, step := range tc.steps {
			before := len(a.logs[slug])
			switch {
			case strings.HasPrefix(step, "key:"):
				tt.Key(0, letters[step[4:]])
			case strings.HasPrefix(step, "type:"):
				tt.Type(step[5:])
			default:
				click(t, tt, step)
			}
			if len(a.logs[slug]) <= before && !quiet[step] && !strings.HasPrefix(step, "type:") {
				t.Errorf("%s: %s logged nothing: %q", slug, step, a.logs[slug])
			}
		}
	}
}

func TestAccountDemosLogTheirActions(t *testing.T) {
	quiet := map[string]bool{"NAME": true, "PASSWORD": true, "Login code": true, "Controls": true, "4/5": true, "type:player1": true, "type:hunter2": true}
	for _, tc := range []struct {
		path  string
		steps []string
	}{
		{"/c/386-settingslist", []string{"Music"}}, {"/c/387-optioncycler", []string{"Text speed"}},
		{"/c/388-keybindings", []string{"Controls", "key:A"}}, {"/c/391-charactercreator", []string{"HAIR"}},
		{"/c/392-loginform", []string{"NAME", "type:player1", "PASSWORD", "type:hunter2", "key:Enter"}},
		{"/c/393-passwordgrid", []string{"SYMBOLS"}}, {"/c/394-sessionlist", []string{"SIGN OUT Old phone"}},
		{"/c/395-languagepicker", []string{"Deutsch"}}, {"/c/396-palettepicker", []string{"Palette"}},
		{"/c/397-volumemixer", []string{"BGM"}}, {"/c/398-difficultyselect", []string{"Difficulty"}},
		{"/c/399-accessibilityoptions", []string{"REDUCE MOTION"}}, {"/c/400-accountswitcher", []string{"P2, Brom"}},
		{"/c/401-playerselect", []string{"Players"}}, {"/c/402-twofactor", []string{"Login code", "type:1"}},
		{"/c/405-feedbackform", []string{"4/5", "SEND ▶"}},
	} {
		a, tt := open(tc.path)
		slug := tc.path[3:]
		for _, step := range tc.steps {
			before := len(a.logs[slug])
			switch {
			case strings.HasPrefix(step, "key:"):
				tt.Key(0, letters[step[4:]])
			case strings.HasPrefix(step, "type:"):
				tt.Type(step[5:])
			default:
				click(t, tt, step)
			}
			if len(a.logs[slug]) <= before && !quiet[step] {
				t.Errorf("%s: %s logged nothing: %q", slug, step, a.logs[slug])
			}
		}
	}
}
