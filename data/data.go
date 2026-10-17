// Package data shows collections: lists, tables, trees, grids.
package data

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var (
	pointer = pixel.MustSprite(`
		#..
		##.
		###
		##.
		#..
	`)
	up = pixel.MustSprite(`
		..#..
		.###.
		#####
	`)
	down = pixel.MustSprite(`
		#####
		.###.
		..#..
	`)
)

// well sinks box and finishes it as st says.
func well(c *ui.Context, box ui.Element, st pixel.State) {
	k := pixel.K(c)
	box.Padding(k.D(2)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
	}).DrawOver(func(p *ui.Painter, r ui.Rect) {
		k.Canvas(p, r).Finish(st)
	})
}

// veil covers an off list; the wheel stops on it.
func veil(c *ui.Context, off bool) {
	if off {
		ui.Box(c).Absolute().Top(0).Left(0).Fill()
	}
}

// row is a shown tree row, with its parent.
type row[K comparable] struct {
	key    K
	depth  int
	branch bool
	parent int
}

// flatten lists the rows shown under roots, inside open branches.
func flatten[K comparable](roots []K, children func(K) []K, opened map[K]bool) []row[K] {
	var out []row[K]
	var walk func(items []K, depth, parent int)
	walk = func(items []K, depth, parent int) {
		for _, it := range items {
			kids := children(it)
			out = append(out, row[K]{it, depth, kids != nil, parent})
			if opened[it] && kids != nil {
				walk(kids, depth+1, len(out)-1)
			}
		}
	}
	walk(roots, 0, -1)
	return out
}

// steer: Right opens or descends; Left closes or climbs.
func steer[K comparable](rows []row[K], at int, opened map[K]bool, right bool) (to int, toggled bool) {
	r := rows[at]
	switch {
	case right && r.branch && !opened[r.key]:
		opened[r.key] = true
		return at, true
	case right && r.branch && at+1 < len(rows):
		return at + 1, false
	case !right && r.branch && opened[r.key]:
		delete(opened, r.key)
		return at, true
	case !right && r.parent >= 0:
		return r.parent, false
	}
	return at, false
}
