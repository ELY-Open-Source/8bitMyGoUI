package game_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/game"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestLetterboxSlidesInSteps(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		on := false
		tt := view(300, 200, wd, func(c *ui.Context) {
			game.Letterbox(c, on, game.LetterboxOptions{}, func() { ui.Box(c).Size(pixel.K(c).D(100), pixel.K(c).D(60)) })
		})
		bar := func() int { return count(t, tt.Image(), 50, 0, 1, 30, pixel.Ink) }
		on = true
		tt.Frame()
		// Two of four steps in, a bar is 8 dots.
		wd.clk.t = wd.clk.t.Add(250 * time.Millisecond)
		tt.Frame()
		if half := bar(); (half == 8) == reduce || reduce && half != 16 {
			t.Fatalf("reduced %v: a quarter second in, the bar is %d dots", reduce, half)
		}
		wd.clk.t = wd.clk.t.Add(time.Second)
		tt.Frame()
		if bar() != 16 {
			t.Fatalf("reduced %v: the bar settled at %d", reduce, bar())
		}
		// Turned off mid-slide, the bars return from there.
		on = false
		tt.Frame()
		wd.clk.t = wd.clk.t.Add(250 * time.Millisecond)
		tt.Frame()
		on = true
		tt.Frame()
		if now := bar(); !reduce && now != 8 || reduce && now != 16 {
			t.Fatalf("reduced %v: turned back on at half, the bar is %d", reduce, now)
		}
	}
}

func TestMissionBriefStarts(t *testing.T) {
	started := false
	tt := view(500, 200, &world{clk: newClock()}, func(c *ui.Context) {
		started = started || game.MissionBrief(c, game.Quest{Title: "Rescue", Objectives: []game.Objective{{Text: "Go", Need: 1}}}, []string{"####", "#..#", "####"}, game.Spot{X: 1, Y: 1}, nil, game.MissionBriefOptions{}).Start
	})
	tt.Click("START")
	if !started {
		t.Fatal("START did not report")
	}
}

func TestAchievementsHideLockedNames(t *testing.T) {
	awards := []game.Award{{Sprite: icon.Trophy, Name: "First blood", How: "Win a fight", Unlocked: true}, {Sprite: icon.Coin, Name: "King", How: "Beat the king"}}
	sel := 0
	tt := view(400, 200, &world{clk: newClock()}, func(c *ui.Context) {
		game.Achievements(c, "Awards", awards, &sel, game.AchievementsOptions{})
	})
	if !tt.HasText("UNLOCKED 1/2") || !tt.HasText("First blood") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Click("Awards")
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyRight)
	if sel != 1 || tt.HasText("King") || !tt.HasText("???") || !tt.HasText("Beat the king") {
		t.Fatalf("a locked award showed %q", tt.Texts())
	}
	r := find(t, tt, "Awards")
	x, y := dots(r)
	// The locked coin is a Shade shadow: no Warn gold.
	if count(t, tt.Image(), x+24, y, 24, 24, pixel.Warn) != 0 || count(t, tt.Image(), x+24, y, 24, 24, pixel.Shade) == 0 {
		t.Fatal("the locked award is not a silhouette")
	}
}

func TestResultScreenCountsThenRanks(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		tt := view(400, 200, wd, func(c *ui.Context) {
			game.ResultScreen(c, "Result", 83*time.Second, 12500, "S", game.ResultScreenOptions{})
		})
		rank := func() int { return count(t, tt.Image(), 0, 0, 200, 100, pixel.Warn) }
		if (rank() > 0) != reduce {
			t.Fatalf("reduced %v: the rank showed before the counts", reduce)
		}
		wd.clk.t = wd.clk.t.Add(2 * time.Second)
		tt.Frame()
		if rank() == 0 {
			t.Fatalf("reduced %v: no rank after the counts", reduce)
		}
	}
}

func TestResultScreenFitsABigScore(t *testing.T) {
	tt := view(600, 200, &world{clk: newClock()}, func(c *ui.Context) {
		ui.Column(c).Padding(pixel.K(c).D(20)).Children(func() {
			game.ResultScreen(c, "Result", 83*time.Second, 1000000000000000, "S", game.ResultScreenOptions{})
		})
	})
	if n := count(t, tt.Image(), 0, 0, 20, 100, pixel.Ink); n != 0 {
		t.Fatalf("%d score dots left of the screen", n)
	}
	if r := find(t, tt, "Result"); int(r.W)/2 < pixel.Measure("SCORE")+12+pixel.Measure("1000000000000000") {
		t.Fatalf("the screen is %v px wide", r.W)
	}
}

