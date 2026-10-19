package datetime_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/datetime"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestCalendarMovesByKeysAndPicksByClick(t *testing.T) {
	date := time.Date(2026, 1, 31, 8, 15, 0, 0, time.UTC)
	picks := 0
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		if datetime.Calendar(c, "Date", &date, datetime.CalendarOptions{}).Picked {
			picks++
		}
	})
	r := find(t, tt, "Date")
	// The grid sits under the header and the day names.
	grid := ui.Rect{X: r.X, Y: r.Y + r.H - 6*32, W: 7 * 32, H: 6 * 32}
	tt.ClickAt(grid.X+3*32+4, grid.Y+4)
	if day(date) != "2025-12-31" || picks != 1 {
		t.Fatalf("a click on the first row's Wednesday chose %s", day(date))
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyPageDown)
	if day(date) != "2026-02-01" {
		t.Fatalf("Right, Page Down gave %s", day(date))
	}
	date = time.Date(2026, 1, 31, 8, 15, 0, 0, time.UTC)
	tt.Frame()
	tt.Key(0, ui.KeyPageDown)
	if day(date) != "2026-02-28" {
		t.Fatalf("Page Down from January 31 gave %s", day(date))
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyHome)
	if day(date) != "2026-03-01" || date.Hour() != 8 {
		t.Fatalf("Page Down, Down, Home gave %s %d:00", day(date), date.Hour())
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	if day(date) != "2026-03-31" || picks != 2 {
		t.Fatalf("End, Enter gave %s, %d picks", day(date), picks)
	}
	keeps(t, tt.Image())
}

func TestCalendarMarksTodayAndTheChosenDay(t *testing.T) {
	date := time.Date(2026, 10, 14, 0, 0, 0, 0, time.UTC)
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		datetime.Calendar(c, "Date", &date, datetime.CalendarOptions{})
	})
	r := find(t, tt, "Date")
	top := r.Y + r.H - 6*32
	// October 2026 starts Thursday: the 11th is row 2.
	today := ui.Rect{X: r.X + 6, Y: top + 2*32 + 28, W: 20, H: 2}
	if n := count(t, tt.Image(), today, pixel.Accent); n == 0 {
		t.Fatal("today is not marked")
	}
	chosen := ui.Rect{X: r.X + 3*32, Y: top + 2*32, W: 32, H: 32}
	if n := count(t, tt.Image(), chosen, pixel.Shade); n == 0 {
		t.Fatal("the chosen day is not sunken")
	}
}

func TestOffCalendarStaysPut(t *testing.T) {
	for _, inert := range []bool{false, true} {
		date := time.Date(2026, 10, 14, 0, 0, 0, 0, time.UTC)
		tt := view(400, 300, newClock(), inert, func(c *ui.Context) {
			datetime.Calendar(c, "Date", &date, datetime.CalendarOptions{Disabled: !inert})
		})
		r := find(t, tt, "Date")
		tt.ClickAt(r.X+4, r.Y+r.H-4)
		tt.Key(0, ui.KeyRight)
		click(t, tt, "NEXT")
		if day(date) != "2026-10-14" {
			t.Fatalf("inert %v: an off calendar moved to %s", inert, day(date))
		}
	}
}

