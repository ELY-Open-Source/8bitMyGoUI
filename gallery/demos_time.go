package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/datetime"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var timeEntries = []Entry{
	{206, "Timer", "定时器", "datetime", "Set minutes and seconds, START, PAUSE. Zero flashes.", demoTimer},
	{207, "SunDial", "日晷", "datetime", "Drag across the sky: the sun by day, the moon by night.", demoSunDial},
	{208, "SeasonWheel", "四季轮", "datetime", "Four seasons on a wheel. Click one, or turn it with the arrows.", demoSeasonWheel},
	{209, "WeekStrip", "周条", "datetime", "Seven days in a row. Left and Right walk into the next week.", demoWeekStrip},
	{210, "Agenda", "日程", "datetime", "Events under the day they fall on.", demoAgenda},
	{211, "EventChip", "事件块", "datetime", "A small chip: time, title, and a colored edge.", demoEventChip},
	{212, "RelativeTime", "相对时间", "datetime", "3 MIN AGO, IN 2 H, TOMORROW. It keeps itself up to date.", demoRelativeTime},
	{213, "Alarm", "闹钟", "datetime", "Alarm rows: time, days, and a lever to arm each one.", demoAlarm},
	{214, "YearView", "年视图", "datetime", "Twelve tiny months; today is lit.", demoYearView},
	{215, "DaySchedule", "日程表", "datetime", "Hour rows with blocks placed by start and length. Overlaps share.", demoDaySchedule},
}

func demoTimer(c *ui.Context, log func(string)) {
	s := timeState(c)
	if datetime.Timer(c, "Egg timer", &s.egg, datetime.TimerOptions{}).Done {
		log("ding! the egg is ready")
	}
}

func demoSunDial(c *ui.Context, log func(string)) {
	s := timeState(c)
	if datetime.SunDial(c, "Hour", &s.hour, datetime.SunDialOptions{}).Changed {
		log(fmt.Sprintf("%02d:00", s.hour))
	}
}

func demoSeasonWheel(c *ui.Context, log func(string)) {
	s := timeState(c)
	if datetime.SeasonWheel(c, "Season", &s.season, datetime.SeasonWheelOptions{}).Changed {
		log([]string{"spring", "summer", "autumn", "winter"}[s.season])
	}
}

func demoWeekStrip(c *ui.Context, log func(string)) {
	s := timeState(c)
	if datetime.WeekStrip(c, "Week", &s.week, datetime.WeekStripOptions{}).Changed {
		log(s.week.Format("Mon 2006-01-02"))
	}
}

// today returns now's date at h:m.
func today(c *ui.Context, h, m int) time.Time {
	n := pixel.K(c).Now()
	return time.Date(n.Year(), n.Month(), n.Day(), h, m, 0, 0, n.Location())
}

func demoAgenda(c *ui.Context, log func(string)) {
	datetime.Agenda(c, []datetime.Event{
		{Start: today(c, 9, 0), Title: "Slime patrol"}, {Start: today(c, 13, 30), Title: "Lunch at the inn"},
		{Start: today(c, 10, 0).AddDate(0, 0, 1), Title: "Boss raid"}, {Start: today(c, 18, 0).AddDate(0, 0, 2), Title: "Fireworks"},
	}, datetime.AgendaOptions{})
}

func demoEventChip(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		datetime.EventChip(c, datetime.Event{Start: today(c, 9, 0), Title: "Slime patrol"}, datetime.EventChipOptions{})
		datetime.EventChip(c, datetime.Event{Start: today(c, 20, 0), Title: "Boss raid"}, datetime.EventChipOptions{Slot: pixel.Bad})
		datetime.EventChip(c, datetime.Event{Start: today(c, 12, 0), Title: "Picnic"}, datetime.EventChipOptions{Slot: pixel.Good})
	})
}

func demoRelativeTime(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	now := k.Now()
	ui.Column(c).Gap(k.D(2)).Children(func() {
		for _, r := range []struct {
			what string
			at   time.Time
		}{{"Saved", now.Add(-3 * time.Minute)}, {"Boss respawns", now.Add(2 * time.Hour)}, {"Last login", now.Add(-30 * time.Hour)}, {"Festival", now.Add(52 * time.Hour)}} {
			ui.Row(c).Gap(k.D(6)).Children(func() {
				pixel.Text(c, r.what, pixel.TextOptions{})
				datetime.RelativeTime(c, r.at, datetime.RelativeTimeOptions{})
			})
		}
	})
}

func demoAlarm(c *ui.Context, log func(string)) {
	s := timeState(c)
	if datetime.Alarm(c, "Alarms", s.alarms, datetime.AlarmOptions{}).Changed {
		log(fmt.Sprintf("7:30 %v, 10:00 %v", s.alarms[0].On, s.alarms[1].On))
	}
}

func demoYearView(c *ui.Context, log func(string)) {
	datetime.YearView(c, pixel.K(c).Now().Year(), datetime.YearViewOptions{})
}

func demoDaySchedule(c *ui.Context, log func(string)) {
	ui.Box(c).Width(pixel.K(c).D(220)).Children(func() {
		datetime.DaySchedule(c, []datetime.Event{
			{Start: today(c, 9, 0), Length: 2 * time.Hour, Title: "Slime patrol"},
			{Start: today(c, 10, 0), Length: time.Hour, Title: "Call mom"},
			{Start: today(c, 13, 0), Length: 90 * time.Minute, Title: "Inn lunch"},
			{Start: today(c, 16, 30), Length: time.Hour, Title: "Training"},
		}, datetime.DayScheduleOptions{})
	})
}
