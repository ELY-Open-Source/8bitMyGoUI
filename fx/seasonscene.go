package fx

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SeasonSceneOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	Label         string
}

// SeasonScene paints a season at an hour, in its palette.
func SeasonScene(c *ui.Context, season pixel.Season, hour int, o SeasonSceneOptions) *ui.Element {
	checkHour("SeasonScene", hour)
	if season > pixel.Winter {
		panic(fmt.Sprintf("fx: SeasonScene season %d", season))
	}
	s := pixel.Current(c)
	s.Palette = pixel.Find(season, hour < 6 || hour >= 18)
	var e *ui.Element
	pixel.Scope(c, s, func() {
		k := pixel.K(c)
		w, h := size("SeasonScene", o.Width, o.Height, 160, 96)
		e = box(c, k, o.Label, w, h).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Clip(0, 0, g.W, g.H, func() { scene(g, k, season, hour, g.Step(8)) })
		})
	})
	return e
}

// scene layers sky, sun or moon, ridges, trees and weather.
func scene(g pixel.Canvas, k pixel.Kit, season pixel.Season, hour, step int) {
	night := k.Palette().Night
	sky(g, k, step/2)
	sunMoon(g, k, hour, step/2)
	if !night {
		clouds(g, step)
	}
	mountains(g, k, step, 2)
	top := g.H - g.H/6
	g.Rect(0, top, g.W, g.H-top, [4]pixel.Slot{pixel.Good, pixel.Good, pixel.Warn, light(k)}[season])
	if night {
		g.Veil(0, top, g.W, g.H-top, dark(k), 8)
	}
	tree := trees[season]
	for x := 4 + rnd(int(season), 190)%12; x < g.W-tree.W; x += 30 + rnd(x, 191)%30 {
		y := top - tree.H + 3
		if !night {
			g.Sprite(x, y, tree, 0, pixel.Ink, 1)
			continue
		}
		// By night trees stand dark; only snow catches the light.
		g.Silhouette(x, y, tree, 0, dark(k), 1)
		for ty := range tree.H {
			for tx := range tree.W {
				if s, ok := tree.At(0, tx, ty, pixel.Ink); ok && s == pixel.Lit {
					g.Dot(x+tx, y+ty, light(k))
				}
			}
		}
	}
	switch season {
	case pixel.Spring:
		petals(g, step)
	case pixel.Summer:
		if night {
			fireflies(g, step)
		}
	case pixel.Autumn:
		leaves(g, step)
	case pixel.Winter:
		snow(g, k, step)
	}
}
