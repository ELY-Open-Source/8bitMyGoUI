package project

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type AssigneePickerOptions struct {
	Disabled bool
}

type AssigneePickerResult struct {
	Element ui.Element
	Changed bool
}

const seatW, seatH = 52, 4 + 32 + 2 + pixel.LineHeight + 2

// AssigneePicker lines up the party; -1 is the empty seat.
func AssigneePicker(c *ui.Context, text string, party []string, chosen *int, o AssigneePickerOptions) AssigneePickerResult {
	k := pixel.K(c)
	if len(party) == 0 || slices.Contains(party, "") || *chosen < -1 || *chosen >= len(party) {
		panic(fmt.Sprintf("project: AssigneePicker has %d members, %d chosen", len(party), *chosen))
	}
	old := *chosen
	var res AssigneePickerResult
	res.Element = ui.Row(c).Gap(k.D(4)).FocusGroup(ui.Horizontal).Role(ui.RoleRadioGroup).Label(text).Children(func() {
		for i := -1; i < len(party); i++ {
			name := k.Msg("nobody")
			if i >= 0 {
				name = party[i]
			}
			// The shell draws the seat; the radio stays bare.
			sel := *chosen
			shell, base, st := pixel.Shell(c, func() ui.Element {
				return ui.RadioBase(c, &sel, i).Label(name).Disabled(o.Disabled)
			})
			// The radio picks before it is told it is off.
			if sel != *chosen && !base.IsDisabled() {
				*chosen = sel
			}
			on := i == *chosen
			shell.Size(k.D(seatW), k.D(seatH)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Frame(0, 0, seatW, seatH, pixel.Raised, pixel.Face)
				if i < 0 {
					g.Veil(10, 4, 32, 32, pixel.Shade, 8)
				} else {
					display.Face(g, 10, 4, name, 2)
				}
				ink := pixel.Ink
				if on {
					ink = pixel.Accent
					g.Rect(0, 0, seatW, 2, pixel.Accent)
					g.Rect(0, seatH-2, seatW, 2, pixel.Accent)
				}
				l := pixel.Wrap(name, seatW-4, 1)[0]
				g.Text((seatW-pixel.Measure(l))/2, 38, l, ink, 1)
				g.Finish(st)
			})
		}
	})
	res.Changed = *chosen != old
	if res.Changed {
		slog.Debug("8bit: assignee", "label", text, "member", *chosen)
	}
	return res
}
