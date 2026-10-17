package agent

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ThinkingOrbOptions struct{}

var (
	// Swirl arms run Lit at the heart, Info, then Magic.
	swirl = []pixel.Slot{pixel.Lit, pixel.Info, pixel.Magic}
	// mist is the glow's half width, row by row.
	mist = func() map[int]int {
		half := map[int]int{}
		for _, d := range ring(15) {
			half[d[1]] = max(half[d[1]], abs(d[0]))
		}
		return half
	}()
)

// ThinkingOrb is a crystal ball that swirls while thinking.
func ThinkingOrb(c *ui.Context, text string, thinking bool, o ThinkingOrbOptions) ui.Element {
	k := pixel.K(c)
	word := k.Msg("idle")
	if thinking {
		word = k.Msg("thinking")
	}
	return ui.Box(c).Size(k.D(44), k.D(48)).Role(ui.RoleStatus).Label(text).Description(word).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Disc(22, 20, 19, pixel.Well)
		// A busy sign: it swirls on under Reduce Motion.
		if thinking {
			for y, h := range mist {
				g.Veil(22-h, 20+y, 2*h+1, 1, pixel.Magic, 4)
			}
			turn := float64(g.Busy(8)%12) * math.Pi / 6
			for arm := range 3 {
				px, py := 22, 20
				for i := 1; i <= 16; i++ {
					t := float64(i)*0.38 + turn + float64(arm)*2*math.Pi/3
					x, y := 22+int(math.Round(float64(i)*math.Cos(t))), 20+int(math.Round(float64(i)*math.Sin(t)))
					g.Line(px, py, x, y, swirl[(i-1)*3/16])
					px, py = x, y
				}
			}
			g.Disc(22, 20, 2, pixel.Lit)
		}
		g.Circle(22, 20, 19, pixel.Edge)
		g.Rect(11, 8, 4, 2, pixel.Lit)
		g.Rect(10, 10, 2, 2, pixel.Lit)
		g.Rect(13, 38, 18, 3, pixel.Edge)
		g.Rect(8, 41, 28, 6, pixel.Shade)
		g.HLine(8, 41, 28, pixel.Edge)
	})
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
