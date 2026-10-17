package datetime_test

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/datetime"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestSunDialDragsHoursAndShowsTheSky(t *testing.T) {
	hour, changes := 9, 0
	tt := view(400, 200, newClock(), false, func(c *ui.Context) {
		if datetime.SunDial(c, "Hour", &hour, datetime.SunDialOptions{}).Changed {
			changes++
		}
	})
	r := find(t, tt, "Hour")
	tt.Press(r.X+r.W/2+2, r.Y+20)
	tt.Release(r.X+r.W/2+2, r.Y+20)
	if hour != 12 {
		t.Fatalf("a press at the middle set %d", hour)
	}
	// Noon puts the sun at the top of the arc.
	if n := count(t, tt.Image(), ui.Rect{X: r.X + r.W/2 - 12, Y: r.Y + 2*(56-6-40) - 12, W: 24, H: 24}, pixel.Warn); n == 0 {
		t.Fatal("no sun at the top at noon")
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyLeft)
	if hour != 10 || changes != 3 {
		t.Fatalf("Left, Left gave %d, %d changes", hour, changes)
	}
	hour = 23
	tt.Frame()
	// At 23:00 the moon sits left of the top.
	moon := ui.Rect{X: r.X + 76 - 12, Y: r.Y + 22 - 12, W: 24, H: 24}
	if count(t, tt.Image(), moon, pixel.Warn) == 0 || count(t, tt.Image(), r, pixel.Shade) < 100 {
		t.Fatal("no moon at eleven at night")
	}
	keeps(t, tt.Image())
}

