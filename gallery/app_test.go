package gallery

import (
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func open(path string) (*App, *ui.Tester) {
	a := New("/")
	if path != "" {
		a.Go(path)
	}
	return a, ui.NewTester(a.View, 1280, 800)
}

func click(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	if err := tt.Click(label); err != nil {
		t.Fatal(err)
	}
}

func TestBootThenStart(t *testing.T) {
	a, tt := open("")
	if a.started || !tt.HasText("Press start") {
		t.Fatal("gallery did not open on the boot screen")
	}
	tt.Key(0, ui.KeyEnter)
	if !a.started || !tt.HasText("Foundations") {
		t.Fatal("Enter did not start the gallery")
	}
}

func TestNavigation(t *testing.T) {
	a, tt := open("/")
	click(t, tt, "Foundations")
	if a.Path() != "/cat/foundations" {
		t.Fatalf("path %q after choosing a category", a.Path())
	}
	click(t, tt, "Palettes")
	if a.Path() != "/f/palettes" {
		t.Fatalf("path %q after choosing an entry", a.Path())
	}
	tt.Key(0, ui.KeyBracketRight)
	if a.Path() != "/f/canvas" {
		t.Fatalf("] went to %q", a.Path())
	}
	tt.Key(0, ui.KeyBracketLeft)
	tt.Key(0, ui.KeyBracketLeft)
	if a.Path() != Entries()[len(Entries())-1].Path() {
		t.Fatalf("[ from the first entry went to %q", a.Path())
	}
	click(t, tt, "home")
	if a.Path() != "/" {
		t.Fatalf("logo went to %q", a.Path())
	}
}

func TestPaletteKeysAndButtons(t *testing.T) {
	a, tt := open("/")
	tt.Key(0, ui.KeyP)
	if a.palette != 1 {
		t.Fatalf("P chose palette %d", a.palette)
	}
	tt.Key(ui.Shift, ui.KeyP)
	tt.Key(ui.Shift, ui.KeyP)
	if a.palette != len(pixel.Palettes)-1 {
		t.Fatalf("Shift+P twice chose palette %d", a.palette)
	}
	click(t, tt, "next palette")
	if a.palette != 0 {
		t.Fatalf("next palette chose %d", a.palette)
	}
}

func TestSearch(t *testing.T) {
	a, tt := open("/")
	tt.Key(0, ui.KeySlash)
	if !tt.Focused("Search components") {
		t.Fatal("/ did not focus search")
	}
	tt.Type("dith")
	if !tt.HasText("1 results") || !tt.HasText("Dither") {
		t.Fatalf("results: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyP)
	if a.palette != 0 {
		t.Fatal("P changed the palette while typing")
	}
	tt.Key(0, ui.KeyEnter)
	if a.Path() != "/f/dither" || a.query != "" {
		t.Fatalf("Enter left path %q query %q", a.Path(), a.query)
	}
	tt.Key(0, ui.KeySlash)
	tt.Type("zz")
	tt.Key(0, ui.KeyEscape)
	if a.query != "" {
		t.Fatalf("Escape left query %q", a.query)
	}
}

func TestDemoLogAndCopy(t *testing.T) {
	a, tt := open("/f/icons")
	click(t, tt, "heart")
	if !slices.Contains(a.logs["icons"], "icon heart") {
		t.Fatalf("log %q", a.logs["icons"])
	}
	click(t, tt, "COPY CODE")
	e, _ := find("icons")
	if tt.Clipboard() != strings.ReplaceAll(e.Code(), "\t", "  ") {
		t.Fatal("COPY CODE copied something else")
	}
}

func TestLostPageAndLocale(t *testing.T) {
	a, tt := open("/nowhere")
	if !tt.HasText("No page at /nowhere") {
		t.Fatalf("lost page texts %q", tt.Texts())
	}
	click(t, tt, "language")
	if a.locale != pixel.Zh || !tt.HasText("重试") {
		t.Fatal("language toggle did not switch built-in copy")
	}
	click(t, tt, "重试")
	if a.Path() != "/" {
		t.Fatalf("continue went to %q", a.Path())
	}
}

func TestGoSkipsBootAndSearch(t *testing.T) {
	a := New("/")
	a.query = "dith"
	a.Go("/f/icons")
	if !a.started || a.Path() != "/f/icons" || a.query != "" {
		t.Fatalf("Go left started %v path %q query %q", a.started, a.Path(), a.query)
	}
}

func TestSearchResultClearsQuery(t *testing.T) {
	a, tt := open("/")
	tt.Key(0, ui.KeySlash)
	tt.Type("dith")
	click(t, tt, "Dither")
	if a.query != "" || a.Path() != "/f/dither" || !tt.HasText("COPY CODE") {
		t.Fatalf("query %q path %q", a.query, a.Path())
	}
}

func TestPathKindMustMatch(t *testing.T) {
	_, tt := open("/c/palettes")
	if !tt.HasText("No page at /c/palettes") {
		t.Fatal("a foundation opened under /c/")
	}
}

func TestShortcutsYieldToFields(t *testing.T) {
	a, tt := open("/f/editor")
	click(t, tt, "Name")
	tt.TypeKey(0, ui.KeyP, "p")
	tt.TypeKey(0, ui.KeyBracketRight, "]")
	if a.palette != 0 || a.Path() != "/f/editor" {
		t.Fatalf("typing changed palette %d or path %q", a.palette, a.Path())
	}
}

func TestFoundationTilesShowFocus(t *testing.T) {
	for path, tile := range map[string]string{"/f/palettes": "Spring Day", "/f/sprites": "slime"} {
		_, tt := open(path)
		r, _ := tt.Find(tile)
		area := func() []uint32 {
			img := tt.Image()
			var out []uint32
			for y := int(r.Y) - 6; y < int(r.Y+r.H)+6; y++ {
				for x := int(r.X) - 30; x < int(r.X+r.W)+6; x++ {
					c := img.RGBAAt(x, y)
					out = append(out, uint32(c.R)<<16|uint32(c.G)<<8|uint32(c.B))
				}
			}
			return out
		}
		before := area()
		for range 40 {
			if tt.Focused(tile) {
				break
			}
			tt.Key(0, ui.KeyTab)
		}
		if !tt.Focused(tile) || slices.Equal(before, area()) {
			t.Fatalf("%s: focus on %s left no mark", path, tile)
		}
	}
}

func TestInputDemosLogEveryClick(t *testing.T) {
	a, tt := open("/c/066-button")
	for _, l := range []string{"Small", "Save", "BIG"} {
		click(t, tt, l)
		if !slices.Contains(a.logs["066-button"], l) {
			t.Fatalf("clicking %s logged %q", l, a.logs["066-button"])
		}
	}
	a, tt = open("/c/067-iconbutton")
	click(t, tt, "Delete")
	if !slices.Contains(a.logs["067-iconbutton"], "Delete") {
		t.Fatalf("Delete logged %q", a.logs["067-iconbutton"])
	}
	for slug, field := range map[string]string{"082-textinput": "Hero name", "083-passwordinput": "Password", "085-searchfield": "Find items"} {
		a, tt = open("/c/" + slug)
		click(t, tt, field)
		tt.Type("ab")
		if len(a.logs[slug]) == 0 {
			t.Fatalf("typing into %s logged nothing", slug)
		}
	}
}

func TestControlDemosLogTheirEvents(t *testing.T) {
	for slug, act := range map[string]func(tt *ui.Tester){
		"088-combobox":        func(tt *ui.Tester) { click(t, tt, "Item"); tt.Type("p") },
		"089-autocomplete":    func(tt *ui.Tester) { click(t, tt, "Spell"); tt.Type("f") },
		"094-pininput":        func(tt *ui.Tester) { click(t, tt, "Vault code"); tt.Type("1") },
		"095-nameentry":       func(tt *ui.Tester) { click(t, tt, "Initials") },
		"096-virtualkeyboard": func(tt *ui.Tester) { click(t, tt, "Message") },
		"103-joystick":        func(tt *ui.Tester) { click(t, tt, "Stick"); tt.Key(0, ui.KeyRight) },
		"105-slotpicker":      func(tt *ui.Tester) { click(t, tt, "Reel 1"); tt.Key(0, ui.KeyDown) },
		"101-fieldset":        func(tt *ui.Tester) { click(t, tt, "Guild name"); tt.Type("a") },
	} {
		a, tt := open("/c/" + slug)
		act(tt)
		if len(a.logs[slug]) == 0 {
			t.Errorf("%s logged nothing", slug)
		}
	}
}

func TestDataDemosLogSubmits(t *testing.T) {
	for slug, row := range map[string]string{"152-table": "Brom, Monk, 5", "154-outline": "Region 3", "174-searchresults": "Fira"} {
		a, tt := open("/c/" + slug)
		click(t, tt, row)
		n := len(a.logs[slug])
		tt.Key(0, ui.KeyEnter)
		if len(a.logs[slug]) == n {
			t.Errorf("%s logged no submit: %q", slug, a.logs[slug])
		}
	}
	a, tt := open("/c/153-tree")
	click(t, tt, "Items")
	n := len(a.logs["153-tree"])
	tt.Key(0, ui.KeyRight)
	if len(a.logs["153-tree"]) == n {
		t.Errorf("opening a branch logged nothing: %q", a.logs["153-tree"])
	}
}

func TestFieldDemoLogsSubmit(t *testing.T) {
	a, tt := open("/c/100-field")
	// The group shares its name; the input sits below.
	r, _ := tt.Find("Call sign")
	tt.ClickAt(r.X+20, r.Y+58)
	tt.Type("Ace")
	n := len(a.logs["100-field"])
	tt.Key(0, ui.KeyEnter)
	if len(a.logs["100-field"]) == n {
		t.Fatal("Enter in the field logged nothing")
	}
}

func TestDataDemosLogTheirEvents(t *testing.T) {
	for slug, act := range map[string]func(tt *ui.Tester){
		"151-list":          func(tt *ui.Tester) { click(t, tt, "Bat") },
		"152-table":         func(tt *ui.Tester) { click(t, tt, "LV") },
		"153-tree":          func(tt *ui.Tester) { click(t, tt, "Items"); tt.Key(0, ui.KeyRight) },
		"154-outline":       func(tt *ui.Tester) { click(t, tt, "Region 2") },
		"174-searchresults": func(tt *ui.Tester) { click(t, tt, "Fira") },
		"155-gridview":      func(tt *ui.Tester) { click(t, tt, "loot 3") },
		"160-accordion":     func(tt *ui.Tester) { click(t, tt, "Where is the key?") },
		"162-carousel":      func(tt *ui.Tester) { click(t, tt, "NEXT") },
		"167-pager":         func(tt *ui.Tester) { click(t, tt, "NEXT") },
		"170-sortable":      func(tt *ui.Tester) { click(t, tt, "Quest order"); tt.Key(ui.Alt, ui.KeyDown) },
		"171-transfer":      func(tt *ui.Tester) { click(t, tt, "RIGHT") },
		"173-filterbar":     func(tt *ui.Tester) { click(t, tt, "FIRE") },
	} {
		a, tt := open("/c/" + slug)
		act(tt)
		if len(a.logs[slug]) == 0 {
			t.Errorf("%s logged nothing", slug)
		}
	}
}
