package layout_test

import (
	"image"
	"image/color"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func use(c *ui.Context) { pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay}) }

func box(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("no %q in %q", label, tt.Texts())
	}
	return r
}

func dragBy(t *testing.T, tt *ui.Tester, label string, dx, dy float32) {
	t.Helper()
	r := box(t, tt, label)
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Move(x+dx, y+dy)
	tt.Release(x+dx, y+dy)
}

func TestSplitDragKeysClamp(t *testing.T) {
	size, changes := 100, 0
	opts := layout.SplitOptions{Label: "Divide"}
	tt := ui.NewTester(func(c *ui.Context) {
		use(c)
		r := layout.Split(c, &size, opts, func() {}, func() {})
		r.Element.Size(400, 100)
		if r.Changed {
			changes++
		}
	}, 400, 100)
	dragBy(t, tt, "Divide", 20, 0)
	if size != 110 || changes == 0 {
		t.Fatalf("drag of 10 dots left size %d, %d changes", size, changes)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyLeft)
	if size != 114 {
		t.Fatalf("keys left size %d, want 114", size)
	}
	size = 1000
	tt.Frame()
	if size != 200-6-24 {
		t.Fatalf("size clamped to %d, want %d", size, 200-6-24)
	}
	size = 3
	tt.Frame()
	if size != 24 {
		t.Fatalf("size clamped to %d, want the minimum 24", size)
	}
	opts.Disabled = true
	dragBy(t, tt, "Divide", 40, 0)
	if size != 24 {
		t.Fatalf("a disabled grip moved to %d", size)
	}
}

func TestSplitRelaysOutWhenItsRoomChanges(t *testing.T) {
	size := 1000
	tt := ui.NewTester(func(c *ui.Context) {
		use(c)
		layout.Split(c, &size, layout.SplitOptions{}, func() {}, func() { ui.Box(c).Fill().Label("SECOND") }).Element.Fill()
	}, 400, 100)
	if second := box(t, tt, "SECOND"); size != 200-6-24 || second.W != 48 {
		t.Fatalf("first frames left size %d and the second pane %v", size, second)
	}
	tt.SetSize(240, 100)
	if second := box(t, tt, "SECOND"); size != 120-6-24 || second.W != 48 {
		t.Fatalf("after shrinking, size %d and the second pane %v", size, second)
	}
}

