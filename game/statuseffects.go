package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Effect is a buff, or a debuff, with turns left.
type Effect struct {
	Sprite *pixel.Sprite
	Name   string
	// Turns left; zero hides the count.
	Turns  int
	Debuff bool
}

type StatusEffectsOptions struct{}

// StatusEffects rows effect icons: buffs edged Good, debuffs Bad.
func StatusEffects(c *ui.Context, text string, effects []Effect, o StatusEffectsOptions) ui.Element {
	k := pixel.K(c)
	parts := make([]string, len(effects))
	for i, ef := range effects {
		if ef.Sprite == nil || ef.Turns < 0 {
			panic(fmt.Sprintf("game: effect %q has no sprite or %d turns", ef.Name, ef.Turns))
		}
		parts[i] = ef.Name
		if ef.Turns > 0 {
			parts[i] += " " + strconv.Itoa(ef.Turns)
		}
	}
	return ui.Box(c).Role(ui.RoleImage).Label(text).Description(strings.Join(parts, ", ")).Size(k.D(len(effects)*17), k.D(16)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, ef := range effects {
			x, edge := i*17, pixel.Good
			if ef.Debuff {
				edge = pixel.Bad
			}
			g.Frame(x, 0, 16, 16, pixel.Flat, pixel.Well)
			g.Rect(x+1, 1, 14, 1, edge)
			g.Sprite(x+(16-ef.Sprite.W)/2, (16-ef.Sprite.H)/2, ef.Sprite, 0, pixel.Ink, 1)
			if ef.Turns > 0 {
				n := strconv.Itoa(ef.Turns)
				g.Rect(x+15-pixel.MicroWidth(n)-2, 9, pixel.MicroWidth(n)+2, 7, edge)
				g.Micro(x+14-pixel.MicroWidth(n), 10, n, pixel.AccentInk, 1)
			}
		}
	})
}
