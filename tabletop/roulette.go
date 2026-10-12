package tabletop

import (
	"fmt"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type RouletteOptions struct {
	// Cue spins it: each new cue lands it on result.
	Cue int
}

const wheel = 28

// Roulette spins its wheel, on a new Cue, to result.
func Roulette(c *ui.Context, text string, segments []string, result int, o RouletteOptions) *ui.Element {
	k := pixel.K(c)
	n := len(segments)
	if n < 2 || n > 64 || result < 0 || result >= n || slices.Contains(segments, "") {
		panic(fmt.Sprintf("tabletop: Roulette lands on %d of %d segments, or one is blank", result, n))
	}
	e := ui.Column(c).Gap(k.D(2)).AlignItems(ui.Center).AlignSelf(ui.Start).Role(ui.RoleGroup).Label(text).Description(segments[result])
	// A spin: a lap plus twelve, a segment a frame.
	total := n + 12
	f := k.Shot(c, e, o.Cue, 15, total)
	under := ((result-(total-f))%n + n) % n
	return e.Children(func() {
		ui.Box(c).Size(k.D(2*wheel+1), k.D(2*wheel+6)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			cx, cy, step := wheel, wheel+5, 2*math.Pi/float64(n)
			for dy := -wheel; dy <= wheel; dy++ {
				for dx := -wheel; dx <= wheel; dx++ {
					if dx*dx+dy*dy > wheel*wheel {
						continue
					}
					// Wedges run clockwise from the top; the pointer reads one.
					t := math.Mod(math.Atan2(float64(dx), float64(-dy))+step/2+2*math.Pi, 2*math.Pi)
					i := (int(t/step)%n + under) % n
					ink := []pixel.Slot{pixel.Bad, pixel.Ink}[i%2]
					if i == 0 {
						ink = pixel.Good
					}
					g.Dot(cx+dx, cy+dy, ink)
				}
			}
			g.Circle(cx, cy, wheel, pixel.Edge)
			g.Disc(cx, cy, 4, pixel.Warn)
			g.Circle(cx, cy, 4, pixel.Edge)
			// The pointer hangs over the rim at twelve o'clock.
			for j := range 4 {
				g.HLine(cx-3+j, j, 7-2*j, pixel.Warn)
			}
		})
		pixel.Text(c, segments[under], pixel.TextOptions{})
	})
}
