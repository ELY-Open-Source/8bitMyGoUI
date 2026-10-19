package datetime

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type TimerOptions struct {
	Disabled bool
}

type TimerResult struct {
	Element ui.Element
	// Done is true in the frame the time runs out.
	Done bool
}

type timerKey struct{}

// ticking is time left when paused, the end, the ring.
type ticking struct {
	left       time.Duration
	ends, rang time.Time
}

// Timer counts *set down; it flashes at zero.
func Timer(c *ui.Context, text string, set *time.Duration, o TimerOptions) TimerResult {
	k := pixel.K(c)
	if *set < 0 || *set >= 100*time.Minute || *set%time.Second != 0 {
		panic(fmt.Sprintf("datetime: Timer set %v is not whole seconds under 100 min", *set))
	}
	col := ui.Column(c).Gap(k.D(4)).Role(ui.RoleGroup).Label(text)
	t := ui.Local(col, timerKey{}, func() ticking { return ticking{} })
	now := k.Now()
	running := !t.ends.IsZero()
	left := t.left
	res := TimerResult{Element: col}
	if running {
		left = max(0, t.ends.Sub(now))
		k.Tick(c, 4)
		if left == 0 {
			t.ends, t.left, t.rang, res.Done, running = time.Time{}, 0, now, true, false
			slog.Debug("8bit: timer done", "label", text)
		}
	}
	// It flashes for three seconds after ringing.
	flash := !t.rang.IsZero() && k.Elapsed(&t.rang) < 3*time.Second
	if flash {
		k.Tick(c, 4)
	}
	if !running && !flash && t.left == 0 {
		left = *set
	}
	col.Children(func() {
		s := clock((left + time.Second - 1).Truncate(time.Second))
		face := ui.Box(c).Size(k.D(pixel.Measure(s)*2), k.D(pixel.LineHeight*2)).Role(ui.RoleStatus).Label(s)
		face.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			if !flash {
				g.Text(0, 0, s, pixel.Ink, 2)
				return
			}
			// Reduced, the ring holds steady in Bad.
			if k.Reduce() || g.Busy(4)%2 == 0 {
				g.Text(0, 0, s, pixel.Bad, 2)
			}
		})
		ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
			locked := o.Disabled || running
			mins, secs := int(*set/time.Minute), int(*set%time.Minute/time.Second)
			for _, f := range []struct {
				name string
				v    *int
				top  int
			}{{k.Msg("minutes"), &mins, 99}, {k.Msg("seconds"), &secs, 59}} {
				pixel.Text(c, f.name, pixel.TextOptions{Slot: pixel.Faint})
				if input.Stepper(c, f.name, f.v, input.StepperOptions{Max: f.top, Disabled: locked}).Changed {
					*set, t.left = time.Duration(mins)*time.Minute+time.Duration(secs)*time.Second, 0
					// The face above, built already, reads the new setting.
					s = clock(*set)
					face.Label(s)
				}
			}
		})
		toggle := k.Msg("start")
		if running {
			toggle = k.Msg("pause")
		}
		if input.Button(c, toggle, input.ButtonOptions{Size: input.Small, Variant: input.Primary, Disabled: o.Disabled}).Clicked {
			switch {
			case running:
				t.left, t.ends = left, time.Time{}
			case t.left > 0:
				t.ends = now.Add(t.left)
			default:
				t.ends = now.Add(*set)
			}
			t.rang = time.Time{}
			slog.Debug("8bit: timer", "label", text, "running", !running)
		}
	})
	return res
}
