package game

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Credit is a heading and the names under it.
type Credit struct {
	Heading string
	Names   []string
}

type CreditsOptions struct {
	// Width and Height are dots; zeros mean 160 by 96.
	Width, Height int
	// Cue rolls them again each time it changes.
	Cue int
}

type CreditsResult struct {
	Element *ui.Element
	// Done is true in the frame the roll ends.
	Done bool
}

type rolledKey struct{}

// Credits roll up until the last name rests mid-window.
func Credits(c *ui.Context, text string, credits []Credit, o CreditsOptions) CreditsResult {
	k := pixel.K(c)
	w, h := pixel.Or("Credits Width", o.Width, 160), pixel.Or("Credits Height", o.Height, 96)
	type line struct {
		s    string
		slot pixel.Slot
	}
	var lines []line
	var all []string
	for i, cr := range credits {
		if i > 0 {
			lines = append(lines, line{})
		}
		lines = append(lines, line{cr.Heading, pixel.Warn})
		all = append(all, cr.Heading+":")
		for _, n := range cr.Names {
			lines, all = append(lines, line{n, pixel.Ink}), append(all, n)
		}
	}
	full := len(lines) * pixel.LineHeight
	e := ui.Box(c).Role(ui.RoleGroup).Label(text).Description(strings.Join(all, " "))
	// A dot a frame, until the last line centers.
	end, roll := h/2+full-pixel.LineHeight/2, start(k, e, o.Cue)
	top, step := 0, end
	if !k.Reduce() {
		step = k.Since(c, roll.at, 15, end)
		top = h - step
	} else {
		h = full
	}
	// Done comes once a roll.
	done := ui.Local(e, rolledKey{}, func() int { return -1 })
	res := CreditsResult{Element: e, Done: step == end && *done != roll.run}
	if res.Done {
		*done = roll.run
	}
	e.Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Clip(0, 0, w, h, func() {
			for i, l := range lines {
				g.Text((w-pixel.Measure(l.s))/2, top+i*pixel.LineHeight, l.s, l.slot, 1)
			}
		})
	})
	return res
}
