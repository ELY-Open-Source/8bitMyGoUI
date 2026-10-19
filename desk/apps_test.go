package desk_test

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/desk"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestAppWindowDragsArrowsAndStaysInside(t *testing.T) {
	x, y, closed, off := 10, 10, 0, false
	tt := view(700, 500, newWorld(), func(c *ui.Context) {
		ui.Box(c).Size(640, 400).Children(func() {
			r := desk.AppWindow(c, "WIN", &x, &y, desk.AppWindowOptions{Width: 100, Height: 40, Room: [2]int{320, 200}, Disabled: off}, func() {})
			if r.Closed {
				closed++
			}
		})
	})
	b := find(t, tt, "WIN")
	// The title bar drags the window by whole dots.
	tt.Press(b.X+20, b.Y+10)
	tt.Move(b.X+20+41, b.Y+10+20)
	tt.Release(b.X+20+41, b.Y+10+20)
	if x != 30 || y != 20 {
		t.Fatalf("a 41×20 DIP drag left it at %d, %d", x, y)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	if x != 34 || y != 24 {
		t.Fatalf("arrows left it at %d, %d", x, y)
	}
	// It stays in its room, title bar and all.
	for range 100 {
		tt.Key(0, ui.KeyRight)
	}
	if x != 320-100-6 {
		t.Fatalf("far right is %d", x)
	}
	// On the bar, a key sent while turning off fails.
	click(t, tt, "CLOSE")
	tt.Key(ui.Shift, ui.KeyTab)
	if !tt.Focused("WIN") {
		t.Fatal("Shift+Tab did not reach the title bar")
	}
	off = true
	tt.Key(0, ui.KeyLeft)
	if closed != 1 || x != 214 {
		t.Fatalf("closed %d; a disabled window moved to %d", closed, x)
	}
}

func TestCalculatorAddsUp(t *testing.T) {
	tt := view(400, 300, newWorld(), func(c *ui.Context) { desk.Calculator(c, "Calc", desk.CalculatorOptions{}) })
	tt.Key(0, ui.KeyTab)
	press := func(keys string) {
		for _, r := range keys {
			switch r {
			case '+':
				tt.Key(ui.Shift, ui.KeyEqual)
			case '*':
				tt.Key(ui.Shift, ui.Key8)
			case '-':
				tt.Key(0, ui.KeyMinus)
			case '/':
				tt.Key(0, ui.KeySlash)
			case '=':
				tt.Key(0, ui.KeyEqual)
			case 'C':
				tt.Key(0, ui.KeyC)
			default:
				tt.Key(0, ui.Key0+ui.Key(r-'0'))
			}
		}
	}
	for keys, want := range map[string]string{
		"12+30=": "42", "7/2=": "3.5", "2*3-1=": "5", "100+0=": "100", "1/3=": "0.3333333333", "0-5=": "-5", "5*0=": "0",
		"1/0=": "GAME OVER", "999999*999999=": "999998000001", "9999999*999999=": "GAME OVER", "100000000000+0=": "100000000000",
	} {
		press("C" + keys)
		if _, ok := tt.Find(want); !ok {
			t.Errorf("%s gave %q, want %q", keys, tt.Texts(), want)
		}
	}
	// After GAME OVER only C continues.
	press("C1/0=5")
	find(t, tt, "GAME OVER")
	press("C")
	find(t, tt, "0")
	// The pad's keys work by arrows and Enter too.
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyEnter)
	find(t, tt, "7")
}

func TestNotepadMenus(t *testing.T) {
	note, saved := "hi", 0
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		if desk.Notepad(c, "Pad", &note, desk.NotepadOptions{}).Saved {
			saved++
		}
	})
	click(t, tt, "EDIT")
	click(t, tt, "TIME/DATE")
	if note != "hi10:30 2026-09-02" {
		t.Fatalf("TIME/DATE gave %q", note)
	}
	click(t, tt, "FILE")
	click(t, tt, "SAVE")
	click(t, tt, "FILE")
	click(t, tt, "NEW")
	if saved != 1 || note != "" {
		t.Fatalf("saved %d, note %q", saved, note)
	}
}

func TestRunDialogRunsOrQuits(t *testing.T) {
	wd := newWorld()
	run, ran, quit := "", 0, 0
	tt := view(600, 300, wd, func(c *ui.Context) {
		r := desk.RunDialog(c, "Run", &run, desk.RunDialogOptions{})
		if r.Submitted {
			ran++
		}
		if r.Canceled {
			quit++
		}
	})
	click(t, tt, "OK")
	if ran != 0 {
		t.Fatal("a blank name ran")
	}
	click(t, tt, "OPEN:")
	log := logged(func() {
		tt.Type("snake")
		tt.Key(0, ui.KeyEnter)
	})
	if ran != 1 || strings.Contains(log, "snake") || !strings.Contains(log, "runes=5") {
		t.Fatalf("ran %d; log:\n%s", ran, log)
	}
	click(t, tt, "OK")
	tt.Key(0, ui.KeyEscape)
	click(t, tt, "CANCEL")
	click(t, tt, "CLOSE")
	if ran != 2 || quit != 3 {
		t.Fatalf("ran %d, quit %d", ran, quit)
	}
	// Veiled, Enter in the field runs nothing.
	click(t, tt, "OPEN:")
	wd.inert = true
	tt.Key(0, ui.KeyEnter)
	if ran != 2 {
		t.Fatal("a veiled dialog ran")
	}
}

func TestClipboardPanelCopiesAndKeepsFocus(t *testing.T) {
	wd := newWorld()
	clips, copied := []string{"alpha", "beta secret"}, -1
	tt := view(600, 300, wd, func(c *ui.Context) {
		if i := desk.ClipboardPanel(c, "Clips", clips, desk.ClipboardPanelOptions{}).Copied; i >= 0 {
			copied = i
		}
	})
	log := logged(func() { click(t, tt, "beta secret") })
	x, y, w, h := dots(find(t, tt, "beta secret"))
	if copied != 1 || tt.Clipboard() != "beta secret" || strings.Contains(log, "secret") || count(tt.Image(), x, y, w, h, pixel.Good) == 0 {
		t.Fatalf("copied %d, clipboard %q, no COPIED!, or log:\n%s", copied, tt.Clipboard(), log)
	}
	// A new clip on top keeps the focus on beta.
	clips = append([]string{"gamma"}, clips...)
	tt.Frame()
	if !tt.Focused("beta secret") {
		t.Fatal("the focus left beta")
	}
	later(tt, wd, 3e9)
	x, y, w, h = dots(find(t, tt, "beta secret"))
	if count(tt.Image(), x, y, w, h, pixel.Good) > 0 {
		t.Fatal("COPIED! outstayed two seconds")
	}
}
