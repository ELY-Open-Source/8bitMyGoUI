package datetime_test

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/datetime"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestKeysDoNotClickWhereThePointerRests(t *testing.T) {
	day := time.Date(2026, 10, 12, 9, 30, 0, 0, time.UTC)
	for name, tc := range map[string]struct {
		build func(c *ui.Context, v *string)
		// fx is where, across the control, the pointer rests.
		fx    float32
		after ui.Key
	}{
		"WeekStrip": {func(c *ui.Context, v *string) {
			d := day
			datetime.WeekStrip(c, "Week", &d, datetime.WeekStripOptions{})
			if !d.Equal(day) {
				*v = d.Format("Jan 2")
			}
		}, 0.95, 0},
		"SeasonWheel": {func(c *ui.Context, v *string) {
			s := pixel.Spring
			datetime.SeasonWheel(c, "Week", &s, datetime.SeasonWheelOptions{})
			if s != pixel.Spring {
				*v = "turned"
			}
		}, 0.9, 0},
		"TimeInput": {func(c *ui.Context, v *string) {
			tm := day
			datetime.TimeInput(c, "Week", &tm, datetime.TimeInputOptions{})
			// Up after the keys must step the hours.
			if tm.Minute() != day.Minute() {
				*v = tm.Format("15:04")
			}
		}, 0.9, ui.KeyUp},
	} {
		got := ""
		tt := view(400, 300, newClock(), false, func(c *ui.Context) { tc.build(c, &got) })
		r := find(t, tt, "Week")
		tt.Move(r.X+r.W*tc.fx, r.Y+r.H/2)
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeySpace)
		tt.Key(0, ui.KeyEnter)
		if tc.after != 0 {
			tt.Key(0, tc.after)
		}
		if got != "" {
			t.Errorf("%s: Space and Enter clicked where the pointer rests: %s", name, got)
		}
	}
}

func TestStopwatchHoldsWhenTheClockStepsBack(t *testing.T) {
	clk := newClock()
	tt := view(400, 200, clk, false, func(c *ui.Context) { datetime.Stopwatch(c, "Watch", datetime.StopwatchOptions{}) })
	if err := tt.Click("START"); err != nil {
		t.Fatal(err)
	}
	clk.t = clk.t.Add(-1200 * time.Millisecond)
	tt.Frame()
	if err := tt.Click("STOP"); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.Join(tt.Texts(), " "), "-") {
		t.Fatalf("a negative time: %q", tt.Texts())
	}
}
