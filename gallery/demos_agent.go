package gallery

import (
	"fmt"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/agent"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var agentEntries = []Entry{
	{451, "AgentCard", "智能体卡", "agent", "A party card. The context window is the HP bar: the more it reads, the less it has left.", demoAgentCard},
	{452, "ToolCall", "工具调用", "agent", "Every tool call is a spell: CAST, a moment of glowing, then OK or FAILED.", demoToolCall},
	{453, "PlanQuest", "计划任务", "agent", "The plan as quest objectives. Finish them all for a fanfare.", demoPlanQuest},
	{454, "ThinkingOrb", "思考水晶球", "agent", "A crystal ball that swirls while the agent thinks.", demoThinkingOrb},
	{455, "TokenMana", "令牌魔力", "agent", "The token budget as MP. Spend it all and the well runs dry.", demoTokenMana},
	{456, "ApprovalPrompt", "授权提示", "agent", "ALLOW SPELL? The agent asks before it runs a command; risky ones come in red.", demoApprovalPrompt},
	{457, "AgentLog", "智能体日志", "agent", "The transcript as a battle log. The newest line types itself out.", demoAgentLog},
	{458, "SubAgentParty", "子智能体", "agent", "Agents summon helpers. The tree shows who summoned whom, and how each is doing.", demoSubAgentParty},
}

type agentDemo struct {
	used, done, mana, next int
	cast                   time.Time
	told, thinking         bool
	log                    []agent.LogStep
	party                  []agent.Member
}

func agentState(c *ui.Context) *agentDemo {
	return ui.Local(c.Root(), demoKey("agent"), func() agentDemo {
		return agentDemo{used: 64_000, done: 1, mana: 28_000, told: true, next: 3, log: slices.Clone(story[:3]),
			party: []agent.Member{{Name: "LEAD", Status: agent.Working, Summoned: []agent.Member{
				{Name: "SCOUT", Status: agent.Done}, {Name: "SMITH", Status: agent.Working, Summoned: []agent.Member{{Name: "IMP", Status: agent.Blocked}}},
			}}}}
	})
}

func demoAgentCard(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := agentState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		agent.AgentCard(c, agent.Agent{Name: "ROBO", Class: "CODER", Level: 12, Used: s.used, Window: 200_000}, agent.AgentCardOptions{})
		agent.AgentCard(c, agent.Agent{Name: "SAGE", Class: "REVIEWER", Level: 30, Used: 182_000, Window: 200_000}, agent.AgentCardOptions{})
		ui.Row(c).Gap(k.D(6)).Children(func() {
			if input.Button(c, "READ A BIG FILE", input.ButtonOptions{Size: input.Small}).Clicked {
				s.used = min(200_000, s.used+24_000)
				log(fmt.Sprintf("ROBO used %d tokens", s.used))
			}
			// Compacting is a night at the inn: HP comes back.
			if input.Button(c, "COMPACT", input.ButtonOptions{Size: input.Small}).Clicked {
				s.used = 12_000
				log("ROBO compacted its context")
			}
		})
	})
}

func demoToolCall(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := agentState(c)
	grep := agent.OK
	if done := s.cast.Add(1500 * time.Millisecond); k.Now().Before(done) {
		grep = agent.Casting
		k.At(c, done)
	} else if !s.told {
		s.told = true
		log("grep: OK")
	}
	ui.Column(c).Gap(k.D(4)).Children(func() {
		agent.ToolCall(c, "read_file", agent.OK, agent.ToolCallOptions{Args: "pixel/canvas.go"})
		agent.ToolCall(c, "grep", grep, agent.ToolCallOptions{Args: `"TODO" in 512 files`})
		agent.ToolCall(c, "bash", agent.Failed, agent.ToolCallOptions{Args: "npm test: 2 tests fainted"})
		if input.Button(c, "CAST GREP", input.ButtonOptions{Size: input.Small, Disabled: grep == agent.Casting}).Clicked {
			s.cast, s.told = k.Now(), false
			log("casting grep")
		}
	})
}

