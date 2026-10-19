package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/game"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/layout"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var rpgEntries = []Entry{
	{258, "SkillTree", "技能树", "game", "Spend points on open skills: click once to choose, twice to learn.", demoSkillTree},
	{259, "SkillButton", "技能按钮", "game", "Cast, then wait out the clockwise cooldown wipe.", demoSkillButton},
	{260, "StatusEffects", "状态效果", "game", "Buffs edged green, debuffs red, turns left in the corner.", demoStatusEffects},
	{261, "BossBar", "首领血条", "game", "A boss in phases. Hit it until it gets angry.", demoBossBar},
	{262, "DamageNumbers", "伤害数字", "game", "Numbers float up and dissolve. Criticals in gold!", demoDamageNumbers},
	{263, "QuestLog", "任务日志", "game", "Choose a quest to track; its objectives show beside it.", demoQuestLog},
	{264, "QuestTracker", "任务追踪", "game", "Compact objectives with counts. Go slay one.", demoQuestTracker},
	{265, "PartyFrames", "队伍面板", "game", "Cards with HP and MP. Brom needs a phoenix down.", demoPartyFrames},
	{266, "CharacterSheet", "角色面板", "game", "Who you are: class, level and stats.", demoCharacterSheet},
	{267, "LevelUpScreen", "升级分配", "game", "LEVEL UP! Spend points with ◀ ▶, then DONE.", demoLevelUpScreen},
	{268, "Shop", "商店", "game", "Buy at the price, sell at half. Too dear? It goes grey.", demoShop},
}

type rpgDemo struct {
	skills                      []game.Skill
	cursor, points, boss, quest int
	slain                       int
	cast                        time.Time
	effects                     []game.Effect
	hits                        []game.Hit
	attrs                       []game.Attribute
	wares                       []game.Ware
	gold, raise                 int
}

func rpgState(c *ui.Context) *rpgDemo {
	return ui.Local(c.Root(), demoKey("rpg"), func() rpgDemo {
		return rpgDemo{skills: talents(), points: 2, boss: 1000, gold: 300, raise: 3,
			effects: []game.Effect{{Sprite: icon.Bolt, Name: "Haste", Turns: 3}, {Sprite: icon.Shield, Name: "Guard", Turns: 2}, {Sprite: icon.Skull, Name: "Poison", Turns: 4, Debuff: true}},
			attrs:   []game.Attribute{{Name: "STR", Base: 12}, {Name: "DEX", Base: 9}, {Name: "INT", Base: 7}, {Name: "LUK", Base: 5}},
			wares: []game.Ware{{Item: data.Item{Sprite: icon.Potion, Name: "Potion"}, Price: 30, Owned: 2}, {Item: data.Item{Sprite: icon.Key, Name: "Skeleton key"}, Price: 120},
				{Item: data.Item{Sprite: icon.Sword, Name: "Flame sword"}, Price: 900}}}
	})
}

// talents is a root, two branches and a capstone.
func talents() []game.Skill {
	return []game.Skill{
		{Name: "Slash", Sprite: icon.Sword, At: game.Spot{X: 1}, Learned: true},
		{Name: "Guard", Sprite: icon.Shield, At: game.Spot{Y: 1}, Needs: []int{0}},
		{Name: "Fireball", Sprite: icon.Flame, At: game.Spot{X: 2, Y: 1}, Needs: []int{0}},
		{Name: "Bulwark", Sprite: icon.Crown, At: game.Spot{Y: 2}, Needs: []int{1}},
		{Name: "Meteor", Sprite: icon.Star, At: game.Spot{X: 2, Y: 2}, Needs: []int{2}},
		{Name: "Starfall", Sprite: icon.Sparkle, At: game.Spot{X: 1, Y: 3}, Needs: []int{3, 4}},
	}
}

func demoSkillTree(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := rpgState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		r := game.SkillTree(c, "Talents", s.skills, &s.cursor, &s.points, game.SkillTreeOptions{})
		if r.Changed {
			log(s.skills[s.cursor].Name)
		}
		if r.Learned >= 0 {
			log("learned " + s.skills[r.Learned].Name)
		}
		if small(c, "RESPEC", s.points == 2) {
			s.skills, s.points = talents(), 2
			log("a clean slate")
		}
	})
}

func demoSkillButton(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := rpgState(c)
	// The spell cools for three seconds, in tenths.
	cool := 1 - float64(k.Since(c, s.cast, 10, 30))/30
	ui.Row(c).Gap(k.D(6)).Children(func() {
		if game.SkillButton(c, "Fireball", icon.Flame, "Q", cool, game.SkillButtonOptions{}).Clicked {
			s.cast = k.Now()
			log("fireball!")
		}
		if game.SkillButton(c, "Heal", icon.Heart, "E", 0, game.SkillButtonOptions{}).Clicked {
			log("you feel better")
		}
		game.SkillButton(c, "Meteor", icon.Star, "R", 0, game.SkillButtonOptions{Disabled: true})
	})
}

