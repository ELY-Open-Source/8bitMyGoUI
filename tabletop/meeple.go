package tabletop

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MeepleOptions struct {
	// Player picks its color, 0 to 5.
	Player int
	// Empty leaves the slot open.
	Empty bool
}

// Meeple is a player's token, standing in its slot.
func Meeple(c *ui.Context, text string, o MeepleOptions) ui.Element {
	k := pixel.K(c)
	if o.Player < 0 || o.Player >= len(players) {
		panic(fmt.Sprintf("tabletop: Meeple player %d", o.Player))
	}
	e := ui.Box(c).Size(k.D(data.SlotSize), k.D(data.SlotSize)).Role(ui.RoleImage).Label(text)
	if o.Empty {
		e.Description(k.Msg("empty"))
	}
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		// Slots look alike everywhere: an empty item slot holds it.
		data.DrawSlot(g, 0, 0, data.Item{}, data.ItemSlotOptions{})
		if o.Empty {
			return
		}
		at := (data.SlotSize - 12) / 2
		for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			g.Silhouette(at+d[0], at+d[1], icon.Meeple, 0, pixel.Edge, 1)
		}
		g.Sprite(at, at, icon.Meeple, 0, players[o.Player], 1)
	})
}