func TestCraftingShowsAKnownRecipe(t *testing.T) {
	grid := make([]data.Item, 9)
	grid[4] = data.Item{Sprite: icon.Leaf, Name: "Herb"}
	recipes := []game.Recipe{{Grid: [9]string{4: "Herb"}, Result: data.Item{Sprite: icon.Potion, Name: "Potion"}}}
	cell, crafted, picked := 0, 0, -1
	tt := view(500, 200, &world{clk: newClock()}, func(c *ui.Context) {
		r := game.Crafting(c, "Craft", grid, recipes, &cell, game.CraftingOptions{})
		if r.Crafted {
			crafted++
		}
		if r.Picked >= 0 {
			picked = r.Picked
		}
	})
	if !tt.HasText("Potion") {
		t.Fatalf("no result shown: %q", tt.Texts())
	}
	tt.Click("CRAFT")
	grid[4] = data.Item{}
	tt.Frame()
	tt.Click("CRAFT")
	if crafted != 1 || tt.HasText("Potion") {
		t.Fatalf("crafted %d times; texts %q", crafted, tt.Texts())
	}
	r := find(t, tt, "Craft")
	tt.ClickAt(r.X+float32(2*(22*2+11)), r.Y+float32(2*(22*2+11)))
	if cell != 8 || picked != 8 {
		t.Fatalf("clicking the last cell picked %d", picked)
	}
}

func TestLootPopupRisesAndCloses(t *testing.T) {
	wd := &world{clk: newClock()}
	open, closed := true, 0
	tt := view(500, 300, wd, func(c *ui.Context) {
		if game.LootPopup(c, data.Item{Sprite: icon.Potion, Name: "POTION", Count: 3}, &open, game.LootPopupOptions{}).Closed {
			closed++
		}
	})
	if !tt.HasText("YOU GOT: POTION ×3!") {
		t.Fatalf("texts %q", tt.Texts())
	}
	r := find(t, tt, "POTION ×3")
	x, y := dots(r)
	top := func() int {
		for j := range int(r.H) / 2 {
			if count(t, tt.Image(), x, y+j, int(r.W)/2, 1, pixel.Bad) > 0 {
				return j
			}
		}
		return -1
	}
	low := top()
	wd.clk.t = wd.clk.t.Add(time.Second)
	tt.Frame()
	if rose := low - top(); rose != 8 {
		t.Fatalf("the potion rose %d dots", rose)
	}
	tt.Click("OK")
	if open || closed != 1 {
		t.Fatal("OK did not close it")
	}
}

func TestTreasureChestOpensOnce(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		open, opened := false, 0
		tt := view(200, 200, wd, func(c *ui.Context) {
			if game.TreasureChest(c, "Chest", data.Item{Sprite: icon.Gem, Name: "Gem"}, &open, game.TreasureChestOptions{}).Opened {
				opened++
			}
		})
		gem := func() int { return count(t, tt.Image(), 0, 0, 24, 21, pixel.Info) }
		tt.Click("Chest")
		if !open || opened != 1 || (gem() > 0) != reduce {
			t.Fatalf("reduced %v: the press showed the gem %v", reduce, gem() > 0)
		}
		wd.clk.t = wd.clk.t.Add(time.Second)
		tt.Frame()
		tt.Click("Chest")
		if gem() == 0 || opened != 1 {
			t.Fatalf("reduced %v: no gem, or opened twice", reduce)
		}
	}
}

func TestTreasureChestGemStaysInside(t *testing.T) {
	open := true
	tt := view(200, 200, &world{clk: newClock()}, func(c *ui.Context) {
		ui.Column(c).Padding(pixel.K(c).D(20)).Children(func() {
			game.TreasureChest(c, "Chest", data.Item{Sprite: icon.Gem, Name: "Gem"}, &open, game.TreasureChestOptions{})
		})
	})
	if n, in := count(t, tt.Image(), 0, 0, 100, 20, pixel.Info), count(t, tt.Image(), 20, 20, 24, 50, pixel.Info); n != 0 || in == 0 {
		t.Fatalf("%d gem dots above the chest, %d in it", n, in)
	}
}

