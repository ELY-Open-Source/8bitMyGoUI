package display

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type Rarity uint8

const (
	Common Rarity = iota
	Uncommon
	Rare
	Epic
	Legendary
)

// Slot is the rarity's color.
func (r Rarity) Slot() pixel.Slot {
	if r > Legendary {
		panic(fmt.Sprintf("display: unknown Rarity %d", r))
	}
	return [...]pixel.Slot{pixel.Faint, pixel.Good, pixel.Info, pixel.Magic, pixel.Warn}[r]
}

func (r Rarity) name(k pixel.Kit) string {
	return k.Msg([...]string{"common", "uncommon", "rare", "epic", "legendary"}[r])
}

var (
	gem = pixel.MustSprite(`
		...eeeee...
		..e#l#l#e..
		.e##l#l##e.
		eeeeeeeeeee
		e#ll###ss#e
		.e#l###s#e.
		..e#l#s#e..
		...e#s#e...
		....e#e....
		.....e.....
	`)
	glint = pixel.MustSprite(`
		..#..
		..#..
		##.##
		..#..
		..#..

		.....
		..#..
		.#.#.
		..#..
		.....
	`)
)

type GemOptions struct {
	Scale int
}

// Gem is a cut gem by rarity, sparkling in steps.
func Gem(c *ui.Context, r Rarity, o GemOptions) ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	ink := r.Slot()
	return art(c, k, r.name(k), (gem.W+3)*sc, (gem.H+2)*sc, func(g pixel.Canvas) {
		g.Sprite(0, 2*sc, gem, 0, ink, sc)
		if f := g.Step(4) % 4; f < 2 {
			g.Sprite(gem.W*sc-2*sc, 0, glint, f, pixel.Lit, sc)
		}
	})
}
