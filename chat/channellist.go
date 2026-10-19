package chat

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Channel is a # channel; Unread marks new messages.
type Channel struct {
	Name   string
	Unread bool
}

type ChannelListResult struct {
	Element ui.Element
	// Submitted is Enter or a double click on a row.
	Changed, Submitted bool
}

type ChannelListOptions struct {
	// Width and Height are dots; zeros mean 140 by 120.
	Width, Height int
	Disabled      bool
}

// ChannelList lists # channels; unread ones glow with a dot.
func ChannelList(c *ui.Context, text string, chans []Channel, selected *int, o ChannelListOptions) ChannelListResult {
	k := pixel.K(c)
	w := pixel.Or("ChannelList Width", o.Width, 140)
	room := w - 14 - 8 - 10
	name := func(i int) string {
		if chans[i].Unread {
			return chans[i].Name + ", " + k.Msg("new_messages")
		}
		return chans[i].Name
	}
	return ChannelListResult(list(c, k, text, len(chans), selected, w, pixel.Or("ChannelList Height", o.Height, 120), o.Disabled, name, func(i int) {
		ch := chans[i]
		ink := pixel.Faint
		if ch.Unread {
			ink = pixel.Ink
		}
		name := "# " + line(ch.Name, room)
		ui.Box(c).FillWidth().Height(k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Text(0, 0, name, ink, 1)
			if ch.Unread {
				g.Disc(g.W-5, pixel.LineHeight/2, 2, pixel.Accent)
			}
		})
	}))
}
