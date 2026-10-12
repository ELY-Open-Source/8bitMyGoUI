package input

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CheatCodeOptions struct {
	// Sequence is the code; nil means ↑↑↓↓←→←→BA.
	Sequence []ui.Key
	Disabled bool
}

type CheatCodeResult struct {
	Element  *ui.Element
	Unlocked bool
}

var konami = []ui.Key{ui.KeyUp, ui.KeyUp, ui.KeyDown, ui.KeyDown, ui.KeyLeft, ui.KeyRight, ui.KeyLeft, ui.KeyRight, ui.KeyB, ui.KeyA}

type cheatKey struct{}

// cheat is the key trail and an unreported win.
type cheat struct {
	trail      []ui.Key
	won, fresh bool
}

// CheatCode watches keys for a sequence and shows the trail.
func CheatCode(c *ui.Context, text string, o CheatCodeOptions) CheatCodeResult {
	k := pixel.K(c)
	seq := o.Sequence
	if seq == nil {
		seq = konami
	}
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Disabled(o.Disabled)
	})
	s := ui.Local(base, cheatKey{}, func() cheat { return cheat{} })
	// Keys come in order here, however fast they are typed.
	base.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputKeyDown || ev.Mods != 0 {
			return false
		}
		s.won, s.trail = false, append(s.trail, ev.Key)
		s.trail = s.trail[max(0, len(s.trail)-len(seq)):]
		if slices.Equal(s.trail, seq) {
			s.won, s.fresh, s.trail = true, true, nil
		}
		// Other keys break the code yet pass on, like Tab.
		return slices.Contains(seq, ev.Key)
	})
	res := CheatCodeResult{Element: shell, Unlocked: s.fresh}
	s.fresh = false
	if res.Unlocked {
		slog.Debug("8bit: cheat unlocked", "label", text)
	}
	got := matched(s.trail, seq)
	base.Description(fmt.Sprintf("%d/%d", got, len(seq)))
	trail, won := s.trail, s.won
	shell.Size(k.D(len(seq)*10+4), k.D(14+pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, 14, pixel.Sunken, pixel.Well)
		// With no trail yet, the code itself waits in Shade.
		keys, ink := trail, pixel.Faint
		switch {
		case won:
			keys, ink = seq, pixel.Good
			g.Text(2, 14, k.Msg("lives"), pixel.Good, 1)
		case len(trail) == 0:
			keys, ink = seq, pixel.Shade
		}
		for i, key := range keys {
			kink := ink
			if len(trail) > 0 && i >= len(trail)-got {
				kink = pixel.Accent
			}
			glyph(g, 2+i*10, 4, key, kink)
		}
		g.Finish(st)
	})
	return res
}

// matched is the longest end of trail that starts seq.
func matched(trail, seq []ui.Key) int {
	for n := min(len(trail), len(seq)); n > 0; n-- {
		if slices.Equal(trail[len(trail)-n:], seq[:n]) {
			return n
		}
	}
	return 0
}

// glyph draws a key as an arrow or micro letter.
func glyph(g pixel.Canvas, x, y int, key ui.Key, ink pixel.Slot) {
	switch {
	case key == ui.KeyUp:
		g.Sprite(x+2, y+1, tinyUp, 0, ink, 1)
	case key == ui.KeyDown:
		g.Sprite(x+2, y+1, tinyDown, 0, ink, 1)
	case key == ui.KeyLeft:
		g.Sprite(x+3, y, tinyLeft, 0, ink, 1)
	case key == ui.KeyRight:
		g.Sprite(x+3, y, tinyRight, 0, ink, 1)
	case key >= ui.KeyA && key <= ui.KeyZ:
		g.Micro(x+2, y, string(rune('A'+key-ui.KeyA)), ink, 1)
	case key >= ui.Key0 && key <= ui.Key9:
		g.Micro(x+2, y, string(rune('0'+key-ui.Key0)), ink, 1)
	default:
		g.Rect(x+2, y+1, 3, 3, ink)
	}
}
