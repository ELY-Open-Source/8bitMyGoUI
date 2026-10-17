package layout

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ResizableOptions struct {
	// MinW and MinH are zero for 24.
	MinW, MinH int
	// MaxW and MaxH are zero for no limit.
	MaxW, MaxH int
	Label      string
	Disabled   bool
}

type ResizableResult struct {
	Element ui.Element
	Changed bool
}

// Resizable is a *w×*h dot box with a ◢ grip.
func Resizable(c *ui.Context, w, h *int, o ResizableOptions, children func()) ResizableResult {
	k := pixel.K(c)
	minW, minH := pixel.Or("Resizable MinW", o.MinW, 24), pixel.Or("Resizable MinH", o.MinH, 24)
	maxW, maxH := pixel.Or("Resizable MaxW", o.MaxW, 0), pixel.Or("Resizable MaxH", o.MaxH, 0)
	if maxW > 0 && maxW < minW || maxH > 0 && maxH < minH {
		panic(fmt.Sprintf("layout: Resizable max %d×%d under min %d×%d", maxW, maxH, minW, minH))
	}
	clamp := func(v, lo, hi int) int {
		if hi > 0 {
			v = min(v, hi)
		}
		return max(v, lo)
	}
	was := [2]int{*w, *h}
	e := ui.Box(c).Clip()
	e.Children(func() {
		children()
		shell, base, st := pixel.Shell(c, func() ui.Element {
			return ui.Box(c).Focusable().Role(ui.RoleSplitter).Label(label(k, o.Label, "resize")).Disabled(o.Disabled).Cursor(ui.CursorResizeNWSE)
		})
		if dx, dy, ok := drag(k, base); ok {
			*w, *h = *w+dx, *h+dy
		}
		for key, d := range map[ui.Key][2]int{ui.KeyLeft: {-4, 0}, ui.KeyRight: {4, 0}, ui.KeyUp: {0, -4}, ui.KeyDown: {0, 4}} {
			if !o.Disabled && pixel.Shortcut(base, 0, key) {
				*w, *h = *w+d[0], *h+d[1]
			}
		}
		shell.Absolute().Right(0).Bottom(0).Size(k.D(9), k.D(9)).Draw(func(p *ui.Painter, r ui.Rect) { corner(k.Canvas(p, r), st) })
	})
	*w, *h = clamp(*w, minW, maxW), clamp(*h, minH, maxH)
	e.Size(k.D(*w), k.D(*h))
	res := ResizableResult{Element: e, Changed: [2]int{*w, *h} != was}
	if res.Changed {
		slog.Debug("8bit: resizable", "w", *w, "h", *h)
	}
	return res
}

// corner is a ◢ of three ridges in the bottom-right.
func corner(g pixel.Canvas, st pixel.State) {
	for _, s := range []int{2, 5, 8} {
		for i := 0; i <= s; i++ {
			x, y := g.W-1-s+i, g.H-1-i
			g.Dot(x, y, pixel.Shade)
			if x > 0 && y > 0 {
				g.Dot(x-1, y-1, pixel.Lit)
			}
		}
	}
	if st.Disabled {
		g.Checker(0, 0, g.W, g.H, pixel.Face)
	}
	if st.Focus {
		g.FocusMark(0, 0, g.W, g.H)
	}
}
