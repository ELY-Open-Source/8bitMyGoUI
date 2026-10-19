package desk_test

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/desk"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
)

func TestStartMenuSurvivesShrinkingPrograms(t *testing.T) {
	progs := []desk.Program{{Name: "A", Icon: icon.File}, {Name: "B", Icon: icon.Folder, Programs: []desk.Program{{Name: "C", Icon: icon.File}}}}
	launched := ""
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		if r := desk.StartMenu(c, "Start", progs, desk.StartMenuOptions{}); r.Launched != "" {
			launched = r.Launched
		}
	})
	click(t, tt, "Start")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyRight)
	progs = progs[:1]
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if launched != "A" {
		t.Fatalf("after B left, Enter launched %q", launched)
	}
}

func TestDisabledAppWindowTurnsItsBodyOff(t *testing.T) {
	x, y, pressed := 0, 0, 0
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		desk.AppWindow(c, "WIN", &x, &y, desk.AppWindowOptions{Disabled: true}, func() {
			if input.Button(c, "PRESS", input.ButtonOptions{}).Clicked {
				pressed++
			}
		})
	})
	click(t, tt, "PRESS")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if pressed != 0 || tt.Focused("PRESS") {
		t.Fatalf("a disabled window's body took %d presses", pressed)
	}
}

func TestAppWindowKeepsItsWholeBodyInTheRoom(t *testing.T) {
	x, y := 0, 190
	tt := view(700, 500, newWorld(), func(c *ui.Context) {
		desk.AppWindow(c, "WIN", &x, &y, desk.AppWindowOptions{Height: 40, Room: [2]int{320, 200}}, func() {})
	})
	_, top, _, _ := dots(find(t, tt, "WIN"))
	if _, wy, _, wh := dots(find(t, tt, "WIN")); wy+wh > 200 || top != y {
		t.Fatalf("the window spans %d to %d, past 200", wy, wy+wh)
	}
}

func TestDesktopIsOneTabStop(t *testing.T) {
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		desk.Desktop(c, "Desk", desk.DesktopOptions{}, func() {
			desk.DesktopIcon(c, "A", icon.File, desk.DesktopIconOptions{})
			desk.DesktopIcon(c, "B", icon.File, desk.DesktopIconOptions{})
		})
		input.Button(c, "AFTER", input.ButtonOptions{})
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("B") {
		t.Fatal("Down did not move to B")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("AFTER") {
		t.Fatal("Tab stopped on another icon")
	}
}

func TestRowsKeepFocusAsSiblingsGo(t *testing.T) {
	apps := []desk.DockApp{{Name: "A", Icon: icon.File}, {Name: "B", Icon: icon.Mail}, {Name: "C", Icon: icon.Music}}
	tt := view(600, 300, newWorld(), func(c *ui.Context) { desk.Dock(c, "Dock", apps, desk.DockOptions{}) })
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	apps = apps[1:]
	tt.Frame()
	if !tt.Focused("B") {
		t.Fatal("the dock lost B's focus as A left")
	}
	active, progs := -1, []desk.Program{{Name: "RUN", Icon: icon.Terminal}}
	tray := []desk.TrayIcon{{Icon: icon.Wifi, Tip: "ONLINE"}}
	tt = view(700, 300, newWorld(), func(c *ui.Context) {
		desk.Taskbar(c, "Bar", nil, &active, desk.TaskbarOptions{Programs: progs, Tray: tray})
	})
	click(t, tt, "ONLINE")
	progs = nil
	tt.Frame()
	if !tt.Focused("ONLINE") {
		t.Fatal("the tray lost its focus as START went")
	}
}

func TestScreensaverWakesOnEveryButtonAndTheWheel(t *testing.T) {
	for _, how := range []string{"right", "wheel"} {
		on := true
		tt := view(600, 400, newWorld(), func(c *ui.Context) { desk.Screensaver(c, "LOGO", &on, desk.ScreensaverOptions{}) })
		if how == "right" {
			tt.RightClickAt(100, 100)
		} else {
			tt.Scroll(100, 100, 0, 30)
		}
		tt.Frame()
		if on {
			t.Errorf("a %s did not wake it", how)
		}
	}
}

func TestStateChangesAreLogged(t *testing.T) {
	run := ""
	log := logged(func() {
		tt := view(700, 1600, newWorld(), func(c *ui.Context) {
			desk.DesktopIcon(c, "ICON", icon.File, desk.DesktopIconOptions{})
			desk.Minesweeper(c, "Mines", desk.MinesweeperOptions{Seed: 1})
			desk.Snake(c, "Snake", desk.SnakeOptions{Seed: 1})
			desk.RunDialog(c, "Run", &run, desk.RunDialogOptions{})
		})
		click(t, tt, "ICON")
		// Past the face, F flags the field's first shut cell.
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyF)
		click(t, tt, "Snake")
		click(t, tt, "OPEN:")
		tt.Type("pwd")
	})
	for _, want := range []string{"desktop icon chosen", "minesweeper flag", "snake start", "run dialog typed"} {
		if !strings.Contains(log, want) {
			t.Errorf("no %q in the log:\n%s", want, log)
		}
	}
	if strings.Contains(log, "pwd") {
		t.Error("the log carries what was typed")
	}
}
