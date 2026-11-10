package diagram

import (
	"bytes"
	"image"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

const flow = "flowchart LR\n  accTitle: Order flow\n  Cart[Cart] --> Pay{Paid?}\n  Pay -- yes --> Ship([Ship order])\n" +
	"  Pay -- no --> Cart\n  style Ship fill:#4f9a3f,stroke:#33283a,color:#fff"

// show opens a window holding one Mermaid of *src.
func show(src *string, o MermaidOptions, w int, s *pixel.Settings) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, *s)
		Mermaid(c, *src, o)
	}, w, 400)
	tt.Advance(time.Second)
	return tt
}

// logged counts log lines holding s.
func logged(t *testing.T) func(s string) int {
	t.Helper()
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(old) })
	return func(s string) int { return strings.Count(buf.String(), s) }
}

func picture(t *testing.T, tt *ui.Tester, name string) ui.Rect {
	t.Helper()
	info, ok := tt.AccessByRole(ui.RoleImage, name)
	if !ok {
		t.Fatalf("no image named %q in %q", name, tt.Texts())
	}
	return info.Bounds
}

// matches holds a box of 2-DIP dots to want.
func matches(t *testing.T, img *image.RGBA, pal *pixel.Palette, box ui.Rect, want *plate) {
	t.Helper()
	for y := range want.h {
		for x := range want.w {
			c := pal.C(pixel.Back)
			if d := want.at(x, y); d != 0 {
				c = pal.C(pixel.Slot(d - 1))
			}
			for _, at := range [2][2]int{{0, 0}, {1, 1}} {
				if got := img.RGBAAt(int(box.X)+2*x+at[0], int(box.Y)+2*y+at[1]); got.R != c.R || got.G != c.G || got.B != c.B {
					t.Fatalf("dot %d,%d draws %v, want %v", x, y, got, c)
				}
			}
		}
	}
}

func scene4(t *testing.T, pal *pixel.Palette) *scene {
	t.Helper()
	sc, err := render(flow, tokens(pal), &pal.Colors)
	if err != nil {
		t.Fatal(err)
	}
	return sc
}

func TestMermaidDrawsItsDots(t *testing.T) {
	src, set := flow, pixel.Settings{Palette: pixel.SpringDay}
	tt := show(&src, MermaidOptions{}, 900, &set)
	sc := scene4(t, pixel.SpringDay)
	box := picture(t, tt, "Order flow")
	if box.X != 0 || box.Y != 0 || box.W != float32(2*sc.w) || box.H != float32(2*sc.h) {
		t.Fatalf("picture %+v, want %d by %d dots of 2 DIPs at the corner", box, sc.w, sc.h)
	}
	matches(t, tt.Image(), pixel.SpringDay, box, &sc.plate)
	if _, _, _, _, n := sc.box(pixel.Good); n < 500 {
		t.Errorf("%d dots of the styled node's green, want its fill", n)
	}
}

func TestMermaidShrinksOnWholeDots(t *testing.T) {
	src, set := flow, pixel.Settings{Palette: pixel.SpringDay}
	sc := scene4(t, pixel.SpringDay)
	for _, w := range []int{500, 300, 150} {
		tt := show(&src, MermaidOptions{}, w, &set)
		box := picture(t, tt, "Order flow")
		if want := float32(w) * float32(sc.h) / float32(sc.w); box.W != float32(w) || box.H < want-0.5 || box.H > want+0.5 {
			t.Fatalf("%d wide: picture %+v, want the scene's shape, %v tall", w, box, want)
		}
		n := ceil(sc.w, w/2)
		d := drawn{colors: pixel.SpringDay.Colors, sc: sc, fit: map[int][]cell{}}
		small := blank(ceil(sc.w, n), ceil(sc.h, n))
		for _, c := range d.cells(n) {
			for y := c.y; y < c.y+c.h; y++ {
				for x := c.x; x < c.x+c.w; x++ {
					small.dots[y*small.w+x] = uint8(c.slot) + 1
				}
			}
		}
		if n < 2 || small.w > w/2 || 2*small.h > int(box.H)+1 {
			t.Fatalf("%d wide: factor %d gives %d by %d dots for a box of %v", w, n, small.w, small.h, box)
		}
		matches(t, tt.Image(), pixel.SpringDay, box, small)
	}
}

func TestShrinkKeepsInkFarFromTheBackground(t *testing.T) {
	p := blank(5, 4)
	edge, face := uint8(pixel.Edge)+1, uint8(pixel.Face)+1
	copy(p.dots, []uint8{
		edge, edge, edge, edge, 0,
		edge, face, face, face, 0,
		edge, face, face, face, 0,
		0, 0, 0, 0, face,
	})
	var rank [pixel.Slots + 1]int
	rank[edge], rank[face] = 9, 2
	got := p.shrink(2, &rank)
	want := []uint8{edge, edge, 0, edge, face, face}
	if got.w != 3 || got.h != 2 || string(got.dots) != string(want) {
		t.Errorf("shrunk to %d by %d: %v, want 3 by 2: %v", got.w, got.h, got.dots, want)
	}
	if same := p.shrink(1, &rank); string(same.dots) != string(p.dots) {
		t.Error("a factor of 1 changed the plate")
	}
}

