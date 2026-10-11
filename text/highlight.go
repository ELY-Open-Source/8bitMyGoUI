package text

import (
	"unicode"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type HighlightOptions struct {
	Slot pixel.Slot
	// Mark is the marker slot; zero means Warn.
	Mark pixel.Slot
}

// Highlight marks each match of query in s, ignoring case.
func Highlight(c *ui.Context, s, query string, o HighlightOptions) *ui.Element {
	k := pixel.K(c)
	mark := slotOr(o.Mark, pixel.Warn)
	rs := []rune(s)
	spans := matches(rs, []rune(query))
	return block(c, k, s, pixel.Measure(s), pixel.LineHeight, func(g pixel.Canvas) {
		for _, sp := range spans {
			x0, x1 := pixel.Measure(string(rs[:sp[0]])), pixel.Measure(string(rs[:sp[1]]))
			g.Veil(x0, 2, x1-x0, pixel.LineHeight-4, mark, 10)
		}
		g.Text(0, 0, s, o.Slot, 1)
	})
}

// matches are the rune spans of q in s.
func matches(s, q []rune) [][2]int {
	var out [][2]int
	if len(q) == 0 {
		return nil
	}
	for i := 0; i+len(q) <= len(s); {
		same := true
		for j := range q {
			if unicode.ToLower(s[i+j]) != unicode.ToLower(q[j]) {
				same = false
				break
			}
		}
		if same {
			out = append(out, [2]int{i, i + len(q)})
			i += len(q)
		} else {
			i++
		}
	}
	return out
}
