package feedback

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/text"
)

type ErrorStateOptions struct {
	Disabled bool
}

type ErrorStateResult struct {
	Element  ui.Element
	Continue bool
}

type startKey struct{}

// ErrorState is GAME OVER, a CONTINUE? countdown, and CONTINUE.
func ErrorState(c *ui.Context, message string, o ErrorStateOptions) ErrorStateResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(8)).Padding(k.D(16)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(k.Msg("game_over"))
	start := ui.Local(col, startKey{}, func() time.Time { return k.Now() })
	left := max(0, 9-max(0, int(k.Now().Sub(*start)/time.Second)))
	if left > 0 {
		k.Tick(c, 10)
	}
	res := ErrorStateResult{Element: col}
	col.Children(func() {
		text.ShadowText(c, k.Msg("game_over"), text.ShadowTextOptions{Slot: pixel.Bad, Scale: 3})
		pixel.Text(c, message, pixel.TextOptions{Slot: pixel.Faint, Wrap: true, MaxLines: 6}).Width(k.D(min(pixel.Measure(message), 280)))
		display.Icon(c, icon.Skull, display.IconOptions{Scale: 3})
		text.ShadowText(c, k.Msg("continue")+" "+strconv.Itoa(left), text.ShadowTextOptions{Slot: pixel.Accent, Scale: 2})
		res.Continue = input.Button(c, k.Msg("continue_game"), input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled}).Clicked
	})
	if res.Continue {
		slog.Debug("8bit: continue", "left", left)
	}
	return res
}
