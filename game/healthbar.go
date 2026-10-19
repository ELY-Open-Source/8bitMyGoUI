package game

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type HealthBarOptions struct {
	// Width is in dots; zero means 96.
	Width int
}

type trailKey struct{}

// trail is the last hp and full, and the trail.
type trail struct {
	seen, full int
	from       float64
	at         time.Time
}

// hold is how long a trail waits before draining.
const hold = 300 * time.Millisecond

// HealthBar shows hp with a lagging trail; low HP blinks.
func HealthBar(c *ui.Context, text string, hp, full int, o HealthBarOptions) ui.Element {
	k := pixel.K(c)
	amount("HealthBar hp", hp, full)
	w := pixel.Or("HealthBar Width", o.Width, 96)
	fits("HealthBar", w, text, ratio(full, full))
	e := ui.Box(c).Role(ui.RoleMeter).Label(text).Description(ratio(hp, full))
	tr := ui.Local(e, trailKey{}, func() trail { return trail{seen: hp, full: full, from: float64(hp)} })
	// A new full makes the old trail meaningless.
	if full != tr.full {
		tr.seen, tr.full, tr.from = hp, full, float64(hp)
	}
	now := k.Now()
	// end is the trail's end: held, then draining fast.
	end := func() float64 {
		// S5: the drain moves in whole frames, fifteen a second.
		drained := max(0, now.Sub(tr.at)-hold).Truncate(time.Second / 15)
		return max(float64(tr.seen), tr.from-drained.Seconds()*float64(full))
	}
	if hp < tr.seen {
		tr.from, tr.at = end(), now
	}
	if hp > tr.seen {
		tr.from = float64(hp)
	}
	tr.seen = hp
	shown := end()
	if k.Reduce() {
		shown = float64(hp)
	}
	if shown > float64(hp) {
		k.Tick(c, 15)
	}
	// Divide, not multiply: huge HP must not overflow.
	ink := pixel.Good
	switch {
	case hp <= full/4:
		ink = pixel.Bad
	case hp <= full/2:
		ink = pixel.Warn
	}
	low := hp > 0 && hp <= full/4
	top := headH(text)
	return e.Size(k.D(w), k.D(top+8)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		head(g, w, text, ratio(hp, full))
		fill := ink
		// Low HP blinks; reduced, it holds steady.
		if low && !k.Reduce() && g.Busy(4)%2 == 1 {
			fill = pixel.Lit
		}
		x := meter(g, 0, top, w, 8, float64(hp)/float64(full), fill)
		if t := 2 + int(shown/float64(full)*float64(w-4)); t > x {
			g.Rect(x, top+2, t-x, 4, pixel.AccentLit)
		}
	})
}
