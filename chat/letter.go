package chat

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LetterOptions struct {
	// Signature closes the letter in Accent.
	Signature string
	// Width is in dots; zero means 180.
	Width int
}

// Letter sets a body on parchment with a signature.
func Letter(c *ui.Context, body string, o LetterOptions) *ui.Element {
	k := pixel.K(c)
	w := pixel.Or("Letter Width", o.Width, 180)
	if w < 40 {
		panic("chat: a Letter needs 40 dots")
	}
	lines := pixel.Wrap(body, w-20, 0)
	h := 20 + len(lines)*pixel.LineHeight
	if o.Signature != "" {
		h += pixel.LineHeight + 4
	}
	label := body
	if o.Signature != "" {
		label += " — " + o.Signature
	}
	return ui.Box(c).Role(ui.RoleText).Label(label).Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, w, h, pixel.Face)
		// The edge darkens outward in three dither rings.
		for i, level := range []int{12, 7, 3} {
			g.Veil(i, i, w-2*i, 1, pixel.Shade, level)
			g.Veil(i, h-1-i, w-2*i, 1, pixel.Shade, level)
			g.Veil(i, i+1, 1, h-2-2*i, pixel.Shade, level)
			g.Veil(w-1-i, i+1, 1, h-2-2*i, pixel.Shade, level)
		}
		for i, l := range lines {
			g.Text(10, 10+i*pixel.LineHeight, l, pixel.Ink, 1)
		}
		if o.Signature == "" {
			return
		}
		y, sw := 14+len(lines)*pixel.LineHeight, pixel.Measure(o.Signature)
		g.Text(w-10-sw, y, o.Signature, pixel.Accent, 1)
		g.HLine(w-12-sw, y+pixel.LineHeight-2, sw+2, pixel.Accent)
		g.Dot(w-9, y+pixel.LineHeight-3, pixel.Accent)
	})
}
