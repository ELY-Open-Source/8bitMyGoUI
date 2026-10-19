package data_test

import (
	"fmt"
	"image"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func newClock() *clock { return &clock{time.Unix(3_000_000, 0)} }

// view builds on a fake clock; inert veils every control.
func view(w, h int, clk *clock, inert bool, build func(c *ui.Context)) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: clk.now, Inert: inert})
		build(c)
	}, w, h)
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

// slot is the palette slot of pixel x, y.
func slot(t *testing.T, img *image.RGBA, x, y int) pixel.Slot {
	t.Helper()
	c := img.RGBAAt(x, y)
	for s, pc := range pixel.SpringDay.Colors {
		if pc.R == c.R && pc.G == c.G && pc.B == c.B {
			return pixel.Slot(s)
		}
	}
	t.Fatalf("pixel %d,%d is %v, off the palette", x, y, c)
	return 0
}

// keeps fails on off-palette colors or split dots.
func keeps(t *testing.T, img *image.RGBA) {
	t.Helper()
	b := img.Bounds()
	for y := 0; y+1 < b.Dy(); y += 2 {
		for x := 0; x+1 < b.Dx(); x += 2 {
			slot(t, img, x, y)
			c := img.RGBAAt(x, y)
			if img.RGBAAt(x+1, y) != c || img.RGBAAt(x, y+1) != c || img.RGBAAt(x+1, y+1) != c {
				t.Fatalf("dot at %d,%d is split", x/2, y/2)
			}
		}
	}
}

var heroes = []string{"Aria", "Brom", "Cyra", "Dax", "Eko"}

func party(c *ui.Context, sel *int, o data.ListOptions, changes, sent *int) {
	r := data.List(c, "Party", len(heroes), sel, o, func(i int) string { return heroes[i] }, func(i int) {
		pixel.Text(c, heroes[i], pixel.TextOptions{})
	})
	if r.Changed {
		*changes++
	}
	if r.Submitted {
		*sent++
	}
}

func TestListChoosesByClickAndKeys(t *testing.T) {
	sel, changes, sent := 0, 0, 0
	tt := view(400, 300, newClock(), false, func(c *ui.Context) { party(c, &sel, data.ListOptions{}, &changes, &sent) })
	click(t, tt, "Cyra")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if sel != 3 || changes != 2 || sent != 1 {
		t.Fatalf("click, Down, Enter: sel %d, %d changes, %d submits", sel, changes, sent)
	}
	keeps(t, tt.Image())
}

func TestListOffTakesNoChoice(t *testing.T) {
	for _, inert := range []bool{false, true} {
		sel, changes, sent := 1, 0, 0
		tt := view(400, 300, newClock(), inert, func(c *ui.Context) {
			party(c, &sel, data.ListOptions{Disabled: !inert}, &changes, &sent)
		})
		click(t, tt, "Dax")
		tt.Key(0, ui.KeyDown)
		if sel != 1 || changes != 0 {
			t.Fatalf("inert %v: sel %d after a click and Down", inert, sel)
		}
	}
}

func roster(c *ui.Context, sort *data.Sort, sel *int, sorted *int) {
	cols := []data.Column{{Title: "Name"}, {Title: "LV", Width: 30}}
	if data.Table(c, "Roster", cols, len(heroes), sort, sel, data.TableOptions{}, func(r, col int) string {
		if col == 0 {
			return heroes[r]
		}
		return fmt.Sprint(r + 1)
	}).Sorted {
		*sorted++
	}
}

func TestTableSortsAndChooses(t *testing.T) {
	sort, sel, sorted := data.Sort{}, -1, 0
	tt := view(500, 300, newClock(), false, func(c *ui.Context) { roster(c, &sort, &sel, &sorted) })
	click(t, tt, "LV")
	click(t, tt, "LV")
	if sort != (data.Sort{Column: 1, Desc: true}) || sorted != 2 {
		t.Fatalf("two clicks on LV gave %+v", sort)
	}
	click(t, tt, "Name")
	if sort != (data.Sort{Column: 0}) {
		t.Fatalf("a click on Name gave %+v", sort)
	}
	click(t, tt, "Cyra, 3")
	tt.Key(0, ui.KeyDown)
	r := find(t, tt, "Dax, 4")
	if bar := slot(t, tt.Image(), int(r.X)+r2i(r.W)-4, int(r.Y)+4); sel != 3 || bar != pixel.Accent {
		t.Fatalf("click on Cyra, Down chose %d, painted %v", sel, bar)
	}
	keeps(t, tt.Image())
}

