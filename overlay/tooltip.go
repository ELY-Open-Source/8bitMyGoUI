package overlay

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type TooltipOptions struct{}

// tipWidth caps a tip's text in dots.
const tipWidth = 160

// Tooltip bubbles text after 500 ms, or on keyboard focus.
func Tooltip(c *ui.Context, text string, o TooltipOptions, trigger func()) ui.Element {
	k := pixel.K(c)
	wrap := ui.Box(c)
	show := hint(c, wrap, 500*time.Millisecond)
	wrap.Children(trigger)
	if show {
		bubble(c, wrap.Bounds(), func() {
			tip := ui.Box(c).Role(ui.RoleTooltip).Label(text)
			tip.Children(func() {
				pixel.Text(c, text, pixel.TextOptions{Wrap: true}).Width(k.D(min(pixel.Measure(text), tipWidth)))
			})
		}, min(pixel.Measure(text), tipWidth)+16)
	}
	return wrap
}

// bubble floats content by b; w keeps it inside.
func bubble(c *ui.Context, b ui.Rect, content func(), w int) ui.Element {
	k := pixel.K(c)
	ww, wh := c.Size()
	x := max(0, min(b.X+b.W/2-k.D(8), ww-k.D(w)))
	at := max(3, int((b.X+b.W/2-x)/k.Dot()))
	var e ui.Element
	ui.Overlay(c, func() {
		e = ui.Box(c).Absolute().Left(x).PassThrough()
		tail := display.Bottom
		if b.Y > k.D(40) {
			e.Bottom(wh - b.Y)
		} else {
			e.Top(b.Y + b.H)
			tail = display.Top
		}
		e.Children(func() { display.Bubble(c, display.BubbleOptions{Tail: tail, At: at}, content) })
	})
	return e
}
