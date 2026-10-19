package account_test

import (
	"fmt"
	"image"
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/account"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// offPalette counts pixels of a box no palette color matches.
func offPalette(img *image.RGBA, r ui.Rect) int {
	n := 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			c, hit := img.RGBAAt(x, y), false
			// PalettePicker shows every palette, so any will do.
			for _, pal := range pixel.Palettes {
				for _, p := range pal.Colors {
					hit = hit || p.R == c.R && p.G == c.G && p.B == c.B
				}
			}
			if !hit {
				n++
			}
		}
	}
	return n
}

func TestPatchNotesReadInChinese(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Locale: pixel.Zh})
		account.PatchNotes(c, "更新", []account.Patch{{Version: "v2", Added: []string{"船"}}}, account.PatchNotesOptions{})
	}, 400, 200)
	find(t, tt, "新增: 船")
}

func TestAccountStaysOnThePaletteWhenOff(t *testing.T) {
	builds := map[string]func(c *ui.Context, off bool){
		"access": func(c *ui.Context, off bool) {
			a := account.Access{Dot: 2}
			account.AccessibilityOptions(c, "X", &a, account.AccessibilityOptionsOptions{Disabled: off})
		},
		"creator": func(c *ui.Context, off bool) {
			l := account.Look{}
			account.CharacterCreator(c, "X", &l, account.CharacterCreatorOptions{Disabled: off})
		},
		"switcher": func(c *ui.Context, off bool) {
			i := 0
			account.AccountSwitcher(c, "X", [4]string{"Aria", "", "", ""}, &i, account.AccountSwitcherOptions{Disabled: off})
		},
		"feedback": func(c *ui.Context, off bool) {
			r, m := 3, "hi"
			account.FeedbackForm(c, "X", &r, &m, account.FeedbackFormOptions{Disabled: off})
		},
		"keys": func(c *ui.Context, off bool) {
			account.KeyBindings(c, "X", []account.Binding{{Action: "Jump", Key: ui.KeyZ}}, account.KeyBindingsOptions{Disabled: off})
		},
		"difficulty": func(c *ui.Context, off bool) {
			i := 1
			account.DifficultySelect(c, "X", &i, account.DifficultySelectOptions{Disabled: off})
		},
		"login": func(c *ui.Context, off bool) {
			n, pw := "a", "b"
			account.LoginForm(c, "X", &n, &pw, account.LoginFormOptions{Disabled: off})
		},
		"cycler": func(c *ui.Context, off bool) {
			i := 0
			account.OptionCycler(c, "X", []string{"A", "B"}, &i, account.OptionCyclerOptions{Disabled: off})
		},
		"palettes": func(c *ui.Context, off bool) {
			i := 0
			account.PalettePicker(c, "X", &i, account.PalettePickerOptions{Disabled: off})
		},
		"players": func(c *ui.Context, off bool) {
			i := 1
			account.PlayerSelect(c, "X", &i, account.PlayerSelectOptions{Disabled: off})
		},
		"grid": func(c *ui.Context, off bool) {
			pw := "AB"
			account.PasswordGrid(c, "X", &pw, account.PasswordGridOptions{Disabled: off})
		},
		"sessions": func(c *ui.Context, off bool) {
			account.SessionList(c, "X", []account.Session{{Device: "Phone", Seen: time.Unix(0, 0)}}, account.SessionListOptions{Disabled: off})
		},
		"mixer": func(c *ui.Context, off bool) {
			m := account.Mix{BGM: 0.5}
			account.VolumeMixer(c, "X", &m, account.VolumeMixerOptions{Disabled: off})
		},
		"twofactor": func(c *ui.Context, off bool) {
			code := "12"
			account.TwoFactor(c, "X", &code, account.TwoFactorOptions{Disabled: off})
		},
	}
	for name, build := range builds {
		for _, inert := range []bool{false, true} {
			wd := newWorld()
			wd.inert = inert
			tt := view(700, 400, wd, func(c *ui.Context) { build(c, !inert) })
			if n := offPalette(tt.Image(), ui.Rect{W: 700, H: 400}); n != 0 {
				t.Errorf("%s, inert %v: %d pixels off the palette", name, inert, n)
			}
		}
	}
}

