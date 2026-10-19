package files

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type QuotaMeterOptions struct {
	// Gems stand for the whole quota, 1–40; zero means 10.
	Gems int
}

var (
	gemFull = pixel.MustSprite(`
		..eeee..
		.e#ll#e.
		e#l####e
		e######e
		.e####e.
		..e##e..
		...ee...`)
	gemEmpty = pixel.MustSprite(`
		..eeee..
		.ewwwwe.
		ewwwwwwe
		ewwwwwwe
		.ewwwwe.
		..ewwe..
		...ee...`)
)

// QuotaMeter shows storage used as gems filling up.
func QuotaMeter(c *ui.Context, text string, used, quota int64, o QuotaMeterOptions) ui.Element {
	k := pixel.K(c)
	n := pixel.Or("QuotaMeter Gems", o.Gems, 10)
	if n > 40 || used < 0 || quota <= 0 {
		panic(fmt.Sprintf("files: QuotaMeter %d gems, %d of %d bytes", n, used, quota))
	}
	frac := float64(used) / float64(quota)
	// Half gems fill only when wholly used; over quota caps.
	halves := int(min(frac, 1) * float64(2*n))
	ink := pixel.Info
	switch {
	case frac >= 1:
		ink = pixel.Bad
	case frac >= 0.8:
		ink = pixel.Warn
	}
	say := fmt.Sprintf(k.Msg("used_of"), size(used), size(quota))
	w := max(n*9-1, pixel.Measure(say))
	return ui.Box(c).Size(k.D(w), k.D(8+pixel.LineHeight)).Role(ui.RoleMeter).Label(text).Description(say).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i := range n {
			x := i * 9
			g.Sprite(x, 0, gemEmpty, 0, ink, 1)
			switch {
			case halves >= 2*i+2:
				g.Sprite(x, 0, gemFull, 0, ink, 1)
			case halves == 2*i+1:
				g.Clip(x, 0, 4, 7, func() { g.Sprite(x, 0, gemFull, 0, ink, 1) })
			}
		}
		g.Text(0, 8, say, pixel.Faint, 1)
	})
}
