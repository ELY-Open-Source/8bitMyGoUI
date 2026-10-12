package game

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/navigation"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TitleScreenOptions struct {
	// Saved turns CONTINUE on.
	Saved    bool
	Disabled bool
}

type TitleScreenResult struct {
	Element *ui.Element
	// Chosen is 0 NEW GAME, 1 CONTINUE, 2 OPTIONS.
	Chosen int
}

// TitleScreen shows the logo and PRESS START, then the menu.
func TitleScreen(c *ui.Context, title string, started *bool, o TitleScreenOptions) TitleScreenResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(12)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(title)
	res := TitleScreenResult{Element: col, Chosen: -1}
	col.Children(func() {
		logo(c, title)
		if !*started {
			if PressStart(c, PressStartOptions{Disabled: o.Disabled}).Pressed {
				*started = true
			}
			return
		}
		items := []navigation.MenuItem{{Label: k.Msg("new_game")}, {Label: k.Msg("continue_game"), Disabled: !o.Saved}, {Label: k.Msg("options")}}
		r := navigation.CommandMenu(c, title, items, navigation.CommandMenuOptions{AutoFocus: true, Disabled: o.Disabled})
		res.Chosen = r.Chosen
		// Escape goes back to PRESS START.
		if r.Canceled {
			*started = false
			slog.Debug("8bit: title screen back", "title", title)
		}
	})
	if res.Chosen >= 0 {
		slog.Debug("8bit: title screen", "title", title, "chosen", res.Chosen)
	}
	return res
}

// logo writes the title big, gold over red, shadowed.
func logo(c *ui.Context, title string) *ui.Element {
	k := pixel.K(c)
	w := 3*pixel.Measure(title) + 2
	return ui.Box(c).Role(ui.RoleText).Label(title).Size(k.D(w), k.D(3*pixel.LineHeight+2)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Text(2, 2, title, pixel.Edge, 3)
		g.TextBlend(0, 0, title, pixel.Warn, pixel.Bad, 3)
	})
}
