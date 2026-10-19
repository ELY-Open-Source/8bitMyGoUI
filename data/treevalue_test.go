package data_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestTreeDrawsValuesInTheirSlot(t *testing.T) {
	sel := "0"
	roots := []data.TreeNode{{Label: "name", Value: `"Aria"`, Slot: pixel.Good}, {Label: "hp", Value: "42", Slot: pixel.Warn}}
	tt := view(300, 200, newClock(), false, func(c *ui.Context) {
		data.Tree(c, "Doc", roots, &sel, data.TreeOptions{})
	})
	r, img := find(t, tt, "Doc"), tt.Image()
	count := func(row int, s pixel.Slot) int {
		n, want := 0, pixel.SpringDay.Colors[s]
		for y := int(r.Y) + row*2*pixel.LineHeight; y < int(r.Y)+(row+1)*2*pixel.LineHeight; y++ {
			for x := int(r.X); x < int(r.X+r.W); x++ {
				if c := img.RGBAAt(x, y); c.R == want.R && c.G == want.G && c.B == want.B {
					n++
				}
			}
		}
		return n
	}
	// The chosen row inks its value; the other keeps Warn.
	if count(0, pixel.Good) != 0 || count(1, pixel.Warn) == 0 {
		t.Fatalf("values: Good on the chosen row %d, Warn below %d", count(0, pixel.Good), count(1, pixel.Warn))
	}
	// A click on the middle chooses the second row.
	click(t, tt, "Doc")
	img = tt.Image()
	if sel != "1" || count(0, pixel.Good) == 0 || count(1, pixel.Warn) != 0 {
		t.Fatalf("choosing hp (%q) did not swap which value is inked", sel)
	}
}

func TestTreeRevealsAChoiceInAClosedBranch(t *testing.T) {
	// Ether sits in Bag, which starts closed.
	sel := "0/1"
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		data.Tree(c, "Items", bag, &sel, data.TreeOptions{})
	})
	if r := find(t, tt, "Items"); sel != "0/1" || r.H != 4*2*pixel.LineHeight {
		t.Fatalf("the choice became %q in %v, want Bag opened on Ether", sel, r)
	}
}
