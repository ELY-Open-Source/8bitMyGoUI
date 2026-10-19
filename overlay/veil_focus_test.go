package overlay_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/overlay"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestPopconfirmGivesPopupFocusBackAfterAVeil(t *testing.T) {
	inert := false
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		s := pixel.Current(c)
		s.Inert = inert
		pixel.Scope(c, s, func() { overlay.Popconfirm(c, "DELETE", "Sure?", overlay.PopconfirmOptions{}) })
		pixel.Text(c, "Elsewhere", pixel.TextOptions{})
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyTab)
	// Focus in the popup, not the root that holds everything.
	if !tt.Focused("NO") || tt.Focused("Elsewhere") {
		t.Fatalf("focus is not in the popup: %q", tt.Texts())
	}
	inert = true
	tt.Frame()
	inert = false
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("Sure?") {
		t.Fatal("the button did not get its focus back after the veil")
	}
}

func TestSheetPullStopsUnderAVeil(t *testing.T) {
	open, inert := true, false
	tt := view(400, 400, newClock(), false, func(c *ui.Context) {
		s := pixel.Current(c)
		s.Inert = inert
		pixel.Scope(c, s, func() {
			overlay.Sheet(c, "Bag", &open, overlay.SheetOptions{}, func() { pixel.Text(c, "Potion", pixel.TextOptions{}) })
		})
	})
	r := find(t, tt, "Bag")
	tt.Press(r.X+r.W/2, r.Y+10)
	tt.Move(r.X+r.W/2, r.Y+60)
	inert = true
	tt.Move(r.X+r.W/2, r.Y+200)
	tt.Release(r.X+r.W/2, r.Y+200)
	if !open {
		t.Fatal("a pull under the veil closed the sheet")
	}
}

func TestPauseMenuPicksNothingAsAVeilFalls(t *testing.T) {
	paused, inert, chose := true, false, -1
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		s := pixel.Current(c)
		s.Inert = inert
		pixel.Scope(c, s, func() {
			if r := overlay.PauseMenu(c, &paused, overlay.PauseMenuOptions{}); r.Chosen >= 0 {
				chose = r.Chosen
			}
		})
	})
	inert = true
	tt.Key(0, ui.KeyEnter)
	if chose >= 0 {
		t.Fatalf("Enter as the veil fell chose %d", chose)
	}
}