func demoPlanQuest(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := agentState(c)
	steps := []agent.PlanStep{{Text: "Reproduce the save bug"}, {Text: "Write a failing test"}, {Text: "Fix the root cause"}, {Text: "Run every test"}, {Text: "Ship it"}}
	for i := range s.done {
		steps[i].Done = true
	}
	ui.Column(c).Gap(k.D(6)).Children(func() {
		agent.PlanQuest(c, "The Lost Save", steps, agent.PlanQuestOptions{})
		if s.done < len(steps) {
			if input.Button(c, "COMPLETE STEP", input.ButtonOptions{Size: input.Small}).Clicked {
				s.done++
				log(fmt.Sprintf("objective %d done", s.done))
			}
			return
		}
		if input.Button(c, "NEW QUEST", input.ButtonOptions{Size: input.Small}).Clicked {
			s.done = 0
			log("a new quest begins")
		}
	})
}

func demoThinkingOrb(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := agentState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		agent.ThinkingOrb(c, "Mind", s.thinking, agent.ThinkingOrbOptions{})
		label := "PONDER"
		if s.thinking {
			label = "SNAP OUT OF IT"
		}
		if input.Button(c, label, input.ButtonOptions{Size: input.Small}).Clicked {
			s.thinking = !s.thinking
			log(fmt.Sprintf("thinking %v", s.thinking))
		}
	})
}

func demoTokenMana(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := agentState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		agent.TokenMana(c, "Token budget", s.mana, 100_000, agent.TokenManaOptions{Width: 160})
		ui.Row(c).Gap(k.D(6)).Children(func() {
			if input.Button(c, "CAST (12K)", input.ButtonOptions{Size: input.Small}).Clicked {
				s.mana += 12_000
				log(fmt.Sprintf("%d tokens spent", s.mana))
			}
			if input.Button(c, "DRINK ETHER", input.ButtonOptions{Size: input.Small}).Clicked {
				s.mana = 0
				log("MP restored")
			}
		})
	})
}

func demoApprovalPrompt(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(8)).Children(func() {
		for _, cmd := range []struct {
			line   string
			danger bool
		}{{"go test ./...", false}, {"rm -rf ~/saves", true}} {
			r := agent.ApprovalPrompt(c, cmd.line, cmd.line, agent.ApprovalPromptOptions{Danger: cmd.danger})
			if r.Allowed {
				log("allowed: " + cmd.line)
			}
			if r.Denied {
				log("denied: " + cmd.line)
			}
		}
	})
}

// story is what the log demo's agent does, in turn.
var story = []agent.LogStep{
	{Kind: agent.Thought, Text: "the save file looks corrupt"}, {Kind: agent.Cast, Text: "read_file"},
	{Kind: agent.Hit, Text: "a stray byte at 0x2A"}, {Kind: agent.Cast, Text: "edit_file"},
	{Kind: agent.Miss, Text: "one test still fails"}, {Kind: agent.Cast, Text: "edit_file"},
	{Kind: agent.Hit, Text: "all 64 tests pass"}, {Kind: agent.Said, Text: "Fixed! Shall I ship it?"},
}

func demoAgentLog(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := agentState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		agent.AgentLog(c, "Transcript", "ROBO", s.log, agent.AgentLogOptions{Width: 260})
		if input.Button(c, "NEXT STEP", input.ButtonOptions{Size: input.Small}).Clicked {
			step := story[s.next%len(story)]
			s.log, s.next = append(s.log, step), s.next+1
			log(fmt.Sprintf("step %d: %s", s.next, step.Text))
		}
	})
}

func demoSubAgentParty(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := agentState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		agent.SubAgentParty(c, "Party", s.party, agent.SubAgentPartyOptions{})
		ui.Row(c).Gap(k.D(6)).Children(func() {
			lead := &s.party[0]
			if input.Button(c, "SUMMON", input.ButtonOptions{Size: input.Small, Disabled: len(lead.Summoned) == 4}).Clicked {
				name := []string{"GOLEM", "WISP"}[len(lead.Summoned)-2]
				lead.Summoned = append(lead.Summoned, agent.Member{Name: name, Status: agent.Working})
				log("LEAD summons " + name)
			}
			// A turn passes: the blocked wake, the busy finish.
			if input.Button(c, "NEXT TURN", input.ButtonOptions{Size: input.Small}).Clicked {
				advance(lead.Summoned)
				log("a turn passes")
			}
		})
	})
}

// advance moves each summoned member on by a turn.
func advance(ms []agent.Member) {
	for i := range ms {
		switch ms[i].Status {
		case agent.Blocked, agent.Idle:
			ms[i].Status = agent.Working
		case agent.Working:
			ms[i].Status = agent.Done
		}
		advance(ms[i].Summoned)
	}
}
