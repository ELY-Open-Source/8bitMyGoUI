package datetime

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type RelativeTimeOptions struct {
	// Slot inks the text; zero means Faint.
	Slot pixel.Slot
}

// RelativeTime says how long since or until t, each minute.
func RelativeTime(c *ui.Context, t time.Time, o RelativeTimeOptions) ui.Element {
	k := pixel.K(c)
	ink := o.Slot
	if ink == 0 {
		ink = pixel.Faint
	}
	// Minutes since t; ahead of it, rounded up and negative.
	d, m := k.Now().Sub(t), 0
	if d >= 0 {
		m = int(d / time.Minute)
	} else {
		m = -int((-d + time.Minute - 1) / time.Minute)
	}
	k.At(c, t.Add(time.Duration(m+1)*time.Minute))
	return pixel.Text(c, words(k, m), pixel.TextOptions{Slot: ink})
}

// words says m minutes: ago when positive, ahead when negative.
func words(k pixel.Kit, m int) string {
	n := max(m, -m)
	h, days := n/60, n/(24*60)
	pick := func(past, future string) string {
		if m > 0 {
			return past
		}
		return future
	}
	switch {
	case n == 0:
		return k.Msg("just_now")
	case h == 0:
		return fmt.Sprintf(k.Msg(pick("min_ago", "in_min")), n)
	case days == 0:
		return fmt.Sprintf(k.Msg(pick("hours_ago", "in_hours")), h)
	case days == 1:
		return k.Msg(pick("yesterday", "tomorrow"))
	}
	return fmt.Sprintf(k.Msg(pick("days_ago", "in_days")), days)
}
