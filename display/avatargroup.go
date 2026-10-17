package display

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type AvatarGroupOptions struct {
	// Max faces shown before +N; zero means 4.
	Max   int
	Scale int
}

// AvatarGroup overlaps faces and counts the rest in a chip.
func AvatarGroup(c *ui.Context, names []string, o AvatarGroupOptions) ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	shown := names[:min(len(names), pixel.Or("AvatarGroup Max", o.Max, 4))]
	rest := len(names) - len(shown)
	more := ""
	if rest > 0 {
		more = "+" + itoa(rest)
	}
	step := 12 * sc
	w := max(0, len(shown)-1)*step + 16*sc
	if more != "" {
		w += pixel.MicroWidth(more) + 8
	}
	return art(c, k, strings.Join(names, ", "), w, 16*sc, func(g pixel.Canvas) {
		for i, name := range shown {
			x := i * step
			Face(g, x, 0, name, sc)
			if i > 0 {
				g.VLine(x, 0, 16*sc, pixel.Back)
			}
		}
		if more != "" {
			x := len(shown)*step + 4*sc
			cw := pixel.MicroWidth(more) + 4
			g.Frame(x, (16*sc-9)/2, cw, 9, pixel.Notched, pixel.Face)
			g.Micro(x+2, (16*sc-9)/2+2, more, pixel.Ink, 1)
		}
	})
}
