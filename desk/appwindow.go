package desk

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type AppWindowOptions struct {
	// The body's Width and Height are dots; zeros mean 160×100.
	Width, Height int
	// Room is its area in dots; zero leaves it open.
	Room     [2]int
	Disabled bool
}

type AppWindowResult struct {
	Element *ui.Element
	// Closed is the close box; Moved, a drag or arrow.
	Closed, Moved bool
}

type dragKey struct{}

// AppWindow floats at x, y; drag or arrow its bar.
func AppWindow(c *ui.Context, title string, x, y *int, o AppWindowOptions, content func()) AppWindowResult {
	k := pixel.K(c)
	w, h := pixel.Or("AppWindow Width", o.Width, 160), pixel.Or("AppWindow Height", o.Height, 100)
	// The bar needs room for a grip and close box.
	if w < 40 {
		panic(fmt.Sprintf("desk: AppWindow Width %d is under 40", w))
	}
	// The window is its body, bar and frame.
	ww, wh := w+6, h+pixel.LineHeight+10
	if o.Room[0] > 0 && o.Room[0] < ww || o.Room[1] > 0 && o.Room[1] < wh {
		panic(fmt.Sprintf("desk: AppWindow %d×%d does not fit its Room %v", ww, wh, o.Room))
	}
	ox, oy := *x, *y
	keep := func() {
		if o.Room[0] > 0 {
			*x = min(*x, o.Room[0]-ww)
		}
		if o.Room[1] > 0 {
			*y = min(*y, o.Room[1]-wh)
		}
		*x, *y = max(0, *x), max(0, *y)
	}
	keep()
	win := ui.Column(c).Absolute().Width(k.D(ww)).Padding(k.D(3)).Gap(k.D(2)).Role(ui.RoleGroup).Label(title)
	// rest keeps a drag's part-dot for the next frame.
	rest := ui.Local(win, dragKey{}, func() [2]float32 { return [2]float32{} })
	active := win.FocusWithin()
	if o.Disabled {
		win.Description(k.Msg("disabled"))
	}
	res := AppWindowResult{Element: win}
	win.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Raised, pixel.Face)
	}).Children(func() {
		ui.Row(c).FillWidth().Gap(k.D(2)).AlignItems(ui.Center).Children(func() {
			shell, base, st := pixel.Shell(c, func() *ui.Element {
				return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(title).Description(k.Msg("arrows_move")).Disabled(o.Disabled)
			})
			if !base.IsDisabled() {
				if dx, dy, ok := base.Dragged(); ok {
					rest[0], rest[1] = rest[0]+dx/k.Dot(), rest[1]+dy/k.Dot()
					*x, *y = *x+int(rest[0]), *y+int(rest[1])
					rest[0], rest[1] = rest[0]-float32(int(rest[0])), rest[1]-float32(int(rest[1]))
				}
				for _, m := range []struct {
					key    ui.Key
					dx, dy int
				}{{ui.KeyLeft, -4, 0}, {ui.KeyRight, 4, 0}, {ui.KeyUp, 0, -4}, {ui.KeyDown, 0, 4}} {
					if base.Shortcut(0, m.key) {
						*x, *y = *x+m.dx, *y+m.dy
					}
				}
			}
			shell.Grow(1).Height(k.D(pixel.LineHeight + 2)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				barLook(g, title, active)
				g.Finish(st)
			})
			res.Closed = input.IconButton(c, icon.Close, k.Msg("close"), input.IconButtonOptions{Size: input.Small, Disabled: o.Disabled}).Clicked
		})
		// Off, the body's controls are off too, and veiled.
		in := pixel.Current(c)
		in.Inert = in.Inert || o.Disabled
		ui.Box(c).Size(k.D(w), k.D(h)).Clip().Children(func() { pixel.Scope(c, in, content) }).DrawOver(func(p *ui.Painter, r ui.Rect) {
			if o.Disabled {
				g := k.Canvas(p, r)
				g.Checker(0, 0, g.W, g.H, pixel.Face)
			}
		})
	})
	keep()
	win.Left(k.D(*x)).Top(k.D(*y))
	res.Moved = *x != ox || *y != oy
	if res.Moved || res.Closed {
		slog.Debug("8bit: app window", "title", title, "x", *x, "y", *y, "closed", res.Closed)
	}
	return res
}
