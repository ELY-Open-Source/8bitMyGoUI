package data_test

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestInventoryChoosesAndSwaps(t *testing.T) {
	items := []data.Item{{Sprite: icon.Potion, Name: "Potion", Count: 3}, {Sprite: icon.Sword, Name: "Sword", Rarity: display.Epic}, {}, {Sprite: icon.Gem, Name: "Gem"}}
	sel, swaps := 0, 0
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		if data.Inventory(c, "Bag", items, &sel, data.InventoryOptions{Columns: 2}).Swapped {
			swaps++
		}
	})
	r := find(t, tt, "Bag")
	tt.ClickAt(r.X+4, r.Y+4)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyRight)
	if sel != 3 {
		t.Fatalf("Down, Right chose %d", sel)
	}
	// One slot is 22 dots, 44 pixels.
	tt.Press(r.X+10, r.Y+10)
	tt.Move(r.X+30, r.Y+30)
	tt.Move(r.X+54, r.Y+10)
	tt.Release(r.X+54, r.Y+10)
	if items[1].Name != "Potion" || items[0].Name != "Sword" || swaps != 1 {
		t.Fatalf("drag swapped to %v, %d swaps", items, swaps)
	}
	tt.ClickAt(r.X+10, r.Y+10)
	if sel != 0 || items[0].Name != "Sword" || swaps != 1 {
		t.Fatalf("a click swapped or chose %d", sel)
	}
	keeps(t, tt.Image())
}

func TestInventoryOffNeitherChoosesNorSwaps(t *testing.T) {
	for _, inert := range []bool{false, true} {
		items := []data.Item{{Sprite: icon.Potion, Name: "Potion"}, {Sprite: icon.Gem, Name: "Gem"}}
		sel := 0
		tt := view(400, 300, newClock(), inert, func(c *ui.Context) {
			data.Inventory(c, "Bag", items, &sel, data.InventoryOptions{Disabled: !inert})
		})
		r := find(t, tt, "Bag")
		tt.Press(r.X+10, r.Y+10)
		tt.Move(r.X+54, r.Y+10)
		tt.Release(r.X+54, r.Y+10)
		tt.Key(0, ui.KeyRight)
		if sel != 0 || items[0].Name != "Potion" {
			t.Fatalf("inert %v: sel %d, items %v", inert, sel, items)
		}
	}
}

func TestItemSlotShowsRarityCountAndCooldown(t *testing.T) {
	tt := view(200, 100, newClock(), false, func(c *ui.Context) {
		ui.Row(c).Children(func() {
			data.ItemSlot(c, data.Item{Sprite: icon.Potion, Name: "Potion", Count: 12, Rarity: display.Legendary}, data.ItemSlotOptions{})
			data.ItemSlot(c, data.Item{Sprite: icon.Potion, Name: "Ether"}, data.ItemSlotOptions{Cooldown: 1})
		})
	})
	r, img := find(t, tt, "Potion ×12"), tt.Image()
	if edge := slot(t, img, int(r.X)+2, int(r.Y)+20); edge != display.Legendary.Slot() {
		t.Fatalf("legendary edge is %v", edge)
	}
	cool := find(t, tt, "Ether")
	shade := 0
	for x := int(cool.X) + 4; x < int(cool.X)+40; x += 2 {
		if slot(t, img, x, int(cool.Y)+6) == pixel.Shade {
			shade++
		}
	}
	if shade == 0 {
		t.Fatal("a full cooldown left no veil")
	}
	keeps(t, img)
}

func TestSortableMovesByAltArrowsAndDrag(t *testing.T) {
	items := []string{"Aria", "Brom", "Cyra"}
	changes := 0
	tt := view(300, 200, newClock(), false, func(c *ui.Context) {
		if data.Sortable(c, "Order", &items, data.SortableOptions{}).Changed {
			changes++
		}
	})
	r := find(t, tt, "Order")
	tt.ClickAt(r.X+20, r.Y+4)
	tt.Key(ui.Alt, ui.KeyDown)
	tt.Key(ui.Alt, ui.KeyDown)
	if !slices.Equal(items, []string{"Brom", "Cyra", "Aria"}) {
		t.Fatalf("Alt+Down twice gave %v", items)
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(ui.Alt, ui.KeyUp)
	if !slices.Equal(items, []string{"Cyra", "Brom", "Aria"}) {
		t.Fatalf("Up, Alt+Up gave %v", items)
	}
	// Rows are 16 dots, 32 pixels.
	tt.Press(r.X+20, r.Y+72)
	tt.Move(r.X+20, r.Y+40)
	tt.Move(r.X+20, r.Y+8)
	tt.Release(r.X+20, r.Y+8)
	if !slices.Equal(items, []string{"Aria", "Cyra", "Brom"}) || changes != 5 {
		t.Fatalf("drag to the top gave %v, %d changes", items, changes)
	}
}

func TestTransferMovesChosenItems(t *testing.T) {
	left, right := []string{"Aria", "Brom"}, []string{"Cyra"}
	changes := 0
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		if data.Transfer(c, "Party", &left, &right, data.TransferOptions{}).Changed {
			changes++
		}
	})
	click(t, tt, "Brom")
	click(t, tt, "RIGHT")
	if !slices.Equal(left, []string{"Aria"}) || !slices.Equal(right, []string{"Cyra", "Brom"}) {
		t.Fatalf("→ moved to %v %v", left, right)
	}
	click(t, tt, "Cyra")
	click(t, tt, "LEFT")
	if !slices.Equal(left, []string{"Aria", "Cyra"}) || !slices.Equal(right, []string{"Brom"}) || changes != 2 {
		t.Fatalf("← moved to %v %v, %d changes", left, right, changes)
	}
}

