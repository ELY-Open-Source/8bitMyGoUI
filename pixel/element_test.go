package pixel_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestSpriteLettersMatchSlots(t *testing.T) {
	sp := pixel.MustSprite("kbflsenwaAigyrcm")
	d := canvas(t, 16, 1, func(g pixel.Canvas) { g.Sprite(0, 0, sp, 0, pixel.Ink, 1) })
	for i := range 16 {
		if got := d.at(i, 0); got != pixel.Slot(i) {
			t.Fatalf("letter %d drew slot %d", i, got)
		}
	}
}

func TestSilhouetteKeepsShapeInOneSlot(t *testing.T) {
	sp := pixel.MustSprite("#a.g")
	d := canvas(t, 8, 2, func(g pixel.Canvas) { g.Silhouette(0, 0, sp, 0, pixel.Shade, 2) })
	if d.count(0, 0, 4, 2, pixel.Shade) != 8 || d.count(4, 0, 2, 2, pixel.Back) != 4 || d.count(6, 0, 2, 2, pixel.Shade) != 4 {
		t.Fatal("Silhouette must fill the shape, scaled, in one slot")
	}
}

func TestSpriteFramesMirrorScale(t *testing.T) {
	sp := pixel.MustSprite(`
		#a.
		..g

		.r.
		.r.
	`)
	if sp.W != 3 || sp.H != 2 || sp.Frames() != 2 {
		t.Fatalf("sprite %d×%d with %d frames", sp.W, sp.H, sp.Frames())
	}
	d := canvas(t, 6, 4, func(g pixel.Canvas) {
		g.Sprite(0, 0, sp, 0, pixel.Magic, 1)
		g.Sprite(3, 0, sp, 3, pixel.Ink, 1)
		g.Sprite(0, 2, sp.Mirror(), 0, pixel.Magic, 1)
	})
	checks := map[[2]int]pixel.Slot{
		{0, 0}: pixel.Magic, {1, 0}: pixel.Accent, {2, 0}: pixel.Back, {2, 1}: pixel.Good,
		{4, 0}: pixel.Bad, {4, 1}: pixel.Bad, {3, 0}: pixel.Back,
		{2, 2}: pixel.Magic, {1, 2}: pixel.Accent, {0, 3}: pixel.Good,
	}
	for p, want := range checks {
		if got := d.at(p[0], p[1]); got != want {
			t.Fatalf("dot %v = %d, want %d", p, got, want)
		}
	}
	big := canvas(t, 6, 4, func(g pixel.Canvas) { g.Sprite(0, 0, sp, 0, pixel.Magic, 2) })
	if big.count(0, 0, 2, 2, pixel.Magic) != 4 || big.count(2, 0, 2, 2, pixel.Accent) != 4 {
		t.Fatal("scale 2 must fill 2×2 dots per art pixel")
	}
	for name, art := range map[string]string{"ragged": "ab\nabc", "letter": "az", "empty": "  ", "frames": "aa\n\naa\naa"} {
		mustPanic(t, name, func() { pixel.MustSprite(art) })
	}
}

func TestClock(t *testing.T) {
	at := time.Unix(1_700_000_001, 250_000_000)
	var step, busy, tick int
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }})
		k := pixel.K(c)
		tick = k.Tick(c, 8)
		ui.Box(c).Size(4, 4).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			step, busy = g.Step(4), g.Busy(4)
		})
	}, 10, 10)
	if want := int(at.UnixNano() / int64(time.Second/4)); step != want || busy != want {
		t.Fatalf("Step %d Busy %d, want %d", step, busy, want)
	}
	if want := int(at.UnixNano() / int64(time.Second/8)); tick != want {
		t.Fatalf("Tick %d, want %d", tick, want)
	}
	pixel.ResetDue()
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if step != 0 || busy == 0 {
		t.Fatalf("under Reduce Motion: Step %d (want 0), Busy %d", step, busy)
	}
	if due := pixel.Due().Sub(tt.Now()); due <= 0 || due > 125*time.Millisecond {
		t.Fatalf("Due in %v, want the 8 fps wait on the frame clock", due)
	}
	pixel.ResetDue()
	if !pixel.Due().IsZero() {
		t.Fatal("ResetDue left a deadline")
	}
	mustPanic(t, "fps 0", func() {
		ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay})
			pixel.K(c).Tick(c, 0)
		}, 10, 10)
	})
	mustPanic(t, "Step 0 under Reduce Motion", func() {
		on := false
		tt := ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay})
			k := pixel.K(c)
			ui.Box(c).Size(4, 4).Draw(func(p *ui.Painter, r ui.Rect) {
				if on {
					k.Canvas(p, r).Step(0)
				}
			})
		}, 10, 10)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		on = true
		tt.Frame()
	})
}

