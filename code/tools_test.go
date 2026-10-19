package code_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/code"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestBreakpointsToggleByClickAndEnter(t *testing.T) {
	wd := newWorld()
	src, marks, changed := "a := 1\nb := 2\nc := 3", map[int]bool{}, false
	tt := view(400, 200, wd, func(c *ui.Context) {
		changed = code.Breakpoints(c, "Code", src, marks, code.BreakpointsOptions{}).Changed || changed
	})
	r := find(t, tt, "Code")
	x, y, _, _ := dots(r)
	tt.ClickAt(r.X+20, r.Y+float32(2*pixel.LineHeight)+12)
	if !marks[2] || !changed || at(tt.Image(), x+5, y+pixel.LineHeight+6) != pixel.Bad {
		t.Fatalf("a click on line 2: %v", marks)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	tt.ClickAt(r.X+20, r.Y+float32(2*pixel.LineHeight)+12)
	if !marks[3] || marks[2] || len(marks) != 1 {
		t.Fatalf("Down, Enter, then line 2 again: %v", marks)
	}
	// Fewer lines keep the cursor on the code.
	src = "a := 1"
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if !marks[1] {
		t.Fatalf("after the code shrank: %v", marks)
	}
	mustPanic(t, "no marks", func() {
		view(400, 200, newWorld(), func(c *ui.Context) { code.Breakpoints(c, "C", "x", nil, code.BreakpointsOptions{}) })
	})
}

func TestCommandHistoryReruns(t *testing.T) {
	wd := newWorld()
	sel := 0
	var got code.CommandHistoryResult
	tt := view(400, 200, wd, func(c *ui.Context) {
		if r := code.CommandHistory(c, "History", []string{"ls", "run"}, &sel, code.CommandHistoryOptions{}); r.Changed || r.Rerun {
			got = r
		}
	})
	click(t, tt, "run")
	if sel != 1 || !got.Changed {
		t.Fatalf("clicking run: %d, %+v", sel, got)
	}
	got = code.CommandHistoryResult{}
	tt.Key(0, ui.KeyEnter)
	if !got.Rerun {
		t.Fatal("Enter did not rerun")
	}
}

func TestShortcutListNamesChords(t *testing.T) {
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		code.ShortcutList(c, "Keys", []code.Shortcut{{Action: "Save", Keys: []string{"Ctrl", "S"}}}, code.ShortcutListOptions{})
	})
	find(t, tt, "Save, Ctrl+S")
	mustPanic(t, "an action without keys", func() {
		view(400, 200, newWorld(), func(c *ui.Context) {
			code.ShortcutList(c, "K", []code.Shortcut{{Action: "Nothing"}}, code.ShortcutListOptions{})
		})
	})
}

func TestDebugGridRulesItsContent(t *testing.T) {
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		code.DebugGrid(c, code.DebugGridOptions{Step: 4}, func() {
			ui.Box(c).Size(40, 40).Label("thing")
		})
	})
	x, y, _, _ := dots(find(t, tt, "thing"))
	img := tt.Image()
	if at(img, x, y) != pixel.Accent || at(img, x+4, y+2) != pixel.Accent || at(img, x+1, y+1) == pixel.Accent || at(img, x-10, y-10) != pixel.Shade {
		t.Fatal("no grid every four dots, or no rulers")
	}
	for _, step := range []int{1, 33} {
		mustPanic(t, "a step outside 2–32", func() {
			view(400, 200, newWorld(), func(c *ui.Context) { code.DebugGrid(c, code.DebugGridOptions{Step: step}, func() {}) })
		})
	}
}

