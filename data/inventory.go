package data

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type InventoryOptions struct {
	// Columns is the grid's width; zero means 6.
	Columns  int
	Disabled bool
}

type InventoryResult struct {
	Element ui.Element
	// Swapped reports two slots traded by a drag.
	Changed, Swapped bool
}

type bagKey struct{}

// bag tracks a drag; a dead one never swaps.
type bag struct {
	from, over int
	held, dead bool
}

// Inventory is a slot grid; arrows choose, drags swap.
func Inventory(c *ui.Context, text string, items []Item, selected *int, o InventoryOptions) InventoryResult {
	k := pixel.K(c)
	cols := pixel.Or("Inventory Columns", o.Columns, 6)
	off := o.Disabled || pixel.Current(c).Inert
	if *selected < 0 || *selected >= len(items) {
		panic(fmt.Sprintf("data: Inventory slot %d of %d", *selected, len(items)))
	}
	old := *selected
	shell, base, st, _ := pixel.Pad(c, text, cols, len(items), SlotSize, SlotSize, selected, o.Disabled)
	b := ui.Local(base, bagKey{}, func() bag { return bag{} })
	res := InventoryResult{Element: shell, Changed: *selected != old}
	px, py, _ := base.PointerPosition()
	at := func() int {
		col, row := int(px/k.Dot())/SlotSize, int(py/k.Dot())/SlotSize
		if px < 0 || py < 0 || col >= cols || row*cols+col >= len(items) {
			return -1
		}
		return row*cols + col
	}
	n := len(items)
	switch _, _, pressed := base.Dragged(); {
	case pressed && (off || b.held && max(b.from, b.over) >= n):
		b.held, b.dead = false, true
	case pressed && b.dead:
	case pressed && !b.held:
		b.held, b.from, b.over = true, at(), at()
	case pressed:
		b.over = at()
	case b.held:
		if !off && min(b.from, b.over) >= 0 && max(b.from, b.over) < n && b.from != b.over {
			items[b.from], items[b.over] = items[b.over], items[b.from]
			res.Swapped = true
		}
		b.held = false
	default:
		b.dead = false
	}
	base.Description(items[*selected].String())
	drag := -1
	if b.held && b.over != b.from {
		drag = b.from
	}
	shown := append([]Item(nil), items...)
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, it := range shown {
			if i == drag {
				it = Item{}
			}
			DrawSlot(g, i%cols*SlotSize, i/cols*SlotSize, it, ItemSlotOptions{Selected: i == *selected})
		}
		// The dragged item rides with the pointer.
		if drag >= 0 && shown[drag].Sprite != nil {
			g.Sprite(int(px/k.Dot())-6, int(py/k.Dot())-6, shown[drag].Sprite, 0, pixel.Ink, 1)
		}
		g.Finish(st)
	})
	if res.Changed {
		slog.Debug("8bit: inventory chose", "label", text, "slot", *selected)
	}
	if res.Swapped {
		slog.Debug("8bit: inventory swap", "label", text, "from", b.from, "to", b.over)
	}
	return res
}
