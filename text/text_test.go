package text_test

import (
	"image"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/text"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func (c *clock) add(d time.Duration) { c.t = c.t.Add(d) }

func view(w, h int, clk *clock, build func(c *ui.Context)) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: clk.now})
		build(c)
	}, w, h)
}

func newClock() *clock { return &clock{time.Unix(1_000_000, 0)} }

// slot is the palette slot of dot x, y.
func slot(t *testing.T, img *image.RGBA, x, y int) pixel.Slot {
	t.Helper()
	c := img.RGBAAt(2*x, 2*y)
	for s, pc := range pixel.SpringDay.Colors {
		if pc.R == c.R && pc.G == c.G && pc.B == c.B {
			return pixel.Slot(s)
		}
	}
	t.Fatalf("dot %d,%d off the palette: %v", x, y, c)
	return 0
}

// count is how many dots in a box hold s.
func count(t *testing.T, img *image.RGBA, x, y, w, h int, s pixel.Slot) int {
	n := 0
	for j := range h {
		for i := range w {
			if slot(t, img, x+i, y+j) == s {
				n++
			}
		}
	}
	return n
}

func find(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("no %q in %q", label, tt.Texts())
	}
	return r
}

func TestHeadingLevels(t *testing.T) {
	tt := view(600, 300, newClock(), func(c *ui.Context) {
		ui.Column(c).Children(func() {
			text.Heading(c, "BIG", text.HeadingOptions{})
			text.Heading(c, "small", text.HeadingOptions{Level: 3})
		})
	})
	if r := find(t, tt, "BIG"); r.W != float32((pixel.Measure("BIG")*3+1)*2) || r.H != float32((pixel.LineHeight*3+1)*2) {
		t.Fatalf("level 1 box %v", r)
	}
	if r := find(t, tt, "small"); r.H != float32((pixel.LineHeight+1)*2) {
		t.Fatalf("level 3 box %v", r)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("level 4 did not panic")
		}
	}()
	view(100, 100, newClock(), func(c *ui.Context) { text.Heading(c, "x", text.HeadingOptions{Level: 4}) })
}

func TestParagraphAlignAndLimit(t *testing.T) {
	words := "one two three four five six seven eight nine ten"
	tt := view(400, 300, newClock(), func(c *ui.Context) {
		ui.Column(c).Width(200).Children(func() {
			text.Paragraph(c, words, text.ParagraphOptions{MaxLines: 1})
			text.Paragraph(c, "hi", text.ParagraphOptions{Align: text.Center})
			text.Paragraph(c, "hi", text.ParagraphOptions{Align: text.Right})
		})
	})
	img := tt.Image()
	if r := find(t, tt, words); r.H != pixel.LineHeight*2 {
		t.Fatalf("MaxLines 1 paragraph is %v DIPs tall", r.H)
	}
	left := func(row int) int {
		for x := range 100 {
			for y := row; y < row+pixel.LineHeight; y++ {
				if slot(t, img, x, y) == pixel.Ink {
					return x
				}
			}
		}
		return -1
	}
	w := pixel.Measure("hi")
	if c, r := left(16), left(32); c < (100-w)/2 || c > (100-w)/2+2 || r < 100-w || r > 100-w+2 {
		t.Fatalf("centered ink at %d, right ink at %d, for width %d", c, r, w)
	}
}

