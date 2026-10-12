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
	o.check()
	return ui.Box(c).Role(ui.RoleImage).Label(it.String()).Size(k.D(SlotSize), k.D(SlotSize)).Draw(func(p *ui.Painter, r ui.Rect) {
		DrawSlot(k.Canvas(p, r), 0, 0, it, o)
	})
}

func (o ItemSlotOptions) check() {
	if !(o.Cooldown >= 0 && o.Cooldown <= 1) {
		panic(fmt.Sprintf("data: Cooldown %v outside 0–1", o.Cooldown))
	}
}

// SlotSize is a slot's side in dots.
const SlotSize = 22

// String reads an item with its count: "Potion ×3".
func (it Item) String() string {
	if it.Count > 1 {
		return it.Name + " ×" + strconv.Itoa(it.Count)
	}
	return it.Name
}

// DrawSlot paints a slot, top left at x, y.
func DrawSlot(g pixel.Canvas, x, y int, it Item, o ItemSlotOptions) {
	o.check()
	edge := pixel.Edge
	if it.Sprite != nil {
		edge = it.Rarity.Slot()
	}
	if o.Selected {
		g.Rect(x, y, SlotSize, SlotSize, pixel.Accent)
	}
	g.Rect(x+1, y+1, SlotSize-2, SlotSize-2, edge)
	g.Rect(x+2, y+2, SlotSize-4, SlotSize-4, pixel.Well)
	if sp := it.Sprite; sp != nil {
		g.Sprite(x+2+(SlotSize-4-sp.W)/2, y+2+(SlotSize-4-sp.H)/2, sp, 0, pixel.Ink, 1)
	}
	if h := int(o.Cooldown * (SlotSize - 4)); h > 0 {
		g.Veil(x+2, y+2, SlotSize-4, h, pixel.Shade, 10)
	}
	if it.Count > 1 {
		n := strconv.Itoa(it.Count)
		g.Micro(x+SlotSize-2-pixel.MicroWidth(n), y+SlotSize-8, n, pixel.Ink, 1)
	}
}
