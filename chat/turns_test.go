package chat_test

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/chat"
)

func TestKeysInAVeilsFirstFrameTurnNothing(t *testing.T) {
	cases := map[string]func(c *ui.Context, v *string){
		"Reactions": func(c *ui.Context, v *string) {
			rs := []chat.Reaction{{Emoji: "love", Count: 1}}
			fmt.Sscan(*v, &rs[0].Count)
			chat.Reactions(c, "Reactions", rs, chat.ReactionsOptions{})
			*v = fmt.Sprint(rs[0].Count)
		},
		"CallControls": func(c *ui.Context, v *string) {
			muted := *v == "true"
			chat.CallControls(c, "Call", &muted, new(bool), chat.CallControlsOptions{})
			*v = fmt.Sprint(muted)
		},
	}
	for name, build := range cases {
		for _, veil := range []bool{false, true} {
			start := map[string]string{"Reactions": "1", "CallControls": "false"}[name]
			v, wd := start, &world{clk: newClock()}
			tt := view(400, 200, wd, func(c *ui.Context) { build(c, &v) })
			tt.Key(0, ui.KeyTab)
			wd.inert = veil
			tt.Key(0, ui.KeySpace)
			if (v != start) == veil {
				t.Errorf("%s: veiled %v, Space made %s", name, veil, v)
			}
		}
	}
}

func TestComposerSendsNothingInAVeilsFirstFrame(t *testing.T) {
	for _, veil := range []bool{false, true} {
		draft, sent, wd := "", "", &world{clk: newClock()}
		tt := view(400, 200, wd, func(c *ui.Context) {
			if r := chat.Composer(c, "Say", &draft, chat.ComposerOptions{}); r.Sent != "" {
				sent = r.Sent
			}
		})
		tt.Key(0, ui.KeyTab)
		tt.Type("hi")
		wd.inert = veil
		tt.Key(0, ui.KeyEnter)
		if (sent == "hi") == veil {
			t.Errorf("veiled %v, Enter sent %q", veil, sent)
		}
	}
}
