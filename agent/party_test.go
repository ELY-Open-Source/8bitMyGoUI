package agent_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/agent"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestAgentCardInksHPByContextLeft(t *testing.T) {
	for used, want := range map[int]pixel.Slot{20_000: pixel.Good, 120_000: pixel.Warn, 180_000: pixel.Bad} {
		a := agent.Agent{Name: "ROBO", Class: "CODER", Level: 7, Used: used, Window: 200_000}
		tt := view(300, 120, newWorld(), func(c *ui.Context) { agent.AgentCard(c, a, agent.AgentCardOptions{}) })
		x, y, _, _ := dots(find(t, tt, "ROBO"))
		// The fill starts two dots into the gauge's well.
		fill := count(tt.Image(), x+38, y+40, 92, 4, want)
		if left := 92 * (200_000 - used) / 200_000; fill < left*4-8 || fill > left*4+8 {
			t.Errorf("%d of 200K used: %d dots of %v, want about %d", used, fill, want, left*4)
		}
	}
}

func TestToolCallStripesItsOutcome(t *testing.T) {
	for out, ink := range map[agent.Outcome]pixel.Slot{agent.Casting: pixel.Info, agent.OK: pixel.Good, agent.Failed: pixel.Bad} {
		wd := newWorld()
		wd.reduce = true
		tt := view(500, 100, wd, func(c *ui.Context) {
			agent.ToolCall(c, "read_file", out, agent.ToolCallOptions{Args: "main.go"})
		})
		x, y, w, h := dots(find(t, tt, "CAST read_file"))
		if got := at(tt.Image(), x+2, y+5); got != ink || h != 38 {
			t.Errorf("outcome %d: stripe %v, want %v; %d dots high", out, got, ink, h)
		}
		// Casting is busy: it moves even under Reduce Motion.
		if moving := moves(tt, wd); moving != (out == agent.Casting) {
			t.Errorf("outcome %d: moving is %v", out, moving)
		}
		if out == agent.Failed && count(tt.Image(), x+w-60, y, 56, 22, pixel.Good) > 0 {
			t.Error("a failed cast shows Good")
		}
	}
}

func TestPlanQuestMarksEachObjective(t *testing.T) {
	steps := []agent.PlanStep{{Text: "Read the code", Done: true}, {Text: "Write a test"}, {Text: "Ship it"}}
	tt := view(400, 200, newWorld(), func(c *ui.Context) { agent.PlanQuest(c, "Fix the bug", steps, agent.PlanQuestOptions{}) })
	x, y, _, h := dots(find(t, tt, "Fix the bug"))
	row := func(i int, s pixel.Slot) int { return count(tt.Image(), x, y+20+i*16, 12, 16, s) }
	if row(0, pixel.Good) == 0 || row(1, pixel.Accent) == 0 || row(2, pixel.Well) == 0 || row(2, pixel.Accent) > 0 {
		t.Fatal("not a check, a ▶ and an open box")
	}
	for i := range steps {
		steps[i].Done = true
	}
	tt.Frame()
	_, _, _, done := dots(find(t, tt, "Fix the bug"))
	if done != h+18 || count(tt.Image(), x, y+h, 200, 18, pixel.Good) == 0 {
		t.Fatalf("all done: %d dots high, was %d, and no QUEST COMPLETE!", done, h)
	}
}

func TestThinkingOrbSwirlsOnlyWhileThinking(t *testing.T) {
	for _, thinking := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = true
		tt := view(100, 100, wd, func(c *ui.Context) { agent.ThinkingOrb(c, "Mind", thinking, agent.ThinkingOrbOptions{}) })
		x, y, _, _ := dots(find(t, tt, "Mind"))
		if moves(tt, wd) != thinking || (count(tt.Image(), x+3, y+1, 39, 39, pixel.Magic) > 0) != thinking {
			t.Errorf("thinking %v: the swirl does not match", thinking)
		}
	}
}

func TestTokenManaRunsDry(t *testing.T) {
	for used, dry := range map[int]bool{0: false, 99_999: false, 100_000: true, 250_000: true} {
		tt := view(300, 100, newWorld(), func(c *ui.Context) {
			agent.TokenMana(c, "Budget", used, 100_000, agent.TokenManaOptions{})
		})
		x, y, _, _ := dots(find(t, tt, "Budget"))
		img := tt.Image()
		if got := count(img, x+2, y+9, 116, 4, pixel.Bad) > 0; got != dry {
			t.Errorf("%d used: dry veil %v", used, got)
		}
		if full := count(img, x+2, y+9, 116, 4, pixel.Info); used == 0 && full < 116*4-9*4 {
			t.Errorf("a full budget fills %d dots", full)
		}
	}
}

