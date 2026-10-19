package code_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/code"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestLogViewerFiltersByLevel(t *testing.T) {
	wd := newWorld()
	at := wd.clk.t
	lines := []code.LogLine{{At: at, Level: code.Debug, Text: "d"}, {At: at, Level: code.Info, Text: "i"}, {At: at, Level: code.Error, Text: "e"}}
	least, changed := code.Info, false
	tt := view(600, 300, wd, func(c *ui.Context) {
		changed = code.LogViewer(c, "Log", lines, &least, code.LogViewerOptions{}).Changed || changed
	})
	if tt.HasText("10:30:00 DEBUG d") || !tt.HasText("10:30:00 INFO i") {
		t.Fatalf("INFO and up shows %q", tt.Texts())
	}
	x, y, w, h := dots(find(t, tt, "10:30:00 ERROR e"))
	if count(tt.Image(), x, y, w, h, pixel.Bad) == 0 {
		t.Fatal("an ERROR line is not Bad")
	}
	click(t, tt, "ERROR")
	if least != code.Error || !changed || tt.HasText("10:30:00 INFO i") {
		t.Fatalf("ERROR filter: least %v, changed %v, %q", least, changed, tt.Texts())
	}
	mustPanic(t, "a level past ERROR", func() {
		view(600, 300, newWorld(), func(c *ui.Context) {
			l := code.Level(4)
			code.LogViewer(c, "Log", nil, &l, code.LogViewerOptions{})
		})
	})
}

func TestHexViewerWalksBytes(t *testing.T) {
	wd := newWorld()
	data := []byte("HELLO\x00WORLD\x00\x01\x02\x03\x04\xff\xfe\x00\x10")
	cur, changed := 0, false
	tt := view(600, 300, wd, func(c *ui.Context) {
		changed = code.HexViewer(c, "Bytes", data, &cur, code.HexViewerOptions{}).Changed || changed
	})
	r := find(t, tt, "Bytes")
	x, y, _, _ := dots(r)
	img := tt.Image()
	// The cursor's cell is Accent; a zero byte is Faint.
	if at(img, x+1, y+2) != pixel.Accent || count(img, x+5*12, y, 11, 12, pixel.Faint) == 0 || count(img, x+12, y, 11, 12, pixel.Faint) != 0 {
		t.Fatal("cursor or zero bytes are not marked")
	}
	click(t, tt, "Bytes")
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	if cur != 9 || !changed {
		t.Fatalf("Right, Down: cursor %d, changed %v", cur, changed)
	}
	tt.Key(0, ui.KeyEnd)
	if cur != len(data)-1 {
		t.Fatalf("End: cursor %d", cur)
	}
	for name, cols := range map[string]int{"3 a row": 3, "17 a row": 17} {
		mustPanic(t, name, func() {
			view(600, 300, newWorld(), func(c *ui.Context) {
				n := 0
				code.HexViewer(c, "B", data, &n, code.HexViewerOptions{Columns: cols})
			})
		})
	}
}

func TestMemoryMapFillsByUse(t *testing.T) {
	wd := newWorld()
	tt := view(600, 300, wd, func(c *ui.Context) {
		code.MemoryMap(c, "Memory", []code.Bank{{Name: "ROM0", Used: 15, Size: 16}, {Name: "ROM1", Used: 8, Size: 16}, {Name: "WRAM", Used: 3, Size: 16}, {Name: "OAM", Size: 16}}, code.MemoryMapOptions{})
	})
	for i, tc := range []struct {
		label string
		floor pixel.Slot
	}{{"ROM0, 0 of 0 KB, 93%", pixel.Bad}, {"ROM1, 0 of 0 KB, 50%", pixel.Warn}, {"WRAM, 0 of 0 KB, 18%", pixel.Good}, {"OAM, 0 of 0 KB, 0%", pixel.Well}} {
		x, y, _, _ := dots(find(t, tt, tc.label))
		if got := at(tt.Image(), x+20, y+26); got != tc.floor {
			t.Errorf("bank %d floor is %v, want %v", i, got, tc.floor)
		}
	}
	for name, b := range map[string]code.Bank{"overfull": {Name: "A", Used: 2, Size: 1}, "empty": {Name: "A"}, "hanzi": {Name: "内存", Size: 1}} {
		mustPanic(t, name, func() {
			view(600, 300, newWorld(), func(c *ui.Context) { code.MemoryMap(c, "M", []code.Bank{b}, code.MemoryMapOptions{}) })
		})
	}
}

func TestJSONTreeOpensAndRejectsBadJSON(t *testing.T) {
	wd := newWorld()
	sel, toggled := "", false
	tt := view(600, 300, wd, func(c *ui.Context) {
		toggled = code.JSONTree(c, "Doc", `{"a": {"b": 1}, "c": "x"}`, &sel, code.JSONTreeOptions{}).Toggled || toggled
	})
	r := find(t, tt, "Doc")
	if r.H != 2*2*pixel.LineHeight {
		t.Fatalf("closed, the tree is %v tall", r.H)
	}
	click(t, tt, "Doc")
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyRight)
	if r := find(t, tt, "Doc"); !toggled || r.H != 3*2*pixel.LineHeight {
		t.Fatalf("Right did not open a: toggled %v, %v tall", toggled, r.H)
	}
	for name, doc := range map[string]string{"broken": `{"a": }`, "two values": `{} {}`, "a stray brace": `{}}`} {
		mustPanic(t, name, func() {
			view(600, 300, newWorld(), func(c *ui.Context) {
				s := ""
				code.JSONTree(c, "J", doc, &s, code.JSONTreeOptions{})
			})
		})
	}
}

