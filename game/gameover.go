package game

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/navigation"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type GameOverOptions struct {
	Disabled bool
}

type GameOverResult struct {
	Element        ui.Element
	Continue, Quit bool
}

// GameOver dithers GAME OVER in, above CONTINUE and QUIT.
func GameOver(c *ui.Context, o GameOverOptions) GameOverResult {
	k := pixel.K(c)
	s := k.Msg("game_over")
	col := ui.Column(c).Gap(k.D(10)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(s)
	res := GameOverResult{Element: col}
	col.Children(func() {
		e := ui.Box(c).Role(ui.RoleText).Label(s).Size(k.D(3*pixel.Measure(s)), k.D(3*pixel.LineHeight))
		level := 16
		// It starts as it appears, dithering in for 2 s.
		if !k.Reduce() {
			level = k.Since(c, start(k, e, 0).at, 8, 16)
		}
		e.Draw(func(p *ui.Painter, r ui.Rect) { k.Canvas(p, r).TextVeil(0, 0, s, pixel.Bad, level, 3) })
		items := []navigation.MenuItem{{Label: k.Msg("continue_game")}, {Label: k.Msg("quit")}}
		r := navigation.CommandMenu(c, s, items, navigation.CommandMenuOptions{AutoFocus: true, Disabled: o.Disabled})
		res.Continue, res.Quit = r.Chosen == 0, r.Chosen == 1
	})
	if res.Continue || res.Quit {
		slog.Debug("8bit: game over", "continue", res.Continue)
	}
	return res
}