func TestCellsRankSlotsByExactDistance(t *testing.T) {
	for _, pal := range pixel.Palettes {
		back := pal.Colors[pixel.Back]
		far := func(s int) int {
			c := pal.Colors[s]
			dr, dg, db := int(c.R)-int(back.R), int(c.G)-int(back.G), int(c.B)-int(back.B)
			return dr*dr + dg*dg + db*db
		}
		for a := range pixel.Slots {
			for b := range pixel.Slots {
				if far(a) <= far(b) {
					continue
				}
				for _, dots := range [][]uint8{{uint8(a) + 1, uint8(b) + 1}, {uint8(b) + 1, uint8(a) + 1}} {
					d := drawn{colors: pal.Colors, sc: &scene{plate: plate{w: 2, h: 1, dots: dots}}, fit: map[int][]cell{}}
					if got := d.cells(2); len(got) != 1 || got[0].slot != pixel.Slot(a) {
						t.Fatalf("%s: slots %d and %d shrink to %+v, want the farther %d", pal.Name, a, b, got, a)
					}
				}
			}
		}
	}
}

func TestMermaidRendersOncePerInput(t *testing.T) {
	count := logged(t)
	src, set := flow, pixel.Settings{Palette: pixel.SpringDay}
	tt := show(&src, MermaidOptions{}, 900, &set)
	for range 5 {
		tt.Frame()
	}
	if count("mermaid drew") != 1 || count("mermaid fit") != 1 {
		t.Fatalf("six frames drew %d times and fit %d times", count("mermaid drew"), count("mermaid fit"))
	}
	tt.SetSize(300, 400)
	tt.Frame()
	tt.Frame()
	if count("mermaid drew") != 1 || count("mermaid fit") != 2 {
		t.Errorf("a narrower window drew %d times and fit %d times, want 1 and 2", count("mermaid drew"), count("mermaid fit"))
	}
	set.Palette = pixel.SpringNight
	tt.Frame()
	tt.Frame()
	if count("mermaid drew") != 2 {
		t.Errorf("a new palette drew %d times, want 2", count("mermaid drew"))
	}
	night := scene4(t, pixel.SpringNight)
	box := picture(t, tt, "Order flow")
	if got, want := tt.Image().RGBAAt(int(box.X)+1, int(box.Y)+1), pixel.SpringNight.C(pixel.Back); night.at(0, 0) != 0 || got.R != want.R || got.G != want.G {
		t.Errorf("the corner draws %v at night, want the night's Back %v", got, want)
	}
	if tt.FramePending() {
		t.Error("still asking for frames")
	}
}

func TestMermaidShowsPelesError(t *testing.T) {
	count := logged(t)
	src, set := "flowchart LR\n  Cart -->", pixel.Settings{Palette: pixel.SpringDay}
	tt := show(&src, MermaidOptions{}, 500, &set)
	for range 3 {
		tt.Frame()
	}
	if !tt.HasText("DIAGRAM FAILED") || count("mermaid failed") != 1 || count("line=3 column=1") != 1 {
		t.Fatalf("texts %q, failed logged %d", tt.Texts(), count("mermaid failed"))
	}
	if count("Cart") != 0 {
		t.Error("the log holds the source")
	}
	told := false
	for _, s := range tt.Texts() {
		told = told || strings.HasPrefix(s, "3:1 Parse error on line 3:") && strings.Contains(s, "Expecting 'NODE_STRING', got 'EOF'")
	}
	if !told {
		t.Errorf("texts %q lack Pele's message", tt.Texts())
	}
	if _, ok := tt.AccessByRole(ui.RoleImage, "DIAGRAM"); ok {
		t.Error("a failed diagram still draws")
	}
	src = "flowchart LR\n  Cart --> Pay"
	tt.Frame()
	tt.Advance(time.Second)
	picture(t, tt, "DIAGRAM")
	if tt.HasText("DIAGRAM FAILED") || count("mermaid drew") != 1 {
		t.Errorf("after the fix: texts %q, drew %d", tt.Texts(), count("mermaid drew"))
	}
	src = "no such diagram"
	tt.Frame()
	if !tt.HasText("No diagram type detected for the given text.") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMermaidLabels(t *testing.T) {
	src, set := flow, pixel.Settings{Palette: pixel.SpringDay}
	picture(t, show(&src, MermaidOptions{Label: "Checkout"}, 500, &set), "Checkout")
	src = "pie\n  \"Tea\" : 3\n  \"Coffee\" : 5"
	picture(t, show(&src, MermaidOptions{}, 500, &set), "DIAGRAM")
	set.Locale = pixel.Zh
	picture(t, show(&src, MermaidOptions{}, 500, &set), "图表")
	src = "not a diagram"
	if tt := show(&src, MermaidOptions{}, 500, &set); !tt.HasText("图表出错") {
		t.Errorf("texts %q", tt.Texts())
	}
}
