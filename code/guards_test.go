package code_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/code"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestDisabledScrollersHoldStill(t *testing.T) {
	var lines []string
	var logs []code.LogLine
	for i := range 30 {
		lines = append(lines, fmt.Sprint("line ", i))
		logs = append(logs, code.LogLine{Level: code.Info, Text: fmt.Sprint("log ", i)})
	}
	least, input := code.Info, ""
	tt := view(700, 500, newWorld(), func(c *ui.Context) {
		ui.Row(c).Children(func() {
			code.Terminal(c, "Shell", lines, &input, code.TerminalOptions{Disabled: true})
			code.LogViewer(c, "Log", logs, &least, code.LogViewerOptions{Disabled: true})
		})
	})
	// Disabled, it neither follows nor turns with the wheel.
	for _, l := range []string{"line 0", "00:00:00 INFO log 0"} {
		r := find(t, tt, l)
		tt.Scroll(r.X+4, r.Y+4, 0, -300)
		if after := find(t, tt, l); after.Y != r.Y || after.H == 0 {
			t.Fatalf("%q moved from %v to %v", l, r, after)
		}
	}
}

func TestPixelEditorSurvivesAShrinkMidDrag(t *testing.T) {
	art := "....\n....\n"
	tt := view(500, 300, newWorld(), func(c *ui.Context) {
		code.PixelEditor(c, "Art", &art, code.PixelEditorOptions{Zoom: 8})
	})
	r := find(t, tt, "Art")
	tt.Press(r.X+16*3+8, r.Y+16+8)
	art = "..\n"
	tt.Frame()
	tt.Move(r.X+8, r.Y+8)
	tt.Release(r.X+8, r.Y+8)
	if art != "k.\n" {
		t.Fatalf("after the shrink the drag left %q", art)
	}
}

func TestPaletteEditorSwatchesWearTheEdit(t *testing.T) {
	pal := *pixel.SpringDay
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		code.PaletteEditor(c, "Copy", &pal, code.PaletteEditorOptions{})
	})
	click(t, tt, "R")
	tt.Key(0, ui.KeyEnd)
	// Accent, slot 8, is the second row's first swatch.
	x, y, _, _ := dots(find(t, tt, "Copy"))
	want := pal.Colors[pixel.Accent]
	if c := tt.Image().RGBAAt(2*(x+6), 2*(y+14+6)); c.R != want.R || c.G != want.G || c.B != want.B {
		t.Fatalf("the Accent swatch is %v, want the edit %v", c, want)
	}
}

func TestJSONTreeShowsAnEmptyRoot(t *testing.T) {
	for _, doc := range []string{"{}", "[]", "42"} {
		sel := ""
		tt := view(400, 200, newWorld(), func(c *ui.Context) {
			code.JSONTree(c, "Doc", doc, &sel, code.JSONTreeOptions{})
		})
		if r := find(t, tt, "Doc"); r.H != 2*pixel.LineHeight {
			t.Errorf("%s shows %v, want one row", doc, r)
		}
	}
}

func TestCrashScreenReadsGameOver(t *testing.T) {
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		code.CrashScreen(c, "panic: oops", nil, code.CrashScreenOptions{})
	})
	// The group and its title both read GAME OVER.
	n := 0
	for _, s := range tt.Texts() {
		if s == "GAME OVER" {
			n++
		}
	}
	if n != 2 {
		t.Fatalf("GAME OVER read %d times in %q", n, tt.Texts())
	}
	find(t, tt, "CONTINUE?")
}

func TestTerminalShowsAnEmptyLineOnceItRuns(t *testing.T) {
	input := ""
	tt := view(500, 300, newWorld(), func(c *ui.Context) {
		code.Terminal(c, "Shell", nil, &input, code.TerminalOptions{})
	})
	click(t, tt, "COMMAND")
	x, y, w, h := dots(find(t, tt, "COMMAND"))
	empty := count(tt.Image(), x, y, w, h, pixel.Ink)
	tt.Type("help")
	tt.Key(0, ui.KeyEnter)
	if n := count(tt.Image(), x, y, w, h, pixel.Ink); n != empty {
		t.Fatalf("after Enter the field shows %d ink dots, empty is %d", n, empty)
	}
}

func TestCrashScreenStepsAtFourFPS(t *testing.T) {
	wd := newWorld()
	tt := view(600, 300, wd, func(c *ui.Context) {
		code.CrashScreen(c, "panic: oops", nil, code.CrashScreenOptions{})
	})
	// Paint asks for the next step on the frame clock.
	pixel.ResetDue()
	before := tt.Now()
	tt.Frame()
	if d := pixel.Due().Sub(before); d < time.Second/4 || d >= time.Second/2 {
		t.Fatalf("the next blink step is %v away, want a quarter second", d)
	}
}

func TestSpriteSheetKeepsS5Rates(t *testing.T) {
	for _, fps := range []int{2, 30} {
		mustPanic(t, fmt.Sprint(fps, " fps"), func() {
			view(600, 300, newWorld(), func(c *ui.Context) {
				f, p := 0, false
				code.SpriteSheet(c, "S", spin, &f, &p, code.SpriteSheetOptions{FPS: fps})
			})
		})
	}
}
