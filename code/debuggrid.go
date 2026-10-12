package code

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DebugGridOptions struct {
	// Step is dots between lines, 2–32; zero means 8.
	Step int
}

// DebugGrid rules content with the dot grid and rulers.
func DebugGrid(c *ui.Context, o DebugGridOptions, content func()) *ui.Element {
	k := pixel.K(c)
	step := pixel.Or("DebugGrid Step", o.Step, 8)
	if step < 2 || step > 32 {
		panic(fmt.Sprintf("code: DebugGrid Step %d outside 2–32", step))
	}
	// ruler ticks every two dots; steps are long, some numbered.
	ruler := func(across bool) func(p *ui.Painter, r ui.Rect) {
		return func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Rect(0, 0, g.W, g.H, pixel.Face)
			n := g.H
			if across {
				n = g.W
			}
			for i := 0; i < n; i += 2 {
				l := 2
				if i%step == 0 {
					l = 4
				}
				if across {
					g.VLine(i, 10-l, l, pixel.Ink)
				} else {
					g.HLine(10-l, i, l, pixel.Ink)
				}
				if s := fmt.Sprint(i); i%(4*step) == 0 && i > 0 {
					if across {
						g.Micro(i+1, 0, s, pixel.Faint, 1)
					} else if pixel.MicroWidth(s) <= 6 {
						g.Micro(0, i+1, s, pixel.Faint, 1)
					}
				}
			}
		}
	}
	// The rulers stretch to the content, no further.
	return ui.Column(c).AlignSelf(ui.Start).AlignItems(ui.Stretch).Children(func() {
		ui.Row(c).Children(func() {
			ui.Box(c).Size(k.D(10), k.D(10)).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Rect(0, 0, 10, 10, pixel.Shade)
			})
			ui.Box(c).Height(k.D(10)).Grow(1).Draw(ruler(true))
		})
		ui.Row(c).AlignItems(ui.Stretch).Children(func() {
			ui.Box(c).Width(k.D(10)).Draw(ruler(false))
			ui.Box(c).Children(content).DrawOver(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				// Dotted lines every step, every other dot lit.
				for x := 0; x < g.W; x += step {
					for y := 0; y < g.H; y += 2 {
						g.Dot(x, y, pixel.Accent)
					}
				}
				for y := 0; y < g.H; y += step {
					for x := 0; x < g.W; x += 2 {
						g.Dot(x, y, pixel.Accent)
					}
				}
			})
		})
	})
}
