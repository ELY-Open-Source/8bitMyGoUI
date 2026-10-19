package chat

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Conversation is a ChatList row; Unread counts new lines.
type Conversation struct {
	Name, Last string
	Unread     int
}

type ChatListResult struct {
	Element ui.Element
	// Submitted is Enter or a double click on a row.
	Changed, Submitted bool
}

type ChatListOptions struct {
	// Width and Height are dots; zeros mean 200 by 120.
	Width, Height int
	Disabled      bool
}

// ChatList lists conversations: avatar, last line, unread badge.
func ChatList(c *ui.Context, text string, convs []Conversation, selected *int, o ChatListOptions) ChatListResult {
	k := pixel.K(c)
	w := pixel.Or("ChatList Width", o.Width, 200)
	// Rows lose the ▶ gutter, the avatar and the badge.
	room := w - 14 - 16 - 4 - 22
	name := func(i int) string {
		cv := convs[i]
		if cv.Unread > 0 {
			return fmt.Sprintf("%s: %s, "+k.Msg("unread"), cv.Name, cv.Last, cv.Unread)
		}
		return cv.Name + ": " + cv.Last
	}
	return ChatListResult(list(c, k, text, len(convs), selected, w, pixel.Or("ChatList Height", o.Height, 120), o.Disabled, name, func(i int) {
		cv := convs[i]
		if cv.Unread < 0 {
			panic("chat: ChatList unread count is negative")
		}
		ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).FillWidth().Children(func() {
			display.Avatar(c, cv.Name, display.AvatarOptions{})
			ui.Column(c).Grow(1).Children(func() {
				pixel.Text(c, line(cv.Name, room), pixel.TextOptions{})
				pixel.Text(c, line(cv.Last, room), pixel.TextOptions{Slot: pixel.Faint})
			})
			if cv.Unread > 0 {
				display.Badge(c, cv.Unread, display.BadgeOptions{})
			}
		})
	}))
}
