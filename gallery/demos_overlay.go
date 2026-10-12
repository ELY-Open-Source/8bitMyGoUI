package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

var overlayEntries = []Entry{
	{176, "Dialog", "对话框", "overlay", "A titled window over a checkered backdrop. Escape or a click outside closes it.", demoDialog},
	{177, "AlertDialog", "警告框", "overlay", "Icon, title, message, buttons. Enter presses the highlighted one.", demoAlertDialog},
	{178, "ConfirmDialog", "确认框", "overlay", "ARE YOU SURE? The hand starts on NO, to be safe.", demoConfirmDialog},
	{179, "Popover", "弹出层", "overlay", "A panel under its button, notched toward it. It flips above near the bottom.", demoPopover},
	{180, "Tooltip", "提示气泡", "overlay", "Rest the pointer half a second, or Tab in, and a bubble says what it does.", demoTooltip},
	{181, "Drawer", "抽屉", "overlay", "A panel that slides in from an edge, four steps at a time.", demoDrawer},
	{182, "Sheet", "底部面板", "overlay", "Rises from the bottom. Drag its bar down to put it away.", demoSheet},
	{183, "HoverCard", "悬停卡片", "overlay", "Point at a name to see the card; the card stays while you read it.", demoHoverCard},
	{184, "Lightbox", "灯箱", "overlay", "Pictures over a dark veil. ◀ ▶ or the arrows turn them.", demoLightbox},
	{185, "DialogueBox", "剧情对话", "overlay", "An RPG text box. A press shows the line, the next turns the page.", demoDialogueBox},
}

type overlayDemo struct {
	dialog, alert, confirm, popover, sheet, lightbox, choice, pause, prompt, saving bool
	drawer                                                                          display.Side
	drawerOpen                                                                      bool
	picture, page, step, cue                                                        int
	name                                                                            string
	notes                                                                           []overlay.Notification
	spot                                                                            bool
	since                                                                           time.Time
}

func overlayState(c *ui.Context) *overlayDemo {
	return ui.Local(c.Root(), demoKey("overlay"), func() overlayDemo { return overlayDemo{step: 1} })
}

func demoDialog(c *ui.Context, log func(string)) {
	s := overlayState(c)
	if input.Button(c, "OPEN QUEST", input.ButtonOptions{}).Clicked {
		s.dialog = true
		log("open")
	}
	r := overlay.Dialog(c, "Quest", &s.dialog, overlay.DialogOptions{Actions: []string{"DECLINE", "ACCEPT"}}, func() {
		text.Paragraph(c, "Bring three slime jellies to the baker. Reward: one warm pie.", text.ParagraphOptions{})
	})
	if r.Action >= 0 {
		log([]string{"declined", "accepted"}[r.Action])
	}
	if r.Closed {
		log("closed")
	}
}

func demoAlertDialog(c *ui.Context, log func(string)) {
	s := overlayState(c)
	if input.Button(c, "FIGHT BOSS", input.ButtonOptions{Variant: input.Danger}).Clicked {
		s.alert = true
		log("open")
	}
	buttons := []string{"NO", "HEAL"}
	r := overlay.AlertDialog(c, "Low HP", "You have 3 HP. Heal before the boss?", &s.alert, overlay.AlertDialogOptions{Buttons: buttons, Default: 1})
	if r.Pressed >= 0 {
		log(buttons[r.Pressed])
	}
	if r.Closed {
		log("closed")
	}
}

func demoConfirmDialog(c *ui.Context, log func(string)) {
	s := overlayState(c)
	if input.Button(c, "DROP SWORD", input.ButtonOptions{}).Clicked {
		s.confirm = true
		log("open")
	}
	r := overlay.ConfirmDialog(c, "Drop the legendary sword in the lava?", &s.confirm, overlay.ConfirmDialogOptions{})
	if r.Confirmed {
		log("dropped. it is gone.")
	}
	if r.Canceled {
		log("kept the sword")
	}
}

func demoPopover(c *ui.Context, log func(string)) {
	s := overlayState(c)
	b := input.Button(c, "MAP", input.ButtonOptions{Icon: icon.Map})
	if b.Clicked {
		s.popover = !s.popover
		log(fmt.Sprintf("open %v", s.popover))
	}
	overlay.Popover(c, "Map", b.Element, &s.popover, overlay.PopoverOptions{}, func() {
		pixel.Text(c, "You are here: Town.", pixel.TextOptions{})
		pixel.Text(c, "Cave: 2 screens west.", pixel.TextOptions{Slot: pixel.Faint})
	})
}

