package feedback

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type LoadingScreenOptions struct {
	// Tips take turns, four seconds each.
	Tips []string
	// Width is in dots; zero means 240.
	Width int
}

// LoadingScreen shows NOW LOADING, stepping dots and a tip.
func LoadingScreen(c *ui.Context, o LoadingScreenOptions) ui.Element {
	k := pixel.K(c)
	w := pixel.Or("LoadingScreen Width", o.Width, 240)
	title := k.Msg("loading")
	return ui.Box(c).Role(ui.RoleStatus).Label(title).Size(k.D(w), k.D(2*pixel.LineHeight+8+2*pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		n := g.Busy(4)
		line := title + strings.Repeat(".", n%4)
		x := (w - pixel.Measure(title+"...")*2) / 2
		g.Text(x+2, 2, line, pixel.Shade, 2)
		g.Text(x, 0, line, pixel.Accent, 2)
		if len(o.Tips) == 0 {
			return
		}
		tip := o.Tips[n/16%len(o.Tips)]
		for i, l := range pixel.Wrap(tip, w, 2) {
			g.Text((w-pixel.Measure(l))/2, 2*pixel.LineHeight+8+i*pixel.LineHeight, l, pixel.Faint, 1)
		}
	})
}
