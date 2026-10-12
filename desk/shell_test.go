package desk_test

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/desk"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestDesktopFillsColumnsFirst(t *testing.T) {
	tt := view(700, 400, newWorld(), func(c *ui.Context) {
		desk.Desktop(c, "Desk", desk.DesktopOptions{Height: 140}, func() {
			for _, name := range []string{"A", "B", "C"} {
				desk.DesktopIcon(c, name, icon.File, desk.DesktopIconOptions{})
			}
		})
	})
	a, b, cc := find(t, tt, "A"), find(t, tt, "B"), find(t, tt, "C")
	// Two icons fill a 140-dot column; the third starts another.
	if b.X != a.X || b.Y <= a.Y || cc.X <= a.X || cc.Y != a.Y {
		t.Fatalf("icons at %v, %v, %v", a, b, cc)
	}
}

func TestDesktopIconChoosesOpensAndInverts(t *testing.T) {
	chosen, clicks, opens, off := false, 0, 0, false
	tt := view(300, 200, newWorld(), func(c *ui.Context) {
		r := desk.DesktopIcon(c, "NOTES", icon.FileText, desk.DesktopIconOptions{Chosen: chosen, Disabled: off})
		if r.Clicked {
			clicks++
			chosen = true
		}
		if r.Opened {
			opens++
		}
	})
	x, y, w, h := dots(find(t, tt, "NOTES"))
	if count(tt.Image(), x, y, w, h, pixel.Accent) > 0 {
		t.Fatal("an unchosen icon shows Accent")
	}
	click(t, tt, "NOTES")
	if clicks != 1 || opens != 0 || count(tt.Image(), x, y, w, h, pixel.Accent) == 0 {
		t.Fatalf("one click: %d clicks, %d opens, or it did not invert", clicks, opens)
	}
	// The second quick click is a double click: it opens.
	click(t, tt, "NOTES")
	tt.Key(0, ui.KeyEnter)
	if opens != 2 {
		t.Fatalf("double click and Enter opened %d times", opens)
	}
	// A key sent as it turns off still does nothing.
	off = true
	tt.Key(0, ui.KeyEnter)
	if opens != 2 {
		t.Fatal("a disabled icon opened")
	}
}

func TestTaskbarSwitchesMinimizesAndTellsTheTime(t *testing.T) {
	active, launched, trayed := 0, "", -1
	wd := newWorld()
	progs := []desk.Program{{Name: "GAMES", Icon: icon.Folder, Programs: []desk.Program{{Name: "SNAKE", Icon: icon.Bolt}}}, {Name: "RUN", Icon: icon.Terminal}}
	tt := view(700, 300, wd, func(c *ui.Context) {
		r := desk.Taskbar(c, "Bar", []string{"ONE", "TWO"}, &active, desk.TaskbarOptions{Programs: progs, Tray: []desk.TrayIcon{{Icon: icon.Wifi, Tip: "ONLINE"}}})
		if r.Launched != "" {
			launched = r.Launched
		}
		if r.Tray >= 0 {
			trayed = r.Tray
		}
	})
	tt.Frame()
	find(t, tt, "10:30")
	x, y, w, h := dots(find(t, tt, "OPEN WINDOWS"))
	// The second button chooses TWO; pressing it again minimizes.
	tt.ClickAt(float32(2*(x+w*3/4)), float32(2*(y+h/2)))
	if active != 1 {
		t.Fatalf("active %d after pressing TWO", active)
	}
	tt.Key(0, ui.KeyEnter)
	if active != -1 {
		t.Fatalf("pressing the active window left %d", active)
	}
	click(t, tt, "ONLINE")
	if trayed != 0 {
		t.Fatalf("tray click gave %d", trayed)
	}
	click(t, tt, "START")
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyEnter)
	if launched != "GAMES/SNAKE" {
		t.Fatalf("START, Enter, Enter launched %q", launched)
	}
}

