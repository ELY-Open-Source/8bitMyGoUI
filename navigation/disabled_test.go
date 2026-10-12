package navigation_test

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/navigation"
)

// Each control, disabled after taking focus, ignores every input.
func TestDisabledNavigationTakesNoInput(t *testing.T) {
	type state struct{ n, events int }
	items := []navigation.MenuItem{{Label: "A"}, {Label: "B"}}
	cases := []struct {
		name, label string
		build       func(c *ui.Context, v *state, off bool)
	}{
		{"Tabs", "Two", func(c *ui.Context, v *state, off bool) {
			navigation.Tabs(c, "T", []string{"One", "Two", "Six"}, &v.n, navigation.TabsOptions{Disabled: off}, func() {})
		}},
		{"Sidebar", "Two", func(c *ui.Context, v *state, off bool) {
			navigation.Sidebar(c, "S", []navigation.SidebarSection{{Items: []navigation.SidebarItem{{Label: "One"}, {Label: "Two"}, {Label: "Six"}}}}, &v.n, navigation.SidebarOptions{Disabled: off})
		}},
		{"PageDots", "Dots", func(c *ui.Context, v *state, off bool) {
			v.n = max(v.n, 1)
			navigation.PageDots(c, "Dots", &v.n, 5, navigation.PageDotsOptions{Disabled: off})
		}},
		{"CommandMenu", "Cmd", func(c *ui.Context, v *state, off bool) {
			if r := navigation.CommandMenu(c, "Cmd", items, navigation.CommandMenuOptions{Disabled: off}); r.Chosen >= 0 || r.Canceled {
				v.events++
			}
		}},
		{"DropdownMenu", "Drop", func(c *ui.Context, v *state, off bool) {
			if navigation.DropdownMenu(c, "Drop", items, navigation.DropdownMenuOptions{Disabled: off}).Chosen >= 0 {
				v.events++
			}
		}},
		{"LevelSelect", "Stages", func(c *ui.Context, v *state, off bool) {
			if navigation.LevelSelect(c, "Stages", []navigation.Stage{{Name: "1"}, {Name: "2"}, {Name: "3"}}, &v.n, navigation.LevelSelectOptions{Disabled: off}).Chosen {
				v.events++
			}
		}},
		{"WorldMap", "World", func(c *ui.Context, v *state, off bool) {
			nodes := []navigation.MapNode{{Name: "A", X: 20, Y: 20}, {Name: "B", X: 80, Y: 20}}
			if navigation.WorldMap(c, "World", nodes, [][2]int{{0, 1}}, &v.n, navigation.WorldMapOptions{Disabled: off}).Entered {
				v.events++
			}
		}},
		{"DocumentTabs", "Two", func(c *ui.Context, v *state, off bool) {
			if navigation.DocumentTabs(c, "D", []string{"One", "Two"}, &v.n, navigation.DocumentTabsOptions{Disabled: off}).Closed >= 0 {
				v.events++
			}
		}},
		{"NavRail", "Map", func(c *ui.Context, v *state, off bool) {
			navigation.NavRail(c, "R", []navigation.NavItem{{Icon: icon.Home, Label: "Home"}, {Icon: icon.Map, Label: "Map"}, {Icon: icon.Gem, Label: "Gem"}}, &v.n, navigation.NavRailOptions{Disabled: off})
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var v state
			off := false
			tt := view(500, 400, newClock(), false, func(c *ui.Context) {
				tc.build(c, &v, off)
				ui.Box(c).Size(4, 4).Label("elsewhere")
			})
			click(t, tt, tc.label)
			off = true
			tt.Frame()
			if tt.Focused(tc.label) && !tt.Focused("elsewhere") {
				t.Fatal("a disabled control kept the focus")
			}
			before := fmt.Sprint(v)
			for _, k := range []ui.Key{ui.KeyDown, ui.KeyRight, ui.KeyUp, ui.KeyLeft, ui.KeyEnd, ui.KeyEnter, ui.KeySpace, ui.KeyEscape, ui.KeyDelete} {
				tt.Key(0, k)
			}
			click(t, tt, tc.label)
			if now := fmt.Sprint(v); now != before {
				t.Fatalf("disabled control changed: %s to %s", before, now)
			}
		})
	}
}
