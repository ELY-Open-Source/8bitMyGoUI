package chat

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ChatLogOptions struct {
	// Width and Height are dots; zeros mean 220 by 140.
	Width, Height int
}

type logKey struct{}

// ChatLog scrolls messages and follows the end while there.
func ChatLog(c *ui.Context, text string, o ChatLogOptions, content func()) *ui.Element {
	k := pixel.K(c)
	w, h := pixel.Or("ChatLog Width", o.Width, 220), pixel.Or("ChatLog Height", o.Height, 140)
	box := ui.Box(c).Size(k.D(w), k.D(h))
	s := ui.Local(box, logKey{}, func() ui.ScrollState { return ui.ScrollState{} })
	return box.Children(func() {
		layout.ScrollArea(c, s, layout.ScrollAreaOptions{Label: text, Follow: true}, func() {
			ui.Column(c).FillWidth().Padding(k.D(4)).Gap(k.D(4)).Children(content)
		}).Element.Fill()
	})
}
