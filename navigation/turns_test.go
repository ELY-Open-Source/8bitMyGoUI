package navigation_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/navigation"
)

func TestKeysInAVeilsFirstFrameTurnNothing(t *testing.T) {
	nav := []navigation.NavItem{{Icon: icon.Home, Label: "Home"}, {Icon: icon.Star, Label: "Stars"}}
	cases := map[string]struct {
		build func(c *ui.Context, v *int)
		key   ui.Key
	}{
		"NavRail":   {func(c *ui.Context, v *int) { navigation.NavRail(c, "Rail", nav, v, navigation.NavRailOptions{}) }, ui.KeyDown},
		"BottomNav": {func(c *ui.Context, v *int) { navigation.BottomNav(c, "Nav", nav, v, navigation.BottomNavOptions{}) }, ui.KeyRight},
		"TableOfContents": {func(c *ui.Context, v *int) {
			navigation.TableOfContents(c, "Contents", []navigation.TocItem{{Title: "Intro"}, {Title: "Rules"}}, v, navigation.TableOfContentsOptions{})
		}, ui.KeyDown},
		"Sidebar": {func(c *ui.Context, v *int) {
			navigation.Sidebar(c, "Side", []navigation.SidebarSection{{Title: "Game", Items: []navigation.SidebarItem{{Label: "Play"}, {Label: "Load"}}}}, v, navigation.SidebarOptions{})
		}, ui.KeyDown},
		"Tabs": {func(c *ui.Context, v *int) {
			navigation.Tabs(c, "Tabs", []string{"Items", "Magic"}, v, navigation.TabsOptions{}, func() {})
		}, ui.KeyRight},
		"DocumentTabs": {func(c *ui.Context, v *int) {
			navigation.DocumentTabs(c, "Docs", []string{"a.txt", "b.txt"}, v, navigation.DocumentTabsOptions{})
		}, ui.KeyRight},
	}
	for name, tc := range cases {
		for _, veil := range []bool{false, true} {
			v, inert := 0, false
			tt := veiled(&inert, func(c *ui.Context) { tc.build(c, &v) })
			tt.Key(0, ui.KeyTab)
			inert = veil
			tt.Key(0, tc.key)
			if turned := v != 0; turned == veil {
				t.Errorf("%s: veiled %v, the key made %s", name, veil, fmt.Sprint(v))
			}
		}
	}
}

func TestDeleteInAVeilsFirstFrameClosesNoTab(t *testing.T) {
	for _, veil := range []bool{false, true} {
		v, closed, inert := 0, -1, false
		tt := veiled(&inert, func(c *ui.Context) {
			if r := navigation.DocumentTabs(c, "Docs", []string{"a.txt", "b.txt"}, &v, navigation.DocumentTabsOptions{}); r.Closed >= 0 {
				closed = r.Closed
			}
		})
		tt.Key(0, ui.KeyTab)
		inert = veil
		tt.Key(0, ui.KeyDelete)
		if (closed == 0) == veil {
			t.Errorf("veiled %v, Delete closed %d", veil, closed)
		}
	}
}

func TestDocumentTabFocusFollowsItsTab(t *testing.T) {
	titles, sel, closed := []string{"a.txt", "b.txt", "c.txt"}, 1, ""
	tt := view(600, 300, newClock(), false, func(c *ui.Context) {
		if r := navigation.DocumentTabs(c, "Docs", titles, &sel, navigation.DocumentTabsOptions{}); r.Closed >= 0 {
			closed = titles[r.Closed]
		}
	})
	tt.Key(0, ui.KeyTab)
	// a.txt closes elsewhere; b.txt keeps its focus.
	titles, sel = titles[1:], 0
	tt.Frame()
	tt.Key(0, ui.KeyDelete)
	if closed != "b.txt" {
		t.Fatalf("Delete after a.txt left closed %q", closed)
	}
}

func TestTwinDocumentTabsPanic(t *testing.T) {
	defer func() {
		if r := fmt.Sprint(recover()); !strings.HasPrefix(r, "navigation: DocumentTabs") {
			t.Fatalf("two a.txt tabs gave %q", r)
		}
	}()
	sel := 0
	view(400, 200, newClock(), false, func(c *ui.Context) {
		navigation.DocumentTabs(c, "Docs", []string{"a.txt", "a.txt"}, &sel, navigation.DocumentTabsOptions{})
	})
}
