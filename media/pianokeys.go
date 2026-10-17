package media

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PianoKeysOptions struct {
	Disabled bool
}

type PianoKeysResult struct {
	Element ui.Element
	// Struck are notes begun this frame, like "C#4".
	Struck []string
}

type pianoKey struct{}

// hand is the key under the pointer and last typed.
type hand struct {
	held, typed int
	at          time.Time
}

// A black key sits right of the white it names.
var (
	whites = [8]int{0, 2, 4, 5, 7, 9, 11, 12}
	blacks = [5][2]int{{1, 0}, {3, 1}, {6, 3}, {8, 4}, {10, 5}}
)

// PianoKeys plays C4–C5 by pointer, A–K and WETYU.
func PianoKeys(c *ui.Context, text string, o PianoKeysOptions) PianoKeysResult {
	k := pixel.K(c)
	shell, base, st := pixel.Shell(c, func() ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Disabled(o.Disabled)
	})
	hd := ui.Local(base, pianoKey{}, func() hand { return hand{held: -1, typed: -1} })
	var struck []int
	held := -1
	if !base.IsDisabled() {
		if base.Pressed() {
			x, y, _ := base.PointerPosition()
			held = keyAt(int(x/k.Dot()), int(y/k.Dot()))
		}
		// Sliding along the keys strikes each new one.
		if held >= 0 && held != hd.held {
			struck = append(struck, held)
		}
		if typed := played(base); len(typed) > 0 {
			struck = append(struck, typed...)
			hd.typed, hd.at = typed[len(typed)-1], k.Now()
		}
	}
	hd.held = held
	lit := hd.held
	if k.Since(c, hd.at, 8, 2) < 2 {
		lit = hd.typed
	}
	res := PianoKeysResult{Element: shell}
	for _, s := range struck {
		res.Struck = append(res.Struck, note(s))
	}
	if len(res.Struck) > 0 {
		base.Description(res.Struck[len(res.Struck)-1])
		slog.Debug("8bit: piano keys", "label", text, "struck", res.Struck)
	}
	shell.Size(k.D(80), k.D(40)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		body, _ := dark(k)
		g.Rect(0, 0, 80, 40, pixel.Edge)
		for i, semi := range whites {
			fill := pixel.Lit
			if semi == lit {
				fill = pixel.AccentLit
			}
			g.Rect(i*10, 0, 9, 40, fill)
			g.HLine(i*10, 37, 9, pixel.Shade)
			g.Micro(i*10+3, 30, keyboard[semi].letter, pixel.Faint, 1)
		}
		for _, b := range blacks {
			x, fill := b[1]*10+7, body
			if b[0] == lit {
				fill = pixel.Accent
			}
			g.Rect(x, 0, 6, 24, fill)
			g.Micro(x+2, 17, keyboard[b[0]].letter, pixel.Faint, 1)
		}
		g.Finish(st)
	})
	return res
}

// keyAt is the semitone at x, y; blacks on top.
func keyAt(x, y int) int {
	if x < 0 || y < 0 || x >= 80 || y >= 40 {
		return -1
	}
	if y < 24 {
		for _, b := range blacks {
			if bx := b[1]*10 + 7; x >= bx && x < bx+6 {
				return b[0]
			}
		}
	}
	return whites[x/10]
}
