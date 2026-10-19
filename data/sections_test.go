package data_test

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestCollapsibleShowsContentWhileOpen(t *testing.T) {
	open, changes := false, 0
	tt := view(300, 200, newClock(), false, func(c *ui.Context) {
		if data.Collapsible(c, "Secrets", &open, data.CollapsibleOptions{}, func() {
			pixel.Text(c, "Konami", pixel.TextOptions{})
		}).Changed {
			changes++
		}
	})
	if shows(tt, "Konami") {
		t.Fatal("closed content shows")
	}
	click(t, tt, "Secrets")
	if !open || !shows(tt, "Konami") {
		t.Fatal("a click did not open it")
	}
	tt.Key(0, ui.KeyEnter)
	if open || shows(tt, "Konami") || changes != 2 {
		t.Fatalf("Enter left open %v, %d changes", open, changes)
	}
	keeps(t, tt.Image())
}

func TestAccordionKeepsOneOpenUnlessMany(t *testing.T) {
	for _, many := range []bool{false, true} {
		open := []bool{true, false, false}
		tt := view(300, 300, newClock(), false, func(c *ui.Context) {
			data.Accordion(c, "Help", []string{"Move", "Jump", "Save"}, open, data.AccordionOptions{Many: many}, func(i int) {
				pixel.Text(c, []string{"Arrows", "Space", "Start"}[i], pixel.TextOptions{})
			})
		})
		click(t, tt, "Save")
		want := []bool{many, false, true}
		if !slices.Equal(open, want) || shows(tt, "Arrows") != many || !shows(tt, "Start") {
			t.Fatalf("many %v: open %v", many, open)
		}
	}
}

func TestCollapsibleOffStaysPut(t *testing.T) {
	for _, inert := range []bool{false, true} {
		open := false
		tt := view(300, 200, newClock(), inert, func(c *ui.Context) {
			data.Collapsible(c, "Secrets", &open, data.CollapsibleOptions{Disabled: !inert}, func() {})
		})
		click(t, tt, "Secrets")
		if open {
			t.Fatalf("inert %v: a click opened it", inert)
		}
	}
}

func deck(c *ui.Context, page *int, o data.CarouselOptions, changes *int) {
	if data.Carousel(c, "Tips", 3, page, o, func(i int) {
		pixel.Text(c, []string{"Tip A", "Tip B", "Tip C"}[i], pixel.TextOptions{})
	}).Changed {
		*changes++
	}
}

func TestCarouselTurnsAndWraps(t *testing.T) {
	page, changes := 0, 0
	tt := view(300, 200, newClock(), false, func(c *ui.Context) { deck(c, &page, data.CarouselOptions{}, &changes) })
	click(t, tt, "PREVIOUS")
	if page != 2 || !shows(tt, "Tip C") {
		t.Fatalf("◀ from the first page went to %d", page)
	}
	click(t, tt, "NEXT")
	click(t, tt, "NEXT")
	if page != 1 || !shows(tt, "Tip B") || changes != 3 {
		t.Fatalf("▶ ▶ went to %d, %d changes", page, changes)
	}
	keeps(t, tt.Image())
}

func TestCarouselTurnsByItselfUnlessReduced(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		clk, page, changes := newClock(), 0, 0
		tt := view(300, 200, clk, false, func(c *ui.Context) { deck(c, &page, data.CarouselOptions{Every: 3}, &changes) })
		if reduce {
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		}
		tt.Move(599, 399)
		clk.t = clk.t.Add(2 * time.Second)
		tt.Frame()
		if page != 0 {
			t.Fatalf("reduce %v: turned after 2 s", reduce)
		}
		clk.t = clk.t.Add(time.Second)
		tt.Frame()
		if want := map[bool]int{false: 1, true: 0}[reduce]; page != want || !shows(tt, []string{"Tip A", "Tip B"}[want]) {
			t.Fatalf("reduce %v: page %d after 3 s", reduce, page)
		}
	}
}