func demoTooltip(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(6)).Children(func() {
		for _, t := range []struct {
			art  *pixel.Sprite
			name string
			tip  string
		}{{icon.Sword, "Attack", "Swing for 12 damage."}, {icon.Shield, "Defend", "Halve damage this turn."}, {icon.Flame, "Magic", "Costs 4 MP."}} {
			overlay.Tooltip(c, t.tip, overlay.TooltipOptions{}, func() {
				if input.IconButton(c, t.art, t.name, input.IconButtonOptions{}).Clicked {
					log(t.name)
				}
			})
		}
	})
}

func demoDrawer(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := overlayState(c)
	ui.Row(c).Gap(k.D(6)).Children(func() {
		if input.Button(c, "◀ LEFT", input.ButtonOptions{}).Clicked {
			s.drawer, s.drawerOpen = display.Left, true
			log("open left")
		}
		if input.Button(c, "RIGHT ▶", input.ButtonOptions{}).Clicked {
			s.drawer, s.drawerOpen = display.Right, true
			log("open right")
		}
	})
	if overlay.Drawer(c, "Party", &s.drawerOpen, overlay.DrawerOptions{Side: s.drawer, Size: 120}, func() {
		pixel.Text(c, "PARTY", pixel.TextOptions{Slot: pixel.Faint})
		for _, h := range []string{"Aria  LV 7", "Brom  LV 5", "Cyra  LV 9"} {
			pixel.Text(c, h, pixel.TextOptions{})
		}
	}).Closed {
		log("drawer closed")
	}
}

func demoSheet(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := overlayState(c)
	if input.Button(c, "OPTIONS", input.ButtonOptions{Icon: icon.Settings}).Clicked {
		s.sheet = true
		log("open")
	}
	if overlay.Sheet(c, "Options", &s.sheet, overlay.SheetOptions{}, func() {
		ui.Column(c).Gap(k.D(4)).Children(func() {
			pixel.Text(c, "Music: ON", pixel.TextOptions{})
			pixel.Text(c, "Text speed: FAST", pixel.TextOptions{})
			pixel.Text(c, "Drag the bar down to close.", pixel.TextOptions{Slot: pixel.Faint})
		})
	}).Closed {
		log("sheet closed")
	}
}

func demoHoverCard(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	overlay.HoverCard(c, "Aria", overlay.HoverCardOptions{}, func() {
		pixel.Text(c, "@aria", pixel.TextOptions{Slot: pixel.Info})
	}, func() {
		ui.Row(c).Gap(k.D(6)).Children(func() {
			display.Avatar(c, "Aria", display.AvatarOptions{})
			ui.Column(c).Children(func() {
				pixel.Text(c, "Aria", pixel.TextOptions{})
				pixel.Text(c, "Knight · LV 7", pixel.TextOptions{Slot: pixel.Faint})
			})
		})
	})
}

var scenes = []struct {
	art  *pixel.Sprite
	name string
}{{icon.Home, "Town"}, {icon.Map, "Field"}, {icon.Skull, "Cave"}}

func demoLightbox(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := overlayState(c)
	ui.Row(c).Gap(k.D(4)).Children(func() {
		for i, sc := range scenes {
			if input.IconButton(c, sc.art, sc.name, input.IconButtonOptions{Size: input.Large}).Clicked {
				s.picture, s.lightbox = i, true
				log("open " + sc.name)
			}
		}
	})
	r := overlay.Lightbox(c, "Postcards", len(scenes), &s.picture, &s.lightbox, overlay.LightboxOptions{}, func(i int) {
		ui.Column(c).AlignItems(ui.Center).Padding(k.D(8)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Well)
		}).Children(func() {
			display.Icon(c, scenes[i].art, display.IconOptions{Scale: 4, Label: scenes[i].name})
			pixel.Text(c, scenes[i].name, pixel.TextOptions{})
		})
	})
	if r.Changed {
		log("showing " + scenes[s.picture].name)
	}
	if r.Closed {
		log("closed")
	}
}

func demoDialogueBox(c *ui.Context, log func(string)) {
	s := overlayState(c)
	pages := []string{"It is dangerous to go alone!", "Take this. It is a stick.", "Mostly harmless. Good luck."}
	r := overlay.DialogueBox(c, pages, &s.page, overlay.DialogueBoxOptions{Name: "OLD MAN", Portrait: icon.User})
	r.Element.Width(pixel.K(c).D(280))
	if r.Changed {
		log(fmt.Sprintf("page %d", s.page+1))
	}
	if r.Done {
		s.page = 0
		log("the old man starts over")
	}
}
