package input_test

import (
	"image"
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func view(w, h int, clk *clock, build func(c *ui.Context)) *ui.Tester {
	if clk == nil {
		clk = &clock{time.Unix(3_000_000, 0)}
	}
	return ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: clk.now})
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

// ink is where a box holds Ink, ignoring its background.
func ink(t *testing.T, img *image.RGBA, r ui.Rect) []bool {
	t.Helper()
	var out []bool
	for _, s := range slots(t, img, r) {
		out = append(out, s == pixel.Ink)
	}
	return out
}

func has(s []pixel.Slot, want pixel.Slot) bool {
	for _, v := range s {
		if v == want {
			return true
		}
	}
	return false
}

func TestButtonClicksKeysLoadingAndSizes(t *testing.T) {
	clicks, opts := 0, input.ButtonOptions{Variant: input.Primary}
	tt := view(400, 200, nil, func(c *ui.Context) {
		if input.Button(c, "GO", opts).Clicked {
			clicks++
		}
		input.Button(c, "S", input.ButtonOptions{Size: input.Small})
		input.Button(c, "L", input.ButtonOptions{Size: input.Large})
	})
	click(t, tt, "GO")
	tt.Key(0, ui.KeyEnter)
	if clicks != 2 {
		t.Fatalf("click and Enter gave %d clicks", clicks)
	}
	if !has(slots(t, tt.Image(), find(t, tt, "GO")), pixel.Accent) {
		t.Fatal("a primary button has no Accent face")
	}
	if s, l := find(t, tt, "S"), find(t, tt, "L"); s.H != 2*pixel.LineHeight || l.H != 2*(2*pixel.LineHeight+4) {
		t.Fatalf("small %v, large %v", s.H, l.H)
	}
	r := find(t, tt, "GO")
	top := func() pixel.Slot { return slots(t, tt.Image(), ui.Rect{X: r.X + r.W/2, Y: r.Y, W: 2, H: 2})[0] }
	up := top()
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	if top() == up {
		t.Fatal("pressing did not sink the face")
	}
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	opts.Loading = true
	tt.Frame()
	click(t, tt, "GO")
	opts.Loading, opts.Disabled = false, true
	tt.Frame()
	click(t, tt, "GO")
	if clicks != 3 {
		t.Fatalf("loading or disabled clicked: %d clicks", clicks)
	}
}

func TestIconButtonNeedsAName(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("a nameless icon button did not panic")
		}
	}()
	view(100, 100, nil, func(c *ui.Context) { input.IconButton(c, icon.Star, "", input.IconButtonOptions{}) })
}

