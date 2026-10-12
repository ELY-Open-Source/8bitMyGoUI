package datetime

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CountdownOptions struct{}

type CountdownResult struct {
	Element *ui.Element
	// Done is true in the frame time runs out.
	Done bool
}

type countKey struct{}

// counting is the deadline counted to, and whether it passed.
type counting struct {
	to   time.Time
	over bool
}

// Countdown shows time left to deadline, Bad near the end.
func Countdown(c *ui.Context, text string, deadline time.Time, o CountdownOptions) CountdownResult {
	k := pixel.K(c)
	e := ui.Box(c).Role(ui.RoleStatus).Label(text)
	st := ui.Local(e, countKey{}, func() counting { return counting{} })
	if !st.to.Equal(deadline) {
		*st = counting{to: deadline}
	}
	left := max(0, deadline.Sub(k.Now()))
	// Seconds round up: zero shows only at the end.
	shown := (left + time.Second - 1).Truncate(time.Second)
	res := CountdownResult{Element: e}
	if left == 0 && !st.over {
		st.over, res.Done = true, true
		slog.Debug("8bit: countdown done", "label", text)
	}
	if left > 0 {
		k.At(c, deadline.Add(-shown+time.Second))
	}
	s, ink := clock(shown), pixel.Ink
	if left < 10*time.Second {
		ink = pixel.Bad
	}
	e.Description(s).Size(k.D(pixel.Measure(s)*2), k.D(pixel.LineHeight*2)).Draw(func(p *ui.Painter, r ui.Rect) {
		k.Canvas(p, r).Text(0, 0, s, ink, 2)
	})
	return res
}
