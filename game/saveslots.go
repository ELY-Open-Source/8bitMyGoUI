package game

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Save is a slot; an empty Place reads NEW GAME.
type Save struct {
	Place    string
	Played   time.Duration
	Portrait *pixel.Sprite
}

type SaveSlotsOptions struct {
	Disabled bool
}

type SaveSlotsResult struct {
	Element *ui.Element
	Changed bool
	// Chosen is the slot opened this frame, or -1.
	Chosen int
}

// card is a save slot's size in dots.
const cardW, cardH = 180, 32

// SaveSlots stacks save cards: portrait, place and play time.
func SaveSlots(c *ui.Context, text string, saves []Save, selected *int, o SaveSlotsOptions) SaveSlotsResult {
	k := pixel.K(c)
	if len(saves) == 0 || *selected < 0 || *selected >= len(saves) {
		panic(fmt.Sprintf("game: SaveSlots slot %d of %d", *selected, len(saves)))
	}
	captions := make([]string, len(saves))
	for i, s := range saves {
		if s.Played < 0 || s.Place == "" && (s.Played != 0 || s.Portrait != nil) {
			panic(fmt.Sprintf("game: save %d is empty but played %v", i, s.Played))
		}
		captions[i] = fmt.Sprintf(k.Msg("file"), i+1) + " " + k.Msg("new_game")
		if s.Place != "" {
			captions[i] = fmt.Sprintf(k.Msg("file"), i+1) + " " + s.Place + " " + clockTime(s.Played)
		}
	}
	old := *selected
	shell, base, st, chosen := pixel.Pad(c, text, 1, len(saves), cardW, cardH+3, selected, o.Disabled)
	res := SaveSlotsResult{Element: shell, Changed: *selected != old, Chosen: chosen}
	base.Description(captions[*selected])
	if res.Changed || res.Chosen >= 0 {
		slog.Debug("8bit: save slots", "label", text, "slot", *selected, "chosen", res.Chosen)
	}
	sel, shown, newGame := *selected, append([]Save(nil), saves...), k.Msg("new_game")
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, s := range shown {
			y, fill := i*(cardH+3), pixel.Face
			if i == sel {
				fill = pixel.Lit
			}
			g.Frame(0, y, cardW, cardH, pixel.Raised, fill)
			g.Micro(4, y+3, fmt.Sprintf("%d", i+1), pixel.Faint, 1)
			if s.Place == "" {
				g.Text((cardW-pixel.Measure(newGame))/2, y+(cardH-pixel.LineHeight)/2, newGame, pixel.Faint, 1)
				continue
			}
			g.Frame(12, y+3, 26, 26, pixel.Sunken, pixel.Well)
			if s.Portrait != nil {
				g.Sprite(13, y+4, s.Portrait, 0, pixel.Ink, 2)
			}
			g.Text(44, y+2, pixel.Wrap(s.Place, cardW-48, 1)[0], pixel.Ink, 1)
			g.Text(44, y+cardH-pixel.LineHeight-1, clockTime(s.Played), pixel.Faint, 1)
		}
		if i := sel; i >= 0 {
			g.Rect(0, i*(cardH+3)+cardH-2, cardW, 2, pixel.Accent)
		}
		g.Finish(st)
	})
	return res
}

// clockTime prints play time as H:MM.
func clockTime(d time.Duration) string {
	return fmt.Sprintf("%d:%02d", int(d.Hours()), int(d.Minutes())%60)
}
