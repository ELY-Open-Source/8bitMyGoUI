package game

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type BossBarOptions struct {
	// Width is in dots; zero means 200.
	Width int
	// Phases are where the boss changes, as fractions: 0.5.
	Phases []float64
}

// BossBar is a wide bar under the name, phases marked.
func BossBar(c *ui.Context, name string, hp, full int, o BossBarOptions) *ui.Element {
	k := pixel.K(c)
	amount("BossBar hp", hp, full)
	for _, f := range o.Phases {
		if !(f > 0 && f < 1) {
			panic(fmt.Sprintf("game: BossBar phase %v outside 0–1", f))
		}
	}
	w := pixel.Or("BossBar Width", o.Width, 200)
	if need := icon.Skull.W + 3 + pixel.Measure(name); w < need {
		panic(fmt.Sprintf("game: BossBar Width %d is under the %d its name needs", w, need))
	}
	return ui.Box(c).Role(ui.RoleMeter).Label(name).Description(ratio(hp, full)).Size(k.D(w), k.D(pixel.LineHeight+12)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Sprite(0, 2, icon.Skull, 0, pixel.Bad, 1)
		g.Text(icon.Skull.W+3, 0, name, pixel.Ink, 1)
		meter(g, 0, pixel.LineHeight, w, 12, float64(hp)/float64(full), pixel.Bad)
		for _, f := range o.Phases {
			x := 2 + int(f*float64(w-4))
			g.VLine(x, pixel.LineHeight-2, 14, pixel.Ink)
			g.Rect(x-1, pixel.LineHeight-2, 3, 1, pixel.Ink)
		}
	})
}
