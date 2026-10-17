package tabletop

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CoinFlipOptions struct {
	// Cue tosses it: a new cue flips, then it lands.
	Cue int
}

// The coin's width and lift, a toss's frame by frame.
var (
	turns = [...]int{24, 16, 6, 2, 6, 16, 24, 16, 6, 2, 6, 16}
	toss  = [...]int{0, 6, 11, 14, 16, 16, 14, 11, 7, 4, 2, 0}
)

const coinUp = 18

// CoinFlip tosses a coin on a new Cue.
func CoinFlip(c *ui.Context, text string, heads bool, o CoinFlipOptions) ui.Element {
	k := pixel.K(c)
	word := k.Msg("tails")
	if heads {
		word = k.Msg("heads")
	}
	e := ui.Column(c).Gap(k.D(2)).AlignItems(ui.Center).AlignSelf(ui.Start).Role(ui.RoleGroup).Label(text).Description(word)
	// A toss turns twelve frames; reduced, it lands at once.
	f := max(0, k.Shot(c, e, o.Cue, 15, len(turns)))
	return e.Children(func() {
		ui.Box(c).Size(k.D(26), k.D(26+coinUp)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			w, y, up := 24, coinUp, heads
			if f < len(turns) {
				// Each half turn shows the other face.
				w, y, up = turns[f], coinUp-toss[f], heads == (f/3%2 == 0)
			}
			for dy := -12; dy <= 12; dy++ {
				half := int(math.Round(float64(w) / 2 * math.Sqrt(1-float64(dy*dy)/144)))
				g.HLine(13-half, y+12+dy, 2*half, pixel.Edge)
				if half > 1 {
					g.HLine(14-half, y+12+dy, 2*half-2, pixel.Warn)
				}
			}
			if w == 24 {
				face := icon.Star
				if up {
					face = icon.Crown
				}
				g.Sprite(7, y+6, face, 0, pixel.Ink, 1)
			}
		})
		if f >= len(turns) {
			pixel.Text(c, word, pixel.TextOptions{})
		} else {
			ui.Box(c).Height(k.D(pixel.LineHeight))
		}
	})
}
