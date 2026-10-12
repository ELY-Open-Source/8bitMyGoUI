package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/game"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var lootEntries = []Entry{
	{280, "ResultScreen", "结算画面", "game", "Time and score count up; then the rank lands.", demoResultScreen},
	{281, "Crafting", "合成台", "game", "Click a cell to swap ingredients. Herb over water makes a potion.", demoCrafting},
	{282, "LootPopup", "获得物品", "game", "YOU GOT: something nice! The item rises out of the box.", demoLootPopup},
	{283, "TreasureChest", "宝箱", "game", "Press the chest. It creaks open; the gem floats up.", demoTreasureChest},
	{284, "GachaReel", "扭蛋转轮", "game", "SPIN: three reels, stopping one by one. No real money involved.", demoGachaReel},
	{285, "InteractPrompt", "互动提示", "game", "Walk up to the mayor and a TALK bubble bobs above him.", demoInteractPrompt},
	{286, "TurnOrder", "行动顺序", "game", "Who acts next. The active one stands raised.", demoTurnOrder},
	{287, "BattleLog", "战斗日志", "game", "Combat lines roll up as the fight goes on.", demoBattleLog},
	{288, "TargetSelector", "目标选择", "game", "A ▼ over the target. Arrows choose; Enter attacks.", demoTargetSelector},
	{289, "ATBGauge", "行动槽", "game", "Each gauge fills at its hero's speed and flashes READY.", demoATBGauge},
	{290, "StaminaBar", "耐力条", "game", "Sprint and it drains; rest and it refills.", demoStaminaBar},
}

type lootDemo struct {
	result, cell, spin, turn, target int
	grid                             []data.Item
	stops                            [3]int
	popup, chest, near, sprint       bool
	lines                            []game.Line
	ready                            [3]time.Time
	stamina                          float64
	last                             time.Time
}

func lootState(c *ui.Context) *lootDemo {
	now := pixel.K(c).Now()
	return ui.Local(c.Root(), demoKey("loot"), func() lootDemo {
		s := lootDemo{grid: make([]data.Item, 9), near: true, stamina: 1, last: now, ready: [3]time.Time{now, now, now}}
		s.grid[1], s.grid[4] = herb, water
		return s
	})
}

var (
	herb  = data.Item{Sprite: icon.Leaf, Name: "Herb"}
	water = data.Item{Sprite: icon.Drop, Name: "Water"}
)

func demoResultScreen(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := lootState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		// A new key makes a fresh screen; counts run again.
		ui.Box(c).Key(s.result).Children(func() {
			game.ResultScreen(c, "Stage clear", 83*time.Second+450*time.Millisecond, 12500, "S", game.ResultScreenOptions{})
		})
		if small(c, "AGAIN", false) {
			s.result++
			log("counting again")
		}
	})
}

func demoCrafting(c *ui.Context, log func(string)) {
	s := lootState(c)
	recipes := []game.Recipe{{Grid: [9]string{1: "Herb", 4: "Water"}, Result: data.Item{Sprite: icon.Potion, Name: "Potion"}}}
	r := game.Crafting(c, "Workbench", s.grid, recipes, &s.cell, game.CraftingOptions{})
	// A picked cell cycles: empty, herb, water.
	if r.Picked >= 0 {
		switch s.grid[r.Picked].Name {
		case "":
			s.grid[r.Picked] = herb
		case "Herb":
			s.grid[r.Picked] = water
		default:
			s.grid[r.Picked] = data.Item{}
		}
		log(fmt.Sprintf("cell %d: %s", r.Picked+1, s.grid[r.Picked].Name))
	}
	if r.Crafted {
		s.grid = make([]data.Item, 9)
		log("made a potion!")
	}
}

func demoLootPopup(c *ui.Context, log func(string)) {
	s := lootState(c)
	if small(c, "OPEN THE BOX", s.popup) {
		s.popup = true
		log("ooh")
	}
	if game.LootPopup(c, data.Item{Sprite: icon.Potion, Name: "POTION", Count: 3}, &s.popup, game.LootPopupOptions{}).Closed {
		log("put it in the bag")
	}
}

func demoTreasureChest(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := lootState(c)
	ui.Row(c).Gap(k.D(10)).AlignItems(ui.End).Children(func() {
		if game.TreasureChest(c, "Chest", data.Item{Sprite: icon.Gem, Name: "Gem"}, &s.chest, game.TreasureChestOptions{}).Opened {
			log("a gem!")
		}
		if small(c, "CLOSE", !s.chest) {
			s.chest = false
			log("closed")
		}
	})
}

// faces are the gacha reel's faces.
var faces = []data.Item{{Sprite: icon.Star, Name: "Star"}, {Sprite: icon.Heart, Name: "Heart"}, {Sprite: icon.Coin, Name: "Coin"}, {Sprite: icon.Skull, Name: "Skull"}, {Sprite: icon.Gem, Name: "Gem"}}

