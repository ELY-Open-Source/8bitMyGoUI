package chat

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

type SpeechBubbleOptions struct {
	// Width wraps the text; zero means 140.
	Width int
}

// SpeechBubble shows who said what, when; tail their side.
func SpeechBubble(c *ui.Context, m Message, o SpeechBubbleOptions) ui.Element {
	k := pixel.K(c)
	w := pixel.Or("SpeechBubble Width", o.Width, 140)
	// Yours sit right on Accent, written in AccentInk.
	side, tail, fill, ink := ui.Start, display.Left, pixel.Lit, pixel.Ink
	if m.Mine {
		side, tail, fill, ink = ui.End, display.Right, pixel.Accent, pixel.AccentInk
	}
	return ui.Column(c).AlignSelf(side).AlignItems(side).Gap(k.D(1)).Children(func() {
		pixel.Text(c, line(m.From, w), pixel.TextOptions{Slot: pixel.Faint})
		display.Bubble(c, display.BubbleOptions{Tail: tail, At: 4, Fill: fill}, func() {
			para(c, k, m.Text, w, 0, ink)
		})
		text.MicroText(c, clock(m.At), text.MicroTextOptions{Slot: pixel.Faint})
	})
}
