package chat

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ChatLogOptions struct {
	// Width and Height are dots; zeros mean 220 by 140.
	Width, Height int
}

type logKey struct{}

// logState is the scroll and the end it last reached.
type logState struct {
	scroll ui.ScrollState
	end    float32
}

// ChatLog scrolls messages and follows the end while there.
func ChatLog(c *ui.Context, text string, o ChatLogOptions, content func()) *ui.Element {
	k := pixel.K(c)
	w, h := pixel.Or("ChatLog Width", o.Width, 220), pixel.Or("ChatLog Height", o.Height, 140)
	box := ui.Box(c).Size(k.D(w), k.D(h))
	st := ui.Local(box, logKey{}, func() logState { return logState{scroll: ui.ScrollState{Y: math.MaxFloat32}} })
	// Content grew under a view at its old end: follow.
	grew := func() bool { return st.scroll.MaxY > st.end && st.scroll.Y >= st.end }
	// Inert freezes the view; the end waits until it wakes.
	frozen := pixel.Current(c).Inert
	if !frozen {
		if grew() {
			st.scroll.Y = st.scroll.MaxY
		}
		st.end = st.scroll.MaxY
	}
	return box.Draw(func(p *ui.Painter, r ui.Rect) {
		// Layout has just measured the content; follow it next frame.
		if !frozen && grew() {
			c.Invalidate()
		}
	}).Children(func() {
		layout.ScrollArea(c, &st.scroll, layout.ScrollAreaOptions{Label: text}, func() {
			ui.Column(c).FillWidth().Padding(k.D(4)).Gap(k.D(4)).Children(content)
		}).Element.Fill()
	})
}
