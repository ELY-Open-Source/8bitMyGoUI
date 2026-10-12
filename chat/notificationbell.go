package chat

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type NotificationBellOptions struct {
	Disabled bool
}

type NotificationBellResult struct {
	Element *ui.Element
	Clicked bool
}

type bellKey struct{}

// bellState is the last count and how often it rose.
type bellState struct{ last, rises int }

// swings tilt the ringing bell a dot a step.
var swings = [...]int{-1, 1, -1, 1, 0}

// NotificationBell counts news and rings when the count rises.
func NotificationBell(c *ui.Context, text string, n int, o NotificationBellOptions) NotificationBellResult {
	k := pixel.K(c)
	if n < 0 {
		panic(fmt.Sprintf("chat: NotificationBell count %d", n))
	}
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.ButtonBase(c).Label(text).Description(fmt.Sprintf(k.Msg("new_count"), n)).Disabled(o.Disabled).FocusRing(false)
	})
	s := ui.Local(shell, bellKey{}, func() bellState { return bellState{last: n} })
	if n > s.last {
		s.rises++
		slog.Debug("8bit: bell rings", "label", text, "count", n)
	}
	s.last = n
	dx := swings[k.Shot(c, shell, s.rises, 10, len(swings)-1)]
	count := strconv.Itoa(n)
	if n > 99 {
		count = "99+"
	}
	shell.Size(k.D(18+pixel.MicroWidth(count)), k.D(18)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Sprite(2+dx, 5, icon.Bell, 0, pixel.Ink, 1)
		if n > 0 {
			bw := pixel.MicroWidth(count) + 5
			g.Frame(g.W-bw, 0, bw, 9, pixel.Notched, pixel.Bad)
			g.Micro(g.W-bw+3, 2, count, pixel.AccentInk, 1)
		}
		g.Finish(st)
	})
	return NotificationBellResult{Element: shell, Clicked: base.Clicked()}
}
