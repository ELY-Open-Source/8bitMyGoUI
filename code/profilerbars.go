package code

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Span is timed work, starting inside its parent's time.
type Span struct {
	Name          string
	Start, Length time.Duration
	Children      []Span
}

type ProfilerBarsOptions struct {
	// Width is in dots; zero means 240.
	Width int
}

var depths = []pixel.Slot{pixel.Accent, pixel.Info, pixel.Good, pixel.Warn, pixel.Magic}

// ProfilerBars stacks spans under their parents, each bar its time.
func ProfilerBars(c *ui.Context, text string, root Span, o ProfilerBarsOptions) *ui.Element {
	k := pixel.K(c)
	w := pixel.Or("ProfilerBars Width", o.Width, 240)
	if root.Length <= 0 {
		panic(fmt.Sprintf("code: ProfilerBars root %q lasts %v", root.Name, root.Length))
	}
	var said []string
	rows := check(root, &said)
	x := func(d time.Duration) int { return int(int64(w) * int64(d-root.Start) / int64(root.Length)) }
	return ui.Box(c).Size(k.D(w), k.D(rows*pixel.LineHeight)).Role(ui.RoleImage).Label(text).Description(strings.Join(said, "; ")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		var bar func(s Span, depth int)
		bar = func(s Span, depth int) {
			x0, y := x(s.Start), depth*pixel.LineHeight
			bw := max(1, x(s.Start+s.Length)-x0)
			g.Rect(x0, y+1, bw, pixel.LineHeight-2, depths[depth%len(depths)])
			// Names show where they fit, times when they also fit.
			for _, s := range []string{fmt.Sprintf("%s %.1fms", s.Name, ms(s.Length)), s.Name} {
				if pixel.Measure(s)+4 <= bw {
					g.Text(x0+2, y, s, pixel.AccentInk, 1)
					break
				}
			}
			for _, ch := range s.Children {
				bar(ch, depth+1)
			}
		}
		bar(root, 0)
	})
}

// check counts rows; a child outside its parent panics.
func check(s Span, said *[]string) int {
	*said = append(*said, fmt.Sprintf("%s %.1f ms", s.Name, ms(s.Length)))
	rows := 1
	for _, ch := range s.Children {
		if ch.Length < 0 || ch.Start < s.Start || ch.Start+ch.Length > s.Start+s.Length {
			panic(fmt.Sprintf("code: ProfilerBars span %q is outside %q", ch.Name, s.Name))
		}
		rows = max(rows, 1+check(ch, said))
	}
	return rows
}
