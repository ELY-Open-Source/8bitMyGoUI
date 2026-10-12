package navigation_test

import (
	"image"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/navigation"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func newClock() *clock { return &clock{time.Unix(3_000_000, 0)} }

// view builds on a fake clock; reduce stops motion.
func view(w, h int, clk *clock, reduce bool, build func(c *ui.Context)) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: clk.now})
		build(c)
	}, w, h)
	if reduce {
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	}
	return tt
}

func find(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("no %q in %q", label, tt.Texts())
	}
	return r
}

func click(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	if err := tt.Click(label); err != nil {
		t.Fatal(err)
	}
}

func shows(tt *ui.Tester, label string) bool {
	_, ok := tt.Find(label)
	return ok
}

// slots are the palette slots of a box's dots.
func slots(t *testing.T, img *image.RGBA, r ui.Rect) []pixel.Slot {
	t.Helper()
	var out []pixel.Slot
	for y := int(r.Y); y < int(r.Y+r.H); y += 2 {
		for x := int(r.X); x < int(r.X+r.W); x += 2 {
			c, found := img.RGBAAt(x, y), false
			for s, pc := range pixel.SpringDay.Colors {
				if pc.R == c.R && pc.G == c.G && pc.B == c.B {
					out, found = append(out, pixel.Slot(s)), true
					break
				}
			}
			if !found {
				t.Fatalf("pixel %d,%d off the palette", x, y)
			}
		}
	}
	return out
}

func count(s []pixel.Slot, want pixel.Slot) int {
	n := 0
	for _, v := range s {
		if v == want {
			n++
		}
	}
	return n
}

func TestTabsChooseAndShowTheirPanel(t *testing.T) {
	for _, side := range []bool{false, true} {
		tab, changes := 0, 0
		tt := view(500, 300, newClock(), false, func(c *ui.Context) {
			build := navigation.Tabs
			if side {
				build = navigation.VerticalTabs
			}
			names := []string{"Items", "Magic", "Status"}
			if build(c, "Menu", names, &tab, navigation.TabsOptions{}, func() {
				pixel.Text(c, "panel "+names[tab], pixel.TextOptions{})
			}).Changed {
				changes++
			}
		})
		click(t, tt, "Magic")
		find(t, tt, "panel Magic")
		// The panel's edge breaks open under the chosen tab only.
		edge := func(l string) pixel.Slot {
			r := find(t, tt, l)
			if side {
				return slots(t, tt.Image(), ui.Rect{X: r.X + r.W, Y: r.Y + r.H/2, W: 2, H: 2})[0]
			}
			return slots(t, tt.Image(), ui.Rect{X: r.X + r.W/2, Y: r.Y + r.H, W: 2, H: 2})[0]
		}
		if edge("Magic") == pixel.Edge || edge("Items") != pixel.Edge {
			t.Fatalf("side %v: the chosen tab does not join the panel", side)
		}
		tt.Key(0, ui.KeyRight)
		if tab != 2 || changes != 2 || !shows(tt, "panel Status") {
			t.Fatalf("side %v: click and Right gave tab %d after %d changes", side, tab, changes)
		}
	}
}

func TestSidebarChoosesAcrossSections(t *testing.T) {
	sel := 0
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		navigation.Sidebar(c, "Places", []navigation.SidebarSection{
			{Title: "WORLD", Items: []navigation.SidebarItem{{Icon: icon.Home, Label: "Town"}, {Icon: icon.Map, Label: "Field"}}},
			{Title: "MAIL", Items: []navigation.SidebarItem{{Icon: icon.Mail, Label: "Inbox", Badge: "3"}}},
		}, &sel, navigation.SidebarOptions{})
	})
	click(t, tt, "Field")
	tt.Key(0, ui.KeyDown)
	if sel != 2 {
		t.Fatalf("click and Down chose %d, want the next section's item", sel)
	}
	if count(slots(t, tt.Image(), find(t, tt, "Inbox")), pixel.Bad) == 0 {
		t.Fatal("the badge does not show")
	}
}

func TestBreadcrumbsLinkAllButTheLast(t *testing.T) {
	clicked := -2
	tt := view(400, 100, newClock(), false, func(c *ui.Context) {
		if r := navigation.Breadcrumbs(c, "Path", []string{"World", "Cave", "B2"}, navigation.BreadcrumbsOptions{}); r.Clicked >= 0 {
			clicked = r.Clicked
		}
	})
	click(t, tt, "B2")
	n := 0
	for _, s := range tt.Texts() {
		if s == "B2" {
			n++
		}
	}
	if clicked != -2 || n != 1 {
		t.Fatalf("the last part linked, or shows %d times", n)
	}
	click(t, tt, "Cave")
	if clicked != 1 {
		t.Fatalf("Cave reported %d", clicked)
	}
}

