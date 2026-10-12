package account_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/account"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestSettingsListSetsRowsUnderTitles(t *testing.T) {
	tt := view(500, 300, newWorld(), func(c *ui.Context) {
		account.SettingsList(c, "Settings", []account.SettingsSection{
			{Title: "AUDIO", Settings: []account.Setting{{Label: "Music", Control: func() {
				ui.ButtonBase(c).Label("toggle").Size(20, 20)
			}}}},
		}, account.SettingsListOptions{})
	})
	// The row label only shows; the control carries the name.
	title, ctl := find(t, tt, "AUDIO"), find(t, tt, "toggle")
	if _, named := tt.Find("Music"); named || !(title.Y+title.H <= ctl.Y && ctl.X > 100) {
		t.Fatalf("title %v, control %v, label named %v", title, ctl, named)
	}
}

func TestOptionCyclerWrapsBothWays(t *testing.T) {
	sel, changed := 0, false
	opts := []string{"SLOW", "NORMAL", "FAST"}
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		changed = account.OptionCycler(c, "Speed", opts, &sel, account.OptionCyclerOptions{}).Changed || changed
	})
	r := find(t, tt, "Speed")
	tt.ClickAt(r.X+4, r.Y+r.H/2)
	if sel != 2 || !changed {
		t.Fatalf("a click on the left half from 0 gave %d", sel)
	}
	tt.Key(0, ui.KeyRight)
	if sel != 0 {
		t.Fatalf("Right from the last gave %d", sel)
	}
	tt.ClickAt(r.X+r.W-4, r.Y+r.H/2)
	tt.Key(0, ui.KeyLeft)
	if sel != 0 {
		t.Fatalf("right half then Left gave %d", sel)
	}
	mustPanic(t, "a choice past the options", func() {
		view(400, 200, newWorld(), func(c *ui.Context) {
			s := 3
			account.OptionCycler(c, "S", opts, &s, account.OptionCyclerOptions{})
		})
	})
}

func TestKeyBindingsRebindAndMarkConflicts(t *testing.T) {
	binds := []account.Binding{{Action: "Jump", Key: ui.KeySpace}, {Action: "Attack", Key: ui.KeyZ}, {Action: "Dash", Key: ui.KeyX}}
	changed := false
	tt := view(500, 300, newWorld(), func(c *ui.Context) {
		changed = account.KeyBindings(c, "Keys", binds, account.KeyBindingsOptions{}).Changed || changed
	})
	r := find(t, tt, "Keys")
	x, y, w, _ := dots(r)
	rh := pixel.LineHeight + 4
	bad := func(row int) int { return count(tt.Image(), x+w-40, y+row*rh, 40, rh, pixel.Bad) }
	if bad(1) != 0 || bad(2) != 0 {
		t.Fatal("no conflict yet, but a key is Bad")
	}
	// Click Dash, press Z: Z is Attack's too.
	tt.ClickAt(r.X+20, r.Y+float32(2*2*rh+rh))
	tt.Key(0, ui.KeyZ)
	if binds[2].Key != ui.KeyZ || !changed || bad(1) == 0 || bad(2) == 0 {
		t.Fatalf("Dash on Z: %v, conflict marks %d %d", binds[2].Key, bad(1), bad(2))
	}
	// Waiting, an arrow binds; Escape gives up.
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyLeft)
	if binds[1].Key != ui.KeyLeft {
		t.Fatalf("Up, Enter, Left bound %v", binds[1].Key)
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyEscape)
	tt.Key(0, ui.KeyQ)
	if binds[1].Key != ui.KeyLeft {
		t.Fatal("Escape did not give up the rebind")
	}
	mustPanic(t, "a key without a cap", func() {
		view(500, 300, newWorld(), func(c *ui.Context) {
			account.KeyBindings(c, "K", []account.Binding{{Action: "Help", Key: ui.KeyEscape}}, account.KeyBindingsOptions{})
		})
	})
}

func TestControllerMapPaintsButtonsAndLabels(t *testing.T) {
	tt := view(700, 300, newWorld(), func(c *ui.Context) {
		account.ControllerMap(c, "Pad", map[account.PadButton]string{account.ButtonA: "JUMP"}, account.ControllerMapOptions{})
	})
	x, y, _, _ := dots(find(t, tt, "Pad"))
	img := tt.Image()
	if at(img, x+70+130, y+45) != pixel.Bad || at(img, x+70+120, y+55) != pixel.Warn || count(img, x+236, y+46, 40, 12, pixel.Ink) == 0 {
		t.Fatal("A is not red, B not gold, or JUMP is missing")
	}
	mustPanic(t, "an unknown button", func() {
		view(700, 300, newWorld(), func(c *ui.Context) {
			account.ControllerMap(c, "P", map[account.PadButton]string{99: "?"}, account.ControllerMapOptions{})
		})
	})
}

func TestProfileCardShowsStats(t *testing.T) {
	tt := view(500, 300, newWorld(), func(c *ui.Context) {
		account.ProfileCard(c, account.Profile{Name: "Aria", Title: "Slayer", Level: 12, Stats: []account.Stat{{Name: "HP", Value: 10, Max: 40}}}, account.ProfileCardOptions{})
	})
	find(t, tt, "Slayer")
	find(t, tt, "LV 12")
	x, y, w, _ := dots(find(t, tt, "HP"))
	// The bar sits between the name and value columns.
	nw, vw := pixel.MicroWidth("HP"), pixel.MicroWidth("10/40")
	inner := w - nw - vw - 12
	if n := count(tt.Image(), x+nw+6, y+4, inner, 1, pixel.Good); n < inner/4-1 || n > inner/4+1 {
		t.Fatalf("HP 10/40 fills %d of %d dots", n, inner)
	}
	mustPanic(t, "a stat past its most", func() {
		view(500, 300, newWorld(), func(c *ui.Context) {
			account.ProfileCard(c, account.Profile{Name: "A", Level: 1, Stats: []account.Stat{{Name: "HP", Value: 5, Max: 4}}}, account.ProfileCardOptions{})
		})
	})
}

func TestCharacterCreatorCyclesTheAvatar(t *testing.T) {
	look, changed := account.Look{}, false
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		changed = account.CharacterCreator(c, "Hero", &look, account.CharacterCreatorOptions{}).Changed || changed
	})
	find(t, tt, "SPIKY, DOTS, TUNIC, DARK, DARK")
	r := find(t, tt, "HAIR COLOR")
	tt.ClickAt(r.X+r.W-4, r.Y+r.H/2)
	tt.ClickAt(r.X+r.W-4, r.Y+r.H/2)
	img := tt.Image()
	pic := find(t, tt, "SPIKY, DOTS, TUNIC, RED, DARK")
	x, y, _, _ := dots(pic)
	// Spiky hair at row 2, column 5, scale 3.
	if look.HairColor != 2 || !changed || at(img, x+4+3*5, y+4+3*2) != pixel.Bad {
		t.Fatalf("hair color %d, hair dot %v", look.HairColor, at(tt.Image(), x+4+3*5, y+4+3*2))
	}
}
