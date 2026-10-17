package chat

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Reaction is an emoji count; Mine marks your own.
type Reaction struct {
	Emoji string
	Count int
	Mine  bool
}

type ReactionsOptions struct {
	Disabled bool
}

type ReactionsResult struct {
	Element ui.Element
	// Toggled is the chip clicked this frame, else -1.
	Toggled int
}

// Reactions shows emoji chips with counts; a click toggles yours.
func Reactions(c *ui.Context, text string, rs []Reaction, o ReactionsOptions) ReactionsResult {
	k := pixel.K(c)
	res := ReactionsResult{Toggled: -1}
	res.Element = ui.Row(c).Gap(k.D(2)).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(text).Children(func() {
		for i := range rs {
			r := &rs[i]
			if r.Count < 0 || r.Mine && r.Count == 0 {
				panic(fmt.Sprintf("chat: reaction %q counts %d, mine %v", r.Emoji, r.Count, r.Mine))
			}
			art := emoji(r.Emoji)
			// The base turns before it learns it is off.
			v := r.Mine
			shell, base, st := pixel.Shell(c, func() ui.Element {
				return ui.ToggleBase(c, &v).Label(fmt.Sprintf("%s %d", r.Emoji, r.Count)).Disabled(o.Disabled).FocusRing(false)
			})
			if v != r.Mine && !base.IsDisabled() {
				r.Mine = v
				if r.Mine {
					r.Count++
				} else {
					r.Count--
				}
				res.Toggled = i
				slog.Debug("8bit: reaction toggled", "label", text, "emoji", r.Emoji, "mine", r.Mine)
			}
			n, mine := strconv.Itoa(r.Count), r.Mine
			shell.Size(k.D(20+pixel.MicroWidth(n)), k.D(16)).Draw(func(p *ui.Painter, rc ui.Rect) {
				g := k.Canvas(p, rc)
				fill, ink := pixel.Face, pixel.Ink
				if mine {
					fill, ink = pixel.Accent, pixel.AccentInk
				}
				g.Frame(0, 0, g.W, g.H, pixel.Notched, fill)
				g.Sprite(2, 2, art.sp, 0, art.ink, 1)
				g.Micro(16, 6, n, ink, 1)
				g.Finish(st)
			})
		}
	})
	return res
}
