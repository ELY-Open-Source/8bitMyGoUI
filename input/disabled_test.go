package input_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Each control, disabled after taking focus, ignores every input.
func TestDisabledHoldButtonIgnoresAHeldPress(t *testing.T) {
	clk := &clock{time.Unix(3_000_000, 0)}
	confirmed := 0
	tt := view(400, 200, clk, func(c *ui.Context) {
		if input.HoldButton(c, "HOLD", input.HoldButtonOptions{Disabled: true}).Confirmed {
			confirmed++
		}
	})
	r := find(t, tt, "HOLD")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	clk.t = clk.t.Add(2 * time.Second)
	tt.Frame()
	if confirmed != 0 {
		t.Fatal("a disabled hold button confirmed")
	}
}

func TestHoldButtonForgetsAHoldWhenDisabled(t *testing.T) {
	clk := &clock{time.Unix(3_000_000, 0)}
	confirmed, off := 0, false
	tt := view(400, 200, clk, func(c *ui.Context) {
		input.Button(c, "before", input.ButtonOptions{})
		if input.HoldButton(c, "HOLD", input.HoldButtonOptions{Disabled: off}).Confirmed {
			confirmed++
		}
	})
	hold := func(d time.Duration) {
		for end := clk.t.Add(d); clk.t.Before(end); clk.t = clk.t.Add(50 * time.Millisecond) {
			tt.Key(0, ui.KeySpace)
		}
	}
	click(t, tt, "before")
	tt.Key(0, ui.KeyTab)
	hold(500 * time.Millisecond)
	off = true
	clk.t = clk.t.Add(100 * time.Millisecond)
	tt.Frame()
	off = false
	tt.Frame()
	click(t, tt, "before")
	tt.Key(0, ui.KeyTab)
	hold(600 * time.Millisecond)
	if confirmed != 0 {
		t.Fatal("a hold from before disabling counted")
	}
}

func TestDisabledWhileFocusedTakesNoInput(t *testing.T) {
	type state struct {
		s, list   string
		n, events int
		f         float64
		tags      []string
		slot      pixel.Slot
	}
	cases := []struct {
		name, label string
		build       func(c *ui.Context, v *state, off bool)
	}{
		{"Select", "Class", func(c *ui.Context, v *state, off bool) {
			input.Select(c, "Class", &v.s, []string{"A", "B"}, input.SelectOptions{Disabled: off})
		}},
		{"MultiSelect", "Gear", func(c *ui.Context, v *state, off bool) {
			input.MultiSelect(c, "Gear", &v.tags, []string{"A", "B"}, input.MultiSelectOptions{Disabled: off})
		}},
		{"Combobox", "Fruit", func(c *ui.Context, v *state, off bool) {
			input.Combobox(c, "Fruit", &v.s, []string{"Apple"}, input.ComboboxOptions{Disabled: off})
		}},
		{"Autocomplete", "Spell", func(c *ui.Context, v *state, off bool) {
			input.Autocomplete(c, "Spell", &v.s, []string{"Fire"}, input.AutocompleteOptions{Disabled: off})
		}},
		{"TokenField", "Tags", func(c *ui.Context, v *state, off bool) {
			input.TokenField(c, "Tags", &v.tags, input.TokenFieldOptions{Disabled: off})
		}},
		{"Segmented", "B", func(c *ui.Context, v *state, off bool) {
			input.Segmented(c, "Mode", []string{"A", "B", "C"}, &v.n, input.SegmentedOptions{Disabled: off})
		}},
		{"Rating", "2/5", func(c *ui.Context, v *state, off bool) {
			input.Rating(c, "Stars", &v.n, input.RatingOptions{Disabled: off})
		}},
		{"Knob", "Gain", func(c *ui.Context, v *state, off bool) {
			input.Knob(c, "Gain", &v.f, input.KnobOptions{Disabled: off})
		}},
		{"PinInput", "PIN", func(c *ui.Context, v *state, off bool) {
			input.PinInput(c, "PIN", &v.s, input.PinInputOptions{Disabled: off})
		}},
		{"NameEntry", "Initials", func(c *ui.Context, v *state, off bool) {
			input.NameEntry(c, "Initials", &v.s, input.NameEntryOptions{Disabled: off})
		}},
		{"VirtualKeyboard", "Keys", func(c *ui.Context, v *state, off bool) {
			input.VirtualKeyboard(c, "Keys", &v.s, input.VirtualKeyboardOptions{Disabled: off})
		}},
		{"DPad", "Pad", func(c *ui.Context, v *state, off bool) {
			if input.DPad(c, "Pad", input.DPadOptions{Disabled: off}).Pressed != input.None {
				v.events++
			}
		}},
		{"HoldButton", "HOLD", func(c *ui.Context, v *state, off bool) {
			if input.HoldButton(c, "HOLD", input.HoldButtonOptions{Duration: time.Millisecond, Disabled: off}).Confirmed {
				v.events++
			}
		}},
		{"EditableText", "Name", func(c *ui.Context, v *state, off bool) {
			input.EditableText(c, "Name", &v.s, input.EditableTextOptions{Disabled: off})
		}},
		{"ColorSlot", "Paint", func(c *ui.Context, v *state, off bool) {
			input.ColorSlot(c, "Paint", &v.slot, input.ColorSlotOptions{Disabled: off})
		}},
		{"Joystick", "Stick", func(c *ui.Context, v *state, off bool) {
			if input.Joystick(c, "Stick", input.JoystickOptions{Disabled: off}).Moved {
				v.events++
			}
		}},
		{"CheatCode", "Cheat", func(c *ui.Context, v *state, off bool) {
			if input.CheatCode(c, "Cheat", input.CheatCodeOptions{Sequence: []ui.Key{ui.KeyUp}, Disabled: off}).Unlocked {
				v.events++
			}
		}},
		{"SlotPicker", "Reel", func(c *ui.Context, v *state, off bool) {
			input.SlotPicker(c, "Reel", &v.n, []string{"A", "B", "C"}, input.SlotPickerOptions{Disabled: off})
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var v state
			off := false
			tt := view(400, 300, nil, func(c *ui.Context) {
				tc.build(c, &v, off)
				ui.Box(c).Size(4, 4).Label("elsewhere")
			})
			click(t, tt, tc.label)
			if tc.name == "EditableText" {
				tt.Key(0, ui.KeyEnter)
			}
			off = true
			tt.Frame()
			if _, open := tt.Find("A"); open && tc.name != "Segmented" && tc.name != "SlotPicker" {
				t.Fatal("a disabled control's popup stayed open")
			}
			// Focused counts ancestors: only a shared one may hold it.
			if tt.Focused(tc.label) && !tt.Focused("elsewhere") {
				t.Fatal("a disabled control kept the focus")
			}
			before := fmt.Sprint(v)
			tt.Type("1a,")
			tt.SetClipboard("12")
			tt.Command("paste")
			for _, k := range []ui.Key{ui.KeyUp, ui.KeyDown, ui.KeyLeft, ui.KeyRight, ui.KeyHome, ui.KeySpace, ui.KeyEnter, ui.KeyBackspace} {
				tt.Key(0, k)
			}
			click(t, tt, tc.label)
			if now := fmt.Sprint(v); now != before {
				t.Fatalf("disabled control changed: %s to %s", before, now)
			}
		})
	}
}
