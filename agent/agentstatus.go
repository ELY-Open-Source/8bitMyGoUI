package agent

import (
	"cmp"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type AgentStatusOptions struct {
	// Portrait fits 12 dots square; zero means a robot.
	Portrait *pixel.Sprite
}

// halo is the status ring's dots from twelve o'clock, clockwise.
var halo = ring(17)

// AgentStatus rings a portrait: idle, working, blocked or done.
func AgentStatus(c *ui.Context, text string, s Status, o AgentStatusOptions) *ui.Element {
	k := pixel.K(c)
	face := portrait("AgentStatus", o.Portrait)
	word, slot := s.look(k)
	// The widest word sets the width, so it never jumps.
	w := 40
	for st := Idle; st <= Done; st++ {
		wd, _ := st.look(k)
		w = max(w, pixel.Measure(wd))
	}
	cx, cy := w/2, 18
	return ui.Box(c).Size(k.D(w), k.D(52)).Role(ui.RoleStatus).Label(text).Description(word).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		switch s {
		case Idle:
			for i, d := range halo {
				if i%2 == 0 {
					g.Dot(cx+d[0], cy+d[1], pixel.Faint)
				}
			}
			// A sleepy Z drifts up; Reduce Motion holds it.
			g.Micro(cx+13, 2-g.Step(4)%3, "Z", pixel.Faint, 1)
		case Working:
			// A busy arc chases round, even under Reduce Motion.
			at := g.Busy(8) % 12 * len(halo) / 12
			for i, d := range halo {
				ink := pixel.Shade
				if (i-at+len(halo))%len(halo) < len(halo)/4 {
					ink = pixel.Info
				}
				g.Dot(cx+d[0], cy+d[1], ink)
			}
		case Blocked:
			g.Circle(cx, cy, 17, pixel.Warn)
			g.Rect(cx+12, 0, 7, 9, pixel.Warn)
			g.Micro(cx+14, 2, "!", pixel.Ink, 1)
		case Done:
			g.Circle(cx, cy, 17, pixel.Good)
			g.Circle(cx, cy, 16, pixel.Good)
		}
		g.Sprite(cx-face.W, cy-face.H, face, 0, pixel.Ink, 2)
		g.Text((w-pixel.Measure(word))/2, 36, word, slot, 1)
	})
}

// ring is a midpoint circle's dots, clockwise from the top.
func ring(r int) [][2]int {
	var pts [][2]int
	x, y, e := r, 0, 1-r
	for x >= y {
		pts = append(pts, [2]int{x, y}, [2]int{y, x}, [2]int{-x, y}, [2]int{-y, x}, [2]int{x, -y}, [2]int{y, -x}, [2]int{-x, -y}, [2]int{-y, -x})
		y++
		if e < 0 {
			e += 2*y + 1
		} else {
			x--
			e += 2*(y-x) + 1
		}
	}
	turn := func(p [2]int) float64 {
		return math.Mod(math.Atan2(float64(p[0]), float64(-p[1]))+2*math.Pi, 2*math.Pi)
	}
	slices.SortFunc(pts, func(a, b [2]int) int { return cmp.Compare(turn(a), turn(b)) })
	return slices.Compact(pts)
}