func TestTypewriterGlyphs(t *testing.T) {
	clk := newClock()
	s, done := "AAAA\nBBBB", 0
	tt := view(200, 100, clk, func(c *ui.Context) {
		ui.Box(c).Width(float32(2 * (4*pixel.Measure("A") + 8))).Children(func() {
			if text.Typewriter(c, s, text.TypewriterOptions{CPS: 10}).Done {
				done++
			}
		})
	})
	w := 4 * pixel.Measure("A")
	ink := func(y int) int { return count(t, tt.Image(), 0, y, w, pixel.LineHeight, pixel.Ink) }
	clk.add(400 * time.Millisecond)
	tt.Frame()
	if ink(0) == 0 || ink(pixel.LineHeight) != 0 || done != 0 {
		t.Fatalf("at 4 of 8 glyphs: line inks %d, %d, done %d", ink(0), ink(pixel.LineHeight), done)
	}
	clk.add(300 * time.Millisecond)
	tt.Frame()
	if ink(pixel.LineHeight) == 0 || done != 0 {
		t.Fatalf("at 7 of 8 glyphs: second line ink %d, done %d", ink(pixel.LineHeight), done)
	}
	clk.add(100 * time.Millisecond)
	tt.Frame()
	if done != 1 {
		t.Fatal("the eighth glyph did not finish the text")
	}
	// A box-wide line would push the mark out.
	s = "AAAAA"
	tt.Frame()
	if err := tt.Click(s); err != nil {
		t.Fatal(err)
	}
	img := tt.Image()
	if count(t, img, 0, 0, w+8, 3*pixel.LineHeight, pixel.Accent) == 0 {
		clk.add(500 * time.Millisecond)
		tt.Frame()
		img = tt.Image()
	}
	if count(t, img, 0, 0, w+8, 3*pixel.LineHeight, pixel.Accent) == 0 {
		t.Fatal("the done mark never showed")
	}
	if n := count(t, img, w+8, 0, 20, 3*pixel.LineHeight, pixel.Accent); n != 0 {
		t.Fatalf("the done mark left the box: %d dots outside", n)
	}
}

func TestTypewriter(t *testing.T) {
	clk := newClock()
	s, done := "ABCDEFGHIJ", 0
	opts := text.TypewriterOptions{CPS: 10}
	tt := view(400, 200, clk, func(c *ui.Context) {
		if text.Typewriter(c, s, opts).Done {
			done++
		}
	})
	clk.add(500 * time.Millisecond)
	tt.Frame()
	if done != 0 {
		t.Fatal("half the time showed the whole text")
	}
	clk.add(600 * time.Millisecond)
	tt.Frame()
	if done != 1 {
		t.Fatalf("after the whole time done fired %d times", done)
	}
	s = "NEW LINE"
	tt.Frame()
	if err := tt.Click(s); err != nil {
		t.Fatal(err)
	}
	if done != 2 {
		t.Fatal("a click did not complete a new line")
	}
	s = "THIRD"
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if done != 3 {
		t.Fatal("Enter did not complete the line")
	}
	s, opts.Disabled = "FOURTH", true
	tt.Frame()
	_ = tt.Click(s)
	if done != 3 {
		t.Fatal("a disabled typewriter completed on a click")
	}
	s, opts.Disabled = "FIFTH", false
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if done != 4 {
		t.Fatal("Reduce Motion did not show the whole line at once")
	}
}