func TestTableZebraDithersOddRows(t *testing.T) {
	sort, sel, sorted := data.Sort{}, -1, 0
	tt := view(500, 300, newClock(), false, func(c *ui.Context) { roster(c, &sort, &sel, &sorted) })
	shaded := func(label string) int {
		r, img, n := find(t, tt, label), tt.Image(), 0
		for x := int(r.X) + 120; x < int(r.X)+160; x += 2 {
			if slot(t, img, x, int(r.Y)+4) == pixel.Shade {
				n++
			}
		}
		return n
	}
	if even, odd := shaded("Aria, 1"), shaded("Brom, 2"); even != 0 || odd == 0 {
		t.Fatalf("Shade dots: even row %d, odd row %d", even, odd)
	}
}

var bag = []data.TreeNode{{Label: "Bag", Children: []data.TreeNode{{Label: "Potion"}, {Label: "Ether"}}}, {Label: "Keys"}}

func TestTreeOpensAndSteers(t *testing.T) {
	sel := "1"
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		data.Tree(c, "Items", bag, &sel, data.TreeOptions{})
	})
	r := find(t, tt, "Items")
	tt.ClickAt(r.X+60, r.Y+8)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	if sel != "0/1" {
		t.Fatalf("click, Right, Right, Down chose %q", sel)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyDown)
	if sel != "1" {
		t.Fatalf("Left, Left, Down chose %q", sel)
	}
	tt.ClickAt(r.X+4, r.Y+8)
	if sel != "0" || !tt.HasText("Items") {
		t.Fatalf("a click on the arrow chose %q", sel)
	}
	tt.Key(0, ui.KeyDown)
	if sel != "0/0" {
		t.Fatalf("the arrow did not open Bag: Down chose %q", sel)
	}
}

func TestOutlineOpensAndChooses(t *testing.T) {
	sel := "/"
	kids := map[string][]string{"/": {"/a", "/b"}, "/a": {"/a/x"}}
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		data.Outline(c, "Files", []string{"/"}, func(k string) []string { return kids[k] }, func(k string) string { return k }, &sel, data.OutlineOptions{})
	})
	click(t, tt, "/")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if sel != "/a/x" || !shows(tt, "/b") {
		t.Fatalf("Right, Down, Right, Right chose %q", sel)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyLeft)
	if sel != "/a" || shows(tt, "/a/x") {
		t.Fatalf("Left, Left chose %q", sel)
	}
	click(t, tt, "/b")
	r := find(t, tt, "/b")
	if bar := slot(t, tt.Image(), int(r.X)+4, int(r.Y)+4); sel != "/b" || bar != pixel.Accent {
		t.Fatalf("a click chose %q, painted %v", sel, bar)
	}
	keeps(t, tt.Image())
}

func TestGridViewMovesInTwoDimensions(t *testing.T) {
	sel, changes, sent := 0, 0, 0
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		r := data.GridView(c, "Tiles", 30, &sel, data.GridViewOptions{CellWidth: 40, CellHeight: 40}, func(i int) {
			ui.Box(c).Grow(1).Label(fmt.Sprintf("tile %d", i))
		})
		if r.Changed {
			changes++
		}
		if r.Submitted {
			sent++
		}
	})
	click(t, tt, "tile 1")
	tt.Key(0, ui.KeyDown)
	a, b := find(t, tt, "tile 1"), find(t, tt, fmt.Sprintf("tile %d", sel))
	if a.X != b.X || b.Y <= a.Y || changes != 2 {
		t.Fatalf("Down from tile 1 chose %d at %v, from %v; %d changes", sel, b, a, changes)
	}
	tt.Key(0, ui.KeyEnter)
	if sent != 1 {
		t.Fatal("Enter did not submit")
	}
	keeps(t, tt.Image())
}