func TestLoginFormLeavesFocusWhenOff(t *testing.T) {
	wd := newWorld()
	name, pass, fails := "a", "b", 0
	tt := view(500, 300, wd, func(c *ui.Context) {
		ui.Column(c).Children(func() {
			ui.ButtonBase(c).Label("elsewhere").Size(40, 20)
			account.LoginForm(c, "WELCOME", &name, &pass, account.LoginFormOptions{Failures: fails})
		})
	})
	click(t, tt, "elsewhere")
	// A failure lands in the frame Inert begins.
	fails, wd.inert = 1, true
	tt.Frame()
	if tt.Focused("PASSWORD") {
		t.Fatal("an inert form took the focus")
	}
}

func TestOptionCyclerKeysIgnoreThePointer(t *testing.T) {
	sel := 0
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		account.OptionCycler(c, "Speed", []string{"A", "B", "C"}, &sel, account.OptionCyclerOptions{})
	})
	r := find(t, tt, "Speed")
	tt.Move(r.X+2, r.Y+r.H/2)
	tt.Key(0, ui.KeyTab)
	for range 2 {
		tt.Key(0, ui.KeyEnter)
		tt.Key(0, ui.KeySpace)
	}
	if sel != 1 {
		t.Fatalf("four presses with the pointer on the left gave %d, want 1", sel)
	}
}

func TestKeyBindingsTakeAnyNextKey(t *testing.T) {
	binds := []account.Binding{{Action: "Help", Key: ui.KeyH}}
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		account.KeyBindings(c, "Keys", binds, account.KeyBindingsOptions{})
	})
	click(t, tt, "Keys")
	for _, k := range []ui.Key{ui.KeyF1, ui.KeySlash, ui.KeyPageDown} {
		tt.Key(0, k)
		if binds[0].Key != k {
			t.Fatalf("the next key, %v, did not bind: %v", k, binds[0].Key)
		}
		tt.Key(0, ui.KeyEnter)
	}
}

func TestTwoFactorTakesASCIIDigitsOnly(t *testing.T) {
	mustPanic(t, "fullwidth digits", func() {
		view(400, 200, newWorld(), func(c *ui.Context) {
			s := "１２"
			account.TwoFactor(c, "C", &s, account.TwoFactorOptions{})
		})
	})
}

func TestTwoFactorBarFollowsAnOddPeriod(t *testing.T) {
	wd := newWorld()
	// At 30 s, 1.1 s periods leave 0.8 s.
	wd.clk.t = time.Unix(30, 0)
	code := ""
	tt := view(400, 200, wd, func(c *ui.Context) {
		account.TwoFactor(c, "Code", &code, account.TwoFactorOptions{Period: 1100 * time.Millisecond})
	})
	x, y, w, _ := dots(find(t, tt, "TIME LEFT"))
	n := count(tt.Image(), x, y+2, w, 1, pixel.Bad)
	if want := (w - 4) * 8 / 11; n < want-1 || n > want+1 {
		t.Fatalf("the bar is %d dots, want %d", n, want)
	}
}

func TestProfileCardWantsMicroStatNames(t *testing.T) {
	mustPanic(t, "a stat named in hanzi", func() {
		view(400, 200, newWorld(), func(c *ui.Context) {
			account.ProfileCard(c, account.Profile{Name: "A", Level: 1, Stats: []account.Stat{{Name: "体力", Value: 1, Max: 2}}}, account.ProfileCardOptions{})
		})
	})
}

func TestKeyBindingsBindSpaceAfterEnter(t *testing.T) {
	binds := []account.Binding{{Action: "Jump", Key: ui.KeyZ}}
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		account.KeyBindings(c, "Keys", binds, account.KeyBindingsOptions{})
	})
	click(t, tt, "Keys")
	tt.Key(0, ui.KeyEscape)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeySpace)
	if binds[0].Key != ui.KeySpace {
		t.Fatalf("Enter, then Space bound %v", binds[0].Key)
	}
}

func TestTwoFactorTakesTextAndThePasteCommand(t *testing.T) {
	code, sent := "", 0
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		if account.TwoFactor(c, "Code", &code, account.TwoFactorOptions{}).Submitted {
			sent++
		}
	})
	click(t, tt, "Code")
	tt.Type("12")
	// Browsers paste through the command, as cmd/web sends it.
	tt.SetClipboard("a3b4c5d6e7")
	tt.Command("paste")
	if _, ok := tt.TextCaret(); !ok || code != "123456" || sent != 1 {
		t.Fatalf("typing and the paste command gave %q, sent %d", code, sent)
	}
}

