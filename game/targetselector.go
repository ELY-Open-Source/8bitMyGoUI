package game

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Enemy is a target; at zero HP, a shadow.
type Enemy struct {
	Name      string
	Sprite    *pixel.Sprite
	HP, MaxHP int
}

type TargetSelectorOptions struct {
	Disabled bool
}

type TargetSelectorResult struct {
	Element ui.Element
	Changed bool
	// Chosen is the target confirmed this frame, or -1.
	Chosen int
}

// foeW and foeH are a target's cell: cursor, sprite, name.
const foeW, foeH = 40, 12 + 24 + 4 + pixel.LineHeight

// TargetSelector rows enemies under a bobbing ▼; arrows choose.
func TargetSelector(c *ui.Context, text string, enemies []Enemy, selected *int, o TargetSelectorOptions) TargetSelectorResult {
	k := pixel.K(c)
	if len(enemies) == 0 || *selected < 0 || *selected >= len(enemies) {
		panic(fmt.Sprintf("game: TargetSelector target %d of %d", *selected, len(enemies)))
	}
	for _, en := range enemies {
		if en.Sprite == nil {
			panic(fmt.Sprintf("game: enemy %q has no sprite", en.Name))
		}
		amount("enemy "+en.Name+" HP", en.HP, en.MaxHP)
	}
	old := *selected
	shell, base, st, chosen := pixel.Pad(c, text, len(enemies), len(enemies), foeW, foeH, selected, o.Disabled)
	res := TargetSelectorResult{Element: shell, Changed: *selected != old, Chosen: chosen}
	en := enemies[*selected]
	base.Description(en.Name + " " + ratio(en.HP, en.MaxHP))
	if res.Changed || res.Chosen >= 0 {
		slog.Debug("8bit: target", "label", text, "enemy", en.Name, "chosen", res.Chosen)
	}
	sel, shown := *selected, append([]Enemy(nil), enemies...)
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, en := range shown {
			x := i * foeW
			sx := x + (foeW-2*en.Sprite.W)/2
			if en.HP == 0 {
				g.Silhouette(sx, 12, en.Sprite, 0, pixel.Shade, 2)
			} else {
				g.Sprite(sx, 12, en.Sprite, 0, pixel.Ink, 2)
			}
			ink := pixel.Ink
			if en.HP == 0 {
				ink = pixel.Faint
			}
			name := pixel.Wrap(en.Name, foeW-2, 1)[0]
			g.Text(x+(foeW-pixel.Measure(name))/2, 12+24+4, name, ink, 1)
			// The cursor bobs a dot; reduced, it holds.
			if i == sel {
				g.Text(x+(foeW-pixel.Measure("▼"))/2, g.Step(4)%2-4, "▼", pixel.Accent, 1)
			}
		}
		g.Finish(st)
	})
	return res
}
