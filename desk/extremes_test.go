package desk_test

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/desk"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestBadInputPanics(t *testing.T) {
	big := pixel.MustSprite(strings.Repeat("#", 13))
	bad := desk.Wallpaper(9)
	n := 0
	for name, build := range map[string]func(c *ui.Context){
		"paper":      func(c *ui.Context) { desk.Desktop(c, "D", desk.DesktopOptions{Wallpaper: bad}, func() {}) },
		"picker":     func(c *ui.Context) { desk.WallpaperPicker(c, "W", &bad, desk.WallpaperPickerOptions{}) },
		"no label":   func(c *ui.Context) { desk.DesktopIcon(c, "", icon.File, desk.DesktopIconOptions{}) },
		"no icon":    func(c *ui.Context) { desk.DesktopIcon(c, "A", nil, desk.DesktopIconOptions{}) },
		"big icon":   func(c *ui.Context) { desk.DesktopIcon(c, "A", big, desk.DesktopIconOptions{}) },
		"no applets": func(c *ui.Context) { desk.ControlPanel(c, "P", nil, desk.ControlPanelOptions{}) },
		"blank applet": func(c *ui.Context) {
			desk.ControlPanel(c, "P", []desk.Applet{{Icon: icon.File}}, desk.ControlPanelOptions{})
		},
		"no programs": func(c *ui.Context) { desk.StartMenu(c, "S", nil, desk.StartMenuOptions{}) },
		"slashed": func(c *ui.Context) {
			desk.StartMenu(c, "S", []desk.Program{{Name: "A/B", Icon: icon.File}}, desk.StartMenuOptions{})
		},
		"twins": func(c *ui.Context) {
			desk.StartMenu(c, "S", []desk.Program{{Name: "A", Icon: icon.File}, {Name: "A", Icon: icon.File}}, desk.StartMenuOptions{})
		},
		"empty submenu": func(c *ui.Context) {
			desk.StartMenu(c, "S", []desk.Program{{Name: "A", Icon: icon.File, Programs: []desk.Program{}}}, desk.StartMenuOptions{})
		},
		"no tip": func(c *ui.Context) {
			desk.TrayIcons(c, "T", []desk.TrayIcon{{Icon: icon.Wifi}}, desk.TrayIconsOptions{})
		},
		"active": func(c *ui.Context) {
			a := 2
			desk.Taskbar(c, "B", []string{"A", "B"}, &a, desk.TaskbarOptions{})
		},
		"no apps":    func(c *ui.Context) { desk.Dock(c, "D", nil, desk.DockOptions{}) },
		"nameless":   func(c *ui.Context) { desk.Dock(c, "D", []desk.DockApp{{Icon: icon.Mail}}, desk.DockOptions{}) },
		"mine field": func(c *ui.Context) { desk.Minesweeper(c, "M", desk.MinesweeperOptions{Cols: 2, Rows: 2, Mines: 4}) },
		"tiny pit":   func(c *ui.Context) { desk.Snake(c, "S", desk.SnakeOptions{Cols: 3, Rows: 3}) },
		"small shut": func(c *ui.Context) { desk.ShutdownScreen(c, desk.ShutdownScreenOptions{Width: 60, Height: 20}) },
		"no title":   func(c *ui.Context) { desk.BootScreen(c, "", nil, desk.BootScreenOptions{}) },
		"long post":  func(c *ui.Context) { desk.BootScreen(c, "B", make([]string, 20), desk.BootScreenOptions{}) },
		"blank clip": func(c *ui.Context) { desk.ClipboardPanel(c, "C", []string{""}, desk.ClipboardPanelOptions{}) },
		"twin clips": func(c *ui.Context) { desk.ClipboardPanel(c, "C", []string{"a", "a"}, desk.ClipboardPanelOptions{}) },
		"no logo": func(c *ui.Context) {
			on := false
			desk.Screensaver(c, "", &on, desk.ScreensaverOptions{})
		},
	} {
		n++
		mustPanic(t, name, func() { view(600, 300, newWorld(), build) })
	}
	if n != 23 {
		t.Fatalf("%d cases", n)
	}
}

func TestExtremesStayInside(t *testing.T) {
	long := strings.Repeat("W", 60)
	many := make([]string, 30)
	for i := range many {
		many[i] = long + string(rune('A'+i))
	}
	for _, tc := range []struct {
		label string
		build func(c *ui.Context)
	}{
		{long, func(c *ui.Context) { desk.DesktopIcon(c, long, icon.File, desk.DesktopIconOptions{Chosen: true}) }},
		{"P", func(c *ui.Context) {
			desk.ControlPanel(c, "P", []desk.Applet{{Name: long, Icon: icon.File}}, desk.ControlPanelOptions{})
		}},
		{"B", func(c *ui.Context) {
			a := 3
			desk.Taskbar(c, "B", many, &a, desk.TaskbarOptions{})
		}},
		{"C", func(c *ui.Context) { desk.ClipboardPanel(c, "C", []string{long + long}, desk.ClipboardPanelOptions{}) }},
		{long, func(c *ui.Context) { desk.BootScreen(c, long, []string{long + long}, desk.BootScreenOptions{}) }},
	} {
		tt := view(1400, 700, newWorld(), tc.build)
		tt.Frame()
		x, y, w, h := dots(find(t, tt, tc.label))
		if n := spill(tt.Image(), x, y, w, h); n > 0 {
			t.Errorf("%.12s spills %d dots", tc.label, n)
		}
	}
}
