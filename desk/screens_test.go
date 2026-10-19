package desk_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/desk"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestShutdownScreenSaysItIsSafe(t *testing.T) {
	tt := view(700, 300, newWorld(), func(c *ui.Context) { desk.ShutdownScreen(c, desk.ShutdownScreenOptions{}) })
	find(t, tt, "IT'S NOW SAFE TO TURN OFF YOUR COMPUTER.")
}

func TestBootScreenCountsThenPrints(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = reduce
		done := 0
		tt := view(700, 400, wd, func(c *ui.Context) {
			if desk.BootScreen(c, "BIOS", []string{"A", "B", "C"}, desk.BootScreenOptions{}).Done {
				done++
			}
		})
		x, y, w, _ := dots(find(t, tt, "BIOS"))
		// POST lines print in Faint, one each, under the count.
		lines := func() int {
			n := 0
			for i := range 3 {
				if count(tt.Image(), x, y+6+(3+i)*16, w, 16, pixel.Faint) > 0 {
					n++
				}
			}
			return n
		}
		first := lines()
		later(tt, wd, 2*time.Second)
		if done != 1 || lines() != 3 || reduce != (first == 3) {
			t.Errorf("reduce %v: done %d, %d lines at first, %d after", reduce, done, first, lines())
		}
		later(tt, wd, time.Second)
		if done != 1 {
			t.Errorf("reduce %v: done reported %d times", reduce, done)
		}
	}
}

func TestScreensaverWakesOnAnyInput(t *testing.T) {
	for _, wake := range []string{"key", "click", "move"} {
		wd := newWorld()
		on, woke := false, 0
		tt := view(600, 400, wd, func(c *ui.Context) {
			if input.Button(c, "SAVER", input.ButtonOptions{}).Clicked {
				on = true
			}
			if desk.Screensaver(c, "LOGO", &on, desk.ScreensaverOptions{}).Woke {
				woke++
			}
		})
		click(t, tt, "SAVER")
		if !on || !tt.Focused("LOGO") {
			t.Fatalf("%s: the saver did not take over", wake)
		}
		// A small nudge of the pointer is not a wake.
		tt.Move(300, 200)
		tt.Move(304, 200)
		before := moves(tt, wd)
		switch wake {
		case "key":
			tt.Key(0, ui.KeyA)
		case "click":
			tt.ClickAt(10, 10)
		case "move":
			tt.Move(380, 300)
		}
		tt.Frame()
		if on || woke != 1 || !before {
			t.Errorf("%s: on %v, woke %d, logo moved %v", wake, on, woke, before)
		}
		if !tt.Focused("SAVER") {
			t.Errorf("%s: the focus did not come back", wake)
		}
	}
}

func TestLockScreenOpensOnAKeyNotAMove(t *testing.T) {
	wd := newWorld()
	locked, open := true, 0
	tt := view(600, 400, wd, func(c *ui.Context) {
		if desk.LockScreen(c, "Lock", &locked, desk.LockScreenOptions{}).Unlocked {
			open++
		}
	})
	tt.Move(10, 10)
	tt.Move(300, 300)
	tt.Frame()
	if !locked {
		t.Fatal("a move unlocked it")
	}
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if locked || open != 1 {
		t.Fatalf("Space left it locked %v, %d unlocks", locked, open)
	}
}