func TestPixelEditorPaintsFillsErasesPicks(t *testing.T) {
	wd := newWorld()
	art, changed := "....\n....\n....\n", false
	tt := view(500, 300, wd, func(c *ui.Context) {
		changed = code.PixelEditor(c, "Art", &art, code.PixelEditorOptions{Zoom: 8}).Changed || changed
	})
	r := find(t, tt, "Art")
	cell := func(cx, cy int) (float32, float32) { return r.X + float32(16*cx+8), r.Y + float32(16*cy+8) }
	// A drag paints every cell it crosses, even jumping.
	tt.Press(cell(0, 0))
	tt.Move(cell(3, 0))
	tt.Release(cell(3, 0))
	if art != "kkkk\n....\n....\n" || !changed {
		t.Fatalf("a drag across row 0 left %q", art)
	}
	click(t, tt, "ERASE")
	tt.ClickAt(cell(1, 0))
	if art != "k.kk\n....\n....\n" {
		t.Fatalf("erasing left %q", art)
	}
	// The erased dot touches the clear rows: one run.
	click(t, tt, "FILL")
	tt.ClickAt(cell(0, 2))
	if art != "kkkk\nkkkk\nkkkk\n" {
		t.Fatalf("filling the clear run left %q", art)
	}
	click(t, tt, "COLOR")
	tt.Key(0, ui.KeyHome)
	for range int(pixel.Good) {
		tt.Key(0, ui.KeyRight)
	}
	click(t, tt, "PEN")
	tt.ClickAt(cell(1, 0))
	// Back to Ink; PICK takes Good, the pen lays it.
	click(t, tt, "COLOR")
	tt.Key(0, ui.KeyHome)
	click(t, tt, "PICK")
	tt.ClickAt(cell(1, 0))
	tt.ClickAt(cell(2, 2))
	if art != "kgkk\nkkkk\nkkgk\n" {
		t.Fatalf("pick then paint left %q", art)
	}
	click(t, tt, "COPY")
	if tt.Clipboard() != art {
		t.Fatalf("COPY gave %q", tt.Clipboard())
	}
	for name, a := range map[string]string{"two frames": "k\n\nk", "bad letter": "z"} {
		mustPanic(t, name, func() {
			view(500, 300, newWorld(), func(c *ui.Context) { code.PixelEditor(c, "A", &a, code.PixelEditorOptions{}) })
		})
	}
}

func TestPixelEditorKeysAndGrid(t *testing.T) {
	wd := newWorld()
	art := "...\n...\n"
	tt := view(500, 300, wd, func(c *ui.Context) {
		code.PixelEditor(c, "Art", &art, code.PixelEditorOptions{Zoom: 4})
	})
	// The click paints the middle; keys paint above it.
	click(t, tt, "Art")
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if art != ".k.\n.k.\n" {
		t.Fatalf("Home, Right, Enter left %q", art)
	}
	x, y, _, _ := dots(find(t, tt, "Art"))
	if at(tt.Image(), x+4, y+1) != pixel.Shade {
		t.Fatal("no grid line between cells")
	}
	click(t, tt, "GRID")
	if at(tt.Image(), x+4, y+1) == pixel.Shade {
		t.Fatal("GRID off still draws lines")
	}
}

func TestPaletteEditorStepsChannelsAndPreviews(t *testing.T) {
	wd := newWorld()
	pal, changed := *pixel.SpringDay, false
	tt := view(600, 300, wd, func(c *ui.Context) {
		changed = code.PaletteEditor(c, "Copy", &pal, code.PaletteEditorOptions{}).Changed || changed
	})
	click(t, tt, "R")
	tt.Key(0, ui.KeyEnd)
	if pal.Colors[pixel.Accent].R != 255 || !changed || pixel.SpringDay.Colors[pixel.Accent].R == 255 {
		t.Fatalf("End on R gave %v; the palette itself must not change", pal.Colors[pixel.Accent])
	}
	tt.Key(0, ui.KeyLeft)
	if pal.Colors[pixel.Accent].R != 238 {
		t.Fatalf("a step down from 255 is 238, got %d", pal.Colors[pixel.Accent].R)
	}
	// The preview's title bar wears the edited Accent.
	x, y, _, _ := dots(find(t, tt, "PREVIEW"))
	img, want := tt.Image(), pal.Colors[pixel.Accent]
	if c := img.RGBAAt(2*(x+10), 2*(y+6)); c.R != want.R || c.G != want.G || c.B != want.B {
		t.Fatalf("preview title %v, want %v", c, want)
	}
	mustPanic(t, "no palette", func() {
		view(600, 300, newWorld(), func(c *ui.Context) { code.PaletteEditor(c, "P", nil, code.PaletteEditorOptions{}) })
	})
}

var spin = pixel.MustSprite("g\n\nr\n\nb\n\ny")

