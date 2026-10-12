package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestGameDemosLogTheirActions(t *testing.T) {
	for _, tc := range []struct {
		path, label string
		key         ui.Key
	}{
		{"/c/248-expbar", "SLAY A SLIME", 0}, {"/c/252-scoreboard", "RESET", 0}, {"/c/253-timerhud", "SKIP TO 33", 0},
		{"/c/257-equipment", "Gear", ui.KeyRight}, {"/c/258-skilltree", "Talents", ui.KeyDown}, {"/c/261-bossbar", "ATTACK", 0},
	} {
		a, tt := open(tc.path)
		slug := tc.path[3:]
		if tc.label == "RESET" {
			click(t, tt, "STOMP")
		}
		before := len(a.logs[slug])
		click(t, tt, tc.label)
		// A key moves a control the click has focused.
		if tc.key != 0 {
			before = len(a.logs[slug])
			tt.Key(0, tc.key)
		}
		if len(a.logs[slug]) <= before {
			t.Errorf("%s: %s logged nothing: %q", slug, tc.label, a.logs[slug])
		}
	}
}

func TestHotbarDemoLogsAnEmptySlotUsed(t *testing.T) {
	a, tt := open("/c/256-hotbar")
	click(t, tt, "Belt")
	tt.Key(0, ui.Key6)
	before := len(a.logs["256-hotbar"])
	tt.Key(0, ui.KeyEnter)
	if len(a.logs["256-hotbar"]) <= before {
		t.Fatalf("using empty slot 6 logged nothing: %q", a.logs["256-hotbar"])
	}
}
