package account

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Mix is music, effects and voice volume, 0 to 1.
type Mix struct {
	BGM, SE, Voice float64
}

type VolumeMixerOptions struct {
	Disabled bool
}

type VolumeMixerResult struct {
	Element ui.Element
	Changed bool
}

// VolumeMixer sets BGM, SE and VOICE on faders.
func VolumeMixer(c *ui.Context, text string, mix *Mix, o VolumeMixerOptions) VolumeMixerResult {
	k := pixel.K(c)
	var res VolumeMixerResult
	res.Element = ui.Row(c).Gap(k.D(10)).Role(ui.RoleGroup).Label(text).Children(func() {
		for _, ch := range []struct {
			key string
			v   *float64
		}{{"bgm", &mix.BGM}, {"se", &mix.SE}, {"voice", &mix.Voice}} {
			if !(*ch.v >= 0 && *ch.v <= 1) {
				panic(fmt.Sprintf("account: VolumeMixer %s at %v", ch.key, *ch.v))
			}
			ui.Column(c).Gap(k.D(2)).AlignItems(ui.Center).Children(func() {
				name := k.Msg(ch.key)
				if input.Fader(c, name, ch.v, input.FaderOptions{Height: 56, Disabled: o.Disabled}).Changed {
					res.Changed = true
					slog.Debug("8bit: volume mixer", "label", text, "channel", ch.key, "level", *ch.v)
				}
				mark(c, k, name, pixel.Faint, 0)
				caption(c, k, fmt.Sprintf("%d%%", int(*ch.v*100+0.5)), pixel.Ink)
			})
		}
	})
	return res
}