func demoGachaReel(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := lootState(c)
	dice := gameState(c).dice
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		if game.GachaReel(c, "Reels", faces, s.stops, s.spin, game.GachaReelOptions{}).Done {
			log(fmt.Sprintf("%s %s %s", faces[s.stops[0]].Name, faces[s.stops[1]].Name, faces[s.stops[2]].Name))
		}
		if small(c, "SPIN", false) {
			s.stops = [3]int{dice.IntN(len(faces)), dice.IntN(len(faces)), dice.IntN(len(faces))}
			s.spin++
			log("spinning")
		}
	})
}

func demoInteractPrompt(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := lootState(c)
	ui.Row(c).Gap(k.D(16)).AlignItems(ui.End).Padding(k.D(24), 0, 0, 0).Children(func() {
		mayor := ui.Box(c).Size(k.D(24), k.D(24)).Label("Mayor").Draw(func(p *ui.Painter, r ui.Rect) {
			k.Canvas(p, r).Sprite(0, 0, game.Hero, 0, pixel.Ink, 2)
		})
		if s.near {
			game.InteractPrompt(c, mayor, "TALK", game.InteractPromptOptions{})
		}
		if input.Toggle(c, "NEAR THE MAYOR", &s.near, input.ToggleOptions{Size: input.Small}).Changed {
			log(fmt.Sprintf("near %v", s.near))
		}
	})
}

// fighters are the turn order demo's lineup.
var fighters = []game.Fighter{{Name: "Aria", Portrait: game.Hero}, {Name: "Slime", Portrait: slime, Foe: true}, {Name: "Brom", Portrait: game.Hero}, {Name: "Bat", Portrait: icon.Bug, Foe: true}}

func demoTurnOrder(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := lootState(c)
	ui.Row(c).Gap(k.D(10)).AlignItems(ui.Center).Children(func() {
		game.TurnOrder(c, "Turns", fighters, s.turn, game.TurnOrderOptions{})
		if small(c, "NEXT TURN", false) {
			s.turn = (s.turn + 1) % len(fighters)
			log(fighters[s.turn].Name + "'s turn")
		}
	})
}

func demoBattleLog(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := lootState(c)
	dice := gameState(c).dice
	ui.Column(c).Gap(k.D(6)).Children(func() {
		game.BattleLog(c, "Battle", s.lines, game.BattleLogOptions{})
		ui.Row(c).Gap(k.D(4)).Children(func() {
			if small(c, "ATTACK", false) {
				line := game.Line{Text: fmt.Sprintf("Slime takes %d.", 3+dice.IntN(9)), Slot: pixel.Bad}
				s.lines = append(s.lines, game.Line{Text: "Aria attacks!"}, line)
				log(line.Text)
			}
			if small(c, "HEAL", false) {
				s.lines = append(s.lines, game.Line{Text: "Brom heals 8.", Slot: pixel.Good})
				log("healed")
			}
		})
	})
}

func demoTargetSelector(c *ui.Context, log func(string)) {
	s := lootState(c)
	enemies := []game.Enemy{{Name: "Slime", Sprite: slime, HP: 6, MaxHP: 9}, {Name: "Bat", Sprite: icon.Bug, HP: 5, MaxHP: 5}, {Name: "Ghost", Sprite: icon.Skull, MaxHP: 7}}
	r := game.TargetSelector(c, "Targets", enemies, &s.target, game.TargetSelectorOptions{})
	if r.Changed {
		log("aim at " + enemies[s.target].Name)
	}
	if r.Chosen >= 0 {
		log("attack " + enemies[r.Chosen].Name)
	}
}

func demoATBGauge(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := lootState(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		for i, h := range []struct {
			name  string
			speed int
		}{{"ARIA", 30}, {"BROM", 45}, {"CYRA", 20}} {
			ui.Row(c).Gap(k.D(6)).AlignItems(ui.End).Children(func() {
				// Each gauge fills in its hero's speed, in tenths.
				fill := k.Since(c, s.ready[i], 10, h.speed)
				game.ATBGauge(c, h.name, float64(fill)/float64(h.speed), game.ATBGaugeOptions{Width: 80})
				if small(c, "ACT", fill < h.speed) {
					s.ready[i] = k.Now()
					log(h.name + " acts")
				}
			})
		}
	})
}

func demoStaminaBar(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := lootState(c)
	now := k.Now()
	// Sprinting drains a quarter a second; rest refills a tenth.
	if dt := max(0, now.Sub(s.last).Seconds()); s.sprint {
		s.stamina = max(0, s.stamina-dt/4)
	} else {
		s.stamina = min(1, s.stamina+dt/10)
	}
	s.last = now
	if s.sprint && s.stamina == 0 {
		s.sprint = false
		log("out of breath")
	}
	if s.sprint || s.stamina < 1 {
		k.Tick(c, 10)
	}
	ui.Row(c).Gap(k.D(10)).AlignItems(ui.End).Children(func() {
		game.StaminaBar(c, "STAMINA", s.stamina, game.StaminaBarOptions{})
		if input.Toggle(c, "SPRINT", &s.sprint, input.ToggleOptions{Size: input.Small, Disabled: s.stamina == 0}).Changed {
			log(fmt.Sprintf("sprint %v", s.sprint))
		}
	})
}
