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

func TestScreenDemosLogTheirActions(t *testing.T) {
	for _, tc := range []struct {
		path  string
		steps []string
	}{
		{"/c/269-saveslots", []string{"Save files"}}, {"/c/270-titlescreen", []string{"PRESS START"}}, {"/c/271-pressstart", []string{"PRESS START"}},
		{"/c/272-optionsmenu", []string{"Options", "Right"}}, {"/c/273-gameover", []string{"Enter"}}, {"/c/274-credits", []string{"ROLL AGAIN"}},
		{"/c/275-endingcard", []string{"AGAIN"}}, {"/c/277-letterbox", []string{"CUTSCENE"}}, {"/c/278-missionbrief", []string{"START"}},
		{"/c/279-achievements", []string{"Awards", "Right"}}, {"/c/280-resultscreen", []string{"AGAIN"}}, {"/c/281-crafting", []string{"CRAFT"}},
		{"/c/282-lootpopup", []string{"OPEN THE BOX", "OK"}}, {"/c/283-treasurechest", []string{"Chest", "CLOSE"}}, {"/c/284-gachareel", []string{"SPIN"}},
		{"/c/285-interactprompt", []string{"NEAR THE MAYOR"}}, {"/c/286-turnorder", []string{"NEXT TURN"}}, {"/c/287-battlelog", []string{"ATTACK", "HEAL"}},
		{"/c/288-targetselector", []string{"Targets", "Right", "Enter"}}, {"/c/290-staminabar", []string{"SPRINT"}},
	} {
		a, tt := open(tc.path)
		slug := tc.path[3:]
		for _, step := range tc.steps {
			before := len(a.logs[slug])
			switch step {
			case "Right":
				tt.Key(0, ui.KeyRight)
			case "Enter":
				tt.Key(0, ui.KeyEnter)
			default:
				click(t, tt, step)
			}
			// A focusing click may change nothing worth a line.
			if len(a.logs[slug]) <= before && step != "Options" && step != "Awards" && step != "Targets" {
				t.Errorf("%s: %s logged nothing: %q", slug, step, a.logs[slug])
			}
		}
	}
}
