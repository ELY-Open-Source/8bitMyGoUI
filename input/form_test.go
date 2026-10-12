package input_test

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestEditableTextEnterKeepsEscapeReverts(t *testing.T) {
	hero, changes := "Hero", 0
	tt := view(400, 200, nil, func(c *ui.Context) {
		input.Button(c, "before", input.ButtonOptions{})
		if input.EditableText(c, "Name", &hero, input.EditableTextOptions{}).Changed {
			changes++
		}
	})
	click(t, tt, "before")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if _, ok := tt.TextCaret(); !ok {
		t.Fatal("Enter did not open a field")
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("Zed")
	tt.Key(0, ui.KeyEnter)
	if hero != "Zed" || changes != 1 || !tt.Focused("Name") {
		t.Fatalf("Enter kept %q after %d changes, focus back %v", hero, changes, tt.Focused("Name"))
	}
	if _, ok := tt.TextCaret(); ok {
		t.Fatal("the field stayed open")
	}
	r := find(t, tt, "Name")
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	if _, ok := tt.TextCaret(); !ok {
		t.Fatal("a double click did not open the field")
	}
	tt.Type("ox")
	tt.Key(0, ui.KeyEscape)
	if hero != "Zed" || changes != 1 {
		t.Fatalf("Escape kept %q", hero)
	}
	if _, ok := tt.TextCaret(); ok {
		t.Fatal("a double click did not open the field, or Escape left it open")
	}
}

func TestFieldAndFieldsetLayOutTheirNotes(t *testing.T) {
	name := ""
	tt := view(400, 300, nil, func(c *ui.Context) {
		input.Fieldset(c, "Hero", input.FieldsetOptions{}, func() {
			input.Field(c, "Call sign", input.FieldOptions{Description: "Shown on the board.", Error: "Taken!"}, func() {
				input.TextInput(c, "sign", &name, input.TextInputOptions{})
			})
		})
	})
	label, field, note, bad := find(t, tt, "Call sign"), find(t, tt, "sign"), find(t, tt, "Shown on the board."), find(t, tt, "Taken!")
	if !(label.Y < field.Y && field.Y < note.Y && note.Y < bad.Y) {
		t.Fatal("not label, control, description, error from top down")
	}
	if !has(slots(t, tt.Image(), ui.Rect{X: bad.X - 30, Y: bad.Y, W: 30, H: bad.H}), pixel.Bad) {
		t.Fatal("the error has no Bad sign")
	}
	set := find(t, tt, "Hero")
	if !has(slots(t, tt.Image(), ui.Rect{X: set.X, Y: set.Y, W: set.W, H: 2 * pixel.LineHeight}), pixel.Accent) || label.X <= set.X {
		t.Fatal("the fieldset has no legend, or its fields sit on its edge")
	}
}

func TestColorSlotPicks(t *testing.T) {
	slot, changes := pixel.Slot(3), 0
	tt := view(400, 200, nil, func(c *ui.Context) {
		if input.ColorSlot(c, "Paint", &slot, input.ColorSlotOptions{}).Changed {
			changes++
		}
	})
	at(t, tt, "Paint", 6, 12+6)
	if slot != pixel.Slot(8) {
		t.Fatalf("a click on the ninth cell gave %v", slot)
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyRight)
	if slot != pixel.Slot(1) || changes != 3 {
		t.Fatalf("Up, Right gave %v after %d changes", slot, changes)
	}
}

func TestJoystickDragsNudgesAndSprings(t *testing.T) {
	clk := &clock{time.Unix(3_000_000, 0)}
	var x, y float64
	tt := view(400, 200, clk, func(c *ui.Context) {
		r := input.Joystick(c, "Stick", input.JoystickOptions{})
		x, y = r.X, r.Y
	})
	r := find(t, tt, "Stick")
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	tt.Press(cx, cy)
	tt.Move(cx+60, cy)
	if math.Abs(x-1) > 1e-9 || y != 0 {
		t.Fatalf("dragging far right gave %v, %v", x, y)
	}
	tt.Move(cx, cy-12)
	if math.Abs(y+0.5) > 1e-9 {
		t.Fatalf("dragging up half way gave y %v", y)
	}
	tt.Release(cx, cy-12)
	if x != 0 || y != 0 {
		t.Fatal("letting go did not spring back")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	if y != 0.5 {
		t.Fatalf("two nudges gave y %v", y)
	}
	clk.t = clk.t.Add(time.Second)
	tt.Frame()
	if y != 0 {
		t.Fatal("a nudge did not spring back")
	}
}

func TestCheatCodeUnlocks(t *testing.T) {
	unlocked := 0
	tt := view(400, 200, nil, func(c *ui.Context) {
		if input.CheatCode(c, "Cheat", input.CheatCodeOptions{}).Unlocked {
			unlocked++
		}
	})
	click(t, tt, "Cheat")
	r := find(t, tt, "Cheat")
	if has(slots(t, tt.Image(), r), pixel.Accent) {
		t.Fatal("an empty trail is lit")
	}
	for _, k := range []ui.Key{ui.KeyUp, ui.KeyUp, ui.KeyUp, ui.KeyDown, ui.KeyDown, ui.KeyLeft, ui.KeyRight, ui.KeyLeft, ui.KeyRight, ui.KeyB} {
		tt.Key(0, k)
	}
	if unlocked != 0 || !has(slots(t, tt.Image(), r), pixel.Accent) {
		t.Fatal("unlocked early, or the trail is not lit")
	}
	tt.Key(0, ui.KeyA)
	if unlocked != 1 || !has(slots(t, tt.Image(), r), pixel.Good) {
		t.Fatal("the code did not unlock")
	}
}

func TestCheatCodeBreaksOnOtherKeys(t *testing.T) {
	unlocked := 0
	tt := view(400, 200, nil, func(c *ui.Context) {
		if input.CheatCode(c, "Cheat", input.CheatCodeOptions{Sequence: []ui.Key{ui.KeyA, ui.KeyB}}).Unlocked {
			unlocked++
		}
	})
	click(t, tt, "Cheat")
	r := find(t, tt, "Cheat")
	tt.Key(0, ui.KeyA)
	tt.Key(0, ui.KeyX)
	if !has(slots(t, tt.Image(), r), pixel.Faint) {
		t.Fatal("the wrong key is not in the trail")
	}
	tt.Key(0, ui.KeyB)
	if unlocked != 0 {
		t.Fatal("A, X, B unlocked A, B")
	}
	tt.Key(0, ui.KeyA)
	tt.Key(0, ui.KeyB)
	if unlocked != 1 {
		t.Fatal("A, B did not unlock")
	}
}

func TestSlotPickerSpins(t *testing.T) {
	reel, changes := 0, 0
	tt := view(400, 200, nil, func(c *ui.Context) {
		if input.SlotPicker(c, "Reel", &reel, []string{"Cherry", "Bell", "Seven"}, input.SlotPickerOptions{}).Changed {
			changes++
		}
	})
	click(t, tt, "Reel")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyUp)
	if reel != 2 || changes != 3 {
		t.Fatalf("Down, Up, Up gave %d after %d changes", reel, changes)
	}
	r := find(t, tt, "Reel")
	tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 2*2*pixel.LineHeight)
	if reel != 1 {
		t.Fatalf("two rows of wheel gave %d", reel)
	}
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(r.X+r.W/2, r.Y+r.H/2+2*pixel.LineHeight+4)
	tt.Release(r.X+r.W/2, r.Y+r.H/2+2*pixel.LineHeight+4)
	if reel != 0 {
		t.Fatalf("dragging a row down gave %d", reel)
	}
}

