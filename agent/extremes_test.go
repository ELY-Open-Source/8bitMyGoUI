package agent_test

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/agent"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestBadInputPanics(t *testing.T) {
	big := pixel.MustSprite("#############\n" + "#############")
	ok := agent.Agent{Name: "R", Class: "C", Level: 1, Window: 10}
	deep := []agent.Member{{Name: "A", Summoned: []agent.Member{{Name: "B", Summoned: []agent.Member{{Name: "C"}}}}}}
	n := 0
	for name, build := range map[string]func(c *ui.Context){
		"over window": func(c *ui.Context) { a := ok; a.Used = 11; agent.AgentCard(c, a, agent.AgentCardOptions{}) },
		"level 0":     func(c *ui.Context) { a := ok; a.Level = 0; agent.AgentCard(c, a, agent.AgentCardOptions{}) },
		"no class":    func(c *ui.Context) { a := ok; a.Class = ""; agent.AgentCard(c, a, agent.AgentCardOptions{}) },
		"big face":    func(c *ui.Context) { agent.AgentCard(c, ok, agent.AgentCardOptions{Portrait: big}) },
		"no tool":     func(c *ui.Context) { agent.ToolCall(c, "", agent.OK, agent.ToolCallOptions{}) },
		"outcome":     func(c *ui.Context) { agent.ToolCall(c, "t", agent.Outcome(9), agent.ToolCallOptions{}) },
		"thin card":   func(c *ui.Context) { agent.ToolCall(c, "t", agent.OK, agent.ToolCallOptions{Width: 40}) },
		"no steps":    func(c *ui.Context) { agent.PlanQuest(c, "Q", nil, agent.PlanQuestOptions{}) },
		"blank step":  func(c *ui.Context) { agent.PlanQuest(c, "Q", []agent.PlanStep{{}}, agent.PlanQuestOptions{}) },
		"thin quest": func(c *ui.Context) {
			agent.PlanQuest(c, "Q", []agent.PlanStep{{Text: "a"}}, agent.PlanQuestOptions{Width: 30})
		},
		"negative mana": func(c *ui.Context) { agent.TokenMana(c, "M", -1, 10, agent.TokenManaOptions{}) },
		"no budget":     func(c *ui.Context) { agent.TokenMana(c, "M", 0, 0, agent.TokenManaOptions{}) },
		"thin mana":     func(c *ui.Context) { agent.TokenMana(c, "M", 0, 10, agent.TokenManaOptions{Width: 10}) },
		"no command":    func(c *ui.Context) { agent.ApprovalPrompt(c, "A", "", agent.ApprovalPromptOptions{}) },
		"thin prompt":   func(c *ui.Context) { agent.ApprovalPrompt(c, "A", "ls", agent.ApprovalPromptOptions{Width: 50}) },
		"no agent":      func(c *ui.Context) { agent.AgentLog(c, "L", "", nil, agent.AgentLogOptions{}) },
		"kind": func(c *ui.Context) {
			agent.AgentLog(c, "L", "R", []agent.LogStep{{Kind: agent.Kind(9), Text: "x"}}, agent.AgentLogOptions{})
		},
		"blank line": func(c *ui.Context) { agent.AgentLog(c, "L", "R", []agent.LogStep{{}}, agent.AgentLogOptions{}) },
		"nobody":     func(c *ui.Context) { agent.SubAgentParty(c, "P", nil, agent.SubAgentPartyOptions{}) },
		"nameless":   func(c *ui.Context) { agent.SubAgentParty(c, "P", []agent.Member{{}}, agent.SubAgentPartyOptions{}) },
		"too deep":   func(c *ui.Context) { agent.SubAgentParty(c, "P", deep, agent.SubAgentPartyOptions{Width: 80}) },
		"status":     func(c *ui.Context) { agent.AgentStatus(c, "S", agent.Status(7), agent.AgentStatusOptions{}) },
		"full shelf": func(c *ui.Context) {
			agent.MemoryCrystals(c, "M", []string{"a", "b"}, agent.MemoryCrystalsOptions{Slots: 1})
		},
		"blank mem": func(c *ui.Context) { agent.MemoryCrystals(c, "M", []string{""}, agent.MemoryCrystalsOptions{}) },
		"no files":  func(c *ui.Context) { agent.DiffScroll(c, "D", nil, agent.DiffScrollOptions{}) },
		"minus lines": func(c *ui.Context) {
			agent.DiffScroll(c, "D", []agent.FileDiff{{Path: "a", Added: -1}}, agent.DiffScrollOptions{})
		},
		"line flood": func(c *ui.Context) {
			agent.DiffScroll(c, "D", []agent.FileDiff{{Path: "a", Added: math.MaxInt}, {Path: "b", Added: 1}}, agent.DiffScrollOptions{Width: 2000})
		},
		"thin scroll": func(c *ui.Context) {
			agent.DiffScroll(c, "D", []agent.FileDiff{{Path: "a"}}, agent.DiffScrollOptions{Width: 60})
		},
		"chosen": func(c *ui.Context) {
			n := 1
			agent.ModelSelect(c, "M", []agent.Model{{Name: "A", Speed: 1, Power: 1, Cost: 1}}, &n, agent.ModelSelectOptions{})
		},
		"stat 6": func(c *ui.Context) {
			n := 0
			agent.ModelSelect(c, "M", []agent.Model{{Name: "A", Speed: 6, Power: 1, Cost: 1}}, &n, agent.ModelSelectOptions{})
		},
		"minus cost": func(c *ui.Context) { agent.CostCoins(c, "C", -1, agent.CostCoinsOptions{}) },
		"undated": func(c *ui.Context) {
			agent.Checkpoints(c, "S", []agent.Checkpoint{{Name: "A"}}, agent.CheckpointsOptions{})
		},
	} {
		n++
		mustPanic(t, name, func() { view(600, 300, newWorld(), build) })
	}
	if n != 32 {
		t.Fatalf("%d cases", n)
	}
}

func TestCheckpointsRefuseTwins(t *testing.T) {
	defer func() {
		if r := fmt.Sprint(recover()); !strings.Contains(r, "agent: Checkpoints") {
			t.Fatalf("twins panicked with %q", r)
		}
	}()
	at := time.Unix(1, 0)
	view(600, 300, newWorld(), func(c *ui.Context) {
		agent.Checkpoints(c, "S", []agent.Checkpoint{{Name: "A", At: at}, {Name: "A", At: at}}, agent.CheckpointsOptions{})
	})
}
