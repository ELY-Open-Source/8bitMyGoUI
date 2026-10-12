package pixel

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
)

type TextOptions struct {
	Slot Slot
	// Scale is 1–4; zero means 1.
	Scale int
	// Wrap breaks lines at the width the layout gives.
	Wrap     bool
	MaxLines int
}

// Text is a pixel-font label that screen readers read.
func Text(c *ui.Context, s string, o TextOptions) *ui.Element {
	k := K(c)
	scale := scaleOr1(o.Scale)
	e := ui.Box(c).Label(s).Role(ui.RoleText)
	width, check := 0, func(ui.Rect) {}
	if o.Wrap {
		width, check = WrapWidth(c, e)
	}
	lines := Wrap(s, width/scale, o.MaxLines)
	if o.Wrap {
		e.Height(k.D(len(lines) * LineHeight * scale))
	} else {
		w := 0
		for _, l := range lines {
			w = max(w, Measure(l))
		}
		e.Size(k.D(w*scale), k.D(len(lines)*LineHeight*scale))
	}
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		check(r)
		g := k.Canvas(p, r)
		for i, l := range lines {
			g.Text(0, i*LineHeight*scale, l, o.Slot, scale)
		}
	})
}

// Wrap wraps s to width; past maxLines, an ellipsis.
func Wrap(s string, width, maxLines int) []string {
	if maxLines < 0 {
		panic(fmt.Sprintf("pixel: MaxLines %d", maxLines))
	}
	lines := Lines(s, width)
	if maxLines == 0 || len(lines) <= maxLines {
		return lines
	}
	lines = lines[:maxLines]
	last := []rune(strings.TrimRight(lines[maxLines-1], " "))
	for width > 0 && len(last) > 0 && Measure(string(last)+"…") > width {
		last = last[:len(last)-1]
	}
	lines[maxLines-1] = string(last) + "…"
	return lines
}

type wrapKey struct{}

// WrapWidth is e's last painted width; Draw must call check.
func WrapWidth(c *ui.Context, e *ui.Element) (int, func(r ui.Rect)) {
	dot := K(c).dot
	w := ui.Local(e, wrapKey{}, func() int { return 0 })
	used := *w
	return used, func(r ui.Rect) {
		if now := int(r.W / dot); now != *w {
			*w = now
			c.Invalidate()
		}
	}
}

func scaleOr1(s int) int {
	if s == 0 {
		return 1
	}
	checkScale(s)
	return s
}

// State is how a control looks this frame.
type State struct{ Hover, Press, Focus, Disabled bool }

// Shell wraps base in an enabled box that paints.
func Shell(c *ui.Context, base func() *ui.Element) (shell, b *ui.Element, st State) {
	shell = ui.Box(c)
	shell.Children(func() { b = base().Fill().FocusRing(false) })
	own, inert := b.IsDisabled(), current(c).cur.Inert
	if inert {
		b.Disabled(true)
	}
	off := own || inert
	Park(c, b, off)
	st = State{Hover: b.Hovered(), Press: b.Pressed() && !off, Focus: b.FocusVisible(), Disabled: own}
	return shell, b, st
}

type parkKey struct{}

// Park moves focus off e while off, back once on.
func Park(c *ui.Context, e *ui.Element, off bool) {
	parked := ui.Local(e, parkKey{}, func() bool { return false })
	switch {
	case off && e.Focused():
		c.Root().Focus()
		*parked = true
	case !off && *parked:
		// Focus comes back unless something else took it.
		if root := c.Root(); root.Focused() || !root.FocusWithin() {
			e.Focus()
		}
		*parked = false
	}
}

// Hand is the pointing cursor of focus and menus.
var Hand = MustSprite(`
	.eee......
	elllleeee.
	elllllllle
	elllleeee.
	elllle....
	.eeee.....
`)

// Finish dithers out a disabled control and marks its focus.
func (g Canvas) Finish(st State) {
	if st.Disabled {
		g.Checker(0, 0, g.W, g.H, Face)
	}
	if st.Focus {
		g.FocusMark(0, 0, g.W, g.H)
	}
}

// FocusMark rings a box with dots and points a hand.
func (g Canvas) FocusMark(x, y, w, h int) {
	rx, ry, rw, rh := x-2, y-2, w+4, h+4
	for i := range rw {
		g.ant(rx+i, ry)
		g.ant(rx+i, ry+rh-1)
	}
	for j := 1; j < rh-1; j++ {
		g.ant(rx, ry+j)
		g.ant(rx+rw-1, ry+j)
	}
	g.Sprite(rx-Hand.W-1+g.Step(4)%2, y+(h-Hand.H)/2, Hand, 0, Edge, 1)
}

func (g Canvas) ant(x, y int) {
	if (g.gx+x+g.gy+y)&1 == 0 {
		g.Dot(x, y, Accent)
	}
}