func TestSeasonWheelPicksByClickAndKeys(t *testing.T) {
	season := pixel.Spring
	tt := view(300, 200, newClock(), false, func(c *ui.Context) {
		datetime.SeasonWheel(c, "Season", &season, datetime.SeasonWheelOptions{})
	})
	r := find(t, tt, "Season")
	tt.ClickAt(r.X+20, r.Y+r.H-20)
	if season != pixel.Autumn {
		t.Fatalf("a click at the bottom left chose %v", season)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if season != pixel.Spring {
		t.Fatalf("Right twice from autumn chose %v", season)
	}
	keeps(t, tt.Image())
}

func TestWeekStripStepsAcrossWeeks(t *testing.T) {
	date := time.Date(2026, 10, 17, 0, 0, 0, 0, time.UTC)
	tt := view(400, 100, newClock(), false, func(c *ui.Context) {
		datetime.WeekStrip(c, "Week", &date, datetime.WeekStripOptions{})
	})
	r := find(t, tt, "Week")
	tt.ClickAt(r.X+4, r.Y+4)
	if day(date) != "2026-10-11" {
		t.Fatalf("a click on Sunday gave %s", day(date))
	}
	tt.Key(0, ui.KeyLeft)
	if day(date) != "2026-10-10" {
		t.Fatalf("Left from Sunday gave %s", day(date))
	}
	keeps(t, tt.Image())
}

var raid = time.Date(2026, 10, 11, 9, 0, 0, 0, time.UTC)

func TestAgendaGroupsByDay(t *testing.T) {
	events := []datetime.Event{
		{Start: raid.Add(26 * time.Hour), Title: "Shop"},
		{Start: raid, Title: "Raid"},
		{Start: raid.Add(3 * time.Hour), Title: "Lunch"},
	}
	tt := view(300, 300, newClock(), false, func(c *ui.Context) {
		datetime.Agenda(c, events, datetime.AgendaOptions{})
	})
	a, b, s := find(t, tt, "SU 11 OCT"), find(t, tt, "09:00 Raid"), find(t, tt, "MO 12 OCT")
	l, sh := find(t, tt, "12:00 Lunch"), find(t, tt, "11:00 Shop")
	if !(a.Y < b.Y && b.Y < l.Y && l.Y < s.Y && s.Y < sh.Y) {
		t.Fatalf("order: %q", tt.Texts())
	}
	keeps(t, tt.Image())
}

func TestAlarmLeverAndBadTimes(t *testing.T) {
	alarms := []datetime.AlarmEntry{{Hour: 7, Minute: 30, Days: [7]bool{1: true, 2: true}}}
	changes := 0
	tt := view(400, 100, newClock(), false, func(c *ui.Context) {
		if datetime.Alarm(c, "Alarms", alarms, datetime.AlarmOptions{}).Changed {
			changes++
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if !alarms[0].On || changes != 1 {
		t.Fatalf("the lever left On %v", alarms[0].On)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("25:00 did not panic")
		}
	}()
	view(100, 100, newClock(), false, func(c *ui.Context) {
		datetime.Alarm(c, "x", []datetime.AlarmEntry{{Hour: 25}}, datetime.AlarmOptions{})
	})
}

func TestYearViewMarksToday(t *testing.T) {
	tt := view(600, 400, newClock(), false, func(c *ui.Context) {
		datetime.YearView(c, 2026, datetime.YearViewOptions{})
	})
	r := find(t, tt, "2026")
	if n := count(t, tt.Image(), r, pixel.Accent); n == 0 {
		t.Fatal("today is not marked")
	}
	keeps(t, tt.Image())
}

func TestDaySchedulePlacesAndSplitsOverlaps(t *testing.T) {
	events := []datetime.Event{
		{Start: raid, Length: 2 * time.Hour, Title: "Raid"},
		{Start: raid.Add(time.Hour), Length: time.Hour, Title: "Call"},
		{Start: raid.Add(4 * time.Hour), Length: 30 * time.Minute, Title: "Lunch"},
	}
	tt := view(400, 500, newClock(), false, func(c *ui.Context) {
		datetime.DaySchedule(c, events, datetime.DayScheduleOptions{})
	})
	tt.Frame()
	a, b, l := find(t, tt, "09:00 Raid"), find(t, tt, "10:00 Call"), find(t, tt, "13:00 Lunch")
	// Hours are 40 pixels; nine o'clock is the second row.
	if a.Y != 40 || b.Y != 80 || b.X <= a.X || l.X != a.X || a.H != 80 || l.W < 2*a.W {
		t.Fatalf("raid %v, call %v, lunch %v", a, b, l)
	}
	keeps(t, tt.Image())
}

func TestEventChipReadsTimeAndTitle(t *testing.T) {
	tt := view(300, 60, newClock(), false, func(c *ui.Context) {
		datetime.EventChip(c, datetime.Event{Start: raid, Title: "Raid"}, datetime.EventChipOptions{})
	})
	if !shows(tt, "09:00 Raid") {
		t.Fatalf("texts %q", tt.Texts())
	}
	keeps(t, tt.Image())
}

func TestOffDateTimeControlsStayPut(t *testing.T) {
	for _, inert := range []bool{false, true} {
		off := !inert
		tm := raid
		week := raid
		month, year, hour, season := raid, 2026, 9, pixel.Spring
		alarms := []datetime.AlarmEntry{{Hour: 7}}
		started := 0
		tt := view(800, 600, newClock(), inert, func(c *ui.Context) {
			ui.Row(c).Gap(8).Children(func() {
				ui.Column(c).Gap(8).Children(func() {
					datetime.TimeInput(c, "Time", &tm, datetime.TimeInputOptions{Disabled: off})
					datetime.WeekStrip(c, "Week", &week, datetime.WeekStripOptions{Disabled: off})
					datetime.MonthPicker(c, "Month", &month, datetime.MonthPickerOptions{Disabled: off})
					datetime.YearPicker(c, "Year", &year, datetime.YearPickerOptions{Disabled: off})
				})
				ui.Column(c).Gap(8).Children(func() {
					datetime.SunDial(c, "Hour", &hour, datetime.SunDialOptions{Disabled: off})
					datetime.SeasonWheel(c, "Season", &season, datetime.SeasonWheelOptions{Disabled: off})
					datetime.Alarm(c, "Alarms", alarms, datetime.AlarmOptions{Disabled: off})
					if datetime.Stopwatch(c, "Run", datetime.StopwatchOptions{Disabled: off}).Started {
						started++
					}
				})
			})
		})
		for _, l := range []string{"Time", "Week", "Month", "Year", "Season", "START"} {
			r := find(t, tt, l)
			tt.ClickAt(r.X+r.W-4, r.Y+r.H-4)
			tt.Key(0, ui.KeyRight)
			tt.Key(0, ui.KeyUp)
			tt.Key(0, ui.KeyEnter)
		}
		h := find(t, tt, "Hour")
		tt.Press(h.X+4, h.Y+4)
		tt.Release(h.X+4, h.Y+4)
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeySpace)
		if !tm.Equal(raid) || !week.Equal(raid) || !month.Equal(raid) || year != 2026 || hour != 9 || season != pixel.Spring || alarms[0].On || started != 0 {
			t.Fatalf("inert %v: off controls changed", inert)
		}
	}
}

func TestDayScheduleKeepsLateBlocksInside(t *testing.T) {
	late := time.Date(2026, 10, 11, 17, 55, 0, 0, time.UTC)
	tt := view(400, 500, newClock(), false, func(c *ui.Context) {
		datetime.DaySchedule(c, []datetime.Event{{Start: late, Length: 5 * time.Minute, Title: "Bell"}}, datetime.DayScheduleOptions{})
	})
	tt.Frame()
	// Five minutes take the three-dot least: six pixels.
	if b := find(t, tt, "17:55 Bell"); b.Y+b.H > 400 || b.H != 6 {
		t.Fatalf("a late five-minute block is %v", b)
	}
}

func TestMonthPickerFitsChineseNames(t *testing.T) {
	month := raid
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: newClock().now, Locale: pixel.Zh})
		datetime.MonthPicker(c, "Month", &month, datetime.MonthPickerOptions{})
	}, 600, 300)
	if w := find(t, tt, "Month").W; w < 4*float32(pixel.Measure("十二月")+2)*2 {
		t.Fatalf("cells of %v pixels for four 十二月", w/4)
	}
	keeps(t, tt.Image())
}