func TestButtonGroupIsOneTabStop(t *testing.T) {
	got := -1
	tt := view(400, 100, nil, func(c *ui.Context) {
		if r := input.ButtonGroup(c, []string{"A", "B", "C"}, input.ButtonGroupOptions{}); r.Clicked >= 0 {
			got = r.Clicked
		}
		input.Button(c, "after", input.ButtonOptions{})
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if got != 1 {
		t.Fatalf("Tab, Right, Enter clicked %d", got)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("after") {
		t.Fatal("Tab did not leave the group")
	}
}

func TestTogglesChecksAndSwitch(t *testing.T) {
	var on, checked, lever bool
	modes := []bool{true, false, false}
	tt := view(500, 200, nil, func(c *ui.Context) {
		input.Toggle(c, "MUSIC", &on, input.ToggleOptions{})
		input.Checkbox(c, "Remember", &checked, input.CheckboxOptions{})
		input.Switch(c, "Power", &lever, input.SwitchOptions{})
		input.ToggleGroup(c, []string{"E", "N", "H"}, modes, input.ToggleGroupOptions{Single: true})
	})
	click(t, tt, "MUSIC")
	click(t, tt, "Remember")
	tt.Key(0, ui.KeySpace)
	click(t, tt, "Power")
	click(t, tt, "H")
	if !on || checked || !lever || modes[0] || !modes[2] {
		t.Fatalf("toggle %v, checkbox %v after Space, switch %v, modes %v", on, checked, lever, modes)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("a short state slice did not panic")
		}
	}()
	view(100, 100, nil, func(c *ui.Context) {
		input.ToggleGroup(c, []string{"a", "b"}, []bool{true}, input.ToggleGroupOptions{})
	})
}

func TestCheckboxGroupParent(t *testing.T) {
	checked := []bool{true, false}
	tt := view(300, 200, nil, func(c *ui.Context) {
		input.CheckboxGroup(c, "Party", []string{"Knight", "Mage"}, checked, input.CheckboxGroupOptions{})
	})
	if !has(slots(t, tt.Image(), find(t, tt, "Party")), pixel.Accent) {
		t.Fatal("a half-checked group shows no mixed mark")
	}
	click(t, tt, "Party")
	if !checked[0] || !checked[1] {
		t.Fatalf("parent on a mixed group checked %v", checked)
	}
	click(t, tt, "Party")
	if checked[0] || checked[1] {
		t.Fatalf("parent on a full group left %v", checked)
	}
}

func TestRadiosChoose(t *testing.T) {
	class, season := "knight", 0
	tt := view(400, 200, nil, func(c *ui.Context) {
		input.Radio(c, "knight", "knight", &class, input.RadioOptions{})
		input.Radio(c, "mage", "mage", &class, input.RadioOptions{})
		input.Button(c, "before", input.ButtonOptions{})
		input.RadioGroup(c, []string{"Spring", "Summer", "Autumn"}, &season, input.RadioGroupOptions{})
	})
	click(t, tt, "mage")
	click(t, tt, "mage")
	if class != "mage" {
		t.Fatalf("radio left %q", class)
	}
	click(t, tt, "Autumn")
	if season != 2 {
		t.Fatalf("click chose %d", season)
	}
	click(t, tt, "before")
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Autumn") || season != 2 {
		t.Fatal("Tab entered elsewhere than the chosen radio, or chose")
	}
	tt.Key(0, ui.KeyUp)
	if season != 1 {
		t.Fatalf("Up chose %d, want 1", season)
	}
}

func TestSlidersDragKeysAndSteps(t *testing.T) {
	v, s, low, high, f := 0.0, 0.0, 20.0, 60.0, 0.0
	disabled := false
	tt := view(500, 300, nil, func(c *ui.Context) {
		ui.Column(c).Gap(10).Children(func() {
			input.Slider(c, "vol", &v, input.SliderOptions{Width: 100, Disabled: disabled})
			input.StepSlider(c, "speed", &s, input.StepSliderOptions{Min: 0, Max: 100, Step: 20, Width: 100})
			input.RangeSlider(c, "lv", &low, &high, input.RangeSliderOptions{Min: 0, Max: 100, Width: 100})
			input.Fader(c, "master", &f, input.FaderOptions{Height: 50})
		})
	})
	r := find(t, tt, "vol")
	tt.ClickAt(r.X+r.W*0.75, r.Y+r.H/2)
	if v < 0.7 || v > 0.8 {
		t.Fatalf("a click three quarters along set %v", v)
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyRight)
	if v != 0.05 {
		t.Fatalf("Home then Right gave %v, want 0.05", v)
	}
	sr := find(t, tt, "speed")
	tt.ClickAt(sr.X+100*2*0.33, sr.Y+sr.H/2)
	if s != 40 {
		t.Fatalf("a third along snapped to %v, want 40", s)
	}
	tt.Key(0, ui.KeyRight)
	if _, shown := tt.Find("60"); s != 60 || !shown {
		t.Fatalf("Right stepped to %v, shown %v, want 60", s, shown)
	}
	lowKnob := find(t, tt, "lv low")
	tt.Press(lowKnob.X+lowKnob.W/2, lowKnob.Y+lowKnob.H/2)
	tt.Move(lowKnob.X+190, lowKnob.Y+lowKnob.H/2)
	tt.Release(lowKnob.X+190, lowKnob.Y+lowKnob.H/2)
	if low != high {
		t.Fatalf("dragging low past high left %v and %v", low, high)
	}
	fr := find(t, tt, "master")
	tt.ClickAt(fr.X+fr.W/2, fr.Y+2)
	if f < 0.9 {
		t.Fatalf("a click at the fader's top set %v", f)
	}
	disabled, v = true, 0.5
	tt.Frame()
	tt.ClickAt(r.X+2, r.Y+r.H/2)
	if v != 0.5 {
		t.Fatal("a disabled slider moved")
	}
}

func TestStepperRepeatsWhileHeld(t *testing.T) {
	clk := &clock{time.Unix(3_000_000, 0)}
	n := 5
	tt := view(300, 100, clk, func(c *ui.Context) {
		input.Stepper(c, "Lives", &n, input.StepperOptions{Min: 0, Max: 20})
	})
	r := find(t, tt, "Lives")
	right := r.X + r.W + 4 + 12
	tt.Press(right, r.Y+r.H/2)
	if n != 6 {
		t.Fatalf("a press stepped to %d", n)
	}
	clk.t = clk.t.Add(300 * time.Millisecond)
	tt.Frame()
	if n != 6 {
		t.Fatalf("repeat began early: %d", n)
	}
	clk.t = clk.t.Add(400 * time.Millisecond)
	tt.Frame()
	tt.Release(right, r.Y+r.H/2)
	if n != 9 {
		t.Fatalf("700 ms held gave %d, want 9", n)
	}
	click(t, tt, "Lives")
	tt.Key(0, ui.KeyDown)
	if n != 8 {
		t.Fatalf("Down stepped to %d", n)
	}
	n = 30
	tt.Key(0, ui.KeyUp)
	if n != 20 {
		t.Fatalf("Up past the top gave %d", n)
	}
}

func TestNumberInputParsesAndSteps(t *testing.T) {
	v := 5.0
	tt := view(300, 100, nil, func(c *ui.Context) {
		input.NumberInput(c, "Price", &v, input.NumberInputOptions{Min: 0, Max: 10})
	})
	click(t, tt, "Price")
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("7.5")
	if v != 7.5 {
		t.Fatalf("typing 7.5 left %v", v)
	}
	border := func() []pixel.Slot {
		r := find(t, tt, "Price")
		return slots(t, tt.Image(), ui.Rect{X: r.X, Y: r.Y - 6, W: r.W, H: 2})
	}
	tt.Type("x")
	if v != 7.5 || !has(border(), pixel.Bad) {
		t.Fatalf("bad text changed the value to %v, or the edge stayed: %v", v, border()[:3])
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("9")
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyUp)
	if v != 10 {
		t.Fatalf("Up twice from 9 gave %v, want the max 10", v)
	}
}

func TestFieldsTypeSubmitClear(t *testing.T) {
	name, query, secret := "", "slime", "abc"
	submits := 0
	tt := view(500, 300, nil, func(c *ui.Context) {
		if input.TextInput(c, "Name", &name, input.TextInputOptions{MaxLength: 3}).Submitted {
			submits++
		}
		input.SearchField(c, "Find", &query, input.SearchFieldOptions{})
		input.PasswordInput(c, "Secret", &secret, input.PasswordInputOptions{})
	})
	click(t, tt, "Name")
	tt.Type("Link!")
	tt.Key(0, ui.KeyEnter)
	if name != "Lin" || submits != 1 {
		t.Fatalf("name %q, submits %d", name, submits)
	}
	click(t, tt, "Find")
	tt.Key(0, ui.KeyEscape)
	if query != "" {
		t.Fatalf("Escape left %q", query)
	}
	secretInk := func() []bool {
		r := find(t, tt, "Secret")
		return ink(t, tt.Image(), ui.Rect{X: r.X, Y: r.Y, W: r.W, H: r.H})
	}
	masked := secretInk()
	click(t, tt, "Show")
	if equal(secretInk(), masked) {
		t.Fatal("the eye did not reveal the password in the same frame")
	}
	click(t, tt, "Find")
	query = "potion"
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	r := find(t, tt, "Find")
	shown := ink(t, tt.Image(), r)
	tt.Frame()
	if query != "" || !equal(shown, ink(t, tt.Image(), r)) {
		t.Fatalf("Escape left %q, or showed the old text a frame", query)
	}
}

func TestTextAreaFollowsTheCaret(t *testing.T) {
	notes := ""
	tt := view(400, 300, nil, func(c *ui.Context) {
		input.TextArea(c, "Notes", &notes, input.TextAreaOptions{Lines: 2, MaxLength: 50})
		pixel.Text(c, "F", pixel.TextOptions{})
	})
	click(t, tt, "Notes")
	tt.Type("A\nB\nC\nD\nE\nF")
	if _, ok := tt.Find("11/50"); !ok {
		t.Fatalf("no character count in %q", tt.Texts())
	}
	area, ref := find(t, tt, "Notes"), find(t, tt, "F")
	img := tt.Image()
	glyph := ink(t, img, ui.Rect{X: ref.X, Y: ref.Y, W: ref.W, H: ref.H})
	last := ink(t, img, ui.Rect{X: area.X, Y: area.Y + area.H - ref.H, W: ref.W, H: ref.H})
	if !equal(glyph, last) {
		t.Fatal("the caret's line F is not the bottom line in view")
	}
}

func equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestStepSliderStaysOnStepsWithinMax(t *testing.T) {
	s := 0.0
	tt := view(300, 100, nil, func(c *ui.Context) {
		input.StepSlider(c, "odd", &s, input.StepSliderOptions{Min: 0, Max: 10, Step: 6})
	})
	click(t, tt, "odd")
	tt.Key(0, ui.KeyEnd)
	if s != 6 {
		t.Fatalf("End went to %v, want the last step 6", s)
	}
	f := 0.0
	tt = view(300, 100, nil, func(c *ui.Context) {
		input.StepSlider(c, "tenths", &f, input.StepSliderOptions{Min: 0, Max: 0.3, Step: 0.1})
	})
	click(t, tt, "tenths")
	tt.Key(0, ui.KeyEnd)
	if math.Abs(f-0.3) > 1e-9 {
		t.Fatalf("End over tenths went to %v, want 0.3", f)
	}
}

func TestNumberInputRejectsNaNAndSyncsAtOnce(t *testing.T) {
	v, changes := 5.0, 0
	tt := view(300, 100, nil, func(c *ui.Context) {
		if input.NumberInput(c, "N", &v, input.NumberInputOptions{Min: 0, Max: 10}).Changed {
			changes++
		}
	})
	r := find(t, tt, "N")
	border := func() []pixel.Slot { return slots(t, tt.Image(), ui.Rect{X: r.X, Y: r.Y - 6, W: r.W, H: 2}) }
	field := func() []bool { return ink(t, tt.Image(), r) }
	eight := 8.0
	ref := view(300, 100, nil, func(c *ui.Context) {
		input.NumberInput(c, "N", &eight, input.NumberInputOptions{Min: 0, Max: 10})
	})
	shows8 := func() bool { return equal(field(), ink(t, ref.Image(), r)) }
	v = 8
	tt.Frame()
	if !shows8() {
		t.Fatal("an outside change did not show at once")
	}
	click(t, tt, "N")
	v = 3
	tt.Frame()
	eight = 3
	ref.Frame()
	click(t, ref, "N")
	if !shows8() {
		t.Fatal("a focused field ignored an outside change")
	}
	click(t, tt, "N")
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("NaN")
	if v != 3 || !has(border(), pixel.Bad) {
		t.Fatalf("NaN set %v or left the edge plain", v)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("10")
	changes = 0
	tt.Key(0, ui.KeyUp)
	if v != 10 || changes != 0 {
		t.Fatalf("Up at the max gave %v with %d changes", v, changes)
	}
}

func TestDisabledTextAreaHoldsStill(t *testing.T) {
	notes, disabled := "A\nB\nC\nD\nE\nF", false
	tt := view(400, 300, nil, func(c *ui.Context) {
		input.TextArea(c, "Notes", &notes, input.TextAreaOptions{Lines: 2, Disabled: disabled})
	})
	disabled = true
	tt.Frame()
	r := find(t, tt, "Notes")
	before := ink(t, tt.Image(), r)
	// The caret's line keeps it at the bottom; wheel up.
	tt.Scroll(r.X+10, r.Y+10, 0, -64)
	if !equal(before, ink(t, tt.Image(), find(t, tt, "Notes"))) {
		t.Fatal("a disabled text area scrolled")
	}
}
