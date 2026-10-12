package input_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestCheckboxIgnoresAKeyAsItTurnsOff(t *testing.T) {
	for _, inert := range []bool{false, true} {
		on, off, changed := false, false, false
		tt := ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: (&clock{time.Unix(3_000_000, 0)}).now, Inert: off && inert})
			changed = input.Checkbox(c, "Box", &on, input.CheckboxOptions{Disabled: off && !inert}).Changed || changed
		}, 300, 100)
		tt.Key(0, ui.KeyTab)
		// Space, sent while on, lands as the box turns off.
		off = true
		tt.Key(0, ui.KeySpace)
		if on || changed {
			t.Errorf("inert %v: Space turned the box to %v, changed %v", inert, on, changed)
		}
	}
}