func TestGachaReelSpinsAndStops(t *testing.T) {
	faces := []data.Item{{Sprite: icon.Star, Name: "Star"}, {Sprite: icon.Heart, Name: "Heart"}, {Sprite: icon.Coin, Name: "Coin"}}
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		spin, done := 0, 0
		tt := view(300, 100, wd, func(c *ui.Context) {
			if game.GachaReel(c, "Reels", faces, [3]int{2, 2, 2}, spin, game.GachaReelOptions{}).Done {
				done++
			}
		})
		coins := func() int { return count(t, tt.Image(), 0, 0, 84, 28, pixel.Warn) }
		still := coins()
		spin = 1
		tt.Frame()
		var seen []int
		for range 4 {
			wd.clk.t = wd.clk.t.Add(100 * time.Millisecond)
			tt.Frame()
			seen = append(seen, coins())
		}
		spun := seen[0] != still || seen[1] != still || seen[2] != still
		if spun == reduce || (done == 1) != reduce {
			t.Fatalf("reduced %v: spun %v, done %d", reduce, spun, done)
		}
		wd.clk.t = wd.clk.t.Add(2 * time.Second)
		tt.Frame()
		if done != 1 || coins() != still {
			t.Fatalf("reduced %v: done %d, coins %d of %d", reduce, done, coins(), still)
		}
		// Spinning back to an old cue lands; reduced, both do.
		spin = 2
		tt.Frame()
		spin = 1
		tt.Frame()
		wd.clk.t = wd.clk.t.Add(2 * time.Second)
		tt.Frame()
		if want := map[bool]int{false: 2, true: 3}[reduce]; done != want {
			t.Fatalf("reduced %v: spinning back to 1 landed %d times in all", reduce, done)
		}
	}
}

func TestGachaReelRefusesAFaceWithoutSprite(t *testing.T) {
	panics(t, "a face with no sprite", func() {
		view(300, 100, &world{clk: newClock()}, func(c *ui.Context) {
			game.GachaReel(c, "x", []data.Item{{Sprite: icon.Star, Name: "Star"}, {Name: "Blank"}}, [3]int{}, 0, game.GachaReelOptions{})
		})
	})
}

func TestInteractPromptSitsAboveItsTarget(t *testing.T) {
	var target *ui.Element
	pad := 40
	tt := view(400, 300, &world{clk: newClock(), reduce: true}, func(c *ui.Context) {
		ui.Column(c).Padding(pixel.K(c).D(pad)).Children(func() {
			target = ui.Box(c).Size(pixel.K(c).D(20), pixel.K(c).D(20)).Label("Mayor")
		})
		game.InteractPrompt(c, target, "TALK", game.InteractPromptOptions{})
	})
	// Right from the start, and after a move.
	for range 2 {
		bubble, mayor := find(t, tt, "A TALK"), find(t, tt, "Mayor")
		if bubble.Y+bubble.H > mayor.Y || bubble.Y+bubble.H < mayor.Y-20 || bubble.X > mayor.X+mayor.W/2 {
			t.Fatalf("the bubble %v does not sit above %v", bubble, mayor)
		}
		pad = 70
		tt.Frame()
	}
	panics(t, "a two-letter key", func() {
		view(300, 200, &world{clk: newClock()}, func(c *ui.Context) {
			game.InteractPrompt(c, ui.Box(c), "TALK", game.InteractPromptOptions{Key: "AB"})
		})
	})
}

func TestTurnOrderRaisesTheActive(t *testing.T) {
	tt := view(200, 60, &world{clk: newClock()}, func(c *ui.Context) {
		game.TurnOrder(c, "Turns", []game.Fighter{{Name: "Aria", Portrait: game.Hero}, {Name: "Slime", Portrait: icon.Skull, Foe: true}}, 1, game.TurnOrderOptions{})
	})
	img := tt.Image()
	if at(t, img, 19+8, 0) != pixel.Accent || at(t, img, 8, 0) == pixel.Accent || at(t, img, 8, 5) != pixel.Good {
		t.Fatal("the active portrait is not raised in Accent")
	}
}

