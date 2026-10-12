package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/datetime"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var datetimeEntries = []Entry{
	{196, "Calendar", "日历", "datetime", "Arrows walk the days, Page Up and Down the months. Today wears a line.", demoCalendar},
	{197, "DateInput", "日期框", "datetime", "YYYY-MM-DD. Click it and a calendar drops down.", demoDateInput},
	{198, "TimeInput", "时间框", "datetime", "HH:MM. Type 0730, or nudge with Up and Down.", demoTimeInput},
	{199, "DateRange", "日期区间", "datetime", "Pick a start, then an end. The days between dither.", demoDateRange},
	{200, "MonthPicker", "月份选择", "datetime", "A year of months; ◀ ▶ turn the year.", demoMonthPicker},
	{201, "YearPicker", "年份选择", "datetime", "A decade at a glance; ◀ ▶ turn decades.", demoYearPicker},
	{202, "AnalogClock", "指针钟", "datetime", "Hands drawn dot by dot with Bresenham lines. It ticks.", demoAnalogClock},
	{203, "DigitalClock", "数字钟", "datetime", "Seven-segment HH:MM:SS. The colons blink, like a VCR.", demoDigitalClock},
	{204, "Countdown", "倒计时", "datetime", "Time left; it turns red for the last ten seconds.", demoCountdown},
	{205, "Stopwatch", "秒表", "datetime", "START, STOP, LAP, RESET. Speedrunners, warm up.", demoStopwatch},
}

type timeDemo struct {
	date, due, wake, from, to, month time.Time
	year, hour                       int
	season                           pixel.Season
	deadline                         time.Time
	week                             time.Time
	alarms                           []datetime.AlarmEntry
	egg                              time.Duration
}

func timeState(c *ui.Context) *timeDemo {
	now := pixel.K(c).Now()
	return ui.Local(c.Root(), demoKey("datetime"), func() timeDemo {
		return timeDemo{date: now, due: now.AddDate(0, 0, 7), wake: time.Date(now.Year(), now.Month(), now.Day(), 7, 30, 0, 0, now.Location()),
			from: now, to: now.AddDate(0, 0, 4), month: now, year: now.Year(), hour: now.Hour(), season: pixel.Autumn,
			deadline: now.Add(30 * time.Second), week: now, egg: 3 * time.Minute,
			alarms: []datetime.AlarmEntry{{Hour: 7, Minute: 30, Days: [7]bool{1: true, 2: true, 3: true, 4: true, 5: true}, On: true}, {Hour: 10, Days: [7]bool{0: true, 6: true}}}}
	})
}

func demoCalendar(c *ui.Context, log func(string)) {
	s := timeState(c)
	r := datetime.Calendar(c, "Date", &s.date, datetime.CalendarOptions{})
	if r.Changed {
		log(s.date.Format("2006-01-02"))
	}
	if r.Picked {
		log("picked " + s.date.Format("2006-01-02"))
	}
}

func demoDateInput(c *ui.Context, log func(string)) {
	s := timeState(c)
	if datetime.DateInput(c, "Quest due", &s.due, datetime.DateInputOptions{}).Changed {
		log("due " + s.due.Format("2006-01-02"))
	}
}

func demoTimeInput(c *ui.Context, log func(string)) {
	s := timeState(c)
	if datetime.TimeInput(c, "Wake up", &s.wake, datetime.TimeInputOptions{}).Changed {
		log("wake at " + s.wake.Format("15:04"))
	}
}

func demoDateRange(c *ui.Context, log func(string)) {
	s := timeState(c)
	if datetime.DateRange(c, "Vacation", &s.from, &s.to, datetime.DateRangeOptions{}).Changed {
		log(s.from.Format("01-02") + " to " + s.to.Format("01-02"))
	}
}

func demoMonthPicker(c *ui.Context, log func(string)) {
	s := timeState(c)
	if datetime.MonthPicker(c, "Month", &s.month, datetime.MonthPickerOptions{}).Changed {
		log(s.month.Format("2006-01"))
	}
}

func demoYearPicker(c *ui.Context, log func(string)) {
	s := timeState(c)
	if datetime.YearPicker(c, "Year", &s.year, datetime.YearPickerOptions{}).Changed {
		log(fmt.Sprint(s.year))
	}
}

func demoAnalogClock(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(8)).Children(func() {
		datetime.AnalogClock(c, datetime.AnalogClockOptions{Size: 64})
		datetime.AnalogClock(c, datetime.AnalogClockOptions{Size: 32})
	})
}

func demoDigitalClock(c *ui.Context, log func(string)) {
	datetime.DigitalClock(c, datetime.DigitalClockOptions{})
}

func demoCountdown(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := timeState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		if datetime.Countdown(c, "Bomb", s.deadline, datetime.CountdownOptions{}).Done {
			log("boom")
		}
		if input.Button(c, "DEFUSE", input.ButtonOptions{Size: input.Small}).Clicked {
			s.deadline = k.Now().Add(30 * time.Second)
			log("thirty more seconds")
		}
	})
}

func demoStopwatch(c *ui.Context, log func(string)) {
	r := datetime.Stopwatch(c, "Speedrun", datetime.StopwatchOptions{})
	for name, on := range map[string]bool{"start": r.Started, "stop": r.Stopped, "lap": r.Lapped, "reset": r.Reset} {
		if on {
			log(name)
		}
	}
}
