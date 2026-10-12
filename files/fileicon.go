package files

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type FileIconOptions struct {
	// Scale is 1–4; zero means 1.
	Scale int
}

// FileIcon draws a kind: a folder or a dog-eared page.
func FileIcon(c *ui.Context, kind Kind, o FileIconOptions) *ui.Element {
	k := pixel.K(c)
	ka := look(kind)
	s := pixel.Or("FileIcon Scale", o.Scale, 1)
	if s > 4 {
		panic(fmt.Sprintf("files: FileIcon Scale %d is past 4", s))
	}
	return ui.Box(c).Size(k.D(12*s), k.D(12*s)).Role(ui.RoleImage).Label(k.Msg(ka.msg)).Draw(func(p *ui.Painter, r ui.Rect) {
		k.Canvas(p, r).Sprite(0, 0, ka.art, 0, pixel.Ink, s)
	})
}
