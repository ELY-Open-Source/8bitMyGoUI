package fx

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ExplosionOptions struct {
	// Cue sets off a blast each time it changes.
	Cue int
	// Size is the square's side in dots; zero means 48.
	Size  int
	Label string
}

const blastSteps = 8

// debrisSlots color the bits thrown out, in turn.
var debrisSlots = []pixel.Slot{pixel.Warn, pixel.Bad, pixel.Faint}

// Explosion bursts, breaks into smoke and throws debris.
func Explosion(c *ui.Context, o ExplosionOptions) *ui.Element {
	k := pixel.K(c)
	n := pixel.Or("Explosion Size", o.Size, 48)
	e := box(c, k, o.Label, n, n)
	step := k.Shot(c, e, o.Cue, 12, blastSteps)
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		if step == blastSteps {
			return
		}
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { blast(g, k, step, o.Cue) })
	})
}

// blast sizes rings and puffs in sixteenths of the box.
func blast(g pixel.Canvas, k pixel.Kit, step, seed int) {
	cx, cy, u := g.W/2, g.H/2, max(1, g.W/16)
	hot := light(k)
	switch step {
	case 0:
		g.Disc(cx, cy, 2*u, pixel.Warn)
		g.Disc(cx, cy, u, hot)
	case 1, 2:
		g.Disc(cx, cy, (3+step)*u, pixel.Bad)
		g.Disc(cx, cy, 3*u, pixel.Warn)
		g.Disc(cx, cy, (3-step)*u, hot)
	default:
		// The fireball breaks into puffs that rise and thin out.
		for i := range 6 {
			a := float64(i)*math.Pi/3 + float64(rnd(seed, 180+i)%40-20)*math.Pi/180
			d := float64((step + 1) * u)
			x, y := cx+int(d*math.Cos(a)), cy+int(d*math.Sin(a))-(step-3)*u
			if r := max(1, (8-step)*u/2); step < 6 {
				g.Disc(x, y, r, pixel.Faint)
			} else {
				g.Circle(x, y, r, pixel.Faint)
			}
		}
		if step < 5 {
			g.Disc(cx, cy, (5-step)*u, pixel.Bad)
			g.Disc(cx, cy, (4-step)*u, pixel.Warn)
		}
	}
	// Debris flies out and falls.
	t := float64(step)
	for i := range 10 {
		a := float64(rnd(seed*16+i, 181)%360) * math.Pi / 180
		v := float64(u) * (1.5 + float64(rnd(seed*16+i, 182)%3)/2)
		x := cx + int(v*t*math.Cos(a))
		y := cy + int(v*t*math.Sin(a)+0.4*float64(u)*t*t)
		g.Rect(x, y, max(1, u/2), max(1, u/2), debrisSlots[i%len(debrisSlots)])
	}
}
