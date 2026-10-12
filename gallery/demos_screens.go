package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/game"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var screenEntries = []Entry{
	{269, "SaveSlots", "存档位", "game", "Three files: where you are and how long you've played. Empty ones start anew.", demoSaveSlots},
	{270, "TitleScreen", "标题画面", "game", "Logo, PRESS START, then the menu. Escape goes back.", demoTitleScreen},
	{271, "PressStart", "按下开始", "game", "Blinks politely until you press something.", demoPressStart},
	{272, "OptionsMenu", "选项菜单", "game", "Up and Down pick a row; ◀ ▶ change it.", demoOptionsMenu},
	{273, "GameOver", "游戏结束", "game", "GAME OVER dithers in. CONTINUE? Of course you will.", demoGameOver},
	{274, "Credits", "制作人员", "game", "Credits roll up and rest on the last name.", demoCredits},
	{275, "EndingCard", "结局卡", "game", "THE END, dithering in on the dark.", demoEndingCard},
	{276, "LoadingTips", "加载提示", "game", "A tip every four seconds while the bar fills.", demoLoadingTips},
	{277, "Letterbox", "宽银幕", "game", "Cutscene bars slide in and out in four steps.", demoLetterbox},
	{278, "MissionBrief", "任务简报", "game", "The map, the objectives, and a big START.", demoMissionBrief},
	{279, "Achievements", "成就墙", "game", "Earned awards shine; locked ones are shadows with a hint.", demoAchievements},
}

type screenDemo struct {
	slot, row, award, credits, ending, load int
	speed, bgm, sfx, rumble                 int
	started, cutscene                       bool
	since                                   time.Time
}

func screenState(c *ui.Context) *screenDemo {
	now := pixel.K(c).Now()
	return ui.Local(c.Root(), demoKey("screens"), func() screenDemo {
		return screenDemo{speed: 1, bgm: 7, sfx: 5, since: now}
	})
}

func demoSaveSlots(c *ui.Context, log func(string)) {
	s := screenState(c)
	saves := []game.Save{{Place: "Slime Woods", Played: 3*time.Hour + 25*time.Minute, Portrait: game.Hero}, {}, {Place: "Castle Gate", Played: 12 * time.Minute, Portrait: game.Hero}}
	r := game.SaveSlots(c, "Save files", saves, &s.slot, game.SaveSlotsOptions{})
	if r.Changed {
		log(fmt.Sprintf("file %d", s.slot+1))
	}
	if r.Chosen >= 0 {
		log(fmt.Sprintf("load file %d", r.Chosen+1))
	}
}

func demoTitleScreen(c *ui.Context, log func(string)) {
	s := screenState(c)
	was := s.started
	r := game.TitleScreen(c, "SLIME QUEST", &s.started, game.TitleScreenOptions{Saved: true})
	if s.started != was {
		log(fmt.Sprintf("started %v", s.started))
	}
	if r.Chosen >= 0 {
		log([]string{"new game", "continue", "options"}[r.Chosen])
	}
}

func demoPressStart(c *ui.Context, log func(string)) {
	if game.PressStart(c, game.PressStartOptions{}).Pressed {
		log("here we go!")
	}
}

func demoOptionsMenu(c *ui.Context, log func(string)) {
	s := screenState(c)
	settings := []game.Setting{
		{Name: "TEXT SPEED", Choices: []string{"SLOW", "NORMAL", "FAST"}, Value: &s.speed},
		{Name: "BGM", Max: 10, Value: &s.bgm}, {Name: "SFX", Max: 10, Value: &s.sfx},
		{Name: "RUMBLE", Choices: []string{"OFF", "ON"}, Value: &s.rumble},
	}
	if game.OptionsMenu(c, "Options", settings, &s.row, game.OptionsMenuOptions{}).Changed {
		log(fmt.Sprintf("speed %d, bgm %d, sfx %d, rumble %d", s.speed, s.bgm, s.sfx, s.rumble))
	}
}

