package chat

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Member is a name and how present they are.
type Member struct {
	Name   string
	Status display.Status
}

type MemberListOptions struct {
	// Width is in dots; zero means 140.
	Width int
}

// MemberList groups members online and offline, with status dots.
func MemberList(c *ui.Context, text string, members []Member, o MemberListOptions) *ui.Element {
	k := pixel.K(c)
	w := pixel.Or("MemberList Width", o.Width, 140)
	var on, off []Member
	for _, m := range members {
		if m.Status == display.Offline {
			off = append(off, m)
		} else {
			on = append(on, m)
		}
	}
	return ui.Column(c).Width(k.D(w)).Gap(k.D(2)).Role(ui.RoleList).Label(text).Children(func() {
		for _, g := range []struct {
			key string
			ms  []Member
			ink pixel.Slot
		}{{"online", on, pixel.Ink}, {"offline", off, pixel.Faint}} {
			if len(g.ms) == 0 {
				continue
			}
			pixel.Text(c, fmt.Sprintf("%s · %d", k.Msg(g.key), len(g.ms)), pixel.TextOptions{Slot: pixel.Faint})
			for _, m := range g.ms {
				ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Padding(0, 0, 0, k.D(4)).Children(func() {
					display.StatusDot(c, m.Status, display.StatusDotOptions{})
					pixel.Text(c, line(m.Name, w-15), pixel.TextOptions{Slot: g.ink})
				})
			}
		}
	})
}
