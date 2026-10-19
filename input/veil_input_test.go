package input_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// veilCase is a control, its value, key and drag.
type veilCase struct {
	label string
	build func(c *ui.Context)
	value func() string
	key   ui.Key
	// dx, dy move the pointer; zeros skip the drag.
	dx, dy float32
}

func veilCases() map[string]func() veilCase {
	return map[string]func() veilCase{
		"Knob": func() veilCase {
			v := 50.0
			return veilCase{"Gain", func(c *ui.Context) { input.Knob(c, "Gain", &v, input.KnobOptions{Max: 100}) },
				func() string { return fmt.Sprint(v) }, ui.KeyUp, 0, -20}
		},
		"RangeSlider": func() veilCase {
			lo, hi := 20.0, 80.0
			return veilCase{"Band low", func(c *ui.Context) {
				input.RangeSlider(c, "Band", &lo, &hi, input.RangeSliderOptions{Max: 100})
			}, func() string { return fmt.Sprint(lo, hi) }, ui.KeyRight, 20, 0}
		},
		"SlotPicker": func() veilCase {
			v := 1
			return veilCase{"Reel", func(c *ui.Context) {
				input.SlotPicker(c, "Reel", &v, []string{"A", "B", "C", "D"}, input.SlotPickerOptions{})
			}, func() string { return fmt.Sprint(v) }, ui.KeyDown, 0, -40}
		},
		"NumberInput": func() veilCase {
			v := 5.0
			return veilCase{"Lives", func(c *ui.Context) { input.NumberInput(c, "Lives", &v, input.NumberInputOptions{Max: 9}) },
				func() string { return fmt.Sprint(v) }, ui.KeyUp, 0, 0}
		},
		"Joystick": func() veilCase {
			var r input.JoystickResult
			return veilCase{"Stick", func(c *ui.Context) { r = input.Joystick(c, "Stick", input.JoystickOptions{}) },
				func() string { return fmt.Sprintf("%.2f %.2f", r.X, r.Y) }, ui.KeyRight, 12, 0}
		},
	}
}

// underVeil builds tc where the test raises a veil.
func underVeil(tc veilCase, inert *bool) *ui.Tester {
	return view(400, 300, nil, func(c *ui.Context) {
		s := pixel.Current(c)
		s.Inert = *inert
		pixel.Scope(c, s, func() { tc.build(c) })
	})
}

func TestVeilsStopDragsAndKeys(t *testing.T) {
	for name, mk := range veilCases() {
		t.Run(name, func(t *testing.T) {
			if tc := mk(); tc.dx != 0 || tc.dy != 0 {
				inert := false
				tt := underVeil(tc, &inert)
				r := find(t, tt, tc.label)
				x, y := r.X+r.W/2, r.Y+r.H/2
				start := tc.value()
				tt.Press(x, y)
				tt.Move(x+tc.dx, y+tc.dy)
				moved := tc.value()
				if moved == start {
					t.Fatalf("the drag did nothing to %s", start)
				}
				inert = true
				tt.Move(x+2*tc.dx, y+2*tc.dy)
				// A stick let go is centered; values stay.
				if got := tc.value(); got != moved && got != "0.00 0.00" {
					t.Errorf("a drag under the veil moved %s to %s", moved, got)
				}
				tt.Release(x+2*tc.dx, y+2*tc.dy)
			}
			tc, inert := mk(), false
			tt := underVeil(tc, &inert)
			tt.Key(0, ui.KeyTab)
			before := tc.value()
			inert = true
			tt.Key(0, tc.key)
			if got := tc.value(); got != before {
				t.Errorf("a key in the veil's first frame moved %s to %s", before, got)
			}
		})
	}
}

func TestSlotPickerWheelStopsUnderAVeil(t *testing.T) {
	v, inert := 1, true
	tt := underVeil(veilCase{build: func(c *ui.Context) {
		input.SlotPicker(c, "Reel", &v, []string{"A", "B", "C", "D"}, input.SlotPickerOptions{})
	}}, &inert)
	r := find(t, tt, "Reel")
	tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 100)
	if v != 1 {
		t.Fatalf("the wheel turned a veiled reel to %d", v)
	}
	inert = false
	tt.Frame()
	tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 100)
	if v == 1 {
		t.Fatal("the wheel did not turn the reel")
	}
}

