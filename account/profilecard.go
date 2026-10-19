package account

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/layout"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/text"
)

// Stat is a number out of its most, as HP.
type Stat struct {
	Name       string
	Value, Max int
}

// Profile is a player as their card shows them.
type Profile struct {
	Name, Title string
	Level       int
	Stats       []Stat
}

type ProfileCardOptions struct {
	// Width is in dots; zero means 180.
	Width int
}

// ProfileCard shows a player: face, name, level, title and stats.
func ProfileCard(c *ui.Context, p Profile, o ProfileCardOptions) ui.Element {
	k := pixel.K(c)
	if p.Level < 1 {
		panic(fmt.Sprintf("account: ProfileCard level %d", p.Level))
	}
	w := pixel.Or("ProfileCard Width", o.Width, 180)
	// Columns fit the widest name and value, so bars align.
	nw, vw := 0, 0
	for _, s := range p.Stats {
		if s.Max < 1 || s.Value < 0 || s.Value > s.Max {
			panic(fmt.Sprintf("account: ProfileCard stat %+v", s))
		}
		nw, vw = max(nw, pixel.MicroWidth(s.Name)), max(vw, pixel.MicroWidth(fmt.Sprintf("%d/%d", s.Value, s.Max)))
	}
	bx, bw := nw+4, w-12-nw-vw-8
	return layout.Panel(c, layout.PanelOptions{Style: pixel.Raised, Padding: 6}, func() {
		ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
			display.Avatar(c, p.Name, display.AvatarOptions{Scale: 2})
			ui.Column(c).Gap(k.D(2)).Children(func() {
				text.ShadowText(c, p.Name, text.ShadowTextOptions{})
				caption(c, k, fmt.Sprintf("LV %d", p.Level), pixel.Accent)
				pixel.Text(c, p.Title, pixel.TextOptions{Slot: pixel.Faint})
			})
		})
		for _, s := range p.Stats {
			v := fmt.Sprintf("%d/%d", s.Value, s.Max)
			ui.Box(c).Size(k.D(w-12), k.D(9)).Role(ui.RoleMeter).Label(s.Name).Description(v).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Micro(0, 2, s.Name, pixel.Faint, 1)
				g.Frame(bx, 1, bw, 7, pixel.Sunken, pixel.Well)
				// Floats, so huge stats cannot overflow.
				if n := int(float64(bw-4) * float64(s.Value) / float64(s.Max)); n > 0 {
					g.Rect(bx+2, 3, n, 3, pixel.Good)
				}
				g.Micro(g.W-pixel.MicroWidth(v), 2, v, pixel.Ink, 1)
			})
		}
	}).Width(k.D(w)).Gap(k.D(4)).Role(ui.RoleGroup).Label(p.Name)
}