func TestPalettePickerArrowsMoveByRows(t *testing.T) {
	chosen := 0
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		account.PalettePicker(c, "Palettes", &chosen, account.PalettePickerOptions{})
	})
	tt.Key(0, ui.KeyTab)
	got := []int{}
	for _, k := range []ui.Key{ui.KeyDown, ui.KeyRight, ui.KeyUp} {
		tt.Key(0, k)
		tt.Key(0, ui.KeyEnter)
		got = append(got, chosen)
	}
	if got[0] != 4 || got[1] != 5 || got[2] != 1 {
		t.Fatalf("Down, Right, Up chose %v, want [4 5 1]", got)
	}
}

func TestNegativeCountsPanic(t *testing.T) {
	mustPanic(t, "a width of -1", func() {
		view(400, 200, newWorld(), func(c *ui.Context) {
			i := 0
			account.OptionCycler(c, "S", []string{"A"}, &i, account.OptionCyclerOptions{Width: -1})
		})
	})
	mustPanic(t, "-1 failures", func() {
		view(400, 200, newWorld(), func(c *ui.Context) {
			n, pw := "", ""
			account.LoginForm(c, "L", &n, &pw, account.LoginFormOptions{Failures: -1})
		})
	})
}

func TestTwoFactorBarSurvivesALongPeriod(t *testing.T) {
	wd := newWorld()
	wd.clk.t = time.Unix(30, 0)
	code := ""
	tt := view(400, 200, wd, func(c *ui.Context) {
		account.TwoFactor(c, "Code", &code, account.TwoFactorOptions{Period: 4 * 365 * 24 * time.Hour})
	})
	x, y, w, _ := dots(find(t, tt, "TIME LEFT"))
	if n := count(tt.Image(), x, y+2, w, 1, pixel.Good); n < w-5 {
		t.Fatalf("a four-year period 30 s in fills %d of %d dots", n, w-4)
	}
}

func TestProfileCardFillsAtHugeValues(t *testing.T) {
	tt := view(1300, 300, newWorld(), func(c *ui.Context) {
		account.ProfileCard(c, account.Profile{Name: "A", Level: 1, Stats: []account.Stat{{Name: "HP", Value: math.MaxInt, Max: math.MaxInt}}}, account.ProfileCardOptions{Width: 600})
	})
	x, y, w, _ := dots(find(t, tt, "HP"))
	v := fmt.Sprintf("%d/%d", math.MaxInt, math.MaxInt)
	inner := w - pixel.MicroWidth("HP") - pixel.MicroWidth(v) - 12
	if n := count(tt.Image(), x, y+4, w, 1, pixel.Good); n != inner {
		t.Fatalf("a full meter of huge numbers fills %d of %d dots", n, inner)
	}
}

func TestProfileCardKeepsLongNamesOffTheBar(t *testing.T) {
	stats := []account.Stat{{Name: "AGILITY", Value: 1, Max: 2}, {Name: "HP", Value: 12500, Max: 12500}}
	tt := view(500, 300, newWorld(), func(c *ui.Context) {
		account.ProfileCard(c, account.Profile{Name: "A", Level: 1, Stats: stats}, account.ProfileCardOptions{})
	})
	// Both bars start after AGILITY and end before 12500/12500.
	nw, vw := pixel.MicroWidth("AGILITY"), pixel.MicroWidth("12500/12500")
	for _, s := range stats {
		x, y, w, _ := dots(find(t, tt, s.Name))
		inner := w - nw - vw - 12
		img := tt.Image()
		if n := count(img, x, y+4, w, 1, pixel.Good); n != inner*s.Value/s.Max || at(img, x+nw+6, y+4) != pixel.Good {
			t.Errorf("%s fills %d dots, want %d from dot %d", s.Name, n, inner*s.Value/s.Max, nw+6)
		}
	}
	// Frames under 3×3 panic, so a card too narrow fails.
	mustPanic(t, "no room for a bar", func() {
		view(500, 300, newWorld(), func(c *ui.Context) {
			account.ProfileCard(c, account.Profile{Name: "A", Level: 1, Stats: stats}, account.ProfileCardOptions{Width: 60})
		})
	})
}
