package feedback_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/feedback"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestCountUpStepsOrJumps(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		clk := newClock()
		v := 0
		tt := view(400, 100, clk, reduce, func(c *ui.Context) {
			feedback.CountUp(c, v, feedback.CountUpOptions{Slot: pixel.Ink})
		})
		v = 800
		wait(tt, clk, 0)
		r := find(t, tt, "800")
		mid := slots(t, tt.Image(), r)
		wait(tt, clk, time.Second)
		end := slots(t, tt.Image(), r)
		if counted := !equal(mid, end); counted == reduce {
			t.Errorf("reduce %v: counted in steps %v", reduce, counted)
		}
	}
}

func TestPulseStopsUnderReduceMotion(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		clk := newClock()
		tt := view(300, 200, clk, reduce, func(c *ui.Context) {
			feedback.Pulse(c, feedback.PulseOptions{}, func() {
				ui.Box(c).Size(20, 20).Label("dot")
			}).Label("ring")
		})
		if moved := moves(t, tt, clk, "ring", 200*time.Millisecond); moved == reduce {
			t.Errorf("reduce %v: the ring moved %v", reduce, moved)
		}
	}
}

func TestLevelUpShowsOnlyOnARise(t *testing.T) {
	clk := newClock()
	lv := 4
	tt := view(500, 200, clk, false, func(c *ui.Context) {
		feedback.LevelUp(c, lv, feedback.LevelUpOptions{})
	})
	lv = 3
	tt.Frame()
	if _, ok := tt.Find("LEVEL UP! LV 3"); ok {
		t.Fatal("LEVEL UP showed without a rise")
	}
	lv = 5
	tt.Frame()
	r := find(t, tt, "LEVEL UP! LV 5")
	if count(slots(t, tt.Image(), r), pixel.Accent) != 0 {
		t.Fatal("the banner did not start off its box")
	}
	wait(tt, clk, time.Second)
	if count(slots(t, tt.Image(), r), pixel.Accent) == 0 {
		t.Fatal("the banner did not slide in")
	}
	wait(tt, clk, 3*time.Second)
	if _, ok := tt.Find("LEVEL UP! LV 5"); ok {
		t.Fatal("the banner stayed")
	}
}

func TestAchievementComesAndGoes(t *testing.T) {
	clk := newClock()
	cue := 0
	tt := view(500, 200, clk, true, func(c *ui.Context) {
		feedback.Achievement(c, "Slime Tamer", feedback.AchievementOptions{Cue: cue})
	})
	if _, ok := tt.Find("ACHIEVEMENT UNLOCKED Slime Tamer"); ok {
		t.Fatal("the card showed before its cue")
	}
	cue++
	tt.Frame()
	find(t, tt, "ACHIEVEMENT UNLOCKED Slime Tamer")
	wait(tt, clk, 3*time.Second)
	find(t, tt, "ACHIEVEMENT UNLOCKED Slime Tamer")
	wait(tt, clk, 2*time.Second)
	if _, ok := tt.Find("ACHIEVEMENT UNLOCKED Slime Tamer"); ok {
		t.Fatal("the card stayed")
	}
}

func TestComboCounterPopsOnARise(t *testing.T) {
	clk := newClock()
	n := 0
	tt := view(400, 200, clk, false, func(c *ui.Context) {
		feedback.ComboCounter(c, n, feedback.ComboCounterOptions{})
	})
	if _, ok := tt.Find("×0 COMBO!"); ok {
		t.Fatal("a zero combo showed")
	}
	n = 12
	tt.Frame()
	r := find(t, tt, "×12 COMBO!")
	popped := count(slots(t, tt.Image(), r), pixel.Warn)
	wait(tt, clk, 500*time.Millisecond)
	if rest := count(slots(t, tt.Image(), r), pixel.Warn); rest >= popped {
		t.Fatalf("no pop: %d Warn dots, then %d", popped, rest)
	}
}

func TestSaveIndicatorStates(t *testing.T) {
	saving := false
	tt := view(300, 100, newClock(), false, func(c *ui.Context) {
		feedback.SaveIndicator(c, saving, feedback.SaveIndicatorOptions{})
	})
	if _, ok := tt.Find("SAVED"); ok {
		t.Fatal("SAVED before any save")
	}
	saving = true
	tt.Frame()
	find(t, tt, "SAVING…")
	saving = false
	tt.Frame()
	find(t, tt, "SAVED")
}

