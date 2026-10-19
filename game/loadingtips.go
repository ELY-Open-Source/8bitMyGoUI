package game

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/feedback"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type LoadingTipsOptions struct {
	// Width is in dots; zero means 200.
	Width int
}

// tipTurn is how long a tip shows.
const tipTurn = 4 * time.Second

// LoadingTips turns tips every 4 s under a progress bar.
func LoadingTips(c *ui.Context, tips []string, value float64, o LoadingTipsOptions) ui.Element {
	k := pixel.K(c)
	w := pixel.Or("LoadingTips Width", o.Width, 200)
	if len(tips) == 0 {
		panic("game: LoadingTips has no tips")
	}
	// Tips turn on the wall clock, like a loader's dots.
	now := k.Now()
	turn := now.Truncate(tipTurn)
	k.At(c, turn.Add(tipTurn))
	tip := tips[int(turn.Unix()/int64(tipTurn/time.Second))%len(tips)]
	name := k.Msg("tip")
	return ui.Column(c).Gap(k.D(6)).Width(k.D(w)).Children(func() {
		feedback.Progress(c, k.Msg("loading"), value, feedback.ProgressOptions{Width: w - 4 - pixel.MicroWidth("100%"), Percent: true})
		lines := pixel.Wrap(tip, w-tagW(name)-4, 3)
		ui.Box(c).Role(ui.RoleStatus).Label(name+": "+tip).Size(k.D(w), k.D(len(lines)*pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			tag(g, 0, 3, name, pixel.Accent)
			for i, l := range lines {
				g.Text(tagW(name)+4, i*pixel.LineHeight, l, pixel.Faint, 1)
			}
		})
	})
}
