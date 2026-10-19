package account

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type DifficultySelectOptions struct {
	Disabled bool
}

type DifficultySelectResult struct {
	Element ui.Element
	Changed bool
}

var levels = []string{"easy", "normal", "hard", "nightmare"}

// DifficultySelect picks EASY to NIGHTMARE, a skull a step.
func DifficultySelect(c *ui.Context, text string, chosen *int, o DifficultySelectOptions) DifficultySelectResult {
	k := pixel.K(c)
	choice("DifficultySelect", chosen, len(levels))
	old := *chosen
	// Arrows choose at once; the cursor is the choice.
	shell, base, st, _ := pixel.Pad(c, text, 1, len(levels), 140, 16, chosen, o.Disabled)
	base.Description(fmt.Sprintf("%s, %d", k.Msg(levels[*chosen]), *chosen+1))
	at := *chosen
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, l := range levels {
			y, ink := i*16, pixel.Ink
			if i == len(levels)-1 {
				ink = pixel.Bad
			}
			if i == at {
				g.Text(2, y+2, "▶", pixel.Accent, 1)
			}
			g.Text(14, y+2, k.Msg(l), ink, 1)
			for s := range i + 1 {
				g.Sprite(140-14*(s+1), y+2, icon.Skull, 0, ink, 1)
			}
		}
		g.Finish(st)
	})
	res := DifficultySelectResult{Element: shell, Changed: *chosen != old}
	if res.Changed {
		slog.Debug("8bit: difficulty", "label", text, "level", levels[*chosen])
	}
	return res
}
