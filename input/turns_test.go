package input_test

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
)

// turnCase builds a control on fresh state, and reads it.
type turnCase func() (build func(c *ui.Context), value func() string)

func TestKeysInAVeilsFirstFrameTurnNothing(t *testing.T) {
	cases := map[string]struct {
		mk   turnCase
		tabs int
		key  ui.Key
	}{
		"Toggle": {func() (func(c *ui.Context), func() string) {
			on := false
			return func(c *ui.Context) { input.Toggle(c, "Mute", &on, input.ToggleOptions{}) }, func() string { return fmt.Sprint(on) }
		}, 1, ui.KeySpace},
		"Switch": {func() (func(c *ui.Context), func() string) {
			on := false
			return func(c *ui.Context) { input.Switch(c, "Sound", &on, input.SwitchOptions{}) }, func() string { return fmt.Sprint(on) }
		}, 1, ui.KeySpace},
		"Checkbox": {func() (func(c *ui.Context), func() string) {
			on := false
			return func(c *ui.Context) { input.Checkbox(c, "Hard", &on, input.CheckboxOptions{}) }, func() string { return fmt.Sprint(on) }
		}, 1, ui.KeySpace},
		"Radio": {func() (func(c *ui.Context), func() string) {
			v := "A"
			return func(c *ui.Context) {
				input.Radio(c, "A", "A", &v, input.RadioOptions{})
				input.Radio(c, "B", "B", &v, input.RadioOptions{})
			}, func() string { return v }
		}, 2, ui.KeySpace},
		"RadioGroup": {func() (func(c *ui.Context), func() string) {
			v := 0
			return func(c *ui.Context) { input.RadioGroup(c, []string{"A", "B", "C"}, &v, input.RadioGroupOptions{}) }, func() string { return fmt.Sprint(v) }
		}, 1, ui.KeyDown},
		"Rating": {func() (func(c *ui.Context), func() string) {
			v := 2
			return func(c *ui.Context) { input.Rating(c, "Stars", &v, input.RatingOptions{}) }, func() string { return fmt.Sprint(v) }
		}, 1, ui.KeyRight},
		"CheckboxGroup": {func() (func(c *ui.Context), func() string) {
			on := []bool{false, false}
			return func(c *ui.Context) {
				input.CheckboxGroup(c, "All", []string{"A", "B"}, on, input.CheckboxGroupOptions{})
			}, func() string { return fmt.Sprint(on) }
		}, 1, ui.KeySpace},
		"ToggleGroup": {func() (func(c *ui.Context), func() string) {
			on := []bool{false, false}
			return func(c *ui.Context) { input.ToggleGroup(c, []string{"B", "I"}, on, input.ToggleGroupOptions{}) }, func() string { return fmt.Sprint(on) }
		}, 1, ui.KeySpace},
		"Segmented": {func() (func(c *ui.Context), func() string) {
			v := 0
			return func(c *ui.Context) { input.Segmented(c, "View", []string{"A", "B"}, &v, input.SegmentedOptions{}) }, func() string { return fmt.Sprint(v) }
		}, 1, ui.KeyRight},
	}
	for name, tc := range cases {
		for _, veil := range []bool{false, true} {
			build, value := tc.mk()
			inert := false
			tt := underVeil(veilCase{build: build}, &inert)
			for range tc.tabs {
				tt.Key(0, ui.KeyTab)
			}
			before := value()
			inert = veil
			tt.Key(0, tc.key)
			if turned := value() != before; turned == veil {
				t.Errorf("%s: veiled %v, the key turned %s to %s", name, veil, before, value())
			}
		}
	}
}
