package text

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type Align uint8

const (
	Left Align = iota
	Center
	Right
)

type ParagraphOptions struct {
	// Gap is extra dots between lines.
	Gap      int
	Align    Align
	MaxLines int
	Slot     pixel.Slot
}

// Paragraph wraps s at the width it had last frame.
func Paragraph(c *ui.Context, s string, o ParagraphOptions) ui.Element {
	k := pixel.K(c)
	gap := pixel.Or("Paragraph Gap", o.Gap, 0)
	if o.Align > Right {
		panic(fmt.Sprintf("text: unknown Align %d", o.Align))
	}
	e := ui.Box(c).FillWidth().Label(s).Role(ui.RoleText)
	w, check := pixel.WrapWidth(c, e)
	lines := pixel.Wrap(s, w, o.MaxLines)
	step := pixel.LineHeight + gap
	return e.Height(k.D(max(0, len(lines)*step-gap))).Draw(func(p *ui.Painter, r ui.Rect) {
		check(r)
		g := k.Canvas(p, r)
		for i, l := range lines {
			x := 0
			switch o.Align {
			case Center:
				x = (g.W - pixel.Measure(l)) / 2
			case Right:
				x = g.W - pixel.Measure(l)
			}
			g.Text(x, i*step, l, o.Slot, 1)
		}
	})
}
