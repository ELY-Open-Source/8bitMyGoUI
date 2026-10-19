package game

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type SkillButtonOptions struct {
	Disabled bool
}

type SkillButtonResult struct {
	Element ui.Element
	// Clicked is a press while the skill is ready.
	Clicked bool
}

// SkillButton shows a skill, its key and a cooldown wipe.
func SkillButton(c *ui.Context, text string, sprite *pixel.Sprite, key string, cooldown float64, o SkillButtonOptions) SkillButtonResult {
	k := pixel.K(c)
	if sprite == nil || !(cooldown >= 0 && cooldown <= 1) || !pixel.MicroOK(key) {
		panic(fmt.Sprintf("game: SkillButton %q needs a sprite, a micro key %q, a cooldown in 0–1 (%v)", text, key, cooldown))
	}
	shell, base, st := pixel.Shell(c, func() ui.Element { return ui.ButtonBase(c).Label(text).Disabled(o.Disabled) })
	res := SkillButtonResult{Element: shell, Clicked: base.Clicked() && cooldown == 0}
	if res.Clicked {
		slog.Debug("8bit: skill", "label", text)
	}
	state := k.Msg("ready")
	if cooldown > 0 {
		state = strconv.Itoa(int(math.Ceil(cooldown*100))) + "%"
	}
	base.Description(state)
	shell.Size(k.D(24), k.D(24)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		style := pixel.Raised
		if st.Press {
			style = pixel.Sunken
		}
		g.Frame(0, 0, 24, 24, style, pixel.Face)
		g.Sprite((24-sprite.W)/2, (24-sprite.H)/2, sprite, 0, pixel.Ink, 1)
		// The wipe clears from twelve, clockwise, as the cooldown ends.
		for y := 2; y < 22; y++ {
			for x := 2; x < 22; x++ {
				a := math.Atan2(float64(x)-11.5, 11.5-float64(y)) / (2 * math.Pi)
				if a < 0 {
					a++
				}
				if a >= 1-cooldown {
					g.Veil(x, y, 1, 1, pixel.Edge, 11)
				}
			}
		}
		kw := pixel.MicroWidth(key)
		g.Rect(1, 1, kw+3, 8, pixel.Lit)
		g.Micro(2, 2, key, pixel.Ink, 1)
		g.Finish(st)
	})
	return res
}
