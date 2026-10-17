package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var popupEntries = []Entry{
	{186, "ChoicePrompt", "选项提示", "overlay", "A question and a hand. Arrows and Enter, or press 1–9.", demoChoicePrompt},
	{187, "Coachmark", "新手引导", "overlay", "A guided tour: a bubble points at each button in turn.", demoCoachmark},
	{188, "Spotlight", "探照灯", "overlay", "Everything but one button sinks into dither.", demoSpotlight},
	{189, "Notification", "通知", "overlay", "Cards stack at the top right and leave on their own. × sends one away.", demoNotification},
	{190, "Popconfirm", "气泡确认", "overlay", "A small YES / NO bubble before the button acts.", demoPopconfirm},
	{191, "PauseMenu", "暂停菜单", "overlay", "PAUSED blinks over the game. Escape resumes.", demoPauseMenu},
	{192, "PromptDialog", "输入对话框", "overlay", "ENTER YOUR NAME, eight letters at most.", demoPromptDialog},
	{193, "ProgressDialog", "进度对话框", "overlay", "Saving, with a way out.", demoProgressDialog},
	{194, "KeyHints", "按键提示", "overlay", "Pad buttons and keycaps beside what they do.", demoKeyHints},
	{195, "TitleCard", "关卡标题", "overlay", "WORLD 1-1 dithers in, waits, and dithers out.", demoTitleCard},
}

func demoChoicePrompt(c *ui.Context, log func(string)) {
	s := overlayState(c)
	choices := []string{"Fight", "Talk", "Flee"}
	if input.Button(c, "A SLIME APPEARS", input.ButtonOptions{}).Clicked {
		s.choice = true
		log("open")
	}
	r := overlay.ChoicePrompt(c, "A wild slime blocks the road. It wobbles.", choices, &s.choice, overlay.ChoicePromptOptions{})
	if r.Chosen >= 0 {
		log(choices[r.Chosen])
	}
	if r.Canceled {
		log("stared at the slime")
	}
}

func demoCoachmark(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := overlayState(c)
	tips := []string{"Your bag. Potions live here.", "The map. You are the red dot.", "Save often. Slimes are tricky."}
	var targets [3]ui.Element
	ui.Row(c).Gap(k.D(6)).Children(func() {
		for i, b := range []struct {
			art  *pixel.Sprite
			name string
		}{{icon.Chest, "Bag"}, {icon.Map, "Map"}, {icon.Save, "Save"}} {
			targets[i] = input.Button(c, b.name, input.ButtonOptions{Icon: b.art}).Element
		}
	})
	if s.step > len(tips) {
		if input.Button(c, "TOUR AGAIN", input.ButtonOptions{Size: input.Small}).Clicked {
			s.step = 1
			log("tour again")
		}
		return
	}
	if overlay.Coachmark(c, targets[s.step-1], tips[s.step-1], s.step, len(tips), overlay.CoachmarkOptions{}).Next {
		s.step++
		log(fmt.Sprintf("step %d", s.step))
	}
}

func demoSpotlight(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := overlayState(c)
	var shop ui.Element
	ui.Row(c).Gap(k.D(6)).Children(func() {
		input.Button(c, "INN", input.ButtonOptions{})
		shop = input.Button(c, "SHOP", input.ButtonOptions{Icon: icon.Cart}).Element
		input.Button(c, "GUILD", input.ButtonOptions{})
	})
	if input.Toggle(c, "SPOTLIGHT", &s.spot, input.ToggleOptions{Size: input.Small}).Changed {
		log(fmt.Sprintf("spotlight %v", s.spot))
	}
	if s.spot {
		overlay.Spotlight(c, shop, "Potions are on sale!", overlay.SpotlightOptions{})
	}
}

func demoNotification(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := overlayState(c)
	ui.Row(c).Gap(k.D(4)).Children(func() {
		if input.Button(c, "NEW MAIL", input.ButtonOptions{Size: input.Small}).Clicked {
			s.notes = append(s.notes, overlay.Notification{Icon: icon.Mail, Title: "Letter from Mom", Text: "Eat your vegetables."})
			log("mail")
		}
		if input.Button(c, "LEVEL UP", input.ButtonOptions{Size: input.Small}).Clicked {
			s.notes = append(s.notes, overlay.Notification{Icon: icon.Star, Title: "Level 8!", Text: "STR +2, LUCK +1.", For: 8 * time.Second})
			log("level up")
		}
	})
	if r := overlay.Notifications(c, &s.notes, overlay.NotificationsOptions{}); r.Dismissed != nil {
		log("dismissed " + r.Dismissed.Title)
	}
}

func demoPopconfirm(c *ui.Context, log func(string)) {
	r := overlay.Popconfirm(c, "SELL CROWN", "Sell it for 5 gold?", overlay.PopconfirmOptions{})
	if r.Confirmed {
		log("sold. what a deal.")
	}
	if r.Canceled {
		log("kept the crown")
	}
}

func demoPauseMenu(c *ui.Context, log func(string)) {
	s := overlayState(c)
	if input.Button(c, "PAUSE", input.ButtonOptions{Icon: icon.Pause}).Clicked {
		s.pause = true
		log("open")
	}
	if r := overlay.PauseMenu(c, &s.pause, overlay.PauseMenuOptions{}); r.Chosen >= 0 {
		log([]string{"resume", "options", "quit"}[r.Chosen])
	}
}

func demoPromptDialog(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := overlayState(c)
	ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		if input.Button(c, "NEW GAME", input.ButtonOptions{}).Clicked {
			s.prompt = true
			log("open")
		}
		if s.name != "" {
			pixel.Text(c, "Hero: "+s.name, pixel.TextOptions{})
		}
	})
	r := overlay.PromptDialog(c, "ENTER YOUR NAME", &s.name, &s.prompt, overlay.PromptDialogOptions{Placeholder: "LINK", MaxLength: 8})
	if r.Submitted {
		log("named " + s.name)
	}
	if r.Canceled {
		log("canceled")
	}
}

func demoProgressDialog(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := overlayState(c)
	if input.Button(c, "SAVE", input.ButtonOptions{Icon: icon.Save}).Clicked {
		s.saving, s.since = true, k.Now()
		log("saving")
	}
	// The pretend save takes three seconds.
	done := 0
	if s.saving {
		done = k.Since(c, s.since, 10, 30)
	}
	if done == 30 {
		s.saving = false
		log("saved")
	}
	if overlay.ProgressDialog(c, "Saving", float64(done)/30, &s.saving, overlay.ProgressDialogOptions{}).Canceled {
		log(fmt.Sprintf("canceled at %d%%", done*100/30))
	}
}

func demoKeyHints(c *ui.Context, log func(string)) {
	overlay.KeyHints(c, []overlay.KeyHint{{Key: "A", Action: "CONFIRM"}, {Key: "B", Action: "BACK"}, {Key: "START", Action: "MENU"}}, overlay.KeyHintsOptions{})
}

func demoTitleCard(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := overlayState(c)
	layout.Layers(c, layout.LayersOptions{}, func() {
		ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Dither(0, 0, g.W, g.H, pixel.Well, pixel.Info, 4)
			g.Rect(0, g.H-12, g.W, 12, pixel.Good)
		})
	}, func() {
		overlay.TitleCard(c, "WORLD 1-1", overlay.TitleCardOptions{Cue: s.cue})
	}).Size(k.D(200), k.D(90))
	if input.Button(c, "ENTER WORLD", input.ButtonOptions{Size: input.Small}).Clicked {
		s.cue++
		log("world 1-1")
	}
}