func demoGameOver(c *ui.Context, log func(string)) {
	r := game.GameOver(c, game.GameOverOptions{})
	if r.Continue {
		log("never give up")
	}
	if r.Quit {
		log("see you next time")
	}
}

func demoCredits(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := screenState(c)
	ui.Column(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		credits := []game.Credit{{Heading: "DIRECTOR", Names: []string{"A Slime"}}, {Heading: "PIXELS", Names: []string{"Aria", "Brom"}},
			{Heading: "MUSIC", Names: []string{"The Bard"}}, {Heading: "THANK YOU", Names: []string{"for playing!"}}}
		if game.Credits(c, "Credits", credits, game.CreditsOptions{Cue: s.credits}).Done {
			log("the end of the credits")
		}
		if small(c, "ROLL AGAIN", false) {
			s.credits++
			log("rolling")
		}
	})
}

func demoEndingCard(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := screenState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		game.EndingCard(c, game.EndingCardOptions{Cue: s.ending})
		if small(c, "AGAIN", false) {
			s.ending++
			log("fade in")
		}
	})
}

func demoLoadingTips(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := screenState(c)
	// The pretend load takes eight seconds, then starts over.
	done := k.Since(c, s.since, 10, 80)
	if done == 80 {
		s.since, s.load = k.Now(), s.load+1
		log(fmt.Sprintf("loaded %d times", s.load))
	}
	tips := []string{"Slimes fear fire. Mostly.", "Save often. Then save again.", "The cat is in the tower."}
	game.LoadingTips(c, tips, float64(done)/80, game.LoadingTipsOptions{})
}

func demoLetterbox(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := screenState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		game.Letterbox(c, s.cutscene, game.LetterboxOptions{}, func() {
			ui.Box(c).Size(k.D(200), k.D(80)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Dither(0, 0, g.W, g.H, pixel.Well, pixel.Info, 4)
				g.Rect(0, g.H-16, g.W, 16, pixel.Good)
				g.Sprite(90, g.H-16-slime.H, slime, 0, pixel.Ink, 1)
			})
		})
		if input.Toggle(c, "CUTSCENE", &s.cutscene, input.ToggleOptions{Size: input.Small}).Changed {
			log(fmt.Sprintf("cutscene %v", s.cutscene))
		}
	})
}

func demoMissionBrief(c *ui.Context, log func(string)) {
	q := game.Quest{Title: "Rescue Mittens", Objectives: []game.Objective{{Text: "Reach the tower", Need: 1}, {Text: "Beat the bats", Have: 1, Need: 3}}}
	area := []string{`############`, `#..""..~~..#`, `#..""...~..#`, `#....##....#`, `############`}
	if game.MissionBrief(c, q, area, game.Spot{X: 2, Y: 3}, []game.Marker{{Spot: game.Spot{X: 9, Y: 1}, Slot: pixel.Warn}}, game.MissionBriefOptions{}).Start {
		log("mission start!")
	}
}

func demoAchievements(c *ui.Context, log func(string)) {
	s := screenState(c)
	awards := []game.Award{
		{Sprite: icon.Sword, Name: "First blood", How: "Win your first fight.", Unlocked: true},
		{Sprite: icon.Coin, Name: "Rich", How: "Carry 1000 gold.", Unlocked: true},
		{Sprite: icon.Crown, Name: "Kingslayer", How: "Beat the Slime King."},
		{Sprite: icon.Heart, Name: "Cat person", How: "Bring Mittens home.", Unlocked: true},
		{Sprite: icon.Skull, Name: "Oops", How: "Fall in a pit 10 times."},
		{Sprite: icon.Star, Name: "Completionist", How: "Earn every award."},
	}
	if game.Achievements(c, "Awards", awards, &s.award, game.AchievementsOptions{Columns: 3}).Changed {
		log(awards[s.award].Name)
	}
}
