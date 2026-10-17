package agent

import (
	"fmt"
	"math"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Kind is what a log step did.
type Kind int

const (
	Thought Kind = iota
	Cast
	Hit
	Miss
	Said
)

var kinds = [...]struct {
	msg  string
	slot pixel.Slot
}{{"log_thought", pixel.Faint}, {"log_cast", pixel.Info}, {"log_hit", pixel.Good}, {"log_miss", pixel.Bad}, {"log_said", pixel.Ink}}

// LogStep is a line of the transcript.
type LogStep struct {
	Kind Kind
	Text string
}

type AgentLogOptions struct {
	// Width and Height are dots; zeros mean 240 by 120.
	Width, Height int
}

type agentLogKey struct{}

// AgentLog tells steps as a battle log, the newest typed.
func AgentLog(c *ui.Context, text, agent string, steps []LogStep, o AgentLogOptions) ui.Element {
	k := pixel.K(c)
	if agent == "" {
		panic("agent: AgentLog has no agent")
	}
	lines := make([]string, len(steps))
	for i, s := range steps {
		if s.Kind < Thought || s.Kind > Said || s.Text == "" {
			panic(fmt.Sprintf("agent: AgentLog step %d is kind %d with %q", i, s.Kind, s.Text))
		}
		// Hits and misses tell the outcome, not who.
		args := []any{agent, s.Text}
		if s.Kind == Hit || s.Kind == Miss {
			args = args[1:]
		}
		lines[i] = fmt.Sprintf(k.Msg(kinds[s.Kind].msg), args...)
	}
	w, h := pixel.Or("AgentLog Width", o.Width, 240), pixel.Or("AgentLog Height", o.Height, 120)
	box := ui.Box(c).Size(k.D(w), k.D(h))
	s := ui.Local(box, agentLogKey{}, func() ui.ScrollState { return ui.ScrollState{} })
	return box.Children(func() {
		layout.ScrollArea(c, s, layout.ScrollAreaOptions{Label: text, Follow: true}, func() {
			col := ui.Column(c).FillWidth().Padding(k.D(4)).Gap(k.D(2))
			// A new line types out, three glyphs a step.
			newest := 0
			if len(lines) > 0 {
				newest = utf8.RuneCountInString(lines[len(lines)-1])
			}
			shown := 3 * k.Shot(c, col, len(steps), 15, (newest+2)/3)
			col.Children(func() {
				for i, l := range lines {
					n := math.MaxInt
					if i == len(lines)-1 {
						n = shown
					}
					typed(c, k, l, kinds[steps[i].Kind].slot, n)
				}
			})
		}).Element.Fill()
	})
}

// typed draws s wrapped, but only its first n glyphs.
func typed(c *ui.Context, k pixel.Kit, s string, ink pixel.Slot, n int) {
	e := ui.Box(c).FillWidth().Role(ui.RoleText).Label(s)
	width, check := pixel.WrapWidth(c, e)
	lines := pixel.Wrap(s, width, 0)
	e.Height(k.D(len(lines) * pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		check(r)
		g, left := k.Canvas(p, r), n
		for i, l := range lines {
			rs := []rune(l)
			g.Text(0, i*pixel.LineHeight, string(rs[:min(left, len(rs))]), ink, 1)
			if left -= len(rs); left <= 0 {
				return
			}
		}
	})
}