func demoStatusEffects(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := rpgState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		game.StatusEffects(c, "Effects", s.effects, game.StatusEffectsOptions{})
		if small(c, "END TURN", len(s.effects) == 0) {
			left := s.effects[:0]
			for _, ef := range s.effects {
				if ef.Turns--; ef.Turns > 0 {
					left = append(left, ef)
				}
			}
			s.effects = left
			log(fmt.Sprintf("%d effects left", len(s.effects)))
		}
	})
}

func demoBossBar(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := rpgState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		game.BossBar(c, "SLIME KING", s.boss, 1000, game.BossBarOptions{Width: 220, Phases: []float64{0.66, 0.33}})
		if small(c, "ATTACK", s.boss == 0) {
			was := s.boss
			s.boss = max(0, s.boss-90)
			log(fmt.Sprintf("hit! %d HP left", s.boss))
			for _, mark := range []int{660, 330} {
				if was > mark && s.boss <= mark {
					log("the Slime King wobbles angrily!")
				}
			}
			if s.boss == 0 {
				log("victory!")
			}
		}
	})
}

func demoDamageNumbers(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s, dice := rpgState(c), gameState(c).dice
	ui.Column(c).Gap(k.D(4)).Children(func() {
		// The numbers float over a slime having a bad day.
		layout.Layers(c, layout.LayersOptions{}, func() {
			ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Sprite((g.W-3*slime.W)/2, g.H-3*slime.H, slime, 0, pixel.Ink, 3)
			})
		}, func() {
			game.DamageNumbers(c, "Damage", &s.hits, game.DamageNumbersOptions{Width: 180, Height: 60})
		})
		if small(c, "ATTACK", false) {
			h := game.Hit{Value: 5 + dice.IntN(30), X: 20 + dice.IntN(140), Y: 25 + dice.IntN(25)}
			// One hit in five is critical, for double.
			if h.Crit = dice.IntN(5) == 0; h.Crit {
				h.Value *= 2
			}
			s.hits = append(s.hits, h)
			log(fmt.Sprintf("hit for %d", h.Value))
		}
	})
}

// adventures are the quest demos' quests.
func adventures(slain int) []game.Quest {
	return []game.Quest{
		{Title: "Slime trouble", Objectives: []game.Objective{{Text: "Slay slimes", Have: slain, Need: 5}, {Text: "Tell the mayor", Have: slain / 5, Need: 1}}},
		{Title: "Lost cat", Objectives: []game.Objective{{Text: "Find Mittens", Have: 1, Need: 1}, {Text: "Bring her home", Need: 1}}},
		{Title: "Herb run", Objectives: []game.Objective{{Text: "Pick moonleaf", Have: 2, Need: 6}}},
	}
}

func demoQuestLog(c *ui.Context, log func(string)) {
	s := rpgState(c)
	quests := adventures(s.slain)
	if game.QuestLog(c, "Quests", quests, &s.quest, game.QuestLogOptions{}).Changed {
		log("tracking " + quests[s.quest].Title)
	}
}

func demoQuestTracker(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := rpgState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		game.QuestTracker(c, adventures(s.slain)[0], game.QuestTrackerOptions{})
		if small(c, "SLAY A SLIME", s.slain == 5) {
			s.slain++
			log(fmt.Sprintf("%d of 5", s.slain))
		}
	})
}

func demoPartyFrames(c *ui.Context, log func(string)) {
	game.PartyFrames(c, "Party", []game.Member{
		{Name: "Aria", Portrait: game.Hero, HP: 34, MaxHP: 40, MP: 18, MaxMP: 20},
		{Name: "Cyra", Portrait: game.Hero, HP: 9, MaxHP: 30, MP: 31, MaxMP: 40},
		{Name: "Brom", Portrait: game.Hero, MaxHP: 50, MP: 4, MaxMP: 10},
	}, game.PartyFramesOptions{})
}

func demoCharacterSheet(c *ui.Context, log func(string)) {
	game.CharacterSheet(c, game.Character{Name: "Aria", Class: "Knight", Level: 7, Portrait: game.Hero,
		Stats: []game.Stat{{Name: "STR", Value: 14}, {Name: "DEX", Value: 9}, {Name: "INT", Value: 7}, {Name: "LUK", Value: 11}}}, game.CharacterSheetOptions{})
}

func demoLevelUpScreen(c *ui.Context, log func(string)) {
	s := rpgState(c)
	r := game.LevelUpScreen(c, "Level up", s.attrs, &s.raise, game.LevelUpScreenOptions{})
	if r.Changed {
		log(fmt.Sprintf("%d points left", s.raise))
	}
	if r.Done {
		for i := range s.attrs {
			s.attrs[i].Base, s.attrs[i].Added = s.attrs[i].Base+s.attrs[i].Added, 0
		}
		s.raise = 3
		log("stronger! three more points, for the demo")
	}
}

func demoShop(c *ui.Context, log func(string)) {
	s := rpgState(c)
	r := game.Shop(c, "Item shop", s.wares, &s.gold, game.ShopOptions{})
	if r.Bought >= 0 {
		log("bought " + s.wares[r.Bought].Item.Name)
	}
	if r.Sold >= 0 {
		log("sold " + s.wares[r.Sold].Item.Name)
	}
}
