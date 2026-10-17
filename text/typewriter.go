package text

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var more = pixel.MustSprite(`
	#####
	.###.
	..#..
`)

type TypewriterOptions struct {
	// CPS is characters per second; zero means 20.
	CPS      int
	Slot     pixel.Slot
	Disabled bool
}

type TypewriterResult struct {
	Element ui.Element
	// Done is true in the frame the whole text shows.
	Done bool
}

type typeKey struct{}

type typeState struct {
	text        string
	start, last int
	done        bool
}

// Typewriter reveals s; a click or Enter shows the rest.
func Typewriter(c *ui.Context, s string, o TypewriterOptions) TypewriterResult {
	k := pixel.K(c)
	cps := pixel.Or("Typewriter CPS", o.CPS, 20)
	if cps > 1000 {
		panic(fmt.Sprintf("text: Typewriter CPS %d over 1000", cps))
	}
	// S5 steps at 4–15 fps; fast typing shows more.
	fps := max(4, min(cps, 15))
	shell, base, st := pixel.Shell(c, func() ui.Element {
		e := ui.Box(c).Role(ui.RoleText).Label(s).Disabled(o.Disabled)
		if !o.Disabled {
			e.Focusable()
		}
		return e
	})
	ts := ui.Local(shell, typeKey{}, func() typeState { return typeState{start: -1} })
	if ts.text != s {
		*ts = typeState{text: s, start: -1}
	}
	// Lines leave room for the mark; glyphs count, not runes.
	w, check := pixel.WrapWidth(c, shell)
	lines := pixel.Wrap(s, max(0, w-8), 0)
	total := 0
	for _, l := range lines {
		total += len([]rune(l))
	}
	res := TypewriterResult{Element: shell}
	if !ts.done {
		ts.last = total
		if !k.Reduce() {
			tick := k.Tick(c, fps)
			if ts.start < 0 {
				ts.start = tick
			}
			ts.last = min(total, (tick-ts.start)*cps/fps)
		}
		if base.Clicked() {
			ts.last = total
		}
		if ts.last == total {
			ts.done, res.Done = true, true
			slog.Debug("8bit: typewriter done", "glyphs", total)
		}
	}
	shown := ts.last
	if ts.done {
		shown = total
	}
	shell.FillWidth().Height(k.D(len(lines) * pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		check(r)
		g := k.Canvas(p, r)
		left, end := shown, 0
		for i, l := range lines {
			lr := []rune(l)
			n := min(len(lr), max(0, left))
			g.Text(0, i*pixel.LineHeight, string(lr[:n]), o.Slot, 1)
			left -= len(lr)
			end = pixel.Measure(string(lr[:n]))
		}
		if shown == total && g.Step(4)%4 < 2 {
			g.Sprite(end+3, len(lines)*pixel.LineHeight-6, more, 0, pixel.Accent, 1)
		}
		if st.Focus {
			g.FocusMark(0, 0, g.W, g.H)
		}
	})
	return res
}
