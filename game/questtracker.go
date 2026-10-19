package game

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type QuestTrackerOptions struct{}

// QuestTracker lists a quest's objectives with counts, met ones ticked.
func QuestTracker(c *ui.Context, q Quest, o QuestTrackerOptions) ui.Element {
	k := pixel.K(c)
	q.check()
	w := icon.Star.W + 2 + pixel.Measure(q.Title)
	parts := make([]string, len(q.Objectives))
	for i, ob := range q.Objectives {
		parts[i] = ob.Text + " " + ratio(ob.Have, ob.Need)
		w = max(w, 11+pixel.Measure(ob.Text)+6+pixel.MicroWidth(ratio(ob.Have, ob.Need)))
	}
	return ui.Box(c).Role(ui.RoleGroup).Label(q.Title).Description(strings.Join(parts, "; ")).Size(k.D(w), k.D((1+len(q.Objectives))*pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Sprite(0, 2, icon.Star, 0, pixel.Warn, 1)
		g.Text(icon.Star.W+2, 0, q.Title, pixel.Warn, 1)
		for i, ob := range q.Objectives {
			y, met := (i+1)*pixel.LineHeight, ob.Have == ob.Need
			tick(g, 1, y+5, met)
			ink := pixel.Ink
			if met {
				ink = pixel.Faint
			}
			g.Text(11, y, ob.Text, ink, 1)
			if ob.Need > 1 {
				n := ratio(ob.Have, ob.Need)
				g.Micro(w-pixel.MicroWidth(n), y+6, n, ink, 1)
			}
		}
	})
}
