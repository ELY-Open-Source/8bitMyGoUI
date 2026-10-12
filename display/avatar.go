package display

import (
	"hash/fnv"
	"image"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var (
	head = pixel.MustSprite(`
		................
		................
		................
		....eeeeeeee....
		...e########e...
		..e##########e..
		..e##########e..
		..e##########e..
		..e##########e..
		..e##########e..
		..e##########e..
		..e##########e..
		...e########e...
		....e######e....
		.....eeeeee.....
		................
	`)
	hair = []*pixel.Sprite{pixel.MustSprite(`
		................
		................
		....########....
		...##########...
		..############..
		..##........##..
		..#..........#..
	`), pixel.MustSprite(`
		................
		....#.#..#.#....
		...##########...
		..############..
		..###.####.###..
		..#..........#..
	`), pixel.MustSprite(`
		................
		................
		....########....
		...##########...
		..############..
		..###......###..
		..##........##..
		..##........##..
		..##........##..
		..##........##..
		..##........##..
		...#........#...
	`)}
	eyes = []*pixel.Sprite{pixel.MustSprite(`
		.....#....#.....
		.....#....#.....
	`), pixel.MustSprite(`
		....##....##....
		................
	`), pixel.MustSprite(`
		....#.#..#.#....
		.....#....#.....
	`)}
	mouths = []*pixel.Sprite{pixel.MustSprite(`
		......#..#......
		.......##.......
	`), pixel.MustSprite(`
		................
		......####......
	`), pixel.MustSprite(`
		.......##.......
		.......##.......
	`)}
	backs  = []pixel.Slot{pixel.Info, pixel.Good, pixel.Magic, pixel.Warn, pixel.AccentLit, pixel.Accent}
	skins  = []pixel.Slot{pixel.Lit, pixel.Face, pixel.Shade, pixel.Well}
	colors = []pixel.Slot{pixel.Ink, pixel.Edge, pixel.Bad, pixel.Warn, pixel.Magic, pixel.Faint, pixel.Info}
)

type AvatarOptions struct {
	Scale int
	// Picture replaces the drawn face.
	Picture image.Image
}

// Avatar draws a face seeded by name, or a picture.
func Avatar(c *ui.Context, name string, o AvatarOptions) *ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	if o.Picture != nil {
		return Picture(c, o.Picture, PictureOptions{W: 16 * sc, H: 16 * sc, Label: name})
	}
	return art(c, k, name, 16*sc, 16*sc, func(g pixel.Canvas) { Face(g, 0, 0, name, sc) })
}

// Face paints the face of name at x, y, scaled.
func Face(g pixel.Canvas, x, y int, name string, sc int) {
	h := fnv.New32a()
	h.Write([]byte(strings.ToLower(name)))
	n := h.Sum32()
	pick := func(k int) int { v := int(n % uint32(k)); n /= uint32(k); return v }
	g.Rect(x, y, 16*sc, 16*sc, backs[pick(len(backs))])
	g.Sprite(x, y, head, 0, skins[pick(len(skins))], sc)
	g.Sprite(x, y, hair[pick(len(hair))], 0, colors[pick(len(colors))], sc)
	g.Sprite(x, y+8*sc, eyes[pick(len(eyes))], 0, pixel.Ink, sc)
	g.Sprite(x, y+11*sc, mouths[pick(len(mouths))], 0, pixel.Ink, sc)
}