func TestScrollAreaArrowsTrackThumb(t *testing.T) {
	s := ui.ScrollState{}
	opts := layout.ScrollAreaOptions{Label: "List"}
	tt := ui.NewTester(func(c *ui.Context) {
		use(c)
		layout.ScrollArea(c, &s, opts, func() {
			for range 40 {
				ui.Box(c).Height(32)
			}
		}).Element.Size(200, 100)
	}, 300, 200)
	tt.Frame()
	if s.MaxY != 40*32-100 {
		t.Fatalf("MaxY %v", s.MaxY)
	}
	click := func(label string) {
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	click("SCROLL DOWN")
	if s.Y != 32 {
		t.Fatalf("down arrow scrolled to %v, want 16 dots", s.Y)
	}
	click("SCROLL UP")
	click("SCROLL UP")
	if s.Y != 0 {
		t.Fatalf("up arrow scrolled to %v, want 0", s.Y)
	}
	up := box(t, tt, "SCROLL UP")
	tt.ClickAt(up.X+up.W/2, 100-up.H-4)
	if s.Y != 100 {
		t.Fatalf("track below the thumb paged to %v, want 100", s.Y)
	}
	// Track 60 DIPs, thumb 16: one thumb DIP scrolls MaxY/44.
	at := up.Y + up.H + 2
	tt.Press(up.X+up.W/2, at+4)
	tt.Move(up.X+up.W/2, at+4+11)
	tt.Release(up.X+up.W/2, at+4+11)
	if want := float32(100) + 11*s.MaxY/44; s.Y != want {
		t.Fatalf("thumb drag scrolled to %v, want %v", s.Y, want)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("List") {
		t.Fatal("Tab did not reach the scroll view")
	}
	onPalette(t, tt.Image())
	opts.Disabled = true
	tt.Frame()
	if tt.Focused("List") {
		t.Fatal("a disabled scroll area kept the focus")
	}
	frozen := s.Y
	_ = tt.Click("SCROLL DOWN")
	tt.Scroll(50, 50, 0, 80)
	tt.Key(0, ui.KeyDown)
	tt.Move(50, 50)
	tt.Key(0, ui.KeyPageDown)
	if s.Y != frozen {
		t.Fatalf("a disabled scroll area moved from %v to %v", frozen, s.Y)
	}
}

// onPalette fails on any pixel off the Spring Day palette.
func onPalette(t *testing.T, img *image.RGBA) {
	t.Helper()
	ok := map[color.RGBA]bool{}
	for _, c := range pixel.SpringDay.Colors {
		ok[color.RGBA{c.R, c.G, c.B, 255}] = true
	}
	for y := range img.Rect.Dy() {
		for x := range img.Rect.Dx() {
			if p := img.RGBAAt(x, y); !ok[p] {
				t.Fatalf("pixel %d,%d is %v, off the palette", x, y, p)
			}
		}
	}
}

func TestTitleBarBoxes(t *testing.T) {
	var closes, minimizes int
	opts := layout.TitleBarOptions{}
	tt := ui.NewTester(func(c *ui.Context) {
		use(c)
		r := layout.TitleBar(c, "Notes", opts)
		closes += boolInt(r.Close)
		minimizes += boolInt(r.Minimize)
	}, 300, 60)
	for _, l := range []string{"CLOSE", "MINIMIZE", "MINIMIZE"} {
		if err := tt.Click(l); err != nil {
			t.Fatal(err)
		}
	}
	if closes != 1 || minimizes != 2 {
		t.Fatalf("closes %d minimizes %d", closes, minimizes)
	}
	opts.NoMinimize = true
	tt.Frame()
	if _, ok := tt.Find("MINIMIZE"); ok {
		t.Fatal("NoMinimize still shows the minimize box")
	}
	opts.Disabled = true
	tt.Frame()
	_ = tt.Click("CLOSE")
	if closes != 1 {
		t.Fatal("a disabled close box closed")
	}
}

func TestResizableGripKeysClamp(t *testing.T) {
	w, h, changes := 100, 60, 0
	opts := layout.ResizableOptions{MaxW: 150}
	tt := ui.NewTester(func(c *ui.Context) {
		use(c)
		if layout.Resizable(c, &w, &h, opts, func() {}).Changed {
			changes++
		}
	}, 600, 400)
	dragBy(t, tt, "RESIZE", 20, 10)
	if w != 110 || h != 65 || changes == 0 {
		t.Fatalf("drag left %d×%d, %d changes", w, h, changes)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	if w != 114 || h != 69 {
		t.Fatalf("keys left %d×%d, want 114×69", w, h)
	}
	w, h = 1000, 1
	tt.Frame()
	if w != 150 || h != 24 {
		t.Fatalf("clamped to %d×%d, want 150×24", w, h)
	}
	opts.Disabled = true
	dragBy(t, tt, "RESIZE", -40, 0)
	if w != 150 {
		t.Fatalf("a disabled grip resized to %d", w)
	}
}

func TestHandheldButtons(t *testing.T) {
	pressed := map[string]int{}
	opts := layout.HandheldOptions{ScreenW: 80, ScreenH: 56}
	locale := pixel.En
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Locale: locale})
		r := layout.Handheld(c, opts, func() {})
		for name, on := range map[string]bool{"UP": r.Up, "DOWN": r.Down, "LEFT": r.Left, "RIGHT": r.Right, "A": r.A, "B": r.B} {
			pressed[name] += boolInt(on)
		}
	}, 600, 500)
	names := []string{"UP", "DOWN", "LEFT", "RIGHT", "A", "B"}
	for _, n := range names {
		if err := tt.Click(n); err != nil {
			t.Fatal(err)
		}
	}
	for _, n := range names {
		if pressed[n] != 1 {
			t.Fatalf("presses %v", pressed)
		}
	}
	locale = pixel.Zh
	tt.Frame()
	if err := tt.Click("左"); err != nil || pressed["LEFT"] != 2 {
		t.Fatalf("Chinese labels: %v, %v", err, pressed)
	}
	opts.Disabled = true
	tt.Frame()
	_ = tt.Click("A")
	if pressed["A"] != 1 {
		t.Fatal("a disabled A button pressed")
	}
}

