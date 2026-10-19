package media

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type GalleryOptions struct {
	// Scale is the large picture's, 1–4; zero means 3.
	Scale    int
	Disabled bool
}

type GalleryResult struct {
	Element ui.Element
	Changed bool
}

// Gallery shows the chosen picture large over its thumbnails.
func Gallery(c *ui.Context, text string, pics []*pixel.Sprite, names []string, chosen *int, o GalleryOptions) GalleryResult {
	k := pixel.K(c)
	sameSize("Gallery", pics)
	if len(names) != len(pics) || *chosen < 0 || *chosen >= len(pics) {
		panic(fmt.Sprintf("media: Gallery chose %d of %d pictures with %d names", *chosen, len(pics), len(names)))
	}
	s := scale("Gallery", o.Scale, 3)
	pw, ph := pics[0].W, pics[0].H
	res := GalleryResult{}
	res.Element = ui.Column(c).Gap(k.D(4)).Children(func() {
		big := pics[*chosen]
		art(c, k, names[*chosen], pw*s+4, ph*s+4).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
			g.Sprite(2, 2, big, 0, pixel.Ink, s)
		})
		old := *chosen
		// Arrows choose at once; the cursor is the choice.
		shell, base, st, _ := pixel.Pad(c, text, len(pics), len(pics), pw+6, ph+6, chosen, o.Disabled)
		base.Description(names[*chosen])
		if res.Changed = *chosen != old; res.Changed {
			slog.Debug("8bit: gallery", "label", text, "chosen", names[*chosen])
		}
		at := *chosen
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, sp := range pics {
				x := i * (pw + 6)
				if i == at {
					g.Cursor(x, 0, pw+6, ph+6)
				}
				g.Sprite(x+3, 3, sp, 0, pixel.Ink, 1)
			}
			g.Finish(st)
		})
	})
	return res
}
