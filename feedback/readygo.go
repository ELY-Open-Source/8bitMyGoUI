package feedback

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ReadyGoOptions struct {
	// Cue starts the count each time it changes.
	Cue int
}

type ReadyGoResult struct {
	Element ui.Element
	// Done comes once, as GO! shows.
	Done bool
}

type goKey struct{}

// ready is a count's cue, start, and Done flag.
type ready struct {
	n    int
	at   time.Time
	said bool
}

// ReadyGo counts 3, 2, 1, GO! in big digits.
func ReadyGo(c *ui.Context, o ReadyGoOptions) ReadyGoResult {
	k := pixel.K(c)
	e := ui.Box(c)
	s := ui.Local(e, goKey{}, func() ready { return ready{n: o.Cue} })
	if s.n != o.Cue {
		*s = ready{n: o.Cue, at: k.Now()}
	}
	res := ReadyGoResult{Element: e}
	sec := k.Since(c, s.at, 10, 40) / 10
	// Done comes once, even when time jumps past GO!.
	if sec >= 3 && !s.said && !s.at.IsZero() {
		s.said, res.Done = true, true
		slog.Debug("8bit: ready go")
	}
	if sec >= 4 {
		return res
	}
	word := [...]string{"3", "2", "1", k.Msg("go")}[sec]
	e.Role(ui.RoleStatus).Label(word).Size(k.D(120), k.D(4*pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		ink := pixel.Warn
		if sec == 3 {
			ink = pixel.Good
		}
		w := pixel.Measure(word) * 4
		g.Text((g.W-w)/2+2, 2, word, pixel.Shade, 4)
		g.Text((g.W-w)/2, 0, word, ink, 4)
	})
	return res
}
