package account

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PlayerSelectOptions struct {
	Disabled bool
}

type PlayerSelectResult struct {
	Element *ui.Element
	// Chosen is Enter, Space or a click on a row.
	Changed, Chosen bool
}

// PlayerSelect points ▶ at 1 or 2 players in *players.
func PlayerSelect(c *ui.Context, text string, players *int, o PlayerSelectOptions) PlayerSelectResult {
	k := pixel.K(c)
	if *players != 1 && *players != 2 {
		panic(fmt.Sprintf("account: PlayerSelect of %d players", *players))
	}
	rows := []string{k.Msg("one_player"), k.Msg("two_players")}
	// Two-scale text after a two-scale pointer and a gap.
	gap := pixel.Measure("▶")*2 + 6
	w := 0
	for _, r := range rows {
		w = max(w, pixel.Measure(r)*2+gap)
	}
	cur, old := *players-1, *players
	shell, base, st, picked := pixel.Pad(c, text, 1, 2, w, 2*pixel.LineHeight+4, &cur, o.Disabled)
	*players = cur + 1
	base.Description(rows[cur])
	res := PlayerSelectResult{Element: shell, Changed: *players != old, Chosen: picked >= 0}
	at := cur
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, s := range rows {
			y := i * (2*pixel.LineHeight + 4)
			if i == at {
				g.Text(0, y+2, "▶", pixel.Accent, 2)
			}
			g.Text(gap, y+2, s, pixel.Ink, 2)
		}
		g.Finish(st)
	})
	if res.Changed || res.Chosen {
		slog.Debug("8bit: player select", "label", text, "players", *players, "chosen", res.Chosen)
	}
	return res
}
