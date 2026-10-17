package agent

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// PlanStep is an objective; the first undone is under way.
type PlanStep struct {
	Text string
	Done bool
}

type PlanQuestOptions struct {
	// Width is in dots; zero means 200.
	Width int
}

// PlanQuest shows a plan's steps as quest objectives.
func PlanQuest(c *ui.Context, title string, steps []PlanStep, o PlanQuestOptions) ui.Element {
	k := pixel.K(c)
	if title == "" || len(steps) == 0 {
		panic(fmt.Sprintf("agent: PlanQuest %q has %d steps", title, len(steps)))
	}
	w := pixel.Or("PlanQuest Width", o.Width, 200)
	if need := max(64, pixel.Measure(k.Msg("quest_complete"))); w < need {
		panic(fmt.Sprintf("agent: PlanQuest Width %d is under %d", w, need))
	}
	cur, done, h := len(steps), 0, pixel.LineHeight+4
	lines, said := make([][]string, len(steps)), make([]string, len(steps))
	for i, s := range steps {
		if s.Text == "" {
			panic(fmt.Sprintf("agent: PlanQuest step %d has no text", i))
		}
		lines[i] = pixel.Wrap(s.Text, w-16, 0)
		h += len(lines[i]) * pixel.LineHeight
		said[i] = s.Text
		switch {
		case s.Done:
			done++
			said[i] += ", " + k.Msg("done")
		case cur == len(steps):
			cur = i
			said[i] += ", " + k.Msg("doing")
		}
	}
	if cur == len(steps) {
		h += pixel.LineHeight + 2
	}
	count := fmt.Sprintf("%d/%d", done, len(steps))
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(title).Description(count + "; " + strings.Join(said, "; ")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		cw := pixel.MicroWidth(count)
		g.Sprite(0, 2, icon.Scroll, 0, pixel.Warn, 1)
		g.Text(16, 0, pixel.Wrap(title, w-20-cw, 1)[0], pixel.Ink, 1)
		g.Micro(w-cw, 6, count, pixel.Ink, 1)
		g.Checker(0, pixel.LineHeight+1, w, 1, pixel.Edge)
		y := pixel.LineHeight + 4
		for i, s := range steps {
			ink := pixel.Ink
			switch {
			case s.Done:
				ink = pixel.Faint
				g.Sprite(0, y+2, icon.Check, 0, pixel.Good, 1)
			case i == cur:
				g.Sprite(0, y+2, icon.TriangleRight, 0, pixel.Accent, 1)
			default:
				g.Rect(2, y+5, 8, 8, pixel.Edge)
				g.Rect(3, y+6, 6, 6, pixel.Well)
			}
			for j, l := range lines[i] {
				g.Text(16, y+j*pixel.LineHeight, l, ink, 1)
			}
			y += len(lines[i]) * pixel.LineHeight
		}
		if cur == len(steps) {
			q := k.Msg("quest_complete")
			g.Text((w-pixel.Measure(q))/2, y+2, q, pixel.Good, 1)
		}
	})
}
