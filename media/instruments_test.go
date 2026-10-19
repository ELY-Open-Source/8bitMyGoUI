package media_test

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/media"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestPianoKeysStrikeByPointerAndLetters(t *testing.T) {
	wd := newWorld()
	var struck []string
	disabled := false
	tt := view(300, 200, wd, func(c *ui.Context) {
		struck = append(struck, media.PianoKeys(c, "Piano", media.PianoKeysOptions{Disabled: disabled}).Struck...)
	})
	x, y, _, _ := dots(find(t, tt, "Piano"))
	px := func(dx, dy int) (float32, float32) { return float32(2*(x+dx) + 1), float32(2*(y+dy) + 1) }
	key := func(dx int) pixel.Slot { return at(tt.Image(), wd.pal, x+dx, y+10) }
	expect := func(what string, want ...string) {
		t.Helper()
		if !slices.Equal(struck, want) {
			t.Fatalf("%s struck %q, want %q", what, struck, want)
		}
		struck = nil
	}
	tt.Press(px(4, 35))
	expect("pressing C", "C4")
	if key(1) != pixel.AccentLit {
		t.Fatal("a held key does not light")
	}
	// Sliding along strikes each new key once.
	tt.Move(px(14, 35))
	tt.Move(px(15, 35))
	expect("sliding to D", "D4")
	tt.Release(px(15, 35))
	tt.Press(px(9, 10))
	tt.Release(px(9, 10))
	expect("a black key on top", "C#4")
	for _, k := range []struct {
		key  ui.Key
		note string
	}{{ui.KeyA, "C4"}, {ui.KeyW, "C#4"}, {ui.KeyH, "A4"}, {ui.KeyK, "C5"}} {
		tt.Key(0, k.key)
		expect("typing", k.note)
	}
	// A typed key lights for a quarter second.
	if key(71) != pixel.AccentLit {
		t.Fatal("typed K does not light C5")
	}
	later(tt, wd, time.Second/4+time.Millisecond)
	if key(71) != pixel.Lit {
		t.Fatal("typed C5 stays lit")
	}
	disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyA)
	tt.Press(px(4, 35))
	tt.Release(px(4, 35))
	expect("disabled")
}

