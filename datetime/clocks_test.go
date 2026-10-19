package datetime_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/datetime"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestAnalogClockHandsPoint(t *testing.T) {
	clk := newClock()
	clk.t = time.Date(2026, 10, 11, 3, 0, 0, 0, time.UTC)
	tt := view(200, 200, clk, false, func(c *ui.Context) {
		datetime.AnalogClock(c, datetime.AnalogClockOptions{Size: 48})
	})
	r, img := find(t, tt, "03:00:00"), tt.Image()
	cx, cy := int(r.X)+48, int(r.Y)+48
	// Three o'clock: the hour hand runs right, the others up.
	if slot(t, img, cx+20, cy) != pixel.Ink || slot(t, img, cx, cy-30) != pixel.Accent {
		t.Fatal("the hands do not point at three o'clock")
	}
	clk.t = clk.t.Add(15 * time.Second)
	tt.Frame()
	if !shows(tt, "03:00:15") || slot(t, tt.Image(), cx+30, cy) != pixel.Accent {
		t.Fatal("the second hand did not reach fifteen")
	}
	keeps(t, tt.Image())
}

func TestDigitalClockBlinksItsColons(t *testing.T) {
	clk := newClock()
	tt := view(400, 100, clk, false, func(c *ui.Context) {
		datetime.DigitalClock(c, datetime.DigitalClockOptions{})
	})
	r := find(t, tt, "09:30:00")
	// Two digits take 64 pixels, then an 8-pixel gap.
	colon := ui.Rect{X: r.X + 72, Y: r.Y, W: 12, H: r.H}
	on := count(t, tt.Image(), colon, pixel.Bad)
	clk.t = clk.t.Add(500 * time.Millisecond)
	tt.Frame()
	if off := count(t, tt.Image(), colon, pixel.Bad); on == 0 || off != 0 {
		t.Fatalf("colon dots %d then %d", on, off)
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if count(t, tt.Image(), colon, pixel.Bad) == 0 {
		t.Fatal("reduced, the colon went dark")
	}
}

func TestCountdownTurnsBadAndEndsOnce(t *testing.T) {
	clk := newClock()
	deadline, done := clk.t.Add(15*time.Second), 0
	tt := view(300, 100, clk, false, func(c *ui.Context) {
		if datetime.Countdown(c, "Bomb", deadline, datetime.CountdownOptions{}).Done {
			done++
		}
	})
	r := find(t, tt, "Bomb")
	if n := count(t, tt.Image(), r, pixel.Bad); n != 0 {
		t.Fatal("15 s left is Bad already")
	}
	clk.t = clk.t.Add(6 * time.Second)
	tt.Frame()
	if count(t, tt.Image(), r, pixel.Bad) == 0 {
		t.Fatal("9 s left is not Bad")
	}
	clk.t = clk.t.Add(10 * time.Second)
	tt.Frame()
	tt.Frame()
	if done != 1 {
		t.Fatalf("Done %d times", done)
	}
	keeps(t, tt.Image())
}

func TestStopwatchTimesAndLaps(t *testing.T) {
	clk := newClock()
	var started, lapped, reset int
	tt := view(300, 300, clk, false, func(c *ui.Context) {
		r := datetime.Stopwatch(c, "Run", datetime.StopwatchOptions{})
		if r.Started {
			started++
		}
		if r.Lapped {
			lapped++
		}
		if r.Reset {
			reset++
		}
	})
	click(t, tt, "START")
	clk.t = clk.t.Add(3500 * time.Millisecond)
	tt.Frame()
	click(t, tt, "LAP")
	if !shows(tt, "00:03.5") || !tt.HasText("LAP 1  00:03.5") {
		t.Fatalf("texts %q", tt.Texts())
	}
	click(t, tt, "STOP")
	clk.t = clk.t.Add(time.Minute)
	tt.Frame()
	if !shows(tt, "00:03.5") {
		t.Fatalf("a stopped watch moved: %q", tt.Texts())
	}
	click(t, tt, "RESET")
	if !shows(tt, "00:00.0") || started != 1 || lapped != 1 || reset != 1 {
		t.Fatalf("reset left %q", tt.Texts())
	}
}

func TestTimerRingsOnceAndFlashes(t *testing.T) {
	clk := newClock()
	done, set := 0, time.Minute
	tt := view(400, 200, clk, false, func(c *ui.Context) {
		if datetime.Timer(c, "Egg", &set, datetime.TimerOptions{}).Done {
			done++
		}
	})
	if !shows(tt, "01:00") {
		t.Fatalf("a new timer shows %q", tt.Texts())
	}
	click(t, tt, "START")
	clk.t = clk.t.Add(30 * time.Second)
	tt.Frame()
	click(t, tt, "PAUSE")
	clk.t = clk.t.Add(time.Hour)
	tt.Frame()
	if !shows(tt, "00:30") {
		t.Fatalf("paused at %q", tt.Texts())
	}
	click(t, tt, "START")
	clk.t = clk.t.Add(31 * time.Second)
	tt.Frame()
	tt.Frame()
	r := find(t, tt, "00:00")
	lit := count(t, tt.Image(), r, pixel.Bad)
	clk.t = clk.t.Add(250 * time.Millisecond)
	tt.Frame()
	if done != 1 || lit+count(t, tt.Image(), r, pixel.Bad) == 0 || lit > 0 && count(t, tt.Image(), r, pixel.Bad) > 0 {
		t.Fatalf("done %d; the zero did not flash", done)
	}
}

func TestRelativeTimeWordsAndUpdates(t *testing.T) {
	clk := newClock()
	for d, want := range map[time.Duration]string{
		-20 * time.Second: "JUST NOW", -5 * time.Minute: "5 MIN AGO", 2 * time.Hour: "IN 2 H",
		-30 * time.Hour: "YESTERDAY", 30 * time.Hour: "TOMORROW", 75 * time.Hour: "IN 3 DAYS",
	} {
		at := clk.t.Add(d)
		tt := view(300, 60, clk, false, func(c *ui.Context) { datetime.RelativeTime(c, at, datetime.RelativeTimeOptions{}) })
		if !shows(tt, want) {
			t.Fatalf("%v: %q, want %q", d, tt.Texts(), want)
		}
	}
	at := clk.t.Add(-90 * time.Second)
	tt := view(300, 60, clk, false, func(c *ui.Context) { datetime.RelativeTime(c, at, datetime.RelativeTimeOptions{}) })
	clk.t = clk.t.Add(time.Minute)
	tt.Frame()
	if !shows(tt, "2 MIN AGO") {
		t.Fatalf("a minute on: %q", tt.Texts())
	}
}

func TestRelativeTimeAsksForTheNextMinute(t *testing.T) {
	clk := newClock()
	at := clk.t.Add(-90 * time.Second)
	tt := view(300, 60, clk, false, func(c *ui.Context) { datetime.RelativeTime(c, at, datetime.RelativeTimeOptions{}) })
	pixel.ResetDue()
	before := tt.Now()
	tt.Frame()
	// The words change at two minutes: thirty seconds on.
	if wait := pixel.Due().Sub(before); wait < 29*time.Second || wait > 31*time.Second {
		t.Fatalf("next repaint in %v", wait)
	}
}

func TestTimerShowsANewSettingAtOnce(t *testing.T) {
	set := time.Minute
	tt := view(400, 200, newClock(), false, func(c *ui.Context) {
		datetime.Timer(c, "Egg", &set, datetime.TimerOptions{})
	})
	r := find(t, tt, "MIN")
	tt.Press(r.X+r.W+12, r.Y+8)
	if set != 0 || !shows(tt, "00:00") {
		t.Fatalf("set %v while the display reads %q", set, tt.Texts())
	}
	tt.Release(r.X+r.W+12, r.Y+8)
}