func TestGridCopyFollowsTheLocale(t *testing.T) {
	render := func(l pixel.Locale) (end, banner []pixel.Slot) {
		name := ""
		tt := ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Locale: l})
			input.NameEntry(c, "Initials", &name, input.NameEntryOptions{})
			input.CheatCode(c, "Cheat", input.CheatCodeOptions{Sequence: []ui.Key{ui.KeyA}})
		}, 400, 300)
		click(t, tt, "Cheat")
		tt.Key(0, ui.KeyA)
		g, ch := find(t, tt, "Initials"), find(t, tt, "Cheat")
		end = slots(t, tt.Image(), ui.Rect{X: g.X + g.W - 28, Y: g.Y + g.H - 32, W: 28, H: 32})
		banner = slots(t, tt.Image(), ui.Rect{X: ch.X, Y: ch.Y + 28, W: ch.W, H: ch.H - 28})
		return end, banner
	}
	enEnd, enBanner := render(pixel.En)
	zhEnd, zhBanner := render(pixel.Zh)
	if equal(enEnd, zhEnd) || equal(enBanner, zhBanner) {
		t.Fatal("END or the unlock banner ignores the locale")
	}
}

func TestSwitchLabelNamesABarePlate(t *testing.T) {
	lever := false
	tt := view(300, 100, nil, func(c *ui.Context) {
		input.Switch(c, "", &lever, input.SwitchOptions{Label: "Fly"})
	})
	click(t, tt, "Fly")
	if !lever || len(tt.Texts()) != 1 {
		t.Fatalf("switch %v, texts %q", lever, tt.Texts())
	}
}
