package game

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// The paper doll's slots, in the order Equipment takes gear.
const (
	Head = iota
	Body
	Weapon
	Shield
	Ring
)

// doll places the slots on a 3×3 grid.
var doll = []struct {
	at   Spot
	name string
	hint *pixel.Sprite
}{
	{Spot{1, 0}, "head", icon.Crown}, {Spot{1, 1}, "body", icon.User}, {Spot{0, 1}, "weapon", icon.Sword},
	{Spot{2, 1}, "shield", icon.Shield}, {Spot{0, 2}, "ring", icon.Gem},
}

// pitch is a doll cell: a slot and its gap.
const pitch = data.SlotSize + 4

type EquipmentOptions struct {
	Disabled bool
}

type EquipmentResult struct {
	Element *ui.Element
	Changed bool
	// Picked is the slot opened this frame, or -1.
	Picked int
}

// Equipment is a paper doll with five slots.
func Equipment(c *ui.Context, text string, gear []data.Item, selected *int, o EquipmentOptions) EquipmentResult {
	k := pixel.K(c)
	if len(gear) != len(doll) || *selected < 0 || *selected >= len(doll) {
		panic(fmt.Sprintf("game: Equipment has %d slots, slot %d; it takes head, body, weapon, shield, ring", len(gear), *selected))
	}
	old := *selected
	spots := make([]Spot, len(doll))
	captions := make([]string, len(doll))
	w := 3*pitch - 4
	for i, d := range doll {
		spots[i], captions[i] = d.at, k.Msg(d.name)
		if gear[i].Sprite != nil {
			captions[i] += ": " + gear[i].String()
		}
		w = max(w, pixel.Measure(captions[i]))
	}
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Disabled(o.Disabled)
	})
	arrows(base, spots, selected)
	res := EquipmentResult{Element: shell, Picked: -1}
	if pixel.Shortcut(base, 0, ui.KeyEnter) || pixel.Shortcut(base, 0, ui.KeySpace) {
		res.Picked = *selected
	}
	left := (w - (3*pitch - 4)) / 2
	// Enter and Space are taken, so clicks come from pointers.
	if base.Clicked() {
		x, y, _ := base.PointerPosition()
		dx, dy := int(x/k.Dot())-left, int(y/k.Dot())
		for i, s := range spots {
			if dx >= s.X*pitch && dx < s.X*pitch+data.SlotSize && dy >= s.Y*pitch && dy < s.Y*pitch+data.SlotSize {
				*selected, res.Picked = i, i
			}
		}
	}
	res.Changed = *selected != old
	base.Description(captions[*selected])
	if res.Changed {
		slog.Debug("8bit: equipment", "label", text, "slot", doll[*selected].name)
	}
	sel, shown := *selected, append([]data.Item(nil), gear...)
	shell.Size(k.D(w), k.D(3*pitch+pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		figure(g, left)
		for i, d := range doll {
			x, y := left+d.at.X*pitch, d.at.Y*pitch
			data.DrawSlot(g, x, y, shown[i], data.ItemSlotOptions{Selected: i == sel})
			if shown[i].Sprite == nil {
				g.Silhouette(x+(data.SlotSize-d.hint.W)/2, y+(data.SlotSize-d.hint.H)/2, d.hint, 0, pixel.Shade, 1)
			}
		}
		g.Text((w-pixel.Measure(captions[sel]))/2, 3*pitch, captions[sel], pixel.Ink, 1)
		g.Finish(st)
	})
	return res
}

// figure sketches the body the slots dress: neck, arms, legs.
func figure(g pixel.Canvas, x int) {
	mid := x + pitch + data.SlotSize/2
	g.Rect(mid-2, data.SlotSize, 4, 4, pixel.Shade)
	g.Rect(x+data.SlotSize/2, pitch+9, 2*pitch, 4, pixel.Shade)
	g.Rect(mid-6, 2*pitch-4, 5, pitch, pixel.Shade)
	g.Rect(mid+1, 2*pitch-4, 5, pitch, pixel.Shade)
}
