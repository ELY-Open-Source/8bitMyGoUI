package input_test

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestRadioGroupArrowsChooseAfterAClick(t *testing.T) {
	season := 0
	tt := view(400, 200, nil, func(c *ui.Context) {
		input.RadioGroup(c, []string{"Spring", "Summer", "Autumn"}, &season, input.RadioGroupOptions{})
	})
	click(t, tt, "Autumn")
	tt.Key(0, ui.KeyUp)
	if season != 1 {
		t.Fatalf("Up after a click chose %d, want 1", season)
	}
	// A press dragged off a radio chooses nothing.
	r := find(t, tt, "Spring")
	tt.Press(r.X+4, r.Y+r.H/2)
	tt.Move(r.X+4, r.Y+200)
	tt.Release(r.X+4, r.Y+200)
	if season != 1 {
		t.Fatalf("a press dragged off chose %d", season)
	}
}

func TestSelectPicksByClickAndKeys(t *testing.T) {
	class, changes, disabled := "Knight", 0, false
	tt := view(400, 300, nil, func(c *ui.Context) {
		if input.Select(c, "Class", &class, []string{"Knight", "Mage", "Thief"}, input.SelectOptions{Disabled: disabled}).Changed {
			changes++
		}
	})
	click(t, tt, "Class")
	mage := find(t, tt, "Mage")
	if !has(slots(t, tt.Image(), find(t, tt, "Knight")), pixel.Accent) {
		t.Fatal("the chosen option is not marked")
	}
	tt.Move(mage.X+mage.W/2, mage.Y+mage.H/2)
	if !has(slots(t, tt.Image(), find(t, tt, "Mage")), pixel.Lit) {
		t.Fatal("the option under the pointer is not lit")
	}
	click(t, tt, "Mage")
	if class != "Mage" || changes != 1 || tt.HasText("Thief") {
		t.Fatalf("click chose %q with %d changes; popup open %v", class, changes, tt.HasText("Thief"))
	}
	if _, ok := tt.Find("Thief"); ok {
		t.Fatal("the popup stayed open")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if class != "Thief" || changes != 2 {
		t.Fatalf("Down, Down, Enter chose %q", class)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEscape)
	if _, ok := tt.Find("Mage"); ok || !tt.Focused("Class") {
		t.Fatal("Escape left the popup open or the focus elsewhere")
	}
	disabled = true
	tt.Frame()
	click(t, tt, "Class")
	if _, ok := tt.Find("Mage"); ok {
		t.Fatal("a disabled select opened")
	}
}

func TestComboboxReachesEveryOption(t *testing.T) {
	v := ""
	letters := []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9"}
	tt := view(400, 400, nil, func(c *ui.Context) {
		input.Combobox(c, "Pick", &v, letters, input.ComboboxOptions{})
	})
	click(t, tt, "Pick")
	tt.Key(0, ui.KeyDown)
	if _, ok := tt.Find("A9"); ok {
		t.Fatal("more than eight rows show")
	}
	// A resting pointer must not steal the light.
	a1 := find(t, tt, "A1")
	tt.Move(a1.X+a1.W/2, a1.Y+a1.H/2)
	for range 8 {
		tt.Key(0, ui.KeyDown)
	}
	_, first := tt.Find("A1")
	if _, last := tt.Find("A9"); first || !last {
		t.Fatal("the rows did not follow the lit option")
	}
	tt.Key(0, ui.KeyEnter)
	if v != "A9" {
		t.Fatalf("eight Downs and Enter took %q", v)
	}
	tt.Key(0, ui.KeyDown)
	// With the pointer on the first row, the wheel moves.
	r := find(t, tt, "A1")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 4*2*pixel.LineHeight)
	tt.Key(0, ui.KeyEnter)
	if v != "A5" {
		t.Fatalf("four rows of wheel then Enter took %q", v)
	}
}

func TestBadOptionsPanic(t *testing.T) {
	for name, build := range map[string]func(c *ui.Context){
		"Select Width": func(c *ui.Context) {
			s := "A"
			input.Select(c, "S", &s, []string{"A"}, input.SelectOptions{Width: -1})
		},
		"Knob Step": func(c *ui.Context) {
			v := 0.5
			input.Knob(c, "K", &v, input.KnobOptions{Step: -0.1})
		},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: a negative value did not panic", name)
				}
			}()
			view(200, 100, nil, build)
		}()
	}
}

func TestInertClosesPopups(t *testing.T) {
	for name, build := range map[string]func(c *ui.Context){
		"Class": func(c *ui.Context) {
			class := "A"
			input.Select(c, "Class", &class, []string{"A", "B"}, input.SelectOptions{})
		},
		"Gear": func(c *ui.Context) {
			var gear []string
			input.MultiSelect(c, "Gear", &gear, []string{"B", "C"}, input.MultiSelectOptions{})
		},
	} {
		inert := false
		tt := view(400, 300, nil, func(c *ui.Context) {
			s := pixel.Current(c)
			s.Inert = inert
			pixel.Scope(c, s, func() { build(c) })
		})
		click(t, tt, name)
		find(t, tt, "B")
		inert = true
		tt.Frame()
		if _, ok := tt.Find("B"); ok {
			t.Fatalf("%s stayed open under Inert", name)
		}
		inert = false
		tt.Frame()
		if _, ok := tt.Find("B"); ok {
			t.Fatalf("%s opened again after Inert", name)
		}
	}
}

