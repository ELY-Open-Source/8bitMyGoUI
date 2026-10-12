package media

import (
	"fmt"
	"math"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CassetteOptions struct {
	// Title is written on the label in the micro font.
	Title   string
	Playing bool
	// Progress, 0 to 1, is tape on the right reel.
	Progress float64
}

// Cassette turns its reels while playing; tape winds right.
func Cassette(c *ui.Context, text string, o CassetteOptions) *ui.Element {
	k := pixel.K(c)
	if !(o.Progress >= 0 && o.Progress <= 1) {
		panic(fmt.Sprintf("media: Cassette Progress %v outside 0–1", o.Progress))
	}
	title := strings.ToUpper(o.Title)
	if !pixel.MicroOK(title) || pixel.MicroWidth(title) > 52 {
		panic(fmt.Sprintf("media: Cassette Title %q does not fit the label", o.Title))
	}
	e := art(c, k, text, 64, 40)
	frame := turned(c, k, e, o.Playing, 8) % 2
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		body, line := dark(k)
		g.Frame(0, 0, 64, 40, pixel.Raised, pixel.Face)
		g.Rect(4, 3, 56, 26, pixel.Lit)
		g.Micro(7, 5, title, pixel.Ink, 1)
		g.Rect(4, 11, 56, 2, pixel.Accent)
		g.Rect(14, 14, 36, 13, body)
		// Packs trade tape as it plays.
		wound := int(math.Round(o.Progress * 3))
		for i, cx := range []int{23, 40} {
			pack := 3 + wound
			if i == 0 {
				pack = 6 - wound
			}
			g.Disc(cx, 20, pack, line)
			g.Disc(cx, 20, 2, pixel.Lit)
			for _, d := range [2][4][2]int{{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}, {{1, 1}, {-1, -1}, {1, -1}, {-1, 1}}}[frame] {
				g.Dot(cx+d[0], 20+d[1], body)
			}
		}
		// The head opening, with its pins.
		for y := range 8 {
			g.HLine(18-y/2, 32+y, 28+y, pixel.Shade)
		}
		for _, x := range []int{21, 30, 33, 42} {
			g.Rect(x, 35, 2, 2, body)
		}
		for _, x := range []int{4, 59} {
			g.HLine(x-1, 34, 3, pixel.Shade)
			g.VLine(x, 33, 3, pixel.Shade)
		}
	})
}
