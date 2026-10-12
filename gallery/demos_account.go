package gallery

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/account"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var accountEntries = []Entry{
	{386, "SettingsList", "设置列表", "account", "Labels on the left, controls on the right, sections on top. Tidy.", demoSettingsList},
	{387, "OptionCycler", "选项切换", "account", "◀ and ▶ walk the choices; past the end it comes around.", demoOptionCycler},
	{388, "KeyBindings", "按键绑定", "account", "Enter, then press the new key. Two actions on one key glow red.", demoKeyBindings},
	{389, "ControllerMap", "手柄图", "account", "Which button does what, with lines to prove it.", demoControllerMap},
	{390, "ProfileCard", "玩家卡", "account", "Face, name, level, title and the numbers that matter.", demoProfileCard},
	{391, "CharacterCreator", "角色创建", "account", "Cycle hair, eyes, outfit and colors; the hero changes as you do.", demoCharacterCreator},
	{392, "LoginForm", "登录表单", "account", "Try player1 and hunter2. Anything else, and the form shakes its head.", demoLoginForm},
	{393, "PasswordGrid", "密码网格", "account", "Spell a password from the grid, as cartridges once asked. Try JUSTINBAILEY.", demoPasswordGrid},
	{394, "SessionList", "会话设备", "account", "Every device still signed in, and a way to kick each out.", demoSessionList},
	{395, "LanguagePicker", "语言选择", "account", "Languages by their own names, flags beside. ▶ marks yours.", demoLanguagePicker},
}

type accountDemo struct {
	music, sfx, fullscreen, scanlines bool
	speed, lang, failures             int
	binds                             []account.Binding
	look                              account.Look
	name, password, grid              string
	sessions                          []account.Session
}

func accountState(c *ui.Context) *accountDemo {
	now := pixel.K(c).Now()
	return ui.Local(c.Root(), demoKey("account"), func() accountDemo {
		return accountDemo{
			music: true, sfx: true, scanlines: true, speed: 1,
			binds: []account.Binding{
				{Action: "Jump", Key: ui.KeySpace}, {Action: "Attack", Key: ui.KeyZ}, {Action: "Dash", Key: ui.KeyZ}, {Action: "Map", Key: ui.KeyM},
			},
			sessions: []account.Session{
				{Device: "This laptop", Icon: icon.Globe, Seen: now, Current: true},
				{Device: "Game console", Icon: icon.Gamepad, Seen: now.Add(-2 * time.Hour)},
				{Device: "Old phone", Icon: icon.Phone, Seen: now.Add(-72 * time.Hour)},
			},
		}
	})
}

var speeds = []string{"SLOW", "NORMAL", "FAST", "INSTANT"}

func demoSettingsList(c *ui.Context, log func(string)) {
	s := accountState(c)
	toggle := func(label string, on *bool) func() {
		return func() {
			// The row shows the label; the switch says it aloud.
			if input.Switch(c, "", on, input.SwitchOptions{Label: label}).Changed {
				log(fmt.Sprintf("%s %v", label, *on))
			}
		}
	}
	account.SettingsList(c, "Settings", []account.SettingsSection{
		{Title: "AUDIO", Settings: []account.Setting{{Label: "Music", Control: toggle("Music", &s.music)}, {Label: "Sound effects", Control: toggle("Sound effects", &s.sfx)}}},
		{Title: "VIDEO", Settings: []account.Setting{{Label: "Full screen", Control: toggle("Full screen", &s.fullscreen)}, {Label: "Scanlines", Control: toggle("Scanlines", &s.scanlines)}}},
		{Title: "GAME", Settings: []account.Setting{{Label: "Text speed", Control: func() {
			if account.OptionCycler(c, "Text speed", speeds, &s.speed, account.OptionCyclerOptions{}).Changed {
				log("text speed " + speeds[s.speed])
			}
		}}}},
	}, account.SettingsListOptions{})
}

func demoOptionCycler(c *ui.Context, log func(string)) {
	s := accountState(c)
	if account.OptionCycler(c, "Text speed", speeds, &s.speed, account.OptionCyclerOptions{}).Changed {
		log("text speed " + speeds[s.speed])
	}
}

func demoKeyBindings(c *ui.Context, log func(string)) {
	s := accountState(c)
	if account.KeyBindings(c, "Controls", s.binds, account.KeyBindingsOptions{}).Changed {
		log("rebound")
	}
}

