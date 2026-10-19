package data

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/text"
)

// Event is a dated line on a timeline.
type Event struct {
	Date, Text string
}

type TimelineOptions struct{}

// Timeline strings dated events down a line of markers.
func Timeline(c *ui.Context, events []Event, o TimelineOptions) ui.Element {
	k := pixel.K(c)
	return ui.Column(c).FillWidth().Children(func() {
		for i, ev := range events {
			first, last := i == 0, i == len(events)-1
			ui.Row(c).FillWidth().Gap(k.D(4)).Children(func() {
				ui.Box(c).Width(k.D(9)).Shrink(0).AlignSelf(ui.Stretch).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					top, foot := 0, g.H
					if first {
						top = 5
					}
					if last {
						foot = 5
					}
					g.VLine(4, top, foot-top, pixel.Edge)
					g.Rect(2, 3, 5, 5, pixel.Edge)
					g.Rect(3, 4, 3, 3, pixel.Accent)
				})
				ui.Column(c).Grow(1).Padding(0, 0, k.D(6), 0).Children(func() {
					pixel.Text(c, ev.Date, pixel.TextOptions{Slot: pixel.Faint})
					text.Paragraph(c, ev.Text, text.ParagraphOptions{})
				})
			})
		}
	})
}