func TestInfiniteListAsksOncePerLength(t *testing.T) {
	n, asks, done := 3, 0, false
	tt := view(300, 300, newClock(), false, func(c *ui.Context) {
		if data.InfiniteList(c, "Log", n, data.InfiniteListOptions{Done: done}, func(i int) {
			pixel.Text(c, fmt.Sprintf("line %d", i), pixel.TextOptions{})
		}).More {
			asks++
		}
	})
	tt.Frame()
	if asks != 1 || !tt.HasText("NOW LOADING") {
		t.Fatalf("%d asks with LOADING %v", asks, tt.HasText("NOW LOADING"))
	}
	n = 5
	tt.Frame()
	tt.Frame()
	if asks != 2 {
		t.Fatalf("%d asks after more rows came", asks)
	}
	done = true
	tt.Frame()
	if tt.HasText("NOW LOADING") || asks != 2 {
		t.Fatal("Done kept the LOADING row")
	}
}

func TestSearchResultsCountAndMark(t *testing.T) {
	sel, opened := 0, 0
	hits := []data.Hit{{Title: "Fire Sword", Text: "A blade of fire"}, {Title: "Ice Bow", Text: "Shoots ice"}}
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		if data.SearchResults(c, "Results", "fire", hits, &sel, data.SearchResultsOptions{}).Submitted {
			opened++
		}
	})
	if !tt.HasText("FOUND: 2") {
		t.Fatalf("texts %q", tt.Texts())
	}
	r, img := find(t, tt, "Fire Sword"), tt.Image()
	if marks := slot(t, img, int(r.X)+4, int(r.Y)+6); marks != pixel.Warn {
		t.Fatalf("the match is marked %v", marks)
	}
	click(t, tt, "Ice Bow")
	tt.Key(0, ui.KeyEnter)
	if sel != 1 || opened != 1 {
		t.Fatalf("sel %d, %d opened", sel, opened)
	}
	if !slices.Contains(tt.Texts(), "Shoots ice") {
		t.Fatal("snippets are not read")
	}
}

func TestListsKeepWholeDotsWhenScrolledByHalfADot(t *testing.T) {
	sel, sort, row, file := 1, data.Sort{}, 1, "row 1"
	many := make([]string, 40)
	for i := range many {
		many[i] = fmt.Sprintf("row %d", i)
	}
	tt := view(600, 300, newClock(), false, func(c *ui.Context) {
		ui.Row(c).Children(func() {
			data.List(c, "Rows", len(many), &sel, data.ListOptions{Height: 60}, func(i int) string { return many[i] }, func(i int) {
				pixel.Text(c, many[i], pixel.TextOptions{})
			}).Element.Width(100)
			data.Table(c, "Grid", []data.Column{{Title: "Row"}}, len(many), &sort, &row, data.TableOptions{Height: 60}, func(r, _ int) string { return many[r] })
			data.Outline(c, "Tree", many, func(string) []string { return nil }, func(k string) string { return k }, &file, data.OutlineOptions{Height: 40})
		})
	})
	for _, l := range []string{"Tree", "Grid", "Rows"} {
		r := find(t, tt, l)
		tt.Move(r.X+20, r.Y+20)
		tt.Scroll(r.X+20, r.Y+20, 0, 1)
	}
	// The pointer stays on a row, hovering it.
	keeps(t, tt.Image())
}

