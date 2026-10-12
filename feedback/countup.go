package feedback

import (
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CountUpOptions struct {
	// Scale is 1–4; zero means 1.
	Scale int
	Slot  pixel.Slot
}

type countKey struct{}

// counting runs from one value to the next.
type counting struct {
	from, to int
	at       time.Time
}

// CountUp counts to a new value in eight steps.
func CountUp(c *ui.Context, value int, o CountUpOptions) *ui.Element {
	k := pixel.K(c)
	sc := pixel.Or("CountUp Scale", o.Scale, 1)
	e := ui.Box(c).Role(ui.RoleText).Label(strconv.Itoa(value))
	s := ui.Local(e, countKey{}, func() counting { return counting{from: value, to: value} })
	shown := func() int {
		step := 8
		if !k.Reduce() {
			step = since(c, k, s.at, 15, 8)
		}
		return s.from + (s.to-s.from)*step/8
	}
	if value != s.to {
		s.from, s.to, s.at = shown(), value, k.Now()
	}
	n := strconv.Itoa(shown())
	w := max(pixel.Measure(strconv.Itoa(s.from)), pixel.Measure(strconv.Itoa(s.to)))
	return e.Size(k.D(w*sc), k.D(pixel.LineHeight*sc)).Draw(func(p *ui.Painter, r ui.Rect) {
		k.Canvas(p, r).Text(0, 0, n, o.Slot, sc)
	})
}