func TestAgentStatusRings(t *testing.T) {
	for s, want := range map[agent.Status]struct {
		ink          pixel.Slot
		moves, still bool
	}{agent.Idle: {pixel.Faint, true, false}, agent.Working: {pixel.Info, true, true}, agent.Blocked: {pixel.Warn, false, false}, agent.Done: {pixel.Good, false, false}} {
		for _, reduce := range []bool{false, true} {
			wd := newWorld()
			wd.reduce = reduce
			tt := view(200, 200, wd, func(c *ui.Context) { agent.AgentStatus(c, "Robo", s, agent.AgentStatusOptions{}) })
			x, y, w, _ := dots(find(t, tt, "Robo"))
			if count(tt.Image(), x, y, w, 36, want.ink) == 0 {
				t.Errorf("status %d: no %v ring", s, want.ink)
			}
			// Working runs on; Idle's Z rests under Reduce Motion.
			if got, wantMove := moves(tt, wd), want.moves && (!reduce || want.still); got != wantMove {
				t.Errorf("status %d, reduce %v: moving %v", s, reduce, got)
			}
		}
	}
}

func TestSubAgentPartyDrawsTheSummoning(t *testing.T) {
	party := []agent.Member{{Name: "LEAD", Status: agent.Working, Summoned: []agent.Member{
		{Name: "SCOUT", Status: agent.Done}, {Name: "SMITH", Status: agent.Blocked, Summoned: []agent.Member{{Name: "IMP", Status: agent.Idle}}},
	}}}
	tt := view(400, 200, newWorld(), func(c *ui.Context) { agent.SubAgentParty(c, "Party", party, agent.SubAgentPartyOptions{}) })
	x, y, _, h := dots(find(t, tt, "Party"))
	img := tt.Image()
	if h != 4*16 {
		t.Fatalf("four members in %d dots", h)
	}
	// Robot inks: LEAD Info, SCOUT Good, SMITH Warn, IMP Faint.
	for i, s := range []pixel.Slot{pixel.Info, pixel.Good, pixel.Warn, pixel.Faint} {
		if count(img, x+[]int{0, 12, 12, 24}[i], y+i*16+2, 12, 12, s) == 0 {
			t.Errorf("row %d has no %v robot", i, s)
		}
	}
	// SCOUT's line runs on to SMITH; SMITH's ends at IMP's.
	if at(img, x+5, y+16+12) != pixel.Edge || at(img, x+17, y+48+12) == pixel.Edge {
		t.Error("the summoning lines run wrong")
	}
}

func TestDiffScrollSharesBlocks(t *testing.T) {
	files := []agent.FileDiff{{Path: "a.go", Added: 10}, {Path: "b.go", Removed: 10}, {Path: "c.go"}}
	tt := view(500, 200, newWorld(), func(c *ui.Context) { agent.DiffScroll(c, "Changes", files, agent.DiffScrollOptions{}) })
	x, y, w, _ := dots(find(t, tt, "Changes"))
	img := tt.Image()
	for i, want := range [][5]pixel.Slot{
		{pixel.Good, pixel.Good, pixel.Good, pixel.Good, pixel.Good}, {pixel.Bad, pixel.Bad, pixel.Bad, pixel.Bad, pixel.Bad},
		{pixel.Shade, pixel.Shade, pixel.Shade, pixel.Shade, pixel.Shade}, {pixel.Good, pixel.Good, pixel.Good, pixel.Bad, pixel.Bad},
	} {
		// The total row sits four dots lower, past the rule.
		ry := y + 10 + i*16 + 6
		if i == 3 {
			ry += 4
		}
		for j, s := range want {
			if got := at(img, x+w-32+j*5+1, ry+1); got != s {
				t.Errorf("row %d block %d is %v, want %v", i, j, got, s)
			}
		}
	}
}

func TestCostCoinsDropOnARise(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = reduce
		cents := 1234
		tt := view(200, 200, wd, func(c *ui.Context) { agent.CostCoins(c, "Spent", cents, agent.CostCoinsOptions{}) })
		x, y, _, _ := dots(find(t, tt, "Spent"))
		// Bronze cents stack third; a fifth coin lands at 26.
		air := func() int {
			img := tt.Image()
			return count(img, x+38, y+21, 10, 2, pixel.Warn) + count(img, x+38, y+21, 10, 2, pixel.Bad)
		}
		if air() > 0 {
			t.Fatal("a coin hangs before any change")
		}
		cents = 1235
		tt.Frame()
		if got := air() > 0; got == reduce {
			t.Errorf("reduce %v: the new coin falling is %v", reduce, got)
		}
		later(tt, wd, time.Second)
		if air() > 0 {
			t.Errorf("reduce %v: the coin never lands", reduce)
		}
	}
}
