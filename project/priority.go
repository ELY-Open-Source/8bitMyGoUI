package project

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PriorityOptions struct {
	// Scale is 1–4; zero means 1.
	Scale int
}

// Bars rise with the level; # dots take its ink.
var bars = [3]*pixel.Sprite{pixel.MustSprite(`
		............
		............
		.........nnn
		.........nnn
		.........nnn
		.....nnn.nnn
		.....nnn.nnn
		.....nnn.nnn
		.###.nnn.nnn
		.###.nnn.nnn
		.###.nnn.nnn
		.###.nnn.nnn`), pixel.MustSprite(`
		............
		............
		.........nnn
		.........nnn
		.........nnn
		.....###.nnn
		.....###.nnn
		.....###.nnn
		.###.###.nnn
		.###.###.nnn
		.###.###.nnn
		.###.###.nnn`), pixel.MustSprite(`
		............
		............
		.........###
		.........###
		.........###
		.....###.###
		.....###.###
		.....###.###
		.###.###.###
		.###.###.###
		.###.###.###
		.###.###.###`)}

var levelInks = [...]pixel.Slot{pixel.Good, pixel.Warn, pixel.Bad, pixel.Bad}

// Priority is a level's icon; critical is on fire.
func Priority(c *ui.Context, l Level, o PriorityOptions) ui.Element {
	k := pixel.K(c)
	say := l.word(k)
	s := pixel.Or("Priority Scale", o.Scale, 1)
	if s > 4 {
		panic(fmt.Sprintf("project: Priority Scale %d is past 4", s))
	}
	return ui.Box(c).Size(k.D(12*s), k.D(12*s)).Role(ui.RoleImage).Label(say).Draw(func(p *ui.Painter, r ui.Rect) {
		mark(k.Canvas(p, r), 0, 0, l, s)
	})
}

// mark draws a level's icon at x, y.
func mark(g pixel.Canvas, x, y int, l Level, s int) {
	art := icon.Flame
	if l < Critical {
		art = bars[l]
	}
	g.Sprite(x, y, art, 0, levelInks[l], s)
}