func TestMatrixTogglesByClickAndKeys(t *testing.T) {
	on := [][]bool{{false, false}, {false, true}}
	changes := 0
	tt := view(300, 200, newClock(), false, func(c *ui.Context) {
		if data.Matrix(c, "Rights", []string{"Read", "Write"}, []string{"Ann", "Bob"}, on, data.MatrixOptions{}).Changed {
			changes++
		}
	})
	r := find(t, tt, "Rights")
	tt.ClickAt(r.X+4, r.Y+4)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if !on[0][0] || !on[1][0] || on[1][1] || changes != 3 {
		t.Fatalf("click, Down Space, Right Enter left %v, %d changes", on, changes)
	}
	keeps(t, tt.Image())
}

func TestInventoryDragEndsWithoutSwapWhenTurnedOff(t *testing.T) {
	for _, inert := range []bool{false, true} {
		items := []data.Item{{Sprite: icon.Potion, Name: "Potion"}, {Sprite: icon.Gem, Name: "Gem"}}
		sel, off, swaps := 0, false, 0
		tt := ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Inert: off && inert})
			if data.Inventory(c, "Bag", items, &sel, data.InventoryOptions{Disabled: off && !inert}).Swapped {
				swaps++
			}
		}, 400, 300)
		r := find(t, tt, "Bag")
		tt.Press(r.X+10, r.Y+10)
		tt.Move(r.X+54, r.Y+10)
		off = true
		tt.Frame()
		off = false
		tt.Release(r.X+54, r.Y+10)
		if items[0].Name != "Potion" || swaps != 0 {
			t.Fatalf("inert %v: turning off mid-drag swapped %v", inert, items)
		}
	}
}

func TestInventoryDragSurvivesShrinkingItems(t *testing.T) {
	items := []data.Item{{Sprite: icon.Potion, Name: "Potion"}, {Sprite: icon.Gem, Name: "Gem"}, {Sprite: icon.Key, Name: "Key"}}
	sel, swaps := 0, 0
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		if data.Inventory(c, "Bag", items, &sel, data.InventoryOptions{}).Swapped {
			swaps++
		}
	})
	r := find(t, tt, "Bag")
	tt.Press(r.X+98, r.Y+10)
	tt.Move(r.X+10, r.Y+10)
	items = items[:1]
	tt.Frame()
	tt.Release(r.X+10, r.Y+10)
	if items[0].Name != "Potion" || swaps != 0 {
		t.Fatalf("a drag from a lost slot swapped %v", items)
	}
}

func TestInventoryReleaseWhileTurningOffDoesNotSwap(t *testing.T) {
	for _, inert := range []bool{false, true} {
		items := []data.Item{{Sprite: icon.Potion, Name: "Potion"}, {Sprite: icon.Gem, Name: "Gem"}}
		sel, off := 0, false
		tt := ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Inert: off && inert})
			data.Inventory(c, "Bag", items, &sel, data.InventoryOptions{Disabled: off && !inert})
		}, 400, 300)
		r := find(t, tt, "Bag")
		tt.Press(r.X+10, r.Y+10)
		tt.Move(r.X+54, r.Y+10)
		off = true
		tt.Release(r.X+54, r.Y+10)
		if items[0].Name != "Potion" {
			t.Fatalf("inert %v: a release as it turned off swapped %v", inert, items)
		}
	}
}

func TestMatrixCursorSurvivesShrinkingData(t *testing.T) {
	rows, cols := []string{"Read", "Write"}, []string{"Ann", "Bob"}
	on := [][]bool{{false, false}, {false, false}}
	tt := view(300, 200, newClock(), false, func(c *ui.Context) {
		data.Matrix(c, "Rights", rows, cols, on, data.MatrixOptions{})
	})
	r := find(t, tt, "Rights")
	tt.ClickAt(r.X+4, r.Y+4)
	tt.Key(0, ui.KeyEnd)
	rows, cols, on = rows[:1], cols[:1], [][]bool{{false}}
	tt.Frame()
	tt.Key(0, ui.KeySpace)
	if !on[0][0] {
		t.Fatalf("Space after shrinking left %v", on)
	}
}
