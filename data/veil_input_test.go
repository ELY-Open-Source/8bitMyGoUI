package data_test

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// underVeil builds fn where the test raises a veil.
func underVeil(inert *bool, fn func(c *ui.Context)) *ui.Tester {
	return view(400, 300, newClock(), false, func(c *ui.Context) {
		s := pixel.Current(c)
		s.Inert = *inert
		pixel.Scope(c, s, func() { fn(c) })
	})
}

var tree = []data.TreeNode{{Label: "Maps", Children: []data.TreeNode{{Label: "Cave"}}}, {Label: "Saves"}}

func TestVeilsStopKeysInTheirFirstFrame(t *testing.T) {
	for name, tc := range map[string]struct {
		build func(c *ui.Context, v *string)
		mods  ui.Modifiers
		key   ui.Key
	}{
		"Sortable": {func(c *ui.Context, v *string) {
			items := []string{"A", "B", "C"}
			if r := data.Sortable(c, "List", &items, data.SortableOptions{}); r.Changed {
				*v = fmt.Sprint(items)
			}
		}, ui.Alt, ui.KeyDown},
		"Tree": {func(c *ui.Context, v *string) { data.Tree(c, "Files", tree, v, data.TreeOptions{}) }, 0, ui.KeyDown},
		"Inventory": {func(c *ui.Context, v *string) {
			at := 0
			fmt.Sscan(*v, &at)
			bag := []data.Item{{Sprite: icon.Sword, Name: "Sword"}, {Sprite: icon.Gem, Name: "Gem"}}
			data.Inventory(c, "Bag", bag, &at, data.InventoryOptions{})
			*v = fmt.Sprint(at)
		}, 0, ui.KeyRight},
	} {
		v, inert := "0", false
		tt := underVeil(&inert, func(c *ui.Context) { tc.build(c, &v) })
		tt.Key(0, ui.KeyTab)
		inert = true
		tt.Key(tc.mods, tc.key)
		if v != "0" {
			t.Errorf("%s: a key in the veil's first frame made %s", name, v)
		}
	}
}

func TestDragsStopUnderAVeil(t *testing.T) {
	items, inert := []string{"A", "B", "C"}, false
	tt := underVeil(&inert, func(c *ui.Context) { data.Sortable(c, "List", &items, data.SortableOptions{}) })
	r := find(t, tt, "List")
	row := r.H / 3
	tt.Press(r.X+10, r.Y+row/2)
	tt.Move(r.X+10, r.Y+row*3/2)
	moved := fmt.Sprint(items)
	inert = true
	tt.Move(r.X+10, r.Y+row*5/2)
	tt.Release(r.X+10, r.Y+row*5/2)
	if moved != "[B A C]" || fmt.Sprint(items) != moved {
		t.Errorf("Sortable: %s before the veil, %v under it", moved, items)
	}
	bag, at := []data.Item{{Sprite: icon.Sword, Name: "Sword"}, {Sprite: icon.Gem, Name: "Gem"}}, 0
	inert = false
	tt = underVeil(&inert, func(c *ui.Context) { data.Inventory(c, "Bag", bag, &at, data.InventoryOptions{}) })
	r = find(t, tt, "Bag")
	tt.Press(r.X+20, r.Y+20)
	tt.Move(r.X+64, r.Y+20)
	inert = true
	tt.Release(r.X+64, r.Y+20)
	if bag[0].Name != "Sword" {
		t.Error("Inventory swapped under the veil")
	}
}

func TestTreeLeavesAnUnknownChoiceAlone(t *testing.T) {
	chosen, changed := "9/9", 0
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		if data.Tree(c, "Files", tree, &chosen, data.TreeOptions{}).Changed {
			changed++
		}
	})
	if chosen != "9/9" || changed != 0 {
		t.Fatalf("an unknown choice became %q, %d changes", chosen, changed)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if chosen != "0" || changed != 1 {
		t.Fatalf("Down chose %q, %d changes", chosen, changed)
	}
}
