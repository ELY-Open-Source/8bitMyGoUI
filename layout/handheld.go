package layout

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type HandheldOptions struct {
	// ScreenW×ScreenH is the screen in dots; zero means 160×144.
	ScreenW, ScreenH int
	// Brand is micro text on the body.
	Brand    string
	Label    string
	Disabled bool
}

type HandheldResult struct {
	Element                     *ui.Element
	Up, Down, Left, Right, A, B bool
}

const bezel = 8

// Handheld is a pocket console around screen; buttons report presses.
func Handheld(c *ui.Context, o HandheldOptions, screen func()) HandheldResult {
	k := pixel.K(c)
	sw, sh := pixel.Or("Handheld ScreenW", o.ScreenW, 160), pixel.Or("Handheld ScreenH", o.ScreenH, 144)
	res := HandheldResult{}
	res.Element = ui.Column(c).Width(k.D(sw + 2*bezel + 40)).Padding(k.D(10)).Gap(k.D(6)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(o.Label).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		box(g, pixel.Raised, pixel.Face)
		g.Micro(10, g.H-9, o.Brand, pixel.Shade, 1)
	}).Children(func() {
		ui.Box(c).Padding(k.D(bezel), k.D(bezel+10), k.D(bezel), k.D(bezel+10)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			box(g, pixel.Sunken, pixel.Shade)
			led := pixel.Good
			if o.Disabled {
				led = pixel.Faint
			}
			g.Rect(5, g.H/2-1, 3, 3, led)
		}).Children(func() {
			ui.Box(c).Size(k.D(sw), k.D(sh)).Children(screen)
		})
		ui.Box(c).Width(k.D(sw + 2*bezel + 20)).Height(k.D(44)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Rect(15, 18, 9, 9, pixel.Edge)
		}).Children(func() {
			pad := func(name string, x, y, w, h int) bool { return key(c, k, name, x, y, w, h, false, o.Disabled) }
			res.Up = pad(k.Msg("up"), 15, 7, 9, 11)
			res.Left = pad(k.Msg("left"), 4, 18, 11, 9)
			res.Right = pad(k.Msg("right"), 24, 18, 11, 9)
			res.Down = pad(k.Msg("down"), 15, 27, 9, 11)
			right := sw + 2*bezel + 20
			res.B = key(c, k, "B", right-44, 18, 15, 15, true, o.Disabled)
			res.A = key(c, k, "A", right-24, 10, 15, 15, true, o.Disabled)
		})
	})
	for name, on := range map[string]bool{"up": res.Up, "down": res.Down, "left": res.Left, "right": res.Right, "a": res.A, "b": res.B} {
		if on {
			slog.Debug("8bit: handheld press", "button", name)
		}
	}
	return res
}

// key is a d-pad arm, or a round button.
func key(c *ui.Context, k pixel.Kit, name string, x, y, w, h int, round, disabled bool) bool {
	shell, base, st := pixel.Shell(c, func() *ui.Element { return ui.ButtonBase(c).Label(name).Disabled(disabled) })
	shell.Absolute().Left(k.D(x)).Top(k.D(y)).Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		dy := 0
		if st.Press {
			dy = 1
		}
		if round {
			g.Disc(w/2, h/2, w/2, pixel.Edge)
			g.Disc(w/2, h/2-1+dy, w/2-1, pixel.Accent)
			g.Micro(w/2-1, h/2-3+dy, name, pixel.AccentInk, 1)
		} else {
			g.Rect(0, dy, g.W, g.H-dy, pixel.Edge)
			g.Rect(1, 1+dy, g.W-2, 1, pixel.Faint)
		}
		if st.Disabled {
			g.Checker(0, 0, g.W, g.H, pixel.Face)
		}
		if st.Focus {
			g.FocusMark(0, 0, g.W, g.H)
		}
	})
	return base.Clicked()
}
