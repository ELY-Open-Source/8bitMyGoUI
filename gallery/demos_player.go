package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/account"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var playerEntries = []Entry{
	{396, "PalettePicker", "调色板选择", "account", "Each palette as a little live scene. Choose one and the sample wears it.", demoPalettePicker},
	{397, "VolumeMixer", "音量混合", "account", "Music, effects and voice, each on its own fader.", demoVolumeMixer},
	{398, "DifficultySelect", "难度选择", "account", "From EASY to NIGHTMARE, one skull a step. Choose wisely.", demoDifficultySelect},
	{399, "AccessibilityOptions", "无障碍", "account", "Reduce motion and the dot size; the sample below obeys at once.", demoAccessibilityOptions},
	{400, "AccountSwitcher", "账号切换", "account", "Four player slots in arcade colors. Empty ones wait for a friend.", demoAccountSwitcher},
	{401, "PlayerSelect", "玩家选择", "account", "1 PLAYER or 2 PLAYERS: the oldest question on the title screen.", demoPlayerSelect},
	{402, "TwoFactor", "两步验证", "account", "Six digits, typed or pasted. The bar is how long the code lives.", demoTwoFactor},
	{403, "About", "关于", "account", "The logo, the version, and the people to thank.", demoAbout},
	{404, "PatchNotes", "更新日志", "account", "What came, what was fixed, what changed, version by version.", demoPatchNotes},
	{405, "FeedbackForm", "反馈表", "account", "Stars and words. SEND waits until you rate.", demoFeedbackForm},
}

type playerDemo struct {
	palette, difficulty, slot, players, rating int
	mix                                        account.Mix
	access                                     account.Access
	code, words                                string
}

func playerState(c *ui.Context) *playerDemo {
	return ui.Local(c.Root(), demoKey("player"), func() playerDemo {
		return playerDemo{difficulty: 1, players: 1, mix: account.Mix{BGM: 0.8, SE: 0.6, Voice: 0.9}, access: account.Access{Dot: 2}}
	})
}

func demoPalettePicker(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := playerState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		if account.PalettePicker(c, "Palette", &s.palette, account.PalettePickerOptions{}).Changed {
			log("palette " + pixel.Palettes[s.palette].Name)
		}
		// Choosing applies: the sample draws in the chosen palette.
		in := pixel.Current(c)
		in.Palette = pixel.Palettes[s.palette]
		pixel.Scope(c, in, func() {
			display.Sprite(c, slime, display.SpriteOptions{FPS: 4, Scale: 3})
		})
	})
}

func demoVolumeMixer(c *ui.Context, log func(string)) {
	s := playerState(c)
	if account.VolumeMixer(c, "Mixer", &s.mix, account.VolumeMixerOptions{}).Changed {
		log(fmt.Sprintf("mix %.2f %.2f %.2f", s.mix.BGM, s.mix.SE, s.mix.Voice))
	}
}

func demoDifficultySelect(c *ui.Context, log func(string)) {
	s := playerState(c)
	if account.DifficultySelect(c, "Difficulty", &s.difficulty, account.DifficultySelectOptions{}).Changed {
		log(fmt.Sprintf("difficulty %d", s.difficulty+1))
	}
}

func demoAccessibilityOptions(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := playerState(c)
	ui.Column(c).Gap(k.D(8)).Children(func() {
		if account.AccessibilityOptions(c, "Access", &s.access, account.AccessibilityOptionsOptions{}).Changed {
			log(fmt.Sprintf("reduce %v, dot %d", s.access.Reduce, s.access.Dot))
		}
		in := pixel.Current(c)
		in.Dot, in.Reduce = s.access.Dot, s.access.Reduce
		pixel.Scope(c, in, func() {
			display.Sprite(c, slime, display.SpriteOptions{FPS: 4, Scale: 2})
		})
	})
}

func demoAccountSwitcher(c *ui.Context, log func(string)) {
	s := playerState(c)
	if account.AccountSwitcher(c, "Players", [4]string{"Aria", "Brom", "", ""}, &s.slot, account.AccountSwitcherOptions{}).Changed {
		log(fmt.Sprintf("P%d", s.slot+1))
	}
}

func demoPlayerSelect(c *ui.Context, log func(string)) {
	s := playerState(c)
	r := account.PlayerSelect(c, "Players", &s.players, account.PlayerSelectOptions{})
	if r.Changed {
		log(fmt.Sprintf("%d players", s.players))
	}
	if r.Chosen {
		log(fmt.Sprintf("start %dP", s.players))
	}
}

func demoTwoFactor(c *ui.Context, log func(string)) {
	s := playerState(c)
	r := account.TwoFactor(c, "Login code", &s.code, account.TwoFactorOptions{})
	if r.Changed {
		log(fmt.Sprintf("%d digits", len(s.code)))
	}
	if r.Submitted {
		log("code sent")
	}
}

func demoAbout(c *ui.Context, log func(string)) {
	account.About(c, "8BIT", "1.0.0", account.AboutOptions{Credits: []account.Credit{
		{Role: "Pixels and code", Name: "The 8bit team"}, {Role: "Native UI", Name: "MyGo"}, {Role: "Font", Name: "Ark Pixel"},
	}})
}

func demoPatchNotes(c *ui.Context, log func(string)) {
	account.PatchNotes(c, "Patch notes", []account.Patch{
		{Version: "v1.1.0", Added: []string{"Palette cycling waterfalls"}, Fixed: []string{"The ghost in level 3 now leaves"}, Changed: []string{"Jump height +1 dot"}},
		{Version: "v1.0.0", Added: []string{"Everything"}},
	}, account.PatchNotesOptions{})
}

func demoFeedbackForm(c *ui.Context, log func(string)) {
	s := playerState(c)
	if account.FeedbackForm(c, "Feedback", &s.rating, &s.words, account.FeedbackFormOptions{}).Sent {
		log(fmt.Sprintf("sent %d stars", s.rating))
		s.rating, s.words = 0, ""
	}
}
