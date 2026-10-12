package overlay_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestDialogActsClosesAndHandsFocusBack(t *testing.T) {
	open, action, closed := false, -1, 0
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		if input.Button(c, "OPEN", input.ButtonOptions{}).Clicked {
			open = true
		}
		r := overlay.Dialog(c, "Save game?", &open, overlay.DialogOptions{Actions: []string{"LATER", "SAVE"}}, func() {
			pixel.Text(c, "Slot 1 is free.", pixel.TextOptions{})
		})
		if r.Action >= 0 {
			action = r.Action
		}
		if r.Closed {
			closed++
		}
	})
	click(t, tt, "OPEN")
	if !shows(tt, "Slot 1 is free.") {
		t.Fatal("the dialog did not open")
	}
	click(t, tt, "SAVE")
	if open || action != 1 || shows(tt, "Slot 1 is free.") {
		t.Fatalf("SAVE left open %v, action %d", open, action)
	}
	if !tt.Focused("OPEN") {
		t.Fatal("closing did not hand the focus back")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyEscape)
	if open || closed != 1 {
		t.Fatalf("Escape left open %v, %d closes", open, closed)
	}
	keeps(t, tt.Image())
}

func TestAlertDialogEnterPressesDefault(t *testing.T) {
	open, pressed := true, -1
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		if r := overlay.AlertDialog(c, "Low HP", "Heal before the boss?", &open, overlay.AlertDialogOptions{Buttons: []string{"NO", "HEAL"}, Default: 1}); r.Pressed >= 0 {
			pressed = r.Pressed
		}
	})
	tt.Key(0, ui.KeyEnter)
	if open || pressed != 1 {
		t.Fatalf("Enter pressed %d, open %v", pressed, open)
	}
	open = true
	tt.Frame()
	click(t, tt, "NO")
	if open || pressed != 0 {
		t.Fatalf("NO pressed %d", pressed)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("a Default past the buttons did not panic")
		}
	}()
	view(100, 100, newClock(), false, func(c *ui.Context) {
		overlay.AlertDialog(c, "x", "y", &open, overlay.AlertDialogOptions{Default: 1})
	})
}

func TestConfirmDialogStartsOnNo(t *testing.T) {
	open, yes, no := true, 0, 0
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		r := overlay.ConfirmDialog(c, "Drop the sword?", &open, overlay.ConfirmDialogOptions{})
		if r.Confirmed {
			yes++
		}
		if r.Canceled {
			no++
		}
	})
	if !shows(tt, "ARE YOU SURE?") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if open || yes != 0 || no != 1 {
		t.Fatalf("Enter on the start gave yes %d no %d", yes, no)
	}
	open = true
	tt.Frame()
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	if open || yes != 1 {
		t.Fatalf("Up, Enter gave yes %d", yes)
	}
	open = true
	tt.Frame()
	click(t, tt, "YES")
	tt.Frame()
	open = true
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	if yes != 2 || no != 2 {
		t.Fatalf("click YES then Escape gave yes %d no %d", yes, no)
	}
	keeps(t, tt.Image())
}

func TestChoicePromptTakesDigitsArrowsAndClicks(t *testing.T) {
	open, chosen, canceled := true, -1, 0
	choices := []string{"Fight", "Talk", "Flee"}
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		r := overlay.ChoicePrompt(c, "A slime appears!", choices, &open, overlay.ChoicePromptOptions{})
		if r.Chosen >= 0 {
			chosen = r.Chosen
		}
		if r.Canceled {
			canceled++
		}
	})
	tt.Key(0, ui.Key3)
	if open || chosen != 2 {
		t.Fatalf("3 chose %d", chosen)
	}
	open = true
	tt.Frame()
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeySpace)
	if chosen != 1 {
		t.Fatalf("Down, Space chose %d", chosen)
	}
	open = true
	tt.Frame()
	click(t, tt, "Flee")
	open = true
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	if chosen != 2 || canceled != 1 || open {
		t.Fatalf("click Flee, Escape: chosen %d, %d cancels", chosen, canceled)
	}
}

func TestPauseMenuResumesAndChooses(t *testing.T) {
	paused, chosen := true, -1
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		if r := overlay.PauseMenu(c, &paused, overlay.PauseMenuOptions{}); r.Chosen >= 0 {
			chosen = r.Chosen
		}
	})
	if !shows(tt, "PAUSED") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if !paused || chosen != 1 {
		t.Fatalf("OPTIONS left paused %v, chosen %d", paused, chosen)
	}
	tt.Key(0, ui.KeyEscape)
	if paused || chosen != 0 {
		t.Fatalf("Escape left paused %v, chosen %d", paused, chosen)
	}
	keeps(t, tt.Image())
}

func TestPromptDialogSubmitsAndCancels(t *testing.T) {
	open, name, sent, canceled := true, "", 0, 0
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		r := overlay.PromptDialog(c, "ENTER YOUR NAME", &name, &open, overlay.PromptDialogOptions{MaxLength: 8})
		if r.Submitted {
			sent++
		}
		if r.Canceled {
			canceled++
		}
	})
	tt.Type("LINK")
	tt.Key(0, ui.KeyEnter)
	if open || name != "LINK" || sent != 1 {
		t.Fatalf("typing, Enter left %q, open %v, %d sent", name, open, sent)
	}
	open = true
	tt.Frame()
	click(t, tt, "CANCEL")
	if open || canceled != 1 {
		t.Fatalf("CANCEL left open %v, %d cancels", open, canceled)
	}
}