func TestCarouselWaitsUnderThePointer(t *testing.T) {
	clk, page, changes := newClock(), 0, 0
	tt := view(300, 200, clk, false, func(c *ui.Context) { deck(c, &page, data.CarouselOptions{Every: 1}, &changes) })
	r := find(t, tt, "Tip A")
	tt.Move(r.X+2, r.Y+2)
	clk.t = clk.t.Add(5 * time.Second)
	tt.Frame()
	if page != 0 {
		t.Fatalf("turned to %d under the pointer", page)
	}
}

func TestPagerStepsAndResizes(t *testing.T) {
	page, size, changes := 0, 20, 0
	total := 312
	tt := view(500, 100, newClock(), false, func(c *ui.Context) {
		if data.Pager(c, "Pages", total, &page, &size, data.PagerOptions{}).Changed {
			changes++
		}
	})
	if !tt.HasText("1–20 OF 312") {
		t.Fatalf("texts %q", tt.Texts())
	}
	click(t, tt, "NEXT")
	if page != 1 || !tt.HasText("21–40 OF 312") {
		t.Fatalf("▶ went to page %d", page)
	}
	click(t, tt, "50")
	if size != 50 || page != 0 || !tt.HasText("1–50 OF 312") {
		t.Fatalf("size 50 left page %d size %d", page, size)
	}
	page = 6
	tt.Frame()
	if !tt.HasText("301–312 OF 312") {
		t.Fatalf("last page texts %q", tt.Texts())
	}
	total = 40
	tt.Frame()
	if page != 0 || changes != 3 {
		t.Fatalf("a shorter total left page %d, %d changes", page, changes)
	}
}

func TestFilterBarTogglesAndClears(t *testing.T) {
	on := []bool{false, false, false}
	changes := 0
	tt := view(400, 100, newClock(), false, func(c *ui.Context) {
		if data.FilterBar(c, "Filters", []string{"Fire", "Ice", "Bolt"}, on, data.FilterBarOptions{}).Changed {
			changes++
		}
	})
	click(t, tt, "CLEAR")
	if changes != 0 {
		t.Fatal("CLEAR worked with nothing on")
	}
	click(t, tt, "Fire")
	click(t, tt, "Bolt")
	if !slices.Equal(on, []bool{true, false, true}) {
		t.Fatalf("chips left %v", on)
	}
	click(t, tt, "CLEAR")
	if slices.Contains(on, true) || changes != 3 {
		t.Fatalf("CLEAR left %v, %d changes", on, changes)
	}
}

func TestPropertyGridEditsEachKind(t *testing.T) {
	name, speed, fly := "Aria", 3.0, false
	changes := 0
	props := []data.Property{{Name: "Name", Text: &name}, {Name: "Speed", Number: &speed}, {Name: "Fly", Switch: &fly}}
	tt := view(400, 200, newClock(), false, func(c *ui.Context) {
		if data.PropertyGrid(c, "Hero", props, data.PropertyGridOptions{}).Changed {
			changes++
		}
	})
	// The switch sits past the 80-dot name column.
	r := find(t, tt, "Fly")
	tt.ClickAt(r.X+180, r.Y+r.H/2)
	if !fly {
		t.Fatalf("the switch row did not switch: %v", r)
	}
	n := find(t, tt, "Name")
	tt.ClickAt(n.X+380, n.Y+n.H/2)
	tt.Type("!")
	if name != "Aria!" || changes != 2 {
		t.Fatalf("typing gave %q, %d changes", name, changes)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("a row with two values did not panic")
		}
	}()
	view(400, 200, newClock(), false, func(c *ui.Context) {
		data.PropertyGrid(c, "Bad", []data.Property{{Name: "Both", Text: &name, Switch: &fly}}, data.PropertyGridOptions{})
	})
}

func TestPagerRefusesBadSizes(t *testing.T) {
	for _, sizes := range [][]int{{20, 0}, {20, -5}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Fatalf("sizes %v did not panic", sizes)
				}
			}()
			page, size := 0, 20
			view(400, 100, newClock(), false, func(c *ui.Context) {
				data.Pager(c, "Pages", 50, &page, &size, data.PagerOptions{Sizes: sizes})
			})
		}()
	}
}
