package account

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type AccountSwitcherOptions struct {
	Disabled bool
}

type AccountSwitcherResult struct {
	Element ui.Element
	Changed bool
}

// AccountSwitcher shows slots P1 to P4; empty names are free.
func AccountSwitcher(c *ui.Context, text string, names [4]string, current *int, o AccountSwitcherOptions) AccountSwitcherResult {
	k := pixel.K(c)
	choice("AccountSwitcher", current, len(names))
	var res AccountSwitcherResult
	res.Element = ui.Row(c).Gap(k.D(4)).FocusGroup(ui.Horizontal).Role(ui.RoleRadioGroup).Label(text).Children(func() {
		for i, n := range names {
			say := n
			if n == "" {
				say = k.Msg("free_slot")
			}
			tag := fmt.Sprintf("P%d", i+1)
			label := tag + ", " + say
			if i == *current {
				label += ", " + k.Msg("in_use")
			}
			// The shell draws the card; the base stays bare.
			shell, base, st := pixel.Shell(c, func() ui.Element {
				return ui.ButtonBase(c).Role(ui.RoleRadio).Label(label).Disabled(o.Disabled)
			})
			if base.Clicked() && *current != i {
				*current, res.Changed = i, true
			}
			on, color := i == *current, players[i]
			shell.Size(k.D(64), k.D(62)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				edge := pixel.Shade
				if on {
					edge = color
				}
				g.Frame(0, 0, 64, 62, pixel.Raised, pixel.Face)
				g.Rect(0, 0, 64, 2, edge)
				g.Rect(0, 60, 64, 2, edge)
				g.Rect(2, 2, 60, 10, color)
				g.Micro(32-pixel.MicroWidth(tag)/2, 4, tag, pixel.AccentInk, 1)
				if n == "" {
					g.Checker(16, 14, 32, 32, pixel.Shade)
				} else {
					display.Face(g, 16, 14, n, 2)
				}
				name := pixel.Wrap(say, 60, 1)[0]
				g.Text(32-pixel.Measure(name)/2, 47, name, pixel.Ink, 1)
				g.Finish(st)
			})
		}
	})
	if res.Changed {
		slog.Debug("8bit: account switcher", "label", text, "slot", *current+1)
	}
	return res
}
