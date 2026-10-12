package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/desk"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var deskEntries = []Entry{
	{466, "Desktop", "桌面", "desk", "A dithered wallpaper with icons down the left, as desktops had.", demoDesktop},
	{467, "DesktopIcon", "桌面图标", "desk", "Click to choose and it inverts; double-click or Enter opens it.", demoDesktopIcon},
	{468, "Taskbar", "任务栏", "desk", "START, a button for each open window, the tray and the clock.", demoTaskbar},
	{469, "StartMenu", "开始菜单", "desk", "Programs in cascading menus. Arrows walk in and out; Enter runs one.", demoStartMenu},
	{470, "AppWindow", "应用框体", "desk", "Drag a window by its title bar, or focus the bar and use the arrows.", demoAppWindow},
	{471, "Calculator", "计算器", "desk", "It really adds up. Divide by zero for a GAME OVER; C continues.", demoCalculator},
	{472, "Notepad", "记事本", "desk", "A text area under FILE and EDIT. TIME/DATE stamps the moment.", demoNotepad},
	{473, "Minesweeper", "扫雷", "desk", "Open with a click or Enter, flag with a right click or F. The face starts over.", demoMinesweeper},
	{474, "ControlPanel", "控制面板", "desk", "Settings as an icon grid. Arrows choose; Enter opens.", demoControlPanel},
	{475, "TrayIcons", "托盘图标", "desk", "Small icons by the clock. Hover or focus one for its tip.", demoTrayIcons},
}

type deskDemo struct {
	chosen, active, wx, wy int
	note                   string
	shut                   bool
}

func deskState(c *ui.Context) *deskDemo {
	return ui.Local(c.Root(), demoKey("desk"), func() deskDemo {
		return deskDemo{chosen: -1, wx: 70, wy: 12, note: "Dear diary,\ntoday I defragged."}
	})
}

// icons sit on the demo desktops.
var icons = []struct {
	name string
	sp   *pixel.Sprite
}{{"MY COMPUTER", icon.Gamepad}, {"RECYCLE BIN", icon.Trash}, {"NOTES.TXT", icon.FileText}, {"GAMES", icon.Folder}}

// programs fill the demo START menus.
var programs = []desk.Program{
	{Name: "PROGRAMS", Icon: icon.Folder, Programs: []desk.Program{
		{Name: "ACCESSORIES", Icon: icon.Folder, Programs: []desk.Program{{Name: "NOTEPAD", Icon: icon.FileText}, {Name: "CALCULATOR", Icon: icon.Grid}}},
		{Name: "GAMES", Icon: icon.Gamepad, Programs: []desk.Program{{Name: "MINESWEEPER", Icon: icon.Flag}, {Name: "SNAKE", Icon: icon.Bolt}}},
	}},
	{Name: "DOCUMENTS", Icon: icon.FileText},
	{Name: "SETTINGS", Icon: icon.Settings, Programs: []desk.Program{{Name: "CONTROL PANEL", Icon: icon.Settings}, {Name: "WALLPAPER", Icon: icon.Image}}},
	{Name: "RUN...", Icon: icon.Terminal},
	{Name: "SHUT DOWN", Icon: icon.Power},
}

func demoDesktop(c *ui.Context, log func(string)) {
	s := deskState(c)
	desk.Desktop(c, "Desktop", desk.DesktopOptions{Width: 300, Height: 150}, func() {
		for i, it := range icons {
			r := desk.DesktopIcon(c, it.name, it.sp, desk.DesktopIconOptions{Chosen: s.chosen == i})
			if r.Clicked {
				s.chosen = i
				log("chose " + it.name)
			}
			if r.Opened {
				log("opened " + it.name)
			}
		}
	})
}

func demoDesktopIcon(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := deskState(c)
	ui.Row(c).Gap(k.D(4)).Children(func() {
		for i, it := range icons[:3] {
			r := desk.DesktopIcon(c, it.name, it.sp, desk.DesktopIconOptions{Chosen: s.chosen == i})
			if r.Clicked {
				s.chosen = i
				log("chose " + it.name)
			}
			if r.Opened {
				log("opened " + it.name)
			}
		}
	})
}

