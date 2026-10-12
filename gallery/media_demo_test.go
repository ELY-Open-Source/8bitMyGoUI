package gallery

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

var letters = map[string]ui.Key{"A": ui.KeyA, "Right": ui.KeyRight, "Up": ui.KeyUp, "Enter": ui.KeyEnter}

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
