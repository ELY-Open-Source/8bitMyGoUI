package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DrawerOptions struct {
	// Side is the edge it slides from.
	Side display.Side
	// Size is its depth in dots; zero means 160.
	Size int
	// Disabled makes the content inert.
	Disabled bool
}

type DrawerResult struct {
	Element *ui.Element
	// Closed is Escape or a click outside.
	Closed bool
}

// Drawer slides a panel in from an edge in steps.
func Drawer(c *ui.Context, text string, open *bool, o DrawerOptions, content func()) DrawerResult {
	k := pixel.K(c)
	if o.Side > display.Right {
		panic("overlay: unknown Drawer Side")
	}
	size := pixel.Or("Drawer Size", o.Size, 160)
	home := ui.Box(c)
	off := k.D(arriving(c, home, *open) * size / 4)
	was := *open
	var res DrawerResult
	res.Element = ui.DialogBase(c, open, func(back, panel *ui.Element) {
		dim(k, back, pixel.Back, 8)
		panel.Absolute().Label(text).Padding(k.D(8)).Focusable().FocusRing(false)
		switch o.Side {
		case display.Left:
			panel.Top(0).Bottom(0).Left(-off).Width(k.D(size))
		case display.Right:
			panel.Top(0).Bottom(0).Right(-off).Width(k.D(size))
		case display.Top:
			panel.Left(0).Right(0).Top(-off).Height(k.D(size))
		case display.Bottom:
			panel.Left(0).Right(0).Bottom(-off).Height(k.D(size))
		}
		panel.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
		}).Children(func() {
			pixel.Harbor(c, panel, func() { body(c, o.Disabled, content) })
		})
	})
	res.Closed = was && !*open
	if res.Closed {
		slog.Debug("8bit: drawer closed", "label", text)
	}
	return res
}