func TestDateInputOpensACalendar(t *testing.T) {
	date := time.Date(2026, 10, 11, 0, 0, 0, 0, time.UTC)
	disabled, changes := false, 0
	tt := view(400, 400, newClock(), false, func(c *ui.Context) {
		if datetime.DateInput(c, "Birthday", &date, datetime.DateInputOptions{Disabled: disabled}).Changed {
			changes++
		}
	})
	click(t, tt, "Birthday")
	if !shows(tt, "OCT 2026") {
		t.Fatalf("no calendar: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if day(date) != "2026-10-12" || shows(tt, "OCT 2026") || changes != 1 {
		t.Fatalf("Right, Enter gave %s", day(date))
	}
	click(t, tt, "Birthday")
	disabled = true
	tt.Frame()
	if shows(tt, "OCT 2026") {
		t.Fatal("the calendar stayed open as its field turned off")
	}
	keeps(t, tt.Image())
}

func TestTimeInputTypesAndSteps(t *testing.T) {
	tm := time.Date(2026, 10, 11, 0, 0, 0, 0, time.UTC)
	tt := view(300, 100, newClock(), false, func(c *ui.Context) {
		datetime.TimeInput(c, "Wake", &tm, datetime.TimeInputOptions{})
	})
	click(t, tt, "Wake")
	tt.Key(0, ui.KeyLeft)
	for _, k := range []ui.Key{ui.Key0, ui.Key7, ui.Key3, ui.Key0} {
		tt.Key(0, k)
	}
	if got := tm.Format("15:04"); got != "07:30" {
		t.Fatalf("0730 typed %s", got)
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyDown)
	if got := tm.Format("15:04"); got != "06:31" || day(tm) != "2026-10-11" {
		t.Fatalf("Up, Left, Down gave %s", got)
	}
	tt.Key(0, ui.Key2)
	tt.Key(0, ui.Key9)
	if got := tm.Format("15:04"); got != "09:31" {
		t.Fatalf("29 on hours gave %s", got)
	}
	keeps(t, tt.Image())
}

func TestDateRangePicksAndDithersBetween(t *testing.T) {
	start := time.Date(2026, 10, 11, 0, 0, 0, 0, time.UTC)
	end := start
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		datetime.DateRange(c, "Trip", &start, &end, datetime.DateRangeOptions{})
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if day(start) != "2026-10-11" || day(end) != "2026-10-14" {
		t.Fatalf("Enter, 3×Right, Enter gave %s to %s", day(start), day(end))
	}
	r := find(t, tt, "Trip")
	top := r.Y + r.H - 6*32
	between := ui.Rect{X: r.X + 32, Y: top + 2*32, W: 32, H: 32}
	if n := count(t, tt.Image(), between, pixel.Accent); n == 0 {
		t.Fatal("the days between are not dithered")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	if day(start) != "2026-10-07" || day(end) != "2026-10-14" {
		t.Fatalf("a start, then an end before it, gave %s to %s", day(start), day(end))
	}
}

func TestMonthAndYearPickers(t *testing.T) {
	month, year := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), 2026
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		datetime.MonthPicker(c, "Month", &month, datetime.MonthPickerOptions{})
		datetime.YearPicker(c, "Year", &year, datetime.YearPickerOptions{})
	})
	if !shows(tt, "2026") || !shows(tt, "2020–2029") {
		t.Fatalf("texts %q", tt.Texts())
	}
	r := find(t, tt, "Month")
	tt.ClickAt(r.X+16, r.Y-20)
	tt.ClickAt(r.X+2*r.W/4+4, r.Y+4)
	if month.Format("2006-01") != "2025-03" {
		t.Fatalf("◀ then MAR gave %s", month.Format("2006-01"))
	}
	y := find(t, tt, "Year")
	tt.ClickAt(y.X+3*56+4, y.Y+32+4)
	if year != 2028 {
		t.Fatalf("a click on the second row's fourth year gave %d", year)
	}
	keeps(t, tt.Image())
}

func TestDateInputParksFocusAsItTurnsOff(t *testing.T) {
	date := raid
	disabled := false
	tt := view(400, 400, newClock(), false, func(c *ui.Context) {
		datetime.DateInput(c, "Due", &date, datetime.DateInputOptions{Disabled: disabled})
		pixel.Text(c, "elsewhere", pixel.TextOptions{})
	})
	click(t, tt, "Due")
	disabled = true
	tt.Frame()
	// Focused counts containers: root focus would light "elsewhere" too.
	if tt.Focused("Due") && !tt.Focused("elsewhere") {
		t.Fatal("the focus came back to the field as it turned off")
	}
}

func TestDateInputGetsItsFocusBackAfterAVeil(t *testing.T) {
	date, inert := raid, false
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: newClock().now, Inert: inert})
		datetime.DateInput(c, "Due", &date, datetime.DateInputOptions{})
		pixel.Text(c, "elsewhere", pixel.TextOptions{})
	}, 400, 400)
	click(t, tt, "Due")
	inert = true
	tt.Frame()
	inert = false
	tt.Frame()
	if !tt.Focused("Due") || tt.Focused("elsewhere") {
		t.Fatal("the field did not get its focus back after the veil")
	}
	tt.Key(0, ui.KeyEnter)
	if !shows(tt, "OCT 2026") {
		t.Fatal("Enter did not reopen the calendar")
	}
}