func TestSpriteSheetPlaysFromTheChosenFrame(t *testing.T) {
	wd := newWorld()
	frame, playing := 0, false
	tt := view(600, 300, wd, func(c *ui.Context) {
		code.SpriteSheet(c, "Sheet", spin, &frame, &playing, code.SpriteSheetOptions{FPS: 4})
	})
	click(t, tt, "PLAY")
	if !playing || frame != 0 {
		t.Fatalf("PLAY: playing %v at %d", playing, frame)
	}
	later(tt, wd, 250*time.Millisecond)
	later(tt, wd, 250*time.Millisecond)
	if frame != 2 {
		t.Fatalf("half a second at 4 fps: frame %d", frame)
	}
	later(tt, wd, 500*time.Millisecond)
	if frame != 0 {
		t.Fatalf("four frames on, it wraps: %d", frame)
	}
	// Choosing a frame while playing plays on from it.
	click(t, tt, "Sheet")
	tt.Key(0, ui.KeyEnd)
	later(tt, wd, 250*time.Millisecond)
	if frame != 0 {
		t.Fatalf("End then a step: frame %d, want 0 after 3", frame)
	}
	click(t, tt, "PLAY")
	later(tt, wd, time.Second)
	if playing || frame != 0 {
		t.Fatalf("paused, it moved to %d", frame)
	}
	mustPanic(t, "a frame past the sheet", func() {
		view(600, 300, newWorld(), func(c *ui.Context) {
			f, p := 4, false
			code.SpriteSheet(c, "S", spin, &f, &p, code.SpriteSheetOptions{})
		})
	})
}

func TestGlyphMapTurnsPages(t *testing.T) {
	wd := newWorld()
	r := 'A'
	tt := view(600, 300, wd, func(c *ui.Context) {
		code.GlyphMap(c, "Glyphs", &r, code.GlyphMapOptions{})
	})
	find(t, tt, "U+0041 A")
	click(t, tt, "U+0041 A")
	click(t, tt, "Glyphs")
	for _, step := range []struct {
		key  ui.Key
		want rune
	}{{ui.KeyHome, 0}, {ui.KeyPageUp, 0}, {ui.KeyEnd, 127}, {ui.KeyPageDown, 255}, {ui.KeyLeft, 254}, {ui.KeyPageUp, 126}} {
		tt.Key(0, step.key)
		if r != step.want {
			t.Fatalf("after key %v: %U, want %U", step.key, r, step.want)
		}
	}
	r = 0x0378
	tt.Frame()
	find(t, tt, "U+0378")
	mustPanic(t, "past the last rune", func() {
		view(600, 300, newWorld(), func(c *ui.Context) {
			x := rune(0x110000)
			code.GlyphMap(c, "G", &x, code.GlyphMapOptions{})
		})
	})
}

func TestRegisterViewLightsWhatChanged(t *testing.T) {
	wd := newWorld()
	cpu := code.CPU{SP: 0xFD, P: 0x24, PC: 0xC000}
	tt := view(600, 300, wd, func(c *ui.Context) { code.RegisterView(c, "CPU", cpu, code.RegisterViewOptions{}) })
	for _, l := range []string{"A $00", "X $00", "SP $FD", "PC $C000", "P $24"} {
		find(t, tt, l)
	}
	// Hex text fills the first 40 dots of a row.
	accent := func(label string) int {
		x, y, _, h := dots(find(t, tt, label))
		return count(tt.Image(), x, y, 40, h, pixel.Accent)
	}
	if accent("A $00") != 0 {
		t.Fatal("nothing changed yet, but A is lit")
	}
	cpu.A, cpu.P = 0x3F, 0x25
	tt.Frame()
	if accent("A $3F") == 0 || accent("X $00") != 0 {
		t.Fatal("A changed but is not lit, or X is")
	}
	// I and the spare bit are set; C just flipped.
	x, y, _, _ := dots(find(t, tt, "P $25"))
	img := tt.Image()
	flag := func(i int) pixel.Slot { return at(img, x+i*9+2, y+3) }
	if flag(5) != pixel.Good || flag(7) != pixel.Good || flag(0) != pixel.Well || at(img, x+7*9, y+3) != pixel.Accent {
		t.Fatal("flags do not show the P bits")
	}
}
