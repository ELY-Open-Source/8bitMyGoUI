// Package overlay floats windows, bubbles and veils over the view.
package overlay

import (
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// dim checkers the backdrop of a dialog.
func dim(k pixel.Kit, back *ui.Element, s pixel.Slot, level int) {
	back.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Veil(0, 0, g.W, g.H, s, level)
	})
}

// dark puts the deeper of Ink and Back first.
func dark(k pixel.Kit) (dark, light pixel.Slot) {
	if k.Palette().Night {
		return pixel.Back, pixel.Ink
	}
	return pixel.Ink, pixel.Back
}

// choose lists items under a hand; digits pick 1–9.
func choose(c *ui.Context, text string, items []string, cur *int, digits, off bool) int {
	k := pixel.K(c)
	n := len(items)
	*cur = max(0, min(*cur, n-1))
	off = off || pixel.Current(c).Inert
	box := ui.Column(c).FocusRing(false).Role(ui.RoleList).Label(text)
	if !off {
		box.Focusable().AutoFocus()
	}
	pixel.Park(c, box, off)
	picked := -1
	if off {
		box.Description(k.Msg("disabled"))
	}
	// The list stays enabled, so off is checked here.
	if !off && pixel.Shortcut(box, 0, ui.KeyDown) {
		*cur = (*cur + 1) % n
	}
	if !off && pixel.Shortcut(box, 0, ui.KeyUp) {
		*cur = (*cur + n - 1) % n
	}
	enter, space := pixel.Shortcut(box, 0, ui.KeyEnter), pixel.Shortcut(box, 0, ui.KeySpace)
	if (enter || space) && !off {
		picked = *cur
	}
	for i := range min(9, n) {
		if digits && !off && pixel.Shortcut(box, 0, ui.Key1+ui.Key(i)) {
			*cur, picked = i, i
		}
	}
	box.Children(func() {
		for i, it := range items {
			row := ui.Row(c).Height(k.D(pixel.LineHeight)).MinWidth(k.D(pixel.Measure(it) + 16)).Role(ui.RoleListItem).Label(it)
			if !off && row.Clicked() {
				*cur, picked = i, i
			}
			row.Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				if i == *cur {
					g.Sprite(g.Step(4)%2, (g.H-pixel.Hand.H)/2, pixel.Hand, 0, pixel.Edge, 1)
				}
				g.Text(14, 0, it, pixel.Ink, 1)
			})
		}
	})
	if !off {
		box.Description(items[*cur])
	}
	return picked
}

// body builds content, inert while off.
func body(c *ui.Context, off bool, content func()) {
	if !off {
		content()
		return
	}
	s := pixel.Current(c)
	s.Inert = true
	pixel.Scope(c, s, content)
}

type hintKey struct{}

// hinting is when hover began, and a press-given focus.
type hinting struct {
	since   time.Time
	pressed bool
}

// hint holds after resting wait, or on keyboard focus.
func hint(c *ui.Context, e *ui.Element, wait time.Duration) bool {
	k := pixel.K(c)
	s := ui.Local(e, hintKey{}, func() hinting { return hinting{} })
	_, _, over := e.PointerPosition()
	hovered, focus := e.Hovered(), e.FocusWithin()
	switch {
	case over && !hovered:
		// A press inside: its focus is not the keyboard's.
		s.pressed = true
	case !focus:
		s.pressed = false
	}
	rested := false
	if hovered {
		if s.since.IsZero() {
			s.since = k.Now()
		}
		steps := int(wait / (100 * time.Millisecond))
		rested = k.Since(c, s.since, 10, steps) >= steps
	} else {
		s.since = time.Time{}
	}
	return rested || focus && !s.pressed
}

// arriving counts the steps, of four, a panel still slides.
func arriving(c *ui.Context, home *ui.Element, open bool) int {
	k := pixel.K(c)
	on := 0
	if open {
		on = 1
	}
	return 4 - k.Shot(c, home, on, 15, 4)
}