func TestStartMenuCascades(t *testing.T) {
	launched, off := "", false
	progs := []desk.Program{
		{Name: "PROGRAMS", Icon: icon.Folder, Programs: []desk.Program{
			{Name: "GAMES", Icon: icon.Folder, Programs: []desk.Program{{Name: "SNAKE", Icon: icon.Bolt}, {Name: "MINES", Icon: icon.Flag}}},
			{Name: "NOTEPAD", Icon: icon.FileText},
		}},
		{Name: "HELP", Icon: icon.Question},
	}
	tt := view(700, 500, newWorld(), func(c *ui.Context) {
		if r := desk.StartMenu(c, "Start", progs, desk.StartMenuOptions{Disabled: off}); r.Launched != "" {
			launched = r.Launched
		}
	})
	click(t, tt, "Start")
	if _, ok := tt.Find("PROGRAMS"); !ok {
		t.Fatal("START did not open")
	}
	// Right goes in, Left comes back out, Enter runs.
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if _, ok := tt.Find("MINES"); !ok {
		t.Fatalf("two Rights did not reach GAMES: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if launched != "PROGRAMS/NOTEPAD" {
		t.Fatalf("launched %q", launched)
	}
	if _, ok := tt.Find("PROGRAMS"); ok {
		t.Fatal("the menu stayed open after a launch")
	}
	// Clicks: a submenu opens, a program runs.
	click(t, tt, "Start")
	click(t, tt, "PROGRAMS")
	click(t, tt, "GAMES")
	click(t, tt, "MINES")
	if launched != "PROGRAMS/GAMES/MINES" {
		t.Fatalf("clicks launched %q", launched)
	}
	// Turned off while open, it closes.
	click(t, tt, "Start")
	off = true
	tt.Frame()
	if _, ok := tt.Find("PROGRAMS"); ok {
		t.Fatal("a START turned off stayed open")
	}
	click(t, tt, "Start")
	if _, ok := tt.Find("PROGRAMS"); ok {
		t.Fatal("a disabled START opened")
	}
}

func TestStartMenuOpensSubmenusUnderThePointer(t *testing.T) {
	progs := []desk.Program{{Name: "TOOLS", Icon: icon.Folder, Programs: []desk.Program{{Name: "HAMMER", Icon: icon.Sword}}}, {Name: "QUIT", Icon: icon.Power}}
	tt := view(700, 500, newWorld(), func(c *ui.Context) { desk.StartMenu(c, "Start", progs, desk.StartMenuOptions{}) })
	click(t, tt, "Start")
	tt.Key(0, ui.KeyDown)
	r := find(t, tt, "TOOLS")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if _, ok := tt.Find("HAMMER"); !ok {
		t.Fatalf("pointing at TOOLS did not open it: %q", tt.Texts())
	}
}

func TestTrayIconsAreOneTabStop(t *testing.T) {
	clicked := -1
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		if i := desk.TrayIcons(c, "Tray", []desk.TrayIcon{{Icon: icon.Wifi, Tip: "ONLINE"}, {Icon: icon.Volume, Tip: "LOUD"}}, desk.TrayIconsOptions{}).Clicked; i >= 0 {
			clicked = i
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if clicked != 1 || !tt.Focused("LOUD") {
		t.Fatalf("Tab, Right, Enter clicked %d", clicked)
	}
	// Focused, the tip shows in a bubble too.
	if n := strings.Count(strings.Join(tt.Texts(), "|"), "LOUD"); n < 2 {
		t.Fatalf("LOUD shows %d times: %q", n, tt.Texts())
	}
}

func TestDockRaisesAndLaunches(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = reduce
		launched := -1
		tt := view(400, 200, wd, func(c *ui.Context) {
			if i := desk.Dock(c, "Dock", []desk.DockApp{{Name: "MAIL", Icon: icon.Mail}, {Name: "MUSIC", Icon: icon.Music}}, desk.DockOptions{}).Launched; i >= 0 {
				launched = i
			}
		})
		x, y, _, _ := dots(find(t, tt, "MAIL"))
		top := func() int {
			for j := y; j < y+40; j++ {
				if count(tt.Image(), x, j, 30, 1, pixel.Ink) > 0 {
					return j - y
				}
			}
			return -1
		}
		rest := top()
		tt.Key(0, ui.KeyTab)
		first := top()
		later(tt, wd, time.Second)
		up := top()
		// Focused, it rises eight dots: stepping, or at once reduced.
		if up != rest-8 || reduce != (first == up) {
			t.Errorf("reduce %v: rest %d, first %d, risen %d", reduce, rest, first, up)
		}
		tt.Key(0, ui.KeyEnter)
		if launched != 0 {
			t.Errorf("Enter launched %d", launched)
		}
	}
}

func TestControlPanelOpensByKeysAndClicks(t *testing.T) {
	opened := -1
	applets := []desk.Applet{{Name: "SOUND", Icon: icon.Volume}, {Name: "CLOCK", Icon: icon.Clock}, {Name: "POWER", Icon: icon.Power}}
	tt := view(500, 200, newWorld(), func(c *ui.Context) {
		if i := desk.ControlPanel(c, "Panel", applets, desk.ControlPanelOptions{}).Opened; i >= 0 {
			opened = i
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	if opened != 2 {
		t.Fatalf("End, Enter opened %d", opened)
	}
	r := find(t, tt, "Panel")
	tt.ClickAt(r.X+2*33, r.Y+2*20)
	if opened != 0 {
		t.Fatalf("a click on SOUND opened %d", opened)
	}
}

func TestWallpaperPickerChooses(t *testing.T) {
	paper, changed := desk.Checks, 0
	tt := view(500, 300, newWorld(), func(c *ui.Context) {
		if desk.WallpaperPicker(c, "Paper", &paper, desk.WallpaperPickerOptions{}).Changed {
			changed++
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if paper != desk.Dots || changed != 1 {
		t.Fatalf("Down, Enter chose %d, %d changes", paper, changed)
	}
	find(t, tt, "DOTS")
}

func TestStartMenuHoldsStillAsSubmenusOpen(t *testing.T) {
	progs := []desk.Program{
		{Name: "TOOLS", Icon: icon.Folder, Programs: []desk.Program{{Name: "SAW", Icon: icon.Sword}, {Name: "AXE", Icon: icon.Sword}, {Name: "HOE", Icon: icon.Leaf}}},
		{Name: "QUIT", Icon: icon.Power},
	}
	for _, low := range []bool{false, true} {
		launched := ""
		tt := view(700, 500, newWorld(), func(c *ui.Context) {
			// Low on the window, the menu opens upward.
			ui.Box(c).Height(map[bool]float32{false: 0, true: 420}[low])
			if r := desk.StartMenu(c, "Start", progs, desk.StartMenuOptions{}); r.Launched != "" {
				launched = r.Launched
			}
		})
		click(t, tt, "Start")
		before := find(t, tt, "TOOLS")
		click(t, tt, "TOOLS")
		if after := find(t, tt, "TOOLS"); after != before {
			t.Errorf("low %v: TOOLS moved from %v to %v", low, before, after)
		}
		click(t, tt, "AXE")
		if launched != "TOOLS/AXE" {
			t.Errorf("low %v: launched %q", low, launched)
		}
	}
}