func TestProgressDialogCancels(t *testing.T) {
	open, canceled := true, 0
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		if overlay.ProgressDialog(c, "Saving", 0.5, &open, overlay.ProgressDialogOptions{}).Canceled {
			canceled++
		}
	})
	if !shows(tt, "Saving") {
		t.Fatalf("texts %q", tt.Texts())
	}
	click(t, tt, "CANCEL")
	if open || canceled != 1 {
		t.Fatalf("CANCEL left open %v", open)
	}
	keeps(t, tt.Image())
}

func TestOffDialogsTakeNoAnswers(t *testing.T) {
	for _, inert := range []bool{false, true} {
		open, action := true, -1
		tt := view(640, 400, newClock(), inert, func(c *ui.Context) {
			if r := overlay.Dialog(c, "Quest", &open, overlay.DialogOptions{Actions: []string{"ACCEPT"}, Disabled: !inert}, func() {}); r.Action >= 0 {
				action = r.Action
			}
		})
		click(t, tt, "ACCEPT")
		if action >= 0 || !open {
			t.Fatalf("inert %v: a dialog took action %d", inert, action)
		}
		alert, pressed := true, -1
		tt = view(640, 400, newClock(), inert, func(c *ui.Context) {
			if r := overlay.AlertDialog(c, "Warn", "Careful.", &alert, overlay.AlertDialogOptions{Disabled: !inert}); r.Pressed >= 0 {
				pressed = r.Pressed
			}
		})
		tt.Key(0, ui.KeyEnter)
		click(t, tt, "OK")
		if pressed >= 0 || !alert {
			t.Fatalf("inert %v: an alert took button %d", inert, pressed)
		}
	}
}

func TestInertMenusTakeNoPicks(t *testing.T) {
	open, paused, chosen, picked := true, true, -1, -1
	tt := view(640, 400, newClock(), true, func(c *ui.Context) {
		if r := overlay.ChoicePrompt(c, "Go?", []string{"Yes", "No"}, &open, overlay.ChoicePromptOptions{}); r.Chosen >= 0 {
			chosen = r.Chosen
		}
	})
	tt.Key(0, ui.Key1)
	click(t, tt, "No")
	tt2 := view(640, 400, newClock(), true, func(c *ui.Context) {
		if r := overlay.PauseMenu(c, &paused, overlay.PauseMenuOptions{}); r.Chosen > 0 {
			picked = r.Chosen
		}
	})
	click(t, tt2, "QUIT")
	if chosen >= 0 || picked >= 0 || !open {
		t.Fatalf("inert menus picked %d and %d", chosen, picked)
	}
}

func TestDisabledMenusAndPromptsTakeNoAnswers(t *testing.T) {
	open, chosen, sure, paused, picked, prompt, name, sent := true, -1, true, true, -1, true, "", 0
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		if r := overlay.ChoicePrompt(c, "Go?", []string{"Yes", "No"}, &open, overlay.ChoicePromptOptions{Disabled: true}); r.Chosen >= 0 {
			chosen = r.Chosen
		}
	})
	tt.Key(0, ui.Key1)
	tt2 := view(640, 400, newClock(), false, func(c *ui.Context) {
		if overlay.ConfirmDialog(c, "Sure?", &sure, overlay.ConfirmDialogOptions{Disabled: true}).Confirmed {
			chosen = 99
		}
	})
	click(t, tt2, "YES")
	tt3 := view(640, 400, newClock(), false, func(c *ui.Context) {
		if r := overlay.PauseMenu(c, &paused, overlay.PauseMenuOptions{Disabled: true}); r.Chosen > 0 {
			picked = r.Chosen
		}
	})
	click(t, tt3, "QUIT")
	tt4 := view(640, 400, newClock(), false, func(c *ui.Context) {
		if overlay.PromptDialog(c, "NAME", &name, &prompt, overlay.PromptDialogOptions{Disabled: true}).Submitted {
			sent++
		}
	})
	tt4.Type("AB")
	click(t, tt4, "OK")
	if chosen >= 0 || picked >= 0 || sent != 0 || name != "" || !open || !sure || !prompt {
		t.Fatalf("disabled answers: chosen %d picked %d sent %d name %q", chosen, picked, sent, name)
	}
}

func TestFocusComesBackAfterAnOffSpellInADialog(t *testing.T) {
	inert, disabled := false, false
	open, action := true, -1
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Inert: inert})
		if r := overlay.Dialog(c, "Quest", &open, overlay.DialogOptions{Actions: []string{"ACCEPT"}}, func() {}); r.Action >= 0 {
			action = r.Action
		}
	}, 640, 400)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("ACCEPT") {
		t.Fatal("Tab did not reach ACCEPT")
	}
	inert = true
	tt.Frame()
	inert = false
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if action != 0 {
		t.Fatalf("Enter after the veil gave action %d", action)
	}
	chosen, choice := -1, true
	tt = view(640, 400, newClock(), false, func(c *ui.Context) {
		if r := overlay.ChoicePrompt(c, "Go?", []string{"Yes", "No"}, &choice, overlay.ChoicePromptOptions{Disabled: disabled}); r.Chosen >= 0 {
			chosen = r.Chosen
		}
	})
	disabled = true
	tt.Frame()
	disabled = false
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if chosen != 0 {
		t.Fatalf("Enter after a disabled spell chose %d", chosen)
	}
	name, prompt := "", true
	tt = view(640, 400, newClock(), false, func(c *ui.Context) {
		overlay.PromptDialog(c, "NAME", &name, &prompt, overlay.PromptDialogOptions{Disabled: disabled})
	})
	disabled = true
	tt.Frame()
	disabled = false
	tt.Frame()
	tt.Type("AB")
	if name != "AB" {
		t.Fatalf("typing after a disabled spell gave %q", name)
	}
}
