package data

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/text"
)

// Activity is one entry of a feed.
type Activity struct {
	Icon       *pixel.Sprite
	Text, Time string
}

type FeedOptions struct{}

// Feed lists activity: icon, text and time, dotted apart.
func Feed(c *ui.Context, entries []Activity, o FeedOptions) ui.Element {
	k := pixel.K(c)
	return ui.Column(c).FillWidth().Children(func() {
		for i, a := range entries {
			last := i == len(entries)-1
			ui.Row(c).FillWidth().Gap(k.D(6)).Padding(k.D(4), 0).Draw(func(p *ui.Painter, r ui.Rect) {
				if !last {
					g := k.Canvas(p, r)
					g.Checker(0, g.H-1, g.W, 1, pixel.Edge)
				}
			}).Children(func() {
				ui.Box(c).Size(k.D(12), k.D(12)).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
					if a.Icon != nil {
						k.Canvas(p, r).Sprite(0, 0, a.Icon, 0, pixel.Accent, 1)
					}
				})
				ui.Column(c).Grow(1).Children(func() {
					text.Paragraph(c, a.Text, text.ParagraphOptions{})
					pixel.Text(c, a.Time, pixel.TextOptions{Slot: pixel.Faint})
				})
			})
		}
	})
}
