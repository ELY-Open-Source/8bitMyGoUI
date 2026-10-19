package datetime

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// AlarmEntry is an alarm: time, days, and if it rings.
type AlarmEntry struct {
	Hour, Minute int
	// Days ring on, Sunday first.
	Days [7]bool
	On   bool
}

type AlarmOptions struct {
	Disabled bool
}

type AlarmResult struct {
	Element ui.Element
	Changed bool
}

// Alarm lists alarms: time, days and a lever each.
func Alarm(c *ui.Context, text string, alarms []AlarmEntry, o AlarmOptions) AlarmResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(4)).Role(ui.RoleList).Label(text)
	res := AlarmResult{Element: col}
	names := weekdays(k)
	col.Children(func() {
		for i := range alarms {
			a := &alarms[i]
			if a.Hour < 0 || a.Hour > 23 || a.Minute < 0 || a.Minute > 59 {
				panic(fmt.Sprintf("datetime: alarm %02d:%02d is no time", a.Hour, a.Minute))
			}
			at := fmt.Sprintf("%02d:%02d", a.Hour, a.Minute)
			ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Role(ui.RoleListItem).Label(at).Children(func() {
				ink := pixel.Faint
				if a.On {
					ink = pixel.Ink
				}
				pixel.Text(c, at, pixel.TextOptions{Slot: ink, Scale: 2})
				ui.Row(c).Gap(k.D(3)).Children(func() {
					for d, n := range names {
						day := pixel.Faint
						if a.Days[d] {
							day = pixel.Accent
						}
						pixel.Text(c, n, pixel.TextOptions{Slot: day})
					}
				})
				if input.Switch(c, "", &a.On, input.SwitchOptions{Label: at, Disabled: o.Disabled}).Changed {
					res.Changed = true
					slog.Debug("8bit: alarm", "label", text, "at", at, "on", a.On)
				}
			})
		}
	})
	return res
}