func TestStepSequencerTogglesAndPlays(t *testing.T) {
	wd := newWorld()
	steps := make([][16]bool, 2)
	steps[0][0] = true
	playing, changed := false, false
	var step int
	tt := view(700, 200, wd, func(c *ui.Context) {
		r := media.StepSequencer(c, "Beat", []string{"KICK", "HAT"}, steps, media.StepSequencerOptions{Playing: playing})
		changed, step = r.Changed || changed, r.Step
	})
	x, y, _, _ := dots(find(t, tt, "Beat"))
	cell := func(col, row int) pixel.Slot { return at(tt.Image(), wd.pal, x+col*12+5, y+row*12+5) }
	tt.ClickAt(float32(2*(x+4*12+6)), float32(2*(y+12+6)))
	if !steps[1][4] || !changed || cell(4, 1) != pixel.Warn {
		t.Fatalf("clicking HAT 5: %v, changed %v; the second beat is Warn", steps[1], changed)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if !steps[1][5] {
		t.Fatal("Right then Enter did not toggle HAT 6")
	}
	led := func(i int) pixel.Slot { return at(tt.Image(), wd.pal, x+i*12+4, y-5) }
	if step != -1 || led(0) != pixel.Shade {
		t.Fatalf("stopped, the playhead is at %d", step)
	}
	// Playing starts at step one; the hit flashes Lit.
	playing = true
	tt.Frame()
	if step != 0 || led(0) != pixel.Bad || at(tt.Image(), wd.pal, x+5, y+5) != pixel.Lit {
		t.Fatalf("playing: step %d, LED %v", step, led(0))
	}
	// At 120 BPM, four steps a beat: 125 ms each.
	later(tt, wd, 125*time.Millisecond)
	if step != 1 || led(1) != pixel.Bad || led(0) != pixel.Shade {
		t.Fatalf("after 125 ms: step %d", step)
	}
	later(tt, wd, 2*time.Second)
	if step != 1 {
		t.Fatalf("sixteen steps on: step %d, want 1 again", step)
	}
	playing = false
	tt.Frame()
	playing = true
	later(tt, wd, time.Millisecond)
	if step != 0 {
		t.Fatalf("play again starts over, at %d", step)
	}
	mustPanic(t, "tracks without rows", func() {
		view(700, 200, newWorld(), func(c *ui.Context) {
			media.StepSequencer(c, "Beat", []string{"A", "B"}, make([][16]bool, 1), media.StepSequencerOptions{})
		})
	})
}

func TestStepSequencerRunsUnderReduceMotion(t *testing.T) {
	wd := newWorld()
	wd.reduce = true
	var step int
	tt := view(700, 200, wd, func(c *ui.Context) {
		step = media.StepSequencer(c, "Beat", []string{"KICK"}, make([][16]bool, 1), media.StepSequencerOptions{Playing: true, Tempo: 60}).Step
	})
	later(tt, wd, 250*time.Millisecond)
	if step != 1 {
		t.Fatalf("at 60 BPM a step is 250 ms; reduced, step %d", step)
	}
}

func TestTrackerWritesNotesAndRests(t *testing.T) {
	wd := newWorld()
	rows := make([]media.Row, 8)
	rows[0] = media.Row{Note: "E-4", Instrument: 3, Volume: 40}
	cur, changed := 0, false
	tt := view(300, 200, wd, func(c *ui.Context) {
		changed = media.Tracker(c, "Pattern", rows, &cur, media.TrackerOptions{}).Changed || changed
	})
	click(t, tt, "Pattern")
	tt.Key(0, ui.KeyHome)
	// S plays D; notes keep the row's voice, step down.
	tt.Key(0, ui.KeyS)
	if rows[0] != (media.Row{Note: "D-4", Instrument: 3, Volume: 40}) || cur != 1 || !changed {
		t.Fatalf("S on row 0: %+v, cursor %d, changed %v", rows[0], cur, changed)
	}
	// A rest borrows the voice above.
	tt.Key(0, ui.KeyW)
	if rows[1] != (media.Row{Note: "C#4", Instrument: 3, Volume: 40}) || cur != 2 {
		t.Fatalf("W on row 1: %+v, cursor %d", rows[1], cur)
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyDelete)
	if rows[1] != (media.Row{}) || cur != 1 {
		t.Fatalf("Delete: %+v, cursor %d", rows[1], cur)
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyK)
	if rows[7].Note != "C-5" || cur != 7 {
		t.Fatalf("K on the last row: %+v, cursor %d", rows[7], cur)
	}
	x, y, _, _ := dots(find(t, tt, "Pattern"))
	img := tt.Image()
	// The cursor row is Accent; every fourth row is Face.
	if at(img, wd.pal, x+1, y+7*7+3) != pixel.Accent || at(img, wd.pal, x+1, y+4*7+3) != pixel.Face ||
		at(img, wd.pal, x+1, y+3*7+3) != pixel.Well || at(img, wd.pal, x+1, y+2*7+3) != pixel.Well {
		t.Fatal("cursor and beat rows are not marked")
	}
	for name, bad := range map[string]media.Row{
		"a note off the scale": {Note: "H-4", Instrument: 1},
		"no instrument":        {Note: "C-4"},
		"too loud":             {Note: "C-4", Instrument: 1, Volume: 65},
	} {
		mustPanic(t, name, func() {
			view(300, 200, newWorld(), func(c *ui.Context) {
				r := 0
				media.Tracker(c, "P", []media.Row{bad}, &r, media.TrackerOptions{})
			})
		})
	}
	mustPanic(t, "a cursor past the rows", func() {
		view(300, 200, newWorld(), func(c *ui.Context) {
			r := 1
			media.Tracker(c, "P", []media.Row{{}}, &r, media.TrackerOptions{})
		})
	})
}

func TestTrackerStartsAnEmptyPatternAtFullVolume(t *testing.T) {
	rows := make([]media.Row, 4)
	cur := 0
	tt := view(300, 200, newWorld(), func(c *ui.Context) {
		media.Tracker(c, "Pattern", rows, &cur, media.TrackerOptions{})
	})
	click(t, tt, "Pattern")
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyA)
	if rows[0] != (media.Row{Note: "C-4", Instrument: 1, Volume: 64}) {
		t.Fatalf("A on an empty pattern: %+v", rows[0])
	}
}

