package navigation_test

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/navigation"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var navItems = []navigation.NavItem{{Icon: icon.Home, Label: "Home"}, {Icon: icon.Map, Label: "Map"}, {Icon: icon.Settings, Label: "Options"}}

func TestNavRailNamesTheHoveredItem(t *testing.T) {
	sel := 0
	tt := view(400, 200, newClock(), false, func(c *ui.Context) {
		navigation.NavRail(c, "Rail", navItems, &sel, navigation.NavRailOptions{})
	})
	m := find(t, tt, "Map")
	label := ui.Rect{X: m.X + m.W + 8, Y: m.Y, W: 60, H: m.H}
	before := slots(t, tt.Image(), label)
	tt.Move(m.X+m.W/2, m.Y+m.H/2)
	if slices.Equal(before, slots(t, tt.Image(), label)) {
		t.Fatal("hovering showed no name")
	}
	click(t, tt, "Map")
	tt.Key(0, ui.KeyDown)
	if sel != 2 {
		t.Fatalf("click and Down gave %d", sel)
	}
}

func TestBottomNavChooses(t *testing.T) {
	sel, changes := 0, 0
	tt := view(400, 200, newClock(), false, func(c *ui.Context) {
		if navigation.BottomNav(c, "Bottom", navItems, &sel, navigation.BottomNavOptions{}).Changed {
			changes++
		}
	})
	click(t, tt, "Options")
	tt.Key(0, ui.KeyLeft)
	if sel != 1 || changes != 2 {
		t.Fatalf("click and Left gave %d after %d changes", sel, changes)
	}
}

func TestLevelSelectMovesAndEnters(t *testing.T) {
	sel, chosen := 0, 0
	stages := []navigation.Stage{{Name: "1-1", Cleared: true, Stars: 3}, {Name: "1-2", Cleared: true, Stars: 1}, {Name: "1-3"}, {Name: "1-4", Locked: true}}
	tt := view(400, 200, newClock(), false, func(c *ui.Context) {
		if navigation.LevelSelect(c, "Stages", stages, &sel, navigation.LevelSelectOptions{Columns: 2}).Chosen {
			chosen++
		}
	})
	click(t, tt, "Stages")
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if sel != 3 || chosen != 0 {
		t.Fatalf("moved to %d; a locked stage was chosen %d times", sel, chosen)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyEnter)
	if sel != 2 || chosen != 1 {
		t.Fatalf("Enter on an open stage: %d, chosen %d", sel, chosen)
	}
	if count(slots(t, tt.Image(), find(t, tt, "Stages")), pixel.Warn) == 0 {
		t.Fatal("cleared stages show no stars")
	}
}

func TestWorldMapWalksAndEnters(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		clk := newClock()
		at, entered := 0, 0
		nodes := []navigation.MapNode{{Name: "Home", X: 30, Y: 70}, {Name: "Woods", X: 110, Y: 70}, {Name: "Cave", X: 110, Y: 20, Locked: true}}
		tt := view(600, 400, clk, reduce, func(c *ui.Context) {
			if navigation.WorldMap(c, "World", nodes, [][2]int{{0, 1}, {1, 2}}, &at, navigation.WorldMapOptions{}).Entered {
				entered++
			}
		})
		click(t, tt, "World")
		tt.Key(0, ui.KeyRight)
		mid := slots(t, tt.Image(), find(t, tt, "World"))
		clk.t = clk.t.Add(time.Second)
		tt.Frame()
		walked := !slices.Equal(mid, slots(t, tt.Image(), find(t, tt, "World")))
		tt.Key(0, ui.KeyUp)
		tt.Key(0, ui.KeyEnter)
		if at != 1 || entered != 1 || walked == reduce {
			t.Fatalf("reduce %v: at %d, entered %d, walked %v", reduce, at, entered, walked)
		}
	}
}

func TestDocumentTabsClose(t *testing.T) {
	sel, closed := 0, -1
	tt := view(500, 100, newClock(), false, func(c *ui.Context) {
		if r := navigation.DocumentTabs(c, "Files", []string{"map.txt", "save.dat"}, &sel, navigation.DocumentTabsOptions{}); r.Closed >= 0 {
			closed = r.Closed
		}
	})
	click(t, tt, "CLOSE save.dat")
	if closed != 1 || sel != 0 {
		t.Fatalf("× closed %d and chose %d", closed, sel)
	}
	click(t, tt, "map.txt")
	tt.Key(0, ui.KeyDelete)
	if closed != 0 {
		t.Fatalf("Delete closed %d", closed)
	}
}
