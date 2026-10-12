package game

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LetterboxOptions struct {
	// Bar is each bar's height in dots; zero means 16.
	Bar int
}

// Letterbox slides cutscene bars over content in four steps.
func Letterbox(c *ui.Context, on bool, o LetterboxOptions, content func()) *ui.Element {
	k := pixel.K(c)
	bar := pixel.Or("Letterbox Bar", o.Bar, 16)
	var cover *ui.Element
	e := layout.Layers(c, content, func() { cover = ui.Box(c).Fill().PassThrough() })
	h := bar * k.Toward(c, e, on, 8, 4) / 4
	cover.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, h, pixel.Ink)
		g.Rect(0, g.H-h, g.W, h, pixel.Ink)
	})
	return e
}