func TestBusyOverlayBlocksInputAndGivesFocusBack(t *testing.T) {
	busy, clicks, raw := false, 0, true
	tt := view(400, 200, newClock(), false, func(c *ui.Context) {
		input.Button(c, "top", input.ButtonOptions{})
		feedback.BusyOverlay(c, busy, feedback.BusyOverlayOptions{}, func() {
			if input.Button(c, "SAVE", input.ButtonOptions{}).Clicked {
				clicks++
			}
			// Plain boxes, not library controls.
			if ui.Box(c).Size(40, 20).Label("raw").Clicked() {
				clicks++
			}
			if raw && ui.Box(c).Size(40, 20).Label("rawkey").Focusable().Shortcut(0, ui.KeyEnter) {
				clicks++
			}
		})
		ui.Box(c).Size(4, 4).Label("elsewhere")
		input.Button(c, "after", input.ButtonOptions{})
	})
	// Focused counts ancestors; a sibling tells the two apart.
	own := func() bool { return tt.Focused("SAVE") && !tt.Focused("elsewhere") }
	click(t, tt, "rawkey")
	busy = true
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if clicks != 0 || tt.Focused("rawkey") && !tt.Focused("elsewhere") {
		t.Fatal("a plain focusable inside kept the focus or its keys")
	}
	busy, raw = false, false
	tt.Frame()
	click(t, tt, "SAVE")
	busy = true
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	// Focus waits on the veil; Tab goes on from it.
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("after") || tt.Focused("elsewhere") {
		t.Fatal("Tab did not pass over the busy content")
	}
	// Presses land on the veil, which keeps the focus.
	click(t, tt, "SAVE")
	click(t, tt, "raw")
	if clicks != 1 || own() {
		t.Fatalf("busy content took input: %d clicks, focus %v", clicks, own())
	}
	busy = false
	tt.Frame()
	if !own() {
		t.Fatal("the focus did not come back")
	}
	tt.Key(0, ui.KeyEnter)
	if clicks != 2 {
		t.Fatal("the content stayed blocked")
	}
}

func TestBusyContentIsOneTabStop(t *testing.T) {
	tt := view(400, 200, newClock(), false, func(c *ui.Context) {
		input.Button(c, "before", input.ButtonOptions{})
		feedback.BusyOverlay(c, true, feedback.BusyOverlayOptions{}, func() {
			ui.Box(c).Size(40, 20).Label("raw").Focusable()
			input.Button(c, "SAVE", input.ButtonOptions{})
		})
		input.Button(c, "after", input.ButtonOptions{})
	})
	click(t, tt, "before")
	tt.Key(0, ui.KeyTab)
	// FocusMark rings the veil two dots outside it.
	v := find(t, tt, "PLEASE WAIT")
	if count(slots(t, tt.Image(), ui.Rect{X: v.X, Y: v.Y + v.H + 2, W: v.W, H: 2}), pixel.Accent) == 0 {
		t.Fatal("the focused veil showed no focus")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("after") || tt.Focused("before") {
		t.Fatal("two Tabs did not cross the busy content")
	}
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(ui.Shift, ui.KeyTab)
	if !tt.Focused("before") || tt.Focused("after") {
		t.Fatal("two Shift+Tabs did not cross back")
	}
}

func TestReadyGoCountsAndReportsDone(t *testing.T) {
	clk := newClock()
	cue, done := 0, 0
	tt := view(400, 300, clk, true, func(c *ui.Context) {
		if feedback.ReadyGo(c, feedback.ReadyGoOptions{Cue: cue}).Done {
			done++
		}
	})
	cue++
	tt.Frame()
	for _, want := range []string{"3", "2", "1", "GO!"} {
		find(t, tt, want)
		wait(tt, clk, 500*time.Millisecond)
		wait(tt, clk, 500*time.Millisecond)
	}
	if _, ok := tt.Find("GO!"); ok || done != 1 {
		t.Fatalf("GO! stayed, or Done came %d times", done)
	}
}

func TestReadyGoReportsDoneAfterAJump(t *testing.T) {
	clk := newClock()
	cue, done := 0, 0
	tt := view(400, 300, clk, false, func(c *ui.Context) {
		if feedback.ReadyGo(c, feedback.ReadyGoOptions{Cue: cue}).Done {
			done++
		}
	})
	cue++
	tt.Frame()
	wait(tt, clk, 5*time.Second)
	wait(tt, clk, time.Second)
	if done != 1 {
		t.Fatalf("a jump past GO! gave Done %d times", done)
	}
}

func TestShakeJoltsUnlessReduced(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		clk := newClock()
		cue := 0
		tt := view(300, 100, clk, reduce, func(c *ui.Context) {
			feedback.Shake(c, feedback.ShakeOptions{Cue: cue}, func() {
				ui.Box(c).Size(20, 20).Label("box")
			})
		})
		x := find(t, tt, "box").X
		cue++
		tt.Frame()
		if jolted := find(t, tt, "box").X != x; jolted == reduce {
			t.Errorf("reduce %v: jolted %v", reduce, jolted)
		}
		wait(tt, clk, time.Second)
		if find(t, tt, "box").X != x {
			t.Error("the box did not settle")
		}
	}
}