// tray sits by the demo taskbars' clocks.
var tray = []desk.TrayIcon{{Icon: icon.Wifi, Tip: "ONLINE"}, {Icon: icon.Volume, Tip: "VOLUME 80%"}, {Icon: icon.Battery, Tip: "BATTERY 99%"}}

func demoTaskbar(c *ui.Context, log func(string)) {
	s := deskState(c)
	windows := []string{"NOTEPAD", "CALC", "MINES"}
	r := desk.Taskbar(c, "Taskbar", windows, &s.active, desk.TaskbarOptions{Programs: programs, Tray: tray})
	if r.Changed {
		log(fmt.Sprintf("active window %d", s.active))
	}
	if r.Launched != "" {
		log("launched " + r.Launched)
	}
	if r.Tray >= 0 {
		log("tray: " + tray[r.Tray].Tip)
	}
}

func demoStartMenu(c *ui.Context, log func(string)) {
	if r := desk.StartMenu(c, "Start", programs, desk.StartMenuOptions{}); r.Launched != "" {
		log("launched " + r.Launched)
	}
}

func demoAppWindow(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := deskState(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		desk.Desktop(c, "Desktop", desk.DesktopOptions{Wallpaper: desk.Dots, Width: 320, Height: 170}, func() {
			if !s.shut {
				r := desk.AppWindow(c, "HELLO.TXT", &s.wx, &s.wy, desk.AppWindowOptions{Width: 120, Height: 48, Room: [2]int{320, 170}}, func() {
					pixel.Text(c, "Drag me around!", pixel.TextOptions{})
				})
				if r.Moved {
					log(fmt.Sprintf("window at %d, %d", s.wx, s.wy))
				}
				if r.Closed {
					s.shut = true
					log("window closed")
				}
			}
		})
		if s.shut && input.Button(c, "REOPEN", input.ButtonOptions{Size: input.Small}).Clicked {
			s.shut = false
			log("window reopened")
		}
	})
}

func demoCalculator(c *ui.Context, log func(string)) {
	if desk.Calculator(c, "Pocket calculator", desk.CalculatorOptions{}).Changed {
		log("key pressed")
	}
}

func demoNotepad(c *ui.Context, log func(string)) {
	s := deskState(c)
	r := desk.Notepad(c, "Notepad", &s.note, desk.NotepadOptions{})
	if r.Changed {
		log(fmt.Sprintf("%d runes", len([]rune(s.note))))
	}
	if r.Saved {
		log("saved")
	}
}

func demoMinesweeper(c *ui.Context, log func(string)) {
	r := desk.Minesweeper(c, "Minesweeper", desk.MinesweeperOptions{})
	if r.Flagged {
		log("flag")
	}
	if r.Restarted {
		log("new game")
	}
	if r.Won {
		log("cleared the field")
	}
	if r.Lost {
		log("boom")
	}
}

func demoControlPanel(c *ui.Context, log func(string)) {
	applets := []desk.Applet{
		{Name: "DISPLAY", Icon: icon.Image}, {Name: "SOUND", Icon: icon.Volume}, {Name: "CLOCK", Icon: icon.Clock}, {Name: "NETWORK", Icon: icon.Wifi},
		{Name: "POWER", Icon: icon.Power}, {Name: "USERS", Icon: icon.Users}, {Name: "FONTS", Icon: icon.FileText}, {Name: "GAMEPAD", Icon: icon.Gamepad},
	}
	if i := desk.ControlPanel(c, "Control panel", applets, desk.ControlPanelOptions{}).Opened; i >= 0 {
		log("opened " + applets[i].Name)
	}
}

func demoTrayIcons(c *ui.Context, log func(string)) {
	if i := desk.TrayIcons(c, "Tray", tray, desk.TrayIconsOptions{}).Clicked; i >= 0 {
		log("clicked " + tray[i].Tip)
	}
}
