package account

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Access is how a player wants motion and dots.
type Access struct {
	// Reduce goes to pixel.Settings.Reduce, Dot to Settings.Dot.
	Reduce bool
	Dot    int
}

type AccessibilityOptionsOptions struct {
	Disabled bool
}

type AccessibilityOptionsResult struct {
	Element *ui.Element
	Changed bool
}

// AccessibilityOptions sets Reduce Motion and the dot size, 1–4.
func AccessibilityOptions(c *ui.Context, text string, a *Access, o AccessibilityOptionsOptions) AccessibilityOptionsResult {
	k := pixel.K(c)
	if a.Dot < 1 || a.Dot > 4 {
		panic(fmt.Sprintf("account: AccessibilityOptions dot %d outside 1–4", a.Dot))
	}
	var res AccessibilityOptionsResult
	res.Element = ui.Column(c).Gap(k.D(6)).Role(ui.RoleGroup).Label(text).Children(func() {
		reduce := input.Switch(c, k.Msg("reduce_motion"), &a.Reduce, input.SwitchOptions{Disabled: o.Disabled}).Changed
		dot, sized := a.Dot-1, false
		ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
			mark(c, k, k.Msg("dot_size"), pixel.Ink, 0)
			if input.Segmented(c, k.Msg("dot_size"), []string{"1", "2", "3", "4"}, &dot, input.SegmentedOptions{Disabled: o.Disabled}).Changed {
				a.Dot, sized = dot+1, true
			}
		})
		res.Changed = reduce || sized
	})
	if res.Changed {
		slog.Debug("8bit: accessibility options", "label", text, "reduce", a.Reduce, "dot", a.Dot)
	}
	return res
}
