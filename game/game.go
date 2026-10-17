// Package game is an RPG's HUD, menus and screens.
package game

import (
	"fmt"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// tag writes s from ink top y, micro if possible.
func tag(g pixel.Canvas, x, y int, s string, ink pixel.Slot) {
	if pixel.MicroOK(s) {
		g.Micro(x, y, s, ink, 1)
		return
	}
	g.Text(x, y-3, s, ink, 1)
}

// tagW is s's width as tag writes it.
func tagW(s string) int {
	if pixel.MicroOK(s) {
		return pixel.MicroWidth(s)
	}
	return pixel.Measure(s)
}

// tagH is s's ink height; font ink spans rows 3–14.
func tagH(s string) int {
	if pixel.MicroOK(s) {
		return pixel.MicroHeight
	}
	return 12
}

// meter frames a bar, fills part and returns its end.
func meter(g pixel.Canvas, x, y, w, h int, part float64, ink pixel.Slot) int {
	g.Frame(x, y, w, h, pixel.Sunken, pixel.Well)
	fill := int(part * float64(w-4))
	g.Rect(x+2, y+2, fill, h-4, ink)
	return x + 2 + fill
}

// amount checks v of most, as hit points are.
func amount(name string, v, most int) {
	if most <= 0 || v < 0 || v > most {
		panic(fmt.Sprintf("game: %s %d of %d", name, v, most))
	}
}

// ratio prints "v/most".
func ratio(v, most int) string { return strconv.Itoa(v) + "/" + strconv.Itoa(most) }

// fits fails when a bar is narrower than its labels.
func fits(name string, w int, text, value string) {
	if need := tagW(text) + 4 + tagW(value); w < need {
		panic(fmt.Sprintf("game: %s Width %d is under the %d its labels need", name, w, need))
	}
}

// headH is the label row's height above a HUD bar.
func headH(text string) int { return tagH(text) + 2 }

// head writes a label left and a micro value right.
func head(g pixel.Canvas, w int, text, value string) {
	tag(g, 0, 0, text, pixel.Ink)
	g.Micro(w-pixel.MicroWidth(value), headH(text)-2-pixel.MicroHeight, value, pixel.Ink, 1)
}

// nearest is the next spot toward dx, dy, or from.
func nearest(spots []Spot, from, dx, dy int) int {
	best, score := from, 0
	for i, s := range spots {
		ahead := (s.X-spots[from].X)*dx + (s.Y-spots[from].Y)*dy
		aside := abs((s.X-spots[from].X)*dy) + abs((s.Y-spots[from].Y)*dx)
		if ahead <= 0 {
			continue
		}
		// Straight ahead beats near but off to the side.
		if d := ahead + 2*aside; best == from || d < score {
			best, score = i, d
		}
	}
	return best
}

// arrows moves cur among spots by arrow keys.
func arrows(base ui.Element, spots []Spot, cur *int) {
	for _, m := range []struct {
		key    ui.Key
		dx, dy int
	}{{ui.KeyLeft, -1, 0}, {ui.KeyRight, 1, 0}, {ui.KeyUp, 0, -1}, {ui.KeyDown, 0, 1}} {
		if pixel.Shortcut(base, 0, m.key) {
			*cur = nearest(spots, *cur, m.dx, m.dy)
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type startKey struct{}

// begun is a cue, its start and its run count.
type begun struct {
	cue, run int
	at       time.Time
}

// start is e's first showing, or its last cue change.
func start(k pixel.Kit, e ui.Element, cue int) begun {
	s := ui.Local(e, startKey{}, func() begun { return begun{cue: cue, at: k.Now()} })
	if s.cue != cue {
		s.cue, s.at, s.run = cue, k.Now(), s.run+1
	}
	k.Elapsed(&s.at)
	return *s
}
