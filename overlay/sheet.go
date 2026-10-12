package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SheetOptions struct {
	// Height is in dots; zero means 120.
	Height   int
	Disabled bool
}

type SheetResult struct {
	Element *ui.Element
	// Closed is a drag down, Escape or a click outside.
	Closed bool
}

type sheetKey struct{}

// Sheet rises from the bottom; drag its bar down.
func Sheet(c *ui.Context, text string, open *bool, o SheetOptions, content func()) SheetResult {
	k := pixel.K(c)
	h := pixel.Or("Sheet Height", o.Height, 120)
	off := o.Disabled || pixel.Current(c).Inert
	home := ui.Box(c)
	pulled := ui.Local(home, sheetKey{}, func() float32 { return 0 })
	rise := k.D(arriving(c, home, *open) * h / 4)
	was := *open
	var res SheetResult
	res.Element = ui.DialogBase(c, open, func(back, panel *ui.Element) {
		dim(k, back, pixel.Back, 8)
		panel.Label(text).Focusable().FocusRing(false).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
		}).Children(func() {
			pixel.Harbor(c, panel, func() {
				bar := ui.Box(c).FillWidth().Height(k.D(10)).Role(ui.RoleNone)
				_, dy, held := bar.Dragged()
				switch {
				case off:
					*pulled = 0
				case held:
					*pulled = max(0, *pulled+dy)
				case *pulled > 0:
					if *pulled > k.D(h/3) {
						// A release is no drag: build again, closed.
						*open = false
						c.Invalidate()
					}
					*pulled = 0
				}
				bar.Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					g.Rect((g.W-24)/2, 4, 24, 2, pixel.Shade)
				})
				ui.Column(c).Grow(1).Padding(0, k.D(8), k.D(8), k.D(8)).Children(func() { body(c, o.Disabled, content) })
			})
		})
		// Placed after the bar, so a release lands at once.
		drop := float32(int(*pulled/k.Dot())) * k.Dot()
		panel.Absolute().Left(0).Right(0).Bottom(-rise - drop).Height(k.D(h))
	})
	res.Closed = was && !*open
	if !*open {
		*pulled = 0
	}
	if res.Closed {
		slog.Debug("8bit: sheet closed", "label", text)
	}
	return res
}
