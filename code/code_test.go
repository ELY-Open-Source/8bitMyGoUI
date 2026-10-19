package code_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/code"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

// world is the settings a test view runs under.
type world struct {
	clk           *clock
	reduce, inert bool
}

func newWorld() *world {
	return &world{clk: &clock{time.Date(2026, 9, 2, 10, 30, 0, 0, time.UTC)}}
}

// view builds on a fake clock, two pixels a dot.
func view(w, h int, wd *world, build func(c *ui.Context)) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: wd.clk.now, Inert: wd.inert})
		build(c)
	}, w, h)
	if wd.reduce {
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tt.Frame()
	}
	return tt
}

// later moves the clock on and builds a frame.
func later(tt *ui.Tester, wd *world, d time.Duration) {
	wd.clk.t = wd.clk.t.Add(d)
	tt.Frame()
}

func click(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	if err := tt.Click(label); err != nil {
		t.Fatalf("%v in %q", err, tt.Texts())
	}
}

func find(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("no %q in %q", label, tt.Texts())
	}
	return r
}

// none is a color outside the palette.
const none pixel.Slot = 255

// at is the slot at dot x, y, or none.
func at(img *image.RGBA, x, y int) pixel.Slot {
	c := img.RGBAAt(2*x, 2*y)
	for s, want := range pixel.SpringDay.Colors {
		if c.R == want.R && c.G == want.G && c.B == want.B {
			return pixel.Slot(s)
		}
	}
	return none
}

// count is how many dots in a box hold s.
func count(img *image.RGBA, x, y, w, h int, s pixel.Slot) int {
	n := 0
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			if at(img, i, j) == s {
				n++
			}
		}
	}
	return n
}

// dots is a found box in dots.
func dots(r ui.Rect) (x, y, w, h int) {
	return int(r.X / 2), int(r.Y / 2), int(r.W / 2), int(r.H / 2)
}

func mustPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s did not panic", name)
		}
	}()
	fn()
}

func TestCodeBlockNumbersColorsAndCopies(t *testing.T) {
	wd := newWorld()
	copied := false
	tt := view(500, 300, wd, func(c *ui.Context) {
		copied = code.CodeBlock(c, "main.go", "package main\n\nfunc main() {\n\tprintln(1)\n}\n", code.CodeBlockOptions{}).Copied || copied
	})
	// Tabs become four spaces; each line is read alone.
	for _, l := range []string{"package main", "func main() {", "    println(1)", "}"} {
		find(t, tt, l)
	}
	x, y, w, h := dots(find(t, tt, "package main"))
	img := tt.Image()
	if count(img, x, y, w, h, pixel.Accent) == 0 || count(img, x, y, 10, h, pixel.Faint) == 0 {
		t.Fatal("no Accent keyword or Faint line number on line 1")
	}
	click(t, tt, "COPY")
	if !copied || tt.Clipboard() != "package main\n\nfunc main() {\n    println(1)\n}" {
		t.Fatalf("copied %v: %q", copied, tt.Clipboard())
	}
	find(t, tt, "COPIED!")
	later(tt, wd, 2*time.Second)
	find(t, tt, "COPY")
}

func TestTerminalRunsAndRecalls(t *testing.T) {
	wd := newWorld()
	lines, input := []string{"READY."}, ""
	var ran []string
	tt := view(500, 300, wd, func(c *ui.Context) {
		if r := code.Terminal(c, "Shell", lines, &input, code.TerminalOptions{}); r.Ran != "" {
			ran = append(ran, r.Ran)
			lines = append(lines, "$ "+r.Ran)
		}
	})
	find(t, tt, "READY.")
	click(t, tt, "COMMAND")
	for _, cmd := range []string{"ls", "pwd", "pwd"} {
		tt.Type(cmd)
		tt.Key(0, ui.KeyEnter)
	}
	if len(ran) != 3 || ran[0] != "ls" || input != "" {
		t.Fatalf("ran %q, input %q", ran, input)
	}
	find(t, tt, "$ pwd")
	// A draft waits below the history; repeats are kept once.
	tt.Type("ec")
	steps := []struct {
		key  ui.Key
		want string
	}{{ui.KeyUp, "pwd"}, {ui.KeyUp, "ls"}, {ui.KeyUp, "ls"}, {ui.KeyDown, "pwd"}, {ui.KeyDown, "ec"}, {ui.KeyDown, "ec"}}
	for i, s := range steps {
		tt.Key(0, s.key)
		if input != s.want {
			t.Fatalf("step %d: input %q, want %q", i, input, s.want)
		}
	}
}

func TestTerminalHoldsUnderInert(t *testing.T) {
	wd := newWorld()
	input := ""
	tt := view(500, 300, wd, func(c *ui.Context) {
		code.Terminal(c, "Shell", nil, &input, code.TerminalOptions{})
	})
	click(t, tt, "COMMAND")
	tt.Type("ls")
	tt.Key(0, ui.KeyEnter)
	wd.inert = true
	tt.Key(0, ui.KeyUp)
	if input != "" {
		t.Fatalf("an inert terminal recalled %q", input)
	}
}
