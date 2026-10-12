package navigation_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/navigation"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var commands = []navigation.MenuItem{{Label: "FIGHT"}, {Label: "MAGIC", Disabled: true}, {Label: "ITEM"}, {Label: "RUN"}}

func TestCommandMenuChoosesAndCancels(t *testing.T) {
	chosen, canceled := -1, 0
	tt := view(300, 200, newClock(), false, func(c *ui.Context) {
		r := navigation.CommandMenu(c, "Command", commands, navigation.CommandMenuOptions{})
		if r.Chosen >= 0 {
			chosen = r.Chosen
		}
		if r.Canceled {
			canceled++
		}
	})
	r := find(t, tt, "Command")
	tt.ClickAt(r.X+20, r.Y+8+2*16*2+8)
	if chosen != 2 {
		t.Fatalf("a click on ITEM chose %d", chosen)
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	if chosen != 2 {
		t.Fatal("Enter took the disabled MAGIC")
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyEscape)
	if chosen != 0 || canceled != 1 {
		t.Fatalf("Up, Enter chose %d; Escape canceled %d times", chosen, canceled)
	}
}

func TestDropdownMenuKeysAndClicks(t *testing.T) {
	chosen, disabled := -1, false
	items := []navigation.MenuItem{{Label: "Save", Icon: icon.Save, Keys: "CTRL+S"}, {Separator: true}, {Label: "Load", Disabled: true}, {Label: "Quit"}}
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		if r := navigation.DropdownMenu(c, "File", items, navigation.DropdownMenuOptions{Disabled: disabled}); r.Chosen >= 0 {
			chosen = r.Chosen
		}
	})
	click(t, tt, "File")
	find(t, tt, "Quit")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if chosen != 3 || shows(tt, "Quit") {
		t.Fatalf("Down skipped to %d, or the menu stayed", chosen)
	}
	tt.Key(0, ui.KeyDown)
	click(t, tt, "Save")
	if chosen != 0 {
		t.Fatalf("a click chose %d", chosen)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEscape)
	if shows(tt, "Quit") || !tt.Focused("File") {
		t.Fatal("Escape left the menu open or lost the focus")
	}
	disabled = true
	tt.Frame()
	click(t, tt, "File")
	if shows(tt, "Quit") {
		t.Fatal("a disabled dropdown opened")
	}
}

func TestMenusIgnoreDisabledItems(t *testing.T) {
	chosen := -1
	items := []navigation.MenuItem{{Label: "Save"}, {Label: "Load", Disabled: true}, {Label: "Quit"}}
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		input.Button(c, "before", input.ButtonOptions{})
		if r := navigation.DropdownMenu(c, "File", items, navigation.DropdownMenuOptions{}); r.Chosen >= 0 {
			chosen = r.Chosen
		}
	})
	click(t, tt, "File")
	click(t, tt, "Load")
	if shows(tt, "Quit") || !tt.Focused("File") {
		t.Fatal("a click on a disabled item left the menu or the focus astray")
	}
	// A click on the panel's edge, below the last row.
	click(t, tt, "File")
	q := find(t, tt, "Quit")
	tt.ClickAt(q.X+q.W/2, q.Y+q.H+3)
	if shows(tt, "Quit") || !tt.Focused("File") {
		t.Fatal("a click on the panel's edge left the menu or the focus astray")
	}
	tt.Key(0, ui.KeyDown)
	items[0].Disabled = true
	tt.Key(0, ui.KeyEnter)
	if chosen != -1 {
		t.Fatalf("Enter took an item disabled while lit: %d", chosen)
	}
}

