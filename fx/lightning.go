package fx

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LightningOptions struct {
	// Cue strikes each time it changes.
	Cue int
	// Width and Height are dots; zeros mean 96 by 96.
	Width, Height int
	Label         string
}

// strikeSteps run flash, half flash, bolt, bolt, dim bolt.
const strikeSteps = 5

// Lightning flashes, then forks a bolt to the ground.
func Lightning(c *ui.Context, o LightningOptions) *ui.Element {
	k := pixel.K(c)
	w, h := size("Lightning", o.Width, o.Height, 96, 96)
	e := box(c, k, o.Label, w, h)
	step := k.Shot(c, e, o.Cue, 15, strikeSteps)
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		if step == strikeSteps {
			return
		}
		g := k.Canvas(p, r)
		g.Clip(0, 0, g.W, g.H, func() { strike(g, k, step, o.Cue) })
	})
}

// strike draws step of a strike; seed shapes the bolt.
func strike(g pixel.Canvas, k pixel.Kit, step, seed int) {
	switch step {
	case 0:
		g.Rect(0, 0, g.W, g.H, light(k))
		return
	case 1:
		g.Veil(0, 0, g.W, g.H, light(k), 6)
	}
	core, glow := pick(k, pixel.Warn, pixel.Ink), pixel.Info
	if step == 4 {
		core = pixel.AccentLit
	}
	x, y := g.W/3+rnd(seed, 140)%max(1, g.W/3), 0
	fork := 1 + rnd(seed, 141)%3
	for i := 0; y < g.H; i++ {
		nx, ny := x+rnd(seed+i, 142)%7-3, y+4+rnd(seed+i, 143)%6
		g.Line(x+1, y, nx+1, ny, glow)
		g.Line(x, y, nx, ny, core)
		// One branch forks off and dies out.
		if i == fork {
			bx, by, dir := nx, ny, 1-2*(rnd(seed, 144)%2)
			for j := range 3 {
				ex, ey := bx+dir*(3+rnd(seed+j, 145)%3), by+3+rnd(seed+j, 146)%4
				g.Line(bx, by, ex, ey, core)
				bx, by = ex, ey
			}
		}
		x, y = nx, ny
	}
}