func TestMultiSelectChecksAndShowsTags(t *testing.T) {
	var gear []string
	changes := 0
	tt := view(400, 300, nil, func(c *ui.Context) {
		if input.MultiSelect(c, "Gear", &gear, []string{"Sword", "Shield", "Potion"}, input.MultiSelectOptions{}).Changed {
			changes++
		}
	})
	trigger := find(t, tt, "Gear")
	empty := slots(t, tt.Image(), trigger)
	click(t, tt, "Gear")
	click(t, tt, "Shield")
	click(t, tt, "Potion")
	click(t, tt, "Shield")
	if !slices.Equal(gear, []string{"Potion"}) || changes != 3 {
		t.Fatalf("checks left %q after %d changes", gear, changes)
	}
	tt.Key(0, ui.KeyEscape)
	if _, ok := tt.Find("Sword"); ok || !tt.Focused("Gear") {
		t.Fatal("Escape left the popup open or the focus in it")
	}
	if equal(empty, slots(t, tt.Image(), trigger)) {
		t.Fatal("the chosen option shows no tag")
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Sword") {
		t.Fatal("opening did not focus the first option")
	}
	tt.Key(0, ui.KeySpace)
	if !slices.Contains(gear, "Sword") {
		t.Fatalf("Space left %q", gear)
	}
}

func TestComboboxFiltersAndChooses(t *testing.T) {
	fruit := ""
	chosen := 0
	tt := view(400, 300, nil, func(c *ui.Context) {
		if input.Combobox(c, "Fruit", &fruit, []string{"Apple", "Apricot", "Banana", "Cherry"}, input.ComboboxOptions{}).Chosen {
			chosen++
		}
		input.Button(c, "after", input.ButtonOptions{})
	})
	click(t, tt, "Fruit")
	tt.Type("ap")
	if _, ok := tt.Find("Banana"); ok {
		t.Fatal("Banana does not match ap")
	}
	find(t, tt, "Apricot")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if fruit != "Apricot" || chosen != 1 {
		t.Fatalf("Down, Enter took %q", fruit)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("b")
	click(t, tt, "Banana")
	if fruit != "Banana" || !tt.Focused("Fruit") {
		t.Fatalf("a click took %q; focus kept %v", fruit, tt.Focused("Fruit"))
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("c")
	find(t, tt, "Cherry")
	tt.Key(0, ui.KeyEscape)
	if _, ok := tt.Find("Cherry"); ok {
		t.Fatal("Escape left the list open")
	}
	tt.Type("h")
	find(t, tt, "Cherry")
	tt.Key(0, ui.KeyTab)
	if _, ok := tt.Find("Cherry"); ok || !tt.Focused("after") {
		t.Fatal("leaving the field left the list open")
	}
}

func TestAutocompleteGhostsAndAccepts(t *testing.T) {
	monster, accepted := "", 0
	tt := view(400, 200, nil, func(c *ui.Context) {
		if input.Autocomplete(c, "Monster", &monster, []string{"Dragon", "Drake", "Dungeon"}, input.AutocompleteOptions{}).Accepted {
			accepted++
		}
		input.Button(c, "after", input.ButtonOptions{})
	})
	click(t, tt, "Monster")
	tt.Type("zz")
	if has(slots(t, tt.Image(), find(t, tt, "Monster")), pixel.Faint) {
		t.Fatal("a ghost shows without a completion")
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("dra")
	if !has(slots(t, tt.Image(), find(t, tt, "Monster")), pixel.Faint) {
		t.Fatal("no ghost completion")
	}
	tt.Key(0, ui.KeyTab)
	if monster != "Dragon" || accepted != 1 || !tt.Focused("Monster") {
		t.Fatalf("Tab took %q, focus kept %v", monster, tt.Focused("Monster"))
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("dun")
	tt.Key(0, ui.KeyEnter)
	if monster != "Dungeon" || accepted != 2 {
		t.Fatalf("Enter took %q", monster)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("after") {
		t.Fatal("Tab without a completion kept the focus")
	}
}

func TestTokenFieldAddsAndRemoves(t *testing.T) {
	var tags []string
	tt := view(400, 200, nil, func(c *ui.Context) {
		input.TokenField(c, "Tags", &tags, input.TokenFieldOptions{})
	})
	click(t, tt, "Tags")
	tt.Type("fire")
	tt.Key(0, ui.KeyEnter)
	tt.Type("ice,wind,")
	tt.Type("fire")
	tt.Key(0, ui.KeyEnter)
	if !slices.Equal(tags, []string{"fire", "ice", "wind"}) {
		t.Fatalf("Enter and commas gave %q", tags)
	}
	tt.Type("x")
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyBackspace)
	if !slices.Equal(tags, []string{"fire", "ice"}) {
		t.Fatalf("Backspace on empty gave %q", tags)
	}
	click(t, tt, "REMOVE fire")
	if !slices.Equal(tags, []string{"ice"}) {
		t.Fatalf("the × gave %q", tags)
	}
}

func TestSegmentedClicksAndArrows(t *testing.T) {
	mode, changes := 0, 0
	tt := view(400, 200, nil, func(c *ui.Context) {
		if input.Segmented(c, "Mode", []string{"EASY", "NORMAL", "HARD"}, &mode, input.SegmentedOptions{}).Changed {
			changes++
		}
	})
	click(t, tt, "HARD")
	tt.Key(0, ui.KeyLeft)
	if mode != 1 || changes != 2 {
		t.Fatalf("click and Left gave %d after %d changes", mode, changes)
	}
	// One dot in from each segment's corner is its fill.
	fill := func(l string) pixel.Slot {
		r := find(t, tt, l)
		return slots(t, tt.Image(), ui.Rect{X: r.X + 4, Y: r.Y + 4, W: 2, H: 2})[0]
	}
	if fill("NORMAL") != pixel.Accent || fill("HARD") == pixel.Accent {
		t.Fatal("the chosen segment is not the Accent one")
	}
}