func TestVeilStopsHoldsAndPicks(t *testing.T) {
	// HoldButton: Space repeats, and the veil falls with the last.
	clk, inert, confirmed := &clock{time.Unix(3_000_000, 0)}, false, 0
	tt := view(400, 200, clk, func(c *ui.Context) {
		s := pixel.Current(c)
		s.Inert = inert
		pixel.Scope(c, s, func() {
			if input.HoldButton(c, "ERASE", input.HoldButtonOptions{}).Confirmed {
				confirmed++
			}
		})
	})
	tt.Key(0, ui.KeyTab)
	for i := range 3 {
		inert = i == 2
		tt.Key(0, ui.KeySpace)
		clk.t = clk.t.Add(500 * time.Millisecond)
	}
	tt.Frame()
	if confirmed != 0 {
		t.Fatal("a hold finished under the veil")
	}
	// Autocomplete and Combobox: Tab or Enter as the veil falls.
	for name, build := range map[string]func(c *ui.Context, v *string){
		"Autocomplete": func(c *ui.Context, v *string) {
			input.Autocomplete(c, "Hero", v, []string{"Quux"}, input.AutocompleteOptions{})
		},
		"Combobox": func(c *ui.Context, v *string) {
			input.Combobox(c, "Hero", v, []string{"Quux"}, input.ComboboxOptions{})
		},
	} {
		v, inert := "", false
		tt := view(400, 300, nil, func(c *ui.Context) {
			s := pixel.Current(c)
			s.Inert = inert
			pixel.Scope(c, s, func() { build(c, &v) })
		})
		tt.Key(0, ui.KeyTab)
		tt.Type("Qu")
		inert = true
		tt.Key(0, map[string]ui.Key{"Autocomplete": ui.KeyTab, "Combobox": ui.KeyEnter}[name])
		if v != "Qu" {
			t.Errorf("%s: a key as the veil fell made %q", name, v)
		}
	}
}

func TestVeilStopsTyping(t *testing.T) {
	for name, tc := range map[string]struct {
		build func(c *ui.Context, v *string, tags *[]string)
		// last is the key or text coming with the veil.
		last string
	}{
		// The editor turns the typed letter away.
		"TextInput": {func(c *ui.Context, v *string, _ *[]string) { input.TextInput(c, "Name", v, input.TextInputOptions{}) }, "Z"},
		"TokenField": {func(c *ui.Context, v *string, tags *[]string) {
			input.TokenField(c, "Tags", tags, input.TokenFieldOptions{})
		}, "enter"},
		"TokenField comma": {func(c *ui.Context, v *string, tags *[]string) {
			input.TokenField(c, "Tags", tags, input.TokenFieldOptions{})
		}, ","},
		"EditableText": {func(c *ui.Context, v *string, _ *[]string) {
			input.EditableText(c, "Hero", v, input.EditableTextOptions{})
		}, "enter"},
	} {
		v, tags, inert := "before", []string{}, false
		tt := view(400, 200, nil, func(c *ui.Context) {
			s := pixel.Current(c)
			s.Inert = inert
			pixel.Scope(c, s, func() { tc.build(c, &v, &tags) })
		})
		tt.Key(0, ui.KeyTab)
		if name == "EditableText" {
			tt.Key(0, ui.KeyEnter)
		}
		tt.Type("ab")
		inert = true
		if tc.last == "enter" {
			tt.Key(0, ui.KeyEnter)
		} else {
			tt.Type(tc.last)
		}
		if want := map[string]string{"TextInput": "beforeab", "EditableText": "before"}[name]; want != "" && v != want || len(tags) != 0 {
			t.Errorf("%s: a key as the veil fell left %q, tags %q", name, v, tags)
		}
	}
}