func TestTodayMarksMoveAtMidnight(t *testing.T) {
	clk := newClock()
	clk.t = time.Date(2026, 10, 11, 23, 59, 59, 0, time.UTC)
	date := clk.t
	for name, build := range map[string]func(c *ui.Context){
		"calendar": func(c *ui.Context) { datetime.Calendar(c, "Date", &date, datetime.CalendarOptions{}) },
		"week":     func(c *ui.Context) { datetime.WeekStrip(c, "Week", &date, datetime.WeekStripOptions{}) },
		"year":     func(c *ui.Context) { datetime.YearView(c, 2026, datetime.YearViewOptions{}) },
	} {
		tt := view(600, 400, clk, false, build)
		pixel.ResetDue()
		before := tt.Now()
		tt.Frame()
		if due := pixel.Due(); due.IsZero() || due.Sub(before) > 2*time.Second {
			t.Fatalf("%s: no rebuild asked for midnight, due %v", name, due)
		}
		marks := accents(t, tt)
		clk.t = clk.t.Add(2 * time.Second)
		tt.Frame()
		after := accents(t, tt)
		clk.t = clk.t.Add(-2 * time.Second)
		if len(marks) == 0 || len(after) == 0 || slices.Equal(marks, after) {
			t.Fatalf("%s: the today mark did not move at midnight", name)
		}
	}
}

func TestRelativeTimeRoundsTheFutureUp(t *testing.T) {
	clk := newClock()
	for _, tc := range []struct {
		ahead time.Duration
		words string
		wait  time.Duration
	}{{30 * time.Second, "IN 1 MIN", 30 * time.Second}, {90 * time.Second, "IN 2 MIN", 30 * time.Second}, {2 * time.Minute, "IN 2 MIN", time.Minute}} {
		at := clk.t.Add(tc.ahead)
		tt := view(300, 60, clk, false, func(c *ui.Context) { datetime.RelativeTime(c, at, datetime.RelativeTimeOptions{}) })
		pixel.ResetDue()
		before := tt.Now()
		tt.Frame()
		if wait := pixel.Due().Sub(before); !shows(tt, tc.words) || wait < tc.wait-time.Second || wait > tc.wait+time.Second {
			t.Fatalf("%v ahead: %q, next repaint in %v", tc.ahead, tt.Texts(), wait)
		}
	}
}

// accents lists where Accent dots are.
func accents(t *testing.T, tt *ui.Tester) []int {
	t.Helper()
	img, out := tt.Image(), []int{}
	for y := 0; y < img.Bounds().Dy(); y += 2 {
		for x := 0; x < img.Bounds().Dx(); x += 2 {
			if slot(t, img, x, y) == pixel.Accent {
				out = append(out, y*img.Bounds().Dx()+x)
			}
		}
	}
	return out
}