func TestWrapFollowsSameFrameResize(t *testing.T) {
	const s = "像素风格组件库，每个点都经过认真摆放，在限制中创造丰富。"
	width := float32(120)
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Dot: 2})
		ui.Column(c).Width(width).Children(func() {
			pixel.Text(c, s, pixel.TextOptions{Wrap: true})
		})
	}, 400, 800)
	for _, w := range []float32{60, 200, 120} {
		width = w
		tt.Frame()
		r, _ := tt.Find(s)
		if want := float32(len(pixel.Lines(s, int(w/2))) * 16 * 2); r.H != want {
			t.Fatalf("width %v: text height %v, want %v", w, r.H, want)
		}
	}
}

func TestTextElement(t *testing.T) {
	var width float32
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Dot: 2})
		pixel.Text(c, "HELLO", pixel.TextOptions{Scale: 2})
		ui.Column(c).Width(120).Children(func() {
			pixel.Text(c, "像素风格组件库，每个点都经过认真摆放。", pixel.TextOptions{Wrap: true})
		})
		pixel.Text(c, "one\ntwo\nthree", pixel.TextOptions{MaxLines: 2})
		width = float32(pixel.Measure("HELLO") * 2 * 2)
	}, 400, 400)
	r, ok := tt.Find("HELLO")
	if !ok || r.W != width || r.H != 2*16*2 {
		t.Fatalf("HELLO box %v (found %v), want %v×64", r, ok, width)
	}
	r, _ = tt.Find("像素风格组件库，每个点都经过认真摆放。")
	if lines := len(pixel.Lines("像素风格组件库，每个点都经过认真摆放。", 60)); r.H != float32(lines*16*2) || lines < 3 {
		t.Fatalf("wrapped box %v, want %d lines", r, lines)
	}
	if r, _ = tt.Find("one\ntwo\nthree"); r.H != 2*16*2 {
		t.Fatalf("MaxLines box height %v, want 64", r.H)
	}
}

func shellView(disabled *bool, clicks *int, st *pixel.State) func(c *ui.Context) {
	at := time.Unix(1_700_000_001, 250_000_000)
	return func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SummerDay, Dot: 2, Now: func() time.Time { return at }})
		k := pixel.K(c)
		ui.Box(c).Size(k.D(16), k.D(4))
		ui.Row(c).Children(func() {
			ui.Box(c).Size(k.D(16), k.D(8))
			shell, base, s := pixel.Shell(c, func() ui.Element {
				return ui.ButtonBase(c).Label("go").Disabled(*disabled)
			})
			*st = s
			if base.Clicked() {
				*clicks++
			}
			shell.Size(k.D(12), k.D(8)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Rect(0, 0, 12, 8, pixel.Accent)
				if s.Focus {
					g.FocusMark(0, 0, 12, 8)
				}
			})
		})
	}
}

func TestShellKeepsLookAndBehavior(t *testing.T) {
	disabled, clicks, st := false, 0, pixel.State{}
	tt := ui.NewTester(shellView(&disabled, &clicks, &st), 100, 100)
	tt.ClickAt(40, 16)
	if clicks != 1 {
		t.Fatalf("enabled shell: %d clicks", clicks)
	}
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	if !st.Focus {
		t.Fatal("Tab did not give the base visible focus")
	}
	d := dots{t: t, img: tt.Image(), pal: pixel.SummerDay}
	if d.at(14, 2) != pixel.Accent || d.at(15, 2) != pixel.Back || d.at(29, 13) != pixel.Accent {
		t.Fatal("focus ring dots missing around the shell")
	}
	if d.count(3, 5, 12, 6, pixel.Edge) == 0 || d.count(3, 5, 12, 6, pixel.Lit) == 0 {
		t.Fatal("hand cursor missing left of the shell")
	}
	disabled = true
	tt.Frame()
	tt.ClickAt(40, 16)
	if clicks != 1 || !st.Disabled || st.Focus {
		t.Fatalf("disabled: clicks %d, state %+v", clicks, st)
	}
	d = dots{t: t, img: tt.Image(), pal: pixel.SummerDay}
	if d.count(16, 4, 12, 8, pixel.Accent) != 96 {
		t.Fatal("disabled base dimmed the shell's look")
	}
}

func TestVeilAnchorsAcrossElements(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SummerDay, Dot: 2})
		k := pixel.K(c)
		ui.Box(c).Size(k.D(1), k.D(1))
		ui.Row(c).Children(func() {
			ui.Box(c).Size(k.D(3), k.D(2))
			ui.Box(c).Size(k.D(4), k.D(2)).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Checker(0, 0, 4, 2, pixel.Ink)
			})
		})
	}, 14, 6)
	d := dots{t: t, img: tt.Image(), pal: pixel.SummerDay}
	if d.at(3, 1) != pixel.Ink || d.at(4, 1) != pixel.Back || d.at(3, 2) != pixel.Back || d.at(4, 2) != pixel.Ink {
		t.Fatal("checker phase follows the element, not the window")
	}
}
