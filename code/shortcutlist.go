package code

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

// Shortcut is an action and its keys, pressed together.
type Shortcut struct {
	Action string
	Keys   []string
}

type ShortcutListOptions struct {
	// Width is in dots; zero means 220.
	Width int
}

// ShortcutList pairs actions with their keycaps.
func ShortcutList(c *ui.Context, label string, items []Shortcut, o ShortcutListOptions) ui.Element {
	k := pixel.K(c)
	for _, it := range items {
		if len(it.Keys) == 0 {
			panic(fmt.Sprintf("code: ShortcutList action %q has no keys", it.Action))
		}
	}
	return ui.Column(c).Width(k.D(pixel.Or("ShortcutList Width", o.Width, 220))).Gap(k.D(3)).Role(ui.RoleGroup).Label(label).Children(func() {
		for _, it := range items {
			ui.Row(c).FillWidth().Gap(k.D(4)).AlignItems(ui.Center).Role(ui.RoleText).Label(it.Action + ", " + strings.Join(it.Keys, "+")).Children(func() {
				pixel.Text(c, it.Action, pixel.TextOptions{}).Grow(1)
				// Long chords wrap rather than run off the edge.
				ui.Row(c).Wrap().Gap(k.D(2)).AlignItems(ui.Center).Children(func() {
					for i, key := range it.Keys {
						if i > 0 {
							pixel.Text(c, "+", pixel.TextOptions{Slot: pixel.Faint})
						}
						text.Kbd(c, key, text.KbdOptions{})
					}
				})
			})
		}
	})
}