func TestOffListsIgnoreTheWheel(t *testing.T) {
	many := make([]string, 40)
	for i := range many {
		many[i] = fmt.Sprintf("row %d", i)
	}
	for _, inert := range []bool{false, true} {
		sel, sort, row, file, tile := 0, data.Sort{}, 0, "row 0", 0
		tt := view(800, 300, newClock(), inert, func(c *ui.Context) {
			ui.Row(c).Children(func() {
				data.List(c, "Rows", len(many), &sel, data.ListOptions{Height: 60, Disabled: !inert}, func(i int) string { return many[i] }, func(i int) {
					pixel.Text(c, many[i], pixel.TextOptions{})
				})
				data.Table(c, "Grid", []data.Column{{Title: "Row"}}, len(many), &sort, &row, data.TableOptions{Height: 60, Disabled: !inert}, func(r, _ int) string { return many[r] })
				data.Outline(c, "Tree", many, func(string) []string { return nil }, func(k string) string { return k }, &file, data.OutlineOptions{Height: 60, Disabled: !inert})
				ui.Box(c).Width(120).Children(func() {
					data.GridView(c, "Cells", 90, &tile, data.GridViewOptions{CellWidth: 20, CellHeight: 20, Height: 40, Disabled: !inert}, func(i int) {
						ui.Box(c).Grow(1).Label(fmt.Sprintf("cell %d", i))
					})
				})
			})
		})
		for _, l := range []string{"Rows", "Grid", "Tree", "Cells"} {
			r := find(t, tt, l)
			before := tt.Texts()
			tt.Move(r.X+20, r.Y+30)
			tt.Scroll(r.X+20, r.Y+30, 0, 80)
			if !slices.Equal(before, tt.Texts()) {
				t.Fatalf("inert %v: the wheel moved %s", inert, l)
			}
		}
	}
}

func TestInertInfiniteListIgnoresTheWheel(t *testing.T) {
	asks := 0
	tt := view(300, 300, newClock(), true, func(c *ui.Context) {
		if data.InfiniteList(c, "Log", 40, data.InfiniteListOptions{Height: 60}, func(i int) {
			pixel.Text(c, fmt.Sprintf("line %d", i), pixel.TextOptions{})
		}).More {
			asks++
		}
	})
	r := find(t, tt, "Log")
	tt.Move(r.X+20, r.Y+30)
	for range 20 {
		tt.Scroll(r.X+20, r.Y+30, 0, 200)
	}
	if !shows(tt, "line 0") || asks != 0 {
		t.Fatalf("the wheel moved an inert list to %q; %d asks", tt.Texts(), asks)
	}
}

func TestListNamesRowsForReadersAndTyping(t *testing.T) {
	fruit := []string{"Apple", "Banana", "Cherry"}
	sel := 0
	tt := view(300, 200, newClock(), false, func(c *ui.Context) {
		data.List(c, "Fruit", len(fruit), &sel, data.ListOptions{}, func(i int) string { return fruit[i] + " row" }, func(i int) {
			pixel.Text(c, fruit[i], pixel.TextOptions{})
		})
	})
	find(t, tt, "Banana row")
	// Type-ahead reads key presses, as a keyboard sends them.
	click(t, tt, "Apple")
	tt.Key(0, ui.KeyC)
	tt.Key(0, ui.KeyH)
	if sel != 2 {
		t.Fatalf("typing ch chose row %d", sel)
	}
}

// MyGo would halve a disabled list's colors off the palette.
func TestOffListsKeepThePalette(t *testing.T) {
	fruit := []string{"Apple", "Banana"}
	for _, tc := range []struct {
		name            string
		disabled, inert bool
	}{{"disabled", true, false}, {"inert", false, true}} {
		sel := 0
		tt := view(300, 200, newClock(), tc.inert, func(c *ui.Context) {
			data.List(c, "Fruit", len(fruit), &sel, data.ListOptions{Disabled: tc.disabled}, func(i int) string { return fruit[i] }, func(i int) {
				pixel.Text(c, fruit[i], pixel.TextOptions{})
			})
		})
		r, img := find(t, tt, "Fruit"), tt.Image()
		for y := int(r.Y); y < int(r.Y+r.H); y++ {
			for x := int(r.X); x < int(r.X+r.W); x++ {
				if c := img.RGBAAt(x, y); !slices.ContainsFunc(pixel.SpringDay.Colors[:], func(p ui.Color) bool { return p.R == c.R && p.G == c.G && p.B == c.B }) {
					t.Fatalf("%s list shows %v at %d, %d, off the palette", tc.name, c, x, y)
				}
			}
		}
	}
}
