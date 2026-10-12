package media

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Row is one tracker line; no Note is a rest.
type Row struct {
	// Note is like "C-4" or "C#4".
	Note string
	// Instrument is 1–99 and Volume 0–64 under a note.
	Instrument, Volume int
}

type TrackerOptions struct {
	Disabled bool
}

type TrackerResult struct {
	Element *ui.Element
	Changed bool
}

var noteForm = regexp.MustCompile(`^[A-G][-#][0-9]$`)

// Tracker edits rows: piano letters write notes, Delete rests.
func Tracker(c *ui.Context, text string, rows []Row, cur *int, o TrackerOptions) TrackerResult {
	k := pixel.K(c)
	if *cur < 0 || *cur >= len(rows) {
		panic(fmt.Sprintf("media: Tracker cursor %d of %d rows", *cur, len(rows)))
	}
	for i, r := range rows {
		if r.Note != "" && (!noteForm.MatchString(r.Note) || r.Instrument < 1 || r.Instrument > 99 || r.Volume < 0 || r.Volume > 64) {
			panic(fmt.Sprintf("media: Tracker row %d is %+v", i, r))
		}
	}
	shell, base, st, _ := pixel.Pad(c, text, 1, len(rows), 52, 7, cur, o.Disabled)
	res := TrackerResult{Element: shell}
	if !base.IsDisabled() {
		for _, semi := range played(base) {
			write(rows, *cur, cell(semi))
			res.Changed = true
			slog.Debug("8bit: tracker", "label", text, "row", *cur, "note", rows[*cur].Note)
			*cur = min(len(rows)-1, *cur+1)
		}
		if base.Shortcut(0, ui.KeyDelete) || base.Shortcut(0, ui.KeyBackspace) {
			rows[*cur], res.Changed = Row{}, true
			slog.Debug("8bit: tracker", "label", text, "row", *cur, "note", "")
		}
	}
	base.Description(line(*cur, rows[*cur]))
	at := *cur
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Well)
		for i, row := range rows {
			y := i * 7
			num, ink := pixel.Faint, pixel.Ink
			switch {
			case i == at:
				g.Rect(0, y, g.W, 7, pixel.Accent)
				num, ink = pixel.AccentInk, pixel.AccentInk
			case i%4 == 0:
				g.Rect(0, y, g.W, 7, pixel.Face)
			}
			g.Micro(2, y+1, fmt.Sprintf("%02d", i), num, 1)
			if row.Note == "" {
				g.Micro(14, y+1, "--- .. ..", num, 1)
				continue
			}
			ins, vol := pixel.Info, pixel.Good
			if i == at {
				ins, vol = ink, ink
			}
			g.Micro(14, y+1, row.Note, ink, 1)
			g.Micro(30, y+1, fmt.Sprintf("%02d", row.Instrument), ins, 1)
			g.Micro(42, y+1, fmt.Sprintf("%02d", row.Volume), vol, 1)
		}
		g.Finish(st)
	})
	return res
}

// cell is a semitone in tracker form, like "C-4".
func cell(semi int) string {
	n := note(semi)
	if len(n) == 2 {
		return n[:1] + "-" + n[1:]
	}
	return n
}

// write notes row i; a rest borrows the voice above.
func write(rows []Row, i int, n string) {
	if rows[i].Note == "" {
		rows[i] = Row{Instrument: 1, Volume: 64}
		for j := i - 1; j >= 0; j-- {
			if rows[j].Note != "" {
				rows[i] = rows[j]
				break
			}
		}
	}
	rows[i].Note = n
}

// line reads a row aloud: number, note, instrument, volume.
func line(i int, r Row) string {
	if r.Note == "" {
		return fmt.Sprintf("%02d ---", i)
	}
	return fmt.Sprintf("%02d %s %02d %02d", i, r.Note, r.Instrument, r.Volume)
}
