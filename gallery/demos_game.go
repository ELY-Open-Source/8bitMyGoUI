package gallery

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/game"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var gameEntries = []Entry{
	{246, "HealthBar", "生命条", "game", "HIT and HEAL. The pink trail lags; under a quarter the bar blinks.", demoHealthBar},
	{247, "ManaBar", "魔力条", "game", "A spell costs 12 MP. Ticks mark every 10.", demoManaBar},
	{248, "ExpBar", "经验条", "game", "Slay slimes for XP. A full bar flashes and the level goes up.", demoExpBar},
	{249, "Hearts", "心形血量", "game", "Damage in quarter hearts. A potion mends a whole one.", demoHearts},
	{250, "Lives", "剩余命数", "game", "Your face, times lives. 1UP!", demoLives},
	{251, "CoinCounter", "金币", "game", "The coin spins; the count ticks up to the new total.", demoCoinCounter},
	{252, "ScoreBoard", "计分板", "game", "Beat the HI-SCORE and its label lights up.", demoScoreBoard},
	{253, "TimerHUD", "关卡计时", "game", "TIME runs out. Under 30 it turns red and blinks. Hurry!", demoTimerHUD},
	{254, "Minimap", "小地图", "game", "Walk with the D-pad; the fog lifts around you.", demoMinimap},
	{255, "Compass", "指南针", "game", "Turn the heading; the quest and town markers slide along.", demoCompass},
	{256, "Hotbar", "快捷栏", "game", "Press 1–9, scroll, or click. Enter uses the slot.", demoHotbar},
	{257, "Equipment", "装备栏", "game", "A paper doll. Arrows walk the slots; Enter opens one.", demoEquipment},
}

type gameDemo struct {
	hp, mp, level, xp, lives, coins, score, slot, gear int
	hearts, heading                                    float64
	deadline                                           time.Time
	pos                                                game.Spot
	seen                                               [][]bool
	dice                                               *rand.Rand
}

func gameState(c *ui.Context) *gameDemo {
	now := pixel.K(c).Now()
	return ui.Local(c.Root(), demoKey("game"), func() gameDemo {
		s := gameDemo{hp: 100, mp: 60, level: 1, lives: 3, coins: 120, hearts: 5, heading: 30, gear: game.Weapon,
			deadline: now.Add(100 * time.Second), pos: game.Spot{X: 2, Y: 2}, dice: rand.New(rand.NewPCG(3, 4))}
		for range dungeon {
			s.seen = append(s.seen, make([]bool, len(dungeon[0])))
		}
		s.look()
		return s
	})
}

// small is a small button.
func small(c *ui.Context, text string, disabled bool) bool {
	return input.Button(c, text, input.ButtonOptions{Size: input.Small, Disabled: disabled}).Clicked
}

func demoHealthBar(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := gameState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		game.HealthBar(c, "HP", s.hp, 100, game.HealthBarOptions{Width: 140})
		ui.Row(c).Gap(k.D(4)).Children(func() {
			if small(c, "HIT", s.hp == 0) {
				s.hp = max(0, s.hp-18)
				log(fmt.Sprintf("ouch! %d HP", s.hp))
			}
			if small(c, "HEAL", s.hp == 100) {
				s.hp = min(100, s.hp+30)
				log(fmt.Sprintf("healed to %d HP", s.hp))
			}
		})
	})
}

func demoManaBar(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := gameState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		game.ManaBar(c, "MP", s.mp, 60, game.ManaBarOptions{Width: 140})
		ui.Row(c).Gap(k.D(4)).Children(func() {
			if small(c, "FIREBALL", s.mp < 12) {
				s.mp -= 12
				log(fmt.Sprintf("fwoosh! %d MP left", s.mp))
			}
			if small(c, "REST", s.mp == 60) {
				s.mp = min(60, s.mp+20)
				log("zzz")
			}
		})
	})
}

func demoExpBar(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := gameState(c)
	next := s.level * 100
	ui.Column(c).Gap(k.D(6)).Children(func() {
		game.ExpBar(c, "EXP", s.level, s.xp, next, game.ExpBarOptions{Width: 140})
		if small(c, "SLAY A SLIME", false) {
			log("+45 XP")
			if s.xp += 45; s.xp >= next {
				s.level, s.xp = s.level+1, s.xp-next
				log(fmt.Sprintf("level up! LV %d", s.level))
			}
		}
	})
}

func demoHearts(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := gameState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		game.Hearts(c, "Life", s.hearts, 5, game.HeartsOptions{})
		ui.Row(c).Gap(k.D(4)).Children(func() {
			if small(c, "OUCH", s.hearts == 0) {
				s.hearts = max(0, s.hearts-0.75)
				log(fmt.Sprintf("%.2f hearts", s.hearts))
			}
			if small(c, "POTION", s.hearts == 5) {
				s.hearts = min(5, s.hearts+1)
				log("glug glug")
			}
		})
	})
}

func demoLives(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := gameState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		game.Lives(c, "Lives", s.lives, game.LivesOptions{})
		if small(c, "1UP", false) {
			s.lives++
			log("1UP!")
		}
		if small(c, "FALL IN A PIT", s.lives == 0) {
			s.lives--
			log("aaaaa...")
		}
	})
}