func TestDiffViewSignsAndColors(t *testing.T) {
	wd := newWorld()
	tt := view(600, 300, wd, func(c *ui.Context) {
		code.DiffView(c, "Diff", "a = 1\nb = 2", "a = 1\nb = 3\nc = 4", code.DiffViewOptions{})
	})
	find(t, tt, "  a = 1")
	for label, s := range map[string]pixel.Slot{"- b = 2": pixel.Bad, "+ b = 3": pixel.Good, "+ c = 4": pixel.Good} {
		x, y, w, h := dots(find(t, tt, label))
		if at(tt.Image(), x, y+2) != s || count(tt.Image(), x+32, y, w-32, h, s) == 0 {
			t.Errorf("%q has no %v bar and text", label, s)
		}
	}
}

func TestCrashScreenBlinksAndContinues(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = reduce
		cont := false
		tt := view(600, 300, wd, func(c *ui.Context) {
			cont = code.CrashScreen(c, "panic: oops", []string{"main.main()"}, code.CrashScreenOptions{}).Continue || cont
		})
		find(t, tt, "panic: oops")
		find(t, tt, "  main.main()")
		// The group is named by its title and found first.
		x, y, _, _ := dots(find(t, tt, "GAME OVER"))
		corner := func() pixel.Slot { return at(tt.Image(), x, y) }
		if corner() != pixel.Bad {
			t.Fatal("the frame is not lit at first")
		}
		later(tt, wd, time.Second/2+time.Millisecond)
		if lit := corner() == pixel.Bad; lit != reduce {
			t.Fatalf("reduce %v: lit %v half a second on", reduce, lit)
		}
		click(t, tt, "CONTINUE?")
		if !cont {
			t.Fatal("CONTINUE did not continue")
		}
	}
	mustPanic(t, "no error", func() {
		view(600, 300, newWorld(), func(c *ui.Context) { code.CrashScreen(c, "", nil, code.CrashScreenOptions{}) })
	})
}

func TestFPSMeterKeepsTheNewestAndColorsByTime(t *testing.T) {
	wd := newWorld()
	ms := time.Millisecond
	// Twenty dots hold eight bars: the first two scroll off.
	frames := []time.Duration{20 * ms, 20 * ms, 40 * ms, 10 * ms, 10 * ms, 10 * ms, 10 * ms, 10 * ms, 10 * ms, 20 * ms}
	tt := view(600, 300, wd, func(c *ui.Context) {
		code.FPSMeter(c, "Frames", frames, code.FPSMeterOptions{Width: 20, Height: 30})
	})
	x, y, _, _ := dots(find(t, tt, "Frames"))
	bar := func(i int) pixel.Slot { return at(tt.Image(), x+2+2*i, y+8+30-3) }
	if bar(0) != pixel.Bad || bar(1) != pixel.Good || bar(7) != pixel.Warn {
		t.Fatalf("bars %v %v %v", bar(0), bar(1), bar(7))
	}
	mustPanic(t, "a negative frame", func() {
		view(600, 300, newWorld(), func(c *ui.Context) {
			code.FPSMeter(c, "F", []time.Duration{-1}, code.FPSMeterOptions{})
		})
	})
}

func TestProfilerBarsPlaceSpansByTime(t *testing.T) {
	wd := newWorld()
	ms := time.Millisecond
	root := code.Span{Name: "frame", Length: 10 * ms, Children: []code.Span{{Name: "a", Length: 5 * ms}, {Name: "b", Start: 5 * ms, Length: 5 * ms, Children: []code.Span{{Name: "c", Start: 5 * ms, Length: ms}}}}}
	tt := view(600, 300, wd, func(c *ui.Context) {
		code.ProfilerBars(c, "Frame", root, code.ProfilerBarsOptions{Width: 100})
	})
	x, y, _, h := dots(find(t, tt, "Frame"))
	img := tt.Image()
	row := func(r int) int { return y + r*pixel.LineHeight + 6 }
	if h != 3*pixel.LineHeight || at(img, x+99, row(0)) != pixel.Accent || at(img, x+30, row(1)) != pixel.Info ||
		at(img, x+52, row(2)) != pixel.Good || at(img, x+80, row(2)) == pixel.Good {
		t.Fatal("bars are not where their spans are")
	}
	mustPanic(t, "a child past its parent", func() {
		view(600, 300, newWorld(), func(c *ui.Context) {
			code.ProfilerBars(c, "F", code.Span{Length: ms, Children: []code.Span{{Start: ms / 2, Length: ms}}}, code.ProfilerBarsOptions{})
		})
	})
}
