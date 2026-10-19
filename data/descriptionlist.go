package data

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/text"
)

type DescriptionListOptions struct {
	// TermWidth is the terms' column in dots; zero means 80.
	TermWidth int
}

// DescriptionList sets terms beside their descriptions.
func DescriptionList(c *ui.Context, pairs [][2]string, o DescriptionListOptions) ui.Element {
	k := pixel.K(c)
	tw := pixel.Or("DescriptionList TermWidth", o.TermWidth, 80)
	return ui.Column(c).FillWidth().Gap(k.D(4)).Children(func() {
		for _, p := range pairs {
			ui.Row(c).FillWidth().Gap(k.D(8)).Children(func() {
				ui.Box(c).Width(k.D(tw)).Shrink(0).Children(func() {
					text.Paragraph(c, p[0], text.ParagraphOptions{Slot: pixel.Faint})
				})
				ui.Column(c).Grow(1).Children(func() { text.Paragraph(c, p[1], text.ParagraphOptions{}) })
			})
		}
	})
}
