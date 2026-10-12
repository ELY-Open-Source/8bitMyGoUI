package navigation

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// hero is the walker, two frames of stepping feet.
var hero = pixel.MustSprite(`
	..eee...
	.eaaae..
	.efffe..
	..eee...
	.eaaae..
	e.aaa.e.
	..e.e...
	.e...e..

	..eee...
	.eaaae..
	.efffe..
	..eee...
	.eaaae..
	.eaaae..
	..e.e...
	..e.e...
`)

type MapNode struct {
	Name string
	// X and Y place it on the map, in dots.
	X, Y   int
	Locked bool
}

type WorldMapOptions struct {
	// Width and Height are dots; zeros mean 240 by 140.
	Width, Height int
	Disabled      bool
}

type WorldMapResult struct {
	Element *ui.Element
	// Entered reports Enter, or a click, at the hero's node.
	Changed, Entered bool
}

type walkKey struct{}

// walk is the node the hero left and when.
type walk struct {
	from int
	at   time.Time
}

// WorldMap joins nodes with dotted paths; arrows walk the hero.
func WorldMap(c *ui.Context, text string, nodes []MapNode, paths [][2]int, at *int, o WorldMapOptions) WorldMapResult {
	k := pixel.K(c)
	if *at < 0 || *at >= len(nodes) {
		panic(fmt.Sprintf("navigation: node %d of %d", *at, len(nodes)))
	}
	w, h := pixel.Or("WorldMap Width", o.Width, 240), pixel.Or("WorldMap Height", o.Height, 140)
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Description(nodes[*at].Name).Disabled(o.Disabled)
	})
	wk := ui.Local(base, walkKey{}, func() walk { return walk{from: *at} })
	res := WorldMapResult{Element: shell}
	walkTo := func(j int) {
		if j != *at && !nodes[j].Locked {
			wk.from, wk.at, *at, res.Changed = *at, k.Now(), j, true
		}
	}
	for _, m := range []struct {
		key    ui.Key
		dx, dy int
	}{{ui.KeyLeft, -1, 0}, {ui.KeyRight, 1, 0}, {ui.KeyUp, 0, -1}, {ui.KeyDown, 0, 1}} {
		if pixel.Shortcut(base, 0, m.key) {
			if j := toward(nodes, paths, *at, m.dx, m.dy); j >= 0 {
				walkTo(j)
			}
		}
	}
	enter, space := pixel.Shortcut(base, 0, ui.KeyEnter), pixel.Shortcut(base, 0, ui.KeySpace)
	res.Entered = enter || space
	// Enter and Space are taken, so clicks come from pointers.
	if base.Clicked() {
		px, py, _ := base.PointerPosition()
		if j := nearest(nodes, int(px/k.Dot()), int(py/k.Dot())); j == *at {
			res.Entered = true
		} else if j >= 0 && linked(paths, *at, j) {
			walkTo(j)
		}
	}
	// Six steps carry the hero; Reduce Motion lands it.
	step := 6
	if !k.Reduce() {
		step = k.Since(c, wk.at, 12, 6)
	}
	a, b := nodes[wk.from], nodes[*at]
	hx, hy := a.X+(b.X-a.X)*step/6, a.Y+(b.Y-a.Y)*step/6
	shell.Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Dither(0, 0, g.W, g.H, pixel.Face, pixel.Good, 2)
		for _, e := range paths {
			trail(g, nodes[e[0]], nodes[e[1]])
		}
		for i, n := range nodes {
			fill := pixel.Accent
			if n.Locked {
				fill = pixel.Shade
			}
			g.Disc(n.X, n.Y, 4, pixel.Edge)
			g.Disc(n.X, n.Y, 3, fill)
			ink := pixel.Faint
			if i == *at {
				ink = pixel.Ink
			}
			g.Text(n.X-pixel.Measure(n.Name)/2, n.Y+6, n.Name, ink, 1)
		}
		g.Sprite(hx-4, hy-12, hero, g.Step(4)%2, pixel.Ink, 1)
		g.Finish(st)
	})
	switch {
	case res.Changed:
		slog.Debug("8bit: world map", "label", text, "node", nodes[*at].Name)
	case res.Entered:
		slog.Debug("8bit: world map enter", "label", text, "node", nodes[*at].Name)
	}
	return res
}

// toward is the linked node lying most along dx, dy.
func toward(nodes []MapNode, paths [][2]int, at, dx, dy int) int {
	best, score := -1, 0.0
	for _, e := range paths {
		j := e[1]
		if e[1] == at {
			j = e[0]
		} else if e[0] != at {
			continue
		}
		vx, vy := float64(nodes[j].X-nodes[at].X), float64(nodes[j].Y-nodes[at].Y)
		if cos := (vx*float64(dx) + vy*float64(dy)) / math.Hypot(vx, vy); cos > 0.3 && cos > score {
			best, score = j, cos
		}
	}
	return best
}

// linked reports a path between nodes a and b.
func linked(paths [][2]int, a, b int) bool {
	for _, e := range paths {
		if e == [2]int{a, b} || e == [2]int{b, a} {
			return true
		}
	}
	return false
}

// nearest is the node within 8 dots, or -1.
func nearest(nodes []MapNode, x, y int) int {
	for i, n := range nodes {
		if (n.X-x)*(n.X-x)+(n.Y-y)*(n.Y-y) <= 64 {
			return i
		}
	}
	return -1
}

// trail dots the path between two nodes.
func trail(g pixel.Canvas, a, b MapNode) {
	n := max(abs(b.X-a.X), abs(b.Y-a.Y))
	for i := 0; i <= n; i += 3 {
		g.Dot(a.X+(b.X-a.X)*i/n, a.Y+(b.Y-a.Y)*i/n, pixel.Edge)
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
