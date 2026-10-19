package project

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type LabelPickerOptions struct {
	// Columns of slots; zero means 3.
	Columns  int
	Disabled bool
}

type LabelPickerResult struct {
	Element ui.Element
	Changed bool
}

type labelKey struct{}

const slotW, slotH = 76, pixel.LineHeight + 6

// LabelPicker toggles labels laid out in slots.
func LabelPicker(c *ui.Context, text string, labels []Label, on []bool, o LabelPickerOptions) LabelPickerResult {
	k := pixel.K(c)
	if len(labels) == 0 || len(on) != len(labels) {
		panic(fmt.Sprintf("project: LabelPicker has %d labels and %d marks", len(labels), len(on)))
	}
	for _, l := range labels {
		if l.Name == "" {
			panic("project: LabelPicker has a label with no name")
		}
	}
	cols := pixel.Or("LabelPicker Columns", o.Columns, 3)
	box := ui.Column(c)
	cur := ui.Local(box, labelKey{}, func() int { return 0 })
	// The labels may have shrunk since the cursor moved.
	*cur = min(*cur, len(labels)-1)
	var res LabelPickerResult
	res.Element = box.Children(func() {
		// Arrows walk the slots; Enter, Space or a click toggles.
		shell, base, st, picked := pixel.Pad(c, text, cols, len(labels), slotW, slotH, cur, o.Disabled)
		if picked >= 0 {
			on[picked], res.Changed = !on[picked], true
			slog.Debug("8bit: label picker", "label", text, "slot", picked, "on", on[picked])
		}
		state := k.Msg("off")
		if on[*cur] {
			state = k.Msg("on")
		}
		base.Description(labels[*cur].Name + ", " + state)
		marks, at := append([]bool(nil), on...), *cur
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, l := range labels {
				x, y := i%cols*slotW, i/cols*slotH
				fill, ink := pixel.Well, pixel.Faint
				if marks[i] {
					fill, ink = pixel.Face, pixel.Ink
				}
				g.Frame(x+1, y+1, slotW-2, slotH-2, pixel.Sunken, fill)
				g.Rect(x+5, y+7, 8, 8, pixel.Edge)
				g.Rect(x+6, y+8, 6, 6, l.Slot)
				room := slotW - 22
				if marks[i] {
					room -= 12
					g.Sprite(x+slotW-15, y+5, icon.Check, 0, pixel.Good, 1)
				}
				g.Text(x+16, y+3, pixel.Wrap(l.Name, room, 1)[0], ink, 1)
				if i == at && st.Focus {
					g.Rect(x, y, slotW, 1, pixel.Accent)
					g.Rect(x, y+slotH-1, slotW, 1, pixel.Accent)
					g.Rect(x, y, 1, slotH, pixel.Accent)
					g.Rect(x+slotW-1, y, 1, slotH, pixel.Accent)
				}
			}
			g.Finish(st)
		})
	})
	return res
}
