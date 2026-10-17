package overlay_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/overlay"
)

func TestDialogsHandTheFocusBack(t *testing.T) {
	kinds := map[string]func(c *ui.Context, open *bool){
		"confirm": func(c *ui.Context, open *bool) {
			overlay.ConfirmDialog(c, "Sure?", open, overlay.ConfirmDialogOptions{})
		},
		"choice": func(c *ui.Context, open *bool) {
			overlay.ChoicePrompt(c, "Pick", []string{"A", "B"}, open, overlay.ChoicePromptOptions{})
		},
		"prompt": func(c *ui.Context, open *bool) {
			v := ""
			overlay.PromptDialog(c, "Name", &v, open, overlay.PromptDialogOptions{})
		},
		"pause": func(c *ui.Context, open *bool) { overlay.PauseMenu(c, open, overlay.PauseMenuOptions{}) },
	}
	for name, show := range kinds {
		open := false
		tt := view(500, 400, newClock(), false, func(c *ui.Context) {
			if input.Button(c, "Opener", input.ButtonOptions{}).Clicked {
				open = true
			}
			show(c, &open)
		})
		if err := tt.Click("Opener"); err != nil || !open || tt.Focused("Opener") {
			t.Fatalf("%s: did not open with the focus inside (%v)", name, err)
		}
		tt.Key(0, ui.KeyEscape)
		if open || !tt.Focused("Opener") {
			t.Fatalf("%s: Escape left it open or the focus lost", name)
		}
	}
}
