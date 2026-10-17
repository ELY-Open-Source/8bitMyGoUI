package text

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var quoteMark = pixel.MustSprite(`
	.##..##.
	#...#...
	##..##..
	##..##..
`)

type QuoteOptions struct {
	Cite string
}

// Quote is a block quote with bar, mark and cite.
func Quote(c *ui.Context, s string, o QuoteOptions) ui.Element {
	k := pixel.K(c)
	return ui.Row(c).Gap(k.D(6)).AlignItems(ui.Stretch).Role(ui.RoleGroup).Children(func() {
		ui.Box(c).Width(k.D(2)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Rect(0, 0, g.W, g.H, pixel.Accent)
		})
		ui.Column(c).Grow(1).Gap(k.D(2)).Children(func() {
			ui.Box(c).Size(k.D(quoteMark.W), k.D(quoteMark.H)).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Sprite(0, 0, quoteMark, 0, pixel.Accent, 1)
			})
			Paragraph(c, s, ParagraphOptions{})
			if o.Cite != "" {
				pixel.Text(c, "— "+o.Cite, pixel.TextOptions{Slot: pixel.Faint})
			}
		})
	})
}
