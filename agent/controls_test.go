package agent_test

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/agent"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestApprovalPromptAnswersAndKeepsTheCommandSecret(t *testing.T) {
	allowed, denied, off := 0, 0, false
	cmd := "curl -H 'Authorization: Bearer SECRET' https://example.com/deploy && rm -rf build/cache/old/artifacts"
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		r := agent.ApprovalPrompt(c, "Approval", cmd, agent.ApprovalPromptOptions{Danger: true, Disabled: off})
		if r.Allowed {
			allowed++
		}
		if r.Denied {
			denied++
		}
	})
	log := logged(func() {
		click(t, tt, "YES")
		click(t, tt, "NO")
	})
	if allowed != 1 || denied != 1 || !strings.Contains(log, "allowed=true") || strings.Contains(log, "SECRET") {
		t.Fatalf("allowed %d, denied %d, log:\n%s", allowed, denied, log)
	}
	// The whole command reads out; Danger inks it Bad.
	x, y, w, h := dots(find(t, tt, cmd))
	if count(tt.Image(), x, y, w, h, pixel.Bad) == 0 || h != 3*16+6 {
		t.Fatalf("the command well is %d dots high, or not Bad", h)
	}
	off = true
	tt.Frame()
	click(t, tt, "YES")
	if allowed != 1 {
		t.Fatal("a disabled prompt took a YES")
	}
}

func TestAgentLogTypesTheNewestStep(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = reduce
		steps := []agent.LogStep{{Kind: agent.Cast, Text: "read_file"}}
		tt := view(600, 300, wd, func(c *ui.Context) { agent.AgentLog(c, "Log", "ROBO", steps, agent.AgentLogOptions{}) })
		ink := func(label string, s pixel.Slot) int {
			x, y, w, h := dots(find(t, tt, label))
			return count(tt.Image(), x, y, w, h, s)
		}
		// Opened, the log shows whole: nothing types.
		if ink("ROBO casts read_file!", pixel.Info) < 20 {
			t.Fatal("the first line is not shown whole")
		}
		steps = append(steps, agent.LogStep{Kind: agent.Hit, Text: "3 files found"})
		tt.Frame()
		first := ink("A hit! 3 files found", pixel.Good)
		later(tt, wd, time.Second)
		full := ink("A hit! 3 files found", pixel.Good)
		if full == 0 || reduce != (first == full) {
			t.Errorf("reduce %v: %d dots typed at once, %d after a second", reduce, first, full)
		}
	}
}

func TestMemoryCrystalsBrowseAndPick(t *testing.T) {
	memories := []string{"Prefers tabs", "Ships on Fridays"}
	picked, off := -1, false
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		if r := agent.MemoryCrystals(c, "Memory", memories, agent.MemoryCrystalsOptions{Slots: 4, Disabled: off}); r.Picked >= 0 {
			picked = r.Picked
		}
	})
	x, y, _, _ := dots(find(t, tt, "Memory"))
	img := tt.Image()
	for i := range 4 {
		if held := count(img, x+i*26, y, 26, 30, pixel.Magic) > 0; held != (i < 2) {
			t.Errorf("crystal %d lit is %v", i, held)
		}
	}
	// A click picks the crystal under it; so does Enter.
	click(t, tt, "Memory")
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyRight)
	if picked != 2 {
		t.Fatalf("the click picked %d", picked)
	}
	tt.Key(0, ui.KeyEnter)
	if picked != 1 {
		t.Fatalf("Home, Right, Enter picked %d", picked)
	}
	off = true
	tt.Frame()
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	if picked != 1 {
		t.Fatalf("a disabled shelf picked %d", picked)
	}
}

func TestModelSelectPicksLikeAFighter(t *testing.T) {
	models := []agent.Model{{Name: "SWIFT", Speed: 5, Power: 2, Cost: 1}, {Name: "TITAN", Speed: 2, Power: 5, Cost: 5}, {Name: "SAGE", Speed: 3, Power: 4, Cost: 3}}
	chosen, changed, off := 0, 0, false
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		if agent.ModelSelect(c, "Models", models, &chosen, agent.ModelSelectOptions{Disabled: off}).Changed {
			changed++
		}
	})
	click(t, tt, "Models")
	if chosen != 1 || changed != 1 {
		t.Fatalf("the click chose %d, %d changes", chosen, changed)
	}
	tt.Key(0, ui.KeyRight)
	x, y, w, _ := dots(find(t, tt, "Models"))
	img := tt.Image()
	// The panel shows SAGE, under the cursor: three SPD pips.
	for j := range 5 {
		if want := j < 3; (at(img, x+w+8+50+j*12+1, y+46+6) == pixel.Info) != want {
			t.Errorf("SPD pip %d lit is not %v", j, want)
		}
	}
	tt.Key(0, ui.KeyEnter)
	if chosen != 2 || changed != 2 || at(tt.Image(), x+60+3, y+3) != pixel.Accent {
		t.Fatalf("Right and Enter chose %d, %d changes, or no Accent corner", chosen, changed)
	}
	off = true
	tt.Frame()
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyEnter)
	if chosen != 2 {
		t.Fatalf("a disabled select chose %d", chosen)
	}
}

func TestPromptBoxCastsOnEnterOrCast(t *testing.T) {
	draft, cast := "", []string{}
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		if r := agent.PromptBox(c, "Prompt", &draft, agent.PromptBoxOptions{}); r.Cast != "" {
			cast = append(cast, r.Cast)
		}
	})
	click(t, tt, "CAST")
	click(t, tt, "Prompt")
	log := logged(func() {
		tt.Type("fix")
		tt.Key(ui.Shift, ui.KeyEnter)
		tt.Type("it")
		tt.Key(0, ui.KeyEnter)
	})
	if len(cast) != 1 || cast[0] != "fix\nit" || draft != "" {
		t.Fatalf("cast %q, draft %q", cast, draft)
	}
	if !strings.Contains(log, "runes=6") || strings.Contains(log, "fix") {
		t.Fatalf("the log should count, not quote:\n%s", log)
	}
	tt.Type("ship")
	click(t, tt, "CAST")
	if len(cast) != 2 || cast[1] != "ship" || draft != "" {
		t.Fatalf("cast %q, draft %q", cast, draft)
	}
}

func TestCheckpointsRestoreByName(t *testing.T) {
	wd := newWorld()
	now := wd.clk.t
	points := []agent.Checkpoint{{Name: "Before refactor", At: now.Add(-2 * time.Hour)}, {Name: "Tests green", At: now.Add(-5 * time.Minute)}}
	restored, off := -1, false
	tt := view(600, 300, wd, func(c *ui.Context) {
		if r := agent.Checkpoints(c, "Saves", points, agent.CheckpointsOptions{Disabled: off}); r.Restored >= 0 {
			restored = r.Restored
		}
	})
	find(t, tt, "5 MIN AGO")
	click(t, tt, "RESTORE Tests green")
	if restored != 1 {
		t.Fatalf("restored %d", restored)
	}
	// Rows are keyed: losing the first keeps the focus put.
	points = points[1:]
	tt.Frame()
	if !tt.Focused("RESTORE Tests green") {
		t.Fatal("the focus left with the removed row")
	}
	off, restored = true, -1
	tt.Frame()
	click(t, tt, "RESTORE Tests green")
	if restored != -1 {
		t.Fatal("a disabled list restored")
	}
	points = nil
	tt.Frame()
	find(t, tt, "NO ITEMS")
}
