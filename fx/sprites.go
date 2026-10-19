package fx

import "github.com/ELY-Open-Source/8bitMyGoUI/pixel"

var smallCloud = pixel.MustSprite(`
.....lll......
..llllllll.l..
.llllllllllll.
llllllllllllll
fllllllllllllf
.ffffffffffff.
`)

var bigCloud = pixel.MustSprite(`
.........llll...........
....llll.llllll.........
..lllllllllllllll.llll..
.lllllllllllllllllllllll
llllllllllllllllllllllll
lllllllllllllllllllllllf
fllllllllllllllllllllllf
.fffllllllllllllllllfff.
...ffffffffffffffffff...
`)

// trees stand in the scene, one per season.
var trees = [4]*pixel.Sprite{
	pixel.MustSprite(`
.....AAAAA.....
...AAAaAAAAA...
..AAAAAAAAaAA..
.AAaAAAAAAAAAA.
.AAAAAAaAAAAAA.
AAAAAAAAAAAaAAA
AAaAAAAAAAAAAAA
AAAAAAaAAAAAAAA
.AAAAAAAAAaAAA.
..AAAaAAAAAAA..
...AAAAnAAAA...
.....AnnnA.....
......nnn......
......nnn......
......nnn......
......nnn......
.....nnnnn.....
`),
	pixel.MustSprite(`
.....ggggg.....
...ggglgggggg..
..ggllggggggg..
.gggllgggggggg.
.ggggggggggggg.
ggggggggggggggg
ggggggggggggggg
.ggggggggggggg.
.ggggggggggggg.
..ggggggggggg..
...gggnggggg...
.....ggnngg....
......nnn......
......nnn......
......nnn......
......nnn......
.....nnnnn.....
`),
	pixel.MustSprite(`
.....yyyyy.....
...yyryyyyay...
..yyyyyyryyyy..
.yyayyyyyyyyry.
.yyyyyryyyyyyy.
yyyryyyyyyayyyy
yyyyyyyyryyyyyy
.yyyayyyyyyyry.
.yyyyyyyyyyyyy.
..yyyryyyyyyy..
...yyyynyyyy...
.....yynnyy....
......nnn......
......nnn......
......nnn......
......nnn......
.....nnnnn.....
`),
	pixel.MustSprite(`
.......l.......
......lgl......
.....lgggl.....
......ggg......
....llgggll....
...lgggggggl...
....ggggggg....
...llgggggll...
..lgggggggggl..
...ggggggggg...
..llgggggggll..
.lgggggggggggl.
lgggggggggggggl
......nnn......
......nnn......
......nnn......
.....nnnnn.....
`),
}

// petal turns: face, tilt, edge, tilt back.
var petal = pixel.MustSprite(`
.A.
AAA
.a.

..A
.Aa
a..

...
AAa
...

A..
aA.
..a
`)

// leaf turns about its stem; # takes the leaf's slot.
var leaf = pixel.MustSprite(`
.###.
#####
#####
.###.
..#..

..##.
.###.
####.
.##..
..#..

..#..
..#..
..#..
..#..
..#..

.##..
.###.
.####
..##.
..#..
`)

// twinkle grows from a dot to a star.
var twinkle = pixel.MustSprite(`
.....
.....
..#..
.....
.....

.....
..#..
.###.
..#..
.....

..#..
..#..
#####
..#..
..#..
`)

// nightLit draws Lit art in Ink, bright at night.
func nightLit(g pixel.Canvas, k pixel.Kit) pixel.Canvas {
	if k.Palette().Night {
		return g.Cycled([]pixel.Slot{pixel.Lit, pixel.Ink}, 1)
	}
	return g
}
