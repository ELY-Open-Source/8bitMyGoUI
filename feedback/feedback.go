// Package feedback tells the player what is going on.
package feedback

import (
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// pingpong runs n back and forth over 0…span.
func pingpong(n, span int) int {
	if span <= 0 {
		return 0
	}
	n %= 2 * span
	if n > span {
		return 2*span - n
	}
	return n
}

// level is Good, then Warn past warn, Bad past bad.
func level(v, warn, bad float64) pixel.Slot {
	switch {
	case bad != 0 && v > bad:
		return pixel.Bad
	case warn != 0 && v > warn:
		return pixel.Warn
	}
	return pixel.Good
}

// fmtNum prints a number without trailing zeros.
func fmtNum(v float64) string { return strconv.FormatFloat(v, 'g', 4, 64) }

type cueKey struct{}

// cued is a cue's last value and when it changed.
type cued struct {
	n  int
	at time.Time
}

// cue is when n last changed; zero until it does.
func cue(e *ui.Element, k pixel.Kit, n int) time.Time {
	s := ui.Local(e, cueKey{}, func() cued { return cued{n: n} })
	if s.n != n {
		s.n, s.at = n, k.Now()
	}
	return s.at
}

// since counts steps at fps from t, up to done.
func since(c *ui.Context, k pixel.Kit, t time.Time, fps, done int) int {
	if t.IsZero() {
		return done
	}
	n := int(k.Now().Sub(t) * time.Duration(fps) / time.Second)
	if n < done {
		k.Tick(c, fps)
	}
	return min(n, done)
}
