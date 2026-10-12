package game

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Member is a party member: portrait, HP and MP.
type Member struct {
	Name                 string
	Portrait             *pixel.Sprite
	HP, MaxHP, MP, MaxMP int
}

type PartyFramesOptions struct{}

// PartyFrames stacks member cards; a fallen member greys out.
func PartyFrames(c *ui.Context, text string, members []Member, o PartyFramesOptions) *ui.Element {
	k := pixel.K(c)
	return ui.Column(c).Gap(k.D(3)).AlignItems(ui.Start).Role(ui.RoleGroup).Label(text).Children(func() {
		for _, m := range members {
			if m.Portrait == nil {
				panic(fmt.Sprintf("game: member %q has no portrait", m.Name))
			}
			down := m.HP == 0
			ui.Row(c).Gap(k.D(4)).Padding(k.D(3)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(m.Name).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Frame(0, 0, g.W, g.H, pixel.Raised, pixel.Face)
			}).Children(func() {
				ui.Box(c).Size(k.D(28), k.D(28)).Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					g.Frame(0, 0, 28, 28, pixel.Sunken, pixel.Well)
					if down {
						g.Silhouette(2, 2, m.Portrait, 0, pixel.Shade, 2)
						return
					}
					g.Sprite(2, 2, m.Portrait, 0, pixel.Ink, 2)
				})
				ui.Column(c).Gap(k.D(1)).Children(func() {
					slot := pixel.Ink
					if down {
						slot = pixel.Faint
					}
					pixel.Text(c, m.Name, pixel.TextOptions{Slot: slot})
					HealthBar(c, "HP", m.HP, m.MaxHP, HealthBarOptions{Width: 72})
					ManaBar(c, "MP", m.MP, m.MaxMP, ManaBarOptions{Width: 72})
				})
			})
		}
	})
}
