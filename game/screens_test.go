package game_test

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/game"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var quests = []game.Quest{
	{Title: "Slime trouble", Objectives: []game.Objective{{Text: "Slay slimes", Have: 3, Need: 5}, {Text: "Tell the mayor", Have: 1, Need: 1}}},
	{Title: "Lost cat", Objectives: []game.Objective{{Text: "Find Mittens", Have: 1, Need: 1}}},
}

func TestQuestLogTracksTheChosenQuest(t *testing.T) {
	tracked, changes := 0, 0
	tt := view(600, 300, &world{clk: newClock()}, func(c *ui.Context) {
		if game.QuestLog(c, "Quests", quests, &tracked, game.QuestLogOptions{}).Changed {
			changes++
		}
	})
	shown := func(title string) int {
		return len(slices.DeleteFunc(tt.Texts(), func(s string) bool { return s != title }))
	}
	// The tracked title shows in the list and the tracker.
	if shown("Slime trouble") != 2 || shown("Lost cat") != 1 || !tt.HasText("TRACKING") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if err := tt.Click("Lost cat"); err != nil {
		t.Fatal(err)
	}
	if tracked != 1 || changes != 1 || shown("Lost cat") != 2 {
		t.Fatalf("clicking Lost cat tracked %d", tracked)
	}
	panics(t, "an objective past its need", func() {
		view(600, 300, &world{clk: newClock()}, func(c *ui.Context) {
			game.QuestTracker(c, game.Quest{Title: "x", Objectives: []game.Objective{{Text: "y", Have: 6, Need: 5}}}, game.QuestTrackerOptions{})
		})
	})
}

func TestQuestTrackerTicksMetObjectives(t *testing.T) {
	tt := view(400, 100, &world{clk: newClock()}, func(c *ui.Context) {
		game.QuestTracker(c, quests[0], game.QuestTrackerOptions{})
	})
	img := tt.Image()
	// Ticks are 7-dot boxes 5 dots into each 16-dot row.
	if count(t, img, 1, 16+5, 7, 7, pixel.Good) != 0 || count(t, img, 1, 32+5, 7, 7, pixel.Good) == 0 {
		t.Fatal("ticks do not follow the objectives")
	}
	keeps(t, img)
}

func TestPartyFramesGreyTheFallen(t *testing.T) {
	tt := view(400, 200, &world{clk: newClock()}, func(c *ui.Context) {
		game.PartyFrames(c, "Party", []game.Member{{Name: "Aria", Portrait: game.Hero, HP: 30, MaxHP: 40, MP: 12, MaxMP: 20},
			{Name: "Brom", Portrait: game.Hero, MaxHP: 50, MaxMP: 10}}, game.PartyFramesOptions{})
	})
	img := tt.Image()
	aria, brom := find(t, tt, "Aria"), find(t, tt, "Brom")
	ax, ay := dots(aria)
	bx, by := dots(brom)
	// Each portrait is a 28-dot box inside a 3-dot pad.
	if count(t, img, ax+3, ay+3, 28, 28, pixel.Bad) == 0 || count(t, img, bx+3, by+3, 28, 28, pixel.Bad) != 0 || count(t, img, bx+3, by+3, 28, 28, pixel.Shade) == 0 {
		t.Fatal("the fallen member is not greyed")
	}
	keeps(t, img)
}

func TestCharacterSheetShowsWhoAndStats(t *testing.T) {
	tt := view(400, 200, &world{clk: newClock(), zh: true}, func(c *ui.Context) {
		game.CharacterSheet(c, game.Character{Name: "Aria", Class: "Knight", Level: 7, Portrait: game.Hero, Stats: []game.Stat{{Name: "STR", Value: 14}}}, game.CharacterSheetOptions{})
	})
	if !tt.HasText("Aria") || !tt.HasText("Knight  等级 7") || !tt.HasText("STR 14") {
		t.Fatalf("texts %q", tt.Texts())
	}
	keeps(t, tt.Image())
}

func TestLevelUpScreenSpendsPoints(t *testing.T) {
	attrs := []game.Attribute{{Name: "STR", Base: 12}, {Name: "DEX", Base: 9}}
	points, done, changes := 2, false, 0
	tt := view(400, 200, &world{clk: newClock()}, func(c *ui.Context) {
		r := game.LevelUpScreen(c, "Level up", attrs, &points, game.LevelUpScreenOptions{})
		done = done || r.Done
		if r.Changed {
			changes++
		}
	})
	for _, label := range []string{"STR UP", "STR UP", "DEX UP", "DEX DOWN", "STR DOWN"} {
		tt.Click(label)
	}
	// With no points, DEX UP and DEX DOWN are off.
	if attrs[0].Added != 1 || attrs[1].Added != 0 || points != 1 || changes != 3 {
		t.Fatalf("STR +%d, DEX +%d, %d points, %d changes", attrs[0].Added, attrs[1].Added, points, changes)
	}
	tt.Click("DONE")
	if !done {
		t.Fatal("DONE did not report")
	}
}

func TestShopBuysAndSellsWithinGold(t *testing.T) {
	wares := []game.Ware{{Item: data.Item{Sprite: icon.Potion, Name: "Potion"}, Price: 30}, {Item: data.Item{Sprite: icon.Shield, Name: "Shield"}, Price: 400}}
	gold, bought, sold := 100, -1, -1
	tt := view(500, 200, &world{clk: newClock()}, func(c *ui.Context) {
		r := game.Shop(c, "Shop", wares, &gold, game.ShopOptions{})
		if r.Bought >= 0 {
			bought = r.Bought
		}
		if r.Sold >= 0 {
			sold = r.Sold
		}
	})
	tt.Click("BUY Potion")
	tt.Click("BUY Potion")
	tt.Click("SELL Potion")
	if gold != 55 || wares[0].Owned != 1 || bought != 0 || sold != 0 {
		t.Fatalf("gold %d, owned %d", gold, wares[0].Owned)
	}
	tt.Click("BUY Shield")
	tt.Click("SELL Shield")
	if gold != 55 || wares[1].Owned != 0 {
		t.Fatal("a ware too dear, or none owned, still traded")
	}
	sx, sy := dots(find(t, tt, "Shield"))
	if count(t, tt.Image(), sx, sy, 30, 16, pixel.Faint) == 0 || count(t, tt.Image(), sx, sy, 30, 16, pixel.Ink) != 0 {
		t.Fatal("a dear ware's name is not Faint")
	}
	keeps(t, tt.Image())
}
