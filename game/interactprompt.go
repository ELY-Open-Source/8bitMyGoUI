package game

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type InteractPromptOptions struct {
	// Key is the pad button's letter; zero means A.
	Key string
}

// InteractPrompt bobs a Ⓐ TALK bubble above target.
func InteractPrompt(c *ui.Context, target *ui.Element, action string, o InteractPromptOptions) *ui.Element {
	k := pixel.K(c)
	key := o.Key
	if key == "" {
		key = "A"
	}
	if len([]rune(key)) != 1 || !pixel.MicroOK(key) {
		panic(fmt.Sprintf("game: InteractPrompt key %q is not one micro letter", key))
	}
	b := target.Bounds()
	ww, wh := c.Size()
	// The bubble's tail points at the target's middle.
	x := max(0, min(b.X+b.W/2-k.D(10), ww-k.D(pixel.Measure(action)+30)))
	at := max(3, int((b.X+b.W/2-x)/k.Dot()))
	bob := 0
	if !k.Reduce() {
		bob = k.Tick(c, 4) / 2 % 2
	}
	var e *ui.Element
	ui.Overlay(c, func() {
		e = ui.Box(c).Absolute().Left(x).Bottom(wh - b.Y + k.D(1+bob)).PassThrough().Role(ui.RoleText).Label(key + " " + action)
		// Bounds lag a layout; if the target moved, build again.
		e.Draw(func(p *ui.Painter, r ui.Rect) {
			if target.Bounds() != b {
				c.Invalidate()
			}
		})
		e.Children(func() {
			display.Bubble(c, display.BubbleOptions{Tail: display.Bottom, At: at}, func() {
				ui.Row(c).Gap(k.D(3)).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(k.D(11), k.D(11)).Draw(func(p *ui.Painter, r ui.Rect) {
						g := k.Canvas(p, r)
						g.Disc(5, 5, 5, pixel.Ink)
						g.Micro((11-pixel.MicroWidth(key))/2, 3, key, pixel.Back, 1)
					})
					pixel.Text(c, action, pixel.TextOptions{})
				})
			})
		})
	})
	return e
}
