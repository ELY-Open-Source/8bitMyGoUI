package chat

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type UnreadDividerOptions struct {
	// Width is in dots; zero means 200.
	Width int
}

// UnreadDivider rules NEW MESSAGES in Bad above the first unread.
func UnreadDivider(c *ui.Context, o UnreadDividerOptions) *ui.Element {
	k := pixel.K(c)
	return rule(c, k, k.Msg("new_messages"), pixel.Or("UnreadDivider Width", o.Width, 200), pixel.Bad, pixel.Bad)
}
