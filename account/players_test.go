package account_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/account"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestLanguagePickerChooses(t *testing.T) {
	chosen := 0
	var got account.LanguagePickerResult
	tt := view(400, 300, newWorld(), func(c *ui.Context) {
		if r := account.LanguagePicker(c, "Language", []account.Language{{Name: "English"}, {Name: "Deutsch"}}, &chosen, account.LanguagePickerOptions{}); r.Changed {
			got = r
		}
	})
	click(t, tt, "Deutsch")
	if chosen != 1 || !got.Changed {
		t.Fatalf("clicking Deutsch chose %d", chosen)
	}
}

func TestPalettePickerTilesWearTheirPalettes(t *testing.T) {
	chosen, changed := 0, false
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		changed = account.PalettePicker(c, "Palettes", &chosen, account.PalettePickerOptions{}).Changed || changed
	})
	r := find(t, tt, "Palettes")
	// Tile 3's title bar is Summer Night's own Accent.
	want := pixel.SummerNight.Colors[pixel.Accent]
	if c := tt.Image().RGBAAt(int(r.X)+2*(3*56+12), int(r.Y)+2*7); c.R != want.R || c.G != want.G || c.B != want.B {
		t.Fatalf("Summer Night's tile title is %v, want %v", c, want)
	}
	tt.ClickAt(r.X+2*(3*56+28), r.Y+2*20)
	x, y, _, _ := dots(r)
	// The chosen tile wears the app's Accent ring.
	if chosen != 3 || !changed || at(tt.Image(), x+3*56+28, y) != pixel.Accent {
		t.Fatalf("clicking Summer Night chose %d", chosen)
	}
}

func TestVolumeMixerFaders(t *testing.T) {
	mix, changed := account.Mix{BGM: 0.5, SE: 0.5, Voice: 0.5}, false
	tt := view(400, 300, newWorld(), func(c *ui.Context) {
		changed = account.VolumeMixer(c, "Mix", &mix, account.VolumeMixerOptions{}).Changed || changed
	})
	r := find(t, tt, "SE")
	tt.Press(r.X+r.W/2, r.Y+2)
	tt.Release(r.X+r.W/2, r.Y+2)
	if mix.SE < 0.9 || mix.BGM != 0.5 || mix.Voice != 0.5 || !changed {
		t.Fatalf("pressing the top of SE gave %+v", mix)
	}
	mustPanic(t, "a level past 1", func() {
		view(400, 300, newWorld(), func(c *ui.Context) {
			m := account.Mix{BGM: 2}
			account.VolumeMixer(c, "M", &m, account.VolumeMixerOptions{})
		})
	})
}

func TestDifficultySelectMovesThePointer(t *testing.T) {
	chosen, changed := 1, false
	tt := view(400, 300, newWorld(), func(c *ui.Context) {
		changed = account.DifficultySelect(c, "Difficulty", &chosen, account.DifficultySelectOptions{}).Changed || changed
	})
	click(t, tt, "Difficulty")
	tt.Key(0, ui.KeyEnd)
	x, y, _, _ := dots(find(t, tt, "Difficulty"))
	img := tt.Image()
	// The pointer is on NIGHTMARE, which is Bad.
	if chosen != 3 || !changed || count(img, x, y+3*16, 12, 16, pixel.Accent) == 0 || count(img, x+14, y+3*16, 60, 16, pixel.Bad) == 0 {
		t.Fatalf("End chose %d", chosen)
	}
}

func TestAccessibilityOptionsSetBoth(t *testing.T) {
	a, changed := account.Access{Dot: 2}, false
	tt := view(500, 300, newWorld(), func(c *ui.Context) {
		changed = account.AccessibilityOptions(c, "Access", &a, account.AccessibilityOptionsOptions{}).Changed || changed
	})
	click(t, tt, "REDUCE MOTION")
	click(t, tt, "3")
	if !a.Reduce || a.Dot != 3 || !changed {
		t.Fatalf("after the switch and 3: %+v", a)
	}
	mustPanic(t, "a dot of 5", func() {
		view(500, 300, newWorld(), func(c *ui.Context) {
			b := account.Access{Dot: 5}
			account.AccessibilityOptions(c, "A", &b, account.AccessibilityOptionsOptions{})
		})
	})
}

func TestAccountSwitcherChoosesSlots(t *testing.T) {
	cur, changed := 0, false
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		changed = account.AccountSwitcher(c, "Players", [4]string{"Aria", "", "", ""}, &cur, account.AccountSwitcherOptions{}).Changed || changed
	})
	find(t, tt, "P1, Aria, IN USE")
	click(t, tt, "P3, OPEN")
	if cur != 2 || !changed {
		t.Fatalf("P3 chose %d", cur)
	}
	// P3 is green; chosen, its frame is too.
	x, y, _, _ := dots(find(t, tt, "P3, OPEN, IN USE"))
	if at(tt.Image(), x+30, y) != pixel.Good {
		t.Fatal("the chosen card does not wear its player color")
	}
	// Aria's eyes sit 16 dots into her 32-dot face.
	x, y, _, _ = dots(find(t, tt, "P1, Aria"))
	if count(tt.Image(), x+16, y+30, 32, 4, pixel.Ink) == 0 {
		t.Fatal("Aria's face has no eyes")
	}
}

func TestPlayerSelectChooses(t *testing.T) {
	players := 1
	var got account.PlayerSelectResult
	tt := view(400, 300, newWorld(), func(c *ui.Context) {
		if r := account.PlayerSelect(c, "Players", &players, account.PlayerSelectOptions{}); r.Changed || r.Chosen {
			got = r
		}
	})
	click(t, tt, "Players")
	tt.Key(0, ui.KeyDown)
	if players != 2 || !got.Changed {
		t.Fatalf("Down gave %d players", players)
	}
	got = account.PlayerSelectResult{}
	tt.Key(0, ui.KeyEnter)
	if !got.Chosen || players != 2 {
		t.Fatalf("Enter: %+v", got)
	}
	mustPanic(t, "three players", func() {
		view(400, 300, newWorld(), func(c *ui.Context) {
			p := 3
			account.PlayerSelect(c, "P", &p, account.PlayerSelectOptions{})
		})
	})
}

func TestAboutAndPatchNotesRead(t *testing.T) {
	tt := view(500, 400, newWorld(), func(c *ui.Context) {
		ui.Row(c).Children(func() {
			account.About(c, "GAME", "1.2", account.AboutOptions{Credits: []account.Credit{{Role: "Art", Name: "Ann"}}})
			account.PatchNotes(c, "Notes", []account.Patch{{Version: "v2", Added: []string{"Boats"}, Fixed: []string{"Rain"}}}, account.PatchNotesOptions{Width: 120})
		})
	})
	for _, l := range []string{"VERSION 1.2", "CREDITS", "Art", "Ann", "ADDED: Boats", "FIXED: Rain"} {
		find(t, tt, l)
	}
	x, y, _, _ := dots(find(t, tt, "ADDED: Boats"))
	if count(tt.Image(), x, y, 30, 12, pixel.Good) == 0 {
		t.Fatal("the ADDED tag is not Good")
	}
	mustPanic(t, "no version", func() {
		view(500, 400, newWorld(), func(c *ui.Context) { account.About(c, "G", "", account.AboutOptions{}) })
	})
}
