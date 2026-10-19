package datetime_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/datetime"
)

func TestPopupsAndPagesAreLogged(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(old)
	date, month, year := raid, raid, 2026
	tt := view(600, 500, newClock(), false, func(c *ui.Context) {
		ui.Row(c).Gap(8).Children(func() {
			datetime.DateInput(c, "Due", &date, datetime.DateInputOptions{})
			datetime.MonthPicker(c, "Month", &month, datetime.MonthPickerOptions{})
			datetime.YearPicker(c, "Year", &year, datetime.YearPickerOptions{})
		})
	})
	click(t, tt, "Due")
	tt.Key(0, ui.KeyEscape)
	for _, l := range []string{"Month", "Year"} {
		r := find(t, tt, l)
		tt.ClickAt(r.X+16, r.Y-20)
	}
	out := buf.String()
	for _, want := range []string{"date input popup", "month picker turned", "year picker turned"} {
		if !strings.Contains(out, want) {
			t.Errorf("no %q in the logs:\n%s", want, out)
		}
	}
}

func TestInertTimerKeepsItsSetting(t *testing.T) {
	set := time.Minute
	tt := view(400, 200, newClock(), true, func(c *ui.Context) {
		datetime.Timer(c, "Egg", &set, datetime.TimerOptions{})
	})
	r := find(t, tt, "MIN")
	tt.Press(r.X+r.W+12, r.Y+8)
	tt.Release(r.X+r.W+12, r.Y+8)
	if set != time.Minute {
		t.Fatalf("an inert timer was set to %v", set)
	}
}
