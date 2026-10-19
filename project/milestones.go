package project

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Milestone is a stop on the way, by date.
type Milestone struct {
	Title string
	Date  time.Time
	Done  bool
}

type MilestonesOptions struct {
	// Width is in dots; zero means 260.
	Width int
}

var (
	// A flag in two frames, to wave.
	flag = pixel.MustSprite(`
		.ee.........
		.e####......
		.e######....
		.e########..
		.e######....
		.e####......
		.e..........
		.e..........
		.e..........
		.e..........
		.e..........
		eee.........

		.ee.........
		.e##........
		.e#####.....
		.e#######...
		.e#######...
		.e#####.....
		.e##........
		.e..........
		.e..........
		.e..........
		.e..........
		eee.........`)
	hero = pixel.MustSprite(`
		..yyyy..
		.yyyyyy.
		.ylkllk.
		..llll..
		.aaaaaa.
		laaaaaal
		.aaaaaa.
		.ee..ee.
		.ee..ee.
		.kk..kk.`)
)

// Milestones plants flags on a path; a hero marks today.
func Milestones(c *ui.Context, text string, ms []Milestone, o MilestonesOptions) ui.Element {
	k := pixel.K(c)
	if len(ms) == 0 {
		panic("project: Milestones has none")
	}
	for i, m := range ms {
		if m.Title == "" || m.Date.IsZero() || i > 0 && m.Date.Before(ms[i-1].Date) {
			panic(fmt.Sprintf("project: Milestones %q is untitled, undated or out of order", m.Title))
		}
	}
	w := pixel.Or("Milestones Width", o.Width, 260)
	first, span := ms[0].Date, ms[len(ms)-1].Date.Sub(ms[0].Date)
	// x places a date on the path; one sits mid-way.
	x := func(t time.Time) int {
		if span == 0 {
			return w / 2
		}
		f := max(0, min(1, float64(t.Sub(first))/float64(span)))
		return 8 + int(f*float64(w-16))
	}
	next, gap := len(ms), w
	said := make([]string, len(ms))
	for i, m := range ms {
		state := ""
		switch {
		case m.Done:
			state = ", " + k.Msg("done")
		case next == len(ms):
			next, state = i, ", "+k.Msg("next")
		}
		said[i] = m.Title + " " + m.Date.Format("01-02") + state
		if i > 0 {
			gap = min(gap, x(m.Date)-x(ms[i-1].Date))
		}
	}
	// Titles alternate rows, so each gets two gaps.
	room := max(24, min(w/2, 2*gap-4))
	now := k.Now()
	return ui.Box(c).Size(k.D(w), k.D(58)).Role(ui.RoleImage).Label(text).Description(strings.Join(said, "; ")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for px := 2; px < w-2; px += 2 {
			g.Dot(px, 16, pixel.Shade)
		}
		wave := g.Step(4) % 2
		for i, m := range ms {
			fx, ink, frame := x(m.Date), pixel.Faint, 0
			switch {
			case m.Done:
				ink = pixel.Good
			case i == next:
				ink, frame = pixel.Accent, wave
			}
			g.Sprite(fx-1, 4, flag, frame, ink, 1)
			d := m.Date.Format("01-02")
			g.Micro(fx-pixel.MicroWidth(d)/2, 19, d, pixel.Faint, 1)
			t := pixel.Wrap(m.Title, room, 1)[0]
			g.Text(max(0, min(w-pixel.Measure(t), fx-pixel.Measure(t)/2)), 26+i%2*16, t, pixel.Ink, 1)
		}
		// The hero bobs where today falls, unless motion is reduced.
		g.Sprite(max(0, min(w-8, x(now)-4)), 6-wave, hero, 0, pixel.Ink, 1)
	})
}
