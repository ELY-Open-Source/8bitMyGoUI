package account

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/datetime"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Session is a signed-in device and when it was seen.
type Session struct {
	Device  string
	Icon    *pixel.Sprite
	Seen    time.Time
	Current bool
}

type SessionListOptions struct {
	// Width is in dots; zero means 240.
	Width    int
	Disabled bool
}

type SessionListResult struct {
	Element ui.Element
	// SignedOut is the session signed out this frame, else -1.
	SignedOut int
}

// SessionList lists signed-in devices, each with SIGN OUT.
func SessionList(c *ui.Context, text string, sessions []Session, o SessionListOptions) SessionListResult {
	k := pixel.K(c)
	// Rows are keyed by device, so focus follows its row.
	seen := map[string]bool{}
	for _, s := range sessions {
		if seen[s.Device] {
			panic(fmt.Sprintf("account: SessionList has %q twice", s.Device))
		}
		seen[s.Device] = true
	}
	res := SessionListResult{SignedOut: -1}
	res.Element = ui.Column(c).Width(k.D(pixel.Or("SessionList Width", o.Width, 240))).Gap(k.D(4)).Role(ui.RoleList).Label(text).Children(func() {
		for i, s := range sessions {
			ui.Row(c).Key(s.Device).FillWidth().Gap(k.D(4)).AlignItems(ui.Center).Role(ui.RoleListItem).Children(func() {
				if s.Icon != nil {
					ui.Box(c).Size(k.D(s.Icon.W), k.D(s.Icon.H)).Draw(func(p *ui.Painter, r ui.Rect) {
						k.Canvas(p, r).Sprite(0, 0, s.Icon, 0, pixel.Ink, 1)
					})
				}
				ui.Column(c).Grow(1).Children(func() {
					pixel.Text(c, s.Device, pixel.TextOptions{MaxLines: 1})
					// This device is here now; others say when.
					if s.Current {
						pixel.Text(c, k.Msg("this_device"), pixel.TextOptions{Slot: pixel.Good})
					} else {
						datetime.RelativeTime(c, s.Seen, datetime.RelativeTimeOptions{})
					}
				})
				if input.Button(c, k.Msg("sign_out"), input.ButtonOptions{Size: input.Small, Variant: input.Danger, Disabled: o.Disabled, Label: k.Msg("sign_out") + " " + s.Device}).Clicked {
					res.SignedOut = i
				}
			})
		}
	})
	if res.SignedOut >= 0 {
		slog.Debug("8bit: session list sign out", "label", text, "session", res.SignedOut)
	}
	return res
}
