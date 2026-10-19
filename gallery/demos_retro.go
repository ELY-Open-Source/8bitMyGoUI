package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/desk"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var retroEntries = []Entry{
	{476, "ShutdownScreen", "关机画面", "desk", "The kindest words a computer ever said.", demoShutdownScreen},
	{477, "BootScreen", "开机自检", "desk", "Memory counts up, POST lines roll by, the logo waits.", demoBootScreen},
	{478, "Screensaver", "屏保", "desk", "A logo bounces round the screen. Any key, click or move wakes it. Watch the corners.", demoScreensaver},
	{479, "Snake", "贪吃蛇", "desk", "Arrows steer, food grows you. The wall is not your friend.", demoSnake},
	{480, "Dock", "程序坞", "desk", "A shelf of apps. The one you point at or focus rises.", demoDock},
	{481, "RunDialog", "运行", "desk", "Type a program's name. OK runs it; CANCEL or Escape walks away.", demoRunDialog},
	{482, "DiskDefrag", "磁盘整理", "desk", "Scattered blocks march into order. Oddly soothing.", demoDiskDefrag},
	{483, "WallpaperPicker", "壁纸选择", "desk", "Dither patterns to dress the desktop.", demoWallpaperPicker},
	{484, "ClipboardPanel", "剪贴板", "desk", "Recent clips. Click one and it is yours to paste again.", demoClipboardPanel},
	{485, "LockScreen", "锁屏", "desk", "The time, the date, and PRESS ANY KEY.", demoLockScreen},
}

type retroDemo struct {
	saver, locked, defrag bool
	boots, disks          int
	run                   string
	paper                 desk.Wallpaper
}

func retroState(c *ui.Context) *retroDemo {
	return ui.Local(c.Root(), demoKey("retro"), func() retroDemo { return retroDemo{run: "notepad"} })
}

func demoShutdownScreen(c *ui.Context, log func(string)) {
	desk.ShutdownScreen(c, desk.ShutdownScreenOptions{Width: 300, Height: 120})
}

func demoBootScreen(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := retroState(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		// A new key is a new screen: REBOOT starts over.
		ui.Box(c).Key(s.boots).Children(func() {
			if desk.BootScreen(c, "8BIT BIOS V4.51", []string{
				"CPU: PIXEL 8086 AT 4.77 MHZ", "FLOPPY A: 1.44 MB", "KEYBOARD... OK", "BOOTING FROM C:",
			}, desk.BootScreenOptions{Width: 300, Height: 140}).Done {
				log("booted")
			}
		})
		if input.Button(c, "REBOOT", input.ButtonOptions{Size: input.Small}).Clicked {
			s.boots++
			log("rebooting")
		}
	})
}

func demoScreensaver(c *ui.Context, log func(string)) {
	s := retroState(c)
	if input.Button(c, "START SAVER", input.ButtonOptions{}).Clicked {
		s.saver = true
		log("saver on")
	}
	if desk.Screensaver(c, "8BIT", &s.saver, desk.ScreensaverOptions{}).Woke {
		log("woke up")
	}
}

func demoSnake(c *ui.Context, log func(string)) {
	r := desk.Snake(c, "Snake pit", desk.SnakeOptions{})
	if r.Started {
		log("go!")
	}
	if r.Over {
		log(fmt.Sprintf("game over at %d", r.Score))
	}
}

func demoDock(c *ui.Context, log func(string)) {
	apps := []desk.DockApp{
		{Name: "FILES", Icon: icon.Folder}, {Name: "MAIL", Icon: icon.Mail}, {Name: "MUSIC", Icon: icon.Music},
		{Name: "CAMERA", Icon: icon.Camera}, {Name: "TERMINAL", Icon: icon.Terminal}, {Name: "TRASH", Icon: icon.Trash},
	}
	if i := desk.Dock(c, "Dock", apps, desk.DockOptions{}).Launched; i >= 0 {
		log("launched " + apps[i].Name)
	}
}

func demoRunDialog(c *ui.Context, log func(string)) {
	s := retroState(c)
	r := desk.RunDialog(c, "Run", &s.run, desk.RunDialogOptions{})
	if r.Changed {
		log(fmt.Sprintf("%d letters typed", len([]rune(s.run))))
	}
	if r.Submitted {
		log(fmt.Sprintf("running a %d-letter program", len([]rune(s.run))))
	}
	if r.Canceled {
		log("canceled")
	}
}

func demoDiskDefrag(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := retroState(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		ui.Box(c).Key(s.disks).Children(func() {
			if desk.DiskDefrag(c, "Drive C:", s.defrag, desk.DiskDefragOptions{Seed: 7}).Done {
				s.defrag = false
				log("drive in order")
			}
		})
		ui.Row(c).Gap(k.D(4)).Children(func() {
			label := "START"
			if s.defrag {
				label = "PAUSE"
			}
			if input.Button(c, label, input.ButtonOptions{Size: input.Small}).Clicked {
				s.defrag = !s.defrag
				log(fmt.Sprintf("defrag running %v", s.defrag))
			}
			if input.Button(c, "NEW DISK", input.ButtonOptions{Size: input.Small}).Clicked {
				s.disks, s.defrag = s.disks+1, false
				log("a fresh mess")
			}
		})
	})
}

func demoWallpaperPicker(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := retroState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Start).Children(func() {
		if desk.WallpaperPicker(c, "Wallpaper", &s.paper, desk.WallpaperPickerOptions{}).Changed {
			log(fmt.Sprintf("wallpaper %d", s.paper))
		}
		desk.Desktop(c, "Preview", desk.DesktopOptions{Wallpaper: s.paper, Width: 120, Height: 90}, func() {})
	})
}

func demoClipboardPanel(c *ui.Context, log func(string)) {
	clips := []string{"git push gitea main", "IT'S NOW SAFE TO TURN OFF", "UP UP DOWN DOWN LEFT RIGHT", "8bitmygoui.workers.dev"}
	if i := desk.ClipboardPanel(c, "Clipboard", clips, desk.ClipboardPanelOptions{Width: 240}).Copied; i >= 0 {
		log(fmt.Sprintf("copied clip %d", i+1))
	}
}

func demoLockScreen(c *ui.Context, log func(string)) {
	s := retroState(c)
	if input.Button(c, "LOCK", input.ButtonOptions{}).Clicked {
		s.locked = true
		log("locked")
	}
	if desk.LockScreen(c, "Lock screen", &s.locked, desk.LockScreenOptions{}).Unlocked {
		log("unlocked")
	}
}
