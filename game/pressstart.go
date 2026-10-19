package game

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type PressStartOptions struct {
	Disabled bool
}

type PressStartResult struct {
	Element ui.Element
	Pressed bool
}

// PressStart blinks PRESS START until pressed.
func PressStart(c *ui.Context, o PressStartOptions) PressStartResult {
	k := pixel.K(c)
	s := k.Msg("press_start")
	shell, base, st := pixel.Shell(c, func() ui.Element { return ui.ButtonBase(c).Label(s).Disabled(o.Disabled) })
	res := PressStartResult{Element: shell, Pressed: base.Clicked()}
	if res.Pressed {
		slog.Debug("8bit: press start")
	}
	shell.Size(k.D(pixel.Measure(s)+8), k.D(pixel.LineHeight+4)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		// It blinks at 1 Hz; reduced, it holds still.
		if g.Step(4)%4 < 2 {
			g.Text(4, 2, s, pixel.Ink, 1)
		}
		g.Finish(st)
	})
	return res
}
