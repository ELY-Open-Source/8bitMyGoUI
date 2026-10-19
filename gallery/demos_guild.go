package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/agent"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var guildEntries = []Entry{
	{459, "MemoryCrystals", "记忆水晶", "agent", "Memory slots as crystals. Lit ones hold something; the arrows peek inside.", demoMemoryCrystals},
	{460, "AgentStatus", "状态环", "agent", "A ring that says it all: sleepy, busy, stuck or done.", demoAgentStatus},
	{461, "DiffScroll", "改动卷轴", "agent", "What changed, on a scroll: each file with lines added and removed.", demoDiffScroll},
	{462, "ModelSelect", "模型选择", "agent", "Pick a model like a fighter: stats on the right, a flag on your pick.", demoModelSelect},
	{463, "PromptBox", "咒语输入", "agent", "Write the spell, then CAST. Enter casts; Shift+Enter starts a new line.", demoPromptBox},
	{464, "CostCoins", "花费", "agent", "What it cost, in coins: gold dollars, silver dimes, bronze cents.", demoCostCoins},
	{465, "Checkpoints", "存档点", "agent", "Save points along the way. RESTORE takes you back.", demoCheckpoints},
}

type guildDemo struct {
	model, cents int
	prompt       string
	saves        []agent.Checkpoint
}

func guildState(c *ui.Context) *guildDemo {
	return ui.Local(c.Root(), demoKey("guild"), func() guildDemo {
		now := pixel.K(c).Now()
		return guildDemo{model: 1, cents: 1234, saves: []agent.Checkpoint{
			{Name: "Project start", At: now.Add(-50 * time.Hour)},
			{Name: "Before the refactor", At: now.Add(-3 * time.Hour)},
			{Name: "Tests green", At: now.Add(-5 * time.Minute)},
		}}
	})
}

func demoMemoryCrystals(c *ui.Context, log func(string)) {
	memories := []string{"Prefers tabs", "Tests before commits", "Never force-push", "Coffee: black", "Hates dark patterns"}
	if i := agent.MemoryCrystals(c, "Memories", memories, agent.MemoryCrystalsOptions{}).Picked; i >= 0 {
		if i < len(memories) {
			log("recalled: " + memories[i])
			return
		}
		log(fmt.Sprintf("crystal %d is empty", i+1))
	}
}

func demoAgentStatus(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(8)).Children(func() {
		for _, s := range []agent.Status{agent.Idle, agent.Working, agent.Blocked, agent.Done} {
			agent.AgentStatus(c, fmt.Sprintf("Robo %d", s+1), s, agent.AgentStatusOptions{})
		}
	})
}

func demoDiffScroll(c *ui.Context, log func(string)) {
	agent.DiffScroll(c, "Changes", []agent.FileDiff{
		{Path: "pixel/canvas.go", Added: 42, Removed: 7},
		{Path: "agent/agentstatus.go", Added: 88},
		{Path: "README.md", Added: 3, Removed: 1},
		{Path: "old/clever_hack.go", Removed: 120},
	}, agent.DiffScrollOptions{Width: 260})
}

func demoModelSelect(c *ui.Context, log func(string)) {
	s := guildState(c)
	models := []agent.Model{
		{Name: "SPRITE", Portrait: icon.Bolt, Speed: 5, Power: 2, Cost: 1},
		{Name: "KNIGHT", Portrait: icon.Shield, Speed: 3, Power: 4, Cost: 3},
		{Name: "WIZARD", Portrait: icon.Sparkle, Speed: 2, Power: 5, Cost: 5},
		{Name: "BARD", Portrait: icon.Music, Speed: 4, Power: 3, Cost: 2},
	}
	if agent.ModelSelect(c, "Models", models, &s.model, agent.ModelSelectOptions{}).Changed {
		log(models[s.model].Name + " joins the party")
	}
}

func demoPromptBox(c *ui.Context, log func(string)) {
	s := guildState(c)
	r := agent.PromptBox(c, "Prompt", &s.prompt, agent.PromptBoxOptions{Placeholder: "Describe the quest…"})
	// Prompts can hold secrets, so the log only counts.
	if r.Changed {
		log(fmt.Sprintf("%d runes written", utf8.RuneCountInString(s.prompt)))
	}
	if r.Cast != "" {
		log(fmt.Sprintf("cast %d runes", utf8.RuneCountInString(r.Cast)))
	}
}

func demoCostCoins(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := guildState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.End).Children(func() {
		agent.CostCoins(c, "Spent", s.cents, agent.CostCoinsOptions{})
		ui.Column(c).Gap(k.D(4)).Children(func() {
			if input.Button(c, "RUN THE AGENT", input.ButtonOptions{Size: input.Small}).Clicked {
				s.cents += 37
				log(fmt.Sprintf("spent $%d.%02d", s.cents/100, s.cents%100))
			}
			if input.Button(c, "NEW MONTH", input.ButtonOptions{Size: input.Small}).Clicked {
				s.cents = 0
				log("the bill is paid")
			}
		})
	})
}

func demoCheckpoints(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := guildState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		if i := agent.Checkpoints(c, "Save points", s.saves, agent.CheckpointsOptions{}).Restored; i >= 0 {
			log("restored: " + s.saves[i].Name)
		}
		if input.Button(c, "SAVE HERE", input.ButtonOptions{Size: input.Small}).Clicked {
			s.saves = append(s.saves, agent.Checkpoint{Name: fmt.Sprintf("Save %d", len(s.saves)+1), At: k.Now()})
			log("saved")
		}
	})
}
