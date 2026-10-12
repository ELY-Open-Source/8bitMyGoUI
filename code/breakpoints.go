package code

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type BreakpointsOptions struct {
	Disabled bool
}

type BreakpointsResult struct {
	Element *ui.Element
	Changed bool
}

type lineKey struct{}

// Breakpoints shows code; a click or Enter toggles a dot.
func Breakpoints(c *ui.Context, text, src string, marks map[int]bool, o BreakpointsOptions) BreakpointsResult {
	k := pixel.K(c)
	if marks == nil {
		panic("code: Breakpoints needs a map of marks")
	}
	src = tidy(src)
	lines, slots := strings.Split(src, "\n"), paint(src)
	gutter := 10 + pixel.Measure(strconv.Itoa(len(lines))) + 6
	w := 0
	for _, l := range lines {
		w = max(w, gutter+pixel.Measure(l)+4)
	}
	res := BreakpointsResult{Element: ui.Box(c)}
	cur := ui.Local(res.Element, lineKey{}, func() int { return 0 })
	*cur = min(*cur, len(lines)-1)
	res.Element.Children(func() {
		shell, base, st, picked := pixel.Pad(c, text, 1, len(lines), w, pixel.LineHeight, cur, o.Disabled)
		if picked >= 0 {
			n := picked + 1
			marks[n], res.Changed = !marks[n], true
			if !marks[n] {
				delete(marks, n)
			}
			slog.Debug("8bit: breakpoints", "label", text, "line", n, "on", marks[n])
		}
		say := strconv.Itoa(*cur + 1)
		if marks[*cur+1] {
			say += ", " + k.Msg("breakpoint")
		}
		base.Description(fmt.Sprintf("%s: %s", say, lines[*cur]))
		at := *cur
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			off := 0
			for i, l := range lines {
				y, from := i*pixel.LineHeight, off
				off += len(l) + 1
				if i == at && st.Focus {
					g.Cursor(0, y, w, pixel.LineHeight)
				}
				if marks[i+1] {
					g.Disc(5, y+6, 3, pixel.Bad)
				}
				n := strconv.Itoa(i + 1)
				g.Text(gutter-6-pixel.Measure(n), y, n, pixel.Faint, 1)
				colored(g, gutter, y, l, slots[from:from+len(l)])
			}
			g.Finish(st)
		})
	})
	return res
}
