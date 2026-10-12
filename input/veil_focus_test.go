package input_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestMultiSelectGivesPopupFocusBackAfterAVeil(t *testing.T) {
	chosen, inert := []string{}, false
	tt := view(400, 300, &clock{time.Unix(3_000_000, 0)}, func(c *ui.Context) {
		s := pixel.Current(c)
		s.Inert = inert
		pixel.Scope(c, s, func() {
			input.MultiSelect(c, "Party", &chosen, []string{"Mage", "Rogue"}, input.MultiSelectOptions{})
		})
		pixel.Text(c, "Elsewhere", pixel.TextOptions{})
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	// Focus on Mage, not the root that holds everything.
	if !tt.Focused("Mage") || tt.Focused("Elsewhere") {
		t.Fatalf("focus is not in the popup: %q", tt.Texts())
	}
	inert = true
	tt.Frame()
	inert = false
	tt.Frame()
	tt.Key(0, ui.KeyDown)
	if !tt.HasText("Rogue") {
		t.Fatal("the field did not get its focus back after the veil")
	}
}