func TestPaginationAndDots(t *testing.T) {
	page, dot := 5, 1
	tt := view(600, 200, newClock(), false, func(c *ui.Context) {
		navigation.Pagination(c, "Pages", &page, 9, navigation.PaginationOptions{})
		navigation.PageDots(c, "Slides", &dot, 4, navigation.PageDotsOptions{})
	})
	for _, l := range []string{"1", "4", "5", "6", "9"} {
		find(t, tt, l)
	}
	if shows(tt, "3") || shows(tt, "7") {
		t.Fatal("pages far from the current one show")
	}
	click(t, tt, "6")
	tt.Key(0, ui.KeySpace)
	if page != 6 || !tt.Focused("6") {
		t.Fatalf("the focus left the clicked page: page %d", page)
	}
	click(t, tt, "NEXT")
	if page != 7 {
		t.Fatalf("6 then NEXT gave %d", page)
	}
	page = 1
	tt.Frame()
	click(t, tt, "PREVIOUS")
	if page != 1 {
		t.Fatal("PREVIOUS moved past page 1")
	}
	// The click lands on the third of four dots.
	click(t, tt, "Slides")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if dot != 4 {
		t.Fatalf("a click and Right gave dot %d", dot)
	}
	r := find(t, tt, "Slides")
	tt.ClickAt(r.X+4, r.Y+r.H/2)
	if dot != 1 {
		t.Fatalf("a click on the first dot gave %d", dot)
	}
}

func TestStepsWaveUnlessReduced(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		clk := newClock()
		tt := view(400, 100, clk, reduce, func(c *ui.Context) {
			navigation.Steps(c, []string{"Gate", "Cave", "Peak"}, 1, navigation.StepsOptions{})
		})
		r := find(t, tt, "Gate, Cave, Peak")
		before := slots(t, tt.Image(), r)
		clk.t = clk.t.Add(300 * time.Millisecond)
		tt.Frame()
		if waved := !slices.Equal(before, slots(t, tt.Image(), r)); waved == reduce {
			t.Errorf("reduce %v: the flag waved %v", reduce, waved)
		}
	}
}

func TestToolbarIsOneTabStop(t *testing.T) {
	tt := view(400, 100, newClock(), false, func(c *ui.Context) {
		input.Button(c, "before", input.ButtonOptions{})
		navigation.Toolbar(c, "Tools", navigation.ToolbarOptions{}, func() {
			input.IconButton(c, icon.Edit, "Pen", input.IconButtonOptions{})
			input.IconButton(c, icon.Trash, "Erase", input.IconButtonOptions{})
		})
		input.Button(c, "after", input.ButtonOptions{})
	})
	click(t, tt, "before")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Erase") {
		t.Fatal("Right did not move along the toolbar")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("after") {
		t.Fatal("Tab did not leave the toolbar")
	}
}

func TestBackForwardWalksHistory(t *testing.T) {
	r := ui.NewRouter("/a")
	r.Push("/b")
	tt := view(300, 100, newClock(), false, func(c *ui.Context) {
		navigation.BackForward(c, r, navigation.BackForwardOptions{})
	})
	click(t, tt, "FORWARD")
	click(t, tt, "BACK")
	if r.Path() != "/a" {
		t.Fatalf("BACK went to %q", r.Path())
	}
	click(t, tt, "BACK")
	click(t, tt, "FORWARD")
	if r.Path() != "/b" {
		t.Fatalf("FORWARD went to %q", r.Path())
	}
}

func TestTableOfContentsMarksTheCurrent(t *testing.T) {
	cur := 0
	tt := view(400, 200, newClock(), false, func(c *ui.Context) {
		navigation.TableOfContents(c, "Manual", []navigation.TocItem{{Title: "Controls"}, {Title: "Jumping", Level: 1}, {Title: "Bosses"}}, &cur, navigation.TableOfContentsOptions{})
	})
	click(t, tt, "Jumping")
	tt.Key(0, ui.KeyDown)
	if cur != 2 {
		t.Fatalf("click and Down gave %d", cur)
	}
	// The marker sits left of the title's indent.
	b := find(t, tt, "Bosses")
	if count(slots(t, tt.Image(), ui.Rect{X: b.X, Y: b.Y, W: 12, H: b.H}), pixel.Accent) == 0 {
		t.Fatal("no ▶ at the current section")
	}
}
