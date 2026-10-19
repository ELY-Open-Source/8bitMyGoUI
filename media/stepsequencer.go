package media

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type StepSequencerOptions struct {
	// Tempo is 1–225 BPM, four steps a beat; zero: 120.
	Tempo    int
	Playing  bool
	Disabled bool
}

type StepSequencerResult struct {
	Element ui.Element
	Changed bool
	// Step is the playhead, 0–15, or -1 when stopped.
	Step int
}

type seqKey struct{}

// seq is when play began and the grid cursor.
type seq struct {
	since time.Time
	cur   int
}

// groups light steps by beat, like the drum machines.
var groups = [4]pixel.Slot{pixel.Bad, pixel.Warn, pixel.Good, pixel.Info}

// StepSequencer toggles 16 steps per track; a playhead runs along.
func StepSequencer(c *ui.Context, text string, tracks []string, steps [][16]bool, o StepSequencerOptions) StepSequencerResult {
	k := pixel.K(c)
	if len(tracks) == 0 || len(tracks) != len(steps) {
		panic(fmt.Sprintf("media: StepSequencer has %d tracks and %d rows", len(tracks), len(steps)))
	}
	tempo := pixel.Or("StepSequencer Tempo", o.Tempo, 120)
	// S5 caps the playhead at fifteen steps a second.
	if tempo > 225 {
		panic(fmt.Sprintf("media: StepSequencer Tempo %d over 225", tempo))
	}
	res := StepSequencerResult{Step: -1}
	labels := 0
	for _, t := range tracks {
		labels = max(labels, pixel.Measure(t)+4)
	}
	res.Element = ui.Column(c).Gap(k.D(2))
	sq := ui.Local(res.Element, seqKey{}, func() seq { return seq{} })
	sq.cur = min(sq.cur, 16*len(tracks)-1)
	switch {
	case !o.Playing:
		sq.since = time.Time{}
	case sq.since.IsZero():
		sq.since = k.Now()
	}
	// The playhead is the music; Reduce Motion keeps it.
	if o.Playing {
		per := time.Minute / time.Duration(4*tempo)
		n := int(k.Elapsed(&sq.since) / per)
		res.Step = n % 16
		k.At(c, sq.since.Add(time.Duration(n+1)*per))
	}
	res.Element.Children(func() {
		ui.Row(c).Children(func() {
			ui.Box(c).Size(k.D(labels), k.D(4))
			ui.Box(c).Size(k.D(16*12), k.D(4)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				for i := range 16 {
					led := pixel.Shade
					if i == res.Step {
						led = pixel.Bad
					}
					g.Rect(i*12+3, 0, 5, 3, led)
				}
			})
		})
		ui.Row(c).Children(func() {
			ui.Box(c).Size(k.D(labels), k.D(12*len(tracks))).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				for i, t := range tracks {
					g.Text(0, i*12, t, pixel.Ink, 1)
				}
			})
			res.Changed = grid(c, k, text, tracks, steps, &sq.cur, res.Step, o.Disabled)
		})
	})
	return res
}

// grid is the Pad of steps; picking one toggles it.
func grid(c *ui.Context, k pixel.Kit, text string, tracks []string, steps [][16]bool, cur *int, step int, disabled bool) bool {
	shell, base, st, picked := pixel.Pad(c, text, 16, 16*len(tracks), 12, 12, cur, disabled)
	if picked >= 0 {
		on := &steps[picked/16][picked%16]
		*on = !*on
		slog.Debug("8bit: step sequencer", "label", text, "track", tracks[picked/16], "step", picked%16+1, "on", *on)
	}
	state := k.Msg("off")
	if steps[*cur/16][*cur%16] {
		state = k.Msg("on")
	}
	base.Description(fmt.Sprintf("%s, %d, %s", tracks[*cur/16], *cur%16+1, state))
	at := *cur
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i := range 16 * len(tracks) {
			x, y, col := i%16*12, i/16*12, i%16
			if i == at && st.Focus {
				g.Cursor(x, y, 12, 12)
			}
			if !steps[i/16][col] {
				g.Frame(x+1, y+1, 10, 10, pixel.Sunken, pixel.Well)
				continue
			}
			g.Frame(x+1, y+1, 10, 10, pixel.Raised, groups[col/4])
			if col == step {
				g.Rect(x+3, y+3, 6, 6, pixel.Lit)
			}
		}
		g.Finish(st)
	})
	return picked >= 0
}