func demoControllerMap(c *ui.Context, log func(string)) {
	account.ControllerMap(c, "Pad", map[account.PadButton]string{
		account.DPad: "MOVE", account.ButtonA: "JUMP", account.ButtonB: "RUN", account.ButtonX: "MAP", account.ButtonY: "ITEM",
		account.ShoulderL: "PREV", account.ShoulderR: "NEXT", account.Start: "PAUSE", account.Select: "MENU",
	}, account.ControllerMapOptions{})
}

func demoProfileCard(c *ui.Context, log func(string)) {
	account.ProfileCard(c, account.Profile{Name: "Aria", Title: "Dragon Slayer", Level: 12, Stats: []account.Stat{
		{Name: "HP", Value: 30, Max: 40}, {Name: "MP", Value: 12, Max: 20}, {Name: "EXP", Value: 640, Max: 1000},
	}}, account.ProfileCardOptions{})
}

func demoCharacterCreator(c *ui.Context, log func(string)) {
	s := accountState(c)
	if account.CharacterCreator(c, "New hero", &s.look, account.CharacterCreatorOptions{}).Changed {
		log(fmt.Sprintf("look %+v", s.look))
	}
}

func demoLoginForm(c *ui.Context, log func(string)) {
	s := accountState(c)
	if account.LoginForm(c, "WELCOME BACK", &s.name, &s.password, account.LoginFormOptions{Failures: s.failures}).Submitted {
		// The demo checks; a real game asks its server.
		if s.name == "player1" && s.password == "hunter2" {
			log("logged in")
		} else {
			s.failures++
			log("wrong name or password")
		}
	}
}

func demoPasswordGrid(c *ui.Context, log func(string)) {
	s := accountState(c)
	r := account.PasswordGrid(c, "Password", &s.grid, account.PasswordGridOptions{})
	if r.Changed {
		log(fmt.Sprintf("%d of 12", len(s.grid)))
	}
	if r.Submitted {
		if s.grid == "JUSTINBAILEY" {
			log("suit off. you know the one.")
		} else {
			log("wrong password")
		}
	}
}

func demoSessionList(c *ui.Context, log func(string)) {
	s := accountState(c)
	if r := account.SessionList(c, "Signed in", s.sessions, account.SessionListOptions{}); r.SignedOut >= 0 {
		log("signed out " + s.sessions[r.SignedOut].Device)
		s.sessions = append(s.sessions[:r.SignedOut:r.SignedOut], s.sessions[r.SignedOut+1:]...)
	}
}

var languages = []account.Language{
	{Name: "English", Flag: flag("llcccrrcccll cllccrrccllc ccllcrrcllcc rrrrrrrrrrrr rrrrrrrrrrrr ccllcrrcllcc cllccrrccllc llcccrrcccll")},
	{Name: "中文", Flag: flag("rryrrrrrrrrr ryyyrrrrrrrr rryrryrrrrrr rrrrrrrrrrrr rrryrrrrrrrr rrrrrrrrrrrr rrrrrrrrrrrr rrrrrrrrrrrr")},
	{Name: "日本語", Flag: flag("llllllllllll llllllllllll lllllrrlllll llllrrrrllll llllrrrrllll lllllrrlllll llllllllllll llllllllllll")},
	{Name: "Français", Flag: flag("ccccllllrrrr ccccllllrrrr ccccllllrrrr ccccllllrrrr ccccllllrrrr ccccllllrrrr ccccllllrrrr ccccllllrrrr")},
	{Name: "Deutsch", Flag: flag("kkkkkkkkkkkk kkkkkkkkkkkk kkkkkkkkkkkk rrrrrrrrrrrr rrrrrrrrrrrr rrrrrrrrrrrr yyyyyyyyyyyy yyyyyyyyyyyy")},
	{Name: "Español", Flag: flag("rrrrrrrrrrrr rrrrrrrrrrrr yyyyyyyyyyyy yyyyyyyyyyyy yyyyyyyyyyyy yyyyyyyyyyyy rrrrrrrrrrrr rrrrrrrrrrrr")},
}

// flag is string art given as rows split by spaces.
func flag(rows string) *pixel.Sprite {
	return pixel.MustSprite(strings.ReplaceAll(rows, " ", "\n"))
}

func demoLanguagePicker(c *ui.Context, log func(string)) {
	s := accountState(c)
	if account.LanguagePicker(c, "Language", languages, &s.lang, account.LanguagePickerOptions{Height: 100}).Changed {
		log("language " + languages[s.lang].Name)
	}
}