func TestRadioGroupSpaceAsAVeilFallsTurnsNothing(t *testing.T) {
	sel, inert := 0, false
	tt := underVeil(veilCase{build: func(c *ui.Context) {
		input.RadioGroup(c, []string{"A", "B", "C"}, &sel, input.RadioGroupOptions{})
	}}, &inert)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	// Focus stays on B as the caller sets A again.
	sel = 0
	tt.Frame()
	inert = true
	tt.Key(0, ui.KeySpace)
	if sel != 0 {
		t.Fatalf("Space as the veil fell chose %d", sel)
	}
}

func TestEditorKeepsTheCallersValueAsAVeilFalls(t *testing.T) {
	v, inert := "before", false
	tt := underVeil(veilCase{build: func(c *ui.Context) { input.TextInput(c, "Text", &v, input.TextInputOptions{}) }}, &inert)
	tt.Key(0, ui.KeyTab)
	// The caller writes as the veil falls on a letter.
	v, inert = "server", true
	tt.Type("x")
	if v != "server" {
		t.Fatalf("the caller's value became %q", v)
	}
}

func TestHandlersActNotAsAVeilFalls(t *testing.T) {
	// PinInput: a digit typed as the veil or Disabled falls.
	for _, kind := range []string{"inert", "disabled"} {
		v, inert, off, complete := "123", false, false, 0
		tt := underVeil(veilCase{build: func(c *ui.Context) {
			if input.PinInput(c, "PIN", &v, input.PinInputOptions{Disabled: off}).Complete {
				complete++
			}
		}}, &inert)
		tt.Key(0, ui.KeyTab)
		inert, off = kind == "inert", kind == "disabled"
		tt.Type("4")
		if v != "123" || complete != 0 {
			t.Errorf("PinInput %s: %q, complete %d", kind, v, complete)
		}
	}
	// CheatCode: the code's last key comes with the veil.
	inert, wins := false, 0
	tt := underVeil(veilCase{build: func(c *ui.Context) {
		if input.CheatCode(c, "Code", input.CheatCodeOptions{Sequence: []ui.Key{ui.KeyB, ui.KeyA}}).Unlocked {
			wins++
		}
	}}, &inert)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyB)
	inert = true
	tt.Key(0, ui.KeyA)
	if wins != 0 {
		t.Errorf("the code unlocked under the veil %d times", wins)
	}
	// SlotPicker: a wheel notch comes with the veil.
	slot := 0
	inert = false
	tt = underVeil(veilCase{build: func(c *ui.Context) {
		input.SlotPicker(c, "Slots", &slot, []string{"A", "B", "C"}, input.SlotPickerOptions{})
	}}, &inert)
	r := find(t, tt, "Slots")
	inert = true
	tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 40)
	if slot != 0 {
		t.Errorf("the wheel turned a veiled reel to %d", slot)
	}
}

func TestComboboxPicksNothingAsAVeilFalls(t *testing.T) {
	for _, kind := range []string{"inert", "disabled"} {
		inert, off, v, chosen := false, false, "", 0
		tt := underVeil(veilCase{build: func(c *ui.Context) {
			if input.Combobox(c, "Combo", &v, []string{"Alpha", "Beta"}, input.ComboboxOptions{Disabled: off}).Chosen {
				chosen++
			}
		}}, &inert)
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyDown)
		r := find(t, tt, "Beta")
		// The press comes before the veil, the release after.
		tt.Press(r.X+r.W/2, r.Y+r.H/2)
		inert, off = kind == "inert", kind == "disabled"
		tt.Release(r.X+r.W/2, r.Y+r.H/2)
		if v != "" || chosen != 0 {
			t.Errorf("%s: picked %q, %d times", kind, v, chosen)
		}
	}
}

func TestStepperHoldNeverRunsBackward(t *testing.T) {
	clk, v := &clock{time.Unix(1_700_000_000, 0)}, 50
	tt := view(400, 100, clk, func(c *ui.Context) { input.Stepper(c, "Step", &v, input.StepperOptions{Max: 100}) })
	r := find(t, tt, "Step")
	tt.Press(r.X+r.W+12, r.Y+r.H/2)
	clk.t = clk.t.Add(time.Second)
	tt.Frame()
	before := v
	clk.t = clk.t.Add(-time.Hour)
	tt.Frame()
	if before <= 51 || v < before {
		t.Fatalf("a held + went %d, then %d", before, v)
	}
}
