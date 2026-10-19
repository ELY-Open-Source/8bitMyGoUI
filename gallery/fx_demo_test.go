package gallery

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestFxDemosLogTheirActions(t *testing.T) {
	for _, tc := range []struct {
		path  string
		steps []string
	}{
		{"/c/293-sunmoon", []string{"Hour", "Right"}}, {"/c/306-lightning", []string{"STRIKE"}},
		{"/c/308-seasonscene", []string{"Winter", "Hour", "Right"}}, {"/c/311-screenshake", []string{"SHAKE"}},
		{"/c/312-flash", []string{"FLASH"}}, {"/c/313-fadedither", []string{"FADE OUT"}}, {"/c/314-iriswipe", []string{"CLOSE"}},
		{"/c/315-dissolve", []string{"DISSOLVE"}}, {"/c/316-mosaic", []string{"PIXELATE"}}, {"/c/317-glitch", []string{"GLITCH"}},
		{"/c/319-explosion", []string{"BOOM"}},
	} {
		a, tt := open(tc.path)
		slug := tc.path[3:]
		for _, step := range tc.steps {
			before := len(a.logs[slug])
			if step == "Right" {
				tt.Key(0, ui.KeyRight)
			} else {
				click(t, tt, step)
			}
			// Focusing the stepper changes nothing worth a line.
			if len(a.logs[slug]) <= before && step != "Hour" {
				t.Errorf("%s: %s logged nothing: %q", slug, step, a.logs[slug])
			}
		}
	}
}

func TestGroupsAreOneTabStopEach(t *testing.T) {
	for _, page := range []struct{ path, query string }{{"/", ""}, {"/cat/fx", ""}, {"/", "e"}} {
		_, tt := open(page.path)
		tt.Key(0, ui.KeySlash)
		tt.Type(page.query)
		// Search, four top-bar controls, then one stop per group.
		stops := 1
		for tt.Key(0, ui.KeyTab); !tt.Focused("Search components"); tt.Key(0, ui.KeyTab) {
			if stops++; stops > 12 {
				t.Fatalf("%s %q: a Tab cycle passes over 12 stops", page.path, page.query)
			}
		}
	}
	_, tt := open("/")
	for range 40 {
		if tt.Focused("Spring Day, in use") {
			break
		}
		tt.Key(0, ui.KeyTab)
	}
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Spring Night") {
		t.Fatal("Right did not move along the palette strip")
	}
}

func TestHomeSceneFollowsThePaletteStrip(t *testing.T) {
	a, tt := open("/")
	if _, ok := tt.Find("Spring Day scene"); !ok {
		t.Fatalf("no spring scene in %q", tt.Texts())
	}
	click(t, tt, "Winter Night")
	if a.palette != 7 {
		t.Fatalf("the strip chose palette %d", a.palette)
	}
	r, ok := tt.Find("Winter Night scene")
	if !ok {
		t.Fatalf("the scene kept its season: %q", tt.Texts())
	}
	// The sky's top row is Winter Night's Well, snow aside.
	want, well, n := pixel.WinterNight.Colors[pixel.Well], 0, 0
	for x := int(r.X); x < int(r.X+r.W); x += 2 {
		got := tt.Image().RGBAAt(x, int(r.Y)+1)
		if got.R == want.R && got.G == want.G && got.B == want.B {
			well++
		}
		n++
	}
	if well*4 < n*3 {
		t.Fatalf("only %d of %d sky dots are Winter Night's Well", well, n)
	}
	if _, ok := tt.Find("Winter Night, in use"); !ok {
		t.Fatal("the strip does not mark the palette in use")
	}
}

func TestChatDemosLogTheirActions(t *testing.T) {
	for _, tc := range []struct {
		path  string
		steps []string
	}{
		{"/c/322-chatlog", []string{"SAY HI"}}, {"/c/323-composer", []string{"Message", "type:hi", "Enter"}},
		{"/c/325-chatlist", []string{"Brom"}}, {"/c/327-reactions", []string{"fire 5"}}, {"/c/329-mentioninput", []string{"Mention", "type:@a"}},
		{"/c/328-emojipicker", []string{"Emoji"}}, {"/c/333-channellist", []string{"quests", "Enter"}}, {"/c/335-usercard", []string{"INVITE"}},
		{"/c/336-mailbox", []string{"Aria: Re: boss fight, NEW MESSAGES", "Enter"}}, {"/c/339-notificationbell", []string{"NEWS!", "Notifications"}},
		{"/c/340-callcontrols", []string{"MUTE", "HANG UP"}}, {"/c/342-attachmentchip", []string{"REMOVE notes.txt"}},
		{"/c/343-linkpreview", []string{"MyGo UI: native widgets, written in Go"}}, {"/c/344-poll", []string{"Bard, 1"}},
		{"/c/345-stickerpicker", []string{"Stickers"}},
	} {
		a, tt := open(tc.path)
		slug := tc.path[3:]
		for _, step := range tc.steps {
			before := len(a.logs[slug])
			switch {
			case step == "Enter":
				tt.Key(0, ui.KeyEnter)
			case strings.HasPrefix(step, "type:"):
				tt.Type(step[5:])
			default:
				click(t, tt, step)
			}
			// Focusing a field changes nothing worth a line.
			if len(a.logs[slug]) <= before && step != "Message" && step != "Mention" {
				t.Errorf("%s: %s logged nothing: %q", slug, step, a.logs[slug])
			}
		}
	}
}
