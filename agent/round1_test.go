package agent_test

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/agent"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestPromptBoxHoldsUnderAVeil(t *testing.T) {
	wd := newWorld()
	draft, cast := "", 0
	tt := view(600, 300, wd, func(c *ui.Context) {
		if agent.PromptBox(c, "Prompt", &draft, agent.PromptBoxOptions{}).Cast != "" {
			cast++
		}
	})
	click(t, tt, "Prompt")
	tt.Type("wait")
	wd.inert = true
	tt.Key(0, ui.KeyEnter)
	if cast != 0 || draft != "wait" {
		t.Fatalf("a veiled box cast %d times, draft %q", cast, draft)
	}
}

func TestFocusedCellsKeepTheirArt(t *testing.T) {
	n := 0
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		agent.MemoryCrystals(c, "Memory", []string{"a"}, agent.MemoryCrystalsOptions{})
		agent.ModelSelect(c, "Models", []agent.Model{{Name: "A", Speed: 1, Power: 1, Cost: 1}}, &n, agent.ModelSelectOptions{})
	})
	tt.Key(0, ui.KeyTab)
	x, y, _, _ := dots(find(t, tt, "Memory"))
	// The cursor lights the cell; the art shows on it.
	if !tt.Focused("Memory") || count(tt.Image(), x, y, 26, 30, pixel.Magic) == 0 || count(tt.Image(), x, y, 26, 30, pixel.Accent) == 0 {
		t.Error("the cursor hides the crystal, or is missing")
	}
	tt.Key(0, ui.KeyTab)
	x, y, _, _ = dots(find(t, tt, "Models"))
	if !tt.Focused("Models") || count(tt.Image(), x, y, 30, 30, pixel.Ink) == 0 || count(tt.Image(), x, y, 30, 30, pixel.Lit) < 200 {
		t.Error("the cursor hides the portrait, or is missing")
	}
}

func TestExtremesStayInside(t *testing.T) {
	long, at := strings.Repeat("W", 80), time.Unix(1, 0)
	quest := pixel.Measure("QUEST COMPLETE!")
	saves := 12 + 4 + 96 + 4 + pixel.Measure("RESTORE") + 8
	for _, tc := range []struct {
		label string
		build func(c *ui.Context)
	}{
		{long, func(c *ui.Context) {
			agent.AgentCard(c, agent.Agent{Name: long, Class: long, Level: 99999, Used: math.MaxInt, Window: math.MaxInt}, agent.AgentCardOptions{})
		}},
		{"CAST " + long, func(c *ui.Context) { agent.ToolCall(c, long, agent.Failed, agent.ToolCallOptions{Args: long}) }},
		{"Q", func(c *ui.Context) {
			agent.PlanQuest(c, "Q", []agent.PlanStep{{Text: long, Done: true}}, agent.PlanQuestOptions{Width: quest})
		}},
		{"M", func(c *ui.Context) { agent.TokenMana(c, "M", math.MaxInt, math.MaxInt, agent.TokenManaOptions{}) }},
		{"A", func(c *ui.Context) {
			agent.ApprovalPrompt(c, "A", long+long, agent.ApprovalPromptOptions{Danger: true})
		}},
		{"P", func(c *ui.Context) {
			agent.SubAgentParty(c, "P", []agent.Member{{Name: long, Summoned: []agent.Member{{Name: long}}}}, agent.SubAgentPartyOptions{})
		}},
		{"S", func(c *ui.Context) { agent.AgentStatus(c, "S", agent.Blocked, agent.AgentStatusOptions{}) }},
		{"D", func(c *ui.Context) {
			agent.DiffScroll(c, "D", []agent.FileDiff{{Path: long, Added: math.MaxInt, Removed: math.MaxInt}}, agent.DiffScrollOptions{Width: 400})
		}},
		{"C", func(c *ui.Context) { agent.CostCoins(c, "C", math.MaxInt, agent.CostCoinsOptions{}) }},
		{"V", func(c *ui.Context) {
			agent.Checkpoints(c, "V", []agent.Checkpoint{{Name: long, At: at}}, agent.CheckpointsOptions{Width: saves})
		}},
	} {
		tt := view(1400, 700, newWorld(), tc.build)
		x, y, w, h := dots(find(t, tt, tc.label))
		if n := spill(tt.Image(), x, y, w, h); n > 0 {
			t.Errorf("%.12s spills %d dots", tc.label, n)
		}
	}
}

func TestTooNarrowPanics(t *testing.T) {
	for name, build := range map[string]func(c *ui.Context){
		"level": func(c *ui.Context) {
			agent.AgentCard(c, agent.Agent{Name: "R", Class: "C", Level: math.MaxInt, Window: 1}, agent.AgentCardOptions{})
		},
		"quest": func(c *ui.Context) {
			agent.PlanQuest(c, "Q", []agent.PlanStep{{Text: "a", Done: true}}, agent.PlanQuestOptions{Width: 64})
		},
		"saves": func(c *ui.Context) { agent.Checkpoints(c, "V", nil, agent.CheckpointsOptions{Width: 1}) },
	} {
		mustPanic(t, name, func() { view(600, 300, newWorld(), build) })
	}
}
