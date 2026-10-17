package display

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type Status uint8

const (
	Online Status = iota
	Away
	Busy
	Offline
)

var (
	statusSlots = [...]pixel.Slot{Online: pixel.Good, Away: pixel.Warn, Busy: pixel.Bad, Offline: pixel.Faint}
	statusKeys  = [...]string{Online: "online", Away: "away", Busy: "busy", Offline: "offline"}
)

type StatusDotOptions struct {
	// Live makes an online dot blink.
	Live  bool
	Label string
}

// StatusDot is a 7-dot status light.
func StatusDot(c *ui.Context, s Status, o StatusDotOptions) ui.Element {
	if int(s) >= len(statusSlots) {
		panic(fmt.Sprintf("display: unknown Status %d", s))
	}
	k := pixel.K(c)
	label := o.Label
	if label == "" {
		label = k.Msg(statusKeys[s])
	}
	return art(c, k, label, 7, 7, func(g pixel.Canvas) {
		g.Disc(3, 3, 3, pixel.Edge)
		if !(o.Live && s == Online && g.Step(4)%4 >= 2) {
			g.Disc(3, 3, 2, statusSlots[s])
			g.Dot(2, 2, pixel.Lit)
		}
	})
}
