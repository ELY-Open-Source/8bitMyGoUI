package desk_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/desk"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestDiskDefragTakesAClockStepBack(t *testing.T) {
	wd := newWorld()
	tt := view(600, 300, wd, func(c *ui.Context) {
		desk.DiskDefrag(c, "Disk", true, desk.DiskDefragOptions{Blocks: 64, Seed: 4})
	})
	later(tt, wd, time.Second)
	later(tt, wd, -2*time.Second)
	later(tt, wd, time.Second)
}

func TestTaskbarCursorStaysInView(t *testing.T) {
	windows := make([]string, 30)
	for i := range windows {
		windows[i] = fmt.Sprint("W", i)
	}
	active := -1
	tt := view(700, 300, newWorld(), func(c *ui.Context) {
		desk.Taskbar(c, "Bar", windows, &active, desk.TaskbarOptions{})
	})
	tt.Frame()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	x, y, _, h := dots(find(t, tt, "Bar"))
	// The cursor's Accent ring shows inside the bar.
	if count(tt.Image(), x, y, 320, h, pixel.Accent) == 0 {
		t.Fatal("End moved the cursor out of sight")
	}
	tt.Key(0, ui.KeyEnter)
	if active != 29 {
		t.Fatalf("Enter chose %d", active)
	}
}

func TestSnakeStartedByAnArrowSaysSo(t *testing.T) {
	started := 0
	tt := view(400, 300, newWorld(), func(c *ui.Context) {
		if desk.Snake(c, "Snake", desk.SnakeOptions{Seed: 2}).Started {
			started++
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyRight)
	if started != 1 {
		t.Fatalf("an arrow start reported %d starts", started)
	}
}

func TestSizesThatCannotFitPanic(t *testing.T) {
	for name, build := range map[string]func(c *ui.Context){
		"boot": func(c *ui.Context) {
			desk.BootScreen(c, "B", []string{"L"}, desk.BootScreenOptions{Height: 64})
		},
		"room": func(c *ui.Context) {
			x, y := 0, 0
			desk.AppWindow(c, "W", &x, &y, desk.AppWindowOptions{Room: [2]int{320, 100}}, func() {})
		},
	} {
		mustPanic(t, name, func() { view(600, 400, newWorld(), build) })
	}
}

func TestClocksSteppingBackHoldAtZero(t *testing.T) {
	wd := newWorld()
	on := true
	tt := view(600, 400, wd, func(c *ui.Context) {
		desk.Screensaver(c, "LOGO", &on, desk.ScreensaverOptions{})
		desk.BootScreen(c, "BIOS", []string{"A"}, desk.BootScreenOptions{})
	})
	later(tt, wd, -time.Hour)
	tt = view(600, 400, wd, func(c *ui.Context) { desk.Minesweeper(c, "Mines", desk.MinesweeperOptions{Seed: 1}) })
	click(t, tt, "MINEFIELD")
	later(tt, wd, -time.Hour)
	if _, ok := tt.Find("000"); !ok {
		t.Fatalf("a clock stepped back shows %q", tt.Texts())
	}
}
