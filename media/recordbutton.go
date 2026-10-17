package media

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type RecordButtonOptions struct {
	Disabled bool
}

type RecordButtonResult struct {
	Element ui.Element
	Changed bool
}

type takeKey struct{}

// take is when recording began and the last take's length.
type take struct {
	since time.Time
	took  time.Duration
}

var (
	recDot = pixel.MustSprite(`
		..rrr..
		.rrrrr.
		rrrrrrr
		rrrrrrr
		rrrrrrr
		.rrrrr.
		..rrr..
	`)
	recOff = pixel.MustSprite(`
		..sss..
		.s...s.
		s.....s
		s.....s
		s.....s
		.s...s.
		..sss..
	`)
)

// RecordButton toggles *rec; the dot blinks, time runs.
func RecordButton(c *ui.Context, text string, rec *bool, o RecordButtonOptions) RecordButtonResult {
	k := pixel.K(c)
	now := k.Now()
	row := ui.Row(c)
	tk := ui.Local(row, takeKey{}, func() take { return take{} })
	// Follow *rec, whoever flips it; a click builds again.
	switch {
	case *rec && tk.since.IsZero():
		tk.since, tk.took = now, 0
	case !*rec && !tk.since.IsZero():
		tk.since, tk.took = time.Time{}, k.Elapsed(&tk.since)
	}
	took, dot := tk.took, recDot
	if *rec {
		took = k.Elapsed(&tk.since)
		k.At(c, tk.since.Add(took.Truncate(time.Second)+time.Second))
		if !k.Reduce() && k.Tick(c, 4)%4 >= 2 {
			dot = recOff
		}
	}
	res := RecordButtonResult{Element: row}
	row.Children(func() {
		res.Changed = input.Toggle(c, k.Msg("rec")+" "+mmss(took), rec, input.ToggleOptions{Icon: dot, Label: text + ", " + mmss(took), Disabled: o.Disabled}).Changed
	})
	if res.Changed {
		slog.Debug("8bit: record button", "label", text, "rec", *rec, "time", mmss(took))
	}
	return res
}
