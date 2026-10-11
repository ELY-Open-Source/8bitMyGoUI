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
	return art(c, k, name, 16*sc, 16*sc, func(g pixel.Canvas) { drawFace(g, 0, name, sc) })
}

// drawFace paints the face of name at x, scaled.
func drawFace(g pixel.Canvas, x int, name string, sc int) {
	h := fnv.New32a()
	h.Write([]byte(strings.ToLower(name)))
	n := h.Sum32()
	pick := func(k int) int { v := int(n % uint32(k)); n /= uint32(k); return v }
	g.Rect(x, 0, 16*sc, 16*sc, backs[pick(len(backs))])
	g.Sprite(x, 0, head, 0, skins[pick(len(skins))], sc)
	g.Sprite(x, 0, hair[pick(len(hair))], 0, colors[pick(len(colors))], sc)
	g.Sprite(x, 8*sc, eyes[pick(len(eyes))], 0, pixel.Ink, sc)
	g.Sprite(x, 11*sc, mouths[pick(len(mouths))], 0, pixel.Ink, sc)
}

type AvatarGroupOptions struct {
	// Max faces shown before +N; zero means 4.
	Max   int
	Scale int
}

// AvatarGroup overlaps faces and counts the rest in a chip.
func AvatarGroup(c *ui.Context, names []string, o AvatarGroupOptions) *ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	shown := names[:min(len(names), orDef("AvatarGroup Max", o.Max, 4))]
	rest := len(names) - len(shown)
	more := ""
	if rest > 0 {
		more = "+" + itoa(rest)
	}
	step := 12 * sc
	w := max(0, len(shown)-1)*step + 16*sc
	if more != "" {
		w += pixel.MicroWidth(more) + 8
	}
	return art(c, k, strings.Join(names, ", "), w, 16*sc, func(g pixel.Canvas) {
		for i, name := range shown {
			x := i * step
			drawFace(g, x, name, sc)
			if i > 0 {
				g.VLine(x, 0, 16*sc, pixel.Back)
			}
		}
		if more != "" {
			x := len(shown)*step + 4*sc
			cw := pixel.MicroWidth(more) + 4
			g.Frame(x, (16*sc-9)/2, cw, 9, pixel.Notched, pixel.Face)
			g.Micro(x+2, (16*sc-9)/2+2, more, pixel.Ink, 1)
		}
	})
}
