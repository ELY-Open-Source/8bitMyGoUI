package layout

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

const barW = 10

var (
	upArrow = pixel.MustSprite(`
		..#..
		.###.
		#####
	`)
	downArrow = pixel.MustSprite(`
		#####
		.###.
		..#..
	`)
)

type ScrollAreaOptions struct {
	// Label names the region for assistive technology.
	Label    string
	Disabled bool
}

type ScrollAreaResult struct {
	Element *ui.Element
	// Changed reports the view moved since the last frame.
	Changed bool
}

type scrollKey struct{}

// ScrollArea scrolls children vertically beside a pixel scrollbar.
func ScrollArea(c *ui.Context, s *ui.ScrollState, o ScrollAreaOptions, children func()) ScrollAreaResult {
	k := pixel.K(c)
	// The bar comes first: TrackScroll reads s when called.
	row := ui.Row(c).AlignItems(ui.Stretch).Reverse()
	last := ui.Local(row, scrollKey{}, func() float32 { return s.Y })
	view := row.Bounds().H
	// Disabled freezes the view: any move is undone before paint.
	if o.Disabled {
		s.Y = *last
	}
	var v *ui.Element
	row.Children(func() {
		scrollBar(c, k, s, view, o.Disabled)
		// A disabled view is a new element, so the focus leaves it.
		v = ui.Scroll(c).Key(o.Disabled).TrackScroll(s).Grow(1).Label(o.Label).FocusRing(false)
		if o.Disabled {
			v.Role(ui.RoleGroup).Description(k.Msg("disabled"))
		} else {
			v.Role(ui.RoleScroll).Focusable()
		}
		v.Children(children)
	})
	focus := v.FocusVisible()
	row.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.H != view {
			c.Invalidate()
		}
	}).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focus {
			g := k.Canvas(p, r)
			w := g.W - barW
			g.Checker(0, 0, w, 1, pixel.Accent)
			g.Checker(0, g.H-1, w, 1, pixel.Accent)
			g.Checker(0, 0, 1, g.H, pixel.Accent)
			g.Checker(w-1, 0, 1, g.H, pixel.Accent)
		}
	})
	res := ScrollAreaResult{Element: row, Changed: s.Y != *last}
	*last = s.Y
	if res.Changed {
		slog.Debug("8bit: scrollarea scrolled", "y", s.Y)
	}
	return res
}

// scrollBar draws arrows, track and thumb, and moves s.
func scrollBar(c *ui.Context, k pixel.Kit, s *ui.ScrollState, view float32, disabled bool) {
	clamp := func(y float32) { s.Y = max(0, min(y, s.MaxY)) }
	ui.Column(c).Width(k.D(barW)).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if disabled {
			g := k.Canvas(p, r)
			g.Checker(0, 0, g.W, g.H, pixel.Face)
		}
	}).Children(func() {
		if arrow(c, k, upArrow, k.Msg("scroll_up"), disabled) {
			clamp(s.Y - k.D(16))
		}
		track := ui.Box(c).Grow(1)
		h := track.Bounds().H
		size, at := thumb(h, view, s.Y, s.MaxY, k.Dot())
		track.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.H != h {
				c.Invalidate()
			}
			box(k.Canvas(p, r), pixel.Sunken, pixel.Well)
		}).Children(func() {
			th := ui.Box(c).Absolute().Top(at).Left(0).FillWidth().Height(size)
			th.Clicked()
			if _, dy, ok := th.Dragged(); ok && !disabled && dy != 0 && h > size {
				clamp(s.Y + dy*s.MaxY/(h-size))
			}
			th.Draw(func(p *ui.Painter, r ui.Rect) { box(k.Canvas(p, r), pixel.Raised, pixel.Face) })
		})
		if track.Clicked() && !disabled {
			if _, y, _ := track.PointerPosition(); y < at {
				clamp(s.Y - view)
			} else if y > at+size {
				clamp(s.Y + view)
			}
		}
		if arrow(c, k, downArrow, k.Msg("scroll_down"), disabled) {
			clamp(s.Y + k.D(16))
		}
	})
}

// thumb sizes and places the thumb in track h.
func thumb(h, view, y, maxY, dot float32) (size, at float32) {
	if maxY <= 0 || view <= 0 {
		return h, 0
	}
	size = min(h, max(8*dot, h*view/(view+maxY)))
	size = float32(int(size/dot)) * dot
	return size, float32(int((h-size)*y/maxY/dot)) * dot
}

// arrow is a scrollbar button; keys scroll the view instead.
func arrow(c *ui.Context, k pixel.Kit, sp *pixel.Sprite, label string, disabled bool) bool {
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Role(ui.RoleButton).Label(label).Disabled(disabled)
	})
	shell.Size(k.D(barW), k.D(barW)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		style, y := pixel.Raised, 0
		if st.Press {
			style, y = pixel.Sunken, 1
		}
		box(g, style, pixel.Face)
		g.Sprite((g.W-sp.W)/2, (g.H-sp.H)/2+y, sp, 0, pixel.Ink, 1)
	})
	return base.Clicked()
}