func TestGridColumnsAreWholeDots(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		use(c)
		layout.Grid(c, 3, layout.GridOptions{Columns: 3, Lines: layout.SolidLines}, func(i int) {
			ui.Box(c).Size(10, 10).Label([]string{"G0", "G1", "G2"}[i])
		}).Width(200)
	}, 300, 100)
	for i, x := range []float32{8, 76, 142} {
		if r := box(t, tt, []string{"G0", "G1", "G2"}[i]); r.X != x {
			t.Fatalf("cell %d content at x %v, want %v", i, r.X, x)
		}
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unknown grid lines did not panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		use(c)
		layout.Grid(c, 1, layout.GridOptions{Lines: 9}, func(int) {})
	}, 100, 100)
}

func TestAppFrameFolds(t *testing.T) {
	text := func(c *ui.Context, s string) func() { return func() { pixel.Text(c, s, pixel.TextOptions{}) } }
	tt := ui.NewTester(func(c *ui.Context) {
		use(c)
		layout.AppFrame(c, layout.AppFrameOptions{}, nil, text(c, "SIDE"), text(c, "MAIN"), nil)
	}, 800, 300)
	tt.Frame()
	side, main := box(t, tt, "SIDE"), box(t, tt, "MAIN")
	if side.X+side.W > main.X {
		t.Fatalf("wide frame: side %v not left of main %v", side, main)
	}
	tt.SetSize(500, 300)
	tt.Frame()
	side, main = box(t, tt, "SIDE"), box(t, tt, "MAIN")
	if side.Y+side.H > main.Y {
		t.Fatalf("narrow frame did not fold: side %v not above main %v", side, main)
	}
}

func TestAspectBoxAndViewport(t *testing.T) {
	var aspect, screen *ui.Element
	tt := ui.NewTester(func(c *ui.Context) {
		use(c)
		ui.Column(c).Children(func() {
			ui.Box(c).Width(300).Children(func() {
				aspect = layout.AspectBox(c, layout.AspectBoxOptions{W: 16, H: 9}, func() {})
			})
			ui.Box(c).Size(400, 300).Children(func() {
				layout.Viewport(c, layout.ViewportOptions{W: 64, H: 40}, func() { screen = ui.Box(c).Fill() })
			})
		})
	}, 500, 700)
	tt.Frame()
	if b := aspect.Bounds(); b.W != 300 || b.H != 2*(150*9/16) {
		t.Fatalf("aspect box %v", b)
	}
	if b := screen.Bounds(); b.X != 72 || b.Y != aspect.Bounds().H+70 || b.W != 256 || b.H != 160 {
		t.Fatalf("viewport screen %v", b)
	}
}

func TestLayersSizes(t *testing.T) {
	var auto, sized, fill *ui.Element
	tt := ui.NewTester(func(c *ui.Context) {
		use(c)
		ui.Column(c).AlignItems(ui.Start).Children(func() {
			auto = layout.Layers(c, func() { ui.Box(c).Size(50, 30) }, func() { ui.Box(c).Fill() })
			sized = layout.Layers(c, func() { fill = ui.Box(c).Fill() }).Size(120, 64)
		})
	}, 400, 300)
	tt.Frame()
	if b := auto.Bounds(); b.W != 50 || b.H != 30 {
		t.Fatalf("first layer did not size the box: %v", b)
	}
	if b := fill.Bounds(); b.W != 120 || b.H != 64 {
		t.Fatalf("first layer did not fill a sized box: %v", b)
	}
	if b := sized.Bounds(); b.H != 64 {
		t.Fatalf("sized box %v", b)
	}
}

func TestTitledContainersAndOrder(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		use(c)
		layout.Window(c, layout.WindowOptions{Title: "ITEMS"}, func() {})
		layout.Section(c, "Sound", layout.SectionOptions{}, func() {})
		layout.Card(c, layout.CardOptions{Title: "QUEST", Footer: func() { pixel.Text(c, "FOOT", pixel.TextOptions{}) }}, func() {
			pixel.Text(c, "BODY", pixel.TextOptions{})
		})
		layout.StatusBar(c, "Ready", "Ln 1")
	}, 400, 400)
	got := tt.Texts()
	for _, l := range []string{"ITEMS", "Sound", "QUEST", "Ready", "Ln 1"} {
		if !slices.Contains(got, l) {
			t.Fatalf("no %q in %q", l, got)
		}
	}
	if slices.Index(got, "BODY") > slices.Index(got, "FOOT") || box(t, tt, "BODY").Y >= box(t, tt, "FOOT").Y {
		t.Fatalf("card footer came before its body: %q", got)
	}
}

func TestDividerStylePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("unknown divider style did not panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		use(c)
		layout.Divider(c, layout.DividerOptions{Style: 9})
	}, 100, 100)
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