func TestMarqueeScrollsAndPauses(t *testing.T) {
	clk := newClock()
	tt := view(400, 100, clk, func(c *ui.Context) {
		text.Marquee(c, "HELLO", text.MarqueeOptions{Width: 60})
	})
	// last is the rightmost ink column inside the window.
	last := func() int {
		img := tt.Image()
		if count(t, img, 60, 0, 140, pixel.LineHeight, pixel.Ink) != 0 {
			t.Fatal("the marquee drew outside its window")
		}
		for x := 59; x >= 0; x-- {
			if count(t, img, x, 0, 1, pixel.LineHeight, pixel.Ink) > 0 {
				return x
			}
		}
		return -1
	}
	// Steps come at 15 a second, a dot each.
	r0 := last()
	clk.add(10 * time.Second / 15)
	tt.Frame()
	r1 := last()
	tt.Move(10, 10)
	clk.add(10 * time.Second / 15)
	tt.Frame()
	r2 := last()
	tt.Move(300, 90)
	tt.Frame()
	clk.add(5 * time.Second / 15)
	tt.Frame()
	if r3 := last(); r0 < 0 || r0-r1 != 10 || r2 != r1 || r1-r3 != 5 {
		t.Fatalf("ends at %d, %d (10 steps), %d (hovered), %d (5 more)", r0, r1, r2, r3)
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	clk.add(time.Second)
	tt.Frame()
	if last() != r0 {
		t.Fatal("under Reduce Motion the text is not shown at rest")
	}
}

func TestBlinkUnderReduceMotion(t *testing.T) {
	clk := newClock()
	tt := view(200, 100, clk, func(c *ui.Context) { text.BlinkText(c, "ON", text.BlinkTextOptions{}) })
	ink := func() int { return count(t, tt.Image(), 0, 0, 30, pixel.LineHeight, pixel.Ink) }
	if ink() == 0 {
		t.Fatal("step 0 hid the text")
	}
	clk.add(500 * time.Millisecond)
	tt.Frame()
	if ink() != 0 {
		t.Fatal("step 1 showed the text")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if ink() == 0 {
		t.Fatal("Reduce Motion hid the text")
	}
}

func TestLinkOpens(t *testing.T) {
	opts := text.LinkOptions{}
	clicks := 0
	tt := view(300, 100, newClock(), func(c *ui.Context) {
		if text.Link(c, "Docs", "https://example.com/docs", opts).Clicked {
			clicks++
		}
	})
	if err := tt.Click("Docs"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnter)
	if clicks != 2 || len(tt.OpenedURLs()) != 2 || tt.OpenedURLs()[0] != "https://example.com/docs" {
		t.Fatalf("clicks %d, opened %q", clicks, tt.OpenedURLs())
	}
	opts.Disabled = true
	tt.Frame()
	_ = tt.Click("Docs")
	if len(tt.OpenedURLs()) != 2 {
		t.Fatal("a disabled link opened")
	}
}

func TestSevenSegmentLightsSegments(t *testing.T) {
	tt := view(200, 100, newClock(), func(c *ui.Context) {
		text.SevenSegment(c, "1", text.SevenSegmentOptions{})
	})
	img := tt.Image()
	if count(t, img, 6, 1, 1, 4, pixel.Bad) != 4 || count(t, img, 6, 6, 1, 4, pixel.Bad) != 4 {
		t.Fatal("1 did not light b and c")
	}
	if count(t, img, 1, 0, 5, 1, pixel.Bad) != 0 || count(t, img, 1, 0, 5, 1, pixel.Faint) == 0 {
		t.Fatal("1 lit a, or left no ghost")
	}
}

func TestOdometerAndMicroPanics(t *testing.T) {
	for name, build := range map[string]func(c *ui.Context){
		"negative":    func(c *ui.Context) { text.Odometer(c, -1, text.OdometerOptions{}) },
		"overflow":    func(c *ui.Context) { text.Odometer(c, 1000, text.OdometerOptions{Digits: 3}) },
		"micro rune":  func(c *ui.Context) { text.MicroText(c, "é", text.MicroTextOptions{}) },
		"big scale":   func(c *ui.Context) { text.ShadowText(c, "x", text.ShadowTextOptions{Scale: 5}) },
		"seven rune":  func(c *ui.Context) { text.SevenSegment(c, "x", text.SevenSegmentOptions{}) },
		"paragraph 9": func(c *ui.Context) { text.Paragraph(c, "x", text.ParagraphOptions{Align: 9}) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s did not panic", name)
				}
			}()
			view(200, 100, newClock(), build)
		}()
	}
}

func TestOdometerRollsToValue(t *testing.T) {
	clk := newClock()
	value := 5
	tt := view(300, 100, clk, func(c *ui.Context) { text.Odometer(c, value, text.OdometerOptions{Digits: 2}) })
	digits := func() []pixel.Slot {
		img := tt.Image()
		var out []pixel.Slot
		for y := range pixel.LineHeight + 2 {
			for x := range 20 {
				out = append(out, slot(t, img, x, y))
			}
		}
		return out
	}
	five := digits()
	value = 6
	tt.Frame()
	// One step at 15 fps, S5's fastest.
	clk.add(67 * time.Millisecond)
	tt.Frame()
	mid := digits()
	clk.add(time.Second)
	tt.Frame()
	six := digits()
	if slices.Equal(six, five) || slices.Equal(mid, five) || slices.Equal(mid, six) {
		t.Fatal("the odometer did not roll from 5 through a step to 6")
	}
}

