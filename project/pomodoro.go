package project

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PomodoroOptions struct {
	// Work and Rest; zeros mean 25 and 5 minutes.
	Work, Rest time.Duration
	Disabled   bool
}

type PomodoroResult struct {
	Element *ui.Element
	// Changed is START, PAUSE or RESET; Rang ends a phase.
	Changed, Rang bool
}

type pomoKey struct{}

// pomo is the phase and its time, banked while paused.
type pomo struct {
	resting bool
	left    time.Duration
	since   time.Time
}

// The tomato is awake at work and asleep at rest.
var tomato = pixel.MustSprite(`
		................
		.......gg.......
		....gg.gg.gg....
		.....gggggg.....
		...errgggggre...
		..errrrggrrrre..
		.errlrrrrrrrrre.
		.erllrrrrrrrrre.
		erlrrkrrrrkrrrre
		errrrkrrrrkrrrre
		errrrrrrrrrrrrre
		errrrrrkkrrrrrre
		.errrrrrrrrrrre.
		.errrrrrrrrrrre.
		..eerrrrrrrree..
		....eeeeeeee....

		................
		.......gg.......
		....gg.gg.gg....
		.....gggggg.....
		...errgggggre...
		..errrrggrrrre..
		.errlrrrrrrrrre.
		.erllrrrrrrrrre.
		erlrrrrrrrrrrrre
		errrkkrrrkkrrrre
		errrrrrrrrrrrrre
		errrrrrkkrrrrrre
		.errrrrrrrrrrre.
		.errrrrrrrrrrre.
		..eerrrrrrrree..
		....eeeeeeee....`)

// Pomodoro times work, then rest, beside a tomato.
func Pomodoro(c *ui.Context, text string, o PomodoroOptions) PomodoroResult {
	k := pixel.K(c)
	work, rest := o.Work, o.Rest
	if work == 0 {
		work = 25 * time.Minute
	}
	if rest == 0 {
		rest = 5 * time.Minute
	}
	if work < time.Second || rest < time.Second || work > day || rest > day {
		panic(fmt.Sprintf("project: Pomodoro work %v, rest %v", work, rest))
	}
	row := ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(text)
	p := ui.Local(row, pomoKey{}, func() pomo { return pomo{left: work} })
	now := k.Now()
	left, running := p.left, !p.since.IsZero()
	if running {
		left -= k.Elapsed(&p.since)
	}
	var res PomodoroResult
	if running && left <= 0 {
		// The phase rings; the next waits for START.
		p.resting = !p.resting
		p.left, p.since, res.Rang = work, time.Time{}, true
		if p.resting {
			p.left = rest
		}
		left, running = p.left, false
		slog.Debug("8bit: pomodoro rang", "label", text, "resting", p.resting)
	}
	if running {
		// Ask for the moment the seconds turn.
		k.At(c, now.Add(left-(left-1)/time.Second*time.Second))
	}
	phase, frame := k.Msg("focus"), 0
	if p.resting {
		phase, frame = k.Msg("rest"), 1
	}
	res.Element = row.Children(func() {
		ui.Box(c).Size(k.D(32), k.D(32)).Role(ui.RoleImage).Label(phase).Draw(func(pt *ui.Painter, r ui.Rect) {
			g := k.Canvas(pt, r)
			// Running, it hops a dot each second, unless reduced.
			dy := 0
			if running && g.Step(4)%4 == 0 {
				dy = 1
			}
			g.Sprite(0, 1-dy, tomato, frame, pixel.Ink, 2)
		})
		ui.Column(c).Gap(k.D(2)).Children(func() {
			pixel.Text(c, phase, pixel.TextOptions{Slot: pixel.Faint})
			// Seconds round up, so 00:00 shows only at the end.
			pixel.Text(c, hms(left+time.Second-1), pixel.TextOptions{Scale: 2})
			ui.Row(c).Gap(k.D(4)).Children(func() {
				label := k.Msg("start")
				if running {
					label = k.Msg("pause")
				}
				if input.Button(c, label, input.ButtonOptions{Size: input.Small, Variant: input.Primary, Disabled: o.Disabled}).Clicked {
					if running {
						p.left, p.since = left, time.Time{}
					} else {
						p.since = now
					}
					res.Changed = true
					slog.Debug("8bit: pomodoro", "label", text, "running", !p.since.IsZero())
				}
				if input.Button(c, k.Msg("reset"), input.ButtonOptions{Size: input.Small, Disabled: o.Disabled}).Clicked {
					*p, res.Changed = pomo{left: work}, true
					slog.Debug("8bit: pomodoro reset", "label", text)
				}
			})
		})
	})
	return res
}
