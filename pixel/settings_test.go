package pixel_test

import (
	"image/color"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func mustPanic(t *testing.T, what string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatalf("%s: no panic", what)
		}
	}()
	fn()
}

func TestUseSetsKitAndBackground(t *testing.T) {
	var k pixel.Kit
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.WinterNight, Dot: 3, Locale: pixel.Zh})
		k = pixel.K(c)
	}, 40, 30)
	if k.Palette() != pixel.WinterNight || k.Dot() != 3 || k.D(5) != 15 || k.Locale() != pixel.Zh {
		t.Fatalf("kit = palette %s dot %v D(5) %v locale %d", k.Palette().Name, k.Dot(), k.D(5), k.Locale())
	}
	if k.C(pixel.Ink) != pixel.WinterNight.C(pixel.Ink) {
		t.Fatalf("C(Ink) = %v", k.C(pixel.Ink))
	}
	back := pixel.WinterNight.C(pixel.Back)
	want := color.RGBA{back.R, back.G, back.B, 255}
	if got := tt.Image().RGBAAt(20, 15); got != want {
		t.Fatalf("window background %v, want %v", got, want)
	}
}

func TestUseBridgesTheme(t *testing.T) {
	for _, p := range []*pixel.Palette{pixel.SummerDay, pixel.AutumnNight} {
		var th ui.Theme
		ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: p})
			th = *c.Theme()
		}, 10, 10)
		pairs := []struct {
			name      string
			got, want ui.Color
		}{
			{"Background", th.Background, p.C(pixel.Back)},
			{"Surface", th.Surface, p.C(pixel.Face)},
			{"SurfaceHover", th.SurfaceHover, p.C(pixel.Lit)},
			{"SurfacePressed", th.SurfacePressed, p.C(pixel.Shade)},
			{"Border", th.Border, p.C(pixel.Edge)},
			{"Text", th.Text, p.C(pixel.Ink)},
			{"TextMuted", th.TextMuted, p.C(pixel.Faint)},
			{"Accent", th.Accent, p.C(pixel.Accent)},
			{"AccentHover", th.AccentHover, p.C(pixel.AccentLit)},
			{"AccentPressed", th.AccentPressed, p.C(pixel.Accent)},
			{"AccentText", th.AccentText, p.C(pixel.AccentInk)},
			{"Danger", th.Danger, p.C(pixel.Bad)},
			{"Warning", th.Warning, p.C(pixel.Warn)},
			{"Success", th.Success, p.C(pixel.Good)},
			{"Selection", th.Selection, p.C(pixel.AccentLit)},
			{"Focus", th.Focus, p.C(pixel.Accent)},
			{"Scrollbar", th.Scrollbar, ui.Transparent},
		}
		for _, pr := range pairs {
			if pr.got != pr.want {
				t.Errorf("%s: theme %s = %v, want %v", p.Name, pr.name, pr.got, pr.want)
			}
		}
		if th.Radius != 0 || th.Dark != p.Night {
			t.Errorf("%s: Radius %v Dark %v", p.Name, th.Radius, th.Dark)
		}
	}
}

func TestLazyKeepsScopeForLayoutRows(t *testing.T) {
	seen := map[int]*pixel.Palette{}
	ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay})
		pixel.Scope(c, pixel.Settings{Palette: pixel.WinterNight}, func() {
			ui.List(c, nil, 100, pixel.Lazy(c, func(i int) {
				seen[i] = pixel.K(c).Palette()
				ui.Box(c).Height(10)
			})).Height(200)
		})
	}, 100, 220)
	if len(seen) < 20 {
		t.Fatalf("list built %d rows, want at least 20", len(seen))
	}
	for i, p := range seen {
		if p != pixel.WinterNight {
			t.Fatalf("row %d built under %s, want Winter Night", i, p.Name)
		}
	}
}

func TestDotDefaultsToTwo(t *testing.T) {
	var k pixel.Kit
	ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay})
		k = pixel.K(c)
	}, 10, 10)
	if k.Dot() != 2 {
		t.Fatalf("Dot = %v, want 2", k.Dot())
	}
}

func TestScopeOverridesAndRestores(t *testing.T) {
	var inside, after *pixel.Palette
	var insideDot float32
	ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.AutumnDay})
		pixel.Scope(c, pixel.Settings{Palette: pixel.SummerNight, Dot: 4}, func() {
			k := pixel.K(c)
			inside, insideDot = k.Palette(), k.Dot()
		})
		after = pixel.K(c).Palette()
	}, 10, 10)
	if inside != pixel.SummerNight || insideDot != 4 {
		t.Fatalf("inside Scope: %s dot %v", inside.Name, insideDot)
	}
	if after != pixel.AutumnDay {
		t.Fatalf("after Scope: %s, want Autumn Day", after.Name)
	}
}

func TestKWithoutUsePanics(t *testing.T) {
	mustPanic(t, "K before Use", func() {
		ui.NewTester(func(c *ui.Context) { pixel.K(c) }, 10, 10)
	})
}

func TestUseIsRequiredEveryFrame(t *testing.T) {
	skip := false
	tt := ui.NewTester(func(c *ui.Context) {
		if !skip {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay})
		}
		pixel.K(c)
	}, 10, 10)
	skip = true
	time.Sleep(time.Millisecond)
	mustPanic(t, "frame without Use", tt.Frame)
}

func TestInvalidSettingsPanic(t *testing.T) {
	bad := map[string]pixel.Settings{
		"nil palette": {},
		"dot 5":       {Palette: pixel.SpringDay, Dot: 5},
		"dot -1":      {Palette: pixel.SpringDay, Dot: -1},
		"locale 9":    {Palette: pixel.SpringDay, Locale: 9},
	}
	for name, s := range bad {
		mustPanic(t, name, func() {
			ui.NewTester(func(c *ui.Context) { pixel.Use(c, s) }, 10, 10)
		})
	}
}

func TestNowOverride(t *testing.T) {
	fixed := time.Date(2026, 10, 11, 9, 30, 0, 0, time.UTC)
	var k pixel.Kit
	ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return fixed }})
		k = pixel.K(c)
	}, 10, 10)
	if !k.Now().Equal(fixed) {
		t.Fatalf("Now = %v, want %v", k.Now(), fixed)
	}
}

func TestReduceMotionFollowsPreferences(t *testing.T) {
	var k pixel.Kit
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay})
		k = pixel.K(c)
	}, 10, 10)
	if k.Reduce() {
		t.Fatal("Reduce on by default")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if !k.Reduce() {
		t.Fatal("Reduce off after ReduceMotion preference")
	}
}
