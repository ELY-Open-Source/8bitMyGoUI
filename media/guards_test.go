package media_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/media"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// inertView builds under Inert while *inert holds.
func inertView(w, h int, inert *bool, build func(c *ui.Context)) *ui.Tester {
	clk := newWorld().clk
	return ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: clk.now, Inert: *inert})
		build(c)
	}, w, h)
}

func TestScrubberKeepsAFractionalEnd(t *testing.T) {
	pos, total := time.Duration(0), 1500*time.Millisecond
	tt := view(400, 100, newWorld(), func(c *ui.Context) {
		media.Scrubber(c, "Seek", &pos, total, media.ScrubberOptions{})
	})
	click(t, tt, "Seek")
	tt.Key(0, ui.KeyEnd)
	tt.Frame()
	if pos != total {
		t.Fatalf("End on a 1.5 s clip: at %v", pos)
	}
}

func TestControlsHoldStillUnderInert(t *testing.T) {
	inert := false
	level, muted := 1, false
	var struck []string
	tt := inertView(500, 300, &inert, func(c *ui.Context) {
		ui.Column(c).Children(func() {
			media.Volume(c, "Volume", &level, &muted, media.VolumeOptions{})
			struck = append(struck, media.PianoKeys(c, "Piano", media.PianoKeysOptions{}).Struck...)
		})
	})
	// A press that began before Inert does not drag on.
	r := find(t, tt, "Volume")
	tt.Press(r.X+2, r.Y+r.H/2)
	inert = true
	tt.Frame()
	tt.Move(r.X+r.W-1, r.Y+r.H/2)
	tt.Release(r.X+r.W-1, r.Y+r.H/2)
	if level != 0 {
		t.Fatalf("an inert Volume dragged to %d", level)
	}
	p := find(t, tt, "Piano")
	tt.Press(p.X+8, p.Y+p.H-6)
	tt.Release(p.X+8, p.Y+p.H-6)
	if len(struck) != 0 {
		t.Fatalf("an inert piano struck %q", struck)
	}
}

func TestStepSequencerIgnoresKeysOnceDisabled(t *testing.T) {
	steps := make([][16]bool, 1)
	disabled := false
	tt := view(700, 200, newWorld(), func(c *ui.Context) {
		media.StepSequencer(c, "Beat", []string{"KICK"}, steps, media.StepSequencerOptions{Disabled: disabled})
	})
	click(t, tt, "Beat")
	steps[0] = [16]bool{}
	// Disabled lands in the frame that reads Enter.
	disabled = true
	tt.Key(0, ui.KeyEnter)
	if steps[0] != ([16]bool{}) {
		t.Fatalf("a disabled sequencer toggled %v", steps[0])
	}
}

func TestStepSequencerKeepsItsCursorAsTracksGo(t *testing.T) {
	tracks, steps := []string{"KICK", "HAT"}, make([][16]bool, 2)
	tt := view(700, 200, newWorld(), func(c *ui.Context) {
		media.StepSequencer(c, "Beat", tracks, steps, media.StepSequencerOptions{})
	})
	click(t, tt, "Beat")
	tt.Key(0, ui.KeyEnd)
	tracks, steps = tracks[:1], steps[:1]
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if !steps[0][15] {
		t.Fatalf("the cursor did not settle on the last step left: %v", steps[0])
	}
}

func TestImageCompareRejectsASplitOutside(t *testing.T) {
	for _, s := range []float64{-0.1, 2} {
		mustPanic(t, "split outside 0–1", func() {
			view(300, 200, newWorld(), func(c *ui.Context) {
				media.ImageCompare(c, "C", solid(2, 2, pixel.Ink), solid(2, 2, pixel.Ink), &s, media.ImageCompareOptions{})
			})
		})
	}
}

func TestComposedControlsLogTheirOwnEvents(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(old)
	pos, sel, rec := time.Duration(0), 0, false
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		ui.Column(c).Children(func() {
			media.Scrubber(c, "Seek", &pos, time.Minute, media.ScrubberOptions{})
			media.Playlist(c, "List", songs, 0, &sel, media.PlaylistOptions{Height: 60})
			media.RecordButton(c, "Record", &rec, media.RecordButtonOptions{})
		})
	})
	click(t, tt, "Seek")
	click(t, tt, "Two, Bob")
	tt.Key(0, ui.KeyEnter)
	click(t, tt, "Record, 0:00")
	for _, want := range []string{"8bit: scrubber", "8bit: playlist", "submitted=true", "8bit: record button"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("no %q in the log:\n%s", want, buf.String())
		}
	}
}

func TestTrackerHoldsStillUnderInert(t *testing.T) {
	inert := false
	rows, cur := make([]media.Row, 4), 0
	rows[1] = media.Row{Note: "E-4", Instrument: 1, Volume: 64}
	tt := inertView(300, 200, &inert, func(c *ui.Context) {
		media.Tracker(c, "Pattern", rows, &cur, media.TrackerOptions{})
	})
	click(t, tt, "Pattern")
	tt.Key(0, ui.KeyHome)
	// Inert lands in the frame that reads the keys.
	inert = true
	tt.Key(0, ui.KeyA)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDelete)
	if rows[0].Note != "" || rows[1].Note != "E-4" || cur != 0 {
		t.Fatalf("an inert tracker took keys: %+v, cursor %d", rows, cur)
	}
}

func TestScrubberStepsAShortClip(t *testing.T) {
	pos, n := time.Duration(0), 0
	tt := view(400, 100, newWorld(), func(c *ui.Context) {
		if media.Scrubber(c, "Seek", &pos, 5*time.Second, media.ScrubberOptions{}).Changed {
			n++
		}
	})
	click(t, tt, "Seek")
	tt.Key(0, ui.KeyHome)
	n = 0
	tt.Key(0, ui.KeyRight)
	if pos != 250*time.Millisecond || n != 1 {
		t.Fatalf("Right on 5 s: at %v after %d changes", pos, n)
	}
	for range 19 {
		tt.Key(0, ui.KeyRight)
	}
	if pos != 5*time.Second || n != 20 {
		t.Fatalf("twenty Rights: at %v after %d changes", pos, n)
	}
}

func TestPlaylistAndVideoFrameRejectBadSizes(t *testing.T) {
	mustPanic(t, "a track of negative length", func() {
		view(500, 300, newWorld(), func(c *ui.Context) {
			s := 0
			media.Playlist(c, "List", []media.Track{{Title: "Back", Artist: "Time", Length: -time.Second}}, -1, &s, media.PlaylistOptions{})
		})
	})
	mustPanic(t, "a screen too narrow for its scrubber", func() {
		view(500, 300, newWorld(), func(c *ui.Context) {
			p, at := false, time.Duration(0)
			media.VideoFrame(c, "Clip", &p, &at, time.Second, media.VideoFrameOptions{Width: 100}, func() {})
		})
	})
}