func TestScrambleSettles(t *testing.T) {
	clk := newClock()
	tt := view(300, 100, clk, func(c *ui.Context) {
		ui.Column(c).Children(func() {
			text.ScrambleText(c, "AB", text.ScrambleTextOptions{})
			pixel.Text(c, "AB", pixel.TextOptions{})
		})
	})
	same := func() bool {
		img := tt.Image()
		for y := range pixel.LineHeight {
			for x := range pixel.Measure("AB") {
				if slot(t, img, x, y) != slot(t, img, x, y+pixel.LineHeight) {
					return false
				}
			}
		}
		return true
	}
	if same() {
		t.Fatal("scramble started settled")
	}
	clk.add(400 * time.Millisecond)
	tt.Frame()
	if !same() {
		t.Fatal("after 6 steps the text had not settled")
	}
}

func TestKbdOutlineRainbow(t *testing.T) {
	tt := view(400, 200, newClock(), func(c *ui.Context) {
		ui.Column(c).AlignItems(ui.Start).Children(func() {
			ui.Row(c).Children(func() {
				text.Kbd(c, "A", text.KbdOptions{})
				text.Kbd(c, "B", text.KbdOptions{Pressed: true})
			})
			text.OutlineText(c, "I", text.OutlineTextOptions{})
			text.RainbowText(c, "AB", text.RainbowTextOptions{})
		})
	})
	img := tt.Image()
	a, b := find(t, tt, "A"), find(t, tt, "B")
	if slot(t, img, int(a.X/2)+4, 0) != pixel.Edge || slot(t, img, int(b.X/2)+4, 0) == pixel.Edge {
		t.Fatal("the pressed key did not sink")
	}
	o := find(t, tt, "I")
	ox, oy := int(o.X/2), int(o.Y/2)
	for y := range int(o.H / 2) {
		for x := range int(o.W / 2) {
			if slot(t, img, ox+x, oy+y) != pixel.Lit {
				continue
			}
			for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				if n := slot(t, img, ox+x+d[0], oy+y+d[1]); n != pixel.Lit && n != pixel.Edge {
					t.Fatalf("outline gap beside %d,%d", x, y)
				}
			}
		}
	}
	r := find(t, tt, "AB")
	wa := pixel.Measure("A")
	if count(t, img, int(r.X/2), int(r.Y/2), wa, pixel.LineHeight, pixel.Bad) == 0 || count(t, img, int(r.X/2)+wa, int(r.Y/2), pixel.Measure("B"), pixel.LineHeight, pixel.Warn) == 0 {
		t.Fatal("rainbow did not give A Bad and B Warn")
	}
}

func TestTypewriterStaysDoneAcrossWraps(t *testing.T) {
	clk := newClock()
	width, done := float32(80), 0
	tt := view(400, 200, clk, func(c *ui.Context) {
		ui.Box(c).Width(width).Children(func() {
			if text.Typewriter(c, "AAAA BBBB", text.TypewriterOptions{}).Done {
				done++
			}
		})
		pixel.Text(c, "AAAA BBBB", pixel.TextOptions{})
	})
	_ = tt.Click("AAAA BBBB")
	if r, _ := tt.Find("AAAA BBBB"); done != 1 || r.H != 2*pixel.LineHeight*2 {
		t.Fatalf("narrow: done %d, box %v", done, r)
	}
	width = 200
	tt.Frame()
	img := tt.Image()
	w := pixel.Measure("AAAA BBBB")
	for y := range pixel.LineHeight {
		for x := range w {
			if slot(t, img, x, y) != slot(t, img, x, pixel.LineHeight+y) {
				t.Fatalf("after widening, dot %d,%d differs from the full text", x, y)
			}
		}
	}
	if done != 1 {
		t.Fatalf("done fired %d times", done)
	}
}
