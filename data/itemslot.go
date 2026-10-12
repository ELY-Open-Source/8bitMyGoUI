package data

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Item is what a slot holds; no Sprite means empty.
type Item struct {
	Sprite *pixel.Sprite
	Name   string
	// Count shows when above one.
	Count  int
	Rarity display.Rarity
}

type ItemSlotOptions struct {
	// Cooldown, 0 to 1, veils the slot downward.
	Cooldown float64
	Selected bool
}

// ItemSlot shows an item, count, rarity edge and cooldown.
func ItemSlot(c *ui.Context, it Item, o ItemSlotOptions) *ui.Element {
	k := pixel.K(c)
	if !(o.Cooldown >= 0 && o.Cooldown <= 1) {
		panic(fmt.Sprintf("data: Cooldown %v outside 0–1", o.Cooldown))
	}
	return ui.Box(c).Role(ui.RoleImage).Label(itemName(it)).Size(k.D(slotSize), k.D(slotSize)).Draw(func(p *ui.Painter, r ui.Rect) {
		slot(k.Canvas(p, r), 0, 0, it, o.Cooldown, o.Selected)
	})
}

const slotSize = 22

// itemName reads an item with its count.
func itemName(it Item) string {
	if it.Count > 1 {
		return it.Name + " ×" + strconv.Itoa(it.Count)
	}
	return it.Name
}

// slot draws a slot with its top-left at x, y.
func slot(g pixel.Canvas, x, y int, it Item, cool float64, chosen bool) {
	edge := pixel.Edge
	if it.Sprite != nil {
		edge = it.Rarity.Slot()
	}
	if chosen {
		g.Rect(x, y, slotSize, slotSize, pixel.Accent)
	}
	g.Rect(x+1, y+1, slotSize-2, slotSize-2, edge)
	g.Rect(x+2, y+2, slotSize-4, slotSize-4, pixel.Well)
	if sp := it.Sprite; sp != nil {
		g.Sprite(x+2+(slotSize-4-sp.W)/2, y+2+(slotSize-4-sp.H)/2, sp, 0, pixel.Ink, 1)
	}
	if h := int(cool * (slotSize - 4)); h > 0 {
		g.Veil(x+2, y+2, slotSize-4, h, pixel.Shade, 10)
	}
	if it.Count > 1 {
		n := strconv.Itoa(it.Count)
		g.Micro(x+slotSize-2-pixel.MicroWidth(n), y+slotSize-8, n, pixel.Ink, 1)
	}
}