func TestMenuBarMovesAcrossAndDown(t *testing.T) {
	m, item := -1, -1
	menus := []navigation.Menu{
		{Title: "File", Items: []navigation.MenuItem{{Label: "New"}, {Label: "Open"}}},
		{Title: "Edit", Items: []navigation.MenuItem{{Label: "Undo"}, {Label: "Redo"}}},
	}
	tt := view(500, 300, newClock(), false, func(c *ui.Context) {
		if r := navigation.MenuBar(c, "Menus", menus, navigation.MenuBarOptions{}); r.Item >= 0 {
			m, item = r.Menu, r.Item
		}
	})
	click(t, tt, "File")
	find(t, tt, "Open")
	tt.Key(0, ui.KeyRight)
	if shows(tt, "Open") || !shows(tt, "Redo") {
		t.Fatal("Right did not move to the Edit menu")
	}
	// The pointer rests on File and must not take over.
	tt.Frame()
	if !shows(tt, "Redo") {
		t.Fatal("a resting pointer took the bar back")
	}
	tt.Key(0, ui.KeyLeft)
	if !shows(tt, "Open") {
		t.Fatal("Left did not move back to File")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if m != 1 || item != 1 {
		t.Fatalf("Down, Enter chose menu %d item %d", m, item)
	}
	click(t, tt, "File")
	e := find(t, tt, "Edit")
	tt.Move(e.X+e.W/2, e.Y+e.H/2)
	if !shows(tt, "Undo") || shows(tt, "New") {
		t.Fatal("moving onto Edit did not open it")
	}
}

func TestContextMenuGivesUpFocusWhenOff(t *testing.T) {
	off := false
	items := []navigation.MenuItem{{Label: "Inspect"}}
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		input.Button(c, "before", input.ButtonOptions{})
		navigation.ContextMenu(c, "Area", items, navigation.ContextMenuOptions{Disabled: off}, func() {
			ui.Box(c).Size(100, 60)
		})
		ui.Box(c).Size(4, 4).Label("elsewhere")
	})
	click(t, tt, "before")
	tt.Key(0, ui.KeyTab)
	off = true
	tt.Frame()
	if tt.Focused("Area") && !tt.Focused("elsewhere") {
		t.Fatal("a disabled area kept the focus")
	}
	off = false
	tt.Frame()
	if !tt.Focused("Area") || tt.Focused("elsewhere") {
		t.Fatal("the focus did not come back")
	}
}

func TestContextMenuOpensAtThePointer(t *testing.T) {
	chosen := -1
	items := []navigation.MenuItem{{Label: "Inspect"}, {Label: "Throw"}}
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		if r := navigation.ContextMenu(c, "Area", items, navigation.ContextMenuOptions{}, func() {
			ui.Box(c).Size(200, 100).Label("Slime")
		}); r.Chosen >= 0 {
			chosen = r.Chosen
		}
	})
	tt.Key(0, ui.KeyTab)
	a := find(t, tt, "Area")
	tt.Move(a.X+120, a.Y+80)
	tt.Key(0, ui.KeyContextMenu)
	if it := find(t, tt, "Inspect"); it.X < a.X+110 || it.Y < a.Y+70 {
		t.Fatalf("the menu key opened at %v, away from the pointer", it)
	}
	tt.Key(0, ui.KeyEscape)
	s := find(t, tt, "Slime")
	if err := tt.RightClick("Slime"); err != nil {
		t.Fatal(err)
	}
	it := find(t, tt, "Inspect")
	if it.X < s.X+s.W/2-2 || it.Y < s.Y+s.H/2-2 {
		t.Fatalf("the menu opened at %v, away from the pointer", it)
	}
	click(t, tt, "Throw")
	if chosen != 1 {
		t.Fatalf("Throw gave %d", chosen)
	}
}

func TestCommandPaletteFiltersAndRuns(t *testing.T) {
	open, ran := false, -1
	cmds := []navigation.MenuItem{{Label: "Save game"}, {Label: "Load game"}, {Label: "Settings", Disabled: true}, {Label: "Quit"}}
	tt := view(500, 400, newClock(), false, func(c *ui.Context) {
		if r := navigation.CommandPalette(c, "Commands", cmds, &open, navigation.CommandPaletteOptions{}); r.Ran >= 0 {
			ran = r.Ran
		}
	})
	tt.Key(ui.Cmd, ui.KeyK)
	if !open {
		t.Fatal("Cmd+K did not open the palette")
	}
	tt.Type("game")
	if shows(tt, "Quit") || !shows(tt, "Load game") {
		t.Fatal("typing did not filter")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if ran != 1 || open {
		t.Fatalf("Down, Enter ran %d, open %v", ran, open)
	}
	tt.Key(ui.Cmd, ui.KeyK)
	tt.Type("set")
	tt.Key(0, ui.KeyEnter)
	if ran != 1 || !open {
		t.Fatal("a disabled command ran")
	}
	tt.Key(0, ui.KeyEscape)
	if open {
		t.Fatal("Escape did not close the palette")
	}
	inert := false
	tt = view(500, 400, newClock(), false, func(c *ui.Context) {
		s := pixel.Current(c)
		s.Inert = inert
		pixel.Scope(c, s, func() { navigation.CommandPalette(c, "Commands", cmds, &open, navigation.CommandPaletteOptions{}) })
	})
	tt.Key(ui.Cmd, ui.KeyK)
	inert = true
	tt.Frame()
	tt.Key(ui.Cmd, ui.KeyK)
	if open {
		t.Fatal("the palette stayed or opened under Inert")
	}
}
