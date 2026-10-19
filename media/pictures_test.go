package media_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/media"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// step is one frame at 8 fps, plus a hair.
const step = time.Second/8 + time.Microsecond

func TestCassetteReelsTurnOnlyWhilePlaying(t *testing.T) {
	for _, tc := range []struct {
		name            string
		playing, reduce bool
		turns           bool
	}{
		{"stopped", false, false, false}, {"playing", true, false, true}, {"reduced", true, true, false},
	} {
		wd := newWorld()
		wd.reduce = tc.reduce
		tt := view(300, 200, wd, func(c *ui.Context) {
			media.Cassette(c, "Tape", media.CassetteOptions{Title: "Side A", Playing: tc.playing})
		})
		x, y, _, _ := dots(find(t, tt, "Tape"))
		// The left hub's spokes stand + then lean ×.
		plus := func() bool {
			img := tt.Image()
			return at(img, wd.pal, x+24, y+20) == pixel.Ink && at(img, wd.pal, x+24, y+21) == pixel.Lit
		}
		if !plus() {
			t.Fatalf("%s: the hub does not start upright", tc.name)
		}
		later(tt, wd, step)
		if plus() == tc.turns {
			t.Fatalf("%s: turned %v, want %v", tc.name, !plus(), tc.turns)
		}
	}
}

func TestCassetteTapeCrossesOver(t *testing.T) {
	for _, tc := range []struct {
		progress    float64
		left, right pixel.Slot
	}{{0, pixel.Faint, pixel.Ink}, {1, pixel.Ink, pixel.Faint}} {
		wd := newWorld()
		tt := view(300, 200, wd, func(c *ui.Context) {
			media.Cassette(c, "Tape", media.CassetteOptions{Progress: tc.progress})
		})
		x, y, _, _ := dots(find(t, tt, "Tape"))
		img := tt.Image()
		// A full pack reaches five dots out; an empty, three.
		if l, r := at(img, wd.pal, x+18, y+20), at(img, wd.pal, x+45, y+20); l != tc.left || r != tc.right {
			t.Fatalf("progress %v: left %v, right %v; want %v, %v", tc.progress, l, r, tc.left, tc.right)
		}
	}
	for name, o := range map[string]media.CassetteOptions{
		"progress past the end":  {Progress: 1.5},
		"a title in hanzi":       {Title: "混音"},
		"a title past the label": {Title: "THE LONGEST MIXTAPE"},
	} {
		mustPanic(t, name, func() {
			view(300, 200, newWorld(), func(c *ui.Context) { media.Cassette(c, "Tape", o) })
		})
	}
}

func TestCartridgeStickerAndCorner(t *testing.T) {
	wd := newWorld()
	art := pixel.SpriteOf(4, 4, func(x, y int) (pixel.Slot, bool) { return pixel.Good, true })
	tt := view(300, 300, wd, func(c *ui.Context) {
		media.Cartridge(c, "Cart", media.CartridgeOptions{Title: "Pixel Quest", Art: art})
		media.Cartridge(c, "Short", media.CartridgeOptions{Title: "Night Run"})
	})
	x, y, _, _ := dots(find(t, tt, "Cart"))
	img := tt.Image()
	if at(img, wd.pal, x+28, y+24) != pixel.Good {
		t.Fatal("the art does not sit in the sticker's middle")
	}
	// PIXEL and QUEST take a line each.
	if count(img, wd.pal, x+8, y+37, 40, 5, pixel.Ink) == 0 || count(img, wd.pal, x+8, y+43, 40, 5, pixel.Ink) == 0 {
		t.Fatal("the title does not fill two lines")
	}
	if at(img, wd.pal, x+53, y) != pixel.Back || at(img, wd.pal, x+40, y) != pixel.Edge {
		t.Fatal("the top right corner is not cut")
	}
	// One line sits in the middle of the band.
	x, y, _, _ = dots(find(t, tt, "Short"))
	if count(img, wd.pal, x+8, y+37, 40, 3, pixel.Ink) != 0 || count(img, wd.pal, x+8, y+40, 40, 5, pixel.Ink) == 0 {
		t.Fatal("NIGHT RUN is not one middle line")
	}
	for name, o := range map[string]media.CartridgeOptions{
		"art past the sticker": {Art: pixel.SpriteOf(40, 4, func(x, y int) (pixel.Slot, bool) { return pixel.Ink, true })},
		"three lines of title": {Title: "ONE TWO THREE FOUR FIVE"},
		"a word past a line":   {Title: "SUPERCALIFRAGILISTIC"},
	} {
		mustPanic(t, name, func() {
			view(300, 300, newWorld(), func(c *ui.Context) { media.Cartridge(c, "Cart", o) })
		})
	}
}

