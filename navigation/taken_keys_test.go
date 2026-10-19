package navigation_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/navigation"
)

func TestPageDotsKeysDoNotClickWhereThePointerRests(t *testing.T) {
	page := 1
	tt := view(400, 200, newClock(), false, func(c *ui.Context) {
		navigation.PageDots(c, "Pages", &page, 9, navigation.PageDotsOptions{})
	})
	r := find(t, tt, "Pages")
	tt.Move(r.X+r.W-4, r.Y+r.H/2)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyEnter)
	if page != 1 {
		t.Fatalf("Space and Enter turned to page %d", page)
	}
}

func TestSpaceTakesTheCursorNotThePointer(t *testing.T) {
	chosen := -1
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		if r := navigation.CommandMenu(c, "Battle", []navigation.MenuItem{{Label: "FIGHT"}, {Label: "MAGIC"}, {Label: "RUN"}}, navigation.CommandMenuOptions{}); r.Chosen >= 0 {
			chosen = r.Chosen
		}
	})
	r := find(t, tt, "Battle")
	tt.Move(r.X+r.W/2, r.Y+r.H-6)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if chosen != 0 {
		t.Fatalf("Space in the menu chose %d", chosen)
	}
	at, entered := 0, 0
	nodes := []navigation.MapNode{{Name: "Town", X: 20, Y: 20}, {Name: "Cave", X: 120, Y: 20}}
	tt = view(600, 400, newClock(), false, func(c *ui.Context) {
		if navigation.WorldMap(c, "Map", nodes, [][2]int{{0, 1}}, &at, navigation.WorldMapOptions{}).Entered {
			entered++
		}
	})
	r = find(t, tt, "Map")
	// The pointer rests on the Cave; the hero is home.
	tt.Move(r.X+2*120, r.Y+2*20)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if at != 0 || entered != 1 {
		t.Fatalf("Space walked to %d, entered %d times", at, entered)
	}
}
