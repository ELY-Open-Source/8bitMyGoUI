package code

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type CrashScreenOptions struct {
	// Width is in dots; zero means 260.
	Width    int
	Disabled bool
}

type CrashScreenResult struct {
	Element  ui.Element
	Continue bool
}

// CrashScreen frames an error and its stack; CONTINUE goes on.
func CrashScreen(c *ui.Context, err string, stack []string, o CrashScreenOptions) CrashScreenResult {
	k := pixel.K(c)
	if err == "" {
		panic("code: CrashScreen needs an error")
	}
	body, ink, dim := dark(k)
	var res CrashScreenResult
	res.Element = ui.Column(c).Width(k.D(pixel.Or("CrashScreen Width", o.Width, 260))).Padding(k.D(8)).Gap(k.D(4)).
		Role(ui.RoleGroup).Label(k.Msg("game_over")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, body)
		// The red frame blinks each second, as the guru's did.
		if g.Step(4)%4 < 2 {
			g.Rect(0, 0, g.W, 3, pixel.Bad)
			g.Rect(0, g.H-3, g.W, 3, pixel.Bad)
			g.Rect(0, 0, 3, g.H, pixel.Bad)
			g.Rect(g.W-3, 0, 3, g.H, pixel.Bad)
		}
	}).Children(func() {
		pixel.Text(c, k.Msg("game_over"), pixel.TextOptions{Slot: pixel.Bad})
		pixel.Text(c, err, pixel.TextOptions{Slot: ink, Wrap: true})
		for _, f := range stack {
			pixel.Text(c, "  "+f, pixel.TextOptions{Slot: dim, MaxLines: 1})
		}
		ui.Row(c).Padding(k.D(4), 0, 0, 0).Children(func() {
			res.Continue = input.Button(c, k.Msg("continue"), input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled}).Clicked
		})
	})
	if res.Continue {
		slog.Debug("8bit: crash screen continue", "error", err)
	}
	return res
}
