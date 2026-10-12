package game

import (
	"fmt"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Objective is a quest step, Have of Need (≥ 1).
type Objective struct {
	Text       string
	Have, Need int
}

// Quest is a title and its objectives.
type Quest struct {
	Title      string
	Objectives []Objective
}

// check fails on objectives past their need.
func (q Quest) check() {
	for _, ob := range q.Objectives {
		if ob.Need < 1 || ob.Have < 0 || ob.Have > ob.Need {
			panic(fmt.Sprintf("game: quest %q objective %q has %d of %d", q.Title, ob.Text, ob.Have, ob.Need))
		}
	}
}

// done is whether every objective is met.
func (q Quest) done() bool {
	for _, ob := range q.Objectives {
		if ob.Have < ob.Need {
			return false
		}
	}
	return true
}

// tick marks a 7×7 box, checked when done.
func tick(g pixel.Canvas, x, y int, done bool) {
	g.Frame(x, y, 7, 7, pixel.Flat, pixel.Well)
	if done {
		g.Line(x+1, y+3, x+2, y+5, pixel.Good)
		g.Line(x+3, y+4, x+5, y+1, pixel.Good)
	}
}