func TestBattleLogKeepsTheNewest(t *testing.T) {
	lines := []game.Line{{Text: "one", Slot: pixel.Good}, {Text: "two"}, {Text: "three"}, {Text: "four"}, {Text: "five", Slot: pixel.Bad}}
	tt := view(400, 200, &world{clk: newClock()}, func(c *ui.Context) { game.BattleLog(c, "Log", lines, game.BattleLogOptions{}) })
	// 68 dots hold four lines: the first rolled off.
	if count(t, tt.Image(), 0, 0, 180, 68, pixel.Good) != 0 || count(t, tt.Image(), 0, 50, 180, 18, pixel.Bad) == 0 {
		t.Fatal("the log does not keep the newest lines at the bottom")
	}
	// Two lines sit at the bottom, not the top.
	lines = lines[3:]
	lines[0].Slot = pixel.Good
	tt.Frame()
	if count(t, tt.Image(), 0, 34, 180, 16, pixel.Good) == 0 || count(t, tt.Image(), 0, 0, 180, 30, pixel.Good) != 0 {
		t.Fatal("a short log is not at the bottom")
	}
}

func TestTargetSelectorMovesTheCursor(t *testing.T) {
	sel, chosen := 0, -1
	enemies := []game.Enemy{{Name: "Slime", Sprite: icon.Skull, HP: 3, MaxHP: 9}, {Name: "Bat", Sprite: icon.Bug, HP: 6, MaxHP: 6}}
	tt := view(300, 200, &world{clk: newClock()}, func(c *ui.Context) {
		if r := game.TargetSelector(c, "Targets", enemies, &sel, game.TargetSelectorOptions{}); r.Chosen >= 0 {
			chosen = r.Chosen
		}
	})
	cursor := func(i int) int { return count(t, tt.Image(), i*40, 0, 40, 12, pixel.Accent) }
	if cursor(0) == 0 || cursor(1) != 0 {
		t.Fatal("the ▼ is not over the first target")
	}
	tt.Click("Targets")
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if sel != 1 || chosen != 1 || cursor(1) == 0 {
		t.Fatalf("Right, Enter chose %d", chosen)
	}
	// A long name stays in its 40-dot cell.
	enemies[0].Name = "Slime of the Very Long Name"
	tt.Frame()
	if n := count(t, tt.Image(), 80, 36, 60, 20, pixel.Ink); n != 0 {
		t.Fatalf("%d name dots spilled past the cells", n)
	}
}

func TestATBGaugeFlashesReady(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		v := 0.5
		tt := view(200, 60, wd, func(c *ui.Context) { game.ATBGauge(c, "ARIA", v, game.ATBGaugeOptions{}) })
		if count(t, tt.Image(), 0, 9, 64, 1, pixel.Info) != 30 {
			t.Fatal("half a gauge is not 30 dots")
		}
		v = 1
		tt.Frame()
		var ready []int
		for range 2 {
			ready = append(ready, count(t, tt.Image(), 30, 0, 34, 6, pixel.Warn))
			wd.clk.t = wd.clk.t.Add(125 * time.Millisecond)
			tt.Frame()
		}
		if flashed := (ready[0] == 0) != (ready[1] == 0); flashed == reduce || ready[0]+ready[1] == 0 {
			t.Fatalf("reduced %v: READY dots %v", reduce, ready)
		}
	}
	zh := view(200, 60, &world{clk: newClock(), zh: true}, func(c *ui.Context) { game.ATBGauge(c, "ARIA", 1, game.ATBGaugeOptions{}) })
	// 就绪 needs a 14-row label line above the 8-row bar.
	if count(t, zh.Image(), 20, 0, 44, 14, pixel.Warn) == 0 || find(t, zh, "ARIA").H != 2*(14+8) {
		t.Fatal("a Chinese READY is missing or crowds the bar")
	}
	panics(t, "a gauge too narrow", func() {
		view(200, 60, &world{clk: newClock()}, func(c *ui.Context) { game.ATBGauge(c, "ARIA", 1, game.ATBGaugeOptions{Width: 10}) })
	})
}

func TestStaminaBarDithersThePartBlock(t *testing.T) {
	tt := view(200, 60, &world{clk: newClock()}, func(c *ui.Context) { game.StaminaBar(c, "STAMINA", 0.65, game.StaminaBarOptions{}) })
	img := tt.Image()
	// Six blocks full, the seventh half, the rest empty.
	if full, part, empty := count(t, img, 0, 7, 36, 6, pixel.Good), count(t, img, 36, 7, 5, 6, pixel.Good), count(t, img, 42, 7, 18, 6, pixel.Good); full != 180 || part == 0 || part == 30 || empty != 0 {
		t.Fatalf("blocks: full %d, part %d, empty %d", full, part, empty)
	}
}