func TestVinylTurnsAndTheArmSwings(t *testing.T) {
	wd := newWorld()
	playing := false
	tt := view(300, 200, wd, func(c *ui.Context) { media.Vinyl(c, "Record", media.VinylOptions{Playing: playing}) })
	x, y, _, _ := dots(find(t, tt, "Record"))
	lit := func(dx, dy int) bool { return at(tt.Image(), wd.pal, x+dx, y+dy) == pixel.AccentLit }
	head := func(dx, dy int) bool { return at(tt.Image(), wd.pal, x+dx, y+dy) == pixel.Edge }
	if !lit(36, 31) || !head(66, 36) {
		t.Fatal("stopped, the mark is not at three o'clock or the arm not at rest")
	}
	later(tt, wd, step)
	if !lit(36, 31) {
		t.Fatal("a stopped record turned")
	}
	playing = true
	tt.Frame()
	later(tt, wd, step)
	if !lit(35, 35) || lit(36, 31) {
		t.Fatal("playing, the mark did not turn an eighth")
	}
	later(tt, wd, time.Second/2)
	if !head(52, 32) || head(66, 36) {
		t.Fatal("the arm did not swing onto the grooves")
	}
	// Five eighths round, it stops where it is.
	playing = false
	later(tt, wd, step)
	if !lit(27, 27) || lit(36, 31) {
		t.Fatal("a stopped record did not hold its turn")
	}
}

func TestFilmstripHoldsEveryFrame(t *testing.T) {
	wd := newWorld()
	sp := pixel.MustSprite("g.\n.g\n\n.g\ng.\n\nr.\n.r")
	tt := view(300, 200, wd, func(c *ui.Context) {
		media.Filmstrip(c, "Film", sp, media.FilmstripOptions{Scale: 2})
	})
	x, y, w, h := dots(find(t, tt, "Film"))
	// Frames are 8 wide, film 4 more; 20 tall.
	if w != 3*8+4 || h != 6+14 {
		t.Fatalf("film is %d by %d dots", w, h)
	}
	img := tt.Image()
	if at(img, wd.pal, x+5, y+8) != pixel.Good || at(img, wd.pal, x+13, y+8) != pixel.Lit || at(img, wd.pal, x+21, y+8) != pixel.Bad {
		t.Fatal("the frames are not in order")
	}
	if at(img, wd.pal, x+15, y+10) != pixel.Lit {
		t.Fatal("a frame's backing is narrower than its sprite")
	}
	if at(img, wd.pal, x+2, y+2) != pixel.Faint || at(img, wd.pal, x+2, y+h-3) != pixel.Faint || at(img, wd.pal, x+5, y+2) != pixel.Ink {
		t.Fatal("no sprocket holes")
	}
}

func TestWaveformAndVisualizerFollowTheirData(t *testing.T) {
	wd := newWorld()
	tt := view(300, 200, wd, func(c *ui.Context) {
		media.Waveform(c, "Wave", []float64{1, -0.5, 0}, media.WaveformOptions{Height: 21})
		media.Visualizer(c, "Bars", []float64{1, 0.5, 0}, media.VisualizerOptions{Height: 20})
	})
	x, y, w, _ := dots(find(t, tt, "Wave"))
	img := tt.Image()
	// Up ten from the middle, down five, then the middle.
	if w != 5 || at(img, wd.pal, x, y) != pixel.Accent || at(img, wd.pal, x+2, y+15) != pixel.Accent ||
		at(img, wd.pal, x+2, y+16) == pixel.Accent || at(img, wd.pal, x+2, y+9) == pixel.Accent || at(img, wd.pal, x+4, y+10) != pixel.Accent {
		t.Fatal("waveform lines do not follow the samples")
	}
	x, y, _, _ = dots(find(t, tt, "Bars"))
	if at(img, wd.pal, x+1, y) != pixel.AccentLit || at(img, wd.pal, x+1, y+19) != pixel.AccentLit || at(img, wd.pal, x+5, y+5) != pixel.AccentLit ||
		at(img, wd.pal, x+5, y+4) == pixel.AccentLit || count(img, wd.pal, x+8, y, 3, 20, pixel.Magic)+count(img, wd.pal, x+8, y, 3, 20, pixel.AccentLit) != 0 {
		t.Fatal("visualizer bars do not reach their levels")
	}
	mustPanic(t, "a sample past 1", func() {
		view(300, 200, newWorld(), func(c *ui.Context) { media.Waveform(c, "W", []float64{1.5}, media.WaveformOptions{}) })
	})
	mustPanic(t, "a level below 0", func() {
		view(300, 200, newWorld(), func(c *ui.Context) { media.Visualizer(c, "V", []float64{-0.1}, media.VisualizerOptions{}) })
	})
}
