package media_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/media"
)

func TestMediaClocksHoldWhenTheClockStepsBack(t *testing.T) {
	wd := newWorld()
	var r media.StepSequencerResult
	tt := view(400, 300, wd, func(c *ui.Context) {
		r = media.StepSequencer(c, "Sequence", []string{"Kick"}, [][16]bool{{}}, media.StepSequencerOptions{Playing: true})
	})
	wd.clk.t = wd.clk.t.Add(-200 * time.Millisecond)
	tt.Frame()
	if r.Step < 0 {
		t.Errorf("the playhead sits at step %d", r.Step)
	}
	rec := true
	tt = view(400, 300, wd, func(c *ui.Context) { media.RecordButton(c, "Record", &rec, media.RecordButtonOptions{}) })
	wd.clk.t = wd.clk.t.Add(-1200 * time.Millisecond)
	tt.Frame()
	if strings.Contains(strings.Join(tt.Texts(), " "), "-") {
		t.Errorf("a negative take: %q", tt.Texts())
	}
}

func TestStepSequencerTempoStaysUnderS5(t *testing.T) {
	defer func() {
		if r := fmt.Sprint(recover()); !strings.HasPrefix(r, "media: StepSequencer Tempo") {
			t.Fatalf("Tempo 300 gave %q", r)
		}
	}()
	// 225 BPM is fifteen steps a second; 300 is twenty.
	view(400, 300, newWorld(), func(c *ui.Context) {
		media.StepSequencer(c, "Seq", []string{"Kick"}, [][16]bool{{}}, media.StepSequencerOptions{Tempo: 300})
	})
}
