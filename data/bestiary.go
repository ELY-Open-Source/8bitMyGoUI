package data

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/layout"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/text"
)

// Monster is a bestiary entry; unseen ones keep their secrets.
type Monster struct {
	Sprite *pixel.Sprite
	Name   string
	Stats  [][2]string
	Notes  string
	Seen   bool
}

type BestiaryOptions struct {
	// Scale enlarges the sprite, 1–4; zero means 2.
	Scale int
}

// Bestiary is a monster card; unseen ones are silhouettes.
func Bestiary(c *ui.Context, m Monster, o BestiaryOptions) ui.Element {
	k := pixel.K(c)
	sc := pixel.Or("Bestiary Scale", o.Scale, 2)
	if m.Sprite == nil || sc > 4 {
		panic(fmt.Sprintf("data: Bestiary %q needs a Sprite and Scale 1–4", m.Name))
	}
	hide := func(s string) string {
		if m.Seen {
			return s
		}
		return k.Msg("unknown")
	}
	return layout.Card(c, layout.CardOptions{Title: hide(m.Name)}, func() {
		ui.Row(c).Gap(k.D(8)).Children(func() {
			ui.Box(c).Size(k.D(m.Sprite.W*sc+8), k.D(m.Sprite.H*sc+8)).Shrink(0).Role(ui.RoleImage).Label(hide(m.Name)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
				if m.Seen {
					g.Sprite(4, 4, m.Sprite, 0, pixel.Ink, sc)
				} else {
					g.Silhouette(4, 4, m.Sprite, 0, pixel.Shade, sc)
				}
			})
			ui.Column(c).Grow(1).Gap(k.D(2)).Children(func() {
				for _, s := range m.Stats {
					display.KeyValue(c, s[0], hide(s[1]), display.KeyValueOptions{})
				}
			})
		})
		ui.Box(c).Height(k.D(6))
		text.Paragraph(c, hide(m.Notes), text.ParagraphOptions{})
	})
}
