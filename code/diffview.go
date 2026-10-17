package code

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DiffViewOptions struct{}

// change is a diff line's sign, text and line numbers.
type change struct {
	sign     byte
	text     string
	old, new int
}

// DiffView marks lines removed with -, added with +.
func DiffView(c *ui.Context, text, before, after string, o DiffViewOptions) ui.Element {
	k := pixel.K(c)
	changes := diff(strings.Split(tidy(before), "\n"), strings.Split(tidy(after), "\n"))
	// Rows share the longest width, so tints line up.
	w := 0
	for _, ch := range changes {
		w = max(w, 32+pixel.Measure(ch.text))
	}
	return layout.Panel(c, layout.PanelOptions{Style: pixel.Sunken, Fill: pixel.Well, Padding: 2}, func() {
		for _, ch := range changes {
			ui.Box(c).Size(k.D(w), k.D(pixel.LineHeight)).Role(ui.RoleText).Label(string(ch.sign) + " " + ch.text).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				// Changed lines take their color and a left bar.
				ink, text := pixel.Faint, pixel.Ink
				switch ch.sign {
				case '+':
					ink, text = pixel.Good, pixel.Good
				case '-':
					ink, text = pixel.Bad, pixel.Bad
				}
				if ch.sign != ' ' {
					g.Rect(0, 0, 2, g.H, ink)
				}
				num := func(x, n int) {
					if n > 0 {
						s := fmt.Sprint(n)
						g.Micro(x-pixel.MicroWidth(s), 4, s, pixel.Faint, 1)
					}
				}
				num(11, ch.old)
				num(23, ch.new)
				g.Text(26, 0, string(ch.sign), ink, 1)
				g.Text(32, 0, ch.text, text, 1)
			})
		}
	}).Role(ui.RoleGroup).Label(text)
}

// diff walks the longest common run of lines.
func diff(a, b []string) []change {
	// lcs[i][j] is the common length of a[i:] and b[j:].
	lcs := make([][]int, len(a)+1)
	for i := range lcs {
		lcs[i] = make([]int, len(b)+1)
	}
	for i := len(a) - 1; i >= 0; i-- {
		for j := len(b) - 1; j >= 0; j-- {
			if a[i] == b[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else {
				lcs[i][j] = max(lcs[i+1][j], lcs[i][j+1])
			}
		}
	}
	var out []change
	i, j := 0, 0
	for i < len(a) || j < len(b) {
		switch {
		case i < len(a) && j < len(b) && a[i] == b[j]:
			out = append(out, change{' ', a[i], i + 1, j + 1})
			i, j = i+1, j+1
		// Taken out before put in, as diffs read.
		case i < len(a) && (j == len(b) || lcs[i+1][j] >= lcs[i][j+1]):
			out = append(out, change{'-', a[i], i + 1, 0})
			i++
		default:
			out = append(out, change{'+', b[j], 0, j + 1})
			j++
		}
	}
	return out
}