func TestTunerSetsWhatThePointerPoints(t *testing.T) {
	wd := newWorld()
	mhz, changed := 98.3, false
	tt := view(500, 200, wd, func(c *ui.Context) {
		changed = media.Tuner(c, "Radio", &mhz, media.TunerOptions{}).Changed || changed
	})
	r := find(t, tt, "Radio")
	x, y, _, _ := dots(r)
	// 98.3 is 10.8 of 20.5 MHz: dot 84 of 160.
	if at(tt.Image(), wd.pal, x+84, y+10) != pixel.Bad {
		t.Fatal("the needle is not at 98.3")
	}
	// The 100 tick is dot 97; pressing it tunes 100.
	tt.Press(r.X+2*97+1, r.Y+r.H/2)
	tt.Release(r.X+2*97+1, r.Y+r.H/2)
	if mhz != 100 || !changed {
		t.Fatalf("pressing the 100 tick tuned %v", mhz)
	}
	tt.Key(0, ui.KeyRight)
	if mhz != 100.1 {
		t.Fatalf("Right tuned %v, want 100.1 exactly", mhz)
	}
	tt.Key(0, ui.KeyEnd)
	changed = false
	tt.Key(0, ui.KeyRight)
	if mhz != 108 || changed {
		t.Fatalf("Right at the top: %v, changed %v", mhz, changed)
	}
	mustPanic(t, "AM on an FM dial", func() {
		view(500, 200, newWorld(), func(c *ui.Context) {
			f := 1.0
			media.Tuner(c, "Radio", &f, media.TunerOptions{})
		})
	})
}

func TestRecordButtonBlinksAndTimesTheTake(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = reduce
		rec, changed := false, false
		tt := view(300, 100, wd, func(c *ui.Context) {
			changed = media.RecordButton(c, "Record", &rec, media.RecordButtonOptions{}).Changed || changed
		})
		red := func(label string) int {
			x, y, w, h := dots(find(t, tt, label))
			return count(tt.Image(), wd.pal, x, y, w, h, pixel.Bad)
		}
		if red("Record, 0:00") != 37 {
			t.Fatal("idle, the dot is not whole")
		}
		click(t, tt, "Record, 0:00")
		if !rec || !changed {
			t.Fatalf("click: rec %v, changed %v", rec, changed)
		}
		// Half a second on, the dot blinks off unless reduced.
		later(tt, wd, time.Second/2+time.Millisecond)
		if n := red("Record, 0:00"); (n == 0) == reduce {
			t.Fatalf("reduce %v: %d red dots half a second in", reduce, n)
		}
		later(tt, wd, time.Second/2)
		find(t, tt, "Record, 0:01")
		click(t, tt, "Record, 0:01")
		later(tt, wd, 3*time.Second)
		if rec || red("Record, 0:01") != 37 {
			t.Fatal("stopped, the take's time and a whole dot stay")
		}
		// The caller may start it too, in a frame.
		rec = true
		tt.Frame()
		later(tt, wd, time.Second)
		find(t, tt, "Record, 0:01")
	}
}