func demoCoinCounter(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := gameState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		game.CoinCounter(c, "Gold", s.coins, game.CoinCounterOptions{Scale: 2})
		if small(c, "OPEN CHEST", false) {
			s.coins += 250
			log("+250 gold")
		}
	})
}

func demoScoreBoard(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := gameState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		game.ScoreBoard(c, s.score, 9900, game.ScoreBoardOptions{})
		ui.Row(c).Gap(k.D(4)).Children(func() {
			if small(c, "STOMP", false) {
				s.score += 1500
				log(fmt.Sprintf("score %d", s.score))
			}
			if small(c, "RESET", s.score == 0) {
				s.score = 0
				log("score reset")
			}
		})
	})
}

func demoTimerHUD(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := gameState(c)
	left := max(0, s.deadline.Sub(k.Now()))
	if left > 0 {
		k.Tick(c, 4)
	}
	ui.Row(c).Gap(k.D(12)).AlignItems(ui.Center).Children(func() {
		game.TimerHUD(c, left, game.TimerHUDOptions{})
		if small(c, "RESTART", false) {
			s.deadline = k.Now().Add(100 * time.Second)
			log("go go go!")
		}
		if small(c, "SKIP TO 33", left <= 33*time.Second) {
			s.deadline = k.Now().Add(33 * time.Second)
			log("33 seconds left")
		}
	})
}

// dungeon is the minimap demo's map.
var dungeon = []string{
	`############`,
	`#....~~....#`,
	`#.""...~...#`,
	`#.""..###..#`,
	`#.....#^^..#`,
	`############`,
}

// look lifts the fog two tiles around the player.
func (s *gameDemo) look() {
	for y := range dungeon {
		for x := range dungeon[y] {
			if max(x-s.pos.X, s.pos.X-x, y-s.pos.Y, s.pos.Y-y) <= 2 {
				s.seen[y][x] = true
			}
		}
	}
}

func demoMinimap(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := gameState(c)
	ui.Row(c).Gap(k.D(10)).AlignItems(ui.Center).Children(func() {
		game.Minimap(c, "Dungeon", dungeon, s.seen, s.pos, []game.Marker{{Spot: game.Spot{X: 9, Y: 4}, Slot: pixel.Warn}}, game.MinimapOptions{Cell: 6})
		step := map[input.Direction]game.Spot{input.Up: {Y: -1}, input.Down: {Y: 1}, input.Left: {X: -1}, input.Right: {X: 1}}
		if d, ok := step[input.DPad(c, "Walk", input.DPadOptions{}).Pressed]; ok {
			to := game.Spot{X: s.pos.X + d.X, Y: s.pos.Y + d.Y}
			if dungeon[to.Y][to.X] != '#' {
				s.pos = to
				s.look()
				log(fmt.Sprintf("at %d, %d", to.X, to.Y))
			}
		}
	})
}

func demoCompass(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := gameState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		game.Compass(c, "Heading", s.heading, []game.Bearing{{Degrees: 120, Slot: pixel.Warn}, {Degrees: 300, Slot: pixel.Info}}, game.CompassOptions{Width: 160})
		if input.Slider(c, "Heading", &s.heading, input.SliderOptions{Max: 360, Width: 160}).Changed {
			log(fmt.Sprintf("%.0f°", s.heading))
		}
	})
}

// pack is the hotbar demo's items.
var pack = []data.Item{
	{Sprite: icon.Sword, Name: "Sword", Rarity: display.Rare}, {Sprite: icon.Shield, Name: "Shield"}, {Sprite: icon.Potion, Name: "Potion", Count: 5},
	{Sprite: icon.Bolt, Name: "Bolt scroll", Rarity: display.Epic}, {Sprite: icon.Key, Name: "Key"}, {}, {Sprite: icon.Map, Name: "Map"}, {}, {Sprite: icon.Gem, Name: "Gem", Count: 3},
}

func demoHotbar(c *ui.Context, log func(string)) {
	s := gameState(c)
	r := game.Hotbar(c, "Belt", pack, &s.slot, game.HotbarOptions{})
	if r.Changed {
		log(fmt.Sprintf("slot %d", s.slot+1))
	}
	if r.Used >= 0 {
		name := pack[r.Used].Name
		if name == "" {
			name = "an empty hand"
		}
		log("used " + name)
	}
}

func demoEquipment(c *ui.Context, log func(string)) {
	s := gameState(c)
	gear := []data.Item{{Sprite: icon.Crown, Name: "Old crown", Rarity: display.Legendary}, {}, {Sprite: icon.Sword, Name: "Iron sword"}, {Sprite: icon.Shield, Name: "Pot lid"}, {}}
	r := game.Equipment(c, "Gear", gear, &s.gear, game.EquipmentOptions{})
	if r.Changed {
		log(fmt.Sprintf("slot %d", s.gear))
	}
	if r.Picked >= 0 {
		log(fmt.Sprintf("open slot %d", r.Picked))
	}
}
